package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/DePavelPo/websocket-chat-server/internal/models"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ID       string
	Nickname string
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
}

func (c *Client) ReadProc() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			logrus.WithError(err).Warn("read message error")
			break
		}

		logrus.WithFields(logrus.Fields{
			"client":  c.Nickname,
			"message": string(msg),
		}).Info("received message")

		// Try to parse as JSON message
		var chatMsg models.ChatMessage
		if err := jsoniter.Unmarshal(msg, &chatMsg); err != nil {
			logrus.WithError(err).Info("message is not JSON, treating as plain text")
			// If not JSON, treat as plain text message
			msgStr := string(msg)
			if strings.HasPrefix(msgStr, "/") {
				c.HandleCommand(msgStr)
				continue
			}
			// Create JSON message from plain text
			chatMsg = models.ChatMessage{
				Type:      "message",
				Content:   msgStr,
				Username:  c.Nickname,
				Timestamp: time.Now(),
			}
		}

		// Check if the content is a command (starts with /)
		if strings.HasPrefix(chatMsg.Content, "/") {
			logrus.WithFields(logrus.Fields{
				"client":  c.Nickname,
				"command": chatMsg.Content,
			}).Info("handling command from JSON message")
			c.HandleCommand(chatMsg.Content)
			continue
		}

		// Handle different message types
		switch chatMsg.Type {
		case "message":
			logrus.WithFields(logrus.Fields{
				"client":  c.Nickname,
				"content": chatMsg.Content,
			}).Info("handling chat message")
			c.HandleChatMessage(chatMsg)
		case "command":
			logrus.WithFields(logrus.Fields{
				"client":  c.Nickname,
				"command": chatMsg.Content,
			}).Info("handling command")
			c.HandleCommand(chatMsg.Content)
		default:
			logrus.WithFields(logrus.Fields{
				"client": c.Nickname,
				"type":   chatMsg.Type,
			}).Warn("unknown message type")
			c.SendErrorMessage("Unknown message type")
		}
	}
}

func (c *Client) HandleChatMessage(msg models.ChatMessage) {
	// Set username from client if not provided
	if msg.Username == "" {
		msg.Username = c.Nickname
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	logrus.WithFields(logrus.Fields{
		"username": msg.Username,
		"content":  msg.Content,
	}).Info("broadcasting chat message")

	// Broadcast the message to all clients
	msgBytes, err := jsoniter.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal chat message")
		return
	}
	c.Hub.Broadcast <- msgBytes
}

func (c *Client) WriteProc() {
	defer func() {
		logrus.Info("closing websocket write connection")
		c.Conn.Close()
	}()

	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			logrus.WithError(err).Warn("write message error")
			break
		}
		logrus.WithField("message", string(msg)).Info("sent message to client")
	}
}

func (c *Client) HandleCommand(msg string) {
	parts := strings.Fields(msg)

	logrus.WithFields(logrus.Fields{
		"client":  c.Nickname,
		"command": parts[0],
		"args":    parts[1:],
	}).Info("handling command")

	switch parts[0] {
	case "/id":
		response := models.CommandMessage{
			Type:      "command_response",
			Command:   "id",
			Content:   fmt.Sprintf("Your ID: %s", c.ID),
			Timestamp: time.Now(),
		}
		c.SendJSONMessage(response)
	case "/nick":
		if len(parts) < 2 {
			response := models.CommandMessage{
				Type:      "command_response",
				Command:   "nick",
				Content:   "Example: /nick your_new_nickname",
				Timestamp: time.Now(),
			}
			c.SendJSONMessage(response)
			return
		}

		oldNickname := c.Nickname
		c.Nickname = parts[1]

		// Send success response to the user
		response := models.CommandMessage{
			Type:      "command_response",
			Command:   "nick",
			Content:   fmt.Sprintf("Your nickname changed to: %s", c.Nickname),
			Timestamp: time.Now(),
		}
		c.SendJSONMessage(response)

		// Broadcast nickname change to all users
		systemMsg := models.SystemMessage{
			Type:      "system",
			Content:   fmt.Sprintf("!!! %s is now known as: %s", oldNickname, c.Nickname),
			Timestamp: time.Now(),
		}
		c.BroadcastJSONMessage(systemMsg)
	case "/who":
		var clientList []string
		for client := range c.Hub.Clients {
			clientList = append(clientList, client.Nickname)
		}

		response := models.UserListMessage{
			Type:      "user_list",
			Users:     clientList,
			Timestamp: time.Now(),
		}
		c.SendJSONMessage(response)
	case "/help":
		response := models.CommandMessage{
			Type:      "command_response",
			Command:   "help",
			Content:   "Available commands:\n/nick <name>\n/id\n/who\n/help",
			Timestamp: time.Now(),
		}
		c.SendJSONMessage(response)
	default:
		response := models.CommandMessage{
			Type:      "command_response",
			Command:   "unknown",
			Content:   "Unknown command. Try /help",
			Timestamp: time.Now(),
		}
		c.SendJSONMessage(response)
	}
}

func (c *Client) SendJSONMessage(msg interface{}) {
	msgBytes, err := jsoniter.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal JSON message")
		return
	}
	logrus.WithFields(logrus.Fields{
		"client":  c.Nickname,
		"message": string(msgBytes),
	}).Info("sending JSON message to client")
	c.Send <- msgBytes
}

func (c *Client) BroadcastJSONMessage(msg interface{}) {
	msgBytes, err := jsoniter.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal broadcast message")
		return
	}
	logrus.WithFields(logrus.Fields{
		"client":  c.Nickname,
		"message": string(msgBytes),
	}).Info("broadcasting JSON message")
	c.Hub.Broadcast <- msgBytes
}

func (c *Client) SendErrorMessage(message string) {
	response := models.CommandMessage{
		Type:      "error",
		Command:   "error",
		Content:   message,
		Timestamp: time.Now(),
	}
	c.SendJSONMessage(response)
}

func (c *Client) StartPing(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := c.Conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(3*time.Second)); err != nil {
			logrus.WithError(err).Warn("ws ping failed")
			return
		}
	}
}

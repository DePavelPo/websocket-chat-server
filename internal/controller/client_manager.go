package controller

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
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

		msgStr := string(msg)
		if strings.HasPrefix(msgStr, "/") {
			c.HandleCommand(msgStr)
			continue
		}
		formattedMsg := fmt.Sprintf("[%s]: %s", c.Nickname, msg)
		c.Hub.Broadcast <- []byte(formattedMsg)
	}
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
	}
}

func (c *Client) HandleCommand(msg string) {
	parts := strings.Fields(msg)

	switch parts[0] {
	case "/id":
		c.Send <- []byte(fmt.Sprintf("Your ID: %s", c.ID))
	case "/nick":
		if len(parts) < 2 {
			c.Send <- []byte("Example: /nick your_new_nickname")
			return
		}

		oldNickname := c.Nickname
		c.Nickname = parts[1]
		c.Send <- []byte(fmt.Sprintf("Your nickname changed to: %s", c.Nickname))
		c.Hub.Broadcast <- []byte(fmt.Sprintf("!!! %s is now known as: %s", oldNickname, c.Nickname))
	case "/who":
		var clientList []string
		for client := range c.Hub.Clients {
			clientList = append(clientList, client.Nickname)
		}
		c.Send <- []byte(fmt.Sprintf("Users online:\n%s", strings.Join(clientList, "\n")))
	case "/help":
		c.Send <- []byte("Available commands:\n/nick <name>\n/id\n/who\n/help")
	default:
		c.Send <- []byte("Unknown command. Try /help")
	}
}

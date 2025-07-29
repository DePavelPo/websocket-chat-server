package models

import "time"

// ChatMessage represents a chat message sent via WebSocket
type ChatMessage struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

// SystemMessage represents a system message (user joined, left, etc.)
type SystemMessage struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// UserListMessage represents a list of online users
type UserListMessage struct {
	Type      string    `json:"type"`
	Users     []string  `json:"users"`
	Timestamp time.Time `json:"timestamp"`
}

// CommandMessage represents a command message
type CommandMessage struct {
	Type      string    `json:"type"`
	Command   string    `json:"command"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

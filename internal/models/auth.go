package models

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type AuthResponse struct {
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

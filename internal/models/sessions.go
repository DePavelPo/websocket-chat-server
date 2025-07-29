package models

import "time"

type Session struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token_hash"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

package repository

import (
	"database/sql"
	"time"

	"github.com/DePavelPo/websocket-chat-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, token_hash, expires_at) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	return r.db.QueryRow(query, session.UserID, session.Token, session.ExpiresAt).Scan(&session.ID)
}

func (r *sessionRepository) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at 
		FROM sessions 
		WHERE token_hash = $1 AND expires_at > $2`

	now := time.Now()
	err := r.db.Get(&session, query, token, now)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) DeleteSession(token string) error {
	query := `DELETE FROM sessions WHERE token_hash = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *sessionRepository) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at < $1`
	now := time.Now()
	_, err := r.db.Exec(query, now)
	return err
}

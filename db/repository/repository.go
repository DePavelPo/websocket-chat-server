package repository

import (
	"github.com/DePavelPo/websocket-chat-server/internal/models"

	"github.com/jmoiron/sqlx"
)

// UserRepository interface for user operations
type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

// SessionRepository interface for session operations
type SessionRepository interface {
	CreateSession(session *models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	DeleteSession(token string) error
	DeleteExpiredSessions() error
}

// Repository struct that implements both interfaces
type Repository struct {
	users    UserRepository
	sessions SessionRepository
}

// NewRepository creates a new repository with separate user and session repositories
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		users:    NewUserRepository(db),
		sessions: NewSessionRepository(db),
	}
}

// UserRepository methods
func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	return r.users.GetUserByUsername(username)
}

func (r *Repository) CreateUser(user *models.User) error {
	return r.users.CreateUser(user)
}

func (r *Repository) UpdateUser(user *models.User) error {
	return r.users.UpdateUser(user)
}

func (r *Repository) DeleteUser(id int) error {
	return r.users.DeleteUser(id)
}

// SessionRepository methods
func (r *Repository) CreateSession(session *models.Session) error {
	return r.sessions.CreateSession(session)
}

func (r *Repository) GetSessionByToken(token string) (*models.Session, error) {
	return r.sessions.GetSessionByToken(token)
}

func (r *Repository) DeleteSession(token string) error {
	return r.sessions.DeleteSession(token)
}

func (r *Repository) DeleteExpiredSessions() error {
	return r.sessions.DeleteExpiredSessions()
}

// GetUserRepository returns the user repository
func (r *Repository) GetUserRepository() UserRepository {
	return r.users
}

// GetSessionRepository returns the session repository
func (r *Repository) GetSessionRepository() SessionRepository {
	return r.sessions
}

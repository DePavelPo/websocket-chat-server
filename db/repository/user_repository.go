package repository

import (
	"database/sql"

	"github.com/DePavelPo/websocket-chat-server/internal/models"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = $1`
	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, password_hash) 
		VALUES ($1, $2) 
		RETURNING id`

	return r.db.QueryRow(query, user.Username, user.Password).Scan(&user.ID)
}

func (r *userRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET username = $1, password_hash = $2 
		WHERE id = $3`

	_, err := r.db.Exec(query, user.Username, user.Password, user.ID)
	return err
}

func (r *userRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

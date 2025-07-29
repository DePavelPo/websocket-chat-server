package service

import (
	"errors"
	"time"

	"github.com/DePavelPo/websocket-chat-server/db/repository"
	"github.com/DePavelPo/websocket-chat-server/internal/auth"
	"github.com/DePavelPo/websocket-chat-server/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	authClient  auth.Auth
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	authClient auth.Auth,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		authClient:  authClient,
	}
}

// Register creates a new user
func (s *AuthService) Register(username, password string) error {
	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return s.userRepo.CreateUser(user)
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (string, error) {
	// Get user from database
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := s.authClient.GenerateJWT(username)
	if err != nil {
		return "", err
	}

	// Create session in database
	session := &models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: timePtr(time.Now().Add(6 * time.Hour)),
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken validates a JWT token and returns the username
func (s *AuthService) ValidateToken(token string) (string, error) {
	// First check if session exists in database
	session, err := s.sessionRepo.GetSessionByToken(token)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", errors.New("session not found or expired")
	}

	// Validate JWT token
	username, err := s.authClient.ValidateJWT(token)
	if err != nil {
		return "", err
	}

	return username, nil
}

// Logout removes a session
func (s *AuthService) Logout(token string) error {
	return s.sessionRepo.DeleteSession(token)
}

// CleanupExpiredSessions removes expired sessions
func (s *AuthService) CleanupExpiredSessions() error {
	return s.sessionRepo.DeleteExpiredSessions()
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"jualan-online/services/auth-service/internal/domain"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, username, password, role string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type authService struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewAuthService(repo domain.UserRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: jwtSecret}
}

func (s *authService) Register(ctx context.Context, username, password, role string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"pvs/internal/domain"
	"pvs/internal/repository"
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.UserRepository, secret []byte) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret}
}

func (s *AuthService) Register(ctx context.Context, email, password, role string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email/password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &domain.User{
		Email:    email,
		Password: string(hash),
		Role:     role,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", err
	}
	return GenerateToken(s.jwtSecret, role)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	return GenerateToken(s.jwtSecret, user.Role)
}

func GenerateToken(secret []byte, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_type": role,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

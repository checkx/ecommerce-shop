package service

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"ecommerce-shop/internal/auth"
)

type AuthService struct {
	DB        *sqlx.DB
	Log       *zap.Logger
	JWTSecret string
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	var id string
	if err := s.DB.GetContext(ctx, &id, `INSERT INTO users(email, password_hash) VALUES ($1,$2) RETURNING id`, email, string(hash)); err != nil {
		return "", "", err
	}
	tok, err := auth.GenerateToken(id, s.JWTSecret, 24*time.Hour)
	return id, tok, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	var id, pw string
	if err := s.DB.QueryRowxContext(ctx, `SELECT id, password_hash FROM users WHERE email=$1`, email).Scan(&id, &pw); err != nil {
		return "", "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(pw), []byte(password)) != nil {
		return "", "", bcrypt.ErrMismatchedHashAndPassword
	}
	tok, err := auth.GenerateToken(id, s.JWTSecret, 24*time.Hour)
	return id, tok, err
}

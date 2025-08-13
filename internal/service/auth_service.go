package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"market/internal/domain"
	"market/internal/repository"
	"market/internal/utils"
)

type AuthService struct {
	users     repository.UserRepository
	jwtSecret string
	ttl       time.Duration
}

type RegisterInput struct {
	Email    string
	Password string
	Role     domain.Role
}

type AuthResult struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
	Role   string `json:"role"`
}

func NewAuthService(users repository.UserRepository, jwtSecret string, ttl time.Duration) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret, ttl: ttl}
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*AuthResult, error) {
	if in.Email == "" || in.Password == "" {
		return nil, errors.New("email and password required")
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if in.Role != domain.RoleBuyer && in.Role != domain.RoleSeller {
		return nil, errors.New("invalid role")
	}
	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	id, err := s.users.Create(ctx, in.Email, hash, in.Role)
	if err != nil {
		return nil, err
	}
	token, err := utils.CreateJWT(id, string(in.Role), s.jwtSecret, s.ttl)
	if err != nil {
		return nil, err
	}
	return &AuthResult{UserID: id, Token: token, Role: string(in.Role)}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := utils.CheckPassword(u.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	token, err := utils.CreateJWT(u.ID, string(u.Role), s.jwtSecret, s.ttl)
	if err != nil {
		return nil, err
	}
	return &AuthResult{UserID: u.ID, Token: token, Role: string(u.Role)}, nil
}

package service

import (
	"context"
	"strings"

	"budget-family/internal/config"
	"budget-family/internal/entity"
	"budget-family/internal/repository"
	"budget-family/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, name, email, phone, password string) (*entity.User, utils.TokenPair, error)
	Login(ctx context.Context, email, password string) (*entity.User, utils.TokenPair, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*entity.User, error)
}

type authService struct {
	cfg        config.Config
	logger     *zap.Logger
	users      repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(cfg config.Config, logger *zap.Logger, users repository.UserRepository, jwtManager *utils.JWTManager) AuthService {
	return &authService{cfg: cfg, logger: logger, users: users, jwtManager: jwtManager}
}

func (s *authService) Register(ctx context.Context, name, email, phone, password string) (*entity.User, utils.TokenPair, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if _, err := s.users.GetByEmail(ctx, email); err == nil {
		return nil, utils.TokenPair{}, utils.NewConflict("email already registered", nil)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, utils.TokenPair{}, utils.NewInternal("failed to check email", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cfg.Auth.PasswordCost)
	if err != nil {
		return nil, utils.TokenPair{}, utils.NewInternal("failed to hash password", err)
	}

	u := &entity.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hash),
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, utils.TokenPair{}, utils.NewInternal("failed to create user", err)
	}

	tokens, err := s.jwtManager.GenerateTokenPair(u.ID, s.cfg.Auth.Issuer)
	if err != nil {
		return nil, utils.TokenPair{}, utils.NewInternal("failed to generate tokens", err)
	}

	return u, tokens, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*entity.User, utils.TokenPair, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.TokenPair{}, utils.NewUnauthorized("invalid email or password", nil)
		}
		return nil, utils.TokenPair{}, utils.NewInternal("failed to fetch user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, utils.TokenPair{}, utils.NewUnauthorized("invalid email or password", nil)
	}

	tokens, err := s.jwtManager.GenerateTokenPair(u.ID, s.cfg.Auth.Issuer)
	if err != nil {
		return nil, utils.TokenPair{}, utils.NewInternal("failed to generate tokens", err)
	}

	return u, tokens, nil
}

func (s *authService) GetMe(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFound("user not found", nil)
		}
		return nil, utils.NewInternal("failed to get user", err)
	}
	return u, nil
}

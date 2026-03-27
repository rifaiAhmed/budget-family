package service

import (
	"context"
	"strings"

	"budget-family/internal/config"
	"budget-family/internal/entity"
	"budget-family/internal/repository"
	"budget-family/pkg/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WalletService interface {
	Create(ctx context.Context, familyID uuid.UUID, name, typ string, balance decimal.Decimal) (*entity.Wallet, error)
	List(ctx context.Context, familyID uuid.UUID, page, limit int) ([]entity.Wallet, utils.PageMeta, error)
	Update(ctx context.Context, id uuid.UUID, name, typ string) (*entity.Wallet, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type walletService struct {
	cfg    config.Config
	logger *zap.Logger
	repo   repository.WalletRepository
}

func NewWalletService(cfg config.Config, logger *zap.Logger, repo repository.WalletRepository) WalletService {
	return &walletService{cfg: cfg, logger: logger, repo: repo}
}

func (s *walletService) Create(ctx context.Context, familyID uuid.UUID, name, typ string, balance decimal.Decimal) (*entity.Wallet, error) {
	w := &entity.Wallet{ID: uuid.New(), FamilyID: familyID, Name: strings.TrimSpace(name), Type: strings.TrimSpace(typ), Balance: balance}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, utils.NewInternal("failed to create wallet", err)
	}
	return w, nil
}

func (s *walletService) List(ctx context.Context, familyID uuid.UUID, page, limit int) ([]entity.Wallet, utils.PageMeta, error) {
	page, limit, offset := utils.NormalizePagination(page, limit)
	wallets, total, err := s.repo.ListByFamily(ctx, familyID, limit, offset)
	if err != nil {
		return nil, utils.PageMeta{}, utils.NewInternal("failed to list wallets", err)
	}
	return wallets, utils.BuildMeta(page, limit, total), nil
}

func (s *walletService) Update(ctx context.Context, id uuid.UUID, name, typ string) (*entity.Wallet, error) {
	w, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFound("wallet not found", nil)
		}
		return nil, utils.NewInternal("failed to get wallet", err)
	}
	w.Name = strings.TrimSpace(name)
	w.Type = strings.TrimSpace(typ)
	if err := s.repo.Update(ctx, w); err != nil {
		return nil, utils.NewInternal("failed to update wallet", err)
	}
	return w, nil
}

func (s *walletService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return utils.NewInternal("failed to delete wallet", err)
	}
	return nil
}

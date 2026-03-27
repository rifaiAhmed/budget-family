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
)

type BillService interface {
	Create(ctx context.Context, familyID uuid.UUID, name string, amount decimal.Decimal, dueDay int, recurring bool) (*entity.Bill, error)
	List(ctx context.Context, familyID uuid.UUID, page, limit int) ([]entity.Bill, utils.PageMeta, error)
}

type billService struct {
	cfg    config.Config
	logger *zap.Logger
	repo   repository.BillRepository
}

func NewBillService(cfg config.Config, logger *zap.Logger, repo repository.BillRepository) BillService {
	return &billService{cfg: cfg, logger: logger, repo: repo}
}

func (s *billService) Create(ctx context.Context, familyID uuid.UUID, name string, amount decimal.Decimal, dueDay int, recurring bool) (*entity.Bill, error) {
	if strings.TrimSpace(name) == "" {
		return nil, utils.NewBadRequest("name is required", nil)
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, utils.NewBadRequest("amount must be greater than 0", nil)
	}
	if dueDay < 1 || dueDay > 31 {
		return nil, utils.NewBadRequest("due_day must be 1-31", nil)
	}

	b := &entity.Bill{ID: uuid.New(), FamilyID: familyID, Name: strings.TrimSpace(name), Amount: amount, DueDay: dueDay, Recurring: recurring}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, utils.NewInternal("failed to create bill", err)
	}
	return b, nil
}

func (s *billService) List(ctx context.Context, familyID uuid.UUID, page, limit int) ([]entity.Bill, utils.PageMeta, error) {
	page, limit, offset := utils.NormalizePagination(page, limit)
	items, total, err := s.repo.ListByFamily(ctx, familyID, limit, offset)
	if err != nil {
		return nil, utils.PageMeta{}, utils.NewInternal("failed to list bills", err)
	}
	return items, utils.BuildMeta(page, limit, total), nil
}

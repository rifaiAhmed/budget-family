package service

import (
	"context"

	"budget-family/internal/config"
	"budget-family/internal/entity"
	"budget-family/internal/repository"
	"budget-family/pkg/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GoalService interface {
	Create(ctx context.Context, familyID uuid.UUID, name string, targetAmount decimal.Decimal, targetDate string) (*entity.Goal, error)
	List(ctx context.Context, familyID uuid.UUID, page, limit int) ([]entity.Goal, utils.PageMeta, error)
	Update(ctx context.Context, id uuid.UUID, name string, targetAmount, currentAmount decimal.Decimal, targetDate string) (*entity.Goal, error)
}

type goalService struct {
	cfg    config.Config
	logger *zap.Logger
	repo   repository.GoalRepository
}

func NewGoalService(cfg config.Config, logger *zap.Logger, repo repository.GoalRepository) GoalService {
	return &goalService{cfg: cfg, logger: logger, repo: repo}
}

func (s *goalService) Create(ctx context.Context, familyID uuid.UUID, name string, targetAmount decimal.Decimal, targetDate string) (*entity.Goal, error) {
	if targetAmount.LessThanOrEqual(decimal.Zero) {
		return nil, utils.NewBadRequest("target_amount must be greater than 0", nil)
	}
	g := &entity.Goal{ID: uuid.New(), FamilyID: familyID, Name: name, TargetAmount: targetAmount, CurrentAmount: decimal.Zero}
	if targetDate != "" {
		t, err := utils.ParseDate(targetDate)
		if err != nil {
			return nil, utils.NewBadRequest("invalid target_date", err)
		}
		g.TargetDate = t
	}
	if err := s.repo.Create(ctx, g); err != nil {
		return nil, utils.NewInternal("failed to create goal", err)
	}
	return g, nil
}

func (s *goalService) List(ctx context.Context, familyID uuid.UUID, page, limit int) ([]entity.Goal, utils.PageMeta, error) {
	page, limit, offset := utils.NormalizePagination(page, limit)
	items, total, err := s.repo.ListByFamily(ctx, familyID, limit, offset)
	if err != nil {
		return nil, utils.PageMeta{}, utils.NewInternal("failed to list goals", err)
	}
	return items, utils.BuildMeta(page, limit, total), nil
}

func (s *goalService) Update(ctx context.Context, id uuid.UUID, name string, targetAmount, currentAmount decimal.Decimal, targetDate string) (*entity.Goal, error) {
	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFound("goal not found", nil)
		}
		return nil, utils.NewInternal("failed to get goal", err)
	}
	if targetAmount.LessThanOrEqual(decimal.Zero) {
		return nil, utils.NewBadRequest("target_amount must be greater than 0", nil)
	}
	if currentAmount.LessThan(decimal.Zero) {
		return nil, utils.NewBadRequest("current_amount must be >= 0", nil)
	}

	g.Name = name
	g.TargetAmount = targetAmount
	g.CurrentAmount = currentAmount
	if targetDate != "" {
		t, err := utils.ParseDate(targetDate)
		if err != nil {
			return nil, utils.NewBadRequest("invalid target_date", err)
		}
		g.TargetDate = t
	}

	if err := s.repo.Update(ctx, g); err != nil {
		return nil, utils.NewInternal("failed to update goal", err)
	}
	return g, nil
}

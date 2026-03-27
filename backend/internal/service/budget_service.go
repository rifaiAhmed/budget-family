package service

import (
	"context"
	"time"

	"budget-family/internal/config"
	"budget-family/internal/entity"
	"budget-family/internal/repository"
	"budget-family/pkg/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type BudgetUsageRow struct {
	FamilyID        uuid.UUID       `json:"family_id"`
	CategoryID      uuid.UUID       `json:"category_id"`
	Month           int             `json:"month"`
	Year            int             `json:"year"`
	BudgetAmount    decimal.Decimal `json:"budget_amount"`
	UsedAmount      decimal.Decimal `json:"used_amount"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	PercentageUsed  decimal.Decimal `json:"percentage_used"`
}

type BudgetService interface {
	Upsert(ctx context.Context, familyID, categoryID uuid.UUID, amount decimal.Decimal, month, year int) (*entity.Budget, error)
	List(ctx context.Context, familyID uuid.UUID, month, year, page, limit int) ([]entity.Budget, utils.PageMeta, error)
	Usage(ctx context.Context, familyID uuid.UUID, month, year int) ([]BudgetUsageRow, error)
	Remaining(ctx context.Context, familyID, categoryID uuid.UUID, month, year int, now time.Time) (decimal.Decimal, decimal.Decimal, error)
}

type budgetService struct {
	cfg     config.Config
	logger  *zap.Logger
	budgets repository.BudgetRepository
	txRepo  repository.TransactionRepository
}

func NewBudgetService(cfg config.Config, logger *zap.Logger, budgets repository.BudgetRepository, txRepo repository.TransactionRepository) BudgetService {
	return &budgetService{cfg: cfg, logger: logger, budgets: budgets, txRepo: txRepo}
}

func (s *budgetService) Upsert(ctx context.Context, familyID, categoryID uuid.UUID, amount decimal.Decimal, month, year int) (*entity.Budget, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, utils.NewBadRequest("amount must be greater than 0", nil)
	}
	if month < 1 || month > 12 {
		return nil, utils.NewBadRequest("month must be 1-12", nil)
	}
	if year < 2000 {
		return nil, utils.NewBadRequest("invalid year", nil)
	}

	b := &entity.Budget{ID: uuid.New(), FamilyID: familyID, CategoryID: categoryID, Amount: amount, Month: month, Year: year}
	if err := s.budgets.Upsert(ctx, b); err != nil {
		return nil, utils.NewInternal("failed to upsert budget", err)
	}
	return b, nil
}

func (s *budgetService) List(ctx context.Context, familyID uuid.UUID, month, year, page, limit int) ([]entity.Budget, utils.PageMeta, error) {
	page, limit, offset := utils.NormalizePagination(page, limit)
	items, total, err := s.budgets.ListByFamily(ctx, familyID, month, year, limit, offset)
	if err != nil {
		return nil, utils.PageMeta{}, utils.NewInternal("failed to list budgets", err)
	}
	return items, utils.BuildMeta(page, limit, total), nil
}

func (s *budgetService) Usage(ctx context.Context, familyID uuid.UUID, month, year int) ([]BudgetUsageRow, error) {
	if month < 1 || month > 12 {
		return nil, utils.NewBadRequest("month must be 1-12", nil)
	}
	if year < 2000 {
		return nil, utils.NewBadRequest("invalid year", nil)
	}

	budgets, _, err := s.budgets.ListByFamily(ctx, familyID, month, year, 1000, 0)
	if err != nil {
		return nil, utils.NewInternal("failed to list budgets", err)
	}

	rows := make([]BudgetUsageRow, 0, len(budgets))
	for _, b := range budgets {
		used, err := s.txRepo.SumByCategoryInMonth(ctx, familyID, b.CategoryID, month, year)
		if err != nil {
			return nil, utils.NewInternal("failed to calculate budget usage", err)
		}

		remaining := b.Amount.Sub(used)
		pct := decimal.Zero
		if b.Amount.GreaterThan(decimal.Zero) {
			pct = used.Div(b.Amount).Mul(decimal.NewFromInt(100))
		}
		rows = append(rows, BudgetUsageRow{
			FamilyID:        familyID,
			CategoryID:      b.CategoryID,
			Month:           month,
			Year:            year,
			BudgetAmount:    b.Amount,
			UsedAmount:      used,
			RemainingAmount: remaining,
			PercentageUsed:  pct.Round(2),
		})
	}

	return rows, nil
}

func (s *budgetService) Remaining(ctx context.Context, familyID, categoryID uuid.UUID, month, year int, now time.Time) (decimal.Decimal, decimal.Decimal, error) {
	budgets, _, err := s.budgets.ListByFamily(ctx, familyID, month, year, 1000, 0)
	if err != nil {
		return decimal.Zero, decimal.Zero, utils.NewInternal("failed to list budgets", err)
	}

	var amount decimal.Decimal
	for _, b := range budgets {
		if b.CategoryID == categoryID {
			amount = b.Amount
			break
		}
	}
	if amount.Equal(decimal.Zero) {
		return decimal.Zero, decimal.Zero, nil
	}

	used, err := s.txRepo.SumByCategoryInMonth(ctx, familyID, categoryID, month, year)
	if err != nil {
		return decimal.Zero, decimal.Zero, utils.NewInternal("failed to calculate budget usage", err)
	}
	return amount, amount.Sub(used), nil
}

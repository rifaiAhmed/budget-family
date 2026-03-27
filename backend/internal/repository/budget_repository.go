package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BudgetRepository interface {
	Upsert(ctx context.Context, b *entity.Budget) error
	ListByFamily(ctx context.Context, familyID uuid.UUID, month, year int, limit, offset int) ([]entity.Budget, int64, error)
	GetByUnique(ctx context.Context, familyID, categoryID uuid.UUID, month, year int) (*entity.Budget, error)
}

type budgetRepository struct{ db *gorm.DB }

func NewBudgetRepository(db *gorm.DB) BudgetRepository { return &budgetRepository{db: db} }

func (r *budgetRepository) Upsert(ctx context.Context, b *entity.Budget) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "family_id"}, {Name: "category_id"}, {Name: "month"}, {Name: "year"}},
		DoUpdates: clause.AssignmentColumns([]string{"amount"}),
	}).Create(b).Error
}

func (r *budgetRepository) ListByFamily(ctx context.Context, familyID uuid.UUID, month, year int, limit, offset int) ([]entity.Budget, int64, error) {
	var budgets []entity.Budget
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Budget{}).Where("family_id = ?", familyID)
	if month > 0 {
		q = q.Where("month = ?", month)
	}
	if year > 0 {
		q = q.Where("year = ?", year)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("year DESC, month DESC").Limit(limit).Offset(offset).Find(&budgets).Error; err != nil {
		return nil, 0, err
	}
	return budgets, total, nil
}

func (r *budgetRepository) GetByUnique(ctx context.Context, familyID, categoryID uuid.UUID, month, year int) (*entity.Budget, error) {
	var b entity.Budget
	if err := r.db.WithContext(ctx).First(&b, "family_id = ? AND category_id = ? AND month = ? AND year = ?", familyID, categoryID, month, year).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

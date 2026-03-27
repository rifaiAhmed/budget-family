package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoalRepository interface {
	Create(ctx context.Context, g *entity.Goal) error
	ListByFamily(ctx context.Context, familyID uuid.UUID, limit, offset int) ([]entity.Goal, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Goal, error)
	Update(ctx context.Context, g *entity.Goal) error
}

type goalRepository struct{ db *gorm.DB }

func NewGoalRepository(db *gorm.DB) GoalRepository { return &goalRepository{db: db} }

func (r *goalRepository) Create(ctx context.Context, g *entity.Goal) error { return r.db.WithContext(ctx).Create(g).Error }

func (r *goalRepository) ListByFamily(ctx context.Context, familyID uuid.UUID, limit, offset int) ([]entity.Goal, int64, error) {
	var goals []entity.Goal
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Goal{}).Where("family_id = ?", familyID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("target_date ASC").Limit(limit).Offset(offset).Find(&goals).Error; err != nil {
		return nil, 0, err
	}
	return goals, total, nil
}

func (r *goalRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Goal, error) {
	var g entity.Goal
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *goalRepository) Update(ctx context.Context, g *entity.Goal) error { return r.db.WithContext(ctx).Save(g).Error }

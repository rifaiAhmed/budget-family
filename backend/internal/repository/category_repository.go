package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *entity.Category) error
	ListByFamily(ctx context.Context, familyID uuid.UUID, typ string, limit, offset int) ([]entity.Category, int64, error)
}

type categoryRepository struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) CategoryRepository { return &categoryRepository{db: db} }

func (r *categoryRepository) Create(ctx context.Context, c *entity.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *categoryRepository) ListByFamily(ctx context.Context, familyID uuid.UUID, typ string, limit, offset int) ([]entity.Category, int64, error) {
	var categories []entity.Category
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Category{}).Where("family_id = ?", familyID)
	if typ != "" {
		q = q.Where("type = ?", typ)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("name ASC").Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

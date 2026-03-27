package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillRepository interface {
	Create(ctx context.Context, b *entity.Bill) error
	ListByFamily(ctx context.Context, familyID uuid.UUID, limit, offset int) ([]entity.Bill, int64, error)
}

type billRepository struct{ db *gorm.DB }

func NewBillRepository(db *gorm.DB) BillRepository { return &billRepository{db: db} }

func (r *billRepository) Create(ctx context.Context, b *entity.Bill) error { return r.db.WithContext(ctx).Create(b).Error }

func (r *billRepository) ListByFamily(ctx context.Context, familyID uuid.UUID, limit, offset int) ([]entity.Bill, int64, error) {
	var bills []entity.Bill
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Bill{}).Where("family_id = ?", familyID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("due_day ASC").Limit(limit).Offset(offset).Find(&bills).Error; err != nil {
		return nil, 0, err
	}
	return bills, total, nil
}

package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FamilyRepository interface {
	CreateFamily(ctx context.Context, family *entity.Family, ownerMember *entity.FamilyMember) error
	ListFamiliesByUser(ctx context.Context, userID uuid.UUID) ([]entity.Family, error)
	IsMember(ctx context.Context, familyID, userID uuid.UUID) (bool, error)
}

type familyRepository struct {
	db *gorm.DB
}

func NewFamilyRepository(db *gorm.DB) FamilyRepository {
	return &familyRepository{db: db}
}

func (r *familyRepository) CreateFamily(ctx context.Context, family *entity.Family, ownerMember *entity.FamilyMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(family).Error; err != nil {
			return err
		}
		if err := tx.Create(ownerMember).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *familyRepository) ListFamiliesByUser(ctx context.Context, userID uuid.UUID) ([]entity.Family, error) {
	var families []entity.Family
	if err := r.db.WithContext(ctx).
		Table("families").
		Select("families.*").
		Joins("JOIN family_members fm ON fm.family_id = families.id").
		Where("fm.user_id = ?", userID).
		Order("families.created_at DESC").
		Scan(&families).Error; err != nil {
		return nil, err
	}
	return families, nil
}

func (r *familyRepository) IsMember(ctx context.Context, familyID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.FamilyMember{}).
		Where("family_id = ? AND user_id = ?", familyID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FamilyRepository interface {
	CreateFamily(ctx context.Context, family *entity.Family, ownerMember *entity.FamilyMember) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Family, error)
	ListFamiliesByUser(ctx context.Context, userID uuid.UUID) ([]entity.Family, error)
	IsMember(ctx context.Context, familyID, userID uuid.UUID) (bool, error)
	AddMember(ctx context.Context, m *entity.FamilyMember) error
	ListMembers(ctx context.Context, familyID uuid.UUID) ([]FamilyMemberRow, error)
}

type FamilyMemberRow struct {
	ID       uuid.UUID `json:"id"`
	FamilyID uuid.UUID `json:"family_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
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

func (r *familyRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Family, error) {
	var f entity.Family
	if err := r.db.WithContext(ctx).First(&f, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &f, nil
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

func (r *familyRepository) AddMember(ctx context.Context, m *entity.FamilyMember) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *familyRepository) ListMembers(ctx context.Context, familyID uuid.UUID) ([]FamilyMemberRow, error) {
	var rows []FamilyMemberRow
	if err := r.db.WithContext(ctx).
		Table("family_members fm").
		Select("fm.id, fm.family_id, fm.user_id, fm.role, u.name, u.email").
		Joins("JOIN users u ON u.id = fm.user_id").
		Where("fm.family_id = ?", familyID).
		Order("CASE WHEN fm.role = 'owner' THEN 0 ELSE 1 END, u.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

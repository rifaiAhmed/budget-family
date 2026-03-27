package repository

import (
	"context"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	Create(ctx context.Context, w *entity.Wallet) error
	ListByFamily(ctx context.Context, familyID uuid.UUID, limit, offset int) ([]entity.Wallet, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Wallet, error)
	Update(ctx context.Context, w *entity.Wallet) error
	Delete(ctx context.Context, id uuid.UUID) error
	LockByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entity.Wallet, error)
	UpdateBalance(ctx context.Context, tx *gorm.DB, id uuid.UUID, newBalance string) error
}

type walletRepository struct{ db *gorm.DB }

func NewWalletRepository(db *gorm.DB) WalletRepository { return &walletRepository{db: db} }

func (r *walletRepository) Create(ctx context.Context, w *entity.Wallet) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *walletRepository) ListByFamily(ctx context.Context, familyID uuid.UUID, limit, offset int) ([]entity.Wallet, int64, error) {
	var wallets []entity.Wallet
	var total int64
	q := r.db.WithContext(ctx).Model(&entity.Wallet{}).Where("family_id = ?", familyID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&wallets).Error; err != nil {
		return nil, 0, err
	}
	return wallets, total, nil
}

func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Wallet, error) {
	var w entity.Wallet
	if err := r.db.WithContext(ctx).First(&w, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *walletRepository) Update(ctx context.Context, w *entity.Wallet) error {
	return r.db.WithContext(ctx).Save(w).Error
}

func (r *walletRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Wallet{}, "id = ?", id).Error
}

func (r *walletRepository) LockByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entity.Wallet, error) {
	var w entity.Wallet
	if err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&w, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *walletRepository) UpdateBalance(ctx context.Context, tx *gorm.DB, id uuid.UUID, newBalance string) error {
	return tx.WithContext(ctx).Model(&entity.Wallet{}).Where("id = ?", id).Update("balance", newBalance).Error
}

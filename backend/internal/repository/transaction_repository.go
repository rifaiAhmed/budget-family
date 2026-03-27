package repository

import (
	"context"
	"time"

	"budget-family/internal/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TransactionFilters struct {
	FamilyID   uuid.UUID
	WalletID   *uuid.UUID
	CategoryID *uuid.UUID
	Type       string
	FromDate   *time.Time
	ToDate     *time.Time
	MinAmount  *decimal.Decimal
	MaxAmount  *decimal.Decimal
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, t *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	List(ctx context.Context, f TransactionFilters, limit, offset int) ([]entity.Transaction, int64, error)
	Update(ctx context.Context, tx *gorm.DB, t *entity.Transaction) error
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	SumByCategoryInMonth(ctx context.Context, familyID, categoryID uuid.UUID, month, year int) (decimal.Decimal, error)
	Summary(ctx context.Context, familyID uuid.UUID, from, to *time.Time) (decimal.Decimal, decimal.Decimal, error)
}

type transactionRepository struct{ db *gorm.DB }

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, t *entity.Transaction) error {
	return tx.WithContext(ctx).Create(t).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	var t entity.Transaction
	if err := r.db.WithContext(ctx).
		Table("transactions t").
		Select("t.*, u.name as created_by_name, u.email as created_by_email").
		Joins("JOIN users u ON u.id = t.created_by").
		Where("t.id = ?", id).
		Scan(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transactionRepository) List(ctx context.Context, f TransactionFilters, limit, offset int) ([]entity.Transaction, int64, error) {
	var txs []entity.Transaction
	var total int64

	q := r.db.WithContext(ctx).Model(&entity.Transaction{}).Where("family_id = ?", f.FamilyID)
	if f.WalletID != nil {
		q = q.Where("wallet_id = ?", *f.WalletID)
	}
	if f.CategoryID != nil {
		q = q.Where("category_id = ?", *f.CategoryID)
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.FromDate != nil {
		q = q.Where("transaction_date >= ?", f.FromDate.Format("2006-01-02"))
	}
	if f.ToDate != nil {
		q = q.Where("transaction_date <= ?", f.ToDate.Format("2006-01-02"))
	}
	if f.MinAmount != nil {
		q = q.Where("amount >= ?", f.MinAmount.String())
	}
	if f.MaxAmount != nil {
		q = q.Where("amount <= ?", f.MaxAmount.String())
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("transaction_date DESC, created_at DESC").Limit(limit).Offset(offset).Find(&txs).Error; err != nil {
		return nil, 0, err
	}
	return txs, total, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *gorm.DB, t *entity.Transaction) error {
	return tx.WithContext(ctx).Save(t).Error
}

func (r *transactionRepository) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	return tx.WithContext(ctx).Delete(&entity.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) SumByCategoryInMonth(ctx context.Context, familyID, categoryID uuid.UUID, month, year int) (decimal.Decimal, error) {
	var sumStr *string
	err := r.db.WithContext(ctx).Table("transactions").
		Select("COALESCE(SUM(amount),0)::text").
		Where("family_id = ? AND category_id = ? AND type = 'expense' AND EXTRACT(MONTH FROM transaction_date) = ? AND EXTRACT(YEAR FROM transaction_date) = ?", familyID, categoryID, month, year).
		Scan(&sumStr).Error
	if err != nil {
		return decimal.Zero, err
	}
	if sumStr == nil {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(*sumStr)
}

func (r *transactionRepository) Summary(ctx context.Context, familyID uuid.UUID, from, to *time.Time) (decimal.Decimal, decimal.Decimal, error) {
	base := r.db.WithContext(ctx).Table("transactions").Where("family_id = ?", familyID)
	if from != nil {
		base = base.Where("transaction_date >= ?", from.Format("2006-01-02"))
	}
	if to != nil {
		base = base.Where("transaction_date <= ?", to.Format("2006-01-02"))
	}

	type row struct {
		Typ string
		Sum string
	}
	var rows []row
	if err := base.Select("type as typ, COALESCE(SUM(amount),0)::text as sum").Group("type").Scan(&rows).Error; err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	income := decimal.Zero
	expense := decimal.Zero
	for _, r := range rows {
		d, _ := decimal.NewFromString(r.Sum)
		if r.Typ == "income" {
			income = d
		}
		if r.Typ == "expense" {
			expense = d
		}
	}
	return income, expense, nil
}

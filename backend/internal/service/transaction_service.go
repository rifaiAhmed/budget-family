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
	"gorm.io/gorm"
)

type TransactionCreateInput struct {
	FamilyID        uuid.UUID
	WalletID        uuid.UUID
	CategoryID      uuid.UUID
	Amount          decimal.Decimal
	Type            string
	Note            string
	TransactionDate time.Time
	CreatedBy       uuid.UUID
}

type TransactionUpdateInput struct {
	Amount          decimal.Decimal
	Type            string
	Note            string
	TransactionDate time.Time
	CategoryID      uuid.UUID
}

type TransactionService interface {
	Create(ctx context.Context, in TransactionCreateInput) (*entity.Transaction, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	List(ctx context.Context, f repository.TransactionFilters, page, limit int) ([]entity.Transaction, utils.PageMeta, error)
	Update(ctx context.Context, id uuid.UUID, in TransactionUpdateInput) (*entity.Transaction, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Summary(ctx context.Context, familyID uuid.UUID, from, to *time.Time) (decimal.Decimal, decimal.Decimal, error)
}

type transactionService struct {
	cfg        config.Config
	logger     *zap.Logger
	txRepo     repository.TransactionRepository
	walletRepo repository.WalletRepository
	budgetRepo repository.BudgetRepository
	db         *gorm.DB
}

func NewTransactionService(cfg config.Config, logger *zap.Logger, db *gorm.DB, txRepo repository.TransactionRepository, walletRepo repository.WalletRepository, budgetRepo repository.BudgetRepository) TransactionService {
	return &transactionService{cfg: cfg, logger: logger, db: db, txRepo: txRepo, walletRepo: walletRepo, budgetRepo: budgetRepo}
}

func (s *transactionService) Create(ctx context.Context, in TransactionCreateInput) (*entity.Transaction, error) {
	if in.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, utils.NewBadRequest("amount must be greater than 0", nil)
	}

	t := &entity.Transaction{
		ID:              uuid.New(),
		FamilyID:        in.FamilyID,
		WalletID:        in.WalletID,
		CategoryID:      in.CategoryID,
		Amount:          in.Amount,
		Type:            in.Type,
		Note:            in.Note,
		TransactionDate: in.TransactionDate,
		CreatedBy:       in.CreatedBy,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		w, err := s.walletRepo.LockByID(ctx, tx, in.WalletID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.NewNotFound("wallet not found", nil)
			}
			return utils.NewInternal("failed to lock wallet", err)
		}

		newBal := w.Balance
		if in.Type == "income" {
			newBal = newBal.Add(in.Amount)
		} else if in.Type == "expense" {
			newBal = newBal.Sub(in.Amount)
		} else {
			return utils.NewBadRequest("type must be income or expense", nil)
		}

		if err := s.txRepo.Create(ctx, tx, t); err != nil {
			return utils.NewInternal("failed to create transaction", err)
		}

		if err := s.walletRepo.UpdateBalance(ctx, tx, w.ID, newBal.StringFixed(2)); err != nil {
			return utils.NewInternal("failed to update wallet balance", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *transactionService) Get(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	t, err := s.txRepo.GetByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewNotFound("transaction not found", nil)
		}
		return nil, utils.NewInternal("failed to get transaction", err)
	}
	return t, nil
}

func (s *transactionService) List(ctx context.Context, f repository.TransactionFilters, page, limit int) ([]entity.Transaction, utils.PageMeta, error) {
	page, limit, offset := utils.NormalizePagination(page, limit)
	txs, total, err := s.txRepo.List(ctx, f, limit, offset)
	if err != nil {
		return nil, utils.PageMeta{}, utils.NewInternal("failed to list transactions", err)
	}
	return txs, utils.BuildMeta(page, limit, total), nil
}

func (s *transactionService) Update(ctx context.Context, id uuid.UUID, in TransactionUpdateInput) (*entity.Transaction, error) {
	if in.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, utils.NewBadRequest("amount must be greater than 0", nil)
	}

	var updated *entity.Transaction
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cur, err := s.txRepo.GetByID(ctx, id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.NewNotFound("transaction not found", nil)
			}
			return utils.NewInternal("failed to get transaction", err)
		}

		w, err := s.walletRepo.LockByID(ctx, tx, cur.WalletID)
		if err != nil {
			return utils.NewInternal("failed to lock wallet", err)
		}

		// reverse current effect
		bal := w.Balance
		if cur.Type == "income" {
			bal = bal.Sub(cur.Amount)
		} else if cur.Type == "expense" {
			bal = bal.Add(cur.Amount)
		}

		// apply new effect
		if in.Type == "income" {
			bal = bal.Add(in.Amount)
		} else if in.Type == "expense" {
			bal = bal.Sub(in.Amount)
		} else {
			return utils.NewBadRequest("type must be income or expense", nil)
		}

		cur.Amount = in.Amount
		cur.Type = in.Type
		cur.Note = in.Note
		cur.TransactionDate = in.TransactionDate
		cur.CategoryID = in.CategoryID

		if err := s.txRepo.Update(ctx, tx, cur); err != nil {
			return utils.NewInternal("failed to update transaction", err)
		}
		if err := s.walletRepo.UpdateBalance(ctx, tx, w.ID, bal.StringFixed(2)); err != nil {
			return utils.NewInternal("failed to update wallet balance", err)
		}
		updated = cur
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *transactionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cur, err := s.txRepo.GetByID(ctx, id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.NewNotFound("transaction not found", nil)
			}
			return utils.NewInternal("failed to get transaction", err)
		}

		w, err := s.walletRepo.LockByID(ctx, tx, cur.WalletID)
		if err != nil {
			return utils.NewInternal("failed to lock wallet", err)
		}

		bal := w.Balance
		if cur.Type == "income" {
			bal = bal.Sub(cur.Amount)
		} else if cur.Type == "expense" {
			bal = bal.Add(cur.Amount)
		}

		if err := s.txRepo.Delete(ctx, tx, id); err != nil {
			return utils.NewInternal("failed to delete transaction", err)
		}
		if err := s.walletRepo.UpdateBalance(ctx, tx, w.ID, bal.StringFixed(2)); err != nil {
			return utils.NewInternal("failed to update wallet balance", err)
		}
		return nil
	})
}

func (s *transactionService) Summary(ctx context.Context, familyID uuid.UUID, from, to *time.Time) (decimal.Decimal, decimal.Decimal, error) {
	income, expense, err := s.txRepo.Summary(ctx, familyID, from, to)
	if err != nil {
		return decimal.Zero, decimal.Zero, utils.NewInternal("failed to get summary", err)
	}
	return income, expense, nil
}

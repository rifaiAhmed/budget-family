package handler

import (
	"net/http"
	"time"

	"budget-family/internal/repository"
	"budget-family/internal/service"
	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type TransactionHandler struct {
	baseHandler
	svc service.TransactionService
}

func NewTransactionHandler(logger *zap.Logger, v *validator.Validate, svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type transactionCreateRequest struct {
	FamilyID        uuid.UUID `json:"family_id" validate:"required"`
	WalletID        uuid.UUID `json:"wallet_id" validate:"required"`
	CategoryID      uuid.UUID `json:"category_id" validate:"required"`
	Amount          string    `json:"amount" validate:"required"`
	Type            string    `json:"type" validate:"required,oneof=income expense"`
	Note            string    `json:"note" validate:"omitempty"`
	TransactionDate string    `json:"transaction_date" validate:"required"` // YYYY-MM-DD
}

type transactionUpdateRequest struct {
	CategoryID      uuid.UUID `json:"category_id" validate:"required"`
	Amount          string    `json:"amount" validate:"required"`
	Type            string    `json:"type" validate:"required,oneof=income expense"`
	Note            string    `json:"note" validate:"omitempty"`
	TransactionDate string    `json:"transaction_date" validate:"required"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	uid, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req transactionCreateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid amount")
		return
	}
	dt, err := utils.ParseDate(req.TransactionDate)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	tx, err := h.svc.Create(c.Request.Context(), service.TransactionCreateInput{
		FamilyID:        req.FamilyID,
		WalletID:        req.WalletID,
		CategoryID:      req.CategoryID,
		Amount:          amount,
		Type:            req.Type,
		Note:            req.Note,
		TransactionDate: dt,
		CreatedBy:       uid,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "success", gin.H{"transaction": tx})
}

func (h *TransactionHandler) List(c *gin.Context) {
	familyIDStr := c.Query("family_id")
	familyID, err := uuid.Parse(familyIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "family_id is required")
		return
	}

	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 20)

	var walletID *uuid.UUID
	if s := c.Query("wallet_id"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			utils.Fail(c, http.StatusBadRequest, "invalid wallet_id")
			return
		}
		walletID = &id
	}
	var categoryID *uuid.UUID
	if s := c.Query("category_id"); s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			utils.Fail(c, http.StatusBadRequest, "invalid category_id")
			return
		}
		categoryID = &id
	}

	from, ok := parseDateQuery(c, "from")
	if !ok {
		return
	}
	to, ok := parseDateQuery(c, "to")
	if !ok {
		return
	}

	filters := repository.TransactionFilters{
		FamilyID:   familyID,
		WalletID:   walletID,
		CategoryID: categoryID,
		Type:       c.Query("type"),
		FromDate:   from,
		ToDate:     to,
	}

	items, meta, err := h.svc.List(c.Request.Context(), filters, page, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"items": items, "meta": meta})
}

func (h *TransactionHandler) Get(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	tx, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"transaction": tx})
}

func (h *TransactionHandler) Update(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	var req transactionUpdateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid amount")
		return
	}
	dt, err := utils.ParseDate(req.TransactionDate)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid transaction_date")
		return
	}

	updated, err := h.svc.Update(c.Request.Context(), id, service.TransactionUpdateInput{
		Amount:          amount,
		Type:            req.Type,
		Note:            req.Note,
		TransactionDate: dt,
		CategoryID:      req.CategoryID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"transaction": updated})
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{})
}

func (h *TransactionHandler) Summary(c *gin.Context) {
	familyIDStr := c.Query("family_id")
	familyID, err := uuid.Parse(familyIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "family_id is required")
		return
	}

	from, ok := parseDateQuery(c, "from")
	if !ok {
		return
	}
	to, ok := parseDateQuery(c, "to")
	if !ok {
		return
	}

	income, expense, err := h.svc.Summary(c.Request.Context(), familyID, from, to)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"income": income, "expense": expense, "net": income.Sub(expense)})
}

var _ = time.Now

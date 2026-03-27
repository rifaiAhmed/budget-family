package handler

import (
	"net/http"
	"time"

	"budget-family/internal/service"
	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type BudgetHandler struct {
	baseHandler
	svc service.BudgetService
}

func NewBudgetHandler(logger *zap.Logger, v *validator.Validate, svc service.BudgetService) *BudgetHandler {
	return &BudgetHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type budgetCreateRequest struct {
	FamilyID   uuid.UUID `json:"family_id" validate:"required"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	Amount     string    `json:"amount" validate:"required"`
	Month      int       `json:"month" validate:"required"`
	Year       int       `json:"year" validate:"required"`
}

func (h *BudgetHandler) Create(c *gin.Context) {
	var req budgetCreateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid amount")
		return
	}

	b, err := h.svc.Upsert(c.Request.Context(), req.FamilyID, req.CategoryID, amount, req.Month, req.Year)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "success", gin.H{"budget": b})
}

func (h *BudgetHandler) List(c *gin.Context) {
	familyIDStr := c.Query("family_id")
	familyID, err := uuid.Parse(familyIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "family_id is required")
		return
	}

	month := parseIntQuery(c, "month", 0)
	year := parseIntQuery(c, "year", 0)
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 20)

	items, meta, err := h.svc.List(c.Request.Context(), familyID, month, year, page, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"items": items, "meta": meta})
}

func (h *BudgetHandler) Usage(c *gin.Context) {
	familyIDStr := c.Query("family_id")
	familyID, err := uuid.Parse(familyIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "family_id is required")
		return
	}

	month := parseIntQuery(c, "month", int(time.Now().Month()))
	year := parseIntQuery(c, "year", time.Now().Year())

	rows, err := h.svc.Usage(c.Request.Context(), familyID, month, year)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"items": rows, "month": month, "year": year})
}

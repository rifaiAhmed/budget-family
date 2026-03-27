package handler

import (
	"net/http"

	"budget-family/internal/service"
	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type BillHandler struct {
	baseHandler
	svc service.BillService
}

func NewBillHandler(logger *zap.Logger, v *validator.Validate, svc service.BillService) *BillHandler {
	return &BillHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type billCreateRequest struct {
	FamilyID  uuid.UUID `json:"family_id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=2,max=200"`
	Amount    string    `json:"amount" validate:"required"`
	DueDay    int       `json:"due_day" validate:"required"`
	Recurring bool      `json:"recurring"`
}

func (h *BillHandler) Create(c *gin.Context) {
	var req billCreateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid amount")
		return
	}

	b, err := h.svc.Create(c.Request.Context(), req.FamilyID, req.Name, amount, req.DueDay, req.Recurring)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "success", gin.H{"bill": b})
}

func (h *BillHandler) List(c *gin.Context) {
	familyIDStr := c.Query("family_id")
	familyID, err := uuid.Parse(familyIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "family_id is required")
		return
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 20)

	items, meta, err := h.svc.List(c.Request.Context(), familyID, page, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"items": items, "meta": meta})
}

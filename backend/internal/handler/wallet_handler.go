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

type WalletHandler struct {
	baseHandler
	svc service.WalletService
}

func NewWalletHandler(logger *zap.Logger, v *validator.Validate, svc service.WalletService) *WalletHandler {
	return &WalletHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type walletCreateRequest struct {
	FamilyID uuid.UUID `json:"family_id" validate:"required"`
	Name     string    `json:"name" validate:"required,min=2,max=200"`
	Type     string    `json:"type" validate:"required,oneof=cash bank ewallet card"`
	Balance  string    `json:"balance" validate:"omitempty"`
}

type walletUpdateRequest struct {
	Name string `json:"name" validate:"required,min=2,max=200"`
	Type string `json:"type" validate:"required,oneof=cash bank ewallet card"`
}

func (h *WalletHandler) List(c *gin.Context) {
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

func (h *WalletHandler) Create(c *gin.Context) {
	var req walletCreateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}
	bal := decimal.Zero
	if req.Balance != "" {
		d, err := decimal.NewFromString(req.Balance)
		if err != nil {
			utils.Fail(c, http.StatusBadRequest, "invalid balance")
			return
		}
		bal = d
	}

	w, err := h.svc.Create(c.Request.Context(), req.FamilyID, req.Name, req.Type, bal)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "success", gin.H{"wallet": w})
}

func (h *WalletHandler) Update(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	var req walletUpdateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}

	w, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Type)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"wallet": w})
}

func (h *WalletHandler) Delete(c *gin.Context) {
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

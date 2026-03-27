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

type GoalHandler struct {
	baseHandler
	svc service.GoalService
}

func NewGoalHandler(logger *zap.Logger, v *validator.Validate, svc service.GoalService) *GoalHandler {
	return &GoalHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type goalCreateRequest struct {
	FamilyID     uuid.UUID `json:"family_id" validate:"required"`
	Name         string    `json:"name" validate:"required,min=2,max=200"`
	TargetAmount string    `json:"target_amount" validate:"required"`
	TargetDate   string    `json:"target_date" validate:"omitempty"`
}

type goalUpdateRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=200"`
	TargetAmount  string `json:"target_amount" validate:"required"`
	CurrentAmount string `json:"current_amount" validate:"required"`
	TargetDate    string `json:"target_date" validate:"omitempty"`
}

func (h *GoalHandler) Create(c *gin.Context) {
	var req goalCreateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}
	ta, err := decimal.NewFromString(req.TargetAmount)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid target_amount")
		return
	}

	g, err := h.svc.Create(c.Request.Context(), req.FamilyID, req.Name, ta, req.TargetDate)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "success", gin.H{"goal": g})
}

func (h *GoalHandler) List(c *gin.Context) {
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

func (h *GoalHandler) Update(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	var req goalUpdateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}
	ta, err := decimal.NewFromString(req.TargetAmount)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid target_amount")
		return
	}
	ca, err := decimal.NewFromString(req.CurrentAmount)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid current_amount")
		return
	}

	g, err := h.svc.Update(c.Request.Context(), id, req.Name, ta, ca, req.TargetDate)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"goal": g})
}

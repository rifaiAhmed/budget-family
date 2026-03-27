package handler

import (
	"net/http"

	"budget-family/internal/service"
	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	baseHandler
	svc service.CategoryService
}

func NewCategoryHandler(logger *zap.Logger, v *validator.Validate, svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type categoryCreateRequest struct {
	FamilyID uuid.UUID `json:"family_id" validate:"required"`
	Name     string    `json:"name" validate:"required,min=2,max=200"`
	Type     string    `json:"type" validate:"required,oneof=income expense"`
	Icon     string    `json:"icon" validate:"omitempty,max=100"`
}

func (h *CategoryHandler) List(c *gin.Context) {
	familyIDStr := c.Query("family_id")
	familyID, err := uuid.Parse(familyIDStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "family_id is required")
		return
	}

	typ := c.Query("type")
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 50)

	items, meta, err := h.svc.List(c.Request.Context(), familyID, typ, page, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, "success", gin.H{"items": items, "meta": meta})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req categoryCreateRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}

	cat, err := h.svc.Create(c.Request.Context(), req.FamilyID, req.Name, req.Type, req.Icon)
	if err != nil {
		h.handleError(c, err)
		return
	}
	utils.Success(c, http.StatusCreated, "success", gin.H{"category": cat})
}

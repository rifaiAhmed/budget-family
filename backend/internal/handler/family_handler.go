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

type FamilyHandler struct {
	baseHandler
	svc service.FamilyService
}

func NewFamilyHandler(logger *zap.Logger, v *validator.Validate, svc service.FamilyService) *FamilyHandler {
	return &FamilyHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type createFamilyRequest struct {
	Name string `json:"name" validate:"required,min=2,max=200"`
}

type inviteFamilyRequest struct {
	FamilyID uuid.UUID `json:"family_id" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
}

func (h *FamilyHandler) Create(c *gin.Context) {
	uid, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createFamilyRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}

	family, err := h.svc.Create(c.Request.Context(), uid, req.Name)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "success", gin.H{"family": family})
}

func (h *FamilyHandler) List(c *gin.Context) {
	uid, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	families, err := h.svc.List(c.Request.Context(), uid)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "success", gin.H{"items": families})
}

func (h *FamilyHandler) Invite(c *gin.Context) {
	uid, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req inviteFamilyRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}

	inv, err := h.svc.Invite(c.Request.Context(), uid, req.FamilyID, req.Email)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "success", gin.H{"invitation": inv})
}

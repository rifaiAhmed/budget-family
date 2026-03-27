package handler

import (
	"net/http"

	"budget-family/internal/service"
	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AuthHandler struct {
	baseHandler
	svc service.AuthService
}

func NewAuthHandler(logger *zap.Logger, v *validator.Validate, svc service.AuthService) *AuthHandler {
	return &AuthHandler{baseHandler: baseHandler{logger: logger, validator: v}, svc: svc}
}

type registerRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=200"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"omitempty,max=30"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}

	u, tokens, err := h.svc.Register(c.Request.Context(), req.Name, req.Email, req.Phone, req.Password)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "success", gin.H{"user": u, "tokens": tokens})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if !h.bindAndValidateJSON(c, &req) {
		return
	}

	u, tokens, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "success", gin.H{"user": u, "tokens": tokens})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid, ok := getUserID(c)
	if !ok {
		utils.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	u, err := h.svc.GetMe(c.Request.Context(), uid)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "success", gin.H{"user": u})
}

package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type baseHandler struct {
	logger    *zap.Logger
	validator *validator.Validate
}

func (h *baseHandler) bindAndValidateJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request body")
		return false
	}
	if err := h.validator.Struct(req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "validation error")
		return false
	}
	return true
}

func (h *baseHandler) handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if appErr, ok := utils.AsAppError(err); ok {
		utils.Fail(c, appErr.HTTPStatus, appErr.Message)
		return
	}
	h.logger.Error("unhandled error", zap.Error(err))
	utils.Fail(c, http.StatusInternalServerError, "internal server error")
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	idStr := c.Param(name)
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid id")
		return uuid.Nil, false
	}
	return id, true
}

func parseIntQuery(c *gin.Context, name string, def int) int {
	s := strings.TrimSpace(c.Query(name))
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func parseDateQuery(c *gin.Context, name string) (*time.Time, bool) {
	s := strings.TrimSpace(c.Query(name))
	if s == "" {
		return nil, true
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid date")
		return nil, false
	}
	return &t, true
}

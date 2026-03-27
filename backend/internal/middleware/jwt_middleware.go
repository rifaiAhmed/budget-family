package middleware

import (
	"net/http"
	"strings"

	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JWTMiddleware struct {
	jwt    *utils.JWTManager
	issuer string
}

func NewJWTMiddleware(jwtManager *utils.JWTManager, issuer string) *JWTMiddleware {
	return &JWTMiddleware{jwt: jwtManager, issuer: issuer}
}

func (m *JWTMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.Fail(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Fail(c, http.StatusUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := m.jwt.Verify(parts[1], m.issuer)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		uid, err := uuid.Parse(claims.UserID)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "invalid token subject")
			c.Abort()
			return
		}

		c.Set("user_id", uid)
		c.Next()
	}
}

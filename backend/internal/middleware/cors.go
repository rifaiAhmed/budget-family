package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowedOrigins []string
}

func CORS(cfg CORSConfig) gin.HandlerFunc {
	allowedAll := false
	allowed := map[string]struct{}{}
	for _, o := range cfg.AllowedOrigins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			allowedAll = true
			continue
		}
		allowed[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			if allowedAll {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				if _, ok := allowed[origin]; ok {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-Id")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateJwt(tokenString, cfg)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func ProductKeyMiddleware(keyRepo repositories.KeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		providedKey := c.GetHeader("X-API-Key")
		if providedKey == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_API_KEY", "Invalid or missing X-API-Key header")
			c.Abort()
			return
		}
		keyHash := utils.HashKey(providedKey)
		key, err := keyRepo.GetActiveKeyByHash(c.Request.Context(), keyHash)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_API_KEY", "Invalid or missing X-API-Key header")
			c.Abort()
			return
		}
		_ = keyRepo.UpdateLastUsedAt(c.Request.Context(), key.ID)
		c.Next()
	}
}

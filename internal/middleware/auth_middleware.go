package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*clientLimiter
	rate     rate.Limit // request per second
	burst    int
}

func NewRateLimiterStore(r rate.Limit, burst int) *RateLimiterStore {
	store := &RateLimiterStore{
		limiters: make(map[string]*clientLimiter),
		rate:     r,
		burst:    burst,
	}
	// Background cleanup: remove entries not seen in 10 minutes
	go store.cleanupLoop()
	return store
}

func (s *RateLimiterStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)

	for range ticker.C {
		s.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for key, entry := range s.limiters {
			if entry.lastSeen.Before(cutoff) {
				delete(s.limiters, key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *RateLimiterStore) get(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.limiters[key]

	if !exists {
		entry = &clientLimiter{
			limiter: rate.NewLimiter(s.rate, s.burst),
		}
		s.limiters[key] = entry
	}

	entry.lastSeen = time.Now()

	return entry.limiter
}

// Allows 20 req/s with a burst of 40 per IP
func IPRateLimiter(store *RateLimiterStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.GetHeader("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(c.Request.RemoteAddr)
		}

		if !store.get(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests from your IP",
				"retry_after": "1s",
			})
			return
		}
		c.Next()
	}
}

// Rate limit for API Key requests
func APIKeyRateLimiter(store *RateLimiterStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		key := fmt.Sprintf("user:%v", userID)
		if !store.get(key).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded for your API key",
				"retry_after": "1s",
			})
			return
		}
		c.Next()
	}
}

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

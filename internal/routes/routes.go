package routes

import (
	"net/http"

	"github.com/coolpythoncodes/nigerian-universities/internal/handlers"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(svc service.InstitutionService, db *gorm.DB, authHandler *handlers.AuthHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	{
		v1.GET("/health", func(c *gin.Context) {
			sqlDB, err := db.DB()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "db-unreachable",
					"error":  "failed to get database connection",
				})
				return
			}

			if err := sqlDB.Ping(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "db-unreachable",
					"error":  "failed to ping database",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// auth
		auth := v1.Group("/auth")

		auth.
			GET("/google", authHandler.GoogleLogin).
			GET("/google/callback", authHandler.GoogleCallback).
			POST("/google/login", authHandler.LoginWithGoogle). // when frontend is using Authjs library
			GET("/github", authHandler.GithubLogin).
			GET("/github/callback", authHandler.GihubCallback).
			POST("/github/login", authHandler.LoginWithGithub)

		// institution
		institution := v1.Group("/institutions")

		institution.GET("", handlers.GetAllInstitutions(svc))

		// api-keys
		keys := v1.Group("/api-keys")

		keys.
			POST("", handlers.CreateApiKey)
	}

	return r


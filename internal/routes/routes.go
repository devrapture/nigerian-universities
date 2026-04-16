package routes

import (
	"net/http"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/handlers"
	"github.com/coolpythoncodes/nigerian-universities/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HandlerDependencies struct {
	AuthHandler        *handlers.AuthHandler
	InstitutionHandler *handlers.InstitutionHandler
	KeyHandler         *handlers.KeyHandlers
}

func Setup(db *gorm.DB, cfg *config.Config, deps HandlerDependencies) *gin.Engine {
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
			GET("/google", deps.AuthHandler.GoogleLogin).
			GET("/google/callback", deps.AuthHandler.GoogleCallback).
			POST("/google/login", deps.AuthHandler.LoginWithGoogle). // when frontend is using Authjs library
			GET("/github", deps.AuthHandler.GithubLogin).
			GET("/github/callback", deps.AuthHandler.GithubCallback).
			POST("/github/login", deps.AuthHandler.LoginWithGithub)

		// institution
		institution := v1.Group("/institutions")

		institution.GET("", deps.InstitutionHandler.GetAllInstitutions)

		// api-keys
		keys := v1.Group("/api-keys")
		keys.Use(middleware.AuthMiddleware(cfg))
		keys.
			POST("/generate", deps.KeyHandler.CreateApiKey).
			GET("", deps.KeyHandler.GetAllKeys).
			POST("/:key_id/revoke", deps.KeyHandler.RevokeKey)

	}

	return r

}

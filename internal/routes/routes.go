package routes

import (
	"net/http"

	"github.com/coolpythoncodes/nigerian-universities/internal/handlers"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/gin-gonic/gin"
)

func Setup(svc service.InstitutionService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})

		// institution
		institution := v1.Group("/institutions")

		institution.GET("", handlers.GetAllInstitutions(svc))
	}

	return r
}

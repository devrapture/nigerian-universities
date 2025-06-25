package main

import (
	"github.com/coolpythoncodes/nigerian-universities/controllers"
	"github.com/coolpythoncodes/nigerian-universities/scraper"
	"github.com/gin-gonic/gin"
)

func main() {
	scraper.ScrapeUniversities()
	r := gin.Default()
	r.GET("/", controllers.GetAllUniversities)
	r.GET("/federal", controllers.GetAllFederalUniversities)
	r.GET("/state", controllers.GetAllStateUniversities)
	r.GET("/private", controllers.GetAllPrivateUniversities)

	// New endpoints
	r.GET("/university/details/:name", controllers.GetUniversityDetailsByNameOrAbbreviation)
	r.GET("/university/city/:city", controllers.GetUniversitiesByCity)
	r.GET("/university/state/:state", controllers.GetUniversitiesByState)
	r.GET("/university/private", controllers.GetAllPrivateUniversities) // alias with limit already applied
	r.GET("/university/private/:state", controllers.GetPrivateUniversitiesByState)

	r.Run(":8080")

}

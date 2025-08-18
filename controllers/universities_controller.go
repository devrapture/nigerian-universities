package controllers

import (
	"net/http"
	"strings"

	"github.com/coolpythoncodes/nigerian-universities/models"

	"github.com/coolpythoncodes/nigerian-universities/utils"
	"github.com/gin-gonic/gin"
)

func GetAllUniversities(c *gin.Context) {
	univeristies, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    univeristies,
	})
}

func GetAllFederalUniversities(c *gin.Context) {
	univeristies, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	federalUniversities := utils.Filter(univeristies, "federal")

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    federalUniversities,
	})
}

func GetAllStateUniversities(c *gin.Context) {
	univeristies, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	stateUniversities := utils.Filter(univeristies, "state")

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    stateUniversities,
	})
}

func GetAllPrivateUniversities(c *gin.Context) {
	univeristies, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	privateUniversities := utils.Filter(univeristies, "private")

	// Limit to first 20 universities
	if len(privateUniversities) > 20 {
		privateUniversities = privateUniversities[:20]
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    privateUniversities,
	})
}

func GetUniversityDetailsByNameOrAbbreviation(c *gin.Context) {
	nameParam := c.Param("name")
	if nameParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "name parameter is required",
		})
		return
	}

	universities, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	// normalise the param to lower for comparison
	target := strings.ToLower(nameParam)
	var result *models.Universities

	for _, uni := range universities {
		if strings.ToLower(uni.Name) == target {
			result = &uni
			break
		}
		if strings.ToLower(utils.GenerateAbbreviation(uni.Name)) == target {
			result = &uni
			break
		}
	}

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "university not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    result,
	})
}

func GetUniversitiesByCity(c *gin.Context) {
	cityParam := c.Param("city")
	if cityParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "city parameter is required",
		})
		return
	}

	universities, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	target := strings.ToLower(cityParam)
	filtered := make([]models.Universities, 0)
	for _, uni := range universities {
		if strings.Contains(strings.ToLower(uni.Name), target) {
			filtered = append(filtered, uni)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    filtered,
	})
}

func GetUniversitiesByState(c *gin.Context) {
	stateParam := c.Param("state")
	if stateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "state parameter is required",
		})
		return
	}

	universities, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	target := strings.ToLower(stateParam)
	filtered := make([]models.Universities, 0)
	for _, uni := range universities {
		if strings.Contains(strings.ToLower(uni.Name), target) {
			filtered = append(filtered, uni)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    filtered,
	})
}

func GetPrivateUniversitiesByState(c *gin.Context) {
	stateParam := c.Param("state")
	if stateParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "state parameter is required",
		})
		return
	}

	universities, err := utils.ReadUniversitiesFromJSONFile(utils.JsonFileName)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	target := strings.ToLower(stateParam)
	filtered := make([]models.Universities, 0)
	for _, uni := range universities {
		if strings.ToLower(uni.Type) == "private" && strings.Contains(strings.ToLower(uni.Name), target) {
			filtered = append(filtered, uni)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    filtered,
	})
}

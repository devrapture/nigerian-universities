package utils

import (
	"strings"

	"github.com/coolpythoncodes/nigerian-universities/models"
)

func Filter(universities []models.Institution, univerisityType string) []models.Institution {
	u := make([]models.Institution, 0)

	for _, uni := range universities {
		if strings.ToLower(uni.Type) == univerisityType {
			u = append(u, uni)
		}
	}

	return u
}

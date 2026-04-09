package utils

import (
	"encoding/json"
	"os"

	"github.com/coolpythoncodes/nigerian-universities/models"
)

func ReadUniversitiesFromJSONFile(fileName string) ([]models.Institution, *RestError) {
	univeristies := []models.Institution{}
	universitiesJson, err := os.ReadFile(fileName)

	if err != nil {
		return nil, NewInternalServerError("unable to get data")
	}
	if err := json.Unmarshal(universitiesJson, &univeristies); err != nil {
		return nil, NewInternalServerError("unable to get data unmarshale")
	}

	return univeristies, nil
}

package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coolpythoncodes/nigerian-universities/models"
)

func WriteJSON(fileName string, data []models.Institution) {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create JSON file")
		return
	}

	if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
		log.Fatalln("Failed to create" + fileName + "file")
	}
}

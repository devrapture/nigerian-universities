package main

import (
	"fmt"
	"log"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/database"
)

func main() {
	cfg := config.Load()

	_, err := database.ConnectDB(cfg)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Server starting on %s", addr)
}

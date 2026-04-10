package main

import (
	"fmt"
	"log"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/database"
	"github.com/coolpythoncodes/nigerian-universities/internal/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	_, err = database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Server starting on %s", addr)

	r := routes.Setup()

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

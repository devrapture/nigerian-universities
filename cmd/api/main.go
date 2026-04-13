package main

import (
	"fmt"
	"log"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/database"
	"github.com/coolpythoncodes/nigerian-universities/internal/handlers"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/routes"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repo := repositories.NewInstitutionRepository(db)
	svc := service.NewInstitutionService(repo)

	userRepo := repositories.NewUserRepository(db)
	userSvc := service.NewUserService(cfg, userRepo)

	authHandler := handlers.NewAuthHandler(userSvc)

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Server starting on %s", addr)

	r := routes.Setup(svc, db, authHandler)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

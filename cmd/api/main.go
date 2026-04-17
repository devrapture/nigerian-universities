package main

import (
	"fmt"
	"log"

	_ "github.com/coolpythoncodes/nigerian-universities/docs"
	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/database"
	"github.com/coolpythoncodes/nigerian-universities/internal/handlers"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/routes"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Nigerian Universities API
// @version 1.0
// @description This is an API for Nigerian Institutions

// @contact.name Rapture Chijioke Godson
// @contact.url https://github.com/devrapture
// @contact.email devrapture@proton.me

// @license.name MIT

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @tag.name Auth
// @tag.description Authentication endpoints
// @tag.name  Institutions
// @tag.description Institution listing endpoints
// @tag.name Keys
// @tag.description API key management endpoints

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Repositories
	institutionRepo := repositories.NewInstitutionRepository(db)
	userRepo := repositories.NewUserRepository(db)
	keyRepo := repositories.NewKeyRepository(db)

	// Services
	institutionService := service.NewInstitutionService(institutionRepo)
	userSvc := service.NewUserService(cfg, userRepo)
	keyService := service.NewKeyService(keyRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(userSvc)
	institutionHandler := handlers.NewInstitutionHandler(institutionService)
	keyHandler := handlers.NewKeyHandler(keyService)

	deps := routes.HandlerDependencies{
		AuthHandler:        authHandler,
		InstitutionHandler: institutionHandler,
		KeyHandler:         keyHandler,
		KeyRepo:            keyRepo,
	}

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Server starting on %s", addr)

	r := routes.Setup(db, cfg, deps)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

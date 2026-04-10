package main

import (
	"context"
	"fmt"
	"log"

	"github.com/coolpythoncodes/nigerian-universities/internal/config"
	"github.com/coolpythoncodes/nigerian-universities/internal/database"
	"github.com/coolpythoncodes/nigerian-universities/internal/repositories"
	"github.com/coolpythoncodes/nigerian-universities/internal/scraper"
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
	s := scraper.NewInstitutionScrapper()
	institutionService := service.NewInstitutionService(repo)
	institutions, err := s.ScrapeAllInstitution()
	fmt.Println("scraped institutions", institutions)
	if err != nil {
		log.Fatalf("Failed to scrape institutions: %v", err)
	}
	if err := institutionService.StoreScrapedInstitutions(context.Background(), institutions); err != nil {
		log.Fatalf("saving scraped data failed: %v", err)
	}

	log.Printf("scraping completed successfully: %d institutions processed", len(institutions))
}

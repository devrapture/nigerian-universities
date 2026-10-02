package handlers

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/gin-gonic/gin"
)

type InstitutionScraper interface {
	ScrapeAllInstitution() ([]model.Institution, error)
}

type ScrapeHandler struct {
	scraper            InstitutionScraper
	institutionService service.InstitutionService
	running            atomic.Bool
}

func NewScrapeHandler(scraper InstitutionScraper, institutionService service.InstitutionService) *ScrapeHandler {
	return &ScrapeHandler{scraper: scraper, institutionService: institutionService}
}

// Trigger starts a scrape in the background and returns without waiting for it.
// @Summary Trigger an institution scrape
// @Description Starts an institution scrape in the background. Overlapping triggers are accepted but do not start duplicate runs.
// @Tags Cron
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer CRON_SECRET"
// @Success 202 {object} schema.ScrapeTriggerResponse
// @Failure 401 {object} schema.InstitutionUnauthorizedResponse
// @Router /cron/scrape [post]
func (h *ScrapeHandler) Trigger(c *gin.Context) {
	if !h.running.CompareAndSwap(false, true) {
		utils.SuccessResponse(c, http.StatusAccepted, "scrape already running", gin.H{"status": "already_running"}, nil)
		return
	}

	go h.run()
	utils.SuccessResponse(c, http.StatusAccepted, "scrape started", gin.H{"status": "started"}, nil)
}

func (h *ScrapeHandler) run() {
	defer h.running.Store(false)

	institutions, err := h.scraper.ScrapeAllInstitution()
	if err != nil {
		log.Printf("background scrape failed: %v", err)
		return
	}

	if err := h.institutionService.StoreScrapedInstitutions(context.Background(), institutions); err != nil {
		log.Printf("saving background scrape failed: %v", err)
		return
	}

	log.Printf("background scrape completed successfully: %d institutions processed", len(institutions))
}

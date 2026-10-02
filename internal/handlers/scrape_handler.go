package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/coolpythoncodes/nigerian-universities/internal/service"
	"github.com/coolpythoncodes/nigerian-universities/internal/utils"
	"github.com/gin-gonic/gin"
)

// scrapeTimeout is slightly under cron-job.org's 25s client timeout so the
// handler can write a response before the caller disconnects.
const scrapeTimeout = 24 * time.Second

type InstitutionScraper interface {
	ScrapeAllInstitution(ctx context.Context) ([]model.Institution, error)
}

type ScrapeHandler struct {
	scraper            InstitutionScraper
	institutionService service.InstitutionService
	timeout            time.Duration
	running            atomic.Bool
}

func NewScrapeHandler(scraper InstitutionScraper, institutionService service.InstitutionService) *ScrapeHandler {
	return &ScrapeHandler{
		scraper:            scraper,
		institutionService: institutionService,
		timeout:            scrapeTimeout,
	}
}

// Trigger runs an institution scrape and waits for it to finish or time out.
// @Summary Trigger an institution scrape
// @Description Runs an institution scrape inside the request with a 24s timeout (cron-job.org's client timeout is 25s). Overlapping triggers are accepted but do not start duplicate runs.
// @Tags Cron
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer CRON_SECRET"
// @Success 200 {object} schema.ScrapeTriggerResponse
// @Success 202 {object} schema.ScrapeTriggerResponse
// @Failure 401 {object} schema.InstitutionUnauthorizedResponse
// @Failure 500 {object} schema.InstitutionInternalServerErrorResponse
// @Failure 504 {object} schema.ScrapeTimeoutResponse
// @Router /cron/scrape [post]
func (h *ScrapeHandler) Trigger(c *gin.Context) {
	if !h.running.CompareAndSwap(false, true) {
		utils.SuccessResponse(c, http.StatusAccepted, "scrape already running", gin.H{"status": "already_running"}, nil)
		return
	}
	defer h.running.Store(false)

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
	defer cancel()

	institutions, err := h.scraper.ScrapeAllInstitution(ctx)
	if err != nil {
		h.handleScrapeErr(c, "scrape failed", err)
		return
	}

	if err := h.institutionService.StoreScrapedInstitutions(ctx, institutions); err != nil {
		h.handleScrapeErr(c, "saving scrape failed", err)
		return
	}

	log.Printf("scrape completed successfully: %d institutions processed", len(institutions))
	utils.SuccessResponse(c, http.StatusOK, "scrape completed", gin.H{
		"status": "completed",
		"count":  len(institutions),
	}, nil)
}

func (h *ScrapeHandler) handleScrapeErr(c *gin.Context, logMsg string, err error) {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		log.Printf("%s: %v", logMsg, err)
		utils.ErrorResponse(c, http.StatusGatewayTimeout, "TIMEOUT", "scrape timed out")
		return
	}

	log.Printf("%s: %v", logMsg, err)
	utils.InternalError(c, logMsg)
}

package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coolpythoncodes/nigerian-universities/internal/dto"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"github.com/gin-gonic/gin"
)

type blockingScraper struct {
	started chan struct{}
	release chan struct{}
}

func (scraper *blockingScraper) ScrapeAllInstitution() ([]model.Institution, error) {
	close(scraper.started)
	<-scraper.release
	return []model.Institution{}, nil
}

type fakeInstitutionService struct {
	stored chan struct{}
}

func (service *fakeInstitutionService) StoreScrapedInstitutions(context.Context, []model.Institution) error {
	close(service.stored)
	return nil
}

func (*fakeInstitutionService) GetAllInstitutions(context.Context, dto.ListInstitutionQuery) ([]model.Institution, int64, error) {
	return nil, 0, nil
}

func (*fakeInstitutionService) GetChangeLogs(context.Context, dto.ListChangeLogQuery) ([]model.ChangeLog, int64, error) {
	return nil, 0, nil
}

func TestScrapeTriggerReturnsImmediatelyAndRejectsOverlap(t *testing.T) {
	gin.SetMode(gin.TestMode)
	scraper := &blockingScraper{started: make(chan struct{}), release: make(chan struct{})}
	service := &fakeInstitutionService{stored: make(chan struct{})}
	handler := NewScrapeHandler(scraper, service)
	router := gin.New()
	router.POST("/scrape", handler.Trigger)

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/scrape", nil))
	if first.Code != http.StatusAccepted {
		t.Fatalf("first trigger status = %d, want %d", first.Code, http.StatusAccepted)
	}
	if body := first.Body.String(); body != `{"success":true,"message":"scrape started","data":{"status":"started"}}` {
		t.Fatalf("first trigger body = %s", body)
	}

	select {
	case <-scraper.started:
	case <-time.After(time.Second):
		t.Fatal("background scraper did not start")
	}

	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/scrape", nil))
	if second.Code != http.StatusAccepted {
		t.Fatalf("overlapping trigger status = %d, want %d", second.Code, http.StatusAccepted)
	}
	if body := second.Body.String(); body != `{"success":true,"message":"scrape already running","data":{"status":"already_running"}}` {
		t.Fatalf("overlapping trigger body = %s", body)
	}

	close(scraper.release)
	select {
	case <-service.stored:
	case <-time.After(time.Second):
		t.Fatal("scraped institutions were not stored")
	}
}

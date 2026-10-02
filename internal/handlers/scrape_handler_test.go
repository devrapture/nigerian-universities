package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
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

func (scraper *blockingScraper) ScrapeAllInstitution(ctx context.Context) ([]model.Institution, error) {
	close(scraper.started)
	select {
	case <-scraper.release:
		return []model.Institution{{Name: "Test University"}}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type waitingScraper struct{}

func (waitingScraper) ScrapeAllInstitution(ctx context.Context) ([]model.Institution, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

type errorScraper struct {
	err error
}

func (s errorScraper) ScrapeAllInstitution(context.Context) ([]model.Institution, error) {
	return nil, s.err
}

type fastScraper struct {
	institutions []model.Institution
}

func (s fastScraper) ScrapeAllInstitution(context.Context) ([]model.Institution, error) {
	return s.institutions, nil
}

type fakeInstitutionService struct {
	stored     chan struct{}
	storeCalls atomic.Int32
	storeErr   error
}

func (service *fakeInstitutionService) StoreScrapedInstitutions(context.Context, []model.Institution) error {
	service.storeCalls.Add(1)
	if service.stored != nil {
		select {
		case <-service.stored:
		default:
			close(service.stored)
		}
	}
	return service.storeErr
}

func (*fakeInstitutionService) GetAllInstitutions(context.Context, dto.ListInstitutionQuery) ([]model.Institution, int64, error) {
	return nil, 0, nil
}

func (*fakeInstitutionService) GetChangeLogs(context.Context, dto.ListChangeLogQuery) ([]model.ChangeLog, int64, error) {
	return nil, 0, nil
}

func newTestRouter(handler *ScrapeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/scrape", handler.Trigger)
	return router
}

func TestScrapeTriggerCompletes(t *testing.T) {
	service := &fakeInstitutionService{}
	handler := NewScrapeHandler(fastScraper{institutions: []model.Institution{{Name: "A"}, {Name: "B"}}}, service)
	router := newTestRouter(handler)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/scrape", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	want := `{"success":true,"message":"scrape completed","data":{"count":2,"status":"completed"}}`
	if body := rec.Body.String(); body != want {
		t.Fatalf("body = %s, want %s", body, want)
	}
	if got := service.storeCalls.Load(); got != 1 {
		t.Fatalf("store calls = %d, want 1", got)
	}
}

func TestScrapeTriggerRejectsOverlap(t *testing.T) {
	scraper := &blockingScraper{started: make(chan struct{}), release: make(chan struct{})}
	service := &fakeInstitutionService{stored: make(chan struct{})}
	handler := NewScrapeHandler(scraper, service)
	router := newTestRouter(handler)

	firstDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/scrape", nil))
		firstDone <- rec
	}()

	select {
	case <-scraper.started:
	case <-time.After(time.Second):
		t.Fatal("scraper did not start")
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
	case first := <-firstDone:
		if first.Code != http.StatusOK {
			t.Fatalf("first trigger status = %d, want %d; body = %s", first.Code, http.StatusOK, first.Body.String())
		}
	case <-time.After(time.Second):
		t.Fatal("first trigger did not complete")
	}

	select {
	case <-service.stored:
	case <-time.After(time.Second):
		t.Fatal("scraped institutions were not stored")
	}
}

func TestScrapeTriggerTimesOut(t *testing.T) {
	service := &fakeInstitutionService{}
	handler := NewScrapeHandler(waitingScraper{}, service)
	handler.timeout = 50 * time.Millisecond
	router := newTestRouter(handler)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/scrape", nil))

	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusGatewayTimeout, rec.Body.String())
	}
	want := `{"success":false,"error":{"code":"TIMEOUT","message":"scrape timed out"}}`
	if body := rec.Body.String(); body != want {
		t.Fatalf("body = %s, want %s", body, want)
	}
	if got := service.storeCalls.Load(); got != 0 {
		t.Fatalf("store calls = %d, want 0", got)
	}
}

func TestScrapeTriggerErrorReleasesLock(t *testing.T) {
	service := &fakeInstitutionService{}
	handler := NewScrapeHandler(errorScraper{err: errors.New("source down")}, service)
	router := newTestRouter(handler)

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/scrape", nil))
	if first.Code != http.StatusInternalServerError {
		t.Fatalf("error status = %d, want %d; body = %s", first.Code, http.StatusInternalServerError, first.Body.String())
	}
	if got := service.storeCalls.Load(); got != 0 {
		t.Fatalf("store calls = %d, want 0", got)
	}

	handler.scraper = fastScraper{institutions: []model.Institution{{Name: "A"}}}
	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/scrape", nil))
	if second.Code != http.StatusOK {
		t.Fatalf("retry status = %d, want %d; body = %s", second.Code, http.StatusOK, second.Body.String())
	}
}

package scraper

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestScrapeAllInstitutionRespectsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewInstitutionScrapper().ScrapeAllInstitution(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}

func TestRemainingTimeoutUsesDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	got := remainingTimeout(ctx)
	if got > 2*time.Second || got < time.Second {
		t.Fatalf("remainingTimeout = %v, want between 1s and 2s", got)
	}
}

func TestRemainingTimeoutWithoutDeadline(t *testing.T) {
	if got := remainingTimeout(context.Background()); got != 60*time.Second {
		t.Fatalf("remainingTimeout = %v, want 60s", got)
	}
}

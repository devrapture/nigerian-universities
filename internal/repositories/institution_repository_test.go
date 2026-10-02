package repositories

import (
	"context"
	"testing"

	"github.com/coolpythoncodes/nigerian-universities/internal/constants"
	"github.com/coolpythoncodes/nigerian-universities/internal/dto"
	"github.com/coolpythoncodes/nigerian-universities/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpsertManyRecordsOnlyRealDataChanges(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:change-log-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.AutoMigrate(&model.Institution{}, &model.ChangeLog{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	repo := NewInstitutionRepository(db)
	ctx := context.Background()
	original := model.Institution{
		Name:                "University of Lagos",
		Type:                constants.FederalUniversity,
		ViceChancellor:      "Professor One",
		Website:             "https://unilag.edu.ng",
		YearOfEstablishment: "1962",
	}

	if err := repo.UpsertMany(ctx, []model.Institution{original}); err != nil {
		t.Fatalf("insert institution: %v", err)
	}
	assertChangeLogCount(t, repo, 1)

	if err := repo.UpsertMany(ctx, []model.Institution{original}); err != nil {
		t.Fatalf("repeat unchanged institution: %v", err)
	}
	assertChangeLogCount(t, repo, 1)

	updated := original
	updated.ViceChancellor = "Professor Two"
	updated.Website = "https://www.unilag.edu.ng"
	if err := repo.UpsertMany(ctx, []model.Institution{updated}); err != nil {
		t.Fatalf("update institution: %v", err)
	}

	logs, total, err := repo.FindChangeLogs(ctx, dto.ListChangeLogQuery{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list change logs: %v", err)
	}
	if total != 2 || len(logs) != 2 {
		t.Fatalf("got total=%d len=%d, want two change-log entries", total, len(logs))
	}
	if logs[0].ChangeType != "updated" {
		t.Fatalf("latest change type = %q, want updated", logs[0].ChangeType)
	}
	if len(logs[0].Changes) != 2 {
		t.Fatalf("updated field count = %d, want 2", len(logs[0].Changes))
	}
	if logs[0].ChangedAt.IsZero() {
		t.Fatal("changed_at must be populated")
	}
}

func assertChangeLogCount(t *testing.T, repo InstitutionRepository, want int64) {
	t.Helper()
	logs, total, err := repo.FindChangeLogs(context.Background(), dto.ListChangeLogQuery{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("list change logs: %v", err)
	}
	if total != want || int64(len(logs)) != want {
		t.Fatalf("got total=%d len=%d, want %d", total, len(logs), want)
	}
}

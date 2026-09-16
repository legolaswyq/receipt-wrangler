package repositories_test

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/repositories"
	"testing"
)

func TestDashboardPeriodRoundTrips(t *testing.T) {
	defer repositories.TruncateTestDb()
	repositories.CreateTestGroupWithUsers()

	repo := repositories.NewDashboardRepository(nil)

	created, err := repo.CreateDashboard(commands.UpsertDashboardCommand{
		Name:    "Home",
		GroupId: "1",
		Period:  "THIS_MONTH",
	}, 1)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.Period != "THIS_MONTH" {
		t.Fatalf("created period = %q, want THIS_MONTH", created.Period)
	}

	// Update to a custom range.
	updated, err := repo.UpdateDashboardById(created.ID, commands.UpsertDashboardCommand{
		Name:            "Home",
		GroupId:         "1",
		Period:          "CUSTOM",
		PeriodStartDate: "2026-01-01T00:00:00Z",
		PeriodEndDate:   "2026-12-31T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Period != "CUSTOM" || updated.PeriodStartDate != "2026-01-01T00:00:00Z" || updated.PeriodEndDate != "2026-12-31T00:00:00Z" {
		t.Fatalf("custom period not persisted: %+v", updated)
	}

	// Switch back to a preset — the custom dates must clear (Select forces the write).
	back, err := repo.UpdateDashboardById(created.ID, commands.UpsertDashboardCommand{
		Name:    "Home",
		GroupId: "1",
		Period:  "ALL_TIME",
	})
	if err != nil {
		t.Fatalf("update back: %v", err)
	}
	if back.Period != "ALL_TIME" || back.PeriodStartDate != "" || back.PeriodEndDate != "" {
		t.Fatalf("preset switch did not clear custom dates: %+v", back)
	}
}

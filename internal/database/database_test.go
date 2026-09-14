package database

import (
	"context"
	"testing"
)

func TestOpenAppliesInitialMigrationIdempotently(t *testing.T) {
	dsn := "file:" + t.TempDir() + "/resonance.db?_pragma=foreign_keys(1)"
	db, err := Open(context.Background(), dsn)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count); err != nil {
		t.Fatalf("query migrations: %v", err)
	}
	if count != 1 {
		t.Fatalf("migration count = %d, want 1", count)
	}
	for _, table := range []string{"resume_bullets", "skills", "bullet_skills", "tags", "bullet_tags"} {
		var exists bool
		if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?)", table).Scan(&exists); err != nil {
			t.Fatalf("query table %s: %v", table, err)
		}
		if !exists {
			t.Errorf("expected metadata table %s to exist", table)
		}
	}

	if err := migrate(context.Background(), db); err != nil {
		t.Fatalf("second migrate() error = %v", err)
	}
}

package db

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var migrationName = regexp.MustCompile(`^(\d{4})_.+\.sql$`)

// TestMigrationsAreContiguous fails when a migration is numbered out of order
// or skipped, which would make the apply-all loop silently miss schema.
func TestMigrationsAreContiguous(t *testing.T) {
	entries, err := os.ReadDir("migrations")
	if err != nil {
		t.Fatalf("read migrations: %v", err)
	}

	var nums []int
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		m := migrationName.FindStringSubmatch(e.Name())
		if m == nil {
			t.Fatalf("migration %q must match NNNN_name.sql", e.Name())
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			t.Fatalf("parse %q: %v", e.Name(), err)
		}
		nums = append(nums, n)
	}
	sort.Ints(nums)

	if len(nums) == 0 {
		t.Fatal("no migrations found")
	}
	for i, n := range nums {
		if n != i+1 {
			t.Fatalf("migration sequence gap: want %04d, got %04d", i+1, n)
		}
	}
}

// TestCIMigrationStepAppliesAllMigrations is the regression guard for BUG-001:
// the test job must invoke the apply-all script, not a single hard-coded file.
func TestCIMigrationStepAppliesAllMigrations(t *testing.T) {
	workflow, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatalf("read ci.yml: %v", err)
	}
	body := string(workflow)

	if !strings.Contains(body, "scripts/apply-migrations.sh") {
		t.Fatal("ci.yml must apply migrations via scripts/apply-migrations.sh")
	}
	if strings.Contains(body, "-f impl/db/migrations/") {
		t.Fatal("ci.yml must not hard-code a single migration file")
	}
}

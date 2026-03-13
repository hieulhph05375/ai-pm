package tests

import (
	"testing"
)

// TestSeedDB is a utility to seed the test database from the command line
// Usage: go test -v ./tests -run TestSeedDB
func TestSeedDB(t *testing.T) {
	if testDB == nil {
		t.Fatal("Database connection not initialized. Run with main_test.go context.")
	}

	err := SeedAll(testDB)
	if err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}
	t.Log("Database seeded successfully")
}

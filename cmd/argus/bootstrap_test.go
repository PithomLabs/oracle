package main

import (
	"os"
	"testing"
)

// --- Unit tests (no CRDB needed) ---

func TestResolveDBURL_Default(t *testing.T) {
	os.Unsetenv("ARGUS_DB_URL")
	got := resolveDBURL("", false)
	if got != defaultDBURL {
		t.Errorf("expected default %q, got %q", defaultDBURL, got)
	}
}

func TestResolveDBURL_FlagOverride(t *testing.T) {
	got := resolveDBURL("postgres://other:5432/test", true)
	if got != "postgres://other:5432/test" {
		t.Errorf("expected flag value, got %q", got)
	}
}

func TestResolveDBURL_EnvOverride(t *testing.T) {
	os.Setenv("ARGUS_DB_URL", "postgres://env-host:9999/mydb")
	defer os.Unsetenv("ARGUS_DB_URL")
	got := resolveDBURL("", false)
	if got != "postgres://env-host:9999/mydb" {
		t.Errorf("expected env value, got %q", got)
	}
}

func TestResolveDBURL_FlagTakesPrecedenceOverEnv(t *testing.T) {
	os.Setenv("ARGUS_DB_URL", "postgres://env:9999/db")
	defer os.Unsetenv("ARGUS_DB_URL")
	got := resolveDBURL("postgres://flag:1111/db", true)
	if got != "postgres://flag:1111/db" {
		t.Errorf("expected flag to take precedence, got %q", got)
	}
}

func TestIsLocalURL_Localhost(t *testing.T) {
	if !isLocalURL("postgres://root@localhost:26257/argus?sslmode=disable") {
		t.Error("expected true for localhost")
	}
}

func TestIsLocalURL_127(t *testing.T) {
	if !isLocalURL("postgres://root@127.0.0.1:26257/argus?sslmode=disable") {
		t.Error("expected true for 127.0.0.1")
	}
}

func TestIsLocalURL_Remote(t *testing.T) {
	if isLocalURL("postgres://root@otherhost:26257/argus?sslmode=disable") {
		t.Error("expected false for remote host")
	}
}

func TestIsLocalURL_Empty(t *testing.T) {
	if isLocalURL("") {
		t.Error("expected false for empty URL")
	}
}

func TestExtractDBName(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"postgres://root@localhost:26257/argus?sslmode=disable", "argus"},
		{"postgres://root@localhost:26257/defaultdb?sslmode=disable", "defaultdb"},
		{"postgres://root@localhost:26257/?sslmode=disable", "defaultdb"},
		{"", "defaultdb"},
	}
	for _, tt := range tests {
		got := extractDBName(tt.url)
		if got != tt.want {
			t.Errorf("extractDBName(%q) = %q, want %q", tt.url, got, tt.want)
		}
	}
}

func TestFindCockroachBinary_NotFound(t *testing.T) {
	// If cockroach is at a hardcoded path, we can't test the not-found case
	home, _ := os.UserHomeDir()
	knownPaths := []string{
		home + "/.local/bin/cockroach",
		"/usr/local/bin/cockroach",
	}
	for _, p := range knownPaths {
		if _, err := os.Stat(p); err == nil {
			t.Skipf("cockroach found at %s, cannot test not-found case", p)
		}
	}
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "/nonexistent")
	defer os.Setenv("PATH", oldPath)
	_, err := findCockroachBinary()
	if err == nil {
		t.Error("expected error when cockroach not in PATH")
	}
}

// --- Integration tests (need CRDB, skip with -short) ---

func TestIsCockroachRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// This tests against the real local CRDB — may or may not be running
	running := isCockroachRunning(crdbHost)
	t.Logf("CRDB running on %s:%s: %v", crdbHost, crdbSQLPort, running)
}

func TestEnsureDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if !isCockroachRunning(crdbHost) {
		t.Skip("CRDB not running, skipping ensureDatabase test")
	}
	if err := ensureDatabase(defaultDBURL); err != nil {
		t.Fatalf("ensureDatabase failed: %v", err)
	}
}

func TestApplyMigrations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if !isCockroachRunning(crdbHost) {
		t.Skip("CRDB not running, skipping migration test")
	}
	if err := applyMigrations(defaultDBURL); err != nil {
		t.Fatalf("applyMigrations failed: %v", err)
	}
	// Apply again — must be idempotent
	if err := applyMigrations(defaultDBURL); err != nil {
		t.Fatalf("applyMigrations not idempotent: %v", err)
	}
}

func TestBootstrap_ExistingDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	if !isCockroachRunning(crdbHost) {
		t.Skip("CRDB not running, skipping bootstrap test")
	}
	resolvedURL, managedCRDB, err := bootstrap("")
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}
	if managedCRDB != nil {
		t.Error("expected nil managed CRDB when CRDB already running")
		managedCRDB.Process.Signal(os.Kill)
		managedCRDB.Wait()
	}
	if resolvedURL != defaultDBURL {
		t.Errorf("expected default URL, got %q", resolvedURL)
	}
}

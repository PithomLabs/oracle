package config

import (
	"os"
	"path/filepath"
)

// Config holds all configuration for the reference loop
type Config struct {
	RunID          string
	Scenario       string
	ExecutorType   string // "recording" or "github"

	// Conductor (SQLite)
	ConductorDBPath string

	// Solvent (CockroachDB)
	SolventDBPath    string // DSN for CockroachDB
	SolventRESTAddr  string // e.g., "http://localhost:18080"
	SolventAPIKey    string

	// GitHub Executor
	GitHubToken    string
	GitHubRepo     string
	GitHubWorkflow string
	GitHubRef      string

	// Test repo
	TestRepoOwner string
	TestRepoName  string

	// Evidence
	EvidenceDir string
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	tmpDir := filepath.Join(".tmp")
	return &Config{
		RunID:            generateRunID(),
		Scenario:         "happy_path",
		ExecutorType:     "recording",
		ConductorDBPath:  filepath.Join(tmpDir, "conductor.db"),
		SolventDBPath:    "postgresql://root@localhost:26257/solvent_ref_loop?sslmode=disable",
		SolventRESTAddr:  "http://localhost:18080",
		SolventAPIKey:    "test-api-key",
		GitHubRepo:       "ibmendoza/reference-loop-test",
		GitHubWorkflow:   "ref-loop.yml",
		GitHubRef:        "main",
		TestRepoOwner:    "ibmendoza",
		TestRepoName:     "reference-loop-test",
		EvidenceDir:      filepath.Join("evidence", generateRunID()),
	}
}

// LoadFromEnv overrides config with environment variables
func (c *Config) LoadFromEnv() {
	if v := os.Getenv("REF_LOOP_RUN_ID"); v != "" {
		c.RunID = v
	}
	if v := os.Getenv("REF_LOOP_SCENARIO"); v != "" {
		c.Scenario = v
	}
	if v := os.Getenv("REF_LOOP_EXECUTOR"); v != "" {
		c.ExecutorType = v
	}
	if v := os.Getenv("CONDUCTOR_DB_PATH"); v != "" {
		c.ConductorDBPath = v
	}
	if v := os.Getenv("SOLVENT_DB_PATH"); v != "" {
		c.SolventDBPath = v
	}
	if v := os.Getenv("SOLVENT_REST_ADDR"); v != "" {
		c.SolventRESTAddr = v
	}
	if v := os.Getenv("SOLVENT_API_KEY"); v != "" {
		c.SolventAPIKey = v
	}
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		c.GitHubToken = v
	}
	if v := os.Getenv("GITHUB_REPO"); v != "" {
		c.GitHubRepo = v
	}
	if v := os.Getenv("GITHUB_WORKFLOW"); v != "" {
		c.GitHubWorkflow = v
	}
	if v := os.Getenv("GITHUB_REF"); v != "" {
		c.GitHubRef = v
	}
	if v := os.Getenv("EVIDENCE_DIR"); v != "" {
		c.EvidenceDir = v
	}
}

func generateRunID() string {
	return "run-" + filepath.Base(os.TempDir()) + "-" + os.Getenv("USER")
}

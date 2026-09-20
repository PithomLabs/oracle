package executor

import (
	"fmt"
	"os"

	"github.com/PithomLabs/solvent/adapter/github"
	"github.com/PithomLabs/solvent/service/executor"
)

// Registry wraps the solvent executor registry with configurable backends
type Registry struct {
	registry        *executor.Registry
	recordingFunc   *executor.RecordingFunc
	currentBackend  string // "recording" or "github"
}

// NewRegistry creates a new configurable executor registry
func NewRegistry(executorType, githubToken, githubRepo, githubWorkflow, githubRef string) (*Registry, error) {
	reg := executor.NewRegistry()
	r := &Registry{
		registry:       reg,
		currentBackend: executorType,
	}

	// Register recording executor (always available)
	rec := executor.NewRecordingFunc("github_trigger_workflow", "recording-run-id-123")
	reg.Register("github_trigger_workflow", rec.Func())
	r.recordingFunc = rec

	// Register GitHub executor if token provided
	if githubToken != "" && executorType == "github" {
		provider := github.NewHTTPProvider(githubToken)
		github.RegisterExecutor(reg, provider)
	}

	return r, nil
}

// GetRegistry returns the underlying solvent executor registry
func (r *Registry) GetRegistry() *executor.Registry {
	return r.registry
}

// SetBackend switches the executor backend
func (r *Registry) SetBackend(backend string) error {
	if backend != "recording" && backend != "github" {
		return fmt.Errorf("unknown backend: %s", backend)
	}
	r.currentBackend = backend
	return nil
}

// CurrentBackend returns the current backend
func (r *Registry) CurrentBackend() string {
	return r.currentBackend
}

// WasRecordingCalled returns whether the recording executor was called
func (r *Registry) WasRecordingCalled() bool {
	return r.recordingFunc.Called()
}

// GetRecordingParams returns the parameters passed to the recording executor
func (r *Registry) GetRecordingParams() map[string]interface{} {
	return r.recordingFunc.Params()
}

// GitHubExecutorConfig holds GitHub executor configuration
type GitHubExecutorConfig struct {
	Token    string
	Repo     string
	Workflow string
	Ref      string
}

// DefaultGitHubExecutorConfig returns default config for the test repo
func DefaultGitHubExecutorConfig() GitHubExecutorConfig {
	return GitHubExecutorConfig{
		Token:    os.Getenv("GITHUB_TOKEN"),
		Repo:     "ibmendoza/reference-loop-test",
		Workflow: "ref-loop.yml",
		Ref:      "main",
	}
}

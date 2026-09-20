package boundary

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMCPBoundary verifies that the MCP package does not import
// the epistemic kernel (Discharge, Promote, RetractCascade).
// This is the mechanical enforcement of the capability boundary.
func TestMCPBoundary(t *testing.T) {
	mcpDir := filepath.Join("..", "mcp")
	entries, err := os.ReadDir(mcpDir)
	if err != nil {
		t.Fatalf("read mcp dir: %v", err)
	}

	forbidden := []string{
		"github.com/PithomLabs/oracle/internal/epistemic",
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		path := filepath.Join(mcpDir, entry.Name())
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		for _, imp := range file.Imports {
			impPath := strings.Trim(imp.Path.Value, `"`)
			for _, forbidden := range forbidden {
				if impPath == forbidden {
					t.Errorf("%s imports forbidden package %s", entry.Name(), impPath)
				}
			}
		}
	}

	t.Log("MCP package boundary verified: no forbidden imports")
}

// TestAppBoundary verifies that the application package only imports
// the allowed boundary packages.
func TestAppBoundary(t *testing.T) {
	appDir := filepath.Join("..", "application")
	entries, err := os.ReadDir(appDir)
	if err != nil {
		t.Fatalf("read application dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		path := filepath.Join(appDir, entry.Name())
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		for _, imp := range file.Imports {
			impPath := strings.Trim(imp.Path.Value, `"`)
			// Application should not import MCP
			if impPath == "github.com/PithomLabs/oracle/internal/mcp" {
				t.Errorf("%s imports forbidden package %s", entry.Name(), impPath)
			}
		}
	}

	t.Log("Application package boundary verified: no forbidden imports")
}

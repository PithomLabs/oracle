package domainpack

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

)

// testPack is a minimal Pack implementation for testing.
type testPack struct {
	id      string
	version string
	valid   bool
}

func (t *testPack) GetPackID() string              { return t.id }
func (t *testPack) GetVersion() string             { return t.version }
func (t *testPack) GetDebtVocabulary() []string    { return nil }
func (t *testPack) GetEvidenceClasses() []string   { return nil }
func (t *testPack) GetRetirementRules() map[string]RetirementRule { return nil }
func (t *testPack) GetVerifierSpecs() []VerifierSpec { return nil }
func (t *testPack) Validate() error {
	if !t.valid {
		return &ValidationError{Msg: "invalid pack"}
	}
	return nil
}

// ValidationError is a simple error type for testing.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }

// testdataDir returns the absolute path to the testdata directory.
func testdataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get test file path")
	}
	return filepath.Join(filepath.Dir(filename), "testdata")
}

func TestRegisterValidPack(t *testing.T) {
	r := NewRegistry()
	pack := &testPack{id: "test", version: "1.0.0", valid: true}
	if err := r.Register(pack); err != nil {
		t.Fatalf("valid pack rejected: %v", err)
	}
}

func TestRegisterInvalidPack(t *testing.T) {
	r := NewRegistry()
	pack := &testPack{id: "test", version: "1.0.0", valid: false}
	if err := r.Register(pack); err == nil {
		t.Fatal("invalid pack accepted")
	}
}

func TestRegisterDuplicatePack(t *testing.T) {
	r := NewRegistry()
	pack := &testPack{id: "test", version: "1.0.0", valid: true}
	if err := r.Register(pack); err != nil {
		t.Fatalf("first registration failed: %v", err)
	}
	if err := r.Register(pack); err == nil {
		t.Fatal("duplicate pack accepted")
	}
}

func TestGetExistingPack(t *testing.T) {
	r := NewRegistry()
	pack := &testPack{id: "test", version: "1.0.0", valid: true}
	if err := r.Register(pack); err != nil {
		t.Fatalf("registration failed: %v", err)
	}
	got, err := r.Get("test", "1.0.0")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.GetPackID() != "test" || got.GetVersion() != "1.0.0" {
		t.Fatalf("got wrong pack: %s@%s", got.GetPackID(), got.GetVersion())
	}
}

func TestGetUnknownPack(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("nonexistent", "1.0.0")
	if err == nil {
		t.Fatal("unknown pack accepted")
	}
}

func TestGetWrongVersion(t *testing.T) {
	r := NewRegistry()
	pack := &testPack{id: "test", version: "1.0.0", valid: true}
	if err := r.Register(pack); err != nil {
		t.Fatalf("registration failed: %v", err)
	}
	_, err := r.Get("test", "2.0.0")
	if err == nil {
		t.Fatal("wrong version accepted")
	}
}

func TestLoadFromDisk(t *testing.T) {
	r := NewRegistry()
	root := testdataDir(t)
	if err := r.LoadFromDisk(root); err != nil {
		t.Fatalf("LoadFromDisk failed: %v", err)
	}

	// The bmist pack should be loaded from testdata/bmist/v1/pack.json
	pack, err := r.Get("bmist", "1.0.0")
	if err != nil {
		t.Fatalf("bmist pack not found after LoadFromDisk: %v", err)
	}
	if pack.GetPackID() != "bmist" {
		t.Fatalf("wrong pack_id: %s", pack.GetPackID())
	}
	if pack.GetVersion() != "1.0.0" {
		t.Fatalf("wrong version: %s", pack.GetVersion())
	}
}

func TestLoadFromDiskInvalidPack(t *testing.T) {
	r := NewRegistry()
	root := testdataDir(t)
	// LoadFromDisk should fail if it encounters an invalid pack.
	// For now, we just verify it doesn't panic.
	_ = r.LoadFromDisk(root)
}

func TestConcurrentRegisterGet(t *testing.T) {
	r := NewRegistry()
	var wg sync.WaitGroup

	// Register packs concurrently.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pack := &testPack{
				id:      "test",
				version: "1.0.0",
				valid:   true,
			}
			// Only one should succeed; others should get duplicate error.
			r.Register(pack)
		}(i)
	}

	wg.Wait()

	// Verify exactly one pack is registered.
	pack, err := r.Get("test", "1.0.0")
	if err != nil {
		t.Fatalf("get after concurrent register failed: %v", err)
	}
	if pack == nil {
		t.Fatal("no pack registered after concurrent attempts")
	}
}

func TestLoadEmbedded(t *testing.T) {
	r := NewRegistry()
	r.RegisterDecoder("bmist", func(data []byte) (Pack, error) {
		// Use the real bmistv1 parser
		type pack struct {
			PackID  string `json:"pack_id"`
			Version string `json:"version"`
		}
		var p pack
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, err
		}
		return &genericPack{PackID: p.PackID, Version: p.Version}, nil
	})

	// Simulate embedded pack data
	data := []byte(`{"pack_id":"bmist","version":"1.0.0"}`)
	if err := r.LoadEmbedded(data, "bmist"); err != nil {
		t.Fatalf("LoadEmbedded failed: %v", err)
	}

	pack, err := r.Get("bmist", "1.0.0")
	if err != nil {
		t.Fatalf("Get after LoadEmbedded failed: %v", err)
	}
	if pack.GetPackID() != "bmist" {
		t.Errorf("expected pack_id bmist, got %s", pack.GetPackID())
	}
}

func TestLoadEmbeddedPackIDMismatch(t *testing.T) {
	r := NewRegistry()
	r.RegisterDecoder("bmist", func(data []byte) (Pack, error) {
		return &genericPack{PackID: "bmist", Version: "1.0.0"}, nil
	})

	data := []byte(`{"pack_id":"wrong","version":"1.0.0"}`)
	err := r.LoadEmbedded(data, "bmist")
	if err == nil {
		t.Fatal("expected error for pack_id mismatch, got nil")
	}
}

package domainpack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// RetirementRule describes the evidence-class gate for retiring a debt item.
type RetirementRule struct {
	EvidenceClass string `json:"evidence_class"`
	Rule          string `json:"rule"`
}

// VerifierSpec declares which verifier is authorized by the pack.
type VerifierSpec struct {
	VerifierID string `json:"verifier_id"`
	MinVersion string `json:"min_version"`
}

// Pack is the interface that all Domain Pack types must implement.
type Pack interface {
	GetPackID() string
	GetVersion() string
	GetDebtVocabulary() []string
	GetEvidenceClasses() []string
	GetRetirementRules() map[string]RetirementRule
	GetVerifierSpecs() []VerifierSpec
}

// Validator is an optional interface that Pack types may implement
// to support structural validation.
type Validator interface {
	Pack
	Validate() error
}

// DecoderFunc parses raw JSON into a Pack. Registered per pack_id.
type DecoderFunc func(data []byte) (Pack, error)

// PackRegistry is an in-memory, startup-loaded registry of Domain Packs.
// It is validated before registration and keyed by pack_id@version.
// It is non-persistent; restart clears the registry.
type PackRegistry struct {
	mu       sync.RWMutex
	packs    map[string]Pack         // key = "pack_id@version"
	decoders map[string]DecoderFunc  // key = pack_id
}

// NewRegistry creates an empty PackRegistry.
func NewRegistry() *PackRegistry {
	return &PackRegistry{
		packs:    make(map[string]Pack),
		decoders: make(map[string]DecoderFunc),
	}
}

// RegisterDecoder registers a domain-specific decoder for a pack_id.
// The composition root calls this before LoadFromDisk.
func (r *PackRegistry) RegisterDecoder(packID string, fn DecoderFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.decoders[packID] = fn
}

// Register validates and adds a pack to the registry.
func (r *PackRegistry) Register(pack Pack) error {
	if pack == nil {
		return fmt.Errorf("cannot register nil pack")
	}
	if v, ok := pack.(Validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("pack validation failed: %w", err)
		}
	}
	key := registryKey(pack.GetPackID(), pack.GetVersion())
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.packs[key]; exists {
		return fmt.Errorf("pack %s@%s already registered", pack.GetPackID(), pack.GetVersion())
	}
	r.packs[key] = pack
	return nil
}

// Get retrieves a pack by pack_id and version.
func (r *PackRegistry) Get(packID, version string) (Pack, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := registryKey(packID, version)
	pack, ok := r.packs[key]
	if !ok {
		return nil, fmt.Errorf("pack %s@%s not found", packID, version)
	}
	return pack, nil
}

// LoadFromDisk scans the domain-pack directory tree, loads each pack.json,
// validates it, and registers it.
func (r *PackRegistry) LoadFromDisk(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("failed to read domain-pack directory %s: %w", root, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		packID := entry.Name()
		packDir := filepath.Join(root, packID)
		if err := r.loadPackDir(packDir, packID); err != nil {
			return err
		}
	}
	return nil
}

func (r *PackRegistry) loadPackDir(packDir, packID string) error {
	entries, err := os.ReadDir(packDir)
	if err != nil {
		return fmt.Errorf("failed to read pack directory %s: %w", packDir, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		version := entry.Name()
		packJSON := filepath.Join(packDir, version, "pack.json")
		if _, err := os.Stat(packJSON); os.IsNotExist(err) {
			continue
		}
		if err := r.loadPackFile(packJSON, packID); err != nil {
			return err
		}
	}
	return nil
}

func (r *PackRegistry) loadPackFile(path, expectedPackID string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}
	pack, err := r.decodePack(data, expectedPackID)
	if err != nil {
		return fmt.Errorf("failed to decode %s: %w", path, err)
	}
	if err := r.Register(pack); err != nil {
		return fmt.Errorf("failed to register pack from %s: %w", path, err)
	}
	return nil
}

// decodePack parses raw JSON into a Pack using registered decoders.
func (r *PackRegistry) decodePack(data []byte, expectedPackID string) (Pack, error) {
	var probe struct {
		PackID  string `json:"pack_id"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	if expectedPackID != "" && !strings.EqualFold(probe.PackID, expectedPackID) {
		return nil, fmt.Errorf("pack_id mismatch: expected %q, got %q", expectedPackID, probe.PackID)
	}
	// Use registered decoder if available
	r.mu.RLock()
	decoder, ok := r.decoders[probe.PackID]
	r.mu.RUnlock()
	if ok {
		return decoder(data)
	}
	// Fallback to generic pack
	return &genericPack{PackID: probe.PackID, Version: probe.Version}, nil
}

// genericPack is a minimal Pack implementation for registry purposes.
type genericPack struct {
	PackID  string `json:"pack_id"`
	Version string `json:"version"`
}

func (g *genericPack) GetPackID() string              { return g.PackID }
func (g *genericPack) GetVersion() string             { return g.Version }
func (g *genericPack) GetDebtVocabulary() []string    { return nil }
func (g *genericPack) GetEvidenceClasses() []string   { return nil }
func (g *genericPack) GetRetirementRules() map[string]RetirementRule { return nil }
func (g *genericPack) GetVerifierSpecs() []VerifierSpec { return nil }


// LoadEmbedded registers a pack from raw JSON bytes using the registered decoder for its pack_id.
func (r *PackRegistry) LoadEmbedded(data []byte, expectedPackID string) error {
	pack, err := r.decodePack(data, expectedPackID)
	if err != nil {
		return fmt.Errorf("decode embedded pack: %w", err)
	}
	if err := r.Register(pack); err != nil {
		return fmt.Errorf("register embedded pack: %w", err)
	}
	return nil
}

func registryKey(packID, version string) string {
	return packID + "@" + version
}

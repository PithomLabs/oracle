package bmistv1

import (
	"testing"
)

func TestLoadEmbeddedReturnsValidJSON(t *testing.T) {
	data := LoadEmbedded()
	if len(data) == 0 {
		t.Fatal("LoadEmbedded returned empty data")
	}

	pack, err := ParsePack(data)
	if err != nil {
		t.Fatalf("ParsePack failed on embedded data: %v", err)
	}

	if pack.GetPackID() != "bmist" {
		t.Errorf("expected pack_id bmist, got %s", pack.GetPackID())
	}
	if pack.GetVersion() != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", pack.GetVersion())
	}
	if len(pack.GetDebtVocabulary()) == 0 {
		t.Error("embedded pack has empty debt vocabulary")
	}
	if len(pack.GetRetirementRules()) == 0 {
		t.Error("embedded pack has empty retirement rules")
	}
}

func TestLoadEmbeddedIsDeterministic(t *testing.T) {
	d1 := LoadEmbedded()
	d2 := LoadEmbedded()
	if len(d1) != len(d2) {
		t.Fatalf("LoadEmbedded returned different lengths: %d vs %d", len(d1), len(d2))
	}
	for i := range d1 {
		if d1[i] != d2[i] {
			t.Fatalf("LoadEmbedded returned different data at byte %d", i)
		}
	}
}

package packetv1

import (
	"fmt"
	"strings"
)

// ReferenceKind indicates whether a reference is local or canonical.
type ReferenceKind string

const (
	RefKindLocal     ReferenceKind = "local"
	RefKindCanonical ReferenceKind = "canonical"
)

// ResolvedReference is the parsed form of a reference string.
type ResolvedReference struct {
	Kind ReferenceKind
	ID   string // for local: the local_id; for canonical: the UUID
	Raw  string // original reference string
}

// ParseReference parses "local:b1" or "canonical:belief:<uuid>"
func ParseReference(ref string) (ResolvedReference, error) {
	if strings.HasPrefix(ref, RefPrefixLocal) {
		id := strings.TrimPrefix(ref, RefPrefixLocal)
		if id == "" {
			return ResolvedReference{}, fmt.Errorf("empty local reference")
		}
		return ResolvedReference{Kind: RefKindLocal, ID: id, Raw: ref}, nil
	}
	if strings.HasPrefix(ref, RefPrefixCanonical) {
		id := strings.TrimPrefix(ref, RefPrefixCanonical)
		if !isValidUUID(id) {
			return ResolvedReference{}, fmt.Errorf("invalid canonical UUID: %s", id)
		}
		return ResolvedReference{Kind: RefKindCanonical, ID: id, Raw: ref}, nil
	}
	return ResolvedReference{}, fmt.Errorf("invalid reference format: %s", ref)
}

// isValidUUID checks if s is a valid UUID v4 string (8-4-4-4-12 hex).
func isValidUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	positions := []int{8, 13, 18, 23}
	for _, p := range positions {
		if s[p] != '-' {
			return false
		}
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue
		}
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

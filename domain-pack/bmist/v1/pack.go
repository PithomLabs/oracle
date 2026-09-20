package bmistv1

import (
	_ "embed"
)

//go:embed pack.json
var embeddedPackJSON []byte

// LoadEmbedded returns the bytes of the compiled-in BM-IST pack.
func LoadEmbedded() []byte {
	return embeddedPackJSON
}

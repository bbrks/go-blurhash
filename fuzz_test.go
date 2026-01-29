package blurhash_test

import (
	"testing"

	"github.com/bbrks/go-blurhash"
)

func FuzzComponents(f *testing.F) {
	// Seed with valid blurhash strings
	for _, test := range testFixtures {
		f.Add(test.hash)
	}

	f.Fuzz(func(t *testing.T, hash string) {
		// Should not panic on any input
		_, _, _ = blurhash.Components(hash)
	})
}

func FuzzDecode(f *testing.F) {
	// Seed with valid blurhash strings and reasonable dimensions
	for _, test := range testFixtures {
		f.Add(test.hash, 32, 32, 1)
	}

	f.Fuzz(func(t *testing.T, hash string, width, height, punch int) {
		// Limit dimensions to avoid OOM on huge allocations
		if width > 1000 || height > 1000 {
			return
		}
		if width < 0 || height < 0 {
			return
		}

		// Should not panic on any input
		_, _ = blurhash.Decode(hash, width, height, punch)
	})
}

func FuzzDecodeRoundtrip(f *testing.F) {
	// Seed with valid blurhash strings
	for _, test := range testFixtures {
		f.Add(test.hash)
	}

	f.Fuzz(func(t *testing.T, hash string) {
		// Try to decode - Components only does minimal validation,
		// so we use Decode as the source of truth for validity
		img, err := blurhash.Decode(hash, 32, 32, 1)
		if err != nil {
			// Invalid hash - skip
			return
		}

		// Get components for re-encoding
		xComp, yComp, err := blurhash.Components(hash)
		if err != nil {
			t.Errorf("decoded successfully but Components failed: %v", err)
			return
		}

		// Re-encode and verify we get a valid hash back
		newHash, err := blurhash.Encode(xComp, yComp, img)
		if err != nil {
			t.Errorf("re-encode failed: %v", err)
			return
		}

		// The new hash should also decode successfully
		_, err = blurhash.Decode(newHash, 32, 32, 1)
		if err != nil {
			t.Errorf("re-encoded hash failed to decode: %v", err)
		}
	})
}

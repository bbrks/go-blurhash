package base83_test

import (
	"testing"

	"github.com/bbrks/go-blurhash/base83"
)

func FuzzDecode(f *testing.F) {
	// Seed with valid base83 strings
	f.Add("0")
	f.Add("~")
	f.Add("00")
	f.Add("~$")
	f.Add("%%%%%%")
	f.Add("LFE.@D9F01_2%L%MIVD*9Goe-;WB")

	f.Fuzz(func(t *testing.T, s string) {
		// Should not panic on any input
		_, _ = base83.Decode(s)
	})
}

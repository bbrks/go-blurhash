package base83_test

import (
	"strings"
	"testing"

	"github.com/bbrks/go-blurhash/base83"
)

var tests = []struct {
	str string
	val int
}{
	{"3", 3},
	{"A", 10},
	{":", 70},
	{"~", 82},
	{"01", 1}, // leading zeros are "trimmed"
	{"11", 84},
	{"33", 252},
	{"~$", 6869},
	{"%%%%%%", 255172974336},
}

func TestDecodeEncode(t *testing.T) {
	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			val, err := base83.Decode(test.str)
			if err != nil {
				t.Fatalf("Decode returned unexpected error: %v", err)
			}
			if val != test.val {
				t.Errorf("Decode got unexpected result: got %d, want %d", val, test.val)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			str, err := base83.Encode(test.val, len(test.str))
			if err != nil {
				t.Fatalf("Encode returned unexpected error: %v", err)
			}
			if str != test.str {
				t.Errorf("Encode got unexpected result: got %q, want %q", str, test.str)
			}
		})
	}
}

func TestDecodeInvalidInput(t *testing.T) {
	tests := []struct {
		str string
		val int
		err error
	}{
		{"&", 0, base83.ErrInvalidInput},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			val, err := base83.Decode(test.str)
			if err == nil {
				t.Fatal("Decode should've returned error for invalid input")
			}
			if !strings.Contains(err.Error(), test.err.Error()) {
				t.Errorf("Decode returned wrong error: got %v, want %v", err, test.err)
			}
			if val != test.val {
				t.Errorf("Decode got unexpected result: got %d, want %d", val, test.val)
			}
		})
	}
}

func TestEncodeInvalidLength(t *testing.T) {
	tests := []struct {
		val    int
		length int
		str    string
	}{
		{255172974336, 3, "%%%"},
		{255172974336, 6, "%%%%%%"},
		{255172974336, 9, "000%%%%%%"},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			output, err := base83.Encode(test.val, test.length)
			if err != nil {
				t.Fatalf("Encode returned unexpected error: %v", err)
			}
			if output != test.str {
				t.Errorf("Encode got unexpected result: got %q, want %q", output, test.str)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, test := range tests {
		b.Run(test.str, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = base83.Decode("~$")
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	for _, test := range tests {
		b.Run(test.str, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = base83.Encode(6869, 2)
			}
		})
	}
}

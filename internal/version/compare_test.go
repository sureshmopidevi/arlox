package version_test

import (
	"testing"

	"github.com/sureshmopidevi/arlox/internal/version"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"0.9.1", "0.9.0", 1},
		{"0.9.0", "0.9.1", -1},
		{"0.9.0", "0.9.0", 0},
		{"1.0.0", "0.9.9", 1},
		{"v0.8.0", "0.9.0", -1},
	}
	for _, tc := range tests {
		got := version.Compare(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("Compare(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestGreater(t *testing.T) {
	if !version.Greater("0.9.1", "0.9.0") {
		t.Fatal("expected 0.9.1 > 0.9.0")
	}
	if version.Greater("0.9.0", "0.9.1") {
		t.Fatal("expected 0.9.0 not > 0.9.1")
	}
}

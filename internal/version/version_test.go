package version_test

import (
	"testing"

	"github.com/sureshmopidevi/arlox/internal/version"
)

func TestVersionNonEmpty(t *testing.T) {
	if version.Version == "" {
		t.Fatal("version.Version is empty — check internal/version/VERSION")
	}
	for _, c := range version.Version {
		if c == '\n' || c == '\r' || c == ' ' {
			t.Fatalf("version.Version has whitespace: %q", version.Version)
		}
	}
}

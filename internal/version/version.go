package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var raw string

// Version is the vibeit release version.
// Bump the VERSION file in this package when cutting a release.
var Version = strings.TrimSpace(raw)

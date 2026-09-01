package naming

import "strings"

// StackNames holds per-stack naming derived from the user workspace slug.
type StackNames struct {
	Name        string // as entered (workspace slug)
	Kebab       string // npm, go module segment
	Snake       string // Flutter package, Postgres DB
	DisplayName string
	Module      string
	WebPackage  string
	FlutterPkg  string
	DBName      string
	Org         string
}

// SplitParts splits a slug on hyphens and underscores.
func SplitParts(name string) []string {
	var parts []string
	var cur strings.Builder
	for _, r := range name {
		if r == '-' || r == '_' {
			if cur.Len() > 0 {
				parts = append(parts, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	if len(parts) == 0 {
		return []string{name}
	}
	return parts
}

// Kebab joins slug parts with hyphens.
func Kebab(name string) string {
	return strings.Join(SplitParts(name), "-")
}

// Snake joins slug parts with underscores.
func Snake(name string) string {
	return strings.Join(SplitParts(name), "_")
}

// DisplayName returns a human-readable title from the slug.
func DisplayName(name string) string {
	parts := SplitParts(name)
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

// DefaultModule returns the Go module path for a backend stack.
func DefaultModule(kebab, custom string) string {
	if custom != "" {
		return custom
	}
	return "github.com/example/" + kebab + "-backend"
}

// FromSlug builds stack names from a validated workspace slug.
func FromSlug(name, module, org string) StackNames {
	kebab := Kebab(name)
	snake := Snake(name)
	return StackNames{
		Name:        name,
		Kebab:       kebab,
		Snake:       snake,
		DisplayName: DisplayName(name),
		Module:      DefaultModule(kebab, module),
		WebPackage:  kebab,
		FlutterPkg:  snake,
		DBName:      snake,
		Org:         org,
	}
}

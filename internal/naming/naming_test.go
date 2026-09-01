package naming

import "testing"

func TestKebab(t *testing.T) {
	tests := map[string]string{
		"my-app":        "my-app",
		"my_app":        "my-app",
		"my-cool_app":   "my-cool-app",
		"a_b-c":         "a-b-c",
		"demo":          "demo",
	}
	for in, want := range tests {
		if got := Kebab(in); got != want {
			t.Errorf("Kebab(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSnake(t *testing.T) {
	tests := map[string]string{
		"my-app":        "my_app",
		"my_app":        "my_app",
		"my-cool_app":   "my_cool_app",
		"a_b-c":         "a_b_c",
		"issue_tracker": "issue_tracker",
	}
	for in, want := range tests {
		if got := Snake(in); got != want {
			t.Errorf("Snake(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDisplayName(t *testing.T) {
	if got := DisplayName("my-cool_app"); got != "My Cool App" {
		t.Fatalf("DisplayName = %q", got)
	}
}

func TestDefaultModule(t *testing.T) {
	if got := DefaultModule("my-app", ""); got != "github.com/example/my-app-backend" {
		t.Fatalf("got %q", got)
	}
	if got := DefaultModule("my-app", "custom/mod"); got != "custom/mod" {
		t.Fatalf("custom module: got %q", got)
	}
}

func TestFromSlug(t *testing.T) {
	n := FromSlug("my-cool_app", "", "com.acme")
	if n.Kebab != "my-cool-app" || n.Snake != "my_cool_app" {
		t.Fatalf("kebab/snake: %q / %q", n.Kebab, n.Snake)
	}
	if n.WebPackage != n.Kebab || n.FlutterPkg != n.Snake || n.DBName != n.Snake {
		t.Fatal("package/db mismatch")
	}
}

// PackageAlias documents Flutter Package field equivalence.
func (s StackNames) PackageAlias() string { return s.FlutterPkg }

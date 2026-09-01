package designsystems

import "testing"

func TestWebCatalogUniqueIDs(t *testing.T) {
	seen := make(map[string]bool)
	for _, s := range WebCatalog {
		if seen[s.ID] {
			t.Fatalf("duplicate id %q", s.ID)
		}
		seen[s.ID] = true
		if s.Label == "" {
			t.Fatalf("empty label for %q", s.ID)
		}
	}
}

func TestValidateWebID(t *testing.T) {
	if err := ValidateWebID("shadcn"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateWebID("invalid"); err == nil {
		t.Fatal("expected error for invalid id")
	}
}

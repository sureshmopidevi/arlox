package generate

import "testing"

func TestParseLearnedREADME(t *testing.T) {
	content := `# Learned

## 2026-01-01 — feature-a
**Pattern:** did thing

## 2026-01-02 — feature-b
**Pattern:** other
**Applied:** 2026-01-03 — promoted to rule
`
	entries := ParseLearnedREADME(content)
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}
	if entries[0].Applied {
		t.Fatal("feature-a should be unapplied")
	}
	if !entries[1].Applied {
		t.Fatal("feature-b should be applied")
	}
}

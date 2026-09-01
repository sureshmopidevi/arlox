package generate

import (
	"os"
	"path/filepath"
	"strings"
)

// LearnedEntryStatus describes one section in a learned/README.md file.
type LearnedEntryStatus struct {
	Title    string
	Applied  bool
	LineHint int // approximate line of ## heading
}

// ParseLearnedREADME splits a learned index into entries and whether each was applied.
func ParseLearnedREADME(content string) []LearnedEntryStatus {
	lines := strings.Split(content, "\n")
	var entries []LearnedEntryStatus
	var cur *LearnedEntryStatus
	var body strings.Builder

	flush := func() {
		if cur == nil {
			return
		}
		cur.Applied = strings.Contains(body.String(), "**Applied:**")
		entries = append(entries, *cur)
		cur = nil
		body.Reset()
	}

	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			flush()
			title := strings.TrimPrefix(line, "## ")
			cur = &LearnedEntryStatus{Title: title, LineHint: i + 1}
			continue
		}
		if cur != nil {
			body.WriteString(line)
			body.WriteByte('\n')
		}
	}
	flush()
	return entries
}

// SkillsStatusLocation is one stack or workspace root with a learned index.
type SkillsStatusLocation struct {
	Label   string
	Path    string
	Entries []LearnedEntryStatus
}

// SkillsStatusSummary aggregates learned/README.md across workspace and stacks.
type SkillsStatusSummary struct {
	Locations []SkillsStatusLocation
}

func (s SkillsStatusSummary) TotalEntries() int {
	n := 0
	for _, loc := range s.Locations {
		n += len(loc.Entries)
	}
	return n
}

func (s SkillsStatusSummary) UnappliedEntries() int {
	n := 0
	for _, loc := range s.Locations {
		for _, e := range loc.Entries {
			if !e.Applied {
				n++
			}
		}
	}
	return n
}

// CollectSkillsStatus scans workspace root and present stacks for learned/README.md files.
func CollectSkillsStatus(root string, present []string) SkillsStatusSummary {
	type locSpec struct {
		label string
		rel   string
	}
	specs := []locSpec{{label: "workspace", rel: ".cursor/skills/add-feature-fullstack/learned/README.md"}}
	for _, stack := range present {
		switch stack {
		case "backend":
			specs = append(specs, locSpec{"backend", "backend/.cursor/skills/add-feature-backend/learned/README.md"})
		case "web":
			specs = append(specs, locSpec{"web", "web/.cursor/skills/add-feature-web/learned/README.md"})
		case "app":
			specs = append(specs, locSpec{"app", "app/.cursor/skills/add-feature-mobile/learned/README.md"})
		}
	}

	var out SkillsStatusSummary
	for _, spec := range specs {
		path := filepath.Join(root, spec.rel)
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		entries := ParseLearnedREADME(string(raw))
		if len(entries) == 0 {
			continue
		}
		out.Locations = append(out.Locations, SkillsStatusLocation{
			Label:   spec.label,
			Path:    path,
			Entries: entries,
		})
	}
	return out
}

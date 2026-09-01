package designsystems

import "fmt"

// WebSystem describes a selectable web UI design system.
type WebSystem struct {
	ID          string
	Label       string
	Description string
}

// DefaultWebID is used when no design system is specified.
const DefaultWebID = "tailwind"

// WebCatalog lists supported web design systems.
var WebCatalog = []WebSystem{
	{ID: "tailwind", Label: "Tailwind only", Description: "Utility classes, no component library"},
	{ID: "shadcn", Label: "shadcn/ui", Description: "Tailwind + Radix primitives and CSS variables"},
	{ID: "antd", Label: "Ant Design", Description: "Enterprise React UI library"},
	{ID: "mui", Label: "Material UI", Description: "Google Material Design components"},
	{ID: "chakra", Label: "Chakra UI", Description: "Accessible component library"},
	{ID: "mantine", Label: "Mantine", Description: "Full-featured React components"},
}

// ValidateWebID returns an error if id is not in the catalog.
func ValidateWebID(id string) error {
	if id == "" {
		return fmt.Errorf("design system id cannot be empty")
	}
	for _, s := range WebCatalog {
		if s.ID == id {
			return nil
		}
	}
	return fmt.Errorf("unknown design system %q (valid: %s)", id, WebIDs())
}

// WebIDs returns comma-separated valid IDs.
func WebIDs() string {
	ids := make([]string, len(WebCatalog))
	for i, s := range WebCatalog {
		ids[i] = s.ID
	}
	out := ids[0]
	for _, id := range ids[1:] {
		out += ", " + id
	}
	return out
}

// WebLabel returns the display label for id, or id if unknown.
func WebLabel(id string) string {
	for _, s := range WebCatalog {
		if s.ID == id {
			return s.Label
		}
	}
	return id
}

// WebOverlayPath returns the embed path for a design system overlay.
func WebOverlayPath(id string) string {
	return "design-systems/web/" + id
}

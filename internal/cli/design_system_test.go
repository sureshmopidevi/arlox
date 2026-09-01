package cli

import (
	"testing"

	"github.com/sureshmopidevi/arlox/internal/designsystems"
	"github.com/sureshmopidevi/arlox/internal/generate"
	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func TestBuildDataWebDesignSystem(t *testing.T) {
	data := buildData("demo", stackFlags{webUI: "shadcn"})
	if data.WebDesignSystem != "shadcn" {
		t.Fatalf("want shadcn, got %q", data.WebDesignSystem)
	}
	if data.WebDesignSystemLabel != "shadcn/ui" {
		t.Fatalf("label %q", data.WebDesignSystemLabel)
	}
}

func TestBuildDataDefaultWebDesignSystem(t *testing.T) {
	data := buildData("demo", stackFlags{})
	if data.WebDesignSystem != designsystems.DefaultWebID {
		t.Fatalf("want %q, got %q", designsystems.DefaultWebID, data.WebDesignSystem)
	}
}

func TestResolveWebDesignSystemNonInteractiveDefault(t *testing.T) {
	id, err := resolveWebDesignSystem(stackFlags{}, []workspace.Stack{workspace.Web})
	if err != nil {
		t.Fatal(err)
	}
	if id != designsystems.DefaultWebID {
		t.Fatalf("want default, got %q", id)
	}
}

func TestResolveWebDesignSystemFlag(t *testing.T) {
	id, err := resolveWebDesignSystem(stackFlags{webUI: "antd"}, []workspace.Stack{workspace.Web})
	if err != nil {
		t.Fatal(err)
	}
	if id != "antd" {
		t.Fatalf("want antd, got %q", id)
	}
}

func TestResolveWebDesignSystemSkipsWithoutWeb(t *testing.T) {
	id, err := resolveWebDesignSystem(stackFlags{}, []workspace.Stack{workspace.Backend})
	if err != nil {
		t.Fatal(err)
	}
	if id != "" {
		t.Fatalf("want empty, got %q", id)
	}
}

var _ = generate.Data{}

package cli

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/sureshmopidevi/arlox/internal/workspace"
)

func TestSelectedVariantsFromFlags(t *testing.T) {
	flags := stackFlags{
		backend: workspace.PyFastAPI,
		web:     workspace.NextJS,
		app:     workspace.NativeIOS,
	}
	got := selectedVariantsFromFlags(flags)
	want := []stackSelection{
		{stack: workspace.Backend, variant: workspace.PyFastAPI},
		{stack: workspace.Web, variant: workspace.NextJS},
		{stack: workspace.App, variant: workspace.NativeIOS},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d selections, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("selection %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestVariantFlagsSupportValuesAndBareDefaults(t *testing.T) {
	var flags stackFlags
	cmd := &cobra.Command{Use: "test"}
	bindAddFlags(cmd, &flags)
	if err := cmd.ParseFlags([]string{"--backend=python", "--web", "--app=ios"}); err != nil {
		t.Fatal(err)
	}
	if flags.backend != workspace.PyFastAPI {
		t.Errorf("backend = %q, want %q", flags.backend, workspace.PyFastAPI)
	}
	if flags.web != workspace.ReactVite {
		t.Errorf("bare web = %q, want %q", flags.web, workspace.ReactVite)
	}
	if flags.app != workspace.NativeIOS {
		t.Errorf("app = %q, want %q", flags.app, workspace.NativeIOS)
	}
}

func TestVariantFlagsSupportSpaceSeparatedValues(t *testing.T) {
	var flags stackFlags
	cmd := &cobra.Command{Use: "test"}
	bindAddFlags(cmd, &flags)
	if err := cmd.ParseFlags([]string{"--backend", "python", "--web", "nextjs"}); err != nil {
		t.Fatal(err)
	}
	remaining, err := resolveVariantArgs(cmd.Flags().Args(), &flags, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining) != 0 {
		t.Fatalf("unexpected positional args: %v", remaining)
	}
	if flags.backend != workspace.PyFastAPI || flags.web != workspace.NextJS {
		t.Fatalf("got backend=%q web=%q", flags.backend, flags.web)
	}
}

func TestBuildDataVariantField(t *testing.T) {
	data := buildData("demo", stackFlags{backend: workspace.NodeFastify})
	if data.Variant != workspace.NodeFastify {
		t.Fatalf("Variant = %q, want %q", data.Variant, workspace.NodeFastify)
	}
}

func TestValidateSelectionsRejectsWrongStack(t *testing.T) {
	err := validateSelections([]stackSelection{{
		stack:   workspace.Web,
		variant: workspace.GoGin,
	}})
	if err == nil {
		t.Fatal("expected invalid cross-stack variant to fail")
	}
}

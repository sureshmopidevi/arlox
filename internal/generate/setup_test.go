package generate

import "testing"

func TestDBNameFromProject(t *testing.T) {
	tests := map[string]string{
		"my-app":         "my_app",
		"issue_tracker":  "issue_tracker",
		"a_b-c":          "a_b_c",
	}
	for in, want := range tests {
		if got := dbNameFromProject(in); got != want {
			t.Errorf("dbNameFromProject(%q) = %q, want %q", in, got, want)
		}
	}
}

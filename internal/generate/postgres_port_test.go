package generate

import "testing"

func TestPostgresPortFromName_stable(t *testing.T) {
	const name = "myapp"
	a := PostgresPortFromName(name)
	b := PostgresPortFromName(name)
	if a != b {
		t.Fatalf("expected stable port, got %d and %d", a, b)
	}
}

func TestPostgresPortFromName_inRange(t *testing.T) {
	for _, name := range []string{"demo", "demo2", "myapp", "x", "long-project-name-here"} {
		p := PostgresPortFromName(name)
		if p < postgresPortMin || p > postgresPortMax {
			t.Fatalf("port %d for %q out of range [%d,%d]", p, name, postgresPortMin, postgresPortMax)
		}
		if p == 5432 {
			t.Fatalf("port must not be default postgres 5432")
		}
	}
}

func TestPostgresPortFromName_differentProjects(t *testing.T) {
	if PostgresPortFromName("demo") == PostgresPortFromName("demo2") {
		t.Fatal("expected demo and demo2 to get different ports")
	}
}

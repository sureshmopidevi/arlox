package version_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sureshmopidevi/arlox/internal/version"
)

func TestRemoteLatest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("1.2.3\n"))
	}))
	t.Cleanup(srv.Close)

	got, err := version.RemoteLatest(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.2.3" {
		t.Fatalf("got %q, want 1.2.3", got)
	}
}

func TestRemoteLatestRejectsNonOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	if _, err := version.RemoteLatest(srv.URL); err == nil {
		t.Fatal("expected error for HTTP 404")
	}
}

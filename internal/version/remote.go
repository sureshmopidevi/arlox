package version

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const DefaultRemoteVERSIONURL = "https://raw.githubusercontent.com/sureshmopidevi/arlox/main/internal/version/VERSION"

// RemoteLatest reads the release version from the public VERSION file.
func RemoteLatest(url string) (string, error) {
	if url == "" {
		url = DefaultRemoteVERSIONURL
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch version: HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return "", err
	}
	v := strings.TrimSpace(string(data))
	if v == "" {
		return "", fmt.Errorf("fetch version: empty response")
	}
	return v, nil
}

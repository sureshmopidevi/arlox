package version

import (
	"strconv"
	"strings"
)

// Greater reports whether a is semantically greater than b (semver x.y.z).
func Greater(a, b string) bool {
	return Compare(a, b) > 0
}

// Compare returns -1, 0, or 1 comparing semver a and b.
func Compare(a, b string) int {
	pa := parseParts(a)
	pb := parseParts(b)
	for i := 0; i < len(pa); i++ {
		if pa[i] != pb[i] {
			return pa[i] - pb[i]
		}
	}
	return 0
}

func parseParts(v string) [3]int {
	var out [3]int
	v = strings.TrimSpace(strings.TrimPrefix(v, "v"))
	for i, part := range strings.Split(v, ".") {
		if i >= len(out) {
			break
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		out[i] = n
	}
	return out
}

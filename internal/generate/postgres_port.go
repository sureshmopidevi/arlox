package generate

import (
	"hash/fnv"
)

const (
	postgresPortMin = 5433
	postgresPortMax = 7432
)

// PostgresPortFromName returns a stable host port for Docker Postgres based on the
// project slug. The same name always maps to the same port; different names
// usually map to different ports (range 5433–7432, avoiding default 5432).
func PostgresPortFromName(name string) int {
	span := postgresPortMax - postgresPortMin + 1
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return postgresPortMin + int(h.Sum32()%uint32(span))
}

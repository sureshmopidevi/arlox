#!/usr/bin/env bash
# Append Go bin to ~/.zshrc if missing. Idempotent.
set -euo pipefail

GOPATH_BIN="$(go env GOPATH)/bin"
MARKER="# Go (go install → \$GOPATH/bin)"
ZSHRC="${HOME}/.zshrc"

if grep -q 'go/bin' "${ZSHRC}" 2>/dev/null || grep -q 'GOPATH/bin' "${ZSHRC}" 2>/dev/null; then
  echo "Go bin already in ${ZSHRC}"
  exit 0
fi

cat >> "${ZSHRC}" <<EOF

${MARKER}
export PATH="\${HOME}/go/bin:\${PATH}"
EOF

echo "Added ${HOME}/go/bin to ${ZSHRC}"
echo "Run: source ~/.zshrc"

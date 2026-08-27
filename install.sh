#!/usr/bin/env bash
# One-step vibeit setup: build, install, add to PATH (~/.zshrc + current shell if sourced).
#
# Usage (recommended — works immediately in this terminal):
#   cd ~/vibeit && source ./install.sh
#
# Or:
#   ./install.sh          # installs; open a new terminal tab after
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]:-$0}")" && pwd)"
GOPATH_BIN="$(go env GOPATH)/bin"
INSTALLED="${GOPATH_BIN}/vibeit"
ZSHRC="${HOME}/.zshrc"
SOURCED=0
[[ "${BASH_SOURCE[0]:-}" != "${0:-install.sh}" ]] && SOURCED=1

echo ""
echo "== vibeit install =="

if ! command -v go >/dev/null 2>&1; then
  echo "error: go not found. Install Go first: https://go.dev/dl/" >&2
  return 1 2>/dev/null || exit 1
fi

echo "→ building..."
(cd "${ROOT}" && go build -o bin/vibeit ./cmd/vibeit)

echo "→ installing to ${INSTALLED}..."
(cd "${ROOT}" && go install ./cmd/vibeit)

if [[ ! -x "${INSTALLED}" ]]; then
  echo "error: install failed — ${INSTALLED} missing" >&2
  return 1 2>/dev/null || exit 1
fi

# Permanent PATH in ~/.zshrc (idempotent)
if [[ -f "${ZSHRC}" ]] && ! grep -qE '(go/bin|GOPATH/bin)' "${ZSHRC}" 2>/dev/null; then
  echo "→ adding ${GOPATH_BIN} to ~/.zshrc"
  cat >> "${ZSHRC}" <<EOF

# Go (go install → \$(go env GOPATH)/bin) — added by vibeit install.sh
export PATH="${GOPATH_BIN}:\$PATH"
EOF
elif [[ -f "${ZSHRC}" ]]; then
  echo "→ ~/.zshrc already has Go bin on PATH"
else
  echo "→ no ~/.zshrc found (skipped PATH persist)"
fi

# Current shell — works when: source ./install.sh
export PATH="${GOPATH_BIN}:${PATH}"

echo ""
echo "✓ vibeit installed: ${INSTALLED}"
echo "✓ version: $("${INSTALLED}" version 2>/dev/null || "${INSTALLED}" -v)"
echo ""

if [[ "${SOURCED}" -eq 1 ]]; then
  echo "Ready in this terminal. Try:"
  echo "  cd ~/Projects && vibeit create myapp"
else
  echo "Installed. For THIS terminal, run:"
  echo "  source ./install.sh"
  echo ""
  echo "Or open a new terminal tab, then:"
  echo "  vibeit create myapp"
fi
echo ""

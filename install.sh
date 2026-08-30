#!/usr/bin/env bash
# One-step arlox installer for macOS & Linux.
#
# Remote 1-line install for any user (recommended):
#   curl -fsSL https://raw.githubusercontent.com/sureshmopidevi/arlox/main/install.sh | bash
#
# Local install from cloned repository:
#   cd ~/arlox && ./install.sh
set -euo pipefail

echo ""
echo "== arlox installer =="

if ! command -v go >/dev/null 2>&1; then
  echo "error: Go is required to install arlox." >&2
  echo "Please install Go first: https://go.dev/dl/" >&2
  exit 1
fi

GOPATH_BIN="$(go env GOPATH)/bin"
INSTALLED="${GOPATH_BIN}/arlox"

# 1. Build and install arlox
SCRIPT_DIR=""
if [[ -f "${BASH_SOURCE[0]:-}" ]]; then
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" 2>/dev/null && pwd || true)"
fi

if [[ -n "${SCRIPT_DIR}" && -f "${SCRIPT_DIR}/cmd/arlox/main.go" ]]; then
  echo "→ building from local source in ${SCRIPT_DIR}..."
  (cd "${SCRIPT_DIR}" && go build -o bin/arlox ./cmd/arlox)
  (cd "${SCRIPT_DIR}" && go install ./cmd/arlox)
else
  echo "→ installing latest release via go install..."
  go install github.com/sureshmopidevi/arlox/cmd/arlox@latest
fi

if [[ ! -x "${INSTALLED}" ]]; then
  echo "error: install failed — ${INSTALLED} not found" >&2
  exit 1
fi

# 2. Automatically link to system bin if writable (instant global PATH access in all terminals)
LINKED=0
for sysdir in "/opt/homebrew/bin" "/usr/local/bin"; do
  if [[ -d "${sysdir}" && -w "${sysdir}" ]]; then
    if ln -sf "${INSTALLED}" "${sysdir}/arlox" 2>/dev/null; then
      echo "→ linked to ${sysdir}/arlox (instant global PATH access)"
      LINKED=1
      break
    fi
  fi
done

# 3. Permanently persist to shell rc files (~/.zshrc, ~/.bashrc)
SHELL_CONFIGS=("${HOME}/.zshrc" "${HOME}/.bashrc" "${HOME}/.config/fish/config.fish")
PERSISTED=0

for rc in "${SHELL_CONFIGS[@]}"; do
  if [[ -f "${rc}" ]]; then
    if ! grep -qE '(go/bin|GOPATH/bin)' "${rc}" 2>/dev/null; then
      echo "→ adding ${GOPATH_BIN} to ${rc}"
      cat >> "${rc}" <<EOF

# Go bin (added by arlox installer)
export PATH="${GOPATH_BIN}:\$PATH"
EOF
      PERSISTED=1
    fi
  fi
done

# Export for current subshell
export PATH="${GOPATH_BIN}:${PATH}"

echo ""
echo "✓ arlox installed successfully!"
echo "✓ location : ${INSTALLED}"
echo "✓ version  : $("${INSTALLED}" version 2>/dev/null || "${INSTALLED}" -v)"
echo ""
echo "Try running:"
echo "  arlox create myapp"
echo ""

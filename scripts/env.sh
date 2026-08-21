#!/usr/bin/env bash
# Source from Makefile or verify scripts to ensure vibeit is on PATH.
set -euo pipefail

_vibeit_repo_root() {
  cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd
}

ensure_vibeit_path() {
  local repo_root
  repo_root="$(_vibeit_repo_root)"
  local gopath_bin
  gopath_bin="$(go env GOPATH)/bin"

  if [[ ":${PATH}:" != *":${gopath_bin}:"* ]]; then
    export PATH="${gopath_bin}:${PATH}"
  fi

  if command -v vibeit >/dev/null 2>&1; then
    return 0
  fi

  if [[ -x "${repo_root}/bin/vibeit" ]]; then
    export PATH="${repo_root}/bin:${PATH}"
    return 0
  fi

  echo "error: vibeit not found. Run: make install" >&2
  return 1
}

# Optional: common local tool paths (only append if dir exists and not already on PATH)
ensure_local_tool_paths() {
  local dir
  for dir in \
    "$(go env GOPATH)/bin" \
    "${HOME}/development/flutter/bin" \
    "/opt/homebrew/bin" \
    "/usr/local/bin"; do
    if [[ -d "${dir}" && ":${PATH}:" != *":${dir}:"* ]]; then
      export PATH="${dir}:${PATH}"
    fi
  done
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  ensure_vibeit_path
  ensure_local_tool_paths
  echo "PATH ok — vibeit=$(command -v vibeit)"
  echo "flutter=$(command -v flutter 2>/dev/null || echo 'not found')"
  echo "go=$(command -v go 2>/dev/null || echo 'not found')"
  echo "node=$(command -v node 2>/dev/null || echo 'not found')"
fi

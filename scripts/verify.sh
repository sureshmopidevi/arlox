#!/usr/bin/env bash
# Local smoke test for vibeit. Uses a TEMP directory — nothing in your cwd.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/env.sh
source "${ROOT}/scripts/env.sh"
ensure_vibeit_path
ensure_local_tool_paths

VIBEIT="$(command -v vibeit)"
WORKDIR="${TMPDIR:-/tmp}/vibeit-verify-$$"
DEMO="${WORKDIR}/demo"
DEMO2="${WORKDIR}/demo2"

cleanup() {
  rm -rf "${WORKDIR}"
}
trap cleanup EXIT

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║  vibeit VERIFY — automated test only                        ║"
echo "║  Creates demo projects in a TEMP folder (not your cwd).       ║"
echo "║  Everything is deleted when this script finishes.             ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "  temp dir: ${WORKDIR}"
echo "  (removed automatically on exit — you do NOT need to delete anything)"
echo ""
echo "  To create a real project interactively, run:"
echo "    cd ~/Projects && vibeit create myapp"
echo ""

echo "== vibeit verify =="
echo "vibeit: ${VIBEIT}"
mkdir -p "${WORKDIR}"
cd "${WORKDIR}"

echo ""
echo "== 1. create all stacks =="
NO_COLOR=1 "${VIBEIT}" create demo --backend --web --app

test -f "${DEMO}/backend/go.mod" || { echo "fail: backend missing"; exit 1; }
test -f "${DEMO}/web/package.json" || { echo "fail: web missing"; exit 1; }
test -f "${DEMO}/app/pubspec.yaml" || { echo "fail: app missing"; exit 1; }
test -f "${DEMO}/demo.code-workspace" || { echo "fail: workspace missing"; exit 1; }
test -f "${DEMO}/Makefile" || { echo "fail: root Makefile missing"; exit 1; }
test -f "${DEMO}/.cursor/skills/add-feature-fullstack/SKILL.md" || { echo "fail: fullstack skill missing"; exit 1; }
echo "ok"

echo ""
echo "== 1b. root make status (partial workspace) =="
NO_COLOR=1 "${VIBEIT}" create demo2 --backend >/dev/null
cd "${DEMO2}"
STATUS="$(NO_COLOR=1 make status)"
echo "${STATUS}" | grep -q "backend  yes"
echo "${STATUS}" | grep -q "web      no"
echo "${STATUS}" | grep -q "app      no"
cd "${WORKDIR}"
echo "ok"

echo ""
echo "== 2. no duplicates =="
OUT="$(NO_COLOR=1 "${VIBEIT}" create demo --backend --web --app 2>&1)"
echo "${OUT}" | grep -q "skipped (already exists)" || { echo "fail: expected skip"; echo "${OUT}"; exit 1; }
echo "ok"

echo ""
echo "== 3b. workspace sync (web on disk but missing from .code-workspace) =="
python3 - <<PY
import json, pathlib
root = pathlib.Path("${DEMO}")
ws = root / "demo.code-workspace"
data = json.loads(ws.read_text())
data["folders"] = [f for f in data["folders"] if f["path"] != "web"]
ws.write_text(json.dumps(data, indent=2) + "\n")
PY
NO_COLOR=1 "${VIBEIT}" create demo --backend --web --app >/dev/null
python3 - <<PY
import json, pathlib, sys
ws = pathlib.Path("${DEMO}/demo.code-workspace")
paths = {f["path"] for f in json.loads(ws.read_text())["folders"]}
if "web" not in paths:
    sys.exit("fail: web not synced into workspace")
PY
echo "ok"

echo ""
echo "== 3c. add after partial create =="
cd "${DEMO2}"
NO_COLOR=1 "${VIBEIT}" add --web
test -f "${DEMO2}/web/package.json" || { echo "fail: web not added"; exit 1; }
grep -q '"web"' "${DEMO2}/demo2.code-workspace" || { echo "fail: workspace not merged"; exit 1; }
echo "ok"

echo ""
echo "== 4. backend build =="
cd "${DEMO}/backend"
go mod tidy
go build -o /dev/null ./cmd/server
echo "ok"

echo ""
echo "== 5. web build =="
cd "${DEMO}/web"
npm install --silent
npm run build
echo "ok"

echo ""
echo "== 6. flutter analyze =="
if ! command -v flutter >/dev/null 2>&1; then
  echo "skip (flutter not on PATH)"
else
  cd "${DEMO}/app"
  flutter pub get >/dev/null
  flutter analyze
  echo "ok"
fi

echo ""
echo "== 7. guardrails =="
test -f "${DEMO}/backend/.cursor/rules/karpathy.mdc"
test -f "${DEMO}/web/.cursor/rules/tailwind.mdc"
test -f "${DEMO}/app/.cursor/skills/add-feature/learned/README.md"
echo "ok"

echo ""
echo "== all checks passed =="
echo ""
echo "  temp dir ${WORKDIR} will be removed now."

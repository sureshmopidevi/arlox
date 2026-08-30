#!/usr/bin/env bash
# Local smoke test for arlox. Uses a TEMP directory — nothing in your cwd.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/env.sh
source "${ROOT}/scripts/env.sh"
ensure_arlox_path
ensure_local_tool_paths

if [[ -x "${ROOT}/bin/arlox" ]]; then
  ARLOX="${ROOT}/bin/arlox"
else
  ARLOX="$(command -v arlox)"
fi
WORKDIR="${TMPDIR:-/tmp}/arlox-verify-$$"
DEMO="${WORKDIR}/demo"
DEMO2="${WORKDIR}/demo2"

cleanup() {
  rm -rf "${WORKDIR}"
}
trap cleanup EXIT

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║  arlox VERIFY — automated test only                          ║"
echo "║  Creates demo projects in a TEMP folder (not your cwd).       ║"
echo "║  Everything is deleted when this script finishes.             ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "  temp dir: ${WORKDIR}"
echo "  (removed automatically on exit — you do NOT need to delete anything)"
echo ""
echo "  To create a real project interactively, run:"
echo "    cd ~/Projects && arlox create myapp"
echo ""

echo "== arlox verify =="
echo "arlox: ${ARLOX}"
EXPECTED_VERSION="$(tr -d '[:space:]' < "${ROOT}/internal/version/VERSION")"
GOT_VERSION="$("${ARLOX}" version)"
echo "${GOT_VERSION}" | grep -qx "arlox ${EXPECTED_VERSION}" || {
  echo "fail: version mismatch — got '${GOT_VERSION}', want 'arlox ${EXPECTED_VERSION}'"
  exit 1
}
echo "$("${ARLOX}" -v)" | grep -qx "arlox ${EXPECTED_VERSION}" || {
  echo "fail: -v mismatch"
  exit 1
}
echo "version: ${GOT_VERSION} ok"
mkdir -p "${WORKDIR}"
cd "${WORKDIR}"

echo ""
echo "== 1. create all stacks =="
NO_COLOR=1 "${ARLOX}" create demo --backend --web --app

test -f "${DEMO}/backend/go.mod" || { echo "fail: backend missing"; exit 1; }
test -f "${DEMO}/backend/go.sum" || { echo "fail: backend go.sum missing"; exit 1; }
test -f "${DEMO}/web/package.json" || { echo "fail: web missing"; exit 1; }
test -d "${DEMO}/web/node_modules" || { echo "fail: web node_modules missing"; exit 1; }
test -f "${DEMO}/app/pubspec.yaml" || { echo "fail: app missing"; exit 1; }
test -f "${DEMO}/demo.code-workspace" || { echo "fail: workspace missing"; exit 1; }
python3 - <<PY
import json, pathlib, sys
ws = pathlib.Path("${DEMO}/demo.code-workspace")
data = json.loads(ws.read_text())
folders = data.get("folders")
if not isinstance(folders, list) or len(folders) < 1:
    sys.exit(f"fail: folders must be non-empty array, got {folders!r}")
paths = {f["path"] for f in folders}
names = {f["path"]: f["name"] for f in folders}
if "." not in paths:
    sys.exit("fail: workspace missing root folder path '.'")
if not {"backend", "web", "app"}.issubset(paths):
    sys.exit(f"fail: missing stack folders in {paths}")
for path, want in {".": "demo", "backend": "demo-backend", "web": "demo-web", "app": "demo-app"}.items():
    if names.get(path) != want:
        sys.exit(f"fail: folder {path} name want {want!r}, got {names.get(path)!r}")
PY
# project identity stamped into each stack
grep -q '^module github.com/example/demo-backend$' "${DEMO}/backend/go.mod" || { echo "fail: backend module name"; exit 1; }
grep -q '"name": "demo"' "${DEMO}/web/package.json" || { echo "fail: web package name"; exit 1; }
grep -q '^name: demo$' "${DEMO}/app/pubspec.yaml" || { echo "fail: app package name"; exit 1; }
grep -q 'APP_NAME=Demo' "${DEMO}/backend/configs/local/app.env.example" || { echo "fail: backend APP_NAME"; exit 1; }
test -f "${DEMO}/backend/configs/local/app.env" || { echo "fail: backend app.env missing"; exit 1; }
test -f "${DEMO}/web/.env" || { echo "fail: web .env missing"; exit 1; }
test -f "${DEMO}/app/.dart_tool/package_config.json" || { echo "fail: app pub get missing"; exit 1; }
test -f "${DEMO}/Makefile" || { echo "fail: root Makefile missing"; exit 1; }
test -f "${DEMO}/.cursor/skills/add-feature-fullstack/SKILL.md" || { echo "fail: fullstack skill missing"; exit 1; }
test -f "${DEMO}/.cursor/rules/versioning.mdc" || { echo "fail: versioning rule missing"; exit 1; }
echo "ok"

echo ""
echo "== 1b. root make status (partial workspace) =="
NO_COLOR=1 "${ARLOX}" create demo2 --backend >/dev/null
cd "${DEMO2}"
STATUS="$(NO_COLOR=1 make status)"
echo "${STATUS}" | grep -q "backend  yes"
echo "${STATUS}" | grep -q "web      no"
echo "${STATUS}" | grep -q "app      no"
cd "${WORKDIR}"
echo "ok"

echo ""
echo "== 2. no duplicates =="
OUT="$(NO_COLOR=1 "${ARLOX}" create demo --backend --web --app 2>&1)"
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
NO_COLOR=1 "${ARLOX}" create demo --backend --web --app >/dev/null
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
NO_COLOR=1 "${ARLOX}" add --web
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
test -f "${DEMO}/backend/.cursor/skills/add-feature-backend/SKILL.md" || { echo "fail: backend add-feature skill missing"; exit 1; }
test -f "${DEMO}/web/.cursor/skills/add-feature-web/SKILL.md" || { echo "fail: web add-feature skill missing"; exit 1; }
test -f "${DEMO}/app/.cursor/skills/add-feature-mobile/learned/README.md" || { echo "fail: app add-feature skill missing"; exit 1; }
echo "ok"

echo ""
echo "== all checks passed =="
echo ""
echo "  temp dir ${WORKDIR} will be removed now."

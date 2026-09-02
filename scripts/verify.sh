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

have_tools() {
  local tool
  MISSING_TOOL=""
  for tool in "$@"; do
    if ! command -v "${tool}" >/dev/null 2>&1; then
      MISSING_TOOL="${tool}"
      return 1
    fi
  done
}

check_variant_manifest() {
  local manifest="$1"
  local stack="$2"
  local variant="$3"
  python3 - "${manifest}" "${stack}" "${variant}" <<'PY'
import json
import pathlib
import sys

path, stack, variant = sys.argv[1:]
data = json.loads(pathlib.Path(path).read_text())
if data.get("stack") != stack:
    raise SystemExit(f"fail: {path} stack is {data.get('stack')!r}, want {stack!r}")
if data.get("variant") != variant:
    raise SystemExit(f"fail: {path} variant is {data.get('variant')!r}, want {variant!r}")
PY
}

package_has_script() {
  local dir="$1"
  local script="$2"
  node -e '
    const pkg = require(process.argv[1]);
    process.exit(pkg.scripts && pkg.scripts[process.argv[2]] ? 0 : 1);
  ' "${dir}/package.json" "${script}"
}

run_package_checks() {
  local dir="$1"
  if package_has_script "${dir}" test; then
    (cd "${dir}" && npm test)
  fi
  (cd "${dir}" && npm run build)
}

run_variant_checks() {
  local stack="$1"
  local variant="$2"
  local dir="$3"

  case "${stack}:${variant}" in
    backend:go)
      (cd "${dir}" && go test ./...)
      ;;
    backend:python)
      (cd "${dir}" && uv run pytest)
      ;;
    backend:node-express|backend:node-fastify)
      run_package_checks "${dir}"
      ;;
    backend:java)
      (cd "${dir}" && mvn test)
      ;;
    web:*)
      run_package_checks "${dir}"
      ;;
    app:flutter)
      (cd "${dir}" && flutter analyze && flutter test)
      ;;
    app:react-pwa)
      run_package_checks "${dir}"
      ;;
    app:ios)
      (cd "${dir}" && swift test)
      ;;
    app:android)
      if [[ -x "${dir}/gradlew" ]]; then
        (cd "${dir}" && ./gradlew test)
      else
        (cd "${dir}" && gradle test)
      fi
      ;;
    *)
      echo "fail: no smoke command for ${stack}:${variant}"
      return 1
      ;;
  esac
}

smoke_variant() {
  local stack="$1"
  local variant="$2"
  local marker="$3"
  local cursor_asset="$4"
  shift 4
  local name="variant-${stack}-${variant}"
  local project="${WORKDIR}/${name}"
  local stack_dir="${project}/${stack}"

  printf '  %-8s %-14s ' "${stack}" "${variant}"
  if ! have_tools "$@"; then
    echo "skip (missing tool: ${MISSING_TOOL})"
    return 0
  fi

  NO_COLOR=1 "${ARLOX}" create "${name}" "--${stack}=${variant}" >/dev/null
  test -f "${stack_dir}/${marker}" || {
    echo "fail: ${variant} marker ${marker} missing"
    return 1
  }
  test -f "${stack_dir}/.origin-manifest.json" || {
    echo "fail: ${variant} manifest missing"
    return 1
  }
  check_variant_manifest "${stack_dir}/.origin-manifest.json" "${stack}" "${variant}"
  test -f "${stack_dir}/${cursor_asset}" || {
    echo "fail: ${variant} cursor asset ${cursor_asset} missing"
    return 1
  }
  if ! run_variant_checks "${stack}" "${variant}" "${stack_dir}" >"${project}/smoke.log" 2>&1; then
    echo "fail: ${variant} generated checks"
    tail -n 40 "${project}/smoke.log"
    return 1
  fi
  echo "ok"
}

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
echo "== 8. per-variant scaffold smoke =="
cd "${WORKDIR}"
BACKEND_CURSOR=".cursor/skills/add-feature-backend/SKILL.md"
WEB_CURSOR=".cursor/skills/add-feature-web/SKILL.md"
APP_CURSOR=".cursor/skills/add-feature-mobile/SKILL.md"

smoke_variant backend go go.mod "${BACKEND_CURSOR}" go
smoke_variant backend python pyproject.toml "${BACKEND_CURSOR}" uv
smoke_variant backend node-express package.json "${BACKEND_CURSOR}" node npm
smoke_variant backend node-fastify package.json "${BACKEND_CURSOR}" node npm
smoke_variant backend java pom.xml "${BACKEND_CURSOR}" java mvn

smoke_variant web react-vite package.json "${WEB_CURSOR}" node npm
smoke_variant web nextjs package.json "${WEB_CURSOR}" node npm
smoke_variant web vue package.json "${WEB_CURSOR}" node npm
smoke_variant web svelte package.json "${WEB_CURSOR}" node npm
smoke_variant web angular angular.json "${WEB_CURSOR}" node npm
smoke_variant web nuxt package.json "${WEB_CURSOR}" node npm

smoke_variant app flutter pubspec.yaml "${APP_CURSOR}" flutter
smoke_variant app react-pwa package.json ".cursor/skills/add-feature-pwa/SKILL.md" node npm
smoke_variant app ios Package.swift ".cursor/skills/add-feature-ios/SKILL.md" swift xcodebuild
smoke_variant app android build.gradle.kts ".cursor/skills/add-feature-android/SKILL.md" gradle

echo "ok"

echo ""
echo "== all checks passed =="
echo ""
echo "  temp dir ${WORKDIR} will be removed now."

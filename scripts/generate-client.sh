#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SWAGGER_URL="${1:-https://api.linkbreakers.com/internal/openapi/api/v1/api.swagger.json}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

echo "Fetching OpenAPI spec from ${SWAGGER_URL}"
curl -fsSL "${SWAGGER_URL}" -o "${TMP_DIR}/openapi-3.2.json"

python3 - "${TMP_DIR}" <<'PY'
import json
import pathlib
import sys

tmp_dir = pathlib.Path(sys.argv[1])
src = tmp_dir / "openapi-3.2.json"
dst = tmp_dir / "openapi-3.1.json"

spec = json.loads(src.read_text())
spec["openapi"] = "3.1.0"
spec.pop("jsonSchemaDialect", None)
dst.write_text(json.dumps(spec, indent=2) + "\n")
PY

echo "Generating Go client"
openapi-generator-cli generate \
  -i "${TMP_DIR}/openapi-3.1.json" \
  -g go \
  -o "${TMP_DIR}/generated" \
  --skip-validate-spec \
  --global-property apiTests=false,modelTests=false,apiDocs=false,modelDocs=false \
  --additional-properties=packageName=linkbreakers,isGoSubmodule=true,enumClassPrefix=true

rm -f "${ROOT_DIR}/internal/client/generated/"*.go
cp "${TMP_DIR}/generated/"*.go "${ROOT_DIR}/internal/client/generated/"

NEW_VERSION="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["info"]["version"])' "${TMP_DIR}/openapi-3.2.json")"
printf '%s\n' "${NEW_VERSION}" > "${ROOT_DIR}/OPENAPI_VERSION"

echo "Updated generated client to ${NEW_VERSION}"

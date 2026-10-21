#!/usr/bin/env bash
# Regenerates every API command from the published OpenAPI document, so the
# CLI gains, changes and loses commands exactly as the API does.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SWAGGER_URL="${1:-https://api.linkbreakers.com/internal/openapi/api/v1/api.swagger.json}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

if [ -f "${SWAGGER_URL}" ]; then
  cp "${SWAGGER_URL}" "${TMP_DIR}/openapi.json"
else
  echo "Fetching OpenAPI spec from ${SWAGGER_URL}"
  curl -fsSL "${SWAGGER_URL}" -o "${TMP_DIR}/openapi.json"
fi

(cd "${ROOT_DIR}" && go run ./internal/gen/commands "${TMP_DIR}/openapi.json" internal/cli/operations_gen.go)

NEW_VERSION="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["info"]["version"])' "${TMP_DIR}/openapi.json")"
printf '%s\n' "${NEW_VERSION}" > "${ROOT_DIR}/OPENAPI_VERSION"
echo "Generated commands for API ${NEW_VERSION}"

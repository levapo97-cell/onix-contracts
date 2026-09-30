#!/usr/bin/env bash
# Tests de contrato: valida los ejemplos contra los JSON Schema.
# Es la red de seguridad del polyrepo: si un servicio (Go o Rust) cambia el shape,
# el ejemplo deja de validar y CI falla. Se ejecuta en local y en CI (`make test`).
# Requisitos: node/npx (ajv-cli se baja con npx).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCHEMAS="$ROOT/schemas"
EXAMPLES="$ROOT/examples"
AJV=(npx --yes ajv-cli@5 validate --spec=draft2020 --strict=false)

fail=0
check() { # <schema> <example>
  echo "→ $(basename "$2")  vs  $(basename "$1")"
  if ! "${AJV[@]}" -s "$1" -d "$2"; then fail=1; fi
}

check "$SCHEMAS/event.raw.schema.json"   "$EXAMPLES/event.raw.example.json"
check "$SCHEMAS/event.clean.schema.json" "$EXAMPLES/event.clean.example.json"

if [ "$fail" -ne 0 ]; then
  echo "✗ Contrato roto: algún ejemplo no valida contra su schema." >&2
  exit 1
fi
echo "✓ Todos los ejemplos validan contra sus schemas."

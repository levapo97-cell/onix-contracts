#!/usr/bin/env bash
# Codegen de onix-contracts: JSON Schema -> tipos Go + tipos Rust.
# Fuente única de verdad = schemas/*.schema.json. NO editar a mano los tipos generados.
# Requisitos: node/npx (quicktype se baja solo con npx). Se ejecuta en CI y en local (`make codegen`).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCHEMAS="$ROOT/schemas"
GO_OUT="$ROOT/go/onixcontracts/events.gen.go"
RUST_OUT="$ROOT/rust/src/events.gen.rs"

QT=(npx --yes quicktype@26)

# Esquemas con tipo de nivel superior (el nombre del tipo sale del "title" de cada schema).
# mcp.methods.schema.json solo define $defs (sub-tipos), no un top-level: se documenta en README.
SRC=(
  "$SCHEMAS/event.raw.schema.json"
  "$SCHEMAS/event.norm.schema.json"
  "$SCHEMAS/event.clean.schema.json"
  "$SCHEMAS/ws.message.schema.json"
)

mkdir -p "$(dirname "$GO_OUT")" "$(dirname "$RUST_OUT")"

echo "→ Generando tipos Go en $GO_OUT"
"${QT[@]}" \
  --src-lang schema \
  --lang go \
  --package onixcontracts \
  "${SRC[@]}" \
  -o "$GO_OUT"

echo "→ Generando tipos Rust en $RUST_OUT"
"${QT[@]}" \
  --src-lang schema \
  --lang rust \
  --visibility public \
  --derive-debug \
  "${SRC[@]}" \
  -o "$RUST_OUT"

echo "✓ Codegen completo. Revisa y commitea los archivos *.gen.*"

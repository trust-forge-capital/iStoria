#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

go build -o bin/istoria .

echo "Built: $PROJECT_ROOT/bin/istoria"

#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT"

echo "Building frontend..."
cd "$FRONTEND_DIR"
npm run build

echo "Copying build to backend..."
rm -rf "$BACKEND_DIR/internal/server/ui/"*
cp -r "$FRONTEND_DIR/build/"* "$BACKEND_DIR/internal/server/ui/"

echo "Done! Run 'cd $BACKEND_DIR && go run cmd/server/main.go' to start the server."

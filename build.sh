#!/bin/bash

set -e
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
    "android/arm64"
)

if [ -d dist ]; then
    echo "Cleaning dist directory..."
    rm -rf dist
fi

echo "Generating code..."
go run ./scripts/ts.go


echo "Building Web UI..."
pnpm web:build
rm -rf ./cli/helpers/web_dist/*
mkdir -p ./cli/helpers/web_dist
cp -r client/web/dist/* cli/helpers/web_dist/

for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    
    echo "----------------------------------------"
    echo "Building: GOOS=$GOOS, GOARCH=$GOARCH"

    EXT=""
    if [ "$GOOS" = "windows" ]; then
        EXT=".exe"
    fi

    OUT_DIR="dist/${GOOS}-${GOARCH}"
    mkdir -p "$OUT_DIR"

    echo "Building CLI..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUT_DIR/nagare$EXT" ./cli

    echo "Running wire for Gateway..."
    (cd gateway && dix wire --workspace)

    echo "Building Gateway..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUT_DIR/nagare-gateway$EXT" ./gateway
    
    echo "Done: $platform"
done


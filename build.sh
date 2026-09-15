#!/bin/bash

set -e
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

if ! command -v dix &> /dev/null; then
    go install github.com/smtdfc/dix@latest
fi

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

    echo "Building Gateway..."
    cd gateway
    dix wire . --workspace
    GOOS=$GOOS GOARCH=$GOARCH go build -o "../$OUT_DIR/nagare-gateway$EXT" .
    cd ..

    echo "Done: $platform"
done


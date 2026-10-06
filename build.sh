#!/usr/bin/env bash
set -e

WEB_DIST_DIR="client/web/dist"
CLI_WEB_DIST_DIR="cli/helpers/web_dist"
DIST_DIR="dist"
MODULES_DIR="$DIST_DIR/modules"
WEB_DIR="$DIST_DIR/web"
RUST_MODULES=(
    "crates/nagare_vector"
)

mkdir -p "$DIST_DIR"
mkdir -p "$MODULES_DIR"
mkdir -p "$WEB_DIR"

echo "Generating TypeScript code..."
go run ./scripts/ts.go

echo "Building Nagare Web UI..."
cd client/web
vite build
cd - > /dev/null

rm -rf "$CLI_WEB_DIST_DIR"
rm -rf $WEB_DIR
cp -r "$WEB_DIST_DIR" "$CLI_WEB_DIST_DIR"
cp -r "$WEB_DIST_DIR" "$WEB_DIR"

echo "Building Rust modules..."
for module in "${RUST_MODULES[@]}"; do
    cargo build \
        --release \
        --manifest-path "$module/Cargo.toml"

    module_name="$(basename "$module")"

    case "$(uname -s)" in
        Linux*)
            library_name="lib${module_name}.so"
            ;;
        Darwin*)
            library_name="lib${module_name}.dylib"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            library_name="${module_name}.dll"
            ;;
        *)
            echo "Unsupported platform: $(uname -s)"
            exit 1
            ;;
    esac

    cp "target/release/$library_name" \
        "$MODULES_DIR/$library_name"
done

echo "Building Nagare CLI..."
cd cli
go build -o "../$DIST_DIR/nagare"
cd - > /dev/null

echo "Building Nagare Gateway..."
cd gateway
dix wire --workspace
go build -o "../$DIST_DIR/nagare-gateway"
cd - > /dev/null

echo "Build complete!"
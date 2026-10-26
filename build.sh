#!/usr/bin/env bash
set -e

JSON_FILE=".release-please-manifest.json"

if [ ! -f "$JSON_FILE" ]; then
    echo "File not found: $JSON_FILE"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "No jq command found. Please install jq to parse JSON files."
    exit 1
fi

VERSION=$(jq -r '.[ "." ]' "$JSON_FILE")
echo "Detected version: $VERSION"

ROOT_DIR="$(pwd)"
CLIENT_WEB_SOURCE="$ROOT_DIR/client/web"
GATEWAY_SOURCE="$ROOT_DIR/gateway"
CLI_SOURCE="$ROOT_DIR/cli"
CORE_SOURCE="$ROOT_DIR/core"
GO_BUILTIN_PLUGINS_SOURCE="$ROOT_DIR/plugin/builtins"
DESKTOP_SOURCE="$ROOT_DIR/client/desktop"
WEB_DIST_DIR="$ROOT_DIR/client/web/dist"
CLI_WEB_DIST_DIR="$ROOT_DIR/cli/helpers/web_dist"
DIST_DIR="$ROOT_DIR/dist"
WEB_DIR="$DIST_DIR/web"
CORE_PLUGINS="$DIST_DIR/plugins/core"
BUILTIN_PLUGINS="$DIST_DIR/plugins/builtins"
CLI_SCRIPT="$CLI_SOURCE/main.go"
RUST_CORE_PLUGIN=(
    "crates/nagare_vector"
)

GO_BUILTIN_PLUGINS=(
    "telegram"
)

mkdir -p "$DIST_DIR"
mkdir -p "$CORE_PLUGINS"
mkdir -p "$WEB_DIR"
mkdir -p "$BUILTIN_PLUGINS"

case "$(uname -s)" in
    Linux*|Darwin*)
        exe_ext=""
        ;;
    MINGW*|MSYS*|CYGWIN*)
        exe_ext=".exe"
        ;;
    *)
        exe_ext=""
        ;;
esac

echo "Generating TypeScript code..."
go run ./scripts/ts.go

echo "Building Nagare Web UI..."
cd "$CLIENT_WEB_SOURCE"
pnpm install --frozen-lockfile
pnpm build
cd - > /dev/null

rm -rf "$CLI_WEB_DIST_DIR"
rm -rf "$WEB_DIR"
cp -r "$WEB_DIST_DIR" "$CLI_WEB_DIST_DIR"
cp -r "$WEB_DIST_DIR" "$WEB_DIR"

echo "Building Rust core plugins ..."
for plugin in "${RUST_CORE_PLUGIN[@]}"; do
    (
        cd "$plugin" || exit 1
        go run "$CLI_SCRIPT" plugin build
        cp ./*.nagare_plugin "$CORE_PLUGINS/"
    )
done

echo "Building Go builtin plugins ..."
for plugin in "${GO_BUILTIN_PLUGINS[@]}"; do
    (
        cd "$GO_BUILTIN_PLUGINS_SOURCE/$plugin" || exit 1
        go run "$CLI_SCRIPT" plugin build
        cp ./*.nagare_plugin "$BUILTIN_PLUGINS/"
    )
done

echo "Generating Go sources ..."
goverter gen "$CORE_SOURCE/mappers/..."
goverter gen "$GATEWAY_SOURCE/..."


echo "Building Desktop app"
pnpm desktop:build
cp "$DESKTOP_SOURCE/bin/nagare-desktop${exe_ext}" "$DIST_DIR"

echo "Building Nagare CLI..."
cd "$CLI_SOURCE"
go build -ldflags "-X github.com/smtdfc/nagare/cli/cmd.Version=$VERSION"  -o "$DIST_DIR/nagare${exe_ext}"
cd - > /dev/null

echo "Building Nagare Gateway..."
cd "$GATEWAY_SOURCE"
dix wire --workspace
go build -o "$DIST_DIR/nagare-gateway${exe_ext}"
cd - > /dev/null


echo "Build complete!"

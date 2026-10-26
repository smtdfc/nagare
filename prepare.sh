ROOT_DIR="$(pwd)"
CORE_SOURCE="$ROOT_DIR/core"
GATEWAY_SOURCE="$ROOT_DIR/gateway"

go install github.com/smtdfc/dix@latest
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
go install github.com/jmattheis/goverter/cmd/goverter@latest
pnpm install

mkdir -p ./cli/helpers/web_dist

echo "Generating Go sources ..."
goverter gen "$CORE_SOURCE/mappers/..."
goverter gen "$GATEWAY_SOURCE/..."

echo "Generating TypeScript code..."
go run ./scripts/ts.go

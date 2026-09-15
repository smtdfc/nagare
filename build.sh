mkdir -p dist/{gateway,cli}

go build -o dist/cli/nagare ./cli

cd gateway
dix wire .  --workspace
go build -o ../dist/gateway/nagare-gateway .
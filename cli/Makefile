APP_NAME := nagare

.PHONY: run
run:
	go run .

.PHONY: build
build:
	go build -o bin/$(APP_NAME) .

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME).exe .

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux .

.PHONY: build-macos
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME)-macos .

.PHONY: build-all
build-all:
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME).exe .
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux .
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME)-macos .

.PHONY: clean
clean:
	rm -rf bin

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...
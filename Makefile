GOLANGCI_VERSION := latest

default: build

build:
	go build -v .

ensure_deps:
	go mod tidy

install_devtools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)

lint:
	golangci-lint run ./...

lint_fix:
	golangci-lint run --fix ./...

# CGO required to run the race detector
test:
	go test -test.v -race -cover ./...

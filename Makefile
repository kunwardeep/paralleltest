export GOSUMDB := off
export GOFLAGS := -v -mod=vendor
GOLANGCI_VERSION := latest

default: build

build:
	go build "$(MAIN_PKG)"

ensure_deps:
	go mod tidy
	go mod vendor
	cd tools  && go mod tidy
	cd tools && go mod vendor

# GOFLAGS=-mod=mod: This ensures Go resolves dependencies via the go.mod file
install_devtools:
	GOFLAGS=-mod=mod go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)

clean:
	go clean $(MAIN_PKG)

lint:
	golangci-lint run ./...

lint_fix:
	golangci-lint run --fix ./...

# CGO required to run the race detector
test:
	go test -test.v -race -cover ./...

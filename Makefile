export GOSUMDB := off
export GOFLAGS := -v -mod=vendor
GOLANGCI_VERSION := latest

default: build

build:
	go build "$(MAIN_PKG)"

ensure_deps:
	go mod tidy
	go mod vendor

install_devtools:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)

clean:
	go clean $(MAIN_PKG)

lint:
	golangci-lint run ./...

lint_fix:
	golangci-lint run --fix ./...

# CGO required to run the race detector
test:
	go test -test.v -race -cover ./...

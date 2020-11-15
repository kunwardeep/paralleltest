export GOSUMDB := off
export GOFLAGS := -v -mod=vendor
GOLANGCI_VERSION := v1.32.2

default: build

build:
	go build "$(MAIN_PKG)"

ensure_deps:
	go mod tidy
	go mod vendor
	cd tools  && go mod tidy
	cd tools && go mod vendor

install_devtools:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCI_VERSION)

clean:
	go clean $(MAIN_PKG)
	rm -f profile_service

lint:
	golangci-lint run ./...

lint_fix:
	golangci-lint run --fix ./...

# CGO required to run the race detector
test:
	go test -test.v -race -cover ./...

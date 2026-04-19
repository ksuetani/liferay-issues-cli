.PHONY: build install clean test lint release snapshot

BINARY := issues
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X github.com/david-truong/liferay-issues-cli/internal/config.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install:
	go install -ldflags "$(LDFLAGS)" .

clean:
	rm -f $(BINARY)

test:
	go test ./...

lint:
	golangci-lint run

# Create a local snapshot release (no publish)
snapshot:
	goreleaser release --snapshot --clean

# Tag and release (requires GITHUB_TOKEN)
release:
	@test -n "$(tag)" || (echo "Usage: make release tag=v1.0.0" && exit 1)
	git tag $(tag)
	git push origin $(tag)
	goreleaser release --clean

VERSION ?= $(if $(GITHUB_REF_NAME),$(GITHUB_REF_NAME),$(shell git describe --tags --always --dirty))
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

BUILD_DIR := build

# Build targets
WHAT ?= minipaas
SUFFIX ?= $(if $(findstring windows,$(GOOS)),.exe,)

OUTPUT ?= $(WHAT)-$(GOOS)-$(GOARCH)$(SUFFIX)

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	@for target in $(WHAT); do \
		fn="$(BUILD_DIR)/$(OUTPUT)"; \
		echo "Building $$fn..."; \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-extldflags -static -X main.Version=$(VERSION)" -o $$fn ./cmd/$$target; \
	done
	@echo "Build complete."


.PHONY: fmt
fmt:
	@go fmt github.com/sombrahq/minipaas/minipaas-cli/...


PHONY: test
test:
	@go test --cover -parallel=1 -v -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | sort -rnk3



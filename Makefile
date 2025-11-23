include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the cmd/server/main.go application
.PHONY: run
run:
	go run ./cmd/server

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

BINARY_NAME=vpn-server
BUILD_DIR=build

PLATFORMS = linux/amd64
#PLATFORMS = \
#	linux/amd64 \
#	linux/arm64 \
#	windows/amd64 \
#	darwin/amd64 \
#	darwin/arm64

## all: clean and build the cmd/server application
.PHONY: build
all: clean build

## build: build the cmd/server application.
## : You can specify platforms using PLATPHORMS variable `make build PLATFORMS="darwin/amd64 linux/amd64"`
.PHONY: build
build:
	@echo "Building binaries..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		ext=""; \
		if [ "$$GOOS" = "windows" ]; then ext=".exe"; fi; \
		output="$(BUILD_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH$$ext"; \
		echo " → $$output"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$output ./cmd/server || exit 1; \
	done
#	@echo 'Building cmd/api...'
#	go build -ldflags='-s' -o=./bin/server ./cmd/server
#	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/server

clean:
	@rm -rf $(BUILD_DIR)
#
#docker: PLATFORMS="linux/amd64" build docker-build
#
#DOCKER_IMAGE_NAME="anpotashev/mpdgo:1.0"
#BUILD_NUM = $(shell date +%Y%m%d%H%M%S)
#
#docker:
#	@echo building docker image $(DOCKER_IMAGE_NAME)-$(BUILD_NUM)
#	@docker build -t $(DOCKER_IMAGE_NAME)-$(BUILD_NUM) .
#	@echo pushing docker image $(DOCKER_IMAGE_NAME)-$(BUILD_NUM)
#	@docker push $(DOCKER_IMAGE_NAME)-$(BUILD_NUM)
#

## test/coverage: run test coverage
#.PHONY: test/coverage
test/coverage:
	go test -v -cover  -coverprofile=cover.txt ./...
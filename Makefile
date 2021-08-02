# =========================================================================== #
# HELPERS
# =========================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# =========================================================================== #
# DEVELOPMENT
# =========================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api


# =========================================================================== #
# QUALITY CONTROL
# =========================================================================== #

## test: test the code
.PHONY: test
test:
	go test ./...

## audit: tidy dependencies, format, vet, and test the code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting the code...'
	go fmt ./...
	@echo 'Vetting the code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# =========================================================================== #
# BUILD
# =========================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/api ./cmd/api

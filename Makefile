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

arch = 386
build_dir = linux_${arch}
service_name = homebugh-api

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -o=./bin/${service_name} ./cmd/api
	GOOS=linux GOARCH=${arch} go build -o=./bin/${build_dir}/${service_name} ./cmd/api


# =========================================================================== #
# PRODUCTION
# =========================================================================== #

production_host = homebugh.info
deploy_user = deploy
deploy_dir = ~/apps/homebugh-api/

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/${build_dir}/${service_name} ${deploy_user}@${production_host}:${deploy_dir}
	ssh -t deploy@${production_host} 'sudo systemctl restart ${service_name}'

## production/deploy/systemdconfig: deploy the systemd config to production
.PHONY: production/deploy/systemdconfig
production/deploy/systemdconfig:
	rsync -P ./remote/production/${service_name}.service ${deploy_user}@${production_host}:~
	ssh -t ${deploy_user}@${production_host} '\
		sudo mv ~/${service_name}.service /etc/systemd/system/ \
		&& sudo systemctl enable ${service_name} \
		&& sudo systemctl restart ${service_name} \
	'

## production/status/api: prints api service status
.PHONY: production/status/api
production/status/api:
	ssh -t ${deploy_user}@${production_host} 'sudo systemctl status ${service_name}'
	@echo "\n------------------------------------------"
	curl https://api.${production_host}:8080/health
	@echo ""

## production/restart/api: restarts api service status
.PHONY: production/restart/api
production/restart/api:
	ssh -t ${deploy_user}@${production_host} 'sudo systemctl restart ${service_name}'
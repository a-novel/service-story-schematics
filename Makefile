# Run tests.
test:
	bash -c "set -m; bash '$(CURDIR)/scripts/test.sh'"

# Check code quality.
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run
	npx prettier . --check
	sqlfluff lint

# Generate mocked interfaces for Go tests.
mocks:
	rm -rf `find . -type d -name mocks`
	go run github.com/vektra/mockery/v3@v3.2.5

# Reformat code so it passes the code style lint checks.
format:
	go mod tidy
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run --fix
	npx prettier . --write
	sqlfluff fix

# Lint OpenAPI specs.
openapi-lint:
	npx @redocly/cli lint ./docs/api.yaml

# Generate OpenAPI docs.
openapi-generate:
	go generate ./...

run-infra:
	podman compose -p "${APP_NAME}" -f "${PWD}/build/podman-compose.yaml" up -d

run-infra-down:
	podman compose -p "${APP_NAME}" -f "${PWD}/build/podman-compose.yaml" down

# Run the API
run-api:
	bash -c "set -m; bash '$(CURDIR)/scripts/run.sh'"

install:
	pipx install sqlfluff

.PHONY: api

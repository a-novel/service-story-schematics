#!/bin/bash

APP_NAME="service-story-schematics-test"
PODMAN_FILE="$PWD/build/podman-compose.test.yaml"
TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# Ensure containers are properly shut down when the program exits abnormally.
int_handler()
{
    podman compose -p "${APP_NAME}" -f "${PODMAN_FILE}" down --volume
}
trap int_handler INT

# Setup test containers.
podman compose -p "${APP_NAME}" -f "${PODMAN_FILE}" up -d --build --pull-always

POSTGRES_DSN=${POSTGRES_DSN_TEST} go run cmd/migrations/main.go

# shellcheck disable=SC2046
go run ${TEST_TOOL_PKG} --format pkgname -- -count=1 -cover $(go list ./... | grep -v /mocks | grep -v /models/api | grep -v /test)

# Normal execution: containers are shut down.
podman compose -p "${APP_NAME}" -f "${PODMAN_FILE}" down --volume

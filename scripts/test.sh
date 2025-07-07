clear
#!/bin/bash

APP_NAME="${APP_NAME}-test"
PODMAN_FILE="$PWD/build/podman-compose.test.yaml"
TEST_TOOL_PKG="gotest.tools/gotestsum@latest"

# Ensure containers are properly shut down when the program exits abnormally.
int_handler()
{
    podman compose -p "${APP_NAME}" -f "${PODMAN_FILE}" down --volume
}
trap int_handler INT

# Setup test containers.
podman compose -p "${APP_NAME}" -f "${PODMAN_FILE}" up -d

# Unlike regular tests, DAO tests require to run in isolated transactions. This is because they are the only
# tests that cannot rely on randomized data (they expect a predictable output).
# Other tests run in integration mode, meaning they use random data for the DAO tests.
export DAO_DSN="${DSN_DAO_TEST}"
export DSN="${DSN_INTEGRATION_TEST}"
export PORT="${PORT_TEST}"
export JSON_KEYS_URL="${JSON_KEYS_SERVICE_TEST_URL}"

# shellcheck disable=SC2046
go run ${TEST_TOOL_PKG} --format pkgname -- -count=1 -cover $(go list ./... | grep -v /mocks | grep -v /codegen)

# Normal execution: containers are shut down.
podman compose -p "${APP_NAME}" -f "${PODMAN_FILE}" down --volume

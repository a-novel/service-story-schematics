---
outline: deep
---

# Containerized

You can import the story-schematics service as a container.

::: info

This service requires an instance of [JSON Keys](https://a-novel.github.io/service-json-keys/service/containerized.html)
to run.

:::

::: code-group

```yaml [podman]
# https://github.com/containers/podman-compose
services:
  json-keys-postgres:
    image: ghcr.io/a-novel/service-json-keys/database:v1
    networks:
      - api
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: json-keys
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - json-keys-postgres-data:/var/lib/postgresql/data/

  story-schematics-postgres:
    image: ghcr.io/a-novel/service-story-schematics/database:v1
    networks:
      - api
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: story-schematics
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - json-keys-postgres-data:/var/lib/postgresql/data/

  json-keys-service:
    image: ghcr.io/a-novel/service-json-keys/standalone:v1
    depends_on:
      json-keys-postgres:
        condition: service_healthy
    environment:
      POSTGRES_DSN: postgres://postgres:postgres@json-keys-postgres:5432/json-keys?sslmode=disable
      APP_MASTER_KEY: fec0681a2f57242211c559ca347721766f8a3acd8ed2e63b36b3768051c702ca
    networks:
      - api

  story-schematics-postgres-migrations:
    image: ghcr.io/a-novel/service-story-schematics/jobs/migrations:v1
    depends_on:
      story-schematics-postgres:
        condition: service_healthy
    networks:
      - api
    environment:
      POSTGRES_DSN: postgres://postgres:postgres@story-schematics-postgres:5432/json-keys?sslmode=disable

  story-schematics-service:
    image: ghcr.io/a-novel/service-story-schematics/api:v1
    depends_on:
      story-schematics-postgres:
        condition: service_healthy
      story-schematics-postgres-migrations:
        condition: service_completed_successfully
    environment:
      POSTGRES_DSN: postgres://postgres:postgres@story-schematics-postgres:5432/story-schematics?sslmode=disable
      JSON_KEYS_SERVICE_URL: http://json-keys-service:8080
      OPENAI_TOKEN: [your_OPENAI_TOKEN]
    networks:
      - api

networks:
  api: {}

volumes:
  json-keys-postgres-data:
  story-schematics-postgres-data:
```

:::

## Standalone image (local)

For local development or CI purposes, you can also load a standalone version that runs all the necessary jobs
before starting the service.

::: warning
The standalone image takes longer to boot, and it is not suited for production use.
:::

::: code-group

```yaml [podman]
# https://github.com/containers/podman-compose
services:
  json-keys-postgres:
    image: ghcr.io/a-novel/service-json-keys/database:v1
    networks:
      - api
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: json-keys
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - json-keys-postgres-data:/var/lib/postgresql/data/

  story-schematics-postgres:
    image: ghcr.io/a-novel/service-story-schematics/database:v1
    networks:
      - api
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: story-schematics
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - json-keys-postgres-data:/var/lib/postgresql/data/

  json-keys-service:
    image: ghcr.io/a-novel/service-json-keys/standalone:v1
    depends_on:
      json-keys-postgres:
        condition: service_healthy
    environment:
      POSTGRES_DSN: postgres://postgres:postgres@json-keys-postgres:5432/json-keys?sslmode=disable
      APP_MASTER_KEY: fec0681a2f57242211c559ca347721766f8a3acd8ed2e63b36b3768051c702ca
    networks:
      - api

  story-schematics-service:
    image: ghcr.io/a-novel/service-story-schematics/standalone:v1
    depends_on:
      story-schematics-postgres:
        condition: service_healthy
    environment:
      POSTGRES_DSN: postgres://postgres:postgres@story-schematics-postgres:5432/story-schematics?sslmode=disable
      JSON_KEYS_SERVICE_URL: http://json-keys-service:8080
      OPENAI_TOKEN: [your_OPENAI_TOKEN]
    networks:
      - api

networks:
  api: {}

volumes:
  json-keys-postgres-data:
  story-schematics-postgres-data:
```

:::

## Configuration

Configuration is done through environment variables.

### Required variables

You must provide the following variables for the service to run correctly.

| Variable                | Description                                                               |
| ----------------------- | ------------------------------------------------------------------------- |
| `POSTGRES_DSN`          | Connection string to the Postgres database.                               |
| `JSON_KEYS_SERVICE_URL` | URL to the JSON Keys service, used for key management.                    |
| `OPENAI_TOKEN`          | Token to authenticate with the [Groq API](https://console.groq.com/home). |

### Optional variables

Generic configuration.

| Variable   | Description                                     | Default                    |
| ---------- | ----------------------------------------------- | -------------------------- |
| `APP_NAME` | Name of the application, used for tracing.      | `story-schematics-service` |
| `ENV`      | Provide information on the current environment. |                            |
| `DEBUG`    | Activate debug mode for logs.                   | `false`                    |

API configuration.

| Variable                     | Description                                                                             | Default |
| ---------------------------- | --------------------------------------------------------------------------------------- | ------- |
| `API_PORT`                   | Port to run the API on.                                                                 | `8080`  |
| `API_MAX_REQUEST_SIZE`       | Maximum request size for the API.<br/>Provided as a number of bytes.                    | `2MB`   |
| `API_TIMEOUT_READ`           | Read timeout for the API.<br/>Provided as a duration string.                            | `5s`    |
| `API_TIMEOUT_READ_HEADER`    | Header read timeout for the API.<br/>Provided as a duration string.                     | `3s`    |
| `API_TIMEOUT_WRITE`          | Write timeout for the API.<br/>Provided as a duration string.                           | `10s`   |
| `API_TIMEOUT_IDLE`           | Idle timeout for the API.<br/>Provided as a duration string.                            | `30s`   |
| `APITimeoutRequest`          | Request timeout for the API.<br/>Provided as a duration string.                         | `15s`   |
| `API_CORS_ALLOWED_ORIGINS`   | CORS allowed origins for the API.<br/>Provided as a list of values separated by commas. | `*`     |
| `API_CORS_ALLOWED_HEADERS`   | CORS allowed headers for the API.<br/>Provided as a list of values separated by commas. | `*`     |
| `API_CORS_ALLOW_CREDENTIALS` | Whether to allow credentials in CORS requests.                                          | `false` |
| `API_CORS_MAX_AGE`           | CORS max age for the API.<br/>Provided as a number of seconds.                          | `3600`  |

Tracing configuration (with [Sentry](https://sentry.io/)).

| Variable               | Description                                                          | Default                               |
| ---------------------- | -------------------------------------------------------------------- | ------------------------------------- |
| `SENTRY_DSN`           | Sentry DSN for tracing.<br/>Tracing will be disabled if omitted.     |                                       |
| `SENTRY_RELEASE`       | Release information for Sentry logs.                                 |                                       |
| `SENTRY_FLUSH_TIMEOUT` | Timeout for flushing Sentry logs.<br/>Provided as a duration string. | `2s`                                  |
| `SENTRY_ENVIRONMENT`   | Which environment to attach logs to.                                 | Uses the value from `ENV` variable.   |
| `SENTRY_DEBUG`         | Activate debug mode for Sentry.                                      | Uses the value from `DEBUG` variable. |

We use [Groq](https://console.groq.com/home) for AI processing. Their platform is based of OpenAI's API, so you can
interchange them with an actual OpenAI token, or any other OpenAI-compatible API.

| Variable          | Description                                                                                    | Default                            |
| ----------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------- |
| `OPENAI_BASE_URL` | Change this URL to use a different OpenAI-compatible API.                                      | `"https://api.groq.com/openai/v1"` |
| `OPENAI_MODEL`    | Select the AI model to use for chat completions. Available options may depend on your provider | `"openai/gpt-oss-120b"`            |

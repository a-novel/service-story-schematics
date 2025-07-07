# Story-Schematics service

[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/agorastoryverse)](https://twitter.com/agorastoryverse)
[![Discord](https://img.shields.io/discord/1315240114691248138?logo=discord)](https://discord.gg/rp4Qr8cA)

<hr />

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-novel/service-story-schematics)
![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/a-novel/service-story-schematics)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/a-novel/service-story-schematics)

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/a-novel/service-story-schematics/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-novel/service-story-schematics)](https://goreportcard.com/report/github.com/a-novel/service-story-schematics)
[![codecov](https://codecov.io/gh/a-novel/service-story-schematics/graph/badge.svg?token=uc71lIIr8G)](https://codecov.io/gh/a-novel/service-story-schematics)

![Coverage graph](https://codecov.io/gh/a-novel/service-story-schematics/graphs/sunburst.svg?token=uc71lIIr8G)

<hr />

This is a quickstart document to test the project locally.

You can find the API documentation on the [repository GitHub page](https://a-novel.github.io/service-story-schematics/).

Want to contribute? Check the [contribution guidelines](CONTRIBUTING.md).

# Use in a project

You can import this application as a docker image. Below is an example using
[podman compose](https://docs.podman.io/en/latest/markdown/podman-compose.1.html).

```yaml
services:
  json-keys-postgres:
    image: docker.io/library/postgres:17
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
    image: docker.io/library/postgres:17
    networks:
      - api
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: story-schematics
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - story-schematics-postgres-data:/var/lib/postgresql/data/

  json-keys-service:
    image: ghcr.io/a-novel/service-json-keys/standalone:v0
    depends_on:
      - json-keys-postgres
    environment:
      PORT: 8080
      ENV: local
      APP_NAME: json-keys-service
      DSN: postgres://postgres:postgres@json-keys-postgres:5432/json-keys?sslmode=disable
      # Dummy key used only for local environment. Consider using a secure, private key in production.
      MASTER_KEY: fec0681a2f57242211c559ca347721766f8a3acd8ed2e63b36b3768051c702ca
      # Used for tracing purposes, can be omitted.
      # SENTRY_DSN: [your_sentry_dsn]
      # SERVER_NAME: json-keys-service-prod
      # RELEASE: v0.1.2
      # Set the following if you want to debug the service locally.
      # DEBUG: true
    networks:
      - api

  story-schematics-service:
    image: ghcr.io/a-novel/service-story-schematics/api:v0
    depends_on:
      - story-schematics-postgres
      - json-keys-service
    ports:
      # Expose the service on port 4001 on the local machine.
      - "4021:8080"
    environment:
      PORT: 8080
      ENV: local
      APP_NAME: story-schematics-service
      DSN: postgres://postgres:postgres@story-schematics-postgres:5432/story-schematics?sslmode=disable
      JSON_KEYS_URL: http://json-keys-service:8080/v1
      # You need a Groq API Key to access the Groq API.
      # https://console.groq.com/keys
      GROQ_TOKEN: "[your_groq_token]"
      # Used for tracing purposes, can be omitted.
      # SENTRY_DSN: [your_sentry_dsn]
      # SERVER_NAME: story-schematics-service-prod
      # RELEASE: v0.1.2
      # Set the following if you want to debug the service locally.
      # DEBUG: true
    networks:
      - api

networks:
  api: {}

volumes:
  story-schematics-postgres-data:
  json-keys-postgres-data:
```

Available tags includes:

- `latest`: latest versioned image
- `vx`: versioned images, pointing to a specific version. Partial versions are supported. When provided, the
  latest subversion is used.\
  examples: `v0`, `v0.1`, `v0.1.2`
- `branch`: get the latest version pushed to a branch. Any valid branch name can be used.\
  examples: `master`, `fix/something`

# Run locally

## Pre-requisites

- [Golang](https://go.dev/doc/install)
- [Node.js](https://nodejs.org/en/download/)
- [Python](https://www.python.org/downloads/)
  - Install [pipx](https://pipx.pypa.io/stable/installation/) to install command-line tools.
- [Podman](https://podman.io/docs/installation)
  - Install [podman-compose](https://github.com/containers/podman-compose)

    ```bash
    # Pipx
    pipx install podman-compose

    # Brew
    brew install podman-compose
    ```

- Make

  ```bash
  # Debian / Ubuntu
  sudo apt-get install build-essential

  # macOS
  brew install make
  ```

  For Windows, you can use [Make for Windows](https://gnuwin32.sourceforge.net/packages/make.htm)

## Setup environment

Create a `.envrc` file from the template:

```bash
cp .envrc.template .envrc
```

Then fill the missing secret variables. Once your file is ready:

```bash
source .envrc
```

> You may use tools such as [direnv](https://direnv.net/), otherwise you'll need to source the env file on each new
> terminal session.

Install the external dependencies:

```bash
make install
```

## Run infrastructure

```bash
make run-infra
# To turn down:
# make run-infra-down
```

> You may skip this step if you already have the global infrastructure running.

## Et Voilà!

```bash
make run-api
```

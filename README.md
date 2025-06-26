# Story-Schematics service

[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/agora_ecrivains)](https://twitter.com/agora_ecrivains)
[![Discord](https://img.shields.io/discord/1315240114691248138?logo=discord)](https://discord.gg/rp4Qr8cA)

<hr />

This is a quickstart document to test the project locally.

You can find the API documentation on the [repository GitHub page](https://a-novel.github.io/service-story-schematics/).

Want to contribute? Check the [contribution guidelines](CONTRIBUTING.md).

# Use in a project

You can import this application as a docker image. Below is an example using
[podman compose](https://docs.podman.io/en/latest/markdown/podman-compose.1.html).

```yaml
services:
  authentication-postgres:
    image: docker.io/library/postgres:17
    networks:
      - api
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: authentication
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - authentication-postgres-data:/var/lib/postgresql/data/

  story-schematics-postgres:
    image: docker.io/library/postgres:17
    networks:
      - api
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: story_schematics
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - story-schematics-postgres-data:/var/lib/postgresql/data/

  # ======================================================================
  # Authentication service is a required dependency for the API.
  # You can see more options in the service documentation.
  # https://github.com/a-novel/service-authentication
  # ======================================================================
  authentication-rotate-keys-job:
    image: ghcr.io/a-novel/service-authentication/jobs/rotatekeys:v0
    depends_on:
      - authentication-postgres
    environment:
      ENV: local
      APP_NAME: authentication-service-rotate-keys-job
      DSN: postgres://postgres:postgres@authentication-postgres:5432/authentication?sslmode=disable
      MASTER_KEY: fec0681a2f57242211c559ca347721766f8a3acd8ed2e63b36b3768051c702ca
      DEBUG: true
    networks:
      - api

  authentication-service:
    image: ghcr.io/a-novel/service-authentication/api:v0
    depends_on:
      - authentication-postgres
    ports:
      - "4001:8080"
    environment:
      PORT: 8080
      ENV: local
      APP_NAME: authentication-service
      DSN: postgres://postgres:postgres@authentication-postgres:5432/authentication?sslmode=disable
      MASTER_KEY: fec0681a2f57242211c559ca347721766f8a3acd8ed2e63b36b3768051c702ca
      SMTP_SANDBOX: true
      AUTH_PLATFORM_URL_UPDATE_EMAIL: http://localhost:4001/update-email
      AUTH_PLATFORM_URL_UPDATE_PASSWORD: http://localhost:4001/update-password
      AUTH_PLATFORM_URL_REGISTER: http://localhost:4001/register
      DEBUG: true
    networks:
      - api

  # ======================================================================
  # Our actual service.
  # ======================================================================
  story-schematics-service:
    image: ghcr.io/a-novel/service-story-schematics/api:v0
    depends_on:
      - authentication-service
    ports:
      - "4011:8080"
    environment:
      PORT: 8080
      ENV: local
      APP_NAME: authentication-service
      DSN: postgres://postgres:postgres@story-schematics-postgres:5432/story_schematics?sslmode=disable
      # Used to handle authentication requests. Make sure it
      # points to a valid authentication service instance.
      AUTH_API_URL: http://localhost:4001
      DEBUG: true
      # Required to use the AI features.
      GROQ_TOKEN: ${GROQ_TOKEN}
    networks:
      - api

networks:
  api: {}

volumes:
  authentication-postgres-data:
  story-schematics-postgres-data:
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

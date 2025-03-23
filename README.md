![Story Schematics Service](./docs/assets/service%20story%20schematics%20banner.png)

[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/agora_ecrivains)](https://twitter.com/agora_ecrivains)
[![Discord](https://img.shields.io/discord/1315240114691248138?logo=discord)](https://discord.gg/rp4Qr8cA)

<hr />

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-novel/story-schematics)
![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/a-novel/story-schematics)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/a-novel/story-schematics)

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/a-novel/story-schematics/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-novel/story-schematics)](https://goreportcard.com/report/github.com/a-novel/story-schematics)
[![codecov](https://codecov.io/gh/a-novel/story-schematics/graph/badge.svg?token=uc71lIIr8G)](https://codecov.io/gh/a-novel/story-schematics)

![Coverage graph](https://codecov.io/gh/a-novel/story-schematics/graphs/sunburst.svg?token=uc71lIIr8G)

<hr />

This is a quickstart document to test the project locally.

You can find the API documentation on the [repository GitHub page](https://a-novel.github.io/story-schematics/).

Want to contribute? Check the [contribution guidelines](CONTRIBUTING.md).

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

## Et Voilà!

```bash
make api
# 3:09PM INF starting application... app=story-schematics
# 3:09PM INF application started! address=:4001 app=story-schematics
```

# GOLDRUSH HighLoad Cup: The task

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Development](#development)
  - [Requirements](#requirements)
  - [Setup](#setup)
  - [Usage](#usage)
- [Deploy](#deploy)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Development

### Requirements

- Go 1.15
- [Docker](https://docs.docker.com/install/) 19.03+
- [Docker Compose](https://docs.docker.com/compose/install/) 1.25+
- Tools used to build/test project (feel free to install these tools using
  your OS package manager or any other way, but please ensure they've
  required versions):

```sh
GO111MODULE=off go get -u github.com/myitcv/gobin
curl -sSfL https://github.com/hadolint/hadolint/releases/download/v1.19.0/hadolint-$(uname)-x86_64 | sudo install /dev/stdin /usr/local/bin/hadolint
curl -sSfL https://github.com/koalaman/shellcheck/releases/download/v0.7.1/shellcheck-v0.7.1.$(uname).x86_64.tar.xz | sudo tar xJf - -C /usr/local/bin --strip-components=1 shellcheck-v0.7.1/shellcheck
```

### Setup

1. After cloning the repo copy `env.sh.dist` to `env.sh`.
2. Review `env.sh` and update for your system as needed.
3. It's recommended to add shell alias `alias dc="if test -f env.sh; then
   source env.sh; fi && docker-compose"` and then run `dc` instead of
   `docker-compose` - this way you won't have to run `source env.sh` after
   changing it.

### Usage

To develop this project you'll need only standard tools: `go generate`,
`go test`, `go build`, `docker build`. Provided scripts are for
convenience only.

- Always load `env.sh` *in every terminal* used to run any project-related
  commands (including `go test`): `source env.sh`.
    - When `env.sh.dist` change (e.g. by `git pull`) next run of `source
      env.sh` will fail and remind you to manually update `env.sh` to
      match current `env.sh.dist`.
- `go generate ./...` - do not forget to run after making changes related
  to auto-generated code
- `go test ./...` - test project (excluding integration tests), fast
- `./scripts/test` - thoroughly test project, slow
- `./scripts/test-ci-circle` - run tests locally like CircleCI will do
- `./scripts/cover` - analyse and show coverage
- `./scripts/build` - build docker image and binaries in `bin/`
    - Then use mentioned above `dc` (or `docker-compose`) to run and
      control the project.
    - Access project at host/port(s) defined in `env.sh`.

As this project isn't a real service but a _one-shot task_ which is
supposed to handle single user for a fixed period of time and then finish,
if you'll use `dc up` to start it while development then you should use
`dc up --force-recreate` to ensure each time it'll start with clean state.
An alternative is to just run `bin/task` or use `docker run`.

## Deploy

```
docker run --name=hlcup2020_task -i -t --rm \
    -e HLCUP2020_DIFFICULTY=normal \
    -v hlcup2020-task:/home/app/var/data \
    <repository>
```

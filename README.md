```
    ____  ____
   /    \/    \                             |
  |            |    ,--.  ,--.   .--  .--.  |,--. ,---.
   \::::::::::;    |   |  |   | |    |    | |   | |---'
    `:::::::;'      '--|  |--'  |     '--'  '---' `---
      `:::;'          /   |
        `'         
```

_gprobe_ is a CLI client for the
[gRPC health-checking protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

## Usage

Assuming server is listening on `localhost:1234`

Check server health (it is considered healthy if it has `grpc.health.v1.Health` service and is able to serve requests)

```bash
gprobe localhost:1234
```

Check specific service health

```bash
gprobe localhost:1234 my.package.MyService
```

Get help

```bash
gprobe -h
```

## Building

Valid _go_ environment is required to build `gprobe` (`go` is in `PATH`, `GOPATH` is set, etc.).

Build distributable tarballs for all OSes

```bash
make release
```

Build binary for current OS

```bash
make bin
```

## Development

This project follows git-flow branching model. All development is done off of the `develop` branch. `HEAD` in
`production` branch should always point to a tagged release. There's no `master` branch to avoid possible confusion.

To contribute:

1. Create a feature branch from the latest `develop`, commit your work there
    ```bash
    git checkout develop
    git pull
    git checkout -b feature/<feature_description>
    ```
2. Run `go fmt` and all the checks before committing any code
    ```bash
    go fmt ./...
    make lint test acctest
    ```
3. When the change is ready in a separate commit update `CHANGELOG.md` describing the change. Follow
[keepachangelog](http://keepachangelog.com/en/1.0.0/) guidelines
4. Create PR to develop

To release:

1. Create a release branch from the latest `develop` and update `CHANGELOG.md` there, setting version and date
    ```bash
    git checkout -b release/1.2.3
    ```
2. Create PR to `production`
3. Once PR is merged, tag `HEAD` commit using annotated tag
    ```bash
    git tag -a 1.2.3 -m "1.2.3"
    ```
4. Merge `production` back to `develop`. Do not use `fast-forward` merges
    ```bash
    git checkout develop
    git merge --no-ff production
    ```

# Contributing

Contributions are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [open source license of this project](../LICENSE).

## Submitting a pull request

1. [Fork](https://github.com/qdm12/log/fork) and clone the repository
1. Create a new branch `git checkout -b my-branch-name`
1. Modify the code
1. Ensure the docker build succeeds `docker build .`
1. Commit your modifications
1. Push to your fork and [submit a pull request](https://github.com/qdm12/log/compare)

## Resources

- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)

## Development

### Setup

#### VSCode + Docker

This is faster to setup since your entire development environment is bundled in a Docker image.

See [.devcontainer/README.md](../.devcontainer/README.md) on how to launch it in VSCode.

#### Locally

1. Install [Go](https://golang.org/dl/)
1. Install [golangci-lint](https://github.com/golangci/golangci-lint#install)
1. Install [Git](https://git-scm.com/downloads)
1. Install Go dependencies with

    ```sh
    go mod download
    ```

1. You might want to use an editor such as [Visual Studio Code](https://code.visualstudio.com/download) with the [Go extension](https://code.visualstudio.com/docs/languages/go).

### Commands

```sh
# Test the code
go test ./...
# Regenerate mocks for tests
go generate -name 'mockgen' ./...
# Lint the code
golangci-lint run
```

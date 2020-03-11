# Fuzzball Agent

[![CI Workflow](https://github.com/sylabs/fuzzball-agent/workflows/ci/badge.svg)](https://github.com/sylabs/fuzzball-agent/actions)
[![Dependabot](https://api.dependabot.com/badges/status?host=github&repo=sylabs/fuzzball-agent&identifier=238560817)](https://app.dependabot.com/accounts/sylabs/repos/238560817)

The Fuzzball Agent enables job execution for the [Fuzzball Service](https://github.com/sylabs/fuzzball-service).

## Quick Start

Ensure that you have one of the two most recent minor versions of Go installed as per the [installation instructions](https://golang.org/doc/install).

Configure your Go environment to pull private Go modules, by forcing `go get` to use `git+ssh` instead of `https`. This lets the Go compiler pull private dependencies using your machine's ssh keys.

```sh
git config --global url."ssh://git@github.com/sylabs".insteadOf "https://github.com/sylabs"
```

Starting with v1.13, the `go` command defaults to downloading modules from the public Go module mirror, and validating downloaded modules against the public Go checksum database. Since private Sylabs projects are not availble in the public mirror nor the public checksum database, we must tell Go about this. One way to do this is to set `GOPRIVATE` in the Go environment:

```sh
go env -w GOPRIVATE=github.com/sylabs
```

To run the agent, you'll need NATS endpoints to point it to. If you don't have one already, you can start one with Docker easy enough:

```sh
docker run -d -p 4222:4222 nats
```

The agent can be built and run using `go` commands like:

```sh
go run ./cmd/agent/
```


## Installation and Packaging

But, installing [Mage](https://github.com/magefile/mage) building and packaging of the agent.

Mage can be installed with:
```sh
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

This will allow you to install a binary named `fuzzball-agent` at your `$GOBIN` path:
```sh
mage install
```

Now we can build packages using:
```sh
mage rpm
```
or
```sh
mage deb
```

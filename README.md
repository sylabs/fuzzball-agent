# Fuzzball Agent

[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)
[![CI Workflow](https://github.com/sylabs/fuzzball-agent/workflows/ci/badge.svg)](https://github.com/sylabs/fuzzball-agent/actions)
[![Dependabot](https://api.dependabot.com/badges/status?host=github&repo=sylabs/fuzzball-agent&identifier=238560817)](https://app.dependabot.com/accounts/sylabs/repos/238560817)

The Fuzzball Agent enables job execution for the [Fuzzball Service](https://github.com/sylabs/fuzzball-service).

## Quick Start

Ensure that you have one of the two most recent minor versions of Go installed as per the [installation instructions](https://golang.org/doc/install).

Install [Mage](https://magefile.org) as per the [installation instructions](https://magefile.org/#installation).

To run the agent, you'll need NATS endpoints to point it to. If you don't have one already, you can start one with Docker easy enough:

```sh
docker run -d -p 4222:4222 nats
```

Finally, run the agent:

```sh
mage run
```

## Testing

Unit tests can be run like so:

```sh
mage test
```

## Installation and Packaging

To install `fuzzball-agent` in `$GOBIN`:

```sh
mage install
```

To build a `.deb` and/or `.rpm`:

```sh
mage deb
mage rpm
```

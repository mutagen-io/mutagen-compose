# Mutagen Compose

Mutagen Compose is a (minimally) modified version of
[Docker Compose](https://github.com/docker/compose) that offers automated
integration with [Mutagen](https://github.com/mutagen-io/mutagen). This allows
you to synchronize files and forward network traffic between your local system
and your Compose services. While this is primarily designed for using Compose to
develop on a remote and/or cloud-based Docker engine, it can also be used with
[Docker Desktop](https://www.docker.com/products/docker-desktop) to provide a
high-performance alternative to the virtual filesystems used in Docker Desktop
bind mounts.


## Usage

Check out the
[Mutagen Compose documentation](https://mutagen.io/documentation/orchestration/compose)
for usage information.


## System requirements

Mutagen Compose is currently supported on macOS, Linux, and Windows and is
available for
[a variety of architectures](https://github.com/mutagen-io/mutagen-compose/releases).

To use Mutagen Compose, you will need the Docker CLI (i.e the `docker` command)
installed and available in your path. You can get this via a number of
[official Docker channels](https://docs.docker.com/engine/install/), via your
system's package manager, or via a third-party package manager such as
[Homebrew](https://brew.sh/).

You will also need a matching version of Mutagen
[installed](https://mutagen.io/documentation/introduction/installation) and
available in your path.


## Installation

The best way to install Mutagen Compose is via Homebrew using:

    brew install mutagen-io/mutagen/mutagen-compose

Alternatively, you can download the
[official release binary](https://github.com/mutagen-io/mutagen-compose/releases)
and put it in your path.


## Community

The [Mutagen Community Slack Workspace](https://mutagen.io/slack) is the place
to go for discussion, questions, and ideas.

For updates about the project and its releases, you can
[follow Mutagen on Twitter](https://twitter.com/mutagen_io).


## Status

Mutagen Compose is built and tested on Windows, macOS, and Linux.

| Tests                               | Report card                         | License                                   |
| :---------------------------------: | :---------------------------------: | :---------------------------------------: |
| [![Tests][tests-badge]][tests-link] | [![Report card][rc-badge]][rc-link] | [![License][license-badge]][license-link] |

[tests-badge]: https://github.com/mutagen-io/mutagen-compose/workflows/CI/badge.svg "Test status"
[tests-link]: https://github.com/mutagen-io/mutagen-compose/actions "Test status"
[rc-badge]: https://goreportcard.com/badge/github.com/mutagen-io/mutagen-compose "Report card status"
[rc-link]: https://goreportcard.com/report/github.com/mutagen-io/mutagen-compose "Report card status"
[license-badge]: https://img.shields.io/github/license/mutagen-io/mutagen-compose.svg "MIT licensed"
[license-link]: LICENSE "MIT licensed"


## Contributing

If you'd like to contribute to Mutagen Compose, please see the
[contribution documentation](CONTRIBUTING.md).


## Security

Mutagen and its related projects take security very seriously. If you believe
you have found a security issue with Mutagen Compose, please practice
responsible disclosure practices and send an email directly to
[security@mutagen.io](mailto:security@mutagen.io) instead of opening a GitHub
issue. For more information, please see the
[security documentation](SECURITY.md).


## Building

Please see the [build instructions](BUILDING.md).

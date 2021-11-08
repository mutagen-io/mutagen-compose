# Building

Mutagen Compose relies on the Go toolchain's module support, so make sure that
you have Go module support enabled.

Mutagen Compose can be built normally using the Go toolchain, but a script is
provided to ensure a normalized build, inject version information into the
Compose package via the linker, and perform code signing on macOS. To see
information about the build script, run:

    go run scripts/build.go --help

The build script can do three different types of builds: `local` (the default -
with support for the local system only), `slim` (with support for a selection of
common platforms used in testing), and `release` (used for generating complete
release artifacts).

All artifacts from the build are placed in a `build` directory at the root of
the Mutagen Compose source tree. As a convenience, artifacts built for the
current platform are placed in the root of the build directory for easy testing,
e.g.:

    go run scripts/build.go
    build/mutagen-compose --help

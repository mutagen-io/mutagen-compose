package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	_ "github.com/docker/compose/v2/cmd/compose"

	"github.com/mutagen-io/mutagen/pkg/mutagen"
)

// composeModule is the name of the Docker Compose Go module.
const composeModule = "github.com/docker/compose/v2"

// getComposeVersion returns the underlying Docker Compose version.
func getComposeVersion() (string, error) {
	// Read build information.
	build, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("unable to read build information")
	}

	// Search for the dependency.
	for _, dependency := range build.Deps {
		if dependency.Path == composeModule {
			return dependency.Version, nil
		}
	}

	// No match was found.
	return "", errors.New("unable to find Docker Compose dependency")
}

func main() {
	// Print the program name.
	fmt.Println("Mutagen Compose")

	// Print the underlying Mutagen version.
	fmt.Println("Mutagen version:", mutagen.Version)

	// Print the underlying Docker Compose version.
	composeVersion, err := getComposeVersion()
	if err != nil {
		err = fmt.Errorf("unable to get Docker Compose version: %w", err)
		fmt.Fprintln(os.Stderr, "error: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Docker Compose version:", composeVersion)
}

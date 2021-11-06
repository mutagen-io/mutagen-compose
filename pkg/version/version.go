package version

import (
	"errors"
	"runtime/debug"
	"strings"

	"github.com/mutagen-io/mutagen/pkg/mutagen"

	// HACK: Use dummy package imports to ensure these modules are included as
	// dependencies. This isn't necessary for Mutagen Compose itself, but it is
	// necessary for the build script, which imports this package to grab
	// version information for dependencies.
	_ "github.com/docker/cli/cli/command"
	_ "github.com/docker/compose/v2/pkg/api"
)

const (
	// composeModuleName is the module name that we'll use to identify the
	// Compose dependency version.
	composeModuleName = "github.com/docker/compose/v2"
	// dockerModuleName is the module name that we'll use to identify the Docker
	// dependency version.
	dockerModuleName = "github.com/docker/cli"
)

// Versions encodes the dependency versions for Mutagen Compose. It is designed
// to be serialized as JSON.
type Versions struct {
	// Mutagen is the Mutagen version.
	Mutagen string `json:"mutagen"`
	// Compose is the Compose version.
	Compose string `json:"compose"`
	// Docker is the Docker version.
	Docker string `json:"docker"`
}

// LoadVersions loads version information.
func LoadVersions() (*Versions, error) {
	// Load build information.
	build, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, errors.New("unable to read build information")
	}

	// Create the result.
	result := &Versions{
		Mutagen: mutagen.Version,
	}

	// Attempt to identify Compose and Docker versions.
	var composeFound, dockerFound bool
	for _, dependency := range build.Deps {
		if composeFound && dockerFound {
			break
		} else if dependency.Path == composeModuleName {
			result.Compose = dependency.Version
			composeFound = true
		} else if dependency.Path == dockerModuleName {
			// HACK: The Docker CLI hasn't yet opted-in to Go modules, so
			// its version will be recorded with a +incompatible tag.
			result.Docker = strings.TrimSuffix(dependency.Version, "+incompatible")
			dockerFound = true
		}
	}

	// Fill in unknown information.
	if !composeFound {
		result.Compose = "unknown"
	}
	if !dockerFound {
		result.Docker = "unknown"
	}

	// Done.
	return result, nil
}

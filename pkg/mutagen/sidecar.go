package mutagen

import (
	"github.com/mutagen-io/mutagen/pkg/mutagen"
	"github.com/mutagen-io/mutagen/pkg/sidecar"
	"github.com/mutagen-io/mutagen/pkg/url"
)

const (
	// sidecarServiceName is the name of the Mutagen sidecar service.
	sidecarServiceName = "mutagen"
	// sidecarURLProtocol is a placeholder URL protocol used to indicate that a
	// URL should point to the Mutagen sidecar. It is used before the sidecar
	// container ID is known and will be converted to a Docker URL protocol.
	sidecarURLProtocol url.Protocol = -1
	// sidecarRoleLabelKey is the name of the label applied to the Mutagen
	// Compose sidecar container to identify it as such.
	sidecarRoleLabelKey = "io.mutagen.compose.role"
	// sidecarRoleLabelValue is the value of the label applied to the Mutagen
	// Compose sidecar container to identify it as such.
	sidecarRoleLabelValue = "sidecar"
	// sidecarVersionLabelKey is the name of the label applied to the Mutagen
	// Compose sidecar container to embed Mutagen Compose version information.
	sidecarVersionLabelKey = "io.mutagen.compose.version"
)

// sidecarImage is the full Mutagen sidecar image tag.
var sidecarImage string

func init() {
	// Compute the sidecar image tag.
	sidecarImage = sidecar.BaseTag + ":" + mutagen.Version
}

// reifySidecarURLIfNecessary converts a sidecar URL to a reified Docker URL
// with the specified Docker host and sidecar container ID. If the target URL is
// not a sidecar URL, then this function is a no-op.
func reifySidecarURLIfNecessary(target *url.URL, dockerHost, sidecarID string) {
	// If this isn't a sidecar URL, then we're done.
	if target.Protocol != sidecarURLProtocol {
		return
	}

	// Convert the protocol.
	target.Protocol = url.Protocol_Docker

	// Set the target container.
	target.Host = sidecarID

	// Set the environment.
	// TODO: Actually, we may need to end up needing to pass in Docker CLI flags
	// here after all, esp. if we're going to handle TLS overrides, but we need
	// to be mindful of locking-in implicit context usage. At the very least,
	// what we probably want to do is convert --tlsverify to DOCKER_TLS_VERIFY.
	// Unfortunately the other TLS flags require more granularity than the
	// DOCKER_CERT_PATH environment variable can provide.
	target.Environment = map[string]string{
		"DOCKER_HOST": dockerHost,
	}
}

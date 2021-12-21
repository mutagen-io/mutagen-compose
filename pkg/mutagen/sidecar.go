package mutagen

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/docker/cli/cli/command"

	"github.com/compose-spec/compose-go/types"

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
// using information from the specified Docker CLI flags, Docker CLI, and
// sidecar container ID. If the target URL is not a sidecar URL, then this
// function is a no-op.
func reifySidecarURLIfNecessary(target *url.URL, dockerFlags *pflag.FlagSet, dockerCLI command.Cli, sidecarID string) {
	// If this isn't a sidecar URL, then we're done.
	if target.Protocol != sidecarURLProtocol {
		return
	}

	// Convert the protocol.
	target.Protocol = url.Protocol_Docker

	// Set the target container.
	target.Host = sidecarID

	// Set the transport parameters so that Mutagen can reliably target the same
	// Docker daemon that Compose is currently targeting.
	//
	// There are two possible modes that we need to consider: host-based and
	// context-based. The most reliable way to determine which mode we're in is
	// to inspect the currently selected context. If this context is "default",
	// then the host-based mode (i.e. the mode determined by command line flags
	// and environment variables) is being used. Note that there's a difference
	// between the "default" context and the active context (the latter of which
	// is a term that's not actually used by the context documentation at the
	// time of writing). The context named "default" always exists and indicates
	// that the host-based mode is being used. The "default" context may be the
	// active context, or it may not be. Unfortunately the Docker CLI help
	// information for contexts is confusing, because the "docker context use"
	// command indicates that it sets the "default" context, but it's actually
	// setting the active context that will be used by default if no host-based
	// information is provided.
	//
	// If we're using the context-based mode, then we just set the context and
	// config parameters, because command-line-based and environment-based TLS
	// settings aren't used in that case (TLS information is always sourced from
	// the context). Obviously we're only pinning on the context name here, and
	// in theory that could switch to point elsewhere, but it's at least as
	// stable as an SSH hostname or a Unix domain socket target.
	//
	// If we're using the host-based mode, then we just set the host and
	// TLS-related parameters (config is only necessary to look up context
	// information), and we can just pull TLS settings from the CLI flags,
	// because they will already incorporate their respective environment
	// variables as default settings.
	//
	// The logic here is designed to minimize the number of parameters we set to
	// accomplish correct targeting. It is largely guided by the implementation
	// of the Docker CLI's CommonOptions type, specifically its InstallFlags
	// method and the way it uses environment variables for defaults.
	usingNonDefaultConfigPath := os.Getenv("DOCKER_CONFIG") != ""
	if dockerCLI.CurrentContext() == command.DefaultContextName {
		target.Parameters = map[string]string{
			"host": dockerCLI.Client().DaemonHost(),
		}
		var tlsInUse bool
		if dockerFlags.Lookup("tls").Value.String() == "true" {
			target.Parameters["tls"] = ""
			tlsInUse = true
		}
		if dockerFlags.Lookup("tlsverify").Value.String() == "true" {
			target.Parameters["tlsverify"] = ""
			tlsInUse = true
		}
		// HACK: Technically we should validate that the following flags are
		// non-empty, otherwise Mutagen will reject them, but the fact that the
		// Docker API client has already connected should effectively validate
		// that they're non-empty.
		usingNonDefaultCertPath := os.Getenv("DOCKER_CERT_PATH") != ""
		tlsCACertFlag := dockerFlags.Lookup("tlscacert")
		if tlsInUse && (usingNonDefaultCertPath || usingNonDefaultConfigPath || tlsCACertFlag.Changed) {
			target.Parameters["tlscacert"] = tlsCACertFlag.Value.String()
		}
		tlsCertFlag := dockerFlags.Lookup("tlscert")
		if tlsInUse && (usingNonDefaultCertPath || usingNonDefaultConfigPath || tlsCertFlag.Changed) {
			target.Parameters["tlscert"] = tlsCertFlag.Value.String()
		}
		tlsKeyFlag := dockerFlags.Lookup("tlskey")
		if tlsInUse && (usingNonDefaultCertPath || usingNonDefaultConfigPath || tlsKeyFlag.Changed) {
			target.Parameters["tlskey"] = tlsKeyFlag.Value.String()
		}
	} else {
		target.Parameters = map[string]string{
			"context": dockerCLI.CurrentContext(),
		}
		configFlag := dockerFlags.Lookup("config")
		if usingNonDefaultConfigPath || configFlag.Changed {
			target.Parameters["config"] = configFlag.Value.String()
		}
	}
}

// isValidRestartPolicy returns true if and only if the provided restart policy
// is non-empty and names a valid restart policy.
func isValidRestartPolicy(restart string) bool {
	return restart == types.RestartPolicyAlways ||
		restart == types.RestartPolicyOnFailure ||
		restart == types.RestartPolicyNo ||
		restart == types.RestartPolicyUnlessStopped
}

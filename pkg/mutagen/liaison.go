package mutagen

import (
	"fmt"

	"github.com/docker/cli/cli/command"

	"github.com/docker/docker/client"

	"github.com/compose-spec/compose-go/types"

	"github.com/docker/compose/v2/pkg/api"

	"github.com/mitchellh/mapstructure"
)

// Liaison is the interface point between Compose and Mutagen. Its zero value is
// initialized and ready to use. It implements the Compose service API. It is a
// single-use type, is not safe for concurrent usage, and its Shutdown method
// should be invoked when usage is complete.
type Liaison struct {
	// dockerCLI is the associated Docker CLI instance.
	dockerCLI command.Cli
	// composeService is the underlying Compose service.
	composeService api.Service
	// configuration is the Mutagen configuration loaded from the x-mutagen
	// extensions of the Compose project. If it's nil after invoking the
	// loadConfiguration method, then no x-mutagen extensions were found.
	configuration *configuration
}

// Shutdown terminates liaison resource usage.
func (l *Liaison) Shutdown() error {
	// TODO: Implement. We'll want to terminate any active Mutagen gRPC clients.
	return nil
}

// RegisterDockerCLI registers the associated Docker CLI instance.
func (l *Liaison) RegisterDockerCLI(cli command.Cli) {
	l.dockerCLI = cli
}

// DockerClient returns a Mutagen-aware version of the Docker API client. It
// must only be called after a Docker CLI is registered with RegisterDockerCLI
// and said CLI can return a valid API client via its Client method (typically
// after flag parsing).
func (l *Liaison) DockerClient() client.APIClient {
	return &dockerAPIClient{l, l.dockerCLI.Client()}
}

// RegisterComposeService registers the underlying Compose service. The Compose
// service must be initialized using the Docker API client returned by the
// liaison's DockerClient method.
func (l *Liaison) RegisterComposeService(service api.Service) {
	l.composeService = service
}

// ComposeService returns a Mutagen-aware version of the Compose Service API. It
// must be called only after a Compose service has been registered with
// RegisterDockerCLI.
func (l *Liaison) ComposeService() api.Service {
	return &composeService{l, l.composeService}
}

// processProject loads Mutagen configuration from the specified project, adds
// the Mutagen sidecar service to the project, and sets project dependencies
// accordingly. If project is nil, this method is a no-op and returns nil.
func (l *Liaison) processProject(project *types.Project) error {
	// If the project is nil, then there's nothing to process.
	if project == nil {
		return nil
	}

	// Grab the Mutagen extension section. If it's not present, then there's
	// nothing to load.
	// TODO: Do we want to create a nil configuraiton and still inject a sidecar
	// in this case?
	xMutagen, ok := project.Extensions["x-mutagen"]
	if !ok {
		return nil
	}

	// Decode the extension section.
	l.configuration = &configuration{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			boolToIgnoreVCSModeHookFunc(),
		),
		ErrorUnused: true,
		Result:      l.configuration,
	})
	if err != nil {
		return fmt.Errorf("unable to create configuration decoder: %w", err)
	} else if err = decoder.Decode(xMutagen); err != nil {
		return fmt.Errorf("unable to decode x-mutagen section: %w", err)
	}

	// TODO: Implement configuration parsing, session specification generation,
	// service injection, and dependency injection for the project.
	project.Services = append(project.Services, types.ServiceConfig{
		Name:  "mutagen",
		Image: sidecarImage,
	})
	return nil
}

// reconcileSessions performs Mutagen session reconciliation for the project
// using the specified sidecar container ID as the target identifier. It also
// ensures that all sessions are unpaused.
func (l *Liaison) reconcileSessions(sidecarID string) error {
	// TODO: Implement.
	fmt.Println("Reconciling Mutagen sessions for", sidecarID)
	return nil
}

// listSessions lists Mutagen sessions for the project using the specified
// sidecar container ID as the target identifier.
func (l *Liaison) listSessions(sidecarID string) error {
	// TODO: Implement.
	fmt.Println("Listing Mutagen sessions for", sidecarID)
	return nil
}

// pauseSessions pauses Mutagen sessions for the project using the specified
// sidecar container ID as the target identifier.
func (l *Liaison) pauseSessions(sidecarID string) error {
	// TODO: Implement.
	fmt.Println("Pausing Mutagen sessions for", sidecarID)
	return nil
}

// resumeSessions resumes Mutagen sessions for the project using the specified
// sidecar container ID as the target identifier.
func (l *Liaison) resumeSessions(sidecarID string) error {
	// TODO: Implement.
	fmt.Println("Resuming Mutagen sessions for", sidecarID)
	return nil
}

// terminateSessions terminates Mutagen sessions for the project using the
// specified sidecar container ID as the target identifier.
func (l *Liaison) terminateSessions(sidecarID string) error {
	// TODO: Implement.
	fmt.Println("Terminating Mutagen sessions for", sidecarID)
	return nil
}

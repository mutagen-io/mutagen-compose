package mutagen

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli/command"

	"github.com/docker/compose/v2/pkg/api"

	"github.com/compose-spec/compose-go/types"

	"github.com/mutagen-io/mutagen-compose/pkg/docker"

	"github.com/mitchellh/mapstructure"
)

// Liaison is the interface point between Compose and Mutagen. Its zero value is
// initialized and ready to use. It implements the Compose service API. It is a
// single-use type, is not safe for concurrent usage, and its Shutdown method
// should be invoked when usage is complete.
type Liaison struct {
	// dockerFlags are the associated top-level Docker CLI flags.
	dockerFlags *docker.Flags
	// dockerCLI is the associated Docker CLI instance.
	dockerCLI command.Cli
	// Service is the underlying Compose service.
	api.Service
	// configuration is the Mutagen configuration loaded from the x-mutagen
	// extensions of the Compose project. If it's nil after invoking the
	// loadConfiguration method, then no x-mutagen extensions were found.
	configuration *configuration
}

// Shutdown terminates liaison resource usage.
func (l *Liaison) Shutdown() error {
	// TODO: Implement. We'll want to terminate any active Mutagen clients.
	return nil
}

// RegisterDockerCLIFlags registers the associated top-level Docker CLI flags.
// It must be called before the liaison is used as a Compose service.
func (l *Liaison) RegisterDockerCLIFlags(dockerFlags *docker.Flags) {
	l.dockerFlags = dockerFlags
}

// RegisterDockerCLI registers the associated Docker CLI instance. It must be
// called before the liaison is used as a Compose service.
func (l *Liaison) RegisterDockerCLI(dockerCLI command.Cli) {
	l.dockerCLI = dockerCLI
}

// RegisterComposeService registers the underlying Compose service. It must be
// called before the liaison is used as a Compose service.
func (l *Liaison) RegisterComposeService(service api.Service) {
	l.Service = service
}

// loadConfiguration loads the Mutagen configuration from the x-mutagen
// extensions in a Compose project. It sets the liaison's configuration field
// according to what is (or isn't) found.
func (l *Liaison) loadConfiguration(project *types.Project) error {
	// Grab the Mutagen extension section. If it's not present, then there's
	// nothing to load.
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

	// Success.
	return nil
}

// Up implements github.com/docker/compose/v2/pkg/api.Service.Up.
func (l *Liaison) Up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	// Process Mutagen extensions. If none are present, then just dispatch
	// directly to the underlying Compose service.
	if err := l.loadConfiguration(project); err != nil {
		return fmt.Errorf("unable to load Mutagen configuration: %w", err)
	} else if l.configuration == nil {
		return l.Service.Up(ctx, project, options)
	}

	// TODO: Handle Mutagen-based operation.
	fmt.Println("Mutagen-extended up not yet implemented")
	return nil
}

// Start implements github.com/docker/compose/v2/pkg/api.Service.Start.
func (l *Liaison) Start(ctx context.Context, project *types.Project, options api.StartOptions) error {
	// Process Mutagen extensions. If none are present, then just dispatch
	// directly to the underlying Compose service.
	if err := l.loadConfiguration(project); err != nil {
		return fmt.Errorf("unable to load Mutagen configuration: %w", err)
	} else if l.configuration == nil {
		return l.Service.Start(ctx, project, options)
	}

	// TODO: Handle Mutagen-based operation.
	fmt.Println("Mutagen-extended start not yet implemented")
	return nil
}

// RunOneOffContainer implements
// github.com/docker/compose/v2/pkg/api.Service.RunOneOffContainer.
func (l *Liaison) RunOneOffContainer(ctx context.Context, project *types.Project, options api.RunOptions) (int, error) {
	// Process Mutagen extensions. If none are present, then just dispatch
	// directly to the underlying Compose service.
	if err := l.loadConfiguration(project); err != nil {
		return 0, fmt.Errorf("unable to load Mutagen configuration: %w", err)
	} else if l.configuration == nil {
		return l.Service.RunOneOffContainer(ctx, project, options)
	}

	// TODO: Handle Mutagen-based operation.
	fmt.Println("Mutagen-extended run not yet implemented")
	return 0, nil
}

// Ps implements github.com/docker/compose/v2/pkg/api.Service.Ps.
func (l *Liaison) Ps(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
	// TODO: Use the project name and daemon client to identify the Mutagen
	// sidecar ID and invoke session listing.

	// Dispatch directly to the underlying Compose service.
	return l.Service.Ps(ctx, projectName, options)
}

// Stop implements github.com/docker/compose/v2/pkg/api.Service.Stop.
func (l *Liaison) Stop(ctx context.Context, project *types.Project, options api.StopOptions) error {
	// Process Mutagen extensions. If none are present, then just dispatch
	// directly to the underlying Compose service.
	if err := l.loadConfiguration(project); err != nil {
		return fmt.Errorf("unable to load Mutagen configuration: %w", err)
	} else if l.configuration == nil {
		return l.Service.Stop(ctx, project, options)
	}

	// TODO: Handle Mutagen-based operation.
	fmt.Println("Mutagen-extended stop not yet implemented")
	return nil
}

// Down implements github.com/docker/compose/v2/pkg/api.Service.Down.
func (l *Liaison) Down(ctx context.Context, projectName string, options api.DownOptions) error {
	// TODO: Use the project name and daemon client to identify the Mutagen
	// sidecar ID and invoke session termination, but only if no services have
	// been explicitly specified.

	// TODO: Figure out how Down is operating if options.Project is nil. How
	// does it know which services to take down in that case? Is it just doing a
	// filter based on projectName?
	return l.Service.Down(ctx, projectName, options)
}

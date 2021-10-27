package main

import (
	"context"
	"fmt"
	"os"

	"github.com/compose-spec/compose-go/types"

	"github.com/docker/compose/v2/pkg/api"
)

// needsLiaison maps Compose command names to a boolean indicating whether or
// not they need a Mutagen liaison.
var needsLiaison = map[string]bool{
	"up":    true,
	"start": true,
	"run":   true,
	"ps":    true,
	"stop":  true,
	"down":  true,
}

// liaison is the interface point between Mutagen and Compose.
type liaison struct {
	// compose is the underlying Compose implementation.
	compose api.Service
	// TODO: Add a Mutagen daemon client.
}

// newLiaison creates a new liaison.
func newLiaison() (*liaison, error) {
	// TODO: Add creation of Mutagen daemon client.
	return &liaison{}, nil
}

// wrap wraps and installs Mutagen hooks into a Compose service. A liaison may
// wrap only a single Compose service. A nil liaison returns the Compose service
// unwrapped and unmodified.
func (l *liaison) wrap(compose api.Service) api.Service {
	// If the liaison is nil, then we don't wrap the service.
	if l == nil {
		return compose
	}

	// Store the underlying service.
	l.compose = compose

	// Create an API proxy.
	proxy := api.NewServiceProxy().WithService(compose)

	// Install overrides.
	proxy.UpFn = l.up
	proxy.StartFn = l.start
	proxy.RunOneOffContainerFn = l.runOneOffContainer
	proxy.PsFn = l.ps
	proxy.StopFn = l.stop
	proxy.DownFn = l.down

	// Done.
	return proxy
}

// up implements github.com/docker/compose/v2/pkg/api.Service.Up.
func (l *liaison) up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted up")
	return l.compose.Up(ctx, project, options)
}

// start implements github.com/docker/compose/v2/pkg/api.Service.Start.
func (l *liaison) start(ctx context.Context, project *types.Project, options api.StartOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted start")
	return l.compose.Start(ctx, project, options)
}

// runOneOffContainer implements
// github.com/docker/compose/v2/pkg/api.Service.RunOneOffContainer.
func (l *liaison) runOneOffContainer(ctx context.Context, project *types.Project, options api.RunOptions) (int, error) {
	fmt.Fprintln(os.Stderr, "Intercepted run")
	return l.compose.RunOneOffContainer(ctx, project, options)
}

// ps implements github.com/docker/compose/v2/pkg/api.Service.Ps.
func (l *liaison) ps(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
	fmt.Fprintln(os.Stderr, "Intercepted ps")
	return l.compose.Ps(ctx, projectName, options)
}

// stop implements github.com/docker/compose/v2/pkg/api.Service.Stop.
func (l *liaison) stop(ctx context.Context, project *types.Project, options api.StopOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted stop")
	return l.compose.Stop(ctx, project, options)
}

// down implements github.com/docker/compose/v2/pkg/api.Service.Down.
func (l *liaison) down(ctx context.Context, projectName string, options api.DownOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted down")
	return l.compose.Down(ctx, projectName, options)
}

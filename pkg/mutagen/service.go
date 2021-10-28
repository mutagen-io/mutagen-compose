package mutagen

import (
	"context"
	"fmt"
	"os"

	"github.com/compose-spec/compose-go/types"

	"github.com/docker/compose/v2/pkg/api"

	"google.golang.org/grpc"
)

// service implements the Compose Service interface with Mutagen enhancements.
type service struct {
	// Service is the underlying Compose service.
	api.Service
	// client is the Mutagen daemon client.
	client *grpc.ClientConn
	// TODO: Add Mutagen service clients.
}

// Wrap wraps an existing Compose service and adds Mutagen enhancements.
func Wrap(compose api.Service) api.Service {
	return &service{
		Service: compose,
	}
}

// TODO: Add function to establish Mutagen connection on-demand.

// Up implements github.com/docker/compose/v2/pkg/api.Service.Up.
func (s *service) Up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted up")
	return s.Service.Up(ctx, project, options)
}

// Start implements github.com/docker/compose/v2/pkg/api.Service.Start.
func (s *service) Start(ctx context.Context, project *types.Project, options api.StartOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted start")
	return s.Service.Start(ctx, project, options)
}

// RunOneOffContainer implements
// github.com/docker/compose/v2/pkg/api.Service.RunOneOffContainer.
func (s *service) RunOneOffContainer(ctx context.Context, project *types.Project, options api.RunOptions) (int, error) {
	fmt.Fprintln(os.Stderr, "Intercepted run")
	return s.Service.RunOneOffContainer(ctx, project, options)
}

// Ps implements github.com/docker/compose/v2/pkg/api.Service.Ps.
func (s *service) Ps(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
	fmt.Fprintln(os.Stderr, "Intercepted ps")
	return s.Service.Ps(ctx, projectName, options)
}

// Stop implements github.com/docker/compose/v2/pkg/api.Service.Stop.
func (s *service) Stop(ctx context.Context, project *types.Project, options api.StopOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted stop")
	return s.Service.Stop(ctx, project, options)
}

// Down implements github.com/docker/compose/v2/pkg/api.Service.Down.
func (s *service) Down(ctx context.Context, projectName string, options api.DownOptions) error {
	fmt.Fprintln(os.Stderr, "Intercepted down")
	return s.Service.Down(ctx, projectName, options)
}

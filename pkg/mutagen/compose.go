package mutagen

import (
	"context"
	"fmt"

	"github.com/compose-spec/compose-go/types"

	"github.com/docker/compose/v2/pkg/api"
)

// composeService is a Mutagen-aware implementation of
// github.com/docker/compose/v2/pkg/api.Service that injects Mutagen services
// and dependencies into the project.
type composeService struct {
	// liaison is the parent Mutagen liaison.
	liaison *Liaison
	// service is the underlying Compose service.
	service api.Service
}

// Build implements github.com/docker/compose/v2/pkg/api.Service.Build.
func (s *composeService) Build(ctx context.Context, project *types.Project, options api.BuildOptions) error {
	return s.service.Build(ctx, project, options)
}

// Push implements github.com/docker/compose/v2/pkg/api.Service.Push.
func (s *composeService) Push(ctx context.Context, project *types.Project, options api.PushOptions) error {
	return s.service.Push(ctx, project, options)
}

// Pull implements github.com/docker/compose/v2/pkg/api.Service.Pull.
func (s *composeService) Pull(ctx context.Context, project *types.Project, options api.PullOptions) error {
	// Process Mutagen extensions for the project.
	if err := s.liaison.processProject(project); err != nil {
		return fmt.Errorf("unable to process project: %w", err)
	}

	// Invoke the underlying implementation.
	return s.service.Pull(ctx, project, options)
}

// Create implements github.com/docker/compose/v2/pkg/api.Service.Create.
func (s *composeService) Create(ctx context.Context, project *types.Project, options api.CreateOptions) error {
	// Process Mutagen extensions for the project.
	if err := s.liaison.processProject(project); err != nil {
		return fmt.Errorf("unable to process project: %w", err)
	}

	// Invoke the underlying implementation.
	return s.service.Create(ctx, project, options)
}

// Start implements github.com/docker/compose/v2/pkg/api.Service.Start.
func (s *composeService) Start(ctx context.Context, project *types.Project, options api.StartOptions) error {
	// Process Mutagen extensions for the project.
	if err := s.liaison.processProject(project); err != nil {
		return fmt.Errorf("unable to process project: %w", err)
	}

	// Invoke the underlying implementation.
	return s.service.Start(ctx, project, options)
}

// Restart implements github.com/docker/compose/v2/pkg/api.Service.Restart.
func (s *composeService) Restart(ctx context.Context, project *types.Project, options api.RestartOptions) error {
	return s.service.Restart(ctx, project, options)
}

// Stop implements github.com/docker/compose/v2/pkg/api.Service.Stop.
func (s *composeService) Stop(ctx context.Context, project *types.Project, options api.StopOptions) error {
	// Process Mutagen extensions for the project.
	if err := s.liaison.processProject(project); err != nil {
		return fmt.Errorf("unable to process project: %w", err)
	}

	// Invoke the underlying implementation.
	return s.service.Stop(ctx, project, options)
}

// Up implements github.com/docker/compose/v2/pkg/api.Service.Up.
func (s *composeService) Up(ctx context.Context, project *types.Project, options api.UpOptions) error {
	// Process Mutagen extensions for the project.
	if err := s.liaison.processProject(project); err != nil {
		return fmt.Errorf("unable to process project: %w", err)
	}

	// Invoke the underlying implementation.
	return s.service.Up(ctx, project, options)
}

// Down implements github.com/docker/compose/v2/pkg/api.Service.Down.
func (s *composeService) Down(ctx context.Context, projectName string, options api.DownOptions) error {
	// Process Mutagen extensions for the project.
	if err := s.liaison.processProject(options.Project); err != nil {
		return fmt.Errorf("unable to process project: %w", err)
	}

	// Invoke the underlying implementation.
	return s.service.Down(ctx, projectName, options)
}

// Logs implements github.com/docker/compose/v2/pkg/api.Service.Logs.
func (s *composeService) Logs(ctx context.Context, projectName string, consumer api.LogConsumer, options api.LogOptions) error {
	return s.service.Logs(ctx, projectName, consumer, options)
}

// Ps implements github.com/docker/compose/v2/pkg/api.Service.Ps.
func (s *composeService) Ps(ctx context.Context, projectName string, options api.PsOptions) ([]api.ContainerSummary, error) {
	// TODO: Get the Mutagen Compose sidecar ID and invoke s.liaison.listSessions.
	return s.service.Ps(ctx, projectName, options)
}

// List implements github.com/docker/compose/v2/pkg/api.Service.List.
func (s *composeService) List(ctx context.Context, options api.ListOptions) ([]api.Stack, error) {
	return s.service.List(ctx, options)
}

// Convert implements github.com/docker/compose/v2/pkg/api.Service.Convert.
func (s *composeService) Convert(ctx context.Context, project *types.Project, options api.ConvertOptions) ([]byte, error) {
	return s.service.Convert(ctx, project, options)
}

// Kill implements github.com/docker/compose/v2/pkg/api.Service.Kill.
func (s *composeService) Kill(ctx context.Context, project *types.Project, options api.KillOptions) error {
	return s.service.Kill(ctx, project, options)
}

// RunOneOffContainer implements
// github.com/docker/compose/v2/pkg/api.Service.RunOneOffContainer.
func (s *composeService) RunOneOffContainer(ctx context.Context, project *types.Project, options api.RunOptions) (int, error) {
	// Process Mutagen extensions for the project.
	if err := s.liaison.processProject(project); err != nil {
		return 0, fmt.Errorf("unable to process project: %w", err)
	}

	// Invoke the underlying implementation.
	return s.service.RunOneOffContainer(ctx, project, options)
}

// Remove implements github.com/docker/compose/v2/pkg/api.Service.Remove.
func (s *composeService) Remove(ctx context.Context, project *types.Project, options api.RemoveOptions) error {
	return s.service.Remove(ctx, project, options)
}

// Exec implements github.com/docker/compose/v2/pkg/api.Service.Exec.
func (s *composeService) Exec(ctx context.Context, projectName string, options api.RunOptions) (int, error) {
	return s.service.Exec(ctx, projectName, options)
}

// Copy implements github.com/docker/compose/v2/pkg/api.Service.Copy.
func (s *composeService) Copy(ctx context.Context, project *types.Project, options api.CopyOptions) error {
	return s.service.Copy(ctx, project, options)
}

// Pause implements github.com/docker/compose/v2/pkg/api.Service.Pause.
func (s *composeService) Pause(ctx context.Context, projectName string, options api.PauseOptions) error {
	return s.service.Pause(ctx, projectName, options)
}

// UnPause implements github.com/docker/compose/v2/pkg/api.Service.UnPause.
func (s *composeService) UnPause(ctx context.Context, projectName string, options api.PauseOptions) error {
	return s.service.UnPause(ctx, projectName, options)
}

// Top implements github.com/docker/compose/v2/pkg/api.Service.Top.
func (s *composeService) Top(ctx context.Context, projectName string, services []string) ([]api.ContainerProcSummary, error) {
	return s.service.Top(ctx, projectName, services)
}

// Events implements github.com/docker/compose/v2/pkg/api.Service.Events.
func (s *composeService) Events(ctx context.Context, projectName string, options api.EventsOptions) error {
	return s.service.Events(ctx, projectName, options)
}

// Port implements github.com/docker/compose/v2/pkg/api.Service.Port.
func (s *composeService) Port(ctx context.Context, projectName string, service string, port int, options api.PortOptions) (string, int, error) {
	return s.service.Port(ctx, projectName, service, port, options)
}

// Images implements github.com/docker/compose/v2/pkg/api.Service.Images.
func (s *composeService) Images(ctx context.Context, projectName string, options api.ImagesOptions) ([]api.ImageSummary, error) {
	return s.service.Images(ctx, projectName, options)
}

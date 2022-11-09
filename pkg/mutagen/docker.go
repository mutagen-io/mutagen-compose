package mutagen

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	containerAPITypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// dockerAPIClient is a Mutagen-aware implementation of
// github.com/docker/docker/client.APIClient that performs special handling of
// Mutagen Compose sidecar containers.
type dockerAPIClient struct {
	// liaison is the parent Mutagen liaison.
	liaison *Liaison
	// APIClient is the underlying API client.
	client.APIClient
}

// isMutagenComposeSidecar checks if the specified container is a Mutagen
// Compose sidecar container. In the context of a Compose project, this is
// equivalent to the container being the only Mutagen Compose sidecar container
// for that project.
func (c *dockerAPIClient) isMutagenComposeSidecar(ctx context.Context, container string) (bool, error) {
	// Grab the container metadata.
	metadata, err := c.APIClient.ContainerInspect(ctx, container)
	if err != nil {
		return false, fmt.Errorf("unable to inspect container: %w", err)
	}

	// Check if this is a Mutagen Compose sidecar container.
	return metadata.Config.Labels[sidecarRoleLabelKey] == sidecarRoleLabelValue, nil
}

// ContainerStart implements
// github.com/docker/docker/client.APIClient.ContainerStart.
func (c *dockerAPIClient) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	// Start the container.
	if err := c.APIClient.ContainerStart(ctx, container, options); err != nil {
		return err
	}

	// If this is a Mutagen compose sidecar container, then either reconcile
	// Mutagen sessions or just resume them, depending on whether or not we have
	// session definitions from the project.
	if sidecar, err := c.isMutagenComposeSidecar(ctx, container); err != nil {
		return fmt.Errorf("unable to determine if container is sidecar: %w", err)
	} else if sidecar {
		if c.liaison.processedProject {
			if err := c.liaison.reconcileSessions(ctx, container); err != nil {
				return fmt.Errorf("unable to reconcile Mutagen sessions: %w", err)
			}
		} else {
			if err := c.liaison.resumeSessions(ctx, container); err != nil {
				return fmt.Errorf("unable to resume Mutagen sessions: %w", err)
			}
		}
	}

	// Success.
	return nil
}

// ContainerPause implements
// github.com/docker/docker/client.APIClient.ContainerPause.
func (c *dockerAPIClient) ContainerPause(ctx context.Context, container string) error {
	// If this is a Mutagen compose sidecar container, then pause associated
	// Mutagen sessions.
	if sidecar, err := c.isMutagenComposeSidecar(ctx, container); err != nil {
		return fmt.Errorf("unable to determine if container is sidecar: %w", err)
	} else if sidecar {
		if err := c.liaison.pauseSessions(ctx, container); err != nil {
			return fmt.Errorf("unable to pause Mutagen sessions: %w", err)
		}
	}

	// Pause the container.
	return c.APIClient.ContainerPause(ctx, container)
}

// ContainerUnpause implements
// github.com/docker/docker/client.APIClient.ContainerUnpause.
func (c *dockerAPIClient) ContainerUnpause(ctx context.Context, container string) error {
	// Unpause the container.
	if err := c.APIClient.ContainerUnpause(ctx, container); err != nil {
		return err
	}

	// If this is a Mutagen compose sidecar container, then resume associated
	// Mutagen sessions.
	if sidecar, err := c.isMutagenComposeSidecar(ctx, container); err != nil {
		return fmt.Errorf("unable to determine if container is sidecar: %w", err)
	} else if sidecar {
		if err := c.liaison.resumeSessions(ctx, container); err != nil {
			return fmt.Errorf("unable to resume Mutagen sessions: %w", err)
		}
	}

	// Success.
	return nil
}

// ContainerStop implements
// github.com/docker/docker/client.APIClient.ContainerStop.
func (c *dockerAPIClient) ContainerStop(ctx context.Context, container string, options containerAPITypes.StopOptions) error {
	// If this is a Mutagen compose sidecar container, then pause associated
	// Mutagen sessions.
	if sidecar, err := c.isMutagenComposeSidecar(ctx, container); err != nil {
		return fmt.Errorf("unable to determine if container is sidecar: %w", err)
	} else if sidecar {
		if err := c.liaison.pauseSessions(ctx, container); err != nil {
			return fmt.Errorf("unable to pause Mutagen sessions: %w", err)
		}
	}

	// Stop the container.
	return c.APIClient.ContainerStop(ctx, container, options)
}

// ContainerRemove implements
// github.com/docker/docker/client.APIClient.ContainerRemove.
func (c *dockerAPIClient) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	// If this is a Mutagen compose sidecar container, then terminate associated
	// Mutagen sessions.
	if sidecar, err := c.isMutagenComposeSidecar(ctx, container); err != nil {
		return fmt.Errorf("unable to determine if container is sidecar: %w", err)
	} else if sidecar {
		if err := c.liaison.terminateSessions(ctx, container); err != nil {
			return fmt.Errorf("unable to terminate Mutagen sessions: %w", err)
		}
	}

	// Remove the container.
	return c.APIClient.ContainerRemove(ctx, container, options)
}

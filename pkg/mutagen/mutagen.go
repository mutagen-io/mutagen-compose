package mutagen

import (
	"context"
	"fmt"

	"github.com/mutagen-io/mutagen/pkg/grpcutil"
	"github.com/mutagen-io/mutagen/pkg/selection"
	forwardingsvc "github.com/mutagen-io/mutagen/pkg/service/forwarding"
	synchronizationsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
)

// forwardCreateWithSpecification creates a forwarding session using the
// provided forwarding service client, session specification, and prompter.
func forwardCreateWithSpecification(
	ctx context.Context,
	forwardingService forwardingsvc.ForwardingClient,
	prompter string,
	specification *forwardingsvc.CreationSpecification,
) (string, error) {
	response, err := forwardingService.Create(ctx, &forwardingsvc.CreateRequest{
		Prompter:      prompter,
		Specification: specification,
	})
	if err != nil {
		return "", grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return "", fmt.Errorf("invalid create response received: %w", err)
	}
	return response.Session, nil
}

// syncCreateWithSpecification creates a synchronization session using the
// provided synchronization service client, session specification, and prompter.
func syncCreateWithSpecification(
	ctx context.Context,
	synchronizationService synchronizationsvc.SynchronizationClient,
	prompter string,
	specification *synchronizationsvc.CreationSpecification,
) (string, error) {
	response, err := synchronizationService.Create(ctx, &synchronizationsvc.CreateRequest{
		Prompter:      prompter,
		Specification: specification,
	})
	if err != nil {
		return "", grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return "", fmt.Errorf("invalid create response received: %w", err)
	}
	return response.Session, nil
}

// syncFlushWithSelection flushes synchronization sessions using the provided
// synchronization service client, session selection, and prompter.
func syncFlushWithSelection(
	ctx context.Context,
	synchronizationService synchronizationsvc.SynchronizationClient,
	prompter string,
	selection *selection.Selection,
) error {
	response, err := synchronizationService.Flush(ctx, &synchronizationsvc.FlushRequest{
		Prompter:  prompter,
		Selection: selection,
	})
	if err != nil {
		return grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return fmt.Errorf("invalid flush response received: %w", err)
	}
	return nil
}

// forwardPauseWithSelection pauses forwarding sessions using the provided
// forwarding service client, session selection, and prompter.
func forwardPauseWithSelection(
	ctx context.Context,
	forwardingService forwardingsvc.ForwardingClient,
	prompter string,
	selection *selection.Selection,
) error {
	response, err := forwardingService.Pause(ctx, &forwardingsvc.PauseRequest{
		Prompter:  prompter,
		Selection: selection,
	})
	if err != nil {
		return grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return fmt.Errorf("invalid pause response received: %w", err)
	}
	return nil
}

// syncPauseWithSelection pauses synchronization sessions using the provided
// synchronization service client, session selection, and prompter.
func syncPauseWithSelection(
	ctx context.Context,
	synchronizationService synchronizationsvc.SynchronizationClient,
	prompter string,
	selection *selection.Selection,
) error {
	response, err := synchronizationService.Pause(ctx, &synchronizationsvc.PauseRequest{
		Prompter:  prompter,
		Selection: selection,
	})
	if err != nil {
		return grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return fmt.Errorf("invalid pause response received: %w", err)
	}
	return nil
}

// forwardResumeWithSelection resumes forwarding sessions using the provided
// forwarding service client, session selection, and prompter.
func forwardResumeWithSelection(
	ctx context.Context,
	forwardingService forwardingsvc.ForwardingClient,
	prompter string,
	selection *selection.Selection,
) error {
	response, err := forwardingService.Resume(ctx, &forwardingsvc.ResumeRequest{
		Prompter:  prompter,
		Selection: selection,
	})
	if err != nil {
		return grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return fmt.Errorf("invalid resume response received: %w", err)
	}
	return nil
}

// syncResumeWithSelection resumes synchronization sessions using the provided
// synchronization service client, session selection, and prompter.
func syncResumeWithSelection(
	ctx context.Context,
	synchronizationService synchronizationsvc.SynchronizationClient,
	prompter string,
	selection *selection.Selection,
) error {
	response, err := synchronizationService.Resume(ctx, &synchronizationsvc.ResumeRequest{
		Prompter:  prompter,
		Selection: selection,
	})
	if err != nil {
		return grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return fmt.Errorf("invalid resume response received: %w", err)
	}
	return nil
}

// forwardTerminateWithSelection terminates forwarding sessions using the provided
// forwarding service client, session selection, and prompter.
func forwardTerminateWithSelection(
	ctx context.Context,
	forwardingService forwardingsvc.ForwardingClient,
	prompter string,
	selection *selection.Selection,
) error {
	response, err := forwardingService.Terminate(ctx, &forwardingsvc.TerminateRequest{
		Prompter:  prompter,
		Selection: selection,
	})
	if err != nil {
		return grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return fmt.Errorf("invalid terminate response received: %w", err)
	}
	return nil
}

// syncTerminateWithSelection terminates synchronization sessions using the provided
// synchronization service client, session selection, and prompter.
func syncTerminateWithSelection(
	ctx context.Context,
	synchronizationService synchronizationsvc.SynchronizationClient,
	prompter string,
	selection *selection.Selection,
) error {
	response, err := synchronizationService.Terminate(ctx, &synchronizationsvc.TerminateRequest{
		Prompter:  prompter,
		Selection: selection,
	})
	if err != nil {
		return grpcutil.PeelAwayRPCErrorLayer(err)
	} else if err = response.EnsureValid(); err != nil {
		return fmt.Errorf("invalid terminate response received: %w", err)
	}
	return nil
}

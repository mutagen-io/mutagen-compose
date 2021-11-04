package mutagen

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mutagen-io/mutagen/pkg/grpcutil"
	"github.com/mutagen-io/mutagen/pkg/selection"
	synchronizationsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
	"github.com/mutagen-io/mutagen/pkg/synchronization"
	"github.com/mutagen-io/mutagen/pkg/url"
)

// volumeURLPrefix is the lowercase version of the volume URL prefix.
const volumeURLPrefix = "volume://"

// isVolumeURL checks if raw URL is a Docker Compose volume pseudo-URL.
func isVolumeURL(raw string) bool {
	return strings.HasPrefix(strings.ToLower(raw), volumeURLPrefix)
}

// mountPathForVolumeInMutagenContainer returns the mount path that will be used
// for a volume inside the Mutagen container. The path will be returned without
// a trailing slash. The volume must be non-empty or this function will panic.
// This function should only be called for supported Docker platforms.
func mountPathForVolumeInMutagenContainer(platform, volume string) string {
	// Verify that the volume is non-empty.
	if volume == "" {
		panic("empty volume name")
	}

	// Compute the path based on the daemon OS.
	switch platform {
	case "linux":
		return "/volumes/" + volume
	case "windows":
		return `c:\volumes\` + volume
	default:
		panic("unsupported Docker platform")
	}
}

// parseVolumeURL parses a Docker Compose volume pseudo-URL, converting it to a
// sidecar URL. This URL will only have kind, protocol, and path information
// set. The protocol will need to be changed to Docker and the container target
// and environment will need to be filled in once known. This function also
// returns the volume dependency for the URL. This function must only be called
// on URLs that have been classified as volume URLs by isVolumeURL, otherwise
// this function may panic.
func parseVolumeURL(raw, platform string) (*url.URL, string, error) {
	// Strip off the prefix
	raw = raw[len(volumeURLPrefix):]

	// Find the first slash, which will indicate the end of the volume name. If
	// no slash is found, then we assume that the volume itself is the target
	// synchronization root.
	var volume, path string
	if slashIndex := strings.IndexByte(raw, '/'); slashIndex < 0 {
		volume = raw
		path = mountPathForVolumeInMutagenContainer(platform, volume)
	} else if slashIndex == 0 {
		return nil, "", errors.New("empty volume name")
	} else {
		volume = raw[:slashIndex]
		path = mountPathForVolumeInMutagenContainer(platform, volume) + raw[slashIndex:]
	}

	// Create a Docker synchronization URL.
	return &url.URL{
		Kind:     url.Kind_Synchronization,
		Protocol: sidecarURLProtocol,
		Path:     path,
	}, volume, nil
}

// synchronizationSessionCurrent determines whether or not an existing
// synchronization session is equivalent to the specification for its creation.
func synchronizationSessionCurrent(
	session *synchronization.Session,
	specification *synchronizationsvc.CreationSpecification,
) bool {
	return session.Alpha.Equal(specification.Alpha) &&
		session.Beta.Equal(session.Beta) &&
		session.Configuration.Equal(specification.Configuration) &&
		session.ConfigurationAlpha.Equal(specification.ConfigurationAlpha) &&
		session.ConfigurationBeta.Equal(specification.ConfigurationBeta)
}

// synchronizationCreateWithSpecification creates a synchronization session
// using the provided synchronization service client, session specification, and
// prompter.
func synchronizationCreateWithSpecification(
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

// synchronizationFlushWithSelection flushes synchronization sessions using the
// provided synchronization service client, session selection, and prompter.
func synchronizationFlushWithSelection(
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

// synchronizationPauseWithSelection pauses synchronization sessions using the
// provided synchronization service client, session selection, and prompter.
func synchronizationPauseWithSelection(
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

// synchronizationResumeWithSelection resumes synchronization sessions using the
// provided synchronization service client, session selection, and prompter.
func synchronizationResumeWithSelection(
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

// synchronizationTerminateWithSelection terminates synchronization sessions
// using the provided synchronization service client, session selection, and
// prompter.
func synchronizationTerminateWithSelection(
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

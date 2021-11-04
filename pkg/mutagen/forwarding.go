package mutagen

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mutagen-io/mutagen/pkg/forwarding"
	"github.com/mutagen-io/mutagen/pkg/grpcutil"
	"github.com/mutagen-io/mutagen/pkg/selection"
	forwardingsvc "github.com/mutagen-io/mutagen/pkg/service/forwarding"
	"github.com/mutagen-io/mutagen/pkg/url"
	forwardingurl "github.com/mutagen-io/mutagen/pkg/url/forwarding"
)

// networkURLPrefix is the lowercase version of the network URL prefix.
const networkURLPrefix = "network://"

// isNetworkURL checks if raw URL is a Docker Compose network pseudo-URL.
func isNetworkURL(raw string) bool {
	return strings.HasPrefix(strings.ToLower(raw), networkURLPrefix)
}

// isTCPForwardingProtocol checks if a forwarding protocol is TCP-based.
func isTCPForwardingProtocol(protocol string) bool {
	switch protocol {
	case "tcp":
		return true
	case "tcp4":
		return true
	case "tcp6":
		return true
	default:
		return false
	}
}

// parseNetworkURL parses a Docker Compose network pseudo-URL, enforces that its
// forwarding endpoint protocol is TCP-based, and converts it to a sidecar
// forwarding URL. This URL will only have kind, protocol, and path information
// set. The protocol will need to be changed to Docker and the container target
// and environment will need to be filled in once known. This function also
// returns the network dependency for the URL. This function must only be called
// on URLs that have been classified as network URLs by isNetworkURL, otherwise
// it may panic.
func parseNetworkURL(raw string) (*url.URL, string, error) {
	// Strip off the prefix
	raw = raw[len(networkURLPrefix):]

	// Find the first colon, which will indicate the end of the network name.
	var network, endpoint string
	if colonIndex := strings.IndexByte(raw, ':'); colonIndex < 0 {
		return nil, "", errors.New("unable to find forwarding endpoint specification")
	} else if colonIndex == 0 {
		return nil, "", errors.New("empty network name")
	} else {
		network = raw[:colonIndex]
		endpoint = raw[colonIndex+1:]
	}

	// Parse the forwarding endpoint URL to ensure that it's valid and supported
	// for use with Docker Compose.
	if protocol, _, err := forwardingurl.Parse(endpoint); err != nil {
		return nil, "", fmt.Errorf("invalid forwarding endpoint URL: %w", err)
	} else if !isTCPForwardingProtocol(protocol) {
		return nil, "", fmt.Errorf("non-TCP-based forwarding endpoint (%s) unsupported", endpoint)
	}

	// Create a sidecar forwarding URL.
	return &url.URL{
		Kind:     url.Kind_Forwarding,
		Protocol: sidecarURLProtocol,
		Path:     endpoint,
	}, network, nil
}

// forwardingSessionCurrent determines whether or not an existing forwarding
// session is equivalent to the specification for its creation.
func forwardingSessionCurrent(
	session *forwarding.Session,
	specification *forwardingsvc.CreationSpecification,
) bool {
	return session.Source.Equal(specification.Source) &&
		session.Destination.Equal(specification.Destination) &&
		session.Configuration.Equal(specification.Configuration) &&
		session.ConfigurationSource.Equal(specification.ConfigurationSource) &&
		session.ConfigurationDestination.Equal(specification.ConfigurationDestination)
}

// forwardingCreateWithSpecification creates a forwarding session using the
// provided forwarding service client, session specification, and prompter.
func forwardingCreateWithSpecification(
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

// forwardingPauseWithSelection pauses forwarding sessions using the provided
// forwarding service client, session selection, and prompter.
func forwardingPauseWithSelection(
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

// forwardingResumeWithSelection resumes forwarding sessions using the provided
// forwarding service client, session selection, and prompter.
func forwardingResumeWithSelection(
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

// forwardingTerminateWithSelection terminates forwarding sessions using the
// provided forwarding service client, session selection, and prompter.
func forwardingTerminateWithSelection(
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

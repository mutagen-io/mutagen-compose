package mutagen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mutagen-io/mutagen/pkg/forwarding"
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

package mutagen

import (
	"errors"
	"strings"

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

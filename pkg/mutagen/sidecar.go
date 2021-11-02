package mutagen

import (
	"github.com/mutagen-io/mutagen/pkg/mutagen"
	"github.com/mutagen-io/mutagen/pkg/sidecar"
)

// sidecarImage is the full Mutagen sidecar image tag.
var sidecarImage string

func init() {
	// Compute the sidecar image tag.
	sidecarImage = sidecar.BaseTag + ":" + mutagen.Version
}

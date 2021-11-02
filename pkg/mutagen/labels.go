package mutagen

const (
	// sessionSidecarLabelKey is the name of the label applied to Mutagen
	// sessions to identify their associated Mutagen Compose sidecar container.
	sessionSidecarLabelKey = "io.mutagen.compose.sidecar"
)

// chopSidecarIdentifier chops off the 128-bit prefix of a 256-bit sidecar
// container identifier (encoded as a hex string) to make it fit into Mutagen
// session label values (which are limited to 63 characters). The first 128 bits
// of entropy should be more than sufficient to avoid collisions.
func chopSidecarIdentifier(sidecarID string) string {
	return sidecarID[:32]
}

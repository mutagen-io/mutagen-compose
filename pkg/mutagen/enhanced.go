package mutagen

const (
	// enhancedTagSuffix is the suffix to append to the sidecar image name to
	// identify the enhanced sidecar image.
	enhancedTagSuffix = "-enhanced"
)

// enhancedCapabilities are the capability specifications needed to enable
// enhanced sidecar features.
var enhancedCapabilities = []string{"SYS_ADMIN", "DAC_READ_SEARCH"}

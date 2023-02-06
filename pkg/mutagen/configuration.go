package mutagen

import (
	"github.com/mutagen-io/mutagen/pkg/api/models/forwarding"
	"github.com/mutagen-io/mutagen/pkg/api/models/synchronization"
)

// sidecarConfiguration encodes sidecar service configuration.
type sidecarConfiguration struct {
	// Features controls the sidecar feature set.
	//
	// Deprecated: This field is no longer used because all builds now target
	// the sidecar image with SSPL-licensed enhancements.
	Features string `mapstructure:"features"`
	// Restart is the restart policy for the sidecar container.
	Restart string `mapstructure:"restart"`
	// ContainerName is the name given to the sidecar container.
	ContainerName string `mapstructure:"container_name"`
}

// forwardingConfiguration encodes a forwarding session specification.
type forwardingConfiguration struct {
	// Source is the source URL for the session.
	Source string `mapstructure:"source"`
	// Destination is the destination URL for the session.
	Destination string `mapstructure:"destination"`
	// Configuration is the configuration for the session.
	Configuration forwarding.Configuration `mapstructure:",squash"`
	// ConfigurationSource is the source-specific configuration for the session.
	ConfigurationSource forwarding.Configuration `mapstructure:"configurationSource"`
	// ConfigurationDestination is the destination-specific configuration for
	// the session.
	ConfigurationDestination forwarding.Configuration `mapstructure:"configurationDestination"`
}

// synchronizationConfiguration encodes a synchronization session specification.
type synchronizationConfiguration struct {
	// Alpha is the alpha URL for the session.
	Alpha string `mapstructure:"alpha"`
	// Beta is the beta URL for the session.
	Beta string `mapstructure:"beta"`
	// Configuration is the configuration for the session.
	Configuration synchronization.Configuration `mapstructure:",squash"`
	// ConfigurationAlpha is the alpha-specific configuration for the session.
	ConfigurationAlpha synchronization.Configuration `mapstructure:"configurationAlpha"`
	// ConfigurationBeta is the beta-specific configuration for the session.
	ConfigurationBeta synchronization.Configuration `mapstructure:"configurationBeta"`
}

// configuration encodes collections of Mutagen forwarding and synchronization
// sessions found under an "x-mutagen" extension field.
type configuration struct {
	// Sidecar represents the sidecar service configuration.
	Sidecar sidecarConfiguration `mapstructure:"sidecar"`
	// Forwarding represents the forwarding sessions to be created. If a
	// "defaults" key is present, it is treated as a template upon which other
	// configurations are layered, thus keeping syntactic compatibility with the
	// global Mutagen configuration file.
	Forwarding map[string]forwardingConfiguration `mapstructure:"forward"`
	// Synchronization represents the forwarding sessions to be created. If a
	// "defaults" key is present, it is treated as a template upon which other
	// configurations are layered, thus keeping syntactic compatibility with the
	// global Mutagen configuration file.
	Synchronization map[string]synchronizationConfiguration `mapstructure:"sync"`
}

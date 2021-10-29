package mutagen

import (
	"github.com/mutagen-io/mutagen/pkg/configuration/forwarding"
	"github.com/mutagen-io/mutagen/pkg/configuration/synchronization"
)

// forwardingConfiguration encodes a forwarding session specification.
type forwardingConfiguration struct {
	// Source is the source URL for the session.
	Source string `yaml:"source"`
	// Destination is the destination URL for the session.
	Destination string `yaml:"destination"`
	// Configuration is the configuration for the session.
	Configuration forwarding.Configuration `mapstructure:",squash" yaml:",inline"`
	// ConfigurationSource is the source-specific configuration for the session.
	ConfigurationSource forwarding.Configuration `yaml:"configurationSource"`
	// ConfigurationDestination is the destination-specific configuration for
	// the session.
	ConfigurationDestination forwarding.Configuration `yaml:"configurationDestination"`
}

// synchronizationConfiguration encodes a synchronization session specification.
type synchronizationConfiguration struct {
	// Alpha is the alpha URL for the session.
	Alpha string `yaml:"alpha"`
	// Beta is the beta URL for the session.
	Beta string `yaml:"beta"`
	// Configuration is the configuration for the session.
	Configuration synchronization.Configuration `mapstructure:",squash" yaml:",inline"`
	// ConfigurationAlpha is the alpha-specific configuration for the session.
	ConfigurationAlpha synchronization.Configuration `yaml:"configurationAlpha"`
	// ConfigurationBeta is the beta-specific configuration for the session.
	ConfigurationBeta synchronization.Configuration `yaml:"configurationBeta"`
}

// configuration encodes collections of Mutagen forwarding and synchronization
// sessions found under an "x-mutagen" extension field.
type configuration struct {
	// Forwarding represents the forwarding sessions to be created. If a
	// "defaults" key is present, it is treated as a template upon which other
	// configurations are layered, thus keeping syntactic compatibility with the
	// global Mutagen configuration file.
	Forwarding map[string]forwardingConfiguration `mapstructure:"forward" yaml:"forward"`
	// Synchronization represents the forwarding sessions to be created. If a
	// "defaults" key is present, it is treated as a template upon which other
	// configurations are layered, thus keeping syntactic compatibility with the
	// global Mutagen configuration file.
	Synchronization map[string]synchronizationConfiguration `mapstructure:"sync" yaml:"sync"`
}

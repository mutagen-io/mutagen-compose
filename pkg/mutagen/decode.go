package mutagen

import (
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/mutagen-io/mutagen/pkg/synchronization/core"
)

// boolToIgnoreVCSModeHookFunc returns a mapstructure.DecodeHookFunc that will
// convert boolean types into an IgnoreVCSMode. This hook is necessary because
// the IgnoreVCSMode.UnmarshalText method won't be invoked by the YAML decoding
// performed by Compose (which doesn't know to dispatch text to custom
// unmarshalling functions), and thus the value we eventually receive will
// likely already be decoded into a boolean (though string representations of
// booleans, i.e. "true" or "false", will still be handled by the
// IgnoreVCSMode.UnmarshalText method).
func boolToIgnoreVCSModeHookFunc() mapstructure.DecodeHookFuncType {
	return func(valueType reflect.Type, storageType reflect.Type, data any) (any, error) {
		// If the incoming type isn't a boolean, then we're done.
		if valueType.Kind() != reflect.Bool {
			return data, nil
		}

		// If the storage isn't an IgnoreVCSMode, then we're done.
		if storageType != reflect.TypeOf(core.IgnoreVCSMode_IgnoreVCSModeDefault) {
			return data, nil
		}

		// Otherwise, perform conversion.
		if data.(bool) {
			return core.IgnoreVCSMode_IgnoreVCSModeIgnore, nil
		} else {
			return core.IgnoreVCSMode_IgnoreVCSModePropagate, nil
		}
	}
}

package docker

import (
	"fmt"

	"github.com/spf13/pflag"
)

// Flags stores top-level Docker CLI flags.
type Flags struct {
	// config stores the value of the --config flag.
	config string
	// context stores the value of the -c/--context flag.
	context string
	// debug indicates the presence of the -D/--debug flag.
	debug bool
	// host stores the value of the -H/--host flag.
	host string
	// logLevel stores the value of the -l/--log-level flag.
	logLevel string
	// tls indicates the presence of the --tls flag.
	tls bool
	// tlsCACert stores the value of the --tlscacert flag.
	tlsCACert string
	// tlsCert stores the value of the --tlscert flag.
	tlsCert string
	// tlsKey stores the value of the --tlskey flag.
	tlsKey string
	// tlsVerify indicates the presence of the --tlsverify flag.
	tlsVerify bool
}

// Register registers the flags into the specified flag set.
func (f *Flags) Register(flags *pflag.FlagSet) {
	flags.StringVar(&f.config, "config", "", "")
	flags.StringVarP(&f.context, "context", "c", "", "")
	flags.BoolVarP(&f.debug, "debug", "D", false, "")
	flags.StringVarP(&f.host, "host", "H", "", "")
	flags.StringVarP(&f.logLevel, "log-level", "l", "", "")
	flags.BoolVar(&f.tls, "tls", false, "")
	flags.StringVar(&f.tlsCACert, "tlscacert", "", "")
	flags.StringVar(&f.tlsCert, "tlscert", "", "")
	flags.StringVar(&f.tlsKey, "tlskey", "", "")
	flags.BoolVar(&f.tlsVerify, "tlsverify", false, "")
}

// Reconstituted constructs a representation of the flags suitable for parsing
// by their native command. It requires access to the flag set with which its
// flags were registered so that it can determine which were set. If any of the
// flags registered with the flag set have been removed, this method will panic.
func (f *Flags) Reconstituted(flags *pflag.FlagSet) (result []string) {
	if flags.Lookup("config").Changed {
		result = append(result, "--config", f.config)
	}
	if flags.Lookup("context").Changed {
		result = append(result, "--context", f.context)
	}
	if flags.Lookup("debug").Changed {
		result = append(result, fmt.Sprintf("--debug=%t", f.debug))
	}
	if flags.Lookup("host").Changed {
		result = append(result, "--host", f.host)
	}
	if flags.Lookup("log-level").Changed {
		result = append(result, "--log-level", f.logLevel)
	}
	if flags.Lookup("tls").Changed {
		result = append(result, fmt.Sprintf("--tls=%t", f.tls))
	}
	if flags.Lookup("tlscacert").Changed {
		result = append(result, "--tlscacert", f.tlsCACert)
	}
	if flags.Lookup("tlscert").Changed {
		result = append(result, "--tlscert", f.tlsCert)
	}
	if flags.Lookup("tlskey").Changed {
		result = append(result, "--tlskey", f.tlsKey)
	}
	if flags.Lookup("tlsverify").Changed {
		result = append(result, fmt.Sprintf("--tlsverify=%t", f.tlsVerify))
	}
	return
}

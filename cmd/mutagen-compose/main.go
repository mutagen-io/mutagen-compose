package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/mutagen-io/mutagen/cmd/external"

	"github.com/mutagen-io/mutagen-compose/pkg/compose"
	"github.com/mutagen-io/mutagen-compose/pkg/docker"
	"github.com/mutagen-io/mutagen-compose/pkg/mutagen"
)

func init() {
	// Set flags for invoking Mutagen cmd packages externally.
	external.UsePathBasedLookupForDaemonStart = true
}

func main() {
	// Create storage for top-level Docker and Compose flags.
	dockerFlags := &docker.Flags{}
	composeFlags := &compose.Flags{}

	// Create top-level flag set for parsing.
	var help bool
	flags := pflag.NewFlagSet("mutagen-compose", pflag.ContinueOnError)
	dockerFlags.Register(flags)
	composeFlags.Register(flags)
	flags.BoolVarP(&help, "help", "h", false, "")

	// Mark the shorthand help flag as deprecated to match the behavior of the
	// Docker CLI. We'll alias any help flags that we parse to their full --help
	// form, so this warning won't be duplicated by Docker. Any -h flags passed
	// after the Compose command name won't be handled by our shadow parsing,
	// but will have a corresponding warning printed by the parsing performed by
	// the plugin framework.
	flags.MarkShorthandDeprecated("help", "please use --help")

	// Disable interspersed arguments. This will cause parsing to terminate once
	// the first non-flag argument is found and all remaining unparsed arguments
	// to be stored. In this case, that will be everything from the Compose
	// command name onward.
	flags.SetInterspersed(false)

	// Perform parsing.
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Extract what remains of the command line arguments, which will be the
	// Compose subcommand and its flags/arguments, if any.
	commandAndArguments := flags.Args()

	// Handle emulation of Cobra's default top-level help command. We use our
	// faux top-level command (as opposed to doing a passthrough to plugin
	// invocation) because it will have the correct command name when invoking
	// things like "mutagen-compose help --help", which we can't otherwise
	// override for the real top-level help command.
	if len(commandAndArguments) > 0 && commandAndArguments[0] == "help" {
		emulatedArgs := []string{"mutagen-compose"}
		if help {
			emulatedArgs = append(emulatedArgs, "--help")
		}
		emulatedArgs = append(emulatedArgs, "help")
		emulatedArgs = append(emulatedArgs, commandAndArguments[1:]...)
		os.Args = emulatedArgs
		if err := fauxTopLevelCommandForHelpAndUsage().Execute(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Compute the emulated arguments that we'll use for the plugin-based
	// invocation of Compose.
	emulatedArgs := []string{"docker"}
	if help {
		emulatedArgs = append(emulatedArgs, "--help")
	}
	emulatedArgs = append(emulatedArgs, dockerFlags.Reconstituted(flags)...)
	emulatedArgs = append(emulatedArgs, "compose")
	emulatedArgs = append(emulatedArgs, composeFlags.Reconstituted(flags)...)
	emulatedArgs = append(emulatedArgs, commandAndArguments...)
	os.Args = emulatedArgs

	// Create the Mutagen liaison.
	liaison := &mutagen.Liaison{}

	// Invoke Compose.
	invokeCompose(liaison)
}

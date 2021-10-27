package main

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/mutagen-io/mutagen/cmd"

	"github.com/docker/cli/cli"

	commands "github.com/docker/compose/v2/cmd/compose"

	"github.com/docker/compose/v2/pkg/api"

	"github.com/mutagen-io/mutagen-compose/pkg/compose"
	"github.com/mutagen-io/mutagen-compose/pkg/docker"
)

const (
	// commandName is the command name for Mutagen Compose.
	commandName = "mutagen-compose"
	// commandDescription is the description for Mutagen Compose.
	commandDescription = "Mutagen Compose"
)

// showTopLevelUsage shows the top-level usage message. It does this by merging
// the top-level Docker and top-level Compose flags into a single Cobra command.
func showTopLevelUsage() {
	// Create a top-level Compose command and replace its command name.
	root := commands.RootCommand(api.NewServiceProxy())
	root.Use = commandName
	root.Short = commandDescription

	// HACK: Set this command up as a Docker root command in order to add the
	// top-level Docker CLI flags and to set usage formatting.
	cli.SetupRootCommand(root)

	// HACK: Our -H/--host flag only supports a single value, but the Docker CLI
	// -H/--host flag supports multiple specifications. To correct this in help
	// output, override the usage message and replace the value storage with one
	// that will have the correct type.
	hostFlag := root.Flags().Lookup("host")
	hostFlag.Usage = "Docker daemon host specification"
	hostFlag.Value = root.Flags().Lookup("context").Value

	// HACK: Disable help annotations.
	root.Annotations = nil

	// Display usage information.
	root.Usage()
}

func main() {
	// Create storage for top-level configuration.
	dockerFlags := &docker.Flags{}
	composeFlags := &compose.Flags{}

	// Create top-level flag set for parsing.
	var help, version bool
	flags := pflag.NewFlagSet("mutagen-compose", pflag.ContinueOnError)
	dockerFlags.Register(flags)
	composeFlags.Register(flags)
	flags.BoolVarP(&help, "help", "h", false, "")
	flags.BoolVarP(&version, "version", "v", false, "")

	// Disable interspersed arguments. This will cause parsing to terminate once
	// the first non-flag argument is found and all remaining unparsed arguments
	// to be stored. In this case, that will be everything from the Compose
	// command name onward.
	flags.SetInterspersed(false)

	// Perform parsing.
	if err := flags.Parse(os.Args[1:]); err != nil {
		cmd.Fatal(err)
	}

	// Extract the command name and arguments.
	commandAndArguments := flags.Args()

	// TODO: Perform some sort of validation that the command is valid. I would
	// do this by checking which commands are registered at the root of a
	// Compose command. We also don't want this to accidentally run as a plugin.
	// If the command is invalid, we'll want to show our own top-level usage
	// information and not delegate that responsibility to Compose.

	// TODO: Figure out how to override the usage line of commands in help. At
	// the moment, they still yield "docker compose <name>", but we want them to
	// read "mutagen-compose <name>". The best idea I have at the moment is to
	// override the usage and help functions on all Compose commands to a custom
	// function that will retrieve their usage template, replace the
	// {{.UseLine}} segment with mutagen-compose {{.Use}} and then use their
	// parent command (the top-level Compose command) to retrieve the default
	// templating behavior (that we'd no longer be able to access due to our
	// help/usage function overrides) and execute that with the new template.

	// Handle help requests (including implicit help requests due to a lack of
	// flags or arguments). To mirror the behavior of Cobra's help dispatching,
	// we'll forward a top-level help request to a command if one is specified,
	// otherwise we'll display a custom top-level usage message. If a help
	// request is sent to the command directly, then we won't see it, but it
	// will still be handled by the dispatching below.
	//
	// TODO: In that case, what should we do if version is also specified? The
	// Docker CLI doesn't handle this well, e.g. docker --version compose up -h
	if help || len(os.Args) == 1 {
		if len(commandAndArguments) > 0 {
			os.Args = []string{"docker", "compose", commandAndArguments[0], "--help"}
			invokeCompose()
		} else {
			showTopLevelUsage()
		}
		os.Exit(0)
	}

	// Handle version requests.
	if version {
		// TODO: Implement. We'll want to display Mutagen, Docker, and Compose
		// version information. We'll likely have to pull in Docker and Compose
		// version information via runtime/debug.ReadBuildInfo (possibly with a
		// sync.Once-based wrapper to amortize the lookup cost).
	}

	// TODO: Override the RunE entrypoint of the Compose version command as well
	// to provide our own version information.

	// Compute the emulated arguments that we'll use for the plugin-based
	// invocation of Compose.
	emulatedArgs := []string{"docker"}
	emulatedArgs = append(emulatedArgs, dockerFlags.Reconstituted(flags)...)
	emulatedArgs = append(emulatedArgs, "compose")
	emulatedArgs = append(emulatedArgs, composeFlags.Reconstituted(flags)...)
	emulatedArgs = append(emulatedArgs, commandAndArguments...)
	os.Args = emulatedArgs

	// Invoke Compose.
	invokeCompose()
}

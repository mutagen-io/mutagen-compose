package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli"

	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/cmd/formatter"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"

	versionpkg "github.com/mutagen-io/mutagen-compose/pkg/version"
)

const (
	// commandName is the command name for Mutagen Compose.
	commandName = "mutagen-compose"
	// commandDescription is the description for Mutagen Compose.
	commandDescription = "Mutagen Compose"
)

// fauxTopLevelCommandForHelpAndUsage returns a faux top-level Compose command
// whose help and usage information will include merged-in top-level Docker CLI
// flags supported by Mutagen Compose.
func fauxTopLevelCommandForHelpAndUsage() *cobra.Command {
	// Create a top-level Compose command and replace its command name.
	root := commands.RootCommand(api.NewServiceProxy())
	root.Use = commandName
	root.Short = commandDescription

	// Adjust the version command like we do for the real command hierarchy.
	adjustVersionCommand(root)

	// Add the legal command like we do for the real command hierarchy.
	root.AddCommand(legalCommand)

	// HACK: Set this command up as a Docker plugin root command in order to add
	// the top-level Docker CLI flags and to set usage formatting. Normally
	// there would be an intermediate command above this that would be the
	// actual root.
	cli.SetupPluginRootCommand(root)

	// Extract the unified flag set.
	flags := root.Flags()

	// HACK: Our -H/--host flag only supports a single value, but the Docker CLI
	// -H/--host flag supports multiple specifications. To correct this in help
	// output, override the usage message and replace the value storage with one
	// that will have the correct type.
	hostFlag := flags.Lookup("host")
	hostFlag.Usage = "Docker daemon host specification"
	hostFlag.Value = root.Flags().Lookup("context").Value

	// HACK: Remove mention of the -v/--version flag (brought in by Docker)
	// since we don't support it. A -v/--version flag was also added by Compose
	// in v2.0.2, but we don't support that either since it's a hidden flag that
	// is only provided for backward compatibility with Compose V1.
	flags.MarkHidden("version")

	// HACK: Disable help annotations.
	root.Annotations = nil

	// Done.
	return root
}

// adjustUsageInformation adjusts the Compose root command (and its subcommands)
// to display usage information that corresponds to Mutagen Compose.
func adjustUsageInformation(cmd *cobra.Command) {
	cmd.SetUsageFunc(func(c *cobra.Command) error {
		// Create a faux top-level command with proper usage information,
		// including merged-in top-level Docker CLI flags that we support.
		faux := fauxTopLevelCommandForHelpAndUsage()

		// If usage information has been requested for the Compose root
		// command, then use the faux command display usage information.
		if c == cmd {
			return faux.Usage()
		}

		// Otherwise, this is a help request for a Compose subcommand, so
		// reparent the subcommand onto the faux top-level command to get a
		// proper command name and then display its usage.
		faux.AddCommand(c)
		return c.Usage()
	})
}

const (
	// unknownCommandErrorPrefix is the error prefix used by unknown command
	// errors in Compose.
	unknownCommandErrorPrefix = `unknown docker command: "compose`
	// replacementUnknownCommandErrorPrefix is the altered error prefix used by
	// unknown command errors in Compose.
	replacementUnknownCommandErrorPrefix = `unknown command: "` + commandName
)

// adjustUnknownCommandErrors adjusts the Compose root command to return unknown
// command errors that correspond to Mutagen Compose.
func adjustUnknownCommandErrors(cmd *cobra.Command) {
	// Extract the original entrypoint.
	originalRunE := cmd.RunE

	// Override the entrypoint with one that changes the error message for
	// unknown command errors.
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := originalRunE(cmd, args)
		if err != nil {
			if statusErr, ok := err.(cli.StatusError); ok {
				if strings.HasPrefix(statusErr.Status, unknownCommandErrorPrefix) {
					err = cli.StatusError{
						StatusCode: compose.CommandSyntaxFailure.ExitCode,
						Status:     replacementUnknownCommandErrorPrefix + statusErr.Status[len(unknownCommandErrorPrefix):],
					}
				}
			}
		}
		return err
	}
}

// adjustVersionCommand adjust the behavior of the version command to correspond
// to Mutagen Compose.
func adjustVersionCommand(cmd *cobra.Command) {
	// Look up the version command.
	version, _, _ := cmd.Find([]string{"version"})

	// Replace its description.
	version.Short = fmt.Sprintf("Show the %s version information", commandDescription)

	// Extract its flags.
	flags := version.Flags()

	// Look up the short flag and replace its description.
	shortFlag := flags.Lookup("short")
	shortFlag.Usage = "Show only the version numbers."

	// Override the command entrypoint.
	version.RunE = func(cmd *cobra.Command, args []string) error {
		// Look up the format flag.
		formatFlag := flags.Lookup("format")

		// Extract flag values.
		format := formatFlag.Value.String()
		short := shortFlag.Value.String() == "true"

		// Load version information.
		versions, err := versionpkg.LoadVersions()
		if err != nil {
			return fmt.Errorf("unable to load version information: %w", err)
		}

		// Print accordingly. We don't perform any validation on format because
		// the built-in version command doesn't either.
		if short {
			fmt.Printf("%s/%s/%s\n", versions.Mutagen, versions.Compose, versions.Docker)
			return nil
		}
		if format == formatter.JSON {
			return json.NewEncoder(os.Stdout).Encode(versions)
		}
		fmt.Println("Mutagen version", versions.Mutagen)
		fmt.Println("Compose version", versions.Compose)
		fmt.Println("Docker version", versions.Docker)
		return nil
	}
}

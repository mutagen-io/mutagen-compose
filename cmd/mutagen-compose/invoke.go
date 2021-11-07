// Docker Compose plugin-based invocation. Based on (but modified from):
// https://github.com/docker/compose/blob/6476e10b9337ba58eabb817570ddd1fe66c48583/cmd/main.go
//
// The original code license:
//
//   Copyright 2020 Docker Compose CLI authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package main

import (
	"github.com/spf13/cobra"

	dockercli "github.com/docker/cli/cli"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"

	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"

	mutageninfo "github.com/mutagen-io/mutagen/pkg/mutagen"

	"github.com/mutagen-io/mutagen-compose/pkg/mutagen"
)

func init() {
	// Set the warning that's presented in the event of failure.
	commands.Warning = "Mutagen Compose is currently experimental. " +
		"Feedback can be provided at https://github.com/mutagen-io/mutagen-compose"
}

// invokeCompose invokes Compose via the plugin infrastructure. It requires that
// os.Args be set in a manner that emulates execution as a plugin.
func invokeCompose(liaison *mutagen.Liaison) {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		lazyInit := api.NewServiceProxy()
		cmd := commands.RootCommand(lazyInit)
		originalPreRun := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			liaison.RegisterDockerFlags(cmd.Root().Flags())
			liaison.RegisterDockerCLI(dockerCli)
			liaison.RegisterComposeService(compose.NewComposeService(liaison.DockerClient(), dockerCli.ConfigFile()))
			lazyInit.WithService(liaison.ComposeService())
			if originalPreRun != nil {
				return originalPreRun(cmd, args)
			}
			return nil
		}
		cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
			return dockercli.StatusError{
				StatusCode: compose.CommandSyntaxFailure.ExitCode,
				Status:     err.Error(),
			}
		})
		adjustUsageInformation(cmd)
		adjustUnknownCommandErrors(cmd)
		adjustVersionCommand(cmd)
		cmd.AddCommand(legalCommand)
		return cmd
	},
		manager.Metadata{
			SchemaVersion: "0.1.0",
			Vendor:        "Mutagen IO, Inc.",
			Version:       mutageninfo.Version,
		})
}

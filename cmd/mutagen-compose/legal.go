package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mutagen-io/mutagen/cmd"

	"github.com/mutagen-io/mutagen-compose/pkg/legal"
)

// legalMain is the entry point for the legal command.
func legalMain(_ *cobra.Command, _ []string) error {
	// Print legal information.
	fmt.Println(legal.Licenses)

	// Success.
	return nil
}

// legalCommand is the legal command.
var legalCommand = &cobra.Command{
	Use:          "legal",
	Short:        "Show legal information",
	Args:         cmd.DisallowArguments,
	RunE:         legalMain,
	SilenceUsage: true,
}

package compose

import (
	"fmt"

	"github.com/spf13/pflag"
)

// Flags stores top-level Compose flags.
type Flags struct {
	// ansi stores the value of the --ansi flag.
	ansi string
	// compatibility indicates the presence of the --compatibility flag.
	compatibility bool
	// dryRun indicates the presence of the --dry-run flag.
	dryRun bool
	// envFile stores the value of the --env-file flag.
	envFile string
	// files stores the value(s) of the -f/--file flag(s).
	files []string
	// profiles stores the value(s) of the --profile flag(s).
	profiles []string
	// projectDirectory stores the value of the --project-directory flag.
	projectDirectory string
	// projectName stores the value of the -p/--project-name flag.
	projectName string
}

// Register registers the flags into the specified flag set.
func (f *Flags) Register(flags *pflag.FlagSet) {
	flags.StringVar(&f.ansi, "ansi", "", "")
	flags.BoolVar(&f.compatibility, "compatibility", false, "")
	flags.BoolVar(&f.dryRun, "dry-run", false, "")
	flags.StringVar(&f.envFile, "env-file", "", "")
	flags.StringSliceVarP(&f.files, "file", "f", nil, "")
	flags.StringSliceVar(&f.profiles, "profile", nil, "")
	flags.StringVar(&f.projectDirectory, "project-directory", "", "")
	flags.StringVarP(&f.projectName, "project-name", "p", "", "")
}

// Reconstituted constructs a representation of the flags suitable for parsing
// by their native command. It requires access to the flag set with which its
// flags were registered so that it can determine which were set. If any of the
// flags registered with the flag set have been removed, this method will panic.
func (f *Flags) Reconstituted(flags *pflag.FlagSet) (result []string) {
	if flags.Lookup("ansi").Changed {
		result = append(result, "--ansi", f.ansi)
	}
	if flags.Lookup("compatibility").Changed {
		result = append(result, fmt.Sprintf("--compatibility=%t", f.compatibility))
	}
	if flags.Lookup("dry-run").Changed {
		result = append(result, fmt.Sprintf("--dry-run=%t", f.dryRun))
	}
	if flags.Lookup("env-file").Changed {
		result = append(result, "--env-file", f.envFile)
	}
	if flags.Lookup("file").Changed {
		for _, file := range f.files {
			result = append(result, "--file", file)
		}
	}
	if flags.Lookup("profile").Changed {
		for _, profile := range f.profiles {
			result = append(result, "--profile", profile)
		}
	}
	if flags.Lookup("project-directory").Changed {
		result = append(result, "--project-directory", f.projectDirectory)
	}
	if flags.Lookup("project-name").Changed {
		result = append(result, "--project-name", f.projectName)
	}
	return
}

package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/mutagen-io/mutagen-compose/pkg/version"
)

const (
	// composePackage is the Go package URL to use for building Mutagen Compose.
	composePackage = "github.com/mutagen-io/mutagen-compose/cmd/mutagen-compose"

	// buildDirectoryName is the name of the build directory to create inside
	// the root of the Mutagen Compose source tree.
	buildDirectoryName = "build"

	// composeBuildSubdirectoryName is the name of the build subdirectory where
	// Mutagen Compose binaries are built.
	composeBuildSubdirectoryName = "compose"
	// releaseBuildSubdirectoryName is the name of the build subdirectory where
	// release bundles are built.
	releaseBuildSubdirectoryName = "release"

	// composeBaseName is the name of the Mutagen Compose binary without any
	// path or extension.
	composeBaseName = "mutagen-compose"

	// minimumARMSupport is the value to pass to the GOARM environment variable
	// when building binaries. We currently specify support for ARMv5. This will
	// enable software-based floating point. For our use case, this is totally
	// fine, because we don't have any floating-point-heavy code, and the
	// resulting binary bloat is very minimal. This won't apply for arm64, which
	// always has hardware-based floating point support. For more information,
	// see: https://github.com/golang/go/wiki/GoArm
	minimumARMSupport = "5"
)

// sourceTreePath computes the path to the source directory.
func sourceTreePath() (string, error) {
	// Compute the path to this script.
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to compute script path")
	}

	// Compute the path to the source directory.
	return filepath.Dir(filepath.Dir(filePath)), nil
}

// Target specifies a GOOS/GOARCH combination.
type Target struct {
	// GOOS is the GOOS environment variable specification for the target.
	GOOS string
	// GOARCH is the GOARCH environment variable specification for the target.
	GOARCH string
}

// String generates a human-readable representation of the target.
func (t Target) String() string {
	return fmt.Sprintf("%s/%s", t.GOOS, t.GOARCH)
}

// Name generates a representation of the target that is suitable for paths and
// file names.
func (t Target) Name() string {
	return fmt.Sprintf("%s_%s", t.GOOS, t.GOARCH)
}

// ExecutableName formats executable names for the target.
func (t Target) ExecutableName(base string) string {
	// If we're on Windows, append a ".exe" extension.
	if t.GOOS == "windows" {
		return fmt.Sprintf("%s.exe", base)
	}

	// Otherwise return the base name unmodified.
	return base
}

// appendGoEnv modifies an environment specification to make the Go toolchain
// generate output for the target. It assumes that the resulting environment
// will be used with os/exec.Cmd and thus doesn't avoid duplicate variables.
func (t Target) appendGoEnv(environment []string) []string {
	// Override GOOS/GOARCH.
	environment = append(environment, fmt.Sprintf("GOOS=%s", t.GOOS))
	environment = append(environment, fmt.Sprintf("GOARCH=%s", t.GOARCH))

	// Disable cgo.
	environment = append(environment, "CGO_ENABLED=0")

	// Set up ARM target support. See notes for definition of minimumARMSupport.
	// We don't need to unset any existing GOARM variables since they simply
	// won't be used if we're not targeting (non-64-bit) ARM systems.
	if t.GOARCH == "arm" {
		environment = append(environment, fmt.Sprintf("GOARM=%s", minimumARMSupport))
	}

	// Done.
	return environment
}

// IsCrossTarget determines whether or not the target represents a
// cross-compilation target (i.e. not the native target for the current Go
// toolchain).
func (t Target) IsCrossTarget() bool {
	return t.GOOS != runtime.GOOS || t.GOARCH != runtime.GOARCH
}

// IncludeInSlimBuildModes indicates whether or not the target should be
// included in slim builds.
func (t Target) IncludeInSlimBuildModes() bool {
	return !t.IsCrossTarget() ||
		(t.GOOS == "darwin") ||
		(t.GOOS == "windows" && t.GOARCH == "amd64") ||
		(t.GOOS == "linux" && (t.GOARCH == "amd64" || t.GOARCH == "arm"))
}

// Build executes a module-aware build of the specified package URL, storing the
// output of the build at the specified path.
func (t Target) Build(url, output string, ldflags string) error {
	// Compute the build command.
	arguments := []string{"build", "-o", output, "-trimpath"}
	if ldflags != "" {
		arguments = append(arguments, "-ldflags="+ldflags)
	}
	arguments = append(arguments, url)

	// Create the build command.
	builder := exec.Command("go", arguments...)

	// Set the environment.
	builder.Env = t.appendGoEnv(builder.Environ())

	// Forward input, output, and error streams.
	builder.Stdin = os.Stdin
	builder.Stdout = os.Stdout
	builder.Stderr = os.Stderr

	// Run the build.
	return builder.Run()
}

// targets encodes which combinations of GOOS and GOARCH we want to use for
// building Mutagen Compose binaries. This list is a subset of the list in the
// Mutagen build script. Unfortunately we can't support every platform supported
// by Mutagen (at least not at the moment) due to patches required by Docker
// dependencies on some platforms. Instead, we target the big three: macOS,
// Linux, and Windows (with the same architecture support as Mutagen).
var targets = []Target{
	// Define macOS targets.
	{"darwin", "amd64"},
	{"darwin", "arm64"},

	// Define Linux targets.
	{"linux", "386"},
	{"linux", "amd64"},
	{"linux", "arm"},
	{"linux", "arm64"},
	{"linux", "ppc64"},
	{"linux", "ppc64le"},
	{"linux", "mips"},
	{"linux", "mipsle"},
	{"linux", "mips64"},
	{"linux", "mips64le"},
	{"linux", "riscv64"},
	{"linux", "s390x"},

	// Define Windows targets.
	{"windows", "386"},
	{"windows", "amd64"},
	{"windows", "arm"},
	{"windows", "arm64"},
}

// macOSCodeSign performs macOS code signing on the specified path using the
// specified signing identity. It performs code signing in a manner suitable for
// later submission to Apple for notarization.
func macOSCodeSign(path, identity string) error {
	// Create the code signing command.
	//
	// We include the --force flag because the Go toolchain won't touch binaries
	// if they don't need to be rebuilt and thus we might have a signature from
	// a previous build. In that case, the code signing operation will fail
	// without --force. When --force is specified, any existing signature will
	// be overwritten, unless it's using the same code signing identity, in
	// which case it will simply be left in place (which is actually optimal for
	// for repeated local usage). Note that the --force flag is not required to
	// override ad hoc signatures (which the Go toolchain will add by default
	// darwin/arm64 binaries).
	//
	// The --options runtime and --timestamp flags are required to enable the
	// hardened runtime (which doesn't affect Mutagen Compose binaries) and to
	// use a secure signing timestamp, both of which are required for
	// notarization.
	codesign := exec.Command("codesign",
		"--sign", identity,
		"--force",
		"--options", "runtime",
		"--timestamp",
		"--verbose",
		path,
	)

	// Forward input, output, and error streams.
	codesign.Stdin = os.Stdin
	codesign.Stdout = os.Stdout
	codesign.Stderr = os.Stderr

	// Run code signing.
	return codesign.Run()
}

// archiveBuilderCopyBufferSize determines the size of the copy buffer used when
// generating archive files.
// TODO: Figure out if we should set this on a per-machine basis. This value is
// taken from Go's io.Copy method, which defaults to allocating a 32k buffer if
// none is provided.
const archiveBuilderCopyBufferSize = 32 * 1024

type ArchiveBuilder struct {
	file       *os.File
	compressor *gzip.Writer
	archiver   *tar.Writer
	copyBuffer []byte
}

func NewArchiveBuilder(bundlePath string) (*ArchiveBuilder, error) {
	// Open the underlying file.
	file, err := os.Create(bundlePath)
	if err != nil {
		return nil, fmt.Errorf("unable to create target file: %w", err)
	}

	// Create the compressor.
	compressor, err := gzip.NewWriterLevel(file, gzip.BestCompression)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("unable to create compressor: %w", err)
	}

	// Success.
	return &ArchiveBuilder{
		file:       file,
		compressor: compressor,
		archiver:   tar.NewWriter(compressor),
		copyBuffer: make([]byte, archiveBuilderCopyBufferSize),
	}, nil
}

func (b *ArchiveBuilder) Close() error {
	// Close in the necessary order to trigger flushes.
	if err := b.archiver.Close(); err != nil {
		b.compressor.Close()
		b.file.Close()
		return fmt.Errorf("unable to close archiver: %w", err)
	} else if err := b.compressor.Close(); err != nil {
		b.file.Close()
		return fmt.Errorf("unable to close compressor: %w", err)
	} else if err := b.file.Close(); err != nil {
		return fmt.Errorf("unable to close file: %w", err)
	}

	// Success.
	return nil
}

func (b *ArchiveBuilder) Add(name, path string, mode int64) error {
	// If the name is empty, use the base name.
	if name == "" {
		name = filepath.Base(path)
	}

	// Open the file and ensure its cleanup.
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	// Compute its size.
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to determine file size: %w", err)
	}
	size := stat.Size()

	// Write the header for the entry.
	header := &tar.Header{
		Name:    name,
		Mode:    mode,
		Size:    size,
		ModTime: time.Now(),
	}
	if err := b.archiver.WriteHeader(header); err != nil {
		return fmt.Errorf("unable to write archive header: %w", err)
	}

	// Copy the file contents.
	if _, err := io.CopyBuffer(b.archiver, file, b.copyBuffer); err != nil {
		return fmt.Errorf("unable to write archive entry: %w", err)
	}

	// Success.
	return nil
}

// copyFile copies the contents at sourcePath to a newly created file at
// destinationPath that inherits the permissions of sourcePath.
func copyFile(sourcePath, destinationPath string) error {
	// Open the source file and defer its closure.
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("unable to open source file: %w", err)
	}
	defer source.Close()

	// Grab source file metadata.
	metadata, err := source.Stat()
	if err != nil {
		return fmt.Errorf("unable to query source file metadata: %w", err)
	}

	// Remove the destination.
	os.Remove(destinationPath)

	// Create the destination file and defer its closure. We open with exclusive
	// creation flags to ensure that we're the ones creating the file so that
	// its permissions are set correctly.
	destination, err := os.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, metadata.Mode()&os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create destination file: %w", err)
	}
	defer destination.Close()

	// Copy contents.
	if count, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("unable to copy data: %w", err)
	} else if count != metadata.Size() {
		return errors.New("copied size does not match expected")
	}

	// Success.
	return nil
}

var usage = `usage: build [-h|--help] [-m|--mode=<mode>]
       [--macos-codesign-identity=<identity>]

The mode flag accepts three values: 'local', 'slim', and 'release'. 'local' will
build Mutagen Compose only for the current platform. 'slim' will build Mutagen
Compose for a common subset of platforms. 'release' will build Mutagen Compose
for all platforms and package for release. The default mode is 'local'.

If --macos-codesign-identity specifies a non-empty value, then it will be used
to perform code signing on all macOS binaries in a fashion suitable for
notarization by Apple. The codesign utility must be able to access the
associated certificate and private keys in Keychain Access without a password if
this script is operated in a non-interactive mode.
`

// build is the primary entry point.
func build() error {
	// Parse command line arguments.
	flagSet := pflag.NewFlagSet("build", pflag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	var mode, macosCodesignIdentity string
	flagSet.StringVarP(&mode, "mode", "m", "local", "specify the build mode")
	flagSet.StringVar(&macosCodesignIdentity, "macos-codesign-identity", "", "specify the macOS code signing identity")
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		if err == pflag.ErrHelp {
			fmt.Fprint(os.Stdout, usage)
			return nil
		} else {
			return fmt.Errorf("unable to parse command line: %w", err)
		}
	}
	if !(mode == "local" || mode == "slim" || mode == "release") {
		return fmt.Errorf("invalid build mode: %s", mode)
	}

	// If a macOS code signing identity has been specified, then make sure we're
	// in a mode where that makes sense.
	if macosCodesignIdentity != "" && runtime.GOOS != "darwin" {
		return errors.New("macOS is required for macOS code signing")
	}

	// Compute the path to the Mutagen Compose source directory.
	mutagenComposeSourcePath, err := sourceTreePath()
	if err != nil {
		return fmt.Errorf("unable to compute source tree path: %w", err)
	}

	// Verify that we're running inside the Mutagen Compose source directory,
	// otherwise we can't rely on Go modules working.
	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to compute working directory: %w", err)
	}
	workingDirectoryRelativePath, err := filepath.Rel(mutagenComposeSourcePath, workingDirectory)
	if err != nil {
		return fmt.Errorf("unable to determine working directory relative path: %w", err)
	}
	if strings.Contains(workingDirectoryRelativePath, "..") {
		return errors.New("build script run outside Mutagen Compose source tree")
	}

	// Compute the path to the build directory and ensure that it exists.
	buildPath := filepath.Join(mutagenComposeSourcePath, buildDirectoryName)
	if err := os.MkdirAll(buildPath, 0700); err != nil {
		return fmt.Errorf("unable to create build directory: %w", err)
	}

	// Create the necessary build directory hierarchy.
	composeBuildSubdirectoryPath := filepath.Join(buildPath, composeBuildSubdirectoryName)
	releaseBuildSubdirectoryPath := filepath.Join(buildPath, releaseBuildSubdirectoryName)
	if err := os.MkdirAll(composeBuildSubdirectoryPath, 0700); err != nil {
		return fmt.Errorf("unable to create Compose build subdirectory: %w", err)
	}
	if mode == "release" {
		if err := os.MkdirAll(releaseBuildSubdirectoryPath, 0700); err != nil {
			return fmt.Errorf("unable to create release build subdirectory: %w", err)
		}
	}

	// Compute the local target.
	localTarget := Target{runtime.GOOS, runtime.GOARCH}

	// Compute active targets.
	var activeTargets []Target
	for _, target := range targets {
		if mode == "local" && target.IsCrossTarget() {
			continue
		} else if mode == "slim" && !target.IncludeInSlimBuildModes() {
			continue
		}
		activeTargets = append(activeTargets, target)
	}

	// Load version information if necessary.
	var versions *version.Versions
	if mode == "release" {
		versions, err = version.LoadVersions()
		if err != nil {
			return fmt.Errorf("unable to load version information: %w", err)
		}
	}

	// Compute ldflags.
	var ldflags string
	if mode == "release" {
		ldflags = "-X github.com/docker/compose/v2/internal.Version=" + versions.Compose
	}

	// Build binaries.
	log.Println("Building binaries...")
	for _, target := range activeTargets {
		log.Println("Build for", target)
		executableBuildPath := filepath.Join(composeBuildSubdirectoryPath, target.Name())
		if err := target.Build(composePackage, executableBuildPath, ldflags); err != nil {
			return fmt.Errorf("unable to build Mutagen Compose: %w", err)
		}
		if macosCodesignIdentity != "" && target.GOOS == "darwin" {
			if err := macOSCodeSign(executableBuildPath, macosCodesignIdentity); err != nil {
				return fmt.Errorf("unable to code sign Mutagen Compose for macOS: %w", err)
			}
		}
	}

	// Build release bundles if necessary.
	if mode == "release" {
		log.Println("Building release bundles...")
		for _, target := range activeTargets {
			// Update status.
			log.Println("Building release bundle for", target)

			// Compute paths.
			executableBuildPath := filepath.Join(composeBuildSubdirectoryPath, target.Name())
			releaseBundlePath := filepath.Join(
				releaseBuildSubdirectoryPath,
				fmt.Sprintf("mutagen-compose_%s_v%s.tar.gz", target.Name(), versions.Mutagen),
			)

			// Build the release bundle.
			if releaseBundle, err := NewArchiveBuilder(releaseBundlePath); err != nil {
				return fmt.Errorf("unable to create release bundle: %w", err)
			} else if err = releaseBundle.Add(target.ExecutableName(composeBaseName), executableBuildPath, 0755); err != nil {
				releaseBundle.Close()
				return fmt.Errorf("unable to add executable to release bundle: %w", err)
			} else if err = releaseBundle.Close(); err != nil {
				return fmt.Errorf("unable to finalize release bundle: %w", err)
			}
		}
	}

	// Relocate the Mutagen Compose executable for the current platform.
	log.Println("Copying binary for testing")
	localExecutableBuildPath := filepath.Join(composeBuildSubdirectoryPath, localTarget.Name())
	localExecutableRelocationPath := filepath.Join(buildPath, localTarget.ExecutableName(composeBaseName))
	if err := copyFile(localExecutableBuildPath, localExecutableRelocationPath); err != nil {
		return fmt.Errorf("unable to copy current platform executable: %w", err)
	}

	// Success.
	return nil
}

func main() {
	if err := build(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"

	"github.com/mutagen-io/mutagen/pkg/mutagen"
)

func main() {
	if mutagen.VersionTag != "" {
		fmt.Printf("%d.%d.%d-%s\n", mutagen.VersionMajor, mutagen.VersionMinor, mutagen.VersionPatch, mutagen.VersionTag)
	} else {
		fmt.Printf("%d.%d.%d\n", mutagen.VersionMajor, mutagen.VersionMinor, mutagen.VersionPatch)
	}
}

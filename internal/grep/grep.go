package grep

import (
	"fmt"
	"io"

	"github.com/spf13/pflag"
)

// Result represents the result of a grep operation
type Result struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}

// ColorMode represents the color output mode
type ColorMode string

const (
	ColorAlways ColorMode = "always"
	ColorNever  ColorMode = "never"
	ColorAuto   ColorMode = "auto"
)

type EmulatedGrepLaunchOptions struct {
	Stdin        []byte
	EmulateInTTY bool
}

// EmulateGrep provides a simplified interface that mimics grep command behavior
func EmulateGrep(args []string, launchOptions EmulatedGrepLaunchOptions) Result {
	flagset := pflag.NewFlagSet("grep", pflag.ContinueOnError)

	// We aren't using this in a command line, so disable usage and error messages
	flagset.SetOutput(io.Discard)

	// Define flags
	recursive := flagset.BoolP("recursive", "r", false, "recursive search")
	onlyMatches := flagset.BoolP("only-matching", "o", false, "print only matching parts")

	// emulated grep always assumes -E flag by default
	_ = flagset.BoolP("extended-regexp", "E", false, "extended regex")
	color := flagset.String("color", "never", "colorize output (always|never|auto)")

	// Parse flags
	err := flagset.Parse(args)
	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Failed to launch grep: %s", err))
	}

	// Validate color option
	colorMode := ColorMode(*color)
	if colorMode != ColorAlways && colorMode != ColorNever && colorMode != ColorAuto {
		panic(fmt.Sprintf("Codecrafters Internal Error - Invalid color mode: %s", *color))
	}

	// After parsing flags, remaining args are pattern + files...
	remaining := flagset.Args()
	if len(remaining) == 0 {
		panic("Codecrafters Internal Error - Grep is launched with neither pattern nor files")
	}

	pattern := remaining[0]
	files := remaining[1:]
	useStdin := len(files) == 0

	if useStdin {
		shouldEnableHighlighting := (colorMode == ColorAlways) || (colorMode == ColorAuto && launchOptions.EmulateInTTY)

		return searchStdin(pattern, string(launchOptions.Stdin), searchOptions{
			onlyMatches:        *onlyMatches,
			enableHighlighting: shouldEnableHighlighting,
		})
	}

	return searchFiles(pattern, files, fileSearchOptions{
		recursive: *recursive,
	})
}

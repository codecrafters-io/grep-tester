package grep

import (
	"flag"
	"fmt"
	"io"
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

type EmulationOptions struct {
	Stdin        []byte
	EmulateInTTY bool
}

// EmulateGrep provides a simplified interface that mimics grep command behavior
func EmulateGrep(args []string, launchOptions EmulationOptions) Result {
	flagset := flag.NewFlagSet("grep", flag.ContinueOnError)

	// Discard error and output messages, because this isn't being used in the command line
	flagset.SetOutput(io.Discard)

	recursive := flagset.Bool("r", false, "recursive search")
	onlyMatches := flagset.Bool("o", false, "print only matching parts")
	color := flagset.String("color", "never", "colorize output (always|never|auto)")

	// emulated grep always assumes -E flag no matter what (ignored, but accepted)
	_ = flagset.Bool("E", false, "extended regex")

	if err := flagset.Parse(args); err != nil {
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

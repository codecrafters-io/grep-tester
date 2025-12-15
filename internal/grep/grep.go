package grep

import (
	"flag"
	"fmt"
	"io"

	"github.com/codecrafters-io/grep-tester/internal/utils"
)

// Result represents the result of a grep operation
type Result struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
}

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
	colorMode := utils.ColorMode(*color)
	if colorMode != utils.ColorAlways && colorMode != utils.ColorNever && colorMode != utils.ColorAuto {
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
		shouldEnableHighlighting :=
			(colorMode == utils.ColorAlways) ||
				(colorMode == utils.ColorAuto && launchOptions.EmulateInTTY)

		return searchStdin(pattern, string(launchOptions.Stdin), searchOptions{
			onlyMatches:        *onlyMatches,
			enableHighlighting: shouldEnableHighlighting,
		})
	}

	return searchFiles(pattern, files, fileSearchOptions{
		recursive: *recursive,
	})
}

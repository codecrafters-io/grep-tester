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

// EmulateGrep provides a simplified interface that mimics grep command behavior
func EmulateGrep(args []string, stdin []byte) Result {
	// Define the flag set
	// Use flagset instead of iterating through args because
	// args can appear in any order, eg
	// "grep -o -E 'abc'" is the same as "grep -E 'abc' -o"
	flagset := flag.NewFlagSet("grep", flag.ContinueOnError)

	// we aren't using this as a command line tool, disable help messages
	flagset.SetOutput(io.Discard)

	recursive := flagset.Bool("r", false, "recursive search")
	onlyMatches := flagset.Bool("o", false, "print only matching parts")

	// emulated grep always assumes -E flag by default
	_ = flagset.Bool("E", false, "extended regex")

	// Parse flags
	err := flagset.Parse(args)

	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Failed to launch grep: %s", err))
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
		return searchStdin(pattern, string(stdin), searchOptions{
			onlyMatches: *onlyMatches,
		})
	}

	return searchFiles(pattern, files, fileSearchOptions{
		recursive: *recursive,
	})
}

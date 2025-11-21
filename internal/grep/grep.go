package grep

import (
	"flag"
	"fmt"
	"io"
	"strings"
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

type searchOptions struct {
	onlyMatches bool
}

func searchStdin(pattern string, input string, searchOptions searchOptions) Result {
	var stdout []string
	var stderr []string
	exitCode := 0

	matcher := backReferenceMatcher{
		pattern: pattern,
	}

	lines := strings.Split(input, "\n")
	matchCount := 0

	for _, line := range lines {
		matchResult := matcher.match(line)

		if !matchResult.success {
			continue
		}

		matchCount++

		if searchOptions.onlyMatches {
			stdout = append(stdout, matchResult.matchedStrings...)
		} else {
			stdout = append(stdout, line)
		}
	}

	if matchCount == 0 {
		exitCode = 1
	}

	return Result{
		ExitCode: exitCode,
		Stdout:   []byte(strings.Join(stdout, "\n")),
		Stderr:   []byte(strings.Join(stderr, "\n")),
	}
}

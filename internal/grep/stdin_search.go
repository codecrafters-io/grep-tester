package grep

import (
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/utils"
)

type searchOptions struct {
	onlyMatches        bool
	enableHighlighting bool
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

		// Append only matched strings
		if searchOptions.onlyMatches {
			stdout = append(stdout, matchResult.matchedStrings...)
			continue
		}

		// Append highlighted line
		if searchOptions.enableHighlighting {
			stdout = append(stdout, highlightMatches(line, matchResult.matchedStrings))
			continue
		}

		// Append the raw line
		stdout = append(stdout, line)
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

// highlightMatches places ascii sequences
func highlightMatches(line string, matchedStrings []string) string {
	consumedLength := 0
	remaining := line
	result := ""

	for _, match := range matchedStrings {
		matchIdx := strings.Index(remaining, match)
		beforeMatch := remaining[:matchIdx]
		result += beforeMatch + utils.HighlightString(match)
		consumedLength += len(beforeMatch) + len(match)
		remaining = line[consumedLength:]
	}

	return result + remaining
}

package grep

import "strings"

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
		Stdout:   []byte(linesToProgramOutput(stdout, exitCode == 0)),
		Stderr:   []byte(linesToProgramOutput(stderr, true)),
	}
}

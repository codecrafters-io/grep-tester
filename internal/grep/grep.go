package grep

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Result struct {
	ExitCode int
	Stdout   []string
	Stderr   []string
}

type Options struct {
	IgnoreCase     bool
	LineNumber     bool
	Recursive      bool
	InvertMatch    bool
	WholeWord      bool
	Count          bool
	FilesWithMatch bool
	Quiet          bool
	ExtendedRegex  bool
}

type Matcher interface {
	Match(text string) bool
}

type RegexMatcher struct {
	regex *regexp.Regexp
}

func (m *RegexMatcher) Match(text string) bool {
	return m.regex.MatchString(text)
}

// BackrefMatcher handles patterns with backreferences (\1, \2, etc.).
// Go regexp uses the RE2 engine, which doesn't support backreferences out of the box.
// This matcher implements backreferences using pattern expansion and validation.
type BackrefMatcher struct {
	pattern string
}

func (m *BackrefMatcher) Match(text string) bool {
	return matchWithBackreferences(m.pattern, text)
}

func SearchStdin(pattern string, input string, opts Options) Result {
	var stdout []string
	var stderr []string
	exitCode := 0

	// Apply whole word wrapping
	if opts.WholeWord {
		pattern = `\b` + pattern + `\b`
	}

	matcher := createMatcher(pattern)
	if matcher == nil {
		return Result{
			ExitCode: 2,
			Stderr:   []string{fmt.Sprintf("Invalid regex pattern: %s", pattern)},
		}
	}

	lines := strings.Split(input, "\n")
	matchCount := 0

	for _, line := range lines {
		isMatch := matcher.Match(line)
		if opts.InvertMatch {
			isMatch = !isMatch
		}

		if isMatch {
			matchCount++
			if opts.Count {
				continue
			}
			if opts.Quiet {
				// In quiet mode, we exit immediately on first match
				return Result{ExitCode: 0, Stdout: stdout, Stderr: stderr}
			}

			// For stdin, we don't include filename
			stdout = append(stdout, line)
		}
	}

	if opts.Count {
		stdout = append(stdout, fmt.Sprintf("%d", matchCount))
	}

	if matchCount == 0 {
		exitCode = 1
	}

	return Result{
		ExitCode: exitCode,
		Stdout:   stdout,
		Stderr:   stderr,
	}
}

func SearchFiles(pattern string, files []string, opts Options) Result {
	var stdout []string
	var stderr []string
	exitCode := 0

	// Apply whole word wrapping
	if opts.WholeWord {
		pattern = `\b` + pattern + `\b`
	}

	matcher := createMatcher(pattern)
	if matcher == nil {
		return Result{
			ExitCode: 2,
			Stderr:   []string{fmt.Sprintf("Invalid regex pattern: %s", pattern)},
		}
	}

	totalMatches := 0
	hasMultipleFiles := len(files) > 1

	for _, filename := range files {
		if opts.Recursive && isDirectory(filename) {
			matches, out, err := searchDirectory(matcher, filename, opts, true)
			totalMatches += matches
			stdout = append(stdout, out...)
			stderr = append(stderr, err...)
		} else {
			matches, out, err := searchFile(matcher, filename, opts, hasMultipleFiles)
			totalMatches += matches
			stdout = append(stdout, out...)
			stderr = append(stderr, err...)
		}
	}

	if totalMatches == 0 {
		exitCode = 1
	}

	return Result{
		ExitCode: exitCode,
		Stdout:   stdout,
		Stderr:   stderr,
	}
}

func hasBackreferences(pattern string) bool {
	return regexp.MustCompile(`\\[1-9]`).MatchString(pattern)
}

func createMatcher(pattern string) Matcher {
	if hasBackreferences(pattern) {
		return &BackrefMatcher{
			pattern: pattern,
		}
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}

	return &RegexMatcher{
		regex: regex,
	}
}

func searchDirectory(matcher Matcher, dirname string, opts Options, hasMultipleFiles bool) (int, []string, []string) {
	var stdout []string
	var stderr []string
	totalMatches := 0

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			matches, out, errOut := searchFile(matcher, path, opts, hasMultipleFiles)
			totalMatches += matches
			stdout = append(stdout, out...)
			stderr = append(stderr, errOut...)
		}
		return nil
	})

	if err != nil {
		stderr = append(stderr, fmt.Sprintf("Error walking directory %s: %v", dirname, err))
	}

	return totalMatches, stdout, stderr
}

func searchFile(matcher Matcher, filename string, opts Options, hasMultipleFiles bool) (int, []string, []string) {
	var stdout []string
	var stderr []string

	file, err := os.Open(filename)
	if err != nil {
		stderr = append(stderr, fmt.Sprintf("Error opening file %s: %v", filename, err))
		return 0, stdout, stderr
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matchCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		isMatch := matcher.Match(line)
		if opts.InvertMatch {
			isMatch = !isMatch
		}

		if isMatch {
			matchCount++
			if opts.Count {
				continue
			}
			if opts.FilesWithMatch {
				stdout = append(stdout, filename)
				break
			}
			if opts.Quiet {
				// In quiet mode, we exit immediately on first match
				return matchCount, stdout, stderr
			}

			// Format output line
			var parts []string
			if filename != "" && hasMultipleFiles {
				parts = append(parts, filename)
			}

			if len(parts) > 0 {
				stdout = append(stdout, fmt.Sprintf("%s:%s", strings.Join(parts, ":"), line))
			} else {
				stdout = append(stdout, line)
			}
		}
	}

	if opts.Count {
		if hasMultipleFiles {
			stdout = append(stdout, fmt.Sprintf("%s:%d", filename, matchCount))
		} else {
			stdout = append(stdout, fmt.Sprintf("%d", matchCount))
		}
	}

	return matchCount, stdout, stderr
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func matchWithBackreferences(pattern, text string) bool {
	// First, create a version of the pattern where backreferences are replaced with (.*)
	// This allows us to use Go's regex engine to find potential matches
	regexPattern := regexp.MustCompile(`\\[1-9]`).ReplaceAllString(pattern, `(.*)`)

	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return false
	}

	// Find all potential matches
	allMatches := regex.FindAllStringSubmatch(text, -1)

	// For each potential match, check if backreferences are satisfied
	for _, match := range allMatches {
		if validateBackreferences(pattern, match) {
			return true
		}
	}

	return false
}

func validateBackreferences(pattern string, match []string) bool {
	if len(match) < 2 {
		return false
	}

	// match[0] is the full match
	// match[1], match[2], etc. are the captured groups
	fullMatch := match[0]

	// Replace backreferences with the actual captured values
	expandedPattern := pattern
	for i := 1; i < len(match); i++ {
		backref := fmt.Sprintf("\\%d", i)
		expandedPattern = strings.ReplaceAll(expandedPattern, backref, regexp.QuoteMeta(match[i]))
	}

	// Now check if this expanded pattern matches the full match
	regex, err := regexp.Compile("^" + expandedPattern + "$")
	if err != nil {
		return false
	}

	return regex.MatchString(fullMatch)
}

package grep

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	if len(args) == 0 {
		return Result{
			ExitCode: 2,
			Stderr:   []byte("Usage: grep [options] pattern [files...]\n"),
		}
	}

	// Parse arguments
	var recursive bool
	var pattern string
	var files []string
	var useStdin = true

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-r":
			recursive = true
		case "-E":
			// Extended regex flag - ignore for now as we always use extended
			if i+1 < len(args) {
				i++
				pattern = args[i]
			}
		default:
			if pattern == "" {
				pattern = arg
			} else {
				files = append(files, arg)
				useStdin = false
			}
		}
	}

	if pattern == "" {
		return Result{
			ExitCode: 2,
			Stderr:   []byte("No pattern specified\n"),
		}
	}

	if useStdin {
		return searchStdin(pattern, string(stdin))
	} else {
		return searchFiles(pattern, files, options{recursive: recursive})
	}
}

type options struct {
	recursive bool
}

// BackrefMatcher handles patterns with backreferences (\1, \2, etc.).
// Go regexp uses the RE2 engine, which doesn't support backreferences out of the box.
// This matcher implements backreferences using pattern expansion and validation.
type backrefMatcher struct {
	pattern string
}

func (m *backrefMatcher) match(text string) bool {
	return matchWithBackreferences(m.pattern, text)
}

func searchStdin(pattern string, input string) Result {
	var stdout []string
	var stderr []string
	exitCode := 0

	matcher := createMatcher(pattern)

	lines := strings.Split(input, "\n")
	matchCount := 0

	for _, line := range lines {
		isMatch := matcher.match(line)

		if isMatch {
			matchCount++
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

func searchFiles(pattern string, files []string, opts options) Result {
	var stdout []string
	var stderr []string
	exitCode := 0

	matcher := createMatcher(pattern)

	totalMatches := 0
	hasMultipleFiles := len(files) > 1

	for _, filename := range files {
		if opts.recursive && isDirectory(filename) {
			matches, out, err := searchDirectory(matcher, filename, true)
			totalMatches += matches
			stdout = append(stdout, out...)
			stderr = append(stderr, err...)
		} else {
			matches, out, err := searchFile(matcher, filename, hasMultipleFiles)
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
		Stdout:   []byte(strings.Join(stdout, "\n")),
		Stderr:   []byte(strings.Join(stderr, "\n")),
	}
}

func createMatcher(pattern string) *backrefMatcher {
	return &backrefMatcher{
		pattern: pattern,
	}
}

func searchDirectory(matcher *backrefMatcher, dirname string, hasMultipleFiles bool) (int, []string, []string) {
	var stdout []string
	var stderr []string
	totalMatches := 0

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			matches, out, errOut := searchFile(matcher, path, hasMultipleFiles)
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

func searchFile(matcher *backrefMatcher, filename string, hasMultipleFiles bool) (int, []string, []string) {
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

		isMatch := matcher.match(line)

		if isMatch {
			matchCount++
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
	// Handle patterns without backreferences first
	if !regexp.MustCompile(`\\[1-9]`).MatchString(pattern) {
		regex := regexp.MustCompile(pattern)
		return regex.MatchString(text)
	}

	// First, create a version of the pattern where backreferences are replaced with (.*)
	// This allows us to use Go's regex engine to find potential matches
	regexPattern := regexp.MustCompile(`\\[1-9]`).ReplaceAllString(pattern, `(.*)`)
	regex := regexp.MustCompile(regexPattern)

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
	regex := regexp.MustCompile("^" + expandedPattern + "$")

	return regex.MatchString(fullMatch)
}

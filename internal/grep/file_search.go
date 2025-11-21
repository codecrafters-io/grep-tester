package grep

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type fileSearchOptions struct {
	recursive bool
}

func searchFiles(pattern string, files []string, opts fileSearchOptions) Result {
	var stdout []string
	var stderr []string
	exitCode := 0

	matcher := newBackReferenceMatcher(pattern)

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

func searchDirectory(matcher *backReferenceMatcher, dirname string, hasMultipleFiles bool) (int, []string, []string) {
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

func searchFile(matcher *backReferenceMatcher, filename string, hasMultipleFiles bool) (int, []string, []string) {
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

		if isMatch.success {
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

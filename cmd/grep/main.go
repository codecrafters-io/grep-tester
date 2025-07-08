package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Matcher interface {
	Match(text string) bool
}

type RegexMatcher struct {
	regex *regexp.Regexp
}

func (m *RegexMatcher) Match(text string) bool {
	return m.regex.MatchString(text)
}

type BackrefMatcher struct {
	pattern    string
	ignoreCase bool
}

func (m *BackrefMatcher) Match(text string) bool {
	return matchWithBackreferences(m.pattern, text, m.ignoreCase)
}

type Config struct {
	Pattern        string
	Files          []string
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

func main() {
	config := parseArgs()

	if config.Pattern == "" {
		fmt.Fprintf(os.Stderr, "Usage: grep [OPTIONS] PATTERN [FILE...]\n")
		os.Exit(1)
	}

	pattern := config.Pattern
	if config.WholeWord {
		pattern = `\b` + pattern + `\b`
	}

	matcher := createMatcher(pattern, config)

	if len(config.Files) == 0 {
		searchStdin(matcher, config)
	} else {
		searchFiles(matcher, config)
	}
}

func hasBackreferences(pattern string) bool {
	return regexp.MustCompile(`\\[1-9]`).MatchString(pattern)
}

func createMatcher(pattern string, config Config) Matcher {
	if hasBackreferences(pattern) {
		return &BackrefMatcher{
			pattern:    pattern,
			ignoreCase: config.IgnoreCase,
		}
	}

	var regex *regexp.Regexp
	var err error

	if config.IgnoreCase {
		regex, err = regexp.Compile("(?i)" + pattern)
	} else {
		regex, err = regexp.Compile(pattern)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid regex pattern: %v\n", err)
		os.Exit(1)
	}

	return &RegexMatcher{
		regex: regex,
	}
}

func parseArgs() Config {
	var config Config

	flag.BoolVar(&config.IgnoreCase, "i", false, "Ignore case")
	flag.BoolVar(&config.LineNumber, "n", false, "Show line numbers")
	flag.BoolVar(&config.Recursive, "r", false, "Recursive search")
	flag.BoolVar(&config.InvertMatch, "v", false, "Invert match")
	flag.BoolVar(&config.WholeWord, "w", false, "Match whole words")
	flag.BoolVar(&config.Count, "c", false, "Count matches")
	flag.BoolVar(&config.FilesWithMatch, "l", false, "Show files with matches")
	flag.BoolVar(&config.Quiet, "q", false, "Quiet mode")
	flag.BoolVar(&config.ExtendedRegex, "E", false, "Extended regex")

	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		config.Pattern = args[0]
		config.Files = args[1:]
	}

	return config
}

func searchStdin(matcher Matcher, config Config) {
	scanner := bufio.NewScanner(os.Stdin)
	matchCount := processScanner(matcher, scanner, config, "", false)

	if config.Count {
		fmt.Println(matchCount)
	}

	if matchCount == 0 {
		os.Exit(1)
	}
}

func searchFiles(matcher Matcher, config Config) {
	totalMatches := 0
	hasMultipleFiles := len(config.Files) > 1

	for _, filename := range config.Files {
		if config.Recursive && isDirectory(filename) {
			totalMatches += searchDirectory(matcher, filename, config, true)
		} else {
			totalMatches += searchFile(matcher, filename, config, hasMultipleFiles)
		}
	}

	if totalMatches == 0 {
		os.Exit(1)
	}
}

func searchDirectory(matcher Matcher, dirname string, config Config, hasMultipleFiles bool) int {
	totalMatches := 0

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			totalMatches += searchFile(matcher, path, config, hasMultipleFiles)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", dirname, err)
	}

	return totalMatches
}

func searchFile(matcher Matcher, filename string, config Config, hasMultipleFiles bool) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matchCount := processScanner(matcher, scanner, config, filename, hasMultipleFiles)

	if config.Count {
		if hasMultipleFiles {
			fmt.Printf("%s:%d\n", filename, matchCount)
		} else {
			fmt.Println(matchCount)
		}
	}

	return matchCount
}

func processScanner(matcher Matcher, scanner *bufio.Scanner, config Config, filename string, hasMultipleFiles bool) int {
	lineNum := 0
	matchCount := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		isMatch := matcher.Match(line)
		if config.InvertMatch {
			isMatch = !isMatch
		}

		if isMatch {
			matchCount++
			if config.Count {
				continue
			}
			if config.FilesWithMatch {
				fmt.Println(filename)
				return matchCount
			}
			if config.Quiet {
				os.Exit(0)
			}
			printMatch(filename, line, lineNum, config, hasMultipleFiles)
		}
	}
	return matchCount
}

func printMatch(filename, line string, lineNum int, config Config, hasMultipleFiles bool) {
	var parts []string

	if filename != "" && hasMultipleFiles {
		parts = append(parts, filename)
	}

	if config.LineNumber {
		parts = append(parts, fmt.Sprintf("%d", lineNum))
	}

	if len(parts) > 0 {
		fmt.Printf("%s:%s\n", strings.Join(parts, ":"), line)
	} else {
		fmt.Println(line)
	}
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func matchWithBackreferences(pattern, text string, ignoreCase bool) bool {
	// First, create a version of the pattern where backreferences are replaced with (.*)
	// This allows us to use Go's regex engine to find potential matches
	regexPattern := regexp.MustCompile(`\\[1-9]`).ReplaceAllString(pattern, `(.*)`)

	var regex *regexp.Regexp
	var err error

	if ignoreCase {
		regex, err = regexp.Compile("(?i)" + regexPattern)
	} else {
		regex, err = regexp.Compile(regexPattern)
	}

	if err != nil {
		return false
	}

	// Find all potential matches
	allMatches := regex.FindAllStringSubmatch(text, -1)

	// For each potential match, check if backreferences are satisfied
	for _, match := range allMatches {
		if validateBackreferencesNew(pattern, match, ignoreCase) {
			return true
		}
	}

	return false
}

func validateBackreferencesNew(pattern string, match []string, ignoreCase bool) bool {
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
	var regex *regexp.Regexp
	var err error

	if ignoreCase {
		regex, err = regexp.Compile("(?i)^" + expandedPattern + "$")
	} else {
		regex, err = regexp.Compile("^" + expandedPattern + "$")
	}

	if err != nil {
		return false
	}

	return regex.MatchString(fullMatch)
}

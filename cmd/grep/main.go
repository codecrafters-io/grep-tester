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

	var regex *regexp.Regexp
	var err error
	
	pattern := config.Pattern
	if config.WholeWord {
		pattern = `\b` + pattern + `\b`
	}
	
	if config.IgnoreCase {
		regex, err = regexp.Compile("(?i)" + pattern)
	} else {
		regex, err = regexp.Compile(pattern)
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid regex pattern: %v\n", err)
		os.Exit(1)
	}

	if len(config.Files) == 0 {
		searchStdin(regex, config)
	} else {
		searchFiles(regex, config)
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

func searchStdin(regex *regexp.Regexp, config Config) {
	scanner := bufio.NewScanner(os.Stdin)
	lineNum := 0
	matchCount := 0
	
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		matches := regex.MatchString(line)
		if config.InvertMatch {
			matches = !matches
		}
		
		if matches {
			matchCount++
			if config.Count {
				continue
			}
			if config.Quiet {
				os.Exit(0)
			}
			
			printMatch("", line, lineNum, config)
		}
	}
	
	if config.Count {
		fmt.Println(matchCount)
	}
	
	if matchCount == 0 {
		os.Exit(1)
	}
}

func searchFiles(regex *regexp.Regexp, config Config) {
	totalMatches := 0
	hasMultipleFiles := len(config.Files) > 1
	
	for _, filename := range config.Files {
		if config.Recursive && isDirectory(filename) {
			totalMatches += searchDirectory(regex, filename, config, hasMultipleFiles)
		} else {
			totalMatches += searchFile(regex, filename, config, hasMultipleFiles)
		}
	}
	
	if totalMatches == 0 {
		os.Exit(1)
	}
}

func searchDirectory(regex *regexp.Regexp, dirname string, config Config, hasMultipleFiles bool) int {
	totalMatches := 0
	
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() {
			totalMatches += searchFile(regex, path, config, hasMultipleFiles)
		}
		return nil
	})
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory %s: %v\n", dirname, err)
	}
	
	return totalMatches
}

func searchFile(regex *regexp.Regexp, filename string, config Config, hasMultipleFiles bool) int {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
		return 0
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineNum := 0
	matchCount := 0
	
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		matches := regex.MatchString(line)
		if config.InvertMatch {
			matches = !matches
		}
		
		if matches {
			matchCount++
			if config.Count {
				continue
			}
			if config.FilesWithMatch {
				fmt.Println(filename)
				break
			}
			if config.Quiet {
				os.Exit(0)
			}
			
			printMatch(filename, line, lineNum, config)
		}
	}
	
	if config.Count {
		if hasMultipleFiles {
			fmt.Printf("%s:%d\n", filename, matchCount)
		} else {
			fmt.Println(matchCount)
		}
	}
	
	return matchCount
}

func printMatch(filename, line string, lineNum int, config Config) {
	var parts []string
	
	if filename != "" {
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
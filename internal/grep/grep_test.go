package grep

import (
	"os"
	"path/filepath"
	"testing"
)

type StdinTestCase struct {
	name     string
	pattern  string
	input    string
	options  Options
	expected Result
}

type FileTestCase struct {
	name     string
	pattern  string
	files    []string
	options  Options
	expected Result
}

type MatcherTestCase struct {
	name     string
	pattern  string
	text     string
	expected bool
}

func runStdinTests(t *testing.T, tests []StdinTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, tt.options)

			if result.ExitCode != tt.expected.ExitCode {
				t.Errorf("Expected exit code %d, got %d", tt.expected.ExitCode, result.ExitCode)
			}

			if len(result.Stdout) != len(tt.expected.Stdout) {
				t.Errorf("Expected %d stdout lines, got %d", len(tt.expected.Stdout), len(result.Stdout))
			}

			for i, expectedLine := range tt.expected.Stdout {
				if i >= len(result.Stdout) || result.Stdout[i] != expectedLine {
					t.Errorf("Expected stdout line %d to be %q, got %q", i, expectedLine, result.Stdout[i])
				}
			}
		})
	}
}

func runFileTests(t *testing.T, tests []FileTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchFiles(tt.pattern, tt.files, tt.options)

			if result.ExitCode != tt.expected.ExitCode {
				t.Errorf("Expected exit code %d, got %d", tt.expected.ExitCode, result.ExitCode)
			}

			if len(result.Stdout) != len(tt.expected.Stdout) {
				t.Errorf("Expected %d stdout lines, got %d", len(tt.expected.Stdout), len(result.Stdout))
			}

			for i, expectedLine := range tt.expected.Stdout {
				if i >= len(result.Stdout) || result.Stdout[i] != expectedLine {
					t.Errorf("Expected stdout line %d to be %q, got %q", i, expectedLine, result.Stdout[i])
				}
			}
		})
	}
}

func runMatcherTests(t *testing.T, tests []MatcherTestCase, createMatcherFunc func(string) Matcher) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := createMatcherFunc(tt.pattern)
			if matcher == nil {
				t.Fatalf("Failed to create matcher for pattern %q", tt.pattern)
			}

			result := matcher.Match(tt.text)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSearchStdin(t *testing.T) {
	tests := []StdinTestCase{
		// Basic matching
		{
			name:     "literal character match",
			pattern:  "d",
			input:    "dog",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"dog"}},
		},
		{
			name:     "literal character no match",
			pattern:  "f",
			input:    "dog",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "digit match",
			pattern:  "\\d",
			input:    "123",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"123"}},
		},
		{
			name:     "digit no match",
			pattern:  "\\d",
			input:    "apple",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "alphanumeric match",
			pattern:  "\\w",
			input:    "apple",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"apple"}},
		},
		{
			name:     "alphanumeric no match",
			pattern:  "\\w",
			input:    "+-$×€÷",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Character groups
		{
			name:     "positive character group match",
			pattern:  "[pineapple]",
			input:    "n",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"n"}},
		},
		{
			name:     "positive character group no match",
			pattern:  "[cdfghijklm]",
			input:    "strawberry",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "negative character group match",
			pattern:  "[^xyz]",
			input:    "apple",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"apple"}},
		},
		{
			name:     "negative character group no match",
			pattern:  "[^anb]",
			input:    "banana",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Anchors
		{
			name:     "start anchor match",
			pattern:  "^log",
			input:    "log",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"log"}},
		},
		{
			name:     "start anchor no match",
			pattern:  "^log",
			input:    "slog",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "end anchor match",
			pattern:  "cat$",
			input:    "cat",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"cat"}},
		},
		{
			name:     "end anchor no match",
			pattern:  "cat$",
			input:    "cats",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Quantifiers
		{
			name:     "one or more match",
			pattern:  "ca+t",
			input:    "cat",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"cat"}},
		},
		{
			name:     "one or more no match",
			pattern:  "ca+t",
			input:    "act",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "zero or one match",
			pattern:  "ca?t",
			input:    "cat",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"cat"}},
		},
		{
			name:     "zero or one alternate match",
			pattern:  "ca?t",
			input:    "act",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"act"}},
		},
		{
			name:     "zero or one no match",
			pattern:  "ca?t",
			input:    "dog",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Wildcard
		{
			name:     "wildcard match",
			pattern:  "c.t",
			input:    "cat",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"cat"}},
		},
		{
			name:     "wildcard no match",
			pattern:  "c.t",
			input:    "car",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "wildcard with quantifier match",
			pattern:  "g.+gol",
			input:    "goøö0Ogol",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"goøö0Ogol"}},
		},
		{
			name:     "wildcard with quantifier no match",
			pattern:  "g.+gol",
			input:    "gol",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Alternation
		{
			name:     "alternation match",
			pattern:  "a (cat|dog)",
			input:    "a cat",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"a cat"}},
		},
		{
			name:     "alternation no match",
			pattern:  "a (cat|dog)",
			input:    "a cow",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Single backreferences
		{
			name:     "basic backreference match",
			pattern:  "(cat) and \\1",
			input:    "cat and cat",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"cat and cat"}},
		},
		{
			name:     "basic backreference no match",
			pattern:  "(cat) and \\1",
			input:    "cat and dog",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "complex backreference match",
			pattern:  "([abcd]+) is \\1, not [^xyz]+",
			input:    "abcd is abcd, not efg",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"abcd is abcd, not efg"}},
		},
		{
			name:     "complex backreference no match",
			pattern:  "([abcd]+) is \\1, not [^xyz]+",
			input:    "efgh is efgh, not efg",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "anchored backreference match",
			pattern:  "^(\\w+) starts and ends with \\1$",
			input:    "this starts and ends with this",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"this starts and ends with this"}},
		},
		{
			name:     "anchored backreference no match",
			pattern:  "^(this) starts and ends with \\1$",
			input:    "that starts and ends with this",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Multiple backreferences
		{
			name:     "multiple backreferences match",
			pattern:  "(\\d+) (\\w+) squares and \\1 \\2 circles",
			input:    "3 red squares and 3 red circles",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"3 red squares and 3 red circles"}},
		},
		{
			name:     "multiple backreferences no match",
			pattern:  "(\\d+) (\\w+) squares and \\1 \\2 circles",
			input:    "3 red squares and 4 red circles",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Nested backreferences
		{
			name:     "nested backreferences match",
			pattern:  "('(cat) and \\2') is the same as \\1",
			input:    "'cat and cat' is the same as 'cat and cat'",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"'cat and cat' is the same as 'cat and cat'"}},
		},
		{
			name:     "nested backreferences no match",
			pattern:  "('(cat) and \\2') is the same as \\1",
			input:    "'cat and cat' is the same as 'cat and dog'",
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
	}

	runStdinTests(t, tests)
}

func TestSearchFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("lemon\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	fruitsFile := filepath.Join(tempDir, "fruits.txt")
	err = os.WriteFile(fruitsFile, []byte("banana\nblueberry\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write fruits file: %v", err)
	}

	vegetablesFile := filepath.Join(tempDir, "vegetables.txt")
	err = os.WriteFile(vegetablesFile, []byte("broccoli\ncarrot\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write vegetables file: %v", err)
	}

	// Create directory structure for recursive testing
	dirPath := filepath.Join(tempDir, "dir")
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	dirFruitsFile := filepath.Join(dirPath, "fruits.txt")
	err = os.WriteFile(dirFruitsFile, []byte("pear\nstrawberry\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dir fruits file: %v", err)
	}

	subdirPath := filepath.Join(dirPath, "subdir")
	err = os.MkdirAll(subdirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	subdirVegetablesFile := filepath.Join(subdirPath, "vegetables.txt")
	err = os.WriteFile(subdirVegetablesFile, []byte("celery\ncarrot\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write subdir vegetables file: %v", err)
	}

	dirVegetablesFile := filepath.Join(dirPath, "vegetables.txt")
	err = os.WriteFile(dirVegetablesFile, []byte("cucumber\ncorn\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dir vegetables file: %v", err)
	}

	tests := []FileTestCase{
		// Single file
		{
			name:     "single file match",
			pattern:  "le.*",
			files:    []string{testFile},
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{"lemon"}},
		},
		{
			name:     "single file no match",
			pattern:  "broccoli",
			files:    []string{testFile},
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},

		// Multiple files
		{
			name:    "multiple files match",
			pattern: "b.*$",
			files:   []string{fruitsFile, vegetablesFile},
			options: Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{
				fruitsFile + ":banana",
				fruitsFile + ":blueberry",
				vegetablesFile + ":broccoli",
			}},
		},
		{
			name:     "multiple files no match",
			pattern:  "missing_fruit",
			files:    []string{fruitsFile, vegetablesFile},
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
		{
			name:     "specific file match",
			pattern:  "carrot",
			files:    []string{fruitsFile, vegetablesFile},
			options:  Options{ExtendedRegex: true},
			expected: Result{ExitCode: 0, Stdout: []string{vegetablesFile + ":carrot"}},
		},

		// Recursive search
		{
			name:    "recursive search match",
			pattern: ".*er",
			files:   []string{dirPath},
			options: Options{ExtendedRegex: true, Recursive: true},
			expected: Result{ExitCode: 0, Stdout: []string{
				dirFruitsFile + ":strawberry",
				subdirVegetablesFile + ":celery",
				dirVegetablesFile + ":cucumber",
			}},
		},
		{
			name:     "recursive search no match",
			pattern:  "missing_fruit",
			files:    []string{dirPath},
			options:  Options{ExtendedRegex: true, Recursive: true},
			expected: Result{ExitCode: 1, Stdout: []string{}},
		},
	}

	runFileTests(t, tests)
}

func TestBackrefMatcher(t *testing.T) {
	tests := []MatcherTestCase{
		{
			name:     "simple backreference match",
			pattern:  "(test) \\1",
			text:     "test test",
			expected: true,
		},
		{
			name:     "simple backreference no match",
			pattern:  "(test) \\1",
			text:     "test different",
			expected: false,
		},
		{
			name:     "multiple backreferences match",
			pattern:  "(\\w+) (\\w+) \\1 \\2",
			text:     "hello world hello world",
			expected: true,
		},
		{
			name:     "multiple backreferences no match",
			pattern:  "(\\w+) (\\w+) \\1 \\2",
			text:     "hello world hello different",
			expected: false,
		},
	}

	createBackrefMatcher := func(pattern string) Matcher {
		return &BackrefMatcher{pattern: pattern}
	}

	runMatcherTests(t, tests, createBackrefMatcher)
}

func TestRegexMatcher(t *testing.T) {
	tests := []MatcherTestCase{
		{
			name:     "simple regex match",
			pattern:  "test",
			text:     "this is a test",
			expected: true,
		},
		{
			name:     "simple regex no match",
			pattern:  "test",
			text:     "this is different",
			expected: false,
		},
		{
			name:     "regex with quantifier match",
			pattern:  "a+",
			text:     "aaa",
			expected: true,
		},
		{
			name:     "regex with quantifier no match",
			pattern:  "a+",
			text:     "bbb",
			expected: false,
		},
	}

	runMatcherTests(t, tests, createMatcher)
}

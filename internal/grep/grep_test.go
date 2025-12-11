package grep

import (
	"os"
	"path/filepath"
	"testing"
)

type StdinTestCase struct {
	name        string
	pattern     string
	input       string
	expected    Result
	onlyMatches bool
}

type FileTestCase struct {
	name      string
	pattern   string
	files     []string
	expected  Result
	recursive bool
}

func runStdinTests(t *testing.T, tests []StdinTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arguments := []string{}

			if tt.onlyMatches {
				arguments = append(arguments, "-o")
			}

			arguments = append(arguments, tt.pattern)

			result := EmulateGrep(arguments, []byte(tt.input))

			if result.ExitCode != tt.expected.ExitCode {
				t.Errorf("Expected exit code %d, got %d", tt.expected.ExitCode, result.ExitCode)
			}

			expectedStr := string(tt.expected.Stdout)
			actualStr := string(result.Stdout)

			if actualStr != expectedStr {
				t.Errorf("Expected stdout %q, got %q", expectedStr, actualStr)
			}
		})
	}
}

func runFileTests(t *testing.T, tests []FileTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{}
			if tt.recursive {
				args = append(args, "-r")
			}
			args = append(args, tt.pattern)
			args = append(args, tt.files...)
			result := EmulateGrep(args, []byte{})

			if result.ExitCode != tt.expected.ExitCode {
				t.Errorf("Expected exit code %d, got %d", tt.expected.ExitCode, result.ExitCode)
			}

			expectedStr := string(tt.expected.Stdout)
			actualStr := string(result.Stdout)

			if actualStr != expectedStr {
				t.Errorf("Expected stdout %q, got %q", expectedStr, actualStr)
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
			expected: Result{ExitCode: 0, Stdout: []byte("dog\n")},
		},
		{
			name:     "literal character no match",
			pattern:  "f",
			input:    "dog",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "digit match",
			pattern:  "\\d",
			input:    "123",
			expected: Result{ExitCode: 0, Stdout: []byte("123\n")},
		},
		{
			name:     "digit no match",
			pattern:  "\\d",
			input:    "apple",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "alphanumeric match",
			pattern:  "\\w",
			input:    "apple",
			expected: Result{ExitCode: 0, Stdout: []byte("apple\n")},
		},
		{
			name:     "alphanumeric no match",
			pattern:  "\\w",
			input:    "+-$×€÷",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Character groups
		{
			name:     "positive character group match",
			pattern:  "[pineapple]",
			input:    "n",
			expected: Result{ExitCode: 0, Stdout: []byte("n\n")},
		},
		{
			name:     "positive character group no match",
			pattern:  "[cdfghijklm]",
			input:    "strawberry",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "negative character group match",
			pattern:  "[^xyz]",
			input:    "apple",
			expected: Result{ExitCode: 0, Stdout: []byte("apple\n")},
		},
		{
			name:     "negative character group no match",
			pattern:  "[^anb]",
			input:    "banana",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Anchors
		{
			name:     "start anchor match",
			pattern:  "^log",
			input:    "log",
			expected: Result{ExitCode: 0, Stdout: []byte("log\n")},
		},
		{
			name:     "start anchor no match",
			pattern:  "^log",
			input:    "slog",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "end anchor match",
			pattern:  "cat$",
			input:    "cat",
			expected: Result{ExitCode: 0, Stdout: []byte("cat\n")},
		},
		{
			name:     "end anchor no match",
			pattern:  "cat$",
			input:    "cats",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Quantifiers
		{
			name:     "one or more match",
			pattern:  "ca+t",
			input:    "cat",
			expected: Result{ExitCode: 0, Stdout: []byte("cat\n")},
		},
		{
			name:     "one or more no match",
			pattern:  "ca+t",
			input:    "act",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "zero or one match",
			pattern:  "ca?t",
			input:    "cat",
			expected: Result{ExitCode: 0, Stdout: []byte("cat\n")},
		},
		{
			name:     "zero or one alternate match",
			pattern:  "ca?t",
			input:    "act",
			expected: Result{ExitCode: 0, Stdout: []byte("act\n")},
		},
		{
			name:     "zero or one no match",
			pattern:  "ca?t",
			input:    "dog",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Wildcard
		{
			name:     "wildcard match",
			pattern:  "c.t",
			input:    "cat",
			expected: Result{ExitCode: 0, Stdout: []byte("cat\n")},
		},
		{
			name:     "wildcard no match",
			pattern:  "c.t",
			input:    "car",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "wildcard with quantifier match",
			pattern:  "g.+gol",
			input:    "goøö0Ogol",
			expected: Result{ExitCode: 0, Stdout: []byte("goøö0Ogol\n")},
		},
		{
			name:     "wildcard with quantifier no match",
			pattern:  "g.+gol",
			input:    "gol",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Alternation
		{
			name:     "alternation match",
			pattern:  "a (cat|dog)",
			input:    "a cat",
			expected: Result{ExitCode: 0, Stdout: []byte("a cat\n")},
		},
		{
			name:     "alternation no match",
			pattern:  "a (cat|dog)",
			input:    "a cow",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Single backreferences
		{
			name:     "basic backreference match",
			pattern:  "(cat) and \\1",
			input:    "cat and cat",
			expected: Result{ExitCode: 0, Stdout: []byte("cat and cat\n")},
		},
		{
			name:     "basic backreference no match",
			pattern:  "(cat) and \\1",
			input:    "cat and dog",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "complex backreference match",
			pattern:  "([abcd]+) is \\1, not [^xyz]+",
			input:    "abcd is abcd, not efg",
			expected: Result{ExitCode: 0, Stdout: []byte("abcd is abcd, not efg\n")},
		},
		{
			name:     "complex backreference no match",
			pattern:  "([abcd]+) is \\1, not [^xyz]+",
			input:    "efgh is efgh, not efg",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "anchored backreference match",
			pattern:  "^(\\w+) starts and ends with \\1$",
			input:    "this starts and ends with this",
			expected: Result{ExitCode: 0, Stdout: []byte("this starts and ends with this\n")},
		},
		{
			name:     "anchored backreference no match",
			pattern:  "^(this) starts and ends with \\1$",
			input:    "that starts and ends with this",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Multiple backreferences
		{
			name:     "multiple backreferences match",
			pattern:  "(\\d+) (\\w+) squares and \\1 \\2 circles",
			input:    "3 red squares and 3 red circles",
			expected: Result{ExitCode: 0, Stdout: []byte("3 red squares and 3 red circles\n")},
		},
		{
			name:     "multiple backreferences no match",
			pattern:  "(\\d+) (\\w+) squares and \\1 \\2 circles",
			input:    "3 red squares and 4 red circles",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Nested backreferences
		{
			name:     "nested backreferences match",
			pattern:  "('(cat) and \\2') is the same as \\1",
			input:    "'cat and cat' is the same as 'cat and cat'",
			expected: Result{ExitCode: 0, Stdout: []byte("'cat and cat' is the same as 'cat and cat'\n")},
		},
		{
			name:     "nested backreferences no match",
			pattern:  "('(cat) and \\2') is the same as \\1",
			input:    "'cat and cat' is the same as 'cat and dog'",
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		// Multiple backreferences with nesting
		{
			name:     "multiple backreferences with nesting_pass_1",
			pattern:  "('(cat) and \\2') or ('(dog) and \\4') is the same as \\3",
			input:    "'cat and cat' or 'dog and dog' is the same as 'dog and dog'",
			expected: Result{ExitCode: 0, Stdout: []byte("'cat and cat' or 'dog and dog' is the same as 'dog and dog'\n")},
		},
		{
			name:     "multiple backreferences with nesting_pass_2",
			pattern:  "('(cat) and \\2') or ('(cat) and \\4') is the same as \\3",
			input:    "'cat and cat' or 'cat and cat' is the same as 'cat and cat'",
			expected: Result{ExitCode: 0, Stdout: []byte("'cat and cat' or 'cat and cat' is the same as 'cat and cat'\n")},
		},
		{
			name:     "multiple backreferences with nesting_fail",
			pattern:  "('(cat) and \\2') or ('(cat) and \\4') is the same as \\3",
			input:    "'cat and cat' or 'cat and cat' is the same as 'dog and dog'",
			expected: Result{ExitCode: 1},
		},
		// Only Matches -o flag
		{
			name:        "single match",
			pattern:     "\\d",
			input:       "a1b",
			onlyMatches: true,
			expected:    Result{ExitCode: 0, Stdout: []byte("1\n")},
		},
		{
			name:        "multiple_matches_single_line",
			pattern:     "\\d",
			input:       "a1b2c3d4",
			onlyMatches: true,
			expected:    Result{ExitCode: 0, Stdout: []byte("1\n2\n3\n4\n")},
		},
		{
			name:        "multiple_matches_multiple_lines",
			pattern:     "\\d\\d",
			input:       "a1b22c3d4\n\n234\n2255\n",
			onlyMatches: true,
			expected:    Result{ExitCode: 0, Stdout: []byte("22\n23\n22\n55\n")},
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
			expected: Result{ExitCode: 0, Stdout: []byte("lemon\n")},
		},
		{
			name:     "single file no match",
			pattern:  "broccoli",
			files:    []string{testFile},
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},

		// Multiple files
		{
			name:    "multiple files match",
			pattern: "b.*$",
			files:   []string{fruitsFile, vegetablesFile},
			expected: Result{ExitCode: 0, Stdout: []byte(
				fruitsFile + ":banana" + "\n" +
					fruitsFile + ":blueberry" + "\n" +
					vegetablesFile + ":broccoli\n",
			)},
		},
		{
			name:     "multiple files no match",
			pattern:  "missing_fruit",
			files:    []string{fruitsFile, vegetablesFile},
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
		{
			name:     "specific file match",
			pattern:  "carrot",
			files:    []string{fruitsFile, vegetablesFile},
			expected: Result{ExitCode: 0, Stdout: []byte(vegetablesFile + ":carrot\n")},
		},

		// Recursive search
		{
			name:    "recursive search match",
			pattern: ".*er",
			files:   []string{dirPath},
			expected: Result{ExitCode: 0, Stdout: []byte(
				dirFruitsFile + ":strawberry" + "\n" +
					subdirVegetablesFile + ":celery" + "\n" +
					dirVegetablesFile + ":cucumber\n",
			)},
			recursive: true,
		},
		{
			name:     "recursive search no match",
			pattern:  "missing_fruit",
			files:    []string{dirPath},
			expected: Result{ExitCode: 1, Stdout: []byte{}},
		},
	}

	runFileTests(t, tests)
}

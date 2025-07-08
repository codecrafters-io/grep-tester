package grep

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchStdin_BasicMatching(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "literal character match",
			pattern: "d",
			input:   "dog",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"dog"},
			},
		},
		{
			name:    "literal character no match",
			pattern: "f",
			input:   "dog",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "digit match",
			pattern: "\\d",
			input:   "123",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"123"},
			},
		},
		{
			name:    "digit no match",
			pattern: "\\d",
			input:   "apple",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "alphanumeric match",
			pattern: "\\w",
			input:   "apple",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"apple"},
			},
		},
		{
			name:    "alphanumeric no match",
			pattern: "\\w",
			input:   "+-$×€÷",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_CharacterGroups(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "positive character group match",
			pattern: "[pineapple]",
			input:   "n",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"n"},
			},
		},
		{
			name:    "positive character group no match",
			pattern: "[cdfghijklm]",
			input:   "strawberry",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "negative character group match",
			pattern: "[^xyz]",
			input:   "apple",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"apple"},
			},
		},
		{
			name:    "negative character group no match",
			pattern: "[^anb]",
			input:   "banana",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_Anchors(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "start anchor match",
			pattern: "^log",
			input:   "log",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"log"},
			},
		},
		{
			name:    "start anchor no match",
			pattern: "^log",
			input:   "slog",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "end anchor match",
			pattern: "cat$",
			input:   "cat",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"cat"},
			},
		},
		{
			name:    "end anchor no match",
			pattern: "cat$",
			input:   "cats",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_Quantifiers(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "one or more match",
			pattern: "ca+t",
			input:   "cat",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"cat"},
			},
		},
		{
			name:    "one or more no match",
			pattern: "ca+t",
			input:   "act",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "zero or one match",
			pattern: "ca?t",
			input:   "cat",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"cat"},
			},
		},
		{
			name:    "zero or one alternate match",
			pattern: "ca?t",
			input:   "act",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"act"},
			},
		},
		{
			name:    "zero or one no match",
			pattern: "ca?t",
			input:   "dog",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_Wildcard(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "wildcard match",
			pattern: "c.t",
			input:   "cat",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"cat"},
			},
		},
		{
			name:    "wildcard no match",
			pattern: "c.t",
			input:   "car",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "wildcard with quantifier match",
			pattern: "g.+gol",
			input:   "goøö0Ogol",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"goøö0Ogol"},
			},
		},
		{
			name:    "wildcard with quantifier no match",
			pattern: "g.+gol",
			input:   "gol",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_Alternation(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "alternation match",
			pattern: "a (cat|dog)",
			input:   "a cat",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"a cat"},
			},
		},
		{
			name:    "alternation no match",
			pattern: "a (cat|dog)",
			input:   "a cow",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_SingleBackreference(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "basic backreference match",
			pattern: "(cat) and \\1",
			input:   "cat and cat",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"cat and cat"},
			},
		},
		{
			name:    "basic backreference no match",
			pattern: "(cat) and \\1",
			input:   "cat and dog",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "complex backreference match",
			pattern: "([abcd]+) is \\1, not [^xyz]+",
			input:   "abcd is abcd, not efg",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"abcd is abcd, not efg"},
			},
		},
		{
			name:    "complex backreference no match",
			pattern: "([abcd]+) is \\1, not [^xyz]+",
			input:   "efgh is efgh, not efg",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
		{
			name:    "anchored backreference match",
			pattern: "^(\\w+) starts and ends with \\1$",
			input:   "this starts and ends with this",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"this starts and ends with this"},
			},
		},
		{
			name:    "anchored backreference no match",
			pattern: "^(this) starts and ends with \\1$",
			input:   "that starts and ends with this",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_MultipleBackreferences(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "multiple backreferences match",
			pattern: "(\\d+) (\\w+) squares and \\1 \\2 circles",
			input:   "3 red squares and 3 red circles",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"3 red squares and 3 red circles"},
			},
		},
		{
			name:    "multiple backreferences no match",
			pattern: "(\\d+) (\\w+) squares and \\1 \\2 circles",
			input:   "3 red squares and 4 red circles",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchStdin_NestedBackreferences(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected Result
	}{
		{
			name:    "nested backreferences match",
			pattern: "('(cat) and \\2') is the same as \\1",
			input:   "'cat and cat' is the same as 'cat and cat'",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"'cat and cat' is the same as 'cat and cat'"},
			},
		},
		{
			name:    "nested backreferences no match",
			pattern: "('(cat) and \\2') is the same as \\1",
			input:   "'cat and cat' is the same as 'cat and dog'",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchStdin(tt.pattern, tt.input, Options{ExtendedRegex: true})

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

func TestSearchFiles_SingleFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// Write test content to file
	content := "lemon\n"
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	tests := []struct {
		name     string
		pattern  string
		expected Result
	}{
		{
			name:    "single file match",
			pattern: "le.*",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"lemon"},
			},
		},
		{
			name:    "single file no match",
			pattern: "broccoli",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchFiles(tt.pattern, []string{testFile}, Options{ExtendedRegex: true})

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

func TestSearchFiles_MultipleFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create fruits.txt
	fruitsFile := filepath.Join(tempDir, "fruits.txt")
	fruitsContent := "banana\nblueberry\n"
	err := os.WriteFile(fruitsFile, []byte(fruitsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write fruits file: %v", err)
	}

	// Create vegetables.txt
	vegetablesFile := filepath.Join(tempDir, "vegetables.txt")
	vegetablesContent := "broccoli\ncarrot\n"
	err = os.WriteFile(vegetablesFile, []byte(vegetablesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write vegetables file: %v", err)
	}

	tests := []struct {
		name     string
		pattern  string
		expected []string
	}{
		{
			name:     "multiple files match",
			pattern:  "b.*$",
			expected: []string{"fruits.txt:banana", "fruits.txt:blueberry", "vegetables.txt:broccoli"},
		},
		{
			name:     "multiple files no match",
			pattern:  "missing_fruit",
			expected: []string{},
		},
		{
			name:     "specific file match",
			pattern:  "carrot",
			expected: []string{"vegetables.txt:carrot"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchFiles(tt.pattern, []string{fruitsFile, vegetablesFile}, Options{ExtendedRegex: true})

			expectedExitCode := 1
			if len(tt.expected) > 0 {
				expectedExitCode = 0
			}

			if result.ExitCode != expectedExitCode {
				t.Errorf("Expected exit code %d, got %d", expectedExitCode, result.ExitCode)
			}

			if len(result.Stdout) != len(tt.expected) {
				t.Errorf("Expected %d stdout lines, got %d", len(tt.expected), len(result.Stdout))
			}

			// Convert absolute paths to relative paths for comparison
			actualLines := make([]string, len(result.Stdout))
			for i, line := range result.Stdout {
				// Extract the part after the last directory separator
				parts := strings.Split(line, "/")
				if len(parts) > 0 {
					// Get the filename:content part
					lastPart := parts[len(parts)-1]
					if len(parts) > 1 {
						// Check if the second-to-last part is a filename too
						secondLastPart := parts[len(parts)-2]
						if strings.Contains(lastPart, ":") {
							actualLines[i] = lastPart
						} else {
							actualLines[i] = secondLastPart + ":" + lastPart
						}
					} else {
						actualLines[i] = lastPart
					}
				}
			}

			// Check that all expected lines are present
			expectedMap := make(map[string]bool)
			for _, line := range tt.expected {
				expectedMap[line] = true
			}

			for _, line := range actualLines {
				if !expectedMap[line] {
					t.Errorf("Unexpected stdout line: %q", line)
				}
			}
		})
	}
}

func TestSearchFiles_Recursive(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create dir/fruits.txt
	dirPath := filepath.Join(tempDir, "dir")
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	fruitsFile := filepath.Join(dirPath, "fruits.txt")
	fruitsContent := "pear\nstrawberry\n"
	err = os.WriteFile(fruitsFile, []byte(fruitsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write fruits file: %v", err)
	}

	// Create dir/subdir/vegetables.txt
	subdirPath := filepath.Join(dirPath, "subdir")
	err = os.MkdirAll(subdirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	vegetablesFile := filepath.Join(subdirPath, "vegetables.txt")
	vegetablesContent := "celery\ncarrot\n"
	err = os.WriteFile(vegetablesFile, []byte(vegetablesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write vegetables file: %v", err)
	}

	// Create dir/vegetables.txt
	dirVegetablesFile := filepath.Join(dirPath, "vegetables.txt")
	dirVegetablesContent := "cucumber\ncorn\n"
	err = os.WriteFile(dirVegetablesFile, []byte(dirVegetablesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write dir vegetables file: %v", err)
	}

	tests := []struct {
		name     string
		pattern  string
		expected Result
	}{
		{
			name:    "recursive search match",
			pattern: ".*er",
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"dir/fruits.txt:strawberry", "dir/subdir/vegetables.txt:celery", "dir/vegetables.txt:cucumber"},
			},
		},
		{
			name:    "recursive search no match",
			pattern: "missing_fruit",
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchFiles(tt.pattern, []string{dirPath}, Options{ExtendedRegex: true, Recursive: true})

			if result.ExitCode != tt.expected.ExitCode {
				t.Errorf("Expected exit code %d, got %d", tt.expected.ExitCode, result.ExitCode)
			}

			if len(result.Stdout) != len(tt.expected.Stdout) {
				t.Errorf("Expected %d stdout lines, got %d", len(tt.expected.Stdout), len(result.Stdout))
			}

			// For recursive search, we need to check that all expected lines are present
			// but order might vary. Also normalize paths to be relative to tempDir
			expectedMap := make(map[string]bool)
			for _, line := range tt.expected.Stdout {
				expectedMap[line] = true
			}

			actualMap := make(map[string]bool)
			for _, line := range result.Stdout {
				// Convert absolute path to relative path starting from tempDir
				relativePath, err := filepath.Rel(tempDir, line)
				if err != nil {
					// If we can't get relative path, try to extract the meaningful part
					if strings.Contains(line, "/dir/") {
						parts := strings.Split(line, "/dir/")
						if len(parts) > 1 {
							relativePath = "dir/" + parts[1]
						} else {
							relativePath = line
						}
					} else {
						relativePath = line
					}
				}
				actualMap[relativePath] = true
			}

			for expectedLine := range expectedMap {
				if !actualMap[expectedLine] {
					t.Errorf("Expected stdout line %q not found in actual output", expectedLine)
				}
			}

			for actualLine := range actualMap {
				if !expectedMap[actualLine] {
					t.Errorf("Unexpected stdout line: %q", actualLine)
				}
			}
		})
	}
}

func TestSearchStdin_Options(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		options  Options
		expected Result
	}{
		{
			name:    "count option",
			pattern: "a",
			input:   "apple\nbanana\navocado",
			options: Options{ExtendedRegex: true, Count: true},
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"3"},
			},
		},
		{
			name:    "count option no matches",
			pattern: "xyz",
			input:   "apple\nbanana\navocado",
			options: Options{ExtendedRegex: true, Count: true},
			expected: Result{
				ExitCode: 1,
				Stdout:   []string{"0"},
			},
		},
		{
			name:    "quiet option with match",
			pattern: "apple",
			input:   "apple\nbanana",
			options: Options{ExtendedRegex: true, Quiet: true},
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{},
			},
		},
		{
			name:    "invert match",
			pattern: "apple",
			input:   "banana",
			options: Options{ExtendedRegex: true, InvertMatch: true},
			expected: Result{
				ExitCode: 0,
				Stdout:   []string{"banana"},
			},
		},
	}

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

func TestBackrefMatcher(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		text     string
		expected bool
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := &BackrefMatcher{pattern: tt.pattern}
			result := matcher.Match(tt.text)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRegexMatcher(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		text     string
		expected bool
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := createMatcher(tt.pattern)
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

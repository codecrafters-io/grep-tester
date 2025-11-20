package utils

import (
	"bufio"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/codecrafters-io/tester-utils/random"
)

// FRUITS, VEGETABLES, and ANIMALS are used to generate test file contents
var FRUITS = []string{"apple", "banana", "blackberry", "blueberry", "cherry", "grape", "lemon", "mango", "orange", "pear", "pineapple", "plum", "raspberry", "strawberry", "watermelon"}
var VEGETABLES = []string{"carrot", "onion", "potato", "tomato", "broccoli", "cauliflower", "cabbage", "lettuce", "spinach", "asparagus", "pea", "corn", "zucchini", "pumpkin"}
var ANIMALS = []string{"cat", "dog", "elephant", "fox", "giraffe", "horse", "lion", "monkey", "panda", "rabbit", "tiger", "wolf", "zebra"}

// RandomFilePrefix returns 4 digit random prefix for test files
func RandomFilePrefix() string {
	return strconv.Itoa(random.RandomInt(1000, 10000))
}

func RandomWordsWithoutSubstrings(n int) []string {
loop:
	for {
		words := random.RandomWords(n)

		for i := 0; i < len(words); i++ {
			for j := i + 1; j < len(words); j++ {
				if strings.Contains(words[j], words[i]) || strings.Contains(words[i], words[j]) {
					continue loop
				}
			}
		}

		return words
	}
}

// getRegex101Link returns the link to regex101 website with given pattern and test string
// Notes for the future: grep -E offers POSIX Extended Regular Expression compatibility (PCRE2 is also supported, but is highly experimental
// Ref. https://superuser.com/questions/269803/which-regular-expression-standard-is-used-in-grep), while the website regex101
// doesn't offer this standard. It uses PCRE2, and other standards.
// The comparison between different standards are given here: https://gist.github.com/CMCDragonkai/6c933f4a7d713ef712145c5eb94a1816
// While POSIX ERE seems a 'subset' of PCRE2 (Every feature supported by POSIX ERE is also supported by PCRE2), if, in any case,
// any discrepancies arise in the future, this comment serves as a note for debugging.
func GetRegex101Link(pattern string, testString string) string {
	return fmt.Sprintf("https://regex101.com/?regex=%s&testString=%s", url.QueryEscape(pattern), url.QueryEscape(testString))
}

// ProgramOutputToLines converts a program's output to a string slice, in which
// each element is an individual lines
// The resulting string slice is exactly what one would expect to find in the terminal
func ProgramOutputToLines(output string) []string {
	sc := bufio.NewScanner(strings.NewReader(output))
	sc.Split(bufio.ScanLines)
	var out []string

	for sc.Scan() {
		out = append(out, sc.Text())
	}

	if out == nil {
		return []string{}
	}

	return out
}

// FormatLineForLogging formats a line such that its suitable for logging
// it escapes backslash sequences and adds int for empty line
func FormatLineForLogging(line string) string {
	if line == "" {
		return fmt.Sprintf("%q (empty line)", line)
	}

	return fmt.Sprintf("%q", line)
}

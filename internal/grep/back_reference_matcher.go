package grep

import (
	"fmt"
	"regexp"
	"strings"
)

// BackrefMatcher handles patterns with backreferences (\1, \2, etc.).
// Go regexp uses the RE2 engine, which doesn't support backreferences out of the box.
// This matcher implements backreferences using pattern expansion and validation.
type backrefMatcher struct {
	pattern string
}

func newBackReferenceMatcher(pattern string) *backrefMatcher {
	return &backrefMatcher{
		pattern: pattern,
	}
}

type matchResult struct {
	success        bool
	matchedStrings []string
}

func (m *backrefMatcher) match(text string) matchResult {
	pattern := m.pattern

	// Handle patterns without backreferences first
	if !regexp.MustCompile(`\\[1-9]`).MatchString(pattern) {
		regex := regexp.MustCompile(pattern)
		matches := regex.FindAllString(text, -1)

		return matchResult{
			success:        len(matches) > 0,
			matchedStrings: matches,
		}

	}

	// First, create a version of the pattern where backreferences are replaced with (.*)
	// This allows us to use Go's regex engine to find potential matches
	regexPattern := regexp.MustCompile(`\\[1-9]`).ReplaceAllString(pattern, `(.*)`)
	regex := regexp.MustCompile(regexPattern)

	candidateMatches := regex.FindAllStringSubmatchIndex(text, -1)
	var finalMatches []string

	for _, submatchIdx := range candidateMatches {
		fullStart, fullEnd := submatchIdx[0], submatchIdx[1]
		fullMatch := text[fullStart:fullEnd]

		// Extract captured groups
		groups := []string{}
		for i := 2; i < len(submatchIdx); i += 2 {
			start, end := submatchIdx[i], submatchIdx[i+1]
			if start >= 0 && end >= 0 {
				groups = append(groups, text[start:end])
			} else {
				groups = append(groups, "")
			}
		}

		// Validate backreferences
		if validateBackreferences(pattern, append([]string{fullMatch}, groups...)) {
			finalMatches = append(finalMatches, fullMatch)
		}
	}

	return matchResult{
		success:        len(finalMatches) > 0,
		matchedStrings: finalMatches,
	}
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

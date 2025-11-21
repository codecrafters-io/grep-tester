package grep

import (
	"fmt"
	"regexp"
	"strings"
)

// BackrefMatcher handles patterns with backreferences (\1, \2, etc.).
// Go regexp uses the RE2 engine, which doesn't support backreferences out of the box.
// This matcher implements backreferences using pattern expansion and validation.
type backReferenceMatcher struct {
	pattern string
}

func newBackReferenceMatcher(pattern string) *backReferenceMatcher {
	return &backReferenceMatcher{
		pattern: pattern,
	}
}

type matchResult struct {
	success        bool
	matchedStrings []string
}

func (m *backReferenceMatcher) match(text string) matchResult {
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

	// First, create a version of the pattern where backreferences are replaced with (?:.*)
	// This allows us to use Go's regex engine to find potential matches
	// We use a non-capture group instead of capture group because we don't want the fake capture groups
	// to appear inside the result of FindAllStringSubmatch()
	regexPattern := regexp.MustCompile(`\\[1-9]`).ReplaceAllString(pattern, `(?:.*)`)
	regex := regexp.MustCompile(regexPattern)

	possibleMatches := regex.FindAllStringSubmatch(text, -1)
	qualifiedMatches := []string{}

	for _, possibleMatch := range possibleMatches {
		fullSubMatch := possibleMatch[0]

		if validateBackreferencesForSubmatch(pattern, possibleMatch) {
			qualifiedMatches = append(qualifiedMatches, fullSubMatch)
		}
	}

	return matchResult{
		success:        len(qualifiedMatches) > 0,
		matchedStrings: qualifiedMatches,
	}
}

func validateBackreferencesForSubmatch(pattern string, subMatch []string) bool {
	if len(subMatch) < 2 {
		return false
	}

	// match[0] is the full match
	// match[1], match[2], etc. are the captured groups
	fullMatch := subMatch[0]

	// Replace backreferences with the actual captured values
	expandedPattern := pattern
	for i := 1; i < len(subMatch); i++ {
		backref := fmt.Sprintf("\\%d", i)
		expandedPattern = strings.ReplaceAll(expandedPattern, backref, regexp.QuoteMeta(subMatch[i]))
	}

	// Now check if this expanded pattern matches the full match
	regex := regexp.MustCompile("^" + expandedPattern + "$")

	return regex.MatchString(fullMatch)
}

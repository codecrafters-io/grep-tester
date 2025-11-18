package test_cases

import (
	"fmt"
	"net/url"
)

// getRegex101Link returns the link to regex101 website with given pattern and test string
// Notes for the future: grep -E offers POSIX Extended Regular Expression compatibility (PCRE2 is also supported, but is highly experimental
// Ref. https://superuser.com/questions/269803/which-regular-expression-standard-is-used-in-grep), while the website regex101
// doesn't offer this standard. It uses PCRE2, and other standards.
// The comparison between different standards are given here: https://gist.github.com/CMCDragonkai/6c933f4a7d713ef712145c5eb94a1816
// While POSIX ERE seems a 'subset' of PCRE2 (Every feature supported by POSIX ERE is also supported by PCRE2), if, in any case,
// any discrepancies arise in the future, this comment serves as a note for debugging.
func getRegex101Link(pattern string, testString string) string {
	return fmt.Sprintf("https://regex101.com/?regex=%s&testString=%s", url.QueryEscape(pattern), url.QueryEscape(testString))
}

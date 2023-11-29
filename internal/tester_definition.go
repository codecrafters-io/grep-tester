package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
	"time"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatTestCases: []testerutils.TestCase{},
	ExecutableFileName: "your_grep.sh",
	TestCases: []testerutils.TestCase{
		{
			Slug:     "init",
			TestFunc: testInit,
		},
		{
			Slug:     "match_digit",
			TestFunc: testMatchDigit,
		},
		{
			Slug:     "match_alphanumeric",
			TestFunc: testMatchAlphanumeric,
		},
		{
			Slug:     "positive_character_groups",
			TestFunc: testPositiveCharacterGroups,
		},
		{
			Slug:     "negative_character_groups",
			TestFunc: testNegativeCharacterGroups,
		},
		{
			Slug:     "combining_character_classes",
			TestFunc: testCombiningCharacterClasses,
			Timeout:  20 * time.Second,
		},
		{
			Slug:     "start_of_string_anchor",
			TestFunc: testStartOfStringAnchor,
		},
		{
			Slug:     "end_of_string_anchor",
			TestFunc: testEndOfStringAnchor,
		},
		{
			Slug:     "one_or_more_quantifier",
			TestFunc: testOneOrMoreQuantifier,
		},
		{
			Slug:     "zero_or_one_quantifier",
			TestFunc: testZeroOrOneQuantifier,
		},
		{
			Slug:     "wildcard",
			TestFunc: testWildcard,
		},
		{
			Slug:     "alternation",
			TestFunc: testAlternation,
		},
		{
			Slug:     "backreferences-single",
			TestFunc: testBackreferencesSingle,
		},
		{
			Slug:     "backreferences-multiple",
			TestFunc: testBackreferencesMultiple,
		},
		{
			Slug:     "backreferences-nested",
			TestFunc: testBackreferencesNested,
		},
	},
}

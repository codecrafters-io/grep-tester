package internal

import (
	"github.com/codecrafters-io/tester-utils/tester_definition"
	"time"
)

var testerDefinition = tester_definition.TesterDefinition{
	AntiCheatTestCases: []tester_definition.TestCase{},
	ExecutableFileName: "your_grep.sh",
	TestCases: []tester_definition.TestCase{
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
			Timeout:  20 * time.Second,
		},
		{
			Slug:     "backreferences-multiple",
			TestFunc: testBackreferencesMultiple,
			Timeout:  20 * time.Second,
		},
		{
			Slug:     "backreferences-nested",
			TestFunc: testBackreferencesNested,
			Timeout:  20 * time.Second,
		},
	},
}

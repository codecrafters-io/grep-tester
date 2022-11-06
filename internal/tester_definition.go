package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
	"time"
)

var testerDefinition = testerutils.TesterDefinition{
	AntiCheatStages:    []testerutils.Stage{},
	ExecutableFileName: "your_grep.sh",
	Stages: []testerutils.Stage{
		{
			Number:                  1,
			Slug:                    "init",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  2,
			Slug:                    "match_digit",
			Title:                   "Match digits",
			TestFunc:                testMatchDigit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  3,
			Slug:                    "match_alphanumeric",
			Title:                   "Match alphanumeric characters",
			TestFunc:                testMatchAlphanumeric,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  4,
			Slug:                    "positive_character_groups",
			Title:                   "Positive Character Groups",
			TestFunc:                testPositiveCharacterGroups,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  5,
			Slug:                    "negative_character_groups",
			Title:                   "Negative Character Groups",
			TestFunc:                testNegativeCharacterGroups,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  6,
			Slug:                    "combining_character_classes",
			Title:                   "Combining Character Classes",
			TestFunc:                testCombiningCharacterClasses,
			ShouldRunPreviousStages: true,
			Timeout:                 20 * time.Second,
		},
		{
			Number:                  7,
			Slug:                    "start_of_string_anchor",
			Title:                   "Start of string anchor",
			TestFunc:                testStartOfStringAnchor,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  8,
			Slug:                    "end_of_string_anchor",
			Title:                   "End of string anchor",
			TestFunc:                testEndOfStringAnchor,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  9,
			Slug:                    "one_or_more_quantifier",
			Title:                   "Match one or more times",
			TestFunc:                testOneOrMoreQuantifier,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  10,
			Slug:                    "zero_or_one_quantifier",
			Title:                   "Match zero or one times",
			TestFunc:                testZeroOrOneQuantifier,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  11,
			Slug:                    "wildcard",
			Title:                   "Wildcard",
			TestFunc:                testWildcard,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  12,
			Slug:                    "alternation",
			Title:                   "Alternation",
			TestFunc:                testAlternation,
			ShouldRunPreviousStages: true,
		},
	},
}

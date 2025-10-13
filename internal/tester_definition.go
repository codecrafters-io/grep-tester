package internal

import (
	"time"

	"github.com/codecrafters-io/tester-utils/tester_definition"
)

var testerDefinition = tester_definition.TesterDefinition{
	AntiCheatTestCases: []tester_definition.TestCase{
		{
			Slug:     "anti-cheat-1",
			TestFunc: testWordBoundariesAsAntiCheat,
		},
	},
	ExecutableFileName:       "your_program.sh",
	LegacyExecutableFileName: "your_grep.sh",
	TestCases: []tester_definition.TestCase{
		{
			Slug:     "cq2",
			TestFunc: testInit,
		},
		{
			Slug:     "oq2",
			TestFunc: testMatchDigit,
		},
		{
			Slug:     "mr9",
			TestFunc: testMatchAlphanumeric,
		},
		{
			Slug:     "tl6",
			TestFunc: testPositiveCharacterGroups,
		},
		{
			Slug:     "rk3",
			TestFunc: testNegativeCharacterGroups,
		},
		{
			Slug:     "sh9",
			TestFunc: testCombiningCharacterClasses,
			Timeout:  20 * time.Second,
		},
		{
			Slug:     "rr8",
			TestFunc: testStartOfStringAnchor,
		},
		{
			Slug:     "ao7",
			TestFunc: testEndOfStringAnchor,
		},
		{
			Slug:     "fz7",
			TestFunc: testOneOrMoreQuantifier,
		},
		{
			Slug:     "ny8",
			TestFunc: testZeroOrOneQuantifier,
		},
		{
			Slug:     "zb3",
			TestFunc: testWildcard,
		},
		{
			Slug:     "zm7",
			TestFunc: testAlternation,
		},
		{
			Slug:     "sb5",
			TestFunc: testBackreferencesSingle,
			Timeout:  20 * time.Second,
		},
		{
			Slug:     "tg1",
			TestFunc: testBackreferencesMultiple,
			Timeout:  20 * time.Second,
		},
		{
			Slug:     "xe5",
			TestFunc: testBackreferencesNested,
			Timeout:  20 * time.Second,
		},
		{
			Slug:     "dr5",
			TestFunc: testSingleLineFileSearch,
		},
		{
			Slug:     "ol9",
			TestFunc: testMultiLineFileSearch,
		},
		{
			Slug:     "is6",
			TestFunc: testMultiFileSearch,
		},
		{
			Slug:     "yx6",
			TestFunc: testRecursiveFileSearch,
		},
		{
			Slug:     "ai9",
			TestFunc: testQuantifierAsterisk,
		},
		{
			Slug:     "wy9",
			TestFunc: testQuantifierExactRepitition,
		},
		{
			Slug:     "hk3",
			TestFunc: testQuantifierMinimumRepitition,
		},
		{
			Slug:     "ug0",
			TestFunc: testQuantifierRangeRepitition,
		},
	},
}

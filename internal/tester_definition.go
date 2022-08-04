package internal

import (
	testerutils "github.com/codecrafters-io/tester-utils"
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
			Title:                   "Match a digit",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  3,
			Slug:                    "match_alphanumeric",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  4,
			Slug:                    "positive_character_groups",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  5,
			Slug:                    "negative_character_groups",
			Title:                   "Negative a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  6,
			Slug:                    "start_of_string_anchor",
			Title:                   "Start of string",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  1,
			Slug:                    "end_of_string_anchor",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  1,
			Slug:                    "one_or_more_quantifier",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  1,
			Slug:                    "zero_or_one_quantifier",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  1,
			Slug:                    "wildcard",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
		{
			Number:                  1,
			Slug:                    "alternation",
			Title:                   "Match a literal character",
			TestFunc:                testInit,
			ShouldRunPreviousStages: true,
		},
	},
}

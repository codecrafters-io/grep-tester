package internal

import (
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"init_pass": {
			UntilStageSlug:      "init",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_fail": {
			UntilStageSlug:      "init",
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"match_digit_pass": {
			UntilStageSlug:      "match_digit",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/match_digit/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"match_alphanumeric_pass": {
			UntilStageSlug:      "match_alphanumeric",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/match_alphanumeric/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"positive_character_groups_pass": {
			UntilStageSlug:      "positive_character_groups",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/positive_character_groups/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"negative_character_groups_pass": {
			UntilStageSlug:      "negative_character_groups",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/negative_character_groups/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"combining_character_classes_pass": {
			UntilStageSlug:      "combining_character_classes",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/combining_character_classes/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"start_of_string_anchor_pass": {
			UntilStageSlug:      "start_of_string_anchor",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/start_of_string_anchor/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"end_of_string_anchor_pass": {
			UntilStageSlug:      "end_of_string_anchor",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/end_of_string_anchor/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"one_or_more_quantifier_pass": {
			UntilStageSlug:      "one_or_more_quantifier",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/one_or_more_quantifier/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"zero_or_one_quantifier_pass": {
			UntilStageSlug:      "zero_or_one_quantifier",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/zero_or_one_quantifier/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"wildcard": {
			UntilStageSlug:      "wildcard",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/wildcard/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"alternation": {
			UntilStageSlug:      "alternation",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/alternation/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	return testerOutput
}

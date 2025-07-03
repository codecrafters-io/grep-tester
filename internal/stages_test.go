package internal

import (
	"os"
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	os.Setenv("CODECRAFTERS_RANDOM_SEED", "1234567890")

	falseVar := false

	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"cheat_innocent": {
			UntilStageSlug:      "cq2",
			CodePath:            "./test_helpers/scenarios/cheat_innocent",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/cheat/cheat_innocent",
			NormalizeOutputFunc: normalizeTesterOutput,
			SkipAntiCheat:       &falseVar,
		},
		"cheat_suspect": {
			UntilStageSlug:      "cq2",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/cheat/cheat_suspect",
			NormalizeOutputFunc: normalizeTesterOutput,
			SkipAntiCheat:       &falseVar,
		},
		"init_pass": {
			UntilStageSlug:      "cq2",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/init/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"init_fail": {
			UntilStageSlug:      "cq2",
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"match_digit_pass": {
			UntilStageSlug:      "oq2",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/match_digit/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"match_alphanumeric_pass": {
			UntilStageSlug:      "mr9",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/match_alphanumeric/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"positive_character_groups_pass": {
			UntilStageSlug:      "tl6",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/positive_character_groups/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"negative_character_groups_pass": {
			UntilStageSlug:      "rk3",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/negative_character_groups/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"combining_character_classes_pass": {
			UntilStageSlug:      "sh9",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/combining_character_classes/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"start_of_string_anchor_pass": {
			UntilStageSlug:      "rr8",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/start_of_string_anchor/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"end_of_string_anchor_pass": {
			UntilStageSlug:      "ao7",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/end_of_string_anchor/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"one_or_more_quantifier_pass": {
			UntilStageSlug:      "fz7",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/one_or_more_quantifier/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"zero_or_one_quantifier_pass": {
			UntilStageSlug:      "ny8",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/zero_or_one_quantifier/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"wildcard_pass": {
			UntilStageSlug:      "zb3",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/wildcard/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"alternation_pass": {
			UntilStageSlug:      "zm7",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/alternation/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"backreferences_single_pass": {
			UntilStageSlug:      "sb5",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/backreferences_single/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"backreferences_multiple_pass": {
			UntilStageSlug:      "tg1",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/backreferences_multiple/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"backreferences_nested_pass": {
			UntilStageSlug:      "xe5",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/backreferences_nested/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	return testerOutput
}

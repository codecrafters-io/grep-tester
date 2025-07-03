package internal

import (
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	falseVar := false

	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"cheat_innocent": {
			StageSlugs:          []string{"cq2"},
			CodePath:            "./test_helpers/scenarios/cheat_innocent",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/cheat/cheat_innocent",
			NormalizeOutputFunc: normalizeTesterOutput,
			SkipAntiCheat:       &falseVar,
		},
		"cheat_suspect": {
			StageSlugs:          []string{"cq2"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/cheat/cheat_suspect",
			NormalizeOutputFunc: normalizeTesterOutput,
			SkipAntiCheat:       &falseVar,
		},
		"init_fail": {
			StageSlugs:          []string{"cq2"},
			CodePath:            "./test_helpers/scenarios/init/failure",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/init/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"base_stages_pass": {
			UntilStageSlug:      "zm7",
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/base_stages/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"backreferences_pass": {
			StageSlugs:          []string{"sb5", "tg1", "xe5"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/backreferences/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"file_search_pass": {
			StageSlugs:          []string{"dr5", "ol9", "is6", "yx6"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/file_search/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	return testerOutput
}

package internal

import (
	"os"
	"regexp"
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	os.Setenv("CODECRAFTERS_RANDOM_SEED", "1234567890")
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
		"extra_logs": {
			StageSlugs:          []string{"yx6"},
			CodePath:            "./test_helpers/scenarios/extra_logs",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/extra_logs/failure",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"missing_and_extra_logs": {
			StageSlugs:          []string{"yx6"},
			CodePath:            "./test_helpers/scenarios/missing_and_extra_logs",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/missing_and_extra_logs/failure",
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
		"quantifiers_pass": {
			StageSlugs:          []string{"ai9", "wy9", "hk3", "ug0"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/quantifiers/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"printing_matches_pass": {
			StageSlugs:          []string{"ku5", "pz6"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/printing_matches/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"printing_matches_unexpected_output": {
			StageSlugs:          []string{"ku5"},
			CodePath:            "./test_helpers/scenarios/printing_matches/unexpected_output",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/printing_matches/unexpected_output",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"printing_matches_extra_lines": {
			StageSlugs:          []string{"ku5"},
			CodePath:            "./test_helpers/scenarios/printing_matches/extra_lines",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/printing_matches/extra_lines",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"printing_matches_swapped_lines": {
			StageSlugs:          []string{"pz6"},
			CodePath:            "./test_helpers/scenarios/printing_matches/swapped_lines",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/printing_matches/swapped_lines",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"printing_matches_highlighted": {
			StageSlugs:          []string{"ku5"},
			CodePath:            "./test_helpers/scenarios/printing_matches/highlighted",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/printing_matches/highlighted",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	replacements := map[string][]*regexp.Regexp{
		// We can't be sure about the order of the output lines here
		"grep_output_with_dir_prefix":   {regexp.MustCompile(`.{5}\[your_program\].{5}dir/.*`)},
		"grep_output_with_dir_prefix_2": {regexp.MustCompile(`.*\[tester::#YX6\].*✓ Found line 'dir/.*`)},
		"grep_output_with_dir_prefix_3": {regexp.MustCompile(`.*\[tester::#YX6\].*⨯ Line not found: "dir/.*`)},
	}

	for replacement, regexes := range replacements {
		for _, regex := range regexes {
			testerOutput = regex.ReplaceAll(testerOutput, []byte(replacement))
		}
	}

	return testerOutput
}

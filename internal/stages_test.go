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
		"file_search_unexpected_output": {
			StageSlugs:          []string{"dr5"},
			CodePath:            "./test_helpers/scenarios/file_search/unexpected_output",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/file_search/unexpected_output",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"file_search_extra_empty_line": {
			StageSlugs:          []string{"is6"},
			CodePath:            "./test_helpers/scenarios/file_search/extra_empty_line",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/file_search/extra_empty_line",
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
		"printing_matches_extra_empty_line": {
			StageSlugs:          []string{"ku5"},
			CodePath:            "./test_helpers/scenarios/printing_matches/extra_empty_line",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/printing_matches/extra_empty_line",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"multiple_matches_pass": {
			StageSlugs:          []string{"cj0", "ss2", "bo4"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/multiple_matches/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_pass": {
			StageSlugs:          []string{"bm2", "eq0", "wg2", "jk4", "na5"},
			CodePath:            "./test_helpers/pass_all",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/success",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_missing_line": {
			StageSlugs:          []string{"bm2"},
			CodePath:            "./test_helpers/scenarios/highlighting/missing_line",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/missing_line",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_no_highlight": {
			StageSlugs:          []string{"bm2"},
			CodePath:            "./test_helpers/scenarios/highlighting/no_highlight",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/no_highlight",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_multiple_highlighting_fail": {
			StageSlugs:          []string{"eq0"},
			CodePath:            "./test_helpers/scenarios/highlighting/multiple_highlighting_fail",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/multiple_highlighting_fail",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_red_but_not_bold": {
			StageSlugs:          []string{"bm2"},
			CodePath:            "./test_helpers/scenarios/highlighting/red_but_not_bold",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/red_but_not_bold",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_bold_but_not_red": {
			StageSlugs:          []string{"bm2"},
			CodePath:            "./test_helpers/scenarios/highlighting/bold_but_not_red",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/bold_but_not_red",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_bold_red_and_italic": {
			StageSlugs:          []string{"bm2"},
			CodePath:            "./test_helpers/scenarios/highlighting/bold_red_and_italic",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/bold_red_and_italic",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_wrong_color_green": {
			StageSlugs:          []string{"bm2"},
			CodePath:            "./test_helpers/scenarios/highlighting/wrong_color_green",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/wrong_color_green",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_unexpected_highlight": {
			StageSlugs:          []string{"jk4"},
			CodePath:            "./test_helpers/scenarios/highlighting/unexpected_highlight",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/unexpected_highlight",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_second_line_missing_highlighting": {
			StageSlugs:          []string{"wg2"},
			CodePath:            "./test_helpers/scenarios/highlighting/second_line_missing_highlighting",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/second_line_missing_highlighting",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_second_line_unexpected_highlighting": {
			StageSlugs:          []string{"jk4"},
			CodePath:            "./test_helpers/scenarios/highlighting/second_line_unexpected_highlighting",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/second_line_unexpected_highlighting",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"highlighting_closing_sequence_missing": {
			StageSlugs:          []string{"na5"},
			CodePath:            "./test_helpers/scenarios/highlighting/closing_sequence_missing",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/highlighting/closing_sequence_missing",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	replacements := map[string][]*regexp.Regexp{
		// We can't be sure about the order of the output lines here
		"grep_output_with_dir_prefix":   {regexp.MustCompile(`.{5}\[your_program\].{5}dir/.*`)},
		"grep_output_with_dir_prefix_2": {regexp.MustCompile(`.*\[tester::#YX6\].*✓ Found line "dir/.*`)},
		"grep_output_with_dir_prefix_3": {regexp.MustCompile(`.*\[tester::#YX6\].*⨯ Line not found: "dir/.*`)},
	}

	for replacement, regexes := range replacements {
		for _, regex := range regexes {
			testerOutput = regex.ReplaceAll(testerOutput, []byte(replacement))
		}
	}

	return testerOutput
}

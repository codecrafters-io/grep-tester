package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testSingleLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []FileSearchTestCase{
		{
			Pattern:          "appl.*",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"apple"},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "apple"},
			},
		},
		{
			Pattern:          "carrot",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "apple"},
			},
		},
		{
			Pattern:          ".*ple",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"apple"},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "apple"},
			},
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}

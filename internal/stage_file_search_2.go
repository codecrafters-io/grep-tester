package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []FileSearchTestCase{
		{
			Pattern:          ".*berry",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"blackberry", "blueberry"},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "banana\ngrape\nblackberry\nblueberry"},
			},
		},
		{
			Pattern:          "carrot",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "banana\ngrape\nblackberry\nblueberry"},
			},
		},
		{
			Pattern:          "grape",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"grape"},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "banana\ngrape\nblackberry\nblueberry"},
			},
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}
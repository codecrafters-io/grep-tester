package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []FileSearchTestCase{
		{
			Pattern:          "b.*$",
			FilePaths:        []string{"fruits.txt", "vegetables.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"fruits.txt:banana", "fruits.txt:blueberry", "vegetables.txt:broccoli"},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "banana\nblueberry"},
				{Path: "vegetables.txt", Content: "broccoli\ncarrot"},
			},
		},
		{
			Pattern:          "missing_fruit",
			FilePaths:        []string{"fruits.txt", "vegetables.txt"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "banana\nblueberry"},
				{Path: "vegetables.txt", Content: "broccoli\ncarrot"},
			},
		},
		{
			Pattern:          "carrot",
			FilePaths:        []string{"fruits.txt", "vegetables.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"vegetables.txt:carrot"},
			TestFiles: []TestFile{
				{Path: "fruits.txt", Content: "banana\nblueberry"},
				{Path: "vegetables.txt", Content: "broccoli\ncarrot"},
			},
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}

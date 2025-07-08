package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRecursiveFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testFiles := []TestFile{
		{Path: "dir/fruits.txt", Content: "pear\nstrawberry"},
		{Path: "dir/subdir/vegetables.txt", Content: "celery\ncarrot"},
		{Path: "dir/vegetables.txt", Content: "cucumber\ncorn"},
	}
	if err := CreateTestFiles(testFiles, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:                   ".*er",
			FilePaths:                 []string{"dir/"},
			ExpectedExitCode:          0,
			ExpectedOutput:            []string{"dir/subdir/vegetables.txt:celery", "dir/vegetables.txt:cucumber", "dir/fruits.txt:strawberry"},
			ShouldEnableRecursiveFlag: true,
		},
		{
			Pattern:                   ".*ar",
			FilePaths:                 []string{"dir/"},
			ExpectedExitCode:          0,
			ExpectedOutput:            []string{"dir/subdir/vegetables.txt:carrot", "dir/fruits.txt:pear"},
			ShouldEnableRecursiveFlag: true,
		},
		{
			Pattern:                   "missing_fruit",
			FilePaths:                 []string{"dir/"},
			ExpectedExitCode:          1,
			ExpectedOutput:            []string{},
			ShouldEnableRecursiveFlag: true,
		},
	}

	return testCases.Run(stageHarness)
}

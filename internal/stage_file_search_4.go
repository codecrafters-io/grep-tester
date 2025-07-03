package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRecursiveFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []FileSearchTestCase{
		{
			Pattern:          ".*er",
			FilePaths:        []string{"dir/"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"dir/fruits.txt:strawberry", "dir/subdir/vegetables.txt:celery", "dir/vegetables.txt:cucumber"},
			Recursive:        true,
			TestFiles: []TestFile{
				{Path: "dir/fruits.txt", Content: "pear\nstrawberry"},
				{Path: "dir/subdir/vegetables.txt", Content: "celery\ncarrot"},
				{Path: "dir/vegetables.txt", Content: "cucumber\ncorn"},
			},
		},
		{
			Pattern:          ".*ar",
			FilePaths:        []string{"dir/"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"dir/fruits.txt:pear", "dir/subdir/vegetables.txt:carrot"},
			Recursive:        true,
			TestFiles: []TestFile{
				{Path: "dir/fruits.txt", Content: "pear\nstrawberry"},
				{Path: "dir/subdir/vegetables.txt", Content: "celery\ncarrot"},
				{Path: "dir/vegetables.txt", Content: "cucumber\ncorn"},
			},
		},
		{
			Pattern:          "missing_fruit",
			FilePaths:        []string{"dir/"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
			Recursive:        true,
			TestFiles: []TestFile{
				{Path: "dir/fruits.txt", Content: "pear\nstrawberry"},
				{Path: "dir/subdir/vegetables.txt", Content: "celery\ncarrot"},
				{Path: "dir/vegetables.txt", Content: "cucumber\ncorn"},
			},
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}
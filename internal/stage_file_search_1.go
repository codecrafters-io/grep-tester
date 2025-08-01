package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testSingleLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	file_name := "fruits-" + randomFilePrefix() + ".txt"
	fruit_1 := random.RandomElementFromArray(FRUITS)
	vegetable_1 := random.RandomElementFromArray(VEGETABLES)
	testFiles := []TestFile{
		{Path: file_name, Content: fruit_1},
	}
	if err := CreateTestFiles(testFiles, stageHarness); err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: Failed to create test files: %v", err))
	}

	testCaseCollection := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:          fruit_1[:len(fruit_1)/2] + ".+",
			FilePaths:        []string{file_name},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          vegetable_1,
			FilePaths:        []string{file_name},
			ExpectedExitCode: 1,
		},
		{
			Pattern:          ".+" + fruit_1[len(fruit_1)/2:],
			FilePaths:        []string{file_name},
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}

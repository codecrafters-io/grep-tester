package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	file_name_1 := "fruits-" + randomFilePrefix() + ".txt"
	file_name_2 := "vegetables-" + randomFilePrefix() + ".txt"
	fruits := random.RandomElementsFromArray(FRUITS, random.RandomInt(2, 3))
	vegetables := random.RandomElementsFromArray(VEGETABLES, random.RandomInt(2, 3))
	testFiles := []TestFile{
		{Path: file_name_1, Content: strings.Join(fruits, "\n")},
		{Path: file_name_2, Content: strings.Join(vegetables, "\n")},
	}
	if err := CreateTestFiles(testFiles, stageHarness); err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: Failed to create test files: %v", err))
	}

	testCaseCollection := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:          fruits[0][:2] + ".+$",
			FilePaths:        []string{file_name_1, file_name_2},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "missing_fruit",
			FilePaths:        []string{file_name_1, file_name_2},
			ExpectedExitCode: 1,
		},
		{
			Pattern:          vegetables[0],
			FilePaths:        []string{file_name_1, file_name_2},
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}

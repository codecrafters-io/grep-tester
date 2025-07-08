package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	file_name := "fruits-" + strconv.Itoa(random.RandomInt(1000, 10000)) + ".txt"
	fruits := append(random.RandomElementsFromArray(FRUITS, 3), "blueberry")
	vegetable_1 := random.RandomElementFromArray(VEGETABLES)
	testFiles := []TestFile{
		{Path: file_name, Content: strings.Join(fruits, "\n")},
	}
	if err := CreateTestFiles(testFiles, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:   ".*berry",
			FilePaths: []string{file_name},
		},
		{
			Pattern:   vegetable_1,
			FilePaths: []string{file_name},
		},
		{
			Pattern:   fruits[0],
			FilePaths: []string{file_name},
		},
	}

	return testCases.Run(stageHarness)
}

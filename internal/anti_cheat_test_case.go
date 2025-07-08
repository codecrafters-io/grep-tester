package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type AntiCheatTestCase struct {
	Pattern string
	Input   string
}

type AntiCheatTestCaseCollection []AntiCheatTestCase

func (testCases AntiCheatTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger

	for _, testCase := range testCases {
		executable := stageHarness.Executable.Clone()
		executable.TimeoutInMilliseconds = 1000

		// Get expected results from internal grep implementation
		expectedResult := grep.SearchStdin(testCase.Pattern, testCase.Input, grep.Options{
			ExtendedRegex: true,
		})

		// Run the actual executable
		actualResult, err := executable.RunWithStdin([]byte(testCase.Input), "-E", testCase.Pattern)
		if err != nil && err.Error() == "execution timed out" {
			continue
		}
		
		// If the executable matches our internal implementation exactly, it's likely cheating
		if actualResult.ExitCode == expectedResult.ExitCode {
			logger.Criticalf("anti-cheat (ac1) failed.")
			logger.Criticalf("Please contact us at hello@codecrafters.io if you think this is a mistake.")
			return fmt.Errorf("anti-cheat (ac1) failed")
		}
	}
	
	return nil
}

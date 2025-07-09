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

func (c AntiCheatTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger

	for _, testCase := range c {
		executable := stageHarness.Executable.Clone()
		executable.TimeoutInMilliseconds = 1000

		expectedResult := grep.SearchStdin(testCase.Pattern, testCase.Input, grep.Options{
			ExtendedRegex: true,
		})
		actualResult, err := executable.RunWithStdin([]byte(testCase.Input), "-E", testCase.Pattern)
		if err != nil && err.Error() == "execution timed out" {
			continue
		}
		if actualResult.ExitCode == expectedResult.ExitCode {
			logger.Criticalf("anti-cheat (ac1) failed.")
			logger.Criticalf("Please contact us at hello@codecrafters.io if you think this is a mistake.")
			return fmt.Errorf("anti-cheat (ac1) failed")
		}
	}

	return nil
}

package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type AntiCheatTestCase struct {
	Pattern          string
	Input            string
	ExpectedExitCode int
}

func (t AntiCheatTestCase) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	executable := stageHarness.Executable.Clone()
	executable.TimeoutInMilliseconds = 1000

	result, err := executable.RunWithStdin([]byte(t.Input), "-E", t.Pattern)
	if err != nil && err.Error() == "execution timed out" {
		return nil
	}
	if result.ExitCode == t.ExpectedExitCode {
		return fmt.Errorf("anti-cheat (ac1) failed")
	}
	return nil
}

type AntiCheatTestCaseCollection []AntiCheatTestCase

func (c AntiCheatTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger

	matchCount := 0
	for _, testCase := range c {
		executable := stageHarness.Executable.Clone()
		executable.TimeoutInMilliseconds = 1000

		actualResult, err := executable.RunWithStdin([]byte(testCase.Input), "-E", testCase.Pattern)
		if err != nil && err.Error() == "execution timed out" {
			continue
		}
		if actualResult.ExitCode == testCase.ExpectedExitCode {
			matchCount++
		}
	}

	// Only if all anti-cheat test cases "fail", we fail the test
	if matchCount == len(c) {
		logger.Criticalf("anti-cheat (ac1) failed.")
		logger.Criticalf("Please contact us at hello@codecrafters.io if you think this is a mistake.")
		return fmt.Errorf("anti-cheat (ac1) failed")
	}

	return nil
}

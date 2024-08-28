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

func (t *AntiCheatTestCase) Run(stageHarness *test_case_harness.TestCaseHarness) error {
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

func RunAntiCheatTestCases(testCases []AntiCheatTestCase, stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger

	for _, testCase := range testCases {
		if testCase.Run(stageHarness) == nil {
			return nil
		}
	}
	logger.Criticalf("anti-cheat (ac1) failed.")
	logger.Criticalf("Please contact us at hello@codecrafters.io if you think this is a mistake.")
	return fmt.Errorf("anti-cheat (ac1) failed")
}

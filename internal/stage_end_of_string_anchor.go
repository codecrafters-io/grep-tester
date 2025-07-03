package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testEndOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	MoveGrepToTemp(stageHarness, stageHarness.Logger)

	testCases := []test_cases.StdinTestCase{
		{
			Pattern:          "cat$",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "cat$",
			Input:            "cats",
			ExpectedExitCode: 1,
		},
	}

	return test_cases.RunStdinTestCases(testCases, stageHarness)
}

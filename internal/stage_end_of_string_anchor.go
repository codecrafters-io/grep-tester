package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testEndOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []TestCase{
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

	return RunTestCases(testCases, stageHarness)
}

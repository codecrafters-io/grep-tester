package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testEndOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCases{
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

	return testCases.Run(stageHarness)
}

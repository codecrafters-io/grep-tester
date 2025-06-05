package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testAlternation(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "a (cat|dog)",
			Input:            "a cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "a (cat|dog)",
			Input:            "a cow",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^((Buffalo|buffalo)[ ]?)+$",
			Input:            "Buffalo buffalo Buffalo buffalo buffalo buffalo Buffalo buffalo",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^((W|w)ill(, | |'s )?)+$",
			Input:            "Will, will Will will Will Will's will",
			ExpectedExitCode: 0,
		},
	}

	return RunTestCases(testCases, stageHarness)
}

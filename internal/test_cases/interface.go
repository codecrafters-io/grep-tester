package test_cases

import "github.com/codecrafters-io/tester-utils/test_case_harness"

type TestCase interface {
	Run(stageHarness *test_case_harness.TestCaseHarness) error
}

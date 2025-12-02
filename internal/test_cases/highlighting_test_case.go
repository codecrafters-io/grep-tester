package test_cases

import (
	"fmt"
	"path"

	"github.com/codecrafters-io/grep-tester/internal/assertions"
	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type HighlightingTestCase struct {
	Pattern          string
	Stdin            string
	ExpectedExitCode int
	RunInsideTty     bool
	HighlightingMode utils.ColorMode
}

type HighlightingTestCaseCollection []HighlightingTestCase

func (c HighlightingTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	grepExecutable := stageHarness.Executable

	for _, testCase := range c {
		allArguments := []string{}

		// Initialize color argument
		colorArgument := ""
		if testCase.HighlightingMode != "" {
			colorArgument = fmt.Sprintf("--color=%s", testCase.HighlightingMode)
		}

		// Add color argument to arguments list
		if colorArgument != "" {
			allArguments = append(allArguments, colorArgument)
		}

		// Add pattern
		allArguments = append(allArguments, []string{"-E", testCase.Pattern}...)

		if testCase.RunInsideTty {
			logger.Infof("Running grep inside TTY")
		}
		logger.Infof("echo '%s' | $ ./%s %s -E '%s'", testCase.Stdin,
			path.Base(grepExecutable.Path),
			colorArgument,
			testCase.Pattern,
		)

		// Emulate grep
		emulatedResult := grep.EmulateGrep(allArguments, grep.EmulationOptions{
			Stdin:        []byte(testCase.Stdin),
			EmulateInTTY: testCase.RunInsideTty,
		})

		// Verify expected code against emulated result
		if testCase.ExpectedExitCode != emulatedResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, emulatedResult.ExitCode))
		}

		var actualResult executable.ExecutableResult
		var err error

		if testCase.RunInsideTty {
			actualResult, err = grepExecutable.RunInPtyWithStdin(executable.PTYOptions{
				UsePipeForStdin:  true,
				UsePipeForStderr: true,
			}, []byte(testCase.Stdin), allArguments...)
		} else {
			actualResult, err = grepExecutable.RunWithStdin([]byte(testCase.Stdin), allArguments...)
		}

		if err != nil {
			return err
		}

		exitCodeAssertion := assertions.ExitCodeAssertion{
			ExpectedExitCode: testCase.ExpectedExitCode,
		}

		if err := exitCodeAssertion.Run(actualResult, logger); err != nil {
			return err
		}

		// Assert stdout contents here
	}

	return nil
}

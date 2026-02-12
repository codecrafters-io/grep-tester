package test_cases

import (
	"fmt"
	"path"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/assertions"
	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type HighlightingTestCase struct {
	Pattern          string
	InputLines       []string
	ExpectedExitCode int
	RunInsideTty     bool
	ColorMode        utils.ColorMode
}

type HighlightingTestCaseCollection []HighlightingTestCase

func (c HighlightingTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	grepExecutable := stageHarness.Executable

	for _, testCase := range c {
		allArguments := []string{}
		allInputLines := strings.Join(testCase.InputLines, "\n")

		// Initialize color argument
		colorArgument := ""
		if testCase.ColorMode != "" {
			colorArgument = fmt.Sprintf("--color=%s", testCase.ColorMode)
		}

		// Add color argument to arguments list
		if colorArgument != "" {
			allArguments = append(allArguments, colorArgument)
		}

		// Add pattern
		allArguments = append(allArguments, []string{"-E", testCase.Pattern}...)
		pipeToCat := ""

		if !testCase.RunInsideTty {
			pipeToCat = "| cat"
		}

		logger.Infof("$ echo -ne %q | ./%s %s -E '%s' %s", allInputLines,
			path.Base(grepExecutable.Path),
			colorArgument,
			testCase.Pattern,
			pipeToCat,
		)

		// Emulate grep
		emulatedResult := grep.EmulateGrep(allArguments, grep.EmulationOptions{
			Stdin:        []byte(allInputLines),
			EmulateInTTY: testCase.RunInsideTty,
		})

		// Verify expected code against emulated result
		if testCase.ExpectedExitCode != emulatedResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, emulatedResult.ExitCode))
		}

		var actualResult executable.ExecutableResult
		var err error

		grepExecutable.ShouldUsePtyOutputStreams = testCase.RunInsideTty

		actualResult, err = grepExecutable.RunWithStdin([]byte(allInputLines), allArguments...)

		if err != nil {
			return err
		}

		exitCodeAssertion := assertions.ExitCodeAssertion{
			ExpectedExitCode: testCase.ExpectedExitCode,
		}

		if err := exitCodeAssertion.Run(actualResult, logger); err != nil {
			return err
		}

		highlightingAssertion := assertions.HighlightingAssertion{
			ExpectedOutput: string(emulatedResult.Stdout),
		}

		if err := highlightingAssertion.Run(actualResult, logger); err != nil {
			return err
		}
	}

	return nil
}

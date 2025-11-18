package assertions

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type ExitCodeAssertion struct {
	ExpectedExitCode int
	FailureHint      string
}

func (a ExitCodeAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	if result.ExitCode != a.ExpectedExitCode {
		hintString := ""
		if a.FailureHint != "" {
			hintString = fmt.Sprintf("\nHint: %s", a.FailureHint)
		}
		return fmt.Errorf("Expected exit code %v, got %v.%s", a.ExpectedExitCode, result.ExitCode, hintString)
	}

	logger.Successf("✓ Received exit code %d.", result.ExitCode)

	return nil
}

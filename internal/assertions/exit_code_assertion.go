package assertions

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type ExitCodeAssertion struct {
	ExpectedExitCode int
	FailureHintLine  string
}

func (a ExitCodeAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	if result.ExitCode != a.ExpectedExitCode {
		return fmt.Errorf("Expected exit code %v, got %v.\n%s", a.ExpectedExitCode, result.ExitCode, a.FailureHintLine)
	}

	logger.Successf("✓ Received exit code %d.", result.ExitCode)

	return nil
}

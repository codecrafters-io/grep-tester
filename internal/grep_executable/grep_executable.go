package grep_executable

import (
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type GrepExecutable struct {
	executable *executable.Executable
}

func NewGrepExecutable(stageHarness *test_case_harness.TestCaseHarness) *GrepExecutable {
	return &GrepExecutable{
		executable: stageHarness.NewExecutable(),
	}
}

func (e *GrepExecutable) Path() string {
	return e.executable.Path
}

func (e *GrepExecutable) Run(args ...string) (executable.ExecutableResult, error) {
	return e.executable.Run(args...)
}

func (e *GrepExecutable) RunWithStdin(stdin []byte, args ...string) (executable.ExecutableResult, error) {
	return e.executable.RunWithStdin(stdin, args...)
}

func (e *GrepExecutable) RunWithStdinInTty(stdin []byte, args ...string) (executable.ExecutableResult, error) {
	panic("Not implemented yet")
}

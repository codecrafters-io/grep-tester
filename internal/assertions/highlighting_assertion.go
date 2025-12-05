package assertions

import (
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type HighlightingAssertion struct {
	ExpectedAsciiSequence []byte
}

func (a HighlightingAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	clrfTranslatedExpectedAsciiSequence := crlfTranslation(a.ExpectedAsciiSequence)
	crlfTranslatedStdout := crlfTranslation(result.Stdout)
	maxTerminalWidth := max(len(clrfTranslatedExpectedAsciiSequence), len(crlfTranslatedStdout))

	virtualTerminal1 := virtual_terminal.NewCustomVT(16, maxTerminalWidth)
	defer virtualTerminal1.Close()

	virtualTerminal2 := virtual_terminal.NewCustomVT(16, maxTerminalWidth)
	defer virtualTerminal2.Close()

	if _, err := virtualTerminal1.Write(clrfTranslatedExpectedAsciiSequence); err != nil {
		return err
	}

	if _, err := virtualTerminal2.Write(clrfTranslatedExpectedAsciiSequence); err != nil {
		return err
	}

	expectedScreenState := virtualTerminal1.GetScreenState()
	actualScreenState := virtualTerminal2.GetScreenState()

	screenStateComparator := screenStateAssertion{
		logger:              logger,
		expectedScreenState: expectedScreenState,
	}

	if err := screenStateComparator.Run(actualScreenState); err != nil {
		return err
	}

	return nil
}

// crlfTranslation duplicates the input, replaces \n by \r\n and returns the new byte slice
// existing \r\n sequences in the input are not modified
func crlfTranslation(input []byte) (translated []byte) {
	translated = make([]byte, 0, len(input))

	for i := range input {

		// For non \n characters, append them as is
		if input[i] != '\n' {
			translated = append(translated, input[i])
			continue
		}

		// If previous character was not \r, replace \n with \r\n
		if i > 0 && input[i-1] != '\r' {
			translated = append(translated, '\r', '\n')
			continue
		}

		// If previous character was \r, leave the sequence as is
		translated = append(translated, '\n')
	}

	return translated
}

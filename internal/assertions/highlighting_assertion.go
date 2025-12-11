package assertions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// HighlightingAssertion is a single line assertion that asserts the output's content
// and highlighting color (bold-red: default color) against the expected sequence
// using a virtual terminal
type HighlightingAssertion struct {
	ExpectedOutput string

	// This is computed based on expected screen state
	matchesShouldbeHighlighted bool
	originalResult             executable.ExecutableResult
	logger                     *logger.Logger
	screenStateComparator      *screenStateComparator
}

func (a *HighlightingAssertion) getExpectedOutputOnRowIdx(lineIdx int) string {
	lines := utils.ProgramOutputToLines(a.ExpectedOutput)
	return lines[lineIdx]
}

func (a *HighlightingAssertion) getActualOutputOnLineIdx(lineIdx int) string {
	lines := utils.ProgramOutputToLines(string(a.originalResult.Stdout))
	return lines[lineIdx]
}

func (a *HighlightingAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	defer func() {
		// Reset the computed value(s) after the execution is over
		a.matchesShouldbeHighlighted = false
		a.originalResult = executable.ExecutableResult{}
		a.logger = nil
		a.screenStateComparator = nil
	}()

	a.originalResult = result
	a.logger = logger.Clone()

	// The dimensions for the VT will be the value that is maximum among the expected and actual output's width and height
	// 1 more than the max length because even in the case of empty input, we need to spawn a vt
	maxTerminalWidth := 1 + max(
		len(a.ExpectedOutput),
		len(result.Stdout),
	)

	expectedLines := utils.ProgramOutputToLines(a.ExpectedOutput)
	actualLines := utils.ProgramOutputToLines(string(result.Stdout))

	// we still need a virtual terminal with a single row
	maxTerminalHeight := 1 + max(
		len(expectedLines),
		len(actualLines),
	)

	virtualTerminal1 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal1.Close()

	virtualTerminal2 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal2.Close()

	if _, err := virtualTerminal1.WriteWithCRLFTranslation([]byte(a.ExpectedOutput)); err != nil {
		return err
	}

	if _, err := virtualTerminal2.WriteWithCRLFTranslation(result.Stdout); err != nil {
		return err
	}

	expectedScreenState := virtualTerminal1.GetScreenState()
	a.panicIfExpectedScreenStateisFlawed(expectedScreenState)
	actualScreenState := virtualTerminal2.GetScreenState()

	// Assert text contents first
	if err := a.assertTextContents(expectedScreenState, actualScreenState); err != nil {
		return err
	}

	// Assert highlighting
	return a.assertHighlighting(expectedScreenState, actualScreenState)
}

func (a *HighlightingAssertion) assertTextContents(expectedScreenState, actualScreenState *virtual_terminal.ScreenState) error {
	expectedLines := expectedScreenState.GetLinesOfTextUptoCursor()
	actualLines := actualScreenState.GetLinesOfTextUptoCursor()

	orderedLinesAssertion := OrderedLinesAssertion{
		ExpectedOutputLines: expectedLines,
	}

	// Simulate the result as if it would have been produced without highlighting it
	actualOutput := strings.Join(actualLines, "\n")

	actualResult := executable.ExecutableResult{
		Stdout: []byte(actualOutput),
	}

	return orderedLinesAssertion.Run(actualResult, a.logger)
}

func (a *HighlightingAssertion) assertHighlighting(expectedScreenState, actualScreenState *virtual_terminal.ScreenState) error {
	// Initialize comparison context
	a.screenStateComparator = newScreenStateComparator()

	expectedLines := expectedScreenState.GetLinesOfTextUptoCursor()
	actualLines := actualScreenState.GetLinesOfTextUptoCursor()

	rowsCount := len(expectedLines)

	// Hand off to comparator
	comparisonError := a.screenStateComparator.CompareHighlightingForNRows(expectedScreenState, actualScreenState, rowsCount)

	lastSuccessFulRowIndex := rowsCount - 1
	if comparisonError != nil {
		lastSuccessFulRowIndex = comparisonError.RowIdx - 1
	}

	// Log success messages for each row
	for rowIdx := 0; rowIdx <= lastSuccessFulRowIndex; rowIdx++ {
		if a.screenStateComparator.matchesShouldbeHighlighted {
			a.logger.Successf("✓ All matches in the line %q are highlighted", actualLines[rowIdx])
		} else {
			a.logger.Successf("✓ Line %q is not highlighted", actualLines[rowIdx])
		}
	}

	if comparisonError != nil {
		// Build error using the information from context
		return a.wrapComparisonError(comparisonError)
	}

	return nil
}

func (a *HighlightingAssertion) wrapComparisonError(compErr *ComparisonError) error {
	// Print comparison with cursor
	a.logger.Plainln(
		buildComparisonErrorMessageWithCursor(
			a.getExpectedOutputOnRowIdx(compErr.RowIdx),
			a.getActualOutputOnLineIdx(compErr.RowIdx),
			compErr.ColumnIdx,
		),
	)

	// Print success messages
	for _, successLog := range compErr.SuccessLogs {
		a.logger.Successln(successLog)
	}

	// Print error message
	a.logger.Errorln(fmt.Sprintf("⨯ %s", compErr.Message))

	// Print ANSI sequence comparison
	a.logger.Plainln(
		buildAnsiCodeMismatchComplaint(
			compErr.ExpectedCell.Style.String(),
			compErr.ActualCell.Style.String(),
		),
	)

	return errors.New("Wrong highlighting pattern")
}

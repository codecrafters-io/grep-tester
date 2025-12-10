package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// HighlightingAssertion is a single line assertion that asserts the output's content
// and highlighting color (bold-red: default color) against the expected sequence
// using a virtual terminal
type HighlightingAssertion struct {
	ExpectedOutput  string
	ExpectedMatches []string

	// This is computed based on expected screen state
	matchesShouldbeHighlighted bool
}

func (a *HighlightingAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	defer func() {
		// Reset the computed value(s) after the execution is over
		a.matchesShouldbeHighlighted = false
	}()

	// The dimensions for the VT will be the value that is maximum among the expected and actual output's width and height
	// 1 more than the max length because even in the case of empty input,
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
	if err := a.assertTextContents(expectedScreenState, actualScreenState, logger); err != nil {
		return err
	}

	// Assert highlighting
	return a.assertHighlighting(expectedScreenState, actualScreenState, result, logger)
}

func (a *HighlightingAssertion) assertTextContents(expectedScreenState, actualScreenState *virtual_terminal.ScreenState, logger *logger.Logger) error {
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

	return orderedLinesAssertion.Run(actualResult, logger)
}

func (a *HighlightingAssertion) assertHighlighting(expectedScreenState, actualScreenState *virtual_terminal.ScreenState, result executable.ExecutableResult, logger *logger.Logger) error {
	// Assert the first line on each terminal
	expectedRow := expectedScreenState.GetRowAtIndex(0)
	actualRow := actualScreenState.GetRowAtIndex(0)

	expectedCells := expectedRow.GetCellsArray()
	actualCells := actualRow.GetCellsArray()

	for cellIdx, expectedCell := range expectedCells {
		actualCell := actualCells[cellIdx]

		err := a.compareCells(expectedCell, actualCell)

		if err != nil {
			// We trim the (\r)\n character from the output so error message can be built
			// We arrive only if there is one row in the output
			// Checks in assertTextContents guarantees this
			expectedOutputLineWithoutCRLF := utils.ProgramOutputToLines(a.ExpectedOutput)[0]
			actualOutputLineWithoutCRLF := utils.ProgramOutputToLines(string(result.Stdout))[0]

			return fmt.Errorf(
				"%s\n%s\n%s\n%s",
				utils.BuildColoredErrorMessage(expectedOutputLineWithoutCRLF, actualOutputLineWithoutCRLF, cellIdx),
				err.Error(),
				// Raw output here
				fmt.Sprintf("Expected ANSI Sequence: %q", a.ExpectedOutput),
				fmt.Sprintf("Received ANSI Sequence: %q", string(result.Stdout)),
			)
		}
	}

	// Print highlighting summary only when highlighting should be done
	for _, match := range a.ExpectedMatches {
		if a.matchesShouldbeHighlighted {
			logger.Successf("✓ Match %q is highlighted", match)
		} else {
			logger.Successf("✓ Match %q is not highlighted", match)
		}
	}

	return nil
}

func (a *HighlightingAssertion) compareCells(expected *uv.Cell, actual *uv.Cell) error {

	// Rare cases: Doesn't occur unless the user intentionally does so
	if actual.Link != expected.Link {
		return fmt.Errorf("Expected hyperlink to be absent, but is present")
	}

	if actual.Width != expected.Width {
		return fmt.Errorf("Expected character to be of single width, got %d", actual.Width)
	}

	if actual.Style.Underline != expected.Style.Underline {
		return fmt.Errorf("Expected no underline, but is underlined")
	}

	if actual.Style.Bg != expected.Style.Bg {
		return fmt.Errorf("Expected no background color, got %s", utils.ColorToName(actual.Style.Bg))
	}

	// Foreground checks
	actualColorString := utils.ColorToName(actual.Style.Fg)

	shouldbeHighlighted := expected.Style.Fg == ansi.Red && expected.Style.Attrs == uv.AttrBold

	// If at least one cell is highlighted, it's expected that matches should highlighted
	if shouldbeHighlighted {
		a.matchesShouldbeHighlighted = true
	}

	expectedHighlighting := "bold-red(ANSI Code: 01;31)"

	// Foreground check 1: If the actual style is other than red or none
	if actual.Style.Fg != ansi.Red && actual.Style.Fg != nil {
		if shouldbeHighlighted {
			return fmt.Errorf("Expected %s, got %s", expectedHighlighting, actualColorString)
		}
		return fmt.Errorf("Expected no highlighting, got %s", actualColorString)
	}

	// Make these four cases as verbose as possible
	isBold := actual.Style.Attrs == uv.AttrBold
	isRed := actual.Style.Fg == ansi.Red

	if shouldbeHighlighted {
		if isBold && isRed {
			return nil
		}

		if !isBold && !isRed {
			return fmt.Errorf("Expected %s, got no highlight", expectedHighlighting)
		}

		if !isBold {
			return fmt.Errorf("Expected %s, got only red", expectedHighlighting)
		}

		return fmt.Errorf("Expected %s, got only bold", expectedHighlighting)
	}

	// No highlighting case
	if !isBold && !isRed {
		return nil
	}

	if isBold && isRed {
		return fmt.Errorf("Expected no highlight, got bold-red")
	}

	if isRed {
		return fmt.Errorf("Expected no highlight, got red")
	}

	return fmt.Errorf("Expected no highlight, got bold")
}

// panicIfExpectedScreenStateisFlawed will panic if:
// the screenstate has more than one line of output, or
// if any of the cells are flawed
// See panicIfExpectedCellIsFlawed() for the checks
func (a HighlightingAssertion) panicIfExpectedScreenStateisFlawed(expectedScreenState *virtual_terminal.ScreenState) {
	linesUptoCursor := expectedScreenState.GetLinesOfTextUptoCursor()

	// Assert that expected screen state only has one line in the output
	if len(linesUptoCursor) > 1 {

		outputLines := strings.Join(linesUptoCursor, "\n")

		panic(fmt.Sprintf(
			"Codecrafters Internal Error - Expected one line in output, grep returned %d:\n%s",
			len(linesUptoCursor),
			outputLines,
		))
	}

	for _, expectedRow := range expectedScreenState.GetAllRows() {
		for _, expectedCell := range expectedRow.GetCellsArray() {
			a.panicIfExpectedCellIsFlawed(expectedCell)
		}
	}

}

// panicIfExpectedCellIsFlawed panics if the expected cell
// contains hyperlink
// is not of mono-width
// contains underline
// has non-empty background color
// has foreground styling other than empty or bold-red (used for highlighting)
func (a HighlightingAssertion) panicIfExpectedCellIsFlawed(expectedCell *uv.Cell) {
	emptyCell := uv.EmptyCell

	// Expected cell should have no link
	if expectedCell.Link != emptyCell.Link {
		panic("Codecrafters Internal Error - Expected cell contains hyperlink")
	}

	// Expected cell should be of mono-width
	if expectedCell.Width != emptyCell.Width {
		panic("Codecrafters Internal Error - Expected cell is not of unit width")
	}

	// Expected cell should have no underline
	if expectedCell.Style.Underline != emptyCell.Style.Underline {
		panic("Codecrafters Internal Error - Expected cell is underlined")
	}

	// Expected cell should have no background color
	if expectedCell.Style.Bg != emptyCell.Style.Bg {
		panic("Codecrafters Internal Error - Expected cell doesn't have plain background")
	}

	// Expected cell should either be bold red, or with no styling
	if !(expectedCell.Style.Fg == ansi.Red && expectedCell.Style.Attrs == uv.AttrBold) &&
		!(expectedCell.Style.Fg == emptyCell.Style.Fg && expectedCell.Style.Attrs == emptyCell.Style.Attrs) {
		panic(
			fmt.Sprintf(
				"Codecrafters Internal Error - Expected cell should be either bold-red, or none, got color: %s, attribute: %d",
				utils.ColorToName(expectedCell.Style.Fg),
				expectedCell.Style.Attrs,
			),
		)
	}
}

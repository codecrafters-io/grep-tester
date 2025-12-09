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
	ExpectedAsciiSequence []byte
	ExpectedMatches       []string
}

func (a HighlightingAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	// The dimensions for the VT will be the value that is maximum among the expected and actual output's width and height
	maxTerminalWidth := max(
		len(a.ExpectedAsciiSequence),
		len(result.Stdout),
	)

	maxTerminalHeight := max(
		len(strings.Split(string(a.ExpectedAsciiSequence), "\n")),
		len(strings.Split(string(result.Stdout), "\n")),
	)

	virtualTerminal1 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal1.Close()

	virtualTerminal2 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal2.Close()

	if _, err := virtualTerminal1.Write(a.ExpectedAsciiSequence); err != nil {
		return err
	}

	if _, err := virtualTerminal2.Write(result.Stdout); err != nil {
		return err
	}

	expectedScreenState := virtualTerminal1.GetScreenState()
	// This is equivalent to panic() after mismatch between the exit codes from
	// expected value and emulated grep's value: For eg. stdin_test_case:36
	a.panicIfExpectedScreenStateisFlawed(expectedScreenState)

	actualScreenState := virtualTerminal2.GetScreenState()

	outputLine, err := a.assertContents(expectedScreenState, actualScreenState, logger)

	if err != nil {
		return err
	}

	return a.assertHighlighting(expectedScreenState, actualScreenState, outputLine, logger)
}

func (a HighlightingAssertion) assertContents(expectedScreenState, actualScreenState *virtual_terminal.ScreenState, logger *logger.Logger) (string, error) {
	expectedLines := expectedScreenState.GetLinesUptoCursor()
	actualLines := actualScreenState.GetLinesUptoCursor()

	orderedLinesAssertion := OrderedLinesAssertion{
		ExpectedOutputLines: expectedLines,
	}

	// Simulate the result as if it would have been produced without highlighting it
	actualOutput := strings.Join(actualLines, "\n")

	actualResult := executable.ExecutableResult{
		Stdout: []byte(actualOutput),
	}

	return strings.Split(actualOutput, "\n")[0], orderedLinesAssertion.Run(actualResult, logger)
}

func (a HighlightingAssertion) assertHighlighting(expectedScreenState, actualScreenState *virtual_terminal.ScreenState, outputLine string, logger *logger.Logger) error {
	// Assert the first line on each terminal
	expectedRow := expectedScreenState.GetRow(0)
	actualRow := actualScreenState.GetRow(0)

	for i, expectedCell := range expectedRow.GetCellsArray() {
		actualCell := actualRow.GetCellsArray()[i]
		err := a.compareCells(expectedCell, actualCell)
		if err != nil {
			return fmt.Errorf(
				"%s\n%s\n%s\n%s",
				utils.BuildColoredErrorMessage(string(a.ExpectedAsciiSequence), outputLine, i),
				err.Error(),
				fmt.Sprintf("Expected ANSI Sequence: %q", string(a.ExpectedAsciiSequence)),
				fmt.Sprintf("Received ANSI Sequence: %q", string(outputLine)),
			)
		}
	}

	for _, match := range a.ExpectedMatches {
		logger.Successf("✓ Match %q is highlighted", match)
	}

	return nil
}

func (a HighlightingAssertion) compareCells(expected *uv.Cell, actual *uv.Cell) error {
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
	expectedHighlighting := "bold-red(ANSI Code: 01;31)"

	// Foreground check 1: If the actual style is other than red or none
	if actual.Style.Fg != ansi.Red && actual.Style.Fg != nil {
		if shouldbeHighlighted {
			return fmt.Errorf("Expected %s, got %s", expectedHighlighting, actualColorString)
		}
		return fmt.Errorf("Expected no highlighting, got %s", actualColorString)
	}

	// Make these four cases as verbose as possible
	isBold := actual.Style.Attrs == 1
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

func (a HighlightingAssertion) panicIfExpectedScreenStateisFlawed(expectedScreenState *virtual_terminal.ScreenState) {
	linesUptoCursor := expectedScreenState.GetLinesUptoCursor()

	// Assert that expected screen state only has one line in the output
	if len(linesUptoCursor) > 1 {

		outputLines := strings.Join(linesUptoCursor, "\n")

		panic(fmt.Sprintf(
			"Codecrafters Internal Error - Expected one line in output, grep returned %d:\n%s",
			len(linesUptoCursor),
			outputLines,
		))
	}

	for _, expectedRow := range expectedScreenState.GetRows() {
		for _, expectedCell := range expectedRow.GetCellsArray() {
			a.panicIfExpectedCellIsFlawed(expectedCell)
		}
	}

}

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
				&expectedCell.Style.Attrs,
			),
		)
	}
}

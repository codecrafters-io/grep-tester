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
	// The dimensions for the VT will be the maximum
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
	actualScreenState := virtualTerminal2.GetScreenState()

	expectedLines := expectedScreenState.GetLinesUptoCursor()
	actualLines := actualScreenState.GetLinesUptoCursor()

	// If expected lines are more than 1, panic: there is something wrong with the expected value
	if len(expectedLines) > 1 {
		expectedLinesString := strings.Join(expectedLines, "\n")
		panic(fmt.Sprintf("Codecrafters Internal Error - Multiple expected lines found in HighlightingAssertion: \n%s", expectedLinesString))
	}

	if err := a.assertContents(expectedLines, actualLines, logger); err != nil {
		return err
	}

	return a.assertHighlighting(expectedScreenState, actualScreenState, result, logger)
}

func (a HighlightingAssertion) assertContents(expectedLines []string, actualLines []string, logger *logger.Logger) error {
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

func (a HighlightingAssertion) assertHighlighting(expectedScreenState, actualScreenState *virtual_terminal.ScreenState, result executable.ExecutableResult, logger *logger.Logger) error {
	// Assert the first line on each terminal
	expectedRow := expectedScreenState.GetRow(0)
	actualRow := actualScreenState.GetRow(0)

	for i, expectedCell := range expectedRow.GetCellsArray() {
		actualCell := actualRow.GetCellsArray()[i]
		err := a.compareCells(expectedCell, actualCell)
		if err != nil {
			resultLine := utils.ProgramOutputToLines(string(result.Stdout))[0]
			logger.Plainf(utils.BuildColoredErrorMessage(string(a.ExpectedAsciiSequence), resultLine, i))
			logger.Errorf("Expected ANSI Sequence: %q", string(a.ExpectedAsciiSequence))
			logger.Errorf("Received ANSI Sequence: %q", string(resultLine))
			return err
		}
	}

	return nil
}

func (a HighlightingAssertion) compareCells(expected *uv.Cell, actual *uv.Cell) error {
	a.panicIfExpectedIsFlawed(expected)

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
		return fmt.Errorf("Expected no background color, got %s", utils.ColorToString(actual.Style.Bg))
	}

	// Check for bold-red combination
	expectedBoldRed := expected.Style.Fg == ansi.Red && expected.Style.Attrs == uv.AttrBold
	actualBoldRed := actual.Style.Fg == ansi.Red && actual.Style.Attrs == uv.AttrBold

	if expectedBoldRed {
		// If we got bold-red, success!
		if actualBoldRed {
			return nil
		}

		// Check what we actually got
		hasBold := actual.Style.Attrs == uv.AttrBold
		hasRed := actual.Style.Fg == ansi.Red
		hasOtherColor := actual.Style.Fg != nil && actual.Style.Fg != ansi.Red

		if !hasBold && !hasRed {
			return fmt.Errorf("Expected bold-red, got none")
		} else if hasBold && !hasRed && !hasOtherColor {
			return fmt.Errorf("Expected bold-red, got only bold")
		} else if hasRed && !hasBold {
			return fmt.Errorf("Expected bold-red, got only red")
		} else if hasOtherColor {
			colorStr := utils.ColorToString(actual.Style.Fg)
			if hasBold {
				return fmt.Errorf("Expected bold-red, got bold-%s", colorStr)
			}
			return fmt.Errorf("Expected bold-red, got %s", colorStr)
		}
	}

	// Handle other cases (non-bold-red expectations)
	if actual.Style.Fg != expected.Style.Fg {
		if expected.Style.Fg == ansi.Red {
			return fmt.Errorf("Expected red foreground color, got %s", utils.ColorToString(actual.Style.Fg))
		}
		if expected.Style.Fg == nil {
			return fmt.Errorf("Expected no foreground color, got %s", utils.ColorToString(actual.Style.Fg))
		}
	}

	if actual.Style.Attrs != expected.Style.Attrs {
		if expected.Style.Attrs == uv.AttrBold {
			return fmt.Errorf("Expected bold attribute: 1, got %d", actual.Style.Attrs)
		}
		return fmt.Errorf("Expected no attributes: 0, got %d", actual.Style.Attrs)
	}

	return nil
}

func (a HighlightingAssertion) panicIfExpectedIsFlawed(e *uv.Cell) {
	emp := uv.EmptyCell
	if e.Link != emp.Link {
		panic("Codecrafters Internal Error - Expected cell contains hyperlink")
	}

	if e.Width != uv.EmptyCell.Width {
		panic("Codecrafters Internal Error - Expected cell is not of unit width")
	}

	if e.Style.Underline != uv.EmptyCell.Style.Underline {
		panic("Codecrafters Internal Error - Expected cell is underlined")
	}

	if e.Style.Bg != uv.EmptyCell.Style.Bg {
		panic("Codecrafters Internal Error - Expected cell doesn't have plain background")
	}

	if e.Style.Fg != ansi.Red && e.Style.Fg != uv.EmptyCell.Style.Fg {
		panic("Codecrafters Internal Error - Expected cell neither has plain foreground, nor red foreground")
	}

	if e.Style.Attrs != uv.EmptyCell.Style.Attrs && e.Style.Attrs != uv.AttrBold {
		panic("Codecrafters Internal Error - Expected cell has attributes other than bold or none")
	}
}

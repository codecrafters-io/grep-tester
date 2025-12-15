package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
)

// validateExpectedScreenState will panic if:
// the screenstate has more than expected number of lines in the output, or
// if any of the cells are flawed
// See validateExpectedCell() for the checks
func (a HighlightingAssertion) validateExpectedScreenState(expectedScreenState *virtual_terminal.ScreenState) {
	linesUptoCursor := expectedScreenState.GetLinesOfTextUptoCursor()

	expectedLinesCount := len(utils.ProgramOutputToLines(a.ExpectedOutput))

	// Assert that expected screen state has exact expected number of line in the output
	if len(linesUptoCursor) != expectedLinesCount {

		outputLines := strings.Join(linesUptoCursor, "\n")

		panic(fmt.Sprintf(
			"Codecrafters Internal Error - Expected %d line in output, grep returned %d:\n%s",
			expectedLinesCount,
			len(linesUptoCursor),
			outputLines,
		))
	}

	for rowIdx := range expectedScreenState.GetRowsCount() {
		for colIdx := range expectedScreenState.GetColumnsCount() {
			a.validateExpectedCell(expectedScreenState.MustGetCellAtPosition(rowIdx, colIdx))
		}
	}

}

// validateExpectedCell panics if the expected cell has any of the following properties
// contains hyperlink
// is not of mono-width
// contains underline
// has non-empty background color
// has foreground styling other than empty or bold-red (used for highlighting)
func (a HighlightingAssertion) validateExpectedCell(expectedCell *uv.Cell) {
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
				getFgColorName(expectedCell.Style.Fg),
				expectedCell.Style.Attrs,
			),
		)
	}
}

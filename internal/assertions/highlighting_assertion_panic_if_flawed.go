package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
)

// panicIfExpectedScreenStateisFlawed will panic if:
// the screenstate has more than one line of output, or
// if any of the cells are flawed
// See panicIfExpectedCellIsFlawed() for the checks
func (a HighlightingAssertion) panicIfExpectedScreenStateisFlawed(expectedScreenState *virtual_terminal.ScreenState) {
	linesUptoCursor := expectedScreenState.GetLinesOfTextUptoCursor()

	expectedLinesCount := len(utils.ProgramOutputToLines(a.ExpectedOutput))

	// Assert that expected screen state only has one line in the output
	if len(linesUptoCursor) > expectedLinesCount {

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

// panicIfExpectedCellIsFlawed panics if the expected cell has any of the following properties
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
				getFgColorName(expectedCell.Style.Fg),
				expectedCell.Style.Attrs,
			),
		)
	}
}

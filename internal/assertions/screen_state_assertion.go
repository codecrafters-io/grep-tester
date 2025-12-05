package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type screenStateAssertion struct {
	logger              *logger.Logger
	expectedScreenState *virtual_terminal.ScreenState
}

func (s *screenStateAssertion) Run(actual *virtual_terminal.ScreenState) error {
	// Dimension check and panic if not same
	if !s.expectedScreenState.HasSameDimensionAs(actual) {
		panic("Codecrafters Internal Error - expected screen state and actual screen states have different dimensions")
	}

	// Compare contents first (even if there's no guarantee that highlighting sequences is correct)
	// This is because printing matching lines is important than highlighting actual matches, and we want to
	// deliver this message first
	if err := s.compareContents(actual); err != nil {
		return err
	}

	rogueCellRow, rogueCellColumn, err := s.checkForValidColors(actual)
	if err != nil {
		return fmt.Errorf("Error on cell at row %d and column %d:", rogueCellRow, rogueCellColumn, err)
	}

	return s.assertMatchHighlight(actual)
}

// compareContents only compares the contents (string) of the two screen states
func (s *screenStateAssertion) compareContents(actual *virtual_terminal.ScreenState) error {

	expectedOutputLines := s.expectedScreenState.StringArray()
	actualOutputLines := actual.StringArray()

	orderedLinesAssertion := OrderedLinesAssertion{
		ExpectedOutputLines: expectedOutputLines,
	}

	executableResult := executable.ExecutableResult{
		Stdout: []byte(strings.Join(actualOutputLines, "\n")),
	}

	return orderedLinesAssertion.Run(executableResult, s.logger)
}

func (s *screenStateAssertion) checkForValidColors(actual *virtual_terminal.ScreenState) (int, int, error) {
	for i, row := range actual.GetRows() {
		for j, cell := range row.GetCellsArray() {
			expectedCell := s.expectedScreenState.GetRow(i).GetCellsArray()[j]
			err := getStylingError(cell, expectedCell)
			if err != nil {
				return i, j, err
			}
		}
	}
	return 0, 0, nil
}

func (s *screenStateAssertion) assertMatchHighlight(actual *virtual_terminal.ScreenState) error {
	for i, row := range actual.GetRows() {
		for j, cell := range row.GetCellsArray() {
			expectedCell := s.expectedScreenState.GetRow(i).GetCellsArray()[j]
			if !cell.Style.Equal(&expectedCell.Style) {
				return fmt.Errorf("Expected cell at (%d, %d) to be highlighted/ or not highlighted", i, j)
			}
		}
	}

	return nil
}

func getStylingError(expected *uv.Cell, actual *uv.Cell) (err error) {
	emptyCell := uv.EmptyCell

	// Width check
	if actual.Width != expected.Width {
		return fmt.Errorf("Character is not of expected normal width")
	}

	// Link check
	if actual.Link != expected.Link {
		return fmt.Errorf("Character has hyperlink associated with it")
	}

	// Underline check
	if actual.Style.Underline != 0 {
		return fmt.Errorf("Character is underlined")
	}

	// Background color check
	if expected.Style.Bg != emptyCell.Style.Bg {
		r, g, b, a := actual.Style.Bg.RGBA()
		re, ge, be, ae := expected.Style.Bg.RGBA()
		return fmt.Errorf(
			"Character has background color with RGBA (%d, %d, %d, %d), expected (%d, %d, %d, %d)",
			r, g, b, a, re, ge, be, ae,
		)
	}

	// Foreground color check
	if expected.Style.Fg != emptyCell.Style.Fg && expected.Style.Fg != ansi.Red {
		r, g, b, a := actual.Style.Bg.RGBA()
		re, ge, be, ae := expected.Style.Bg.RGBA()
		return fmt.Errorf(
			"Character has foreground color with RGBA (%d, %d, %d, %d), expected RGBA (%d, %d, %d, %d)", r, g, b, a, re, ge, be, ae,
		)
	}

	// Attribute check
	if expected.Style.Attrs != 1 && expected.Style.Attrs != 0 {
		return fmt.Errorf(
			"Character has attribute is %d, expected %d", actual.Style.Attrs, expected.Style.Attrs,
		)
	}

	// Foreground color check
	return nil
}

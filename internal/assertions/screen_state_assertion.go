package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type screenStateAssertion struct {
	logger              *logger.Logger
	expectedScreenState *virtual_terminal.ScreenState
}

func (s *screenStateAssertion) Run(actual *virtual_terminal.ScreenState, actualStdout string, expectedStdout string, highlights []string) error {
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

	_, rogueCellColumn, err := s.checkForValidColors(actual)
	if err != nil {
		// need original ascii sequence
		// Need simulated ascii sequence
		// rogueCellRow := rogueCellRow
		rogueCellColumn := rogueCellColumn
		s.logger.Plainf("Expected:%s", expectedStdout)
		s.logger.Plainf("Found   :%s", actualStdout[:len(actualStdout)-1])
		s.logger.Plainf("%s^", strings.Repeat(" ", rogueCellColumn+9))
		return fmt.Errorf("Highlight missing")
	}

	for _, h := range highlights {
		s.logger.Successf("✓ Match %q is highlighted", h)
	}

	return nil
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

func getStylingError(expected *uv.Cell, actual *uv.Cell) (err error) {

	// Width check
	if actual.Width != expected.Width {
		return fmt.Errorf("Character is not of expected normal width")
	}

	// Link check
	if actual.Link != expected.Link {
		return fmt.Errorf("Character has hyperlink associated with it")
	}

	// Underline check
	if actual.Style.Underline != expected.Style.Underline {
		return fmt.Errorf("Character is underlined")
	}

	// Background color check
	if actual.Style.Bg != expected.Style.Bg {
		return fmt.Errorf(
			"Character has background color with %v, expected %v", actual.Style.Bg, expected.Style.Bg,
		)
	}

	// Foreground color check
	if actual.Style.Fg != expected.Style.Fg {
		// r, g, b, a := actual.Style.Bg.RGBA()
		// re, ge, be, ae := expected.Style.Bg.RGBA()
		return fmt.Errorf(
			"Character has foreground color %v, expected %v", actual.Style.Fg, expected.Style.Fg,
		)
	}

	// Attribute check
	if actual.Style.Attrs != expected.Style.Attrs {
		return fmt.Errorf(
			"Character has attribute is %d, expected %d", actual.Style.Attrs, expected.Style.Attrs,
		)
	}

	// Foreground color check
	return nil
}

package virtual_terminal

import (
	"fmt"

	uv "github.com/charmbracelet/ultraviolet"
)

type CursorPosition struct {
	RowIndex    int
	ColumnIndex int
}

// ScreenState is a representation of the screen state at a given point in time
type ScreenState struct {
	rows           []*row
	cursorPosition CursorPosition
}

func NewScreenState(rawCellMatrix []*row, cursorPosition CursorPosition) *ScreenState {
	columnsCount := rawCellMatrix[0].GetCellsCount()

	// Dimensions check to ensure rectangular shape
	for _, row := range rawCellMatrix {
		if row.GetCellsCount() != columnsCount {
			panic("Codecrafters Internal Error - NewScreenState: rawCellMatrix is not rectangular")
		}
	}

	return &ScreenState{
		rows:           rawCellMatrix,
		cursorPosition: cursorPosition,
	}
}

// GetRowsCount returns the number of rows in the ScreenState
func (s *ScreenState) GetRowsCount() int {
	return len(s.rows)
}

// GetColumnsCount returns the number of columns in the ScreenState
func (s *ScreenState) GetColumnsCount() int {
	return s.rows[0].GetCellsCount()
}

// GetCursorPosition returns the cursor position in the screen state
func (s *ScreenState) GetCursorPosition() CursorPosition {
	return CursorPosition{
		RowIndex:    s.cursorPosition.RowIndex,
		ColumnIndex: s.cursorPosition.ColumnIndex,
	}
}

// MustGetCellAtPosition returns a copy of the cell at (rowIdx, colIdx)
func (s *ScreenState) MustGetCellAtPosition(rowIdx int, colIdx int) *uv.Cell {
	row := s.mustGetRowAtIndex(rowIdx)

	if colIdx < 0 || colIdx >= len(row.cells) {
		panic(fmt.Sprintf("Codecrafters Internal Error - Cannot get cell at (RowIdx=%d,ColIdx=%d) - Insufficient columns", rowIdx, colIdx))
	}

	return row.cells[colIdx].Clone()
}

// GetLinesOfTextUptoCursor returns the content of all the rows up to the row in which
// the cursor is present. For the row where cursor is present, content up to the cursor is returned
func (s *ScreenState) GetLinesOfTextUptoCursor() []string {
	result := []string{}

	for i := range s.cursorPosition.RowIndex + 1 {
		currentRowContent := s.mustGetRowAtIndex(i).getTextContents()
		result = append(result, currentRowContent)
	}

	// Exclude the row if the cursor is present in the beginning of the line
	if len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}

	return result
}

func (s *ScreenState) mustGetRowAtIndex(idx int) *row {
	if idx >= len(s.rows) || idx < 0 {
		panic(fmt.Sprintf("Codecrafters Internal Error - Cannot get row at index %d ", idx))
	}

	return s.rows[idx]
}

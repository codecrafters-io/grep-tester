package virtual_terminal

type CursorPosition struct {
	RowIndex    int
	ColumnIndex int
}

// ScreenState is a representation of the screen state at a given point in time
type ScreenState struct {
	rows           []*Row
	cursorPosition CursorPosition
}

func NewScreenState(rawCellMatrix []*Row, cursorPosition CursorPosition) *ScreenState {
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

// GetRowAtIndex returns a pointer to a copy of the row at index 'rowIndex'
// If the index is negative, or greater than, or equal to the number of rows, nil will be returned
func (s *ScreenState) GetRowAtIndex(idx int) *Row {
	if idx >= len(s.rows) || idx < 0 {
		return nil
	}

	return s.rows[idx].Clone()
}

// GetRows returns a slice of copy of all the rows in the Screenstate
func (s *ScreenState) GetAllRows() []*Row {
	rows := make([]*Row, len(s.rows))

	for i, row := range s.rows {
		rows[i] = row.Clone()
	}

	return rows
}

// GetLinesOfTextUptoCursor returns the content of all the rows up to the row in which
// the cursor is present. For the row where cursor is present, content up to the cursor is returned
func (s *ScreenState) GetLinesOfTextUptoCursor() []string {
	result := []string{}

	for i := range s.cursorPosition.RowIndex {
		currentRowContent := s.GetRowAtIndex(i).GetContents()
		result = append(result, currentRowContent)
	}

	lastRowIndex := s.cursorPosition.RowIndex
	lastRowContent := s.GetRowAtIndex(lastRowIndex).GetContents()
	lastRowContentBeforeCursor := lastRowContent[:s.cursorPosition.ColumnIndex]

	if lastRowContentBeforeCursor != "" {
		result = append(result, lastRowContentBeforeCursor)
	}

	return result
}

// HasSameDimensionAs returns true if the receiver's dimensions are the same
// as the expected screen state's dimensions
func (s *ScreenState) HasSameDimensionAs(expectedScreenState *ScreenState) bool {
	// Verify rows count
	if s.GetRowsCount() != expectedScreenState.GetRowsCount() {
		return false
	}

	// Checking the cells count of first row suffices, because of the checks in the constructor
	return s.GetRowAtIndex(0).GetCellsCount() == expectedScreenState.GetRowAtIndex(0).GetCellsCount()
}

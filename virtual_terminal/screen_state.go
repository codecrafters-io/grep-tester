package virtual_terminal

type CursorPosition struct {
	RowIndex    int
	ColumnIndex int
}

// ScreenState is a representation of the screen state at a given point in time
type ScreenState struct {
	// rows is always of size [rows x cols].
	//
	// Empty cells are represented by " " (space)
	rows []*Row

	cursorPosition CursorPosition
}

func NewScreenState(rawCellMatrix []*Row, cursorPosition CursorPosition) *ScreenState {

	return &ScreenState{
		rows:           rawCellMatrix,
		cursorPosition: cursorPosition,
	}
}

func (s *ScreenState) GetRow(rowIndex int) *Row {
	cursorCellIndex := -1

	if s.cursorPosition.RowIndex == rowIndex {
		cursorCellIndex = s.cursorPosition.ColumnIndex
	}

	return &Row{
		cellsArray:      s.rows[rowIndex].cellsArray,
		cursorCellIndex: cursorCellIndex,
	}
}

func (s *ScreenState) GetRows() []*Row {
	return s.rows
}

func (s *ScreenState) GetRowCount() int {
	return len(s.rows)
}

func (s *ScreenState) StringArray() []string {
	result := []string{}

	for i := range s.cursorPosition.RowIndex {
		currentRowContent := s.GetRow(i).String()
		result = append(result, currentRowContent)
	}

	lastRowIndex := s.cursorPosition.RowIndex
	lastRowContent := s.GetRow(lastRowIndex).String()
	lastRowContentBeforeCursor := lastRowContent[:s.cursorPosition.ColumnIndex]

	if lastRowContentBeforeCursor != "" {
		result = append(result, lastRowContentBeforeCursor)
	}

	return result
}

func (s *ScreenState) HasSameDimensionAs(base *ScreenState) bool {
	if s.GetRowCount() != base.GetRowCount() {
		return false
	}

	if len(s.GetRow(0).cellsArray) != len(base.GetRow(0).cellsArray) {
		return false
	}

	return true
}

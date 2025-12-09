package virtual_terminal

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
)

type VirtualTerminal struct {
	vt   *vt.Emulator
	rows int
	cols int
}

func NewCustomVT(rows, cols int) *VirtualTerminal {
	return &VirtualTerminal{
		vt:   vt.NewEmulator(cols, rows),
		rows: rows,
		cols: cols,
	}
}

func (vt *VirtualTerminal) Close() {
	vt.vt.Close()
}

func (vt *VirtualTerminal) WriteWithCRLFTranslation(p []byte) (n int, err error) {
	// The vt package does not provide a way to enable ONLCR on the virtual terminal
	// so, we perform the CRLF translation here and write to the vt
	tr := crlfTranslation(p)
	return vt.vt.Write(tr)
}

func (vt *VirtualTerminal) GetScreenState() *ScreenState {
	rows := make([]*Row, vt.rows)
	cursorRowIndex := vt.vt.CursorPosition().Y
	cursorColumnIndex := vt.vt.CursorPosition().X

	for i := 0; i < vt.rows; i++ {
		rows[i] = &Row{
			cells:           make([]*uv.Cell, vt.cols),
			cursorCellIndex: -1,
		}

		if cursorRowIndex == i {
			rows[i].cursorCellIndex = cursorColumnIndex
		}

		for j := 0; j < vt.cols; j++ {
			rows[i].cells[j] = vt.vt.CellAt(j, i)
		}
	}

	return NewScreenState(rows, CursorPosition{
		RowIndex:    vt.vt.CursorPosition().Y,
		ColumnIndex: vt.vt.CursorPosition().X,
	})
}

// crlfTranslation duplicates the input, replaces \n by \r\n and returns the new byte slice
// existing \r\n sequences in the input are not modified
// This is needed because the vt package doesn't provide us a method to enable ONLCR mode
func crlfTranslation(input []byte) (translated []byte) {
	translated = make([]byte, 0, len(input))

	for i := range input {

		// For non \n characters, append them as is
		if input[i] != '\n' {
			translated = append(translated, input[i])
			continue
		}

		// If previous character was not \r, replace \n with \r\n
		// Or, if the \n is present at the beginning, make the replacement
		if (i == 0) || (i > 0 && input[i-1] != '\r') {
			translated = append(translated, '\r', '\n')
			continue
		}

		// If previous character was \r, leave the sequence as is
		translated = append(translated, '\n')
	}

	return translated
}

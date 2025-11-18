package assertions

import "fmt"

func escapeLine(line string) string {
	if line == "" {
		return fmt.Sprintf("%q (empty)", line)
	}
	return fmt.Sprintf("%q", line)
}

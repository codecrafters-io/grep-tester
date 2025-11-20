package utils

import (
	"reflect"
	"testing"
)

func TestProgramOutputToLines(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", []string{}},
		{"abc\n", []string{"abc"}},
		{"\n\n", []string{"", ""}},
		{"abc\ndef", []string{"abc", "def"}},
		{"abc\ndef\nghi\n", []string{"abc", "def", "ghi"}},
	}

	for _, tt := range tests {
		got := ProgramOutputToLines(tt.in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("input %q: got %#v, want %#v", tt.in, got, tt.want)
		}
	}
}

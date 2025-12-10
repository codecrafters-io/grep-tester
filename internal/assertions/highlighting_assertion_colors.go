package assertions

import (
	"fmt"
	"image/color"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	fatihColor "github.com/fatih/color"
)

func buildComparisonErrorMessageWithCursor(expectedOutput string, actualOutput string, cursorIndex int) string {
	errorMsg := colorizeString(fatihColor.FgHiGreen, "\nExpected:")
	errorMsg += " \"" + expectedOutput + "\""
	errorMsg += "\n"
	errorMsg += colorizeString(fatihColor.FgHiRed, "Received:")
	errorMsg += " \"" + actualOutput + "\""
	offset := 11
	errorMsg += "\n" + strings.Repeat(" ", cursorIndex+offset) + "↑"
	return errorMsg
}

func buildAnsiCodeMismatchComplaint(expectedPattern string, actualPattern string) string {
	complaint := colorizeString(fatihColor.FgHiGreen, fmt.Sprintf("Expected ANSI code: %q\n", expectedPattern))
	complaint += colorizeString(fatihColor.FgHiRed, fmt.Sprintf("Received ANSI code: %q\n", actualPattern))
	return complaint
}

func colorizeString(colorToUse fatihColor.Attribute, msg string) string {
	c := fatihColor.New(colorToUse)
	return c.Sprint(msg)
}

// Map of ANSI color codes to color names
var ansiColorNames = map[string]string{
	// 3/4-bit foreground colors (30-37, 90-97)
	"30": "black",
	"31": "red",
	"32": "green",
	"33": "yellow",
	"34": "blue",
	"35": "magenta",
	"36": "cyan",
	"37": "white",
	"90": "bright black",
	"91": "bright red",
	"92": "bright green",
	"93": "bright yellow",
	"94": "bright blue",
	"95": "bright magenta",
	"96": "bright cyan",
	"97": "bright white",
	// 3/4-bit background colors (40-47, 100-107)
	"40":  "black",
	"41":  "red",
	"42":  "green",
	"43":  "yellow",
	"44":  "blue",
	"45":  "magenta",
	"46":  "cyan",
	"47":  "white",
	"100": "bright black",
	"101": "bright red",
	"102": "bright green",
	"103": "bright yellow",
	"104": "bright blue",
	"105": "bright magenta",
	"106": "bright cyan",
	"107": "bright white",
}

// colorCodeTocolorName returns the name of a color if it's a standard ANSI color,
// otherwise returns the raw color code
func colorCodeTocolorName(colorCode string) string {
	if name, ok := ansiColorNames[colorCode]; ok {
		return name + " (ANSI code " + colorCode + ")"
	}
	return fmt.Sprintf("color with ANSI code %q", colorCode)
}

func getFgColorName(c color.Color) string {
	if c == nil {
		return "white"
	}

	var colorCode ansi.Style
	colorCode = colorCode.ForegroundColor(c)
	colorCodeString := strings.TrimPrefix(colorCode.String(), "\x1b[")
	colorCodeString = strings.TrimSuffix(colorCodeString, "m")
	return colorCodeTocolorName(colorCodeString)
}

func getBgColorName(c color.Color) string {
	if c == nil {
		return "no color"
	}

	var colorCode ansi.Style
	colorCode = colorCode.BackgroundColor(c)
	colorCodeString := strings.TrimPrefix(colorCode.String(), "\x1b[")
	colorCodeString = strings.TrimSuffix(colorCodeString, "m")
	return colorCodeTocolorName(colorCodeString)
}

var attrMap = []struct {
	flag uint8
	name string
}{
	{uv.AttrBold, "Bold(1)"},
	{uv.AttrFaint, "Faint(2)"},
	{uv.AttrItalic, "Italic(3)"},
	{uv.AttrBlink, "Blink(5)"},
	{uv.AttrRapidBlink, "Rapid blink(6)"},
	{uv.AttrReverse, "Reverse(7)"},
	{uv.AttrConceal, "Conceal(8)"},
	{uv.AttrStrikethrough, "Strikethrough(9)"},
}

// attributesToNames converts a bitfield of attributes to their string names
func attributesToNames(attributes uint8) []string {
	if attributes == uv.AttrReset {
		return []string{}
	}

	var names []string

	for _, attr := range attrMap {
		if attributes&attr.flag != 0 {
			names = append(names, attr.name)
		}
	}

	return names
}

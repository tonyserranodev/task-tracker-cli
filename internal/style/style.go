// Package style provides ANSI terminal styling helpers for building the CLI UI.
package style

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ansiPattern matches ANSI escape sequences like \x1b[31m or \x1b[1;32m.
// These add bytes to the string but do not take up visible space in the terminal.
var ansiPattern = regexp.MustCompile("\x1b\\[[0-9;]*[a-zA-Z]")

// visibleWidth returns the number of runes that actually appear on screen,
// ignoring any ANSI color/style escape sequences.
func visibleWidth(s string) int {
	stripped := ansiPattern.ReplaceAllString(s, "")
	return utf8.RuneCountInString(stripped)
}

// PadRight pads s on the right with spaces until it reaches the given visible
// width. This is useful when s contains ANSI escape codes, because the padding
// is based on what a human sees rather than the raw byte length.
func PadRight(s string, width int) string {
	extra := width - visibleWidth(s)
	if extra <= 0 {
		return s
	}
	return s + strings.Repeat(" ", extra)
}

// Color represents an ANSI terminal color.
type Color int

const (
	_ Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

var colors = map[string]Color{
	"black":   Black,
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"blue":    Blue,
	"magenta": Magenta,
	"cyan":    Cyan,
	"white":   White,
}

// FG returns the ANSI escape sequence for this color as a foreground color.
// The zero value returns an empty string.
func (c Color) FG() string {
	if c == 0 {
		return ""
	}

	return fmt.Sprintf("\x1b[3%dm", c-1)
}

// BG returns the ANSI escape sequence for this color as a background color.
// The zero value returns an empty string.
func (c Color) BG() string {
	if c == 0 {
		return ""
	}

	return fmt.Sprintf("\x1b[4%dm", c-1)
}

// Decoration represents an ANSI text decoration.
type Decoration int

const (
	_ Decoration = iota
	Bold
	Dim
	Underline
)

var decorations = map[string]Decoration{
	"bold":      Bold,
	"dim":       Dim,
	"underline": Underline,
}

// Clear removes all styling and colors.
const Clear = "\x1b[0m"

// Decorate returns the ANSI escape sequence for this decoration.
func (d Decoration) Decorate() string {
	switch d {
	case Bold:
		return "\x1b[1m"
	case Dim:
		return "\x1b[2m"
	case Underline:
		return "\x1b[4m"
	default:
		return ""
	}
}

// Style holds foreground/background colors and a decoration for terminal output.
type Style struct {
	Foreground Color
	Background Color
	Decoration Decoration
}

// Apply returns text wrapped in the ANSI escape codes described by s.
func (s Style) Apply(text string) string {
	var b strings.Builder

	if s.Decoration != 0 {
		b.WriteString(s.Decoration.Decorate())
	}

	if s.Foreground != 0 {
		b.WriteString(s.Foreground.FG())
	}

	if s.Background != 0 {
		b.WriteString(s.Background.BG())
	}

	b.WriteString(text)
	b.WriteString(Clear)

	return b.String()
}

// Render applies styling arguments to text and returns the result.
// Valid arguments are color names ("red"), background colors ("bg-red"),
// and decorations ("bold", "dim", "underline"). One to three arguments are accepted.
func Render(text string, args ...string) (string, error) {
	if len(args) < 1 || len(args) > 3 {
		return "", errors.New("must pass between 1 and 3 arguments")
	}

	var style Style

	for _, arg := range args {
		if strings.HasPrefix(arg, "bg-") {
			colorName := strings.TrimPrefix(arg, "bg-")
			if color, ok := colors[colorName]; ok {
				style.Background = color
				continue
			}
		}

		if color, ok := colors[arg]; ok {
			style.Foreground = color
			continue
		}

		if dec, ok := decorations[arg]; ok {
			style.Decoration = dec
			continue
		}

		return "", fmt.Errorf("unknown argument passed %s", arg)
	}
	return style.Apply(text), nil
}

// Borders defines the runes used to draw a box in the terminal.
type Borders struct {
	TopLeft     rune // Corner at the top left.
	TopRight    rune // Corner at the top right.
	BottomLeft  rune // Corner at the bottom left.
	BottomRight rune // Corner at the bottom right.
	Horizontal  rune // Character used for the top and bottom edges.
	Vertical    rune // Character used for the left and right edges.
}

// SingleBorders draws boxes with single-line box-drawing characters.
var SingleBorders = Borders{
	TopLeft:     '┌',
	TopRight:    '┐',
	BottomLeft:  '└',
	BottomRight: '┘',
	Horizontal:  '─',
	Vertical:    '│',
}

// Box draws a bordered box around the given lines.
// A width of zero or less sizes the box to fit the longest line.
func Box(width int, lines []string, b Borders) string {
	innerWidth := width - 2

	// If width is zero or negative, size the box to fit the longest line.
	// The total box width becomes longest line + 2 (one border on each side).
	if width <= 0 {
		innerWidth = 0
		for _, line := range lines {
			if l := visibleWidth(line); l > innerWidth {
				innerWidth = l
			}
		}
	}

	top := string(b.TopLeft) + strings.Repeat(string(b.Horizontal), innerWidth) + string(b.TopRight)
	bottom := string(b.BottomLeft) + strings.Repeat(string(b.Horizontal), innerWidth) + string(b.BottomRight)

	var out strings.Builder

	out.WriteString(top)
	out.WriteString("\n")

	for _, line := range lines {
		pad := max(innerWidth-visibleWidth(line), 0)

		out.WriteString(string(b.Vertical))
		out.WriteString(line)
		out.WriteString(strings.Repeat(" ", pad))
		out.WriteString(string(b.Vertical))
		out.WriteString("\n")
	}

	out.WriteString(bottom)

	return out.String()
}

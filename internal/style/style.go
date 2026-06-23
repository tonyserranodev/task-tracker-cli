// Package style provides ANSI terminal styling helpers for building the CLI UI.
package style

import (
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
	Reset   Color = iota // Reset restores the terminal's default color.
	Black                // Black is the black ANSI color.
	Red                  // Red is the red ANSI color.
	Green                // Green is the green ANSI color.
	Yellow               // Yellow is the yellow ANSI color.
	Blue                 // Blue is the blue ANSI color.
	Magenta              // Magenta is the magenta ANSI color.
	Cyan                 // Cyan is the cyan ANSI color.
	White                // White is the white ANSI color.
)

// FG returns the ANSI escape sequence for this color as a foreground color.
// Reset returns the default foreground sequence rather than black.
func (c Color) FG() string {
	if c == Reset {
		return "\x1b[39m"
	}
	return fmt.Sprintf("\x1b[3%dm", c-1)
}

// BG returns the ANSI escape sequence for this color as a background color.
// Reset returns the default background sequence rather than black.
func (c Color) BG() string {
	if c == Reset {
		return "\x1b[49m"
	}
	return fmt.Sprintf("\x1b[4%dm", c-1)
}

const (
	Bold      = "\x1b[1m" // Bold enables bold text.
	Dim       = "\x1b[2m" // Dim enables dim text.
	Underline = "\x1b[4m" // Underline enables underlined text.
	Clear     = "\x1b[0m" // Clear resets all styles and colors.
)

// Style holds foreground/background colors and bold formatting for terminal output.
type Style struct {
	Foreground Color
	Background Color
	Bold       bool
}

// Apply returns text wrapped in the ANSI escape codes described by s.
func (s Style) Apply(text string) string {
	var b strings.Builder

	// Only emit color codes when a color was explicitly requested.
	// Reset is skipped so it does not set the foreground/background to black.
	if s.Foreground != Reset {
		b.WriteString(s.Foreground.FG())
	}
	if s.Background != Reset {
		b.WriteString(s.Background.BG())
	}

	if s.Bold {
		b.WriteString(Bold)
	}

	b.WriteString(text)
	b.WriteString(Clear)

	return b.String()
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

	//Draw top border
	out.WriteString(top)
	out.WriteString("\n")

	for _, line := range lines {
		// Prevent pad from going negative
		pad := max(
			innerWidth-visibleWidth(line), 0)

		out.WriteString(string(b.Vertical))
		out.WriteString(line)
		out.WriteString(strings.Repeat(" ", pad))
		out.WriteString(string(b.Vertical))
		out.WriteString("\n")
	}

	//Draw bottom border
	out.WriteString(bottom)

	return out.String()
}

package style

import (
	"strings"
	"testing"
)

func TestVisibleWidth(t *testing.T) {
	tt := map[string]struct {
		input string
		want  int
	}{
		"plain ascii":       {"hello", 5},
		"with ansi escapes": {"\x1b[31mhello\x1b[0m", 5},
		"empty":             {"", 0},
		"only ansi":         {"\x1b[1;32m", 0},
		"unicode runes":     {"こんにちは", 5},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := visibleWidth(tc.input)
			if got != tc.want {
				t.Errorf("visibleWidth(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tt := map[string]struct {
		input string
		width int
		want  string
	}{
		"pad shorter":   {"hi", 5, "hi   "},
		"no pad exact":  {"hello", 5, "hello"},
		"no pad longer": {"hello world", 5, "hello world"},
		"pad with ansi": {"\x1b[31mhi\x1b[0m", 5, "\x1b[31mhi\x1b[0m   "},
		"pad unicode":   {"こ", 3, "こ  "},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := PadRight(tc.input, tc.width)
			if got != tc.want {
				t.Errorf("PadRight(%q, %d) = %q, want %q", tc.input, tc.width, got, tc.want)
			}
		})
	}
}

func TestColorFG(t *testing.T) {
	tt := map[string]struct {
		color Color
		want  string
	}{
		"reset": {Reset, "\x1b[39m"},
		"red":   {Red, "\x1b[31m"},
		"green": {Green, "\x1b[32m"},
		"blue":  {Blue, "\x1b[34m"},
		"white": {White, "\x1b[37m"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := tc.color.FG()
			if got != tc.want {
				t.Errorf("%v.FG() = %q, want %q", tc.color, got, tc.want)
			}
		})
	}
}

func TestColorBG(t *testing.T) {
	tt := map[string]struct {
		color Color
		want  string
	}{
		"reset": {Reset, "\x1b[49m"},
		"red":   {Red, "\x1b[41m"},
		"blue":  {Blue, "\x1b[44m"},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := tc.color.BG()
			if got != tc.want {
				t.Errorf("%v.BG() = %q, want %q", tc.color, got, tc.want)
			}
		})
	}
}

func TestStyleApply(t *testing.T) {
	tt := map[string]struct {
		style Style
		input string
		want  string
	}{
		"foreground only": {
			style: Style{Foreground: Red},
			input: "x",
			want:  "\x1b[31mx\x1b[0m",
		},
		"background only": {
			style: Style{Background: Blue},
			input: "x",
			want:  "\x1b[44mx\x1b[0m",
		},
		"bold only": {
			style: Style{Bold: true},
			input: "x",
			want:  "\x1b[1mx\x1b[0m",
		},
		"combined": {
			style: Style{Foreground: Green, Background: Black, Bold: true},
			input: "x",
			want:  "\x1b[32m\x1b[40m\x1b[1mx\x1b[0m",
		},
		"reset skipped": {
			style: Style{Foreground: Reset, Background: Reset},
			input: "x",
			want:  "x\x1b[0m",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := tc.style.Apply(tc.input)
			if got != tc.want {
				t.Errorf("Apply(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestBox(t *testing.T) {
	b := SingleBorders

	tt := map[string]struct {
		width int
		lines []string
		want  string
	}{
		"auto width": {
			width: 0,
			lines: []string{"ab", "c"},
			want: strings.Join([]string{
				"┌──┐",
				"│ab│",
				"│c │",
				"└──┘",
			}, "\n"),
		},
		"fixed width": {
			width: 5,
			lines: []string{"x"},
			want: strings.Join([]string{
				"┌───┐",
				"│x  │",
				"└───┘",
			}, "\n"),
		},
		"empty lines": {
			width: 4,
			lines: []string{},
			want: strings.Join([]string{
				"┌──┐",
				"└──┘",
			}, "\n"),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := Box(tc.width, tc.lines, b)
			if got != tc.want {
				t.Errorf("Box(%d, %v) = %q, want %q", tc.width, tc.lines, got, tc.want)
			}
		})
	}
}

package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	var (
		bgFlag    = flag.String("bg", "", "custom full-screen background color (hex, e.g. #0d1117)")
		themeFlag = flag.String("theme", "default", "color theme: default")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of speed:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n  speed\n  speed --bg '#0d1117'\n  speed --bg '#f5f5f5'\n")
	}
	flag.Parse()

	theme := DefaultTheme
	hasBg := false
	if *bgFlag != "" {
		// Accept with or without a leading '#'.
		c := *bgFlag
		if c[0] != '#' {
			c = "#" + c
		}
		if !validHex(c) {
			fmt.Fprintf(os.Stderr, "speed: invalid --bg color %q (expected hex like #0d1117)\n", *bgFlag)
			os.Exit(1)
		}
		theme.Background = lipgloss.Color(c)
		hasBg = true
	}
	_ = themeFlag // reserved for future palettes

	m := newModel(theme, hasBg)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "speed: %v\n", err)
		os.Exit(1)
	}
}

// validHex checks that s is a #rrggbb hex color.
func validHex(s string) bool {
	if len(s) != 7 || s[0] != '#' {
		return false
	}
	for _, c := range s[1:] {
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		case c >= 'A' && c <= 'F':
		default:
			return false
		}
	}
	return true
}

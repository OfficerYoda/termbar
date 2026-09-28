// Package termbar renders compact progress bars for terminals and tmux.
package termbar

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	fullBlock  = "█"
	emptyBlock = " "
)

var blocks = [...]string{" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}

// Output selects the formatting understood by the receiving terminal.
type Output string

const (
	Terminal Output = "terminal"
	Tmux     Output = "tmux"
)

// Options configures a progress bar. Its zero value renders a 10-column terminal bar.
type Options struct {
	Size          int
	Output        Output
	ShowPercent   bool
	ShowDelimiter bool
	BoldDelimiter bool
	Delimiter     string
}

// Render returns a progress bar with text aligned at its right edge.
func Render(percent float64, text string, options Options) (string, error) {
	if percent < 0 || percent > 100 {
		return "", fmt.Errorf("percent must be between 0 and 100: %g", percent)
	}

	if options.Size == 0 {
		options.Size = 10
	}
	if options.Size < 1 {
		return "", fmt.Errorf("size must be positive: %d", options.Size)
	}
	if utf8.RuneCountInString(text) > options.Size {
		return "", fmt.Errorf("text does not fit in bar: %q is wider than %d columns", text, options.Size)
	}
	if options.Output == "" {
		options.Output = Terminal
	}
	if options.Output != Terminal && options.Output != Tmux {
		return "", fmt.Errorf("unknown output format: %q", options.Output)
	}
	if options.Delimiter == "" {
		options.Delimiter = " |"
	}

	full := int(percent * float64(options.Size) / 100)
	remainder := percent/100*float64(options.Size) - float64(full)
	textWidth := utf8.RuneCountInString(text)
	textStart := options.Size - textWidth
	var bar strings.Builder

	for i := range options.Size {
		switch {
		case i >= textStart:
			if i == textStart && full > textStart {
				bar.WriteString(styleText(options.Output, text))
			} else if i == textStart {
				bar.WriteString(text)
			}
		case i < full:
			bar.WriteString(fullBlock)
		case i == full && full < options.Size:
			bar.WriteString(blocks[int(remainder*float64(len(blocks)))])
		default:
			bar.WriteString(emptyBlock)
		}
	}

	result := styleBold(options.Output) + styleBar(options.Output, bar.String())
	if options.ShowPercent {
		result += fmt.Sprintf(" %.0f%%", percent)
	}
	if options.ShowDelimiter {
		result += styleDelimiter(options.Output, options.Delimiter, options.BoldDelimiter)
	}
	return result + resetStyle(options.Output), nil
}

func styleBar(output Output, bar string) string {
	if output == Tmux {
		return "[#[bg=color238]" + bar + "#[bg=default,fg=default]]"
	}
	return "[\x1b[48;5;238m" + bar + "\x1b[49;39m]"
}

func styleText(output Output, text string) string {
	if output == Tmux {
		return "#[bg=lightgrey,fg=color237]" + text + "#[bg=color238,fg=lightgrey]"
	}
	return "\x1b[48;5;250m\x1b[38;5;237m" + text + "\x1b[48;5;238m\x1b[38;5;250m"
}

func styleDelimiter(output Output, delimiter string, bold bool) string {
	if output == Tmux {
		if !bold {
			return "#[nobold,fg=colour242]" + delimiter
		}
		return "#[bold,fg=colour242]" + delimiter
	}
	if !bold {
		return "\x1b[22;38;5;242m" + delimiter
	}
	return "\x1b[38;5;242m" + delimiter
}

func styleBold(output Output) string {
	if output == Tmux {
		return "#[bold]"
	}
	return "\x1b[1m"
}

func resetStyle(output Output) string {
	if output == Tmux {
		return "#[default]"
	}
	return "\x1b[0m"
}

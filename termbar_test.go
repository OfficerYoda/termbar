package termbar

import (
	"strings"
	"testing"
)

func TestRenderTerminal(t *testing.T) {
	got, err := Render(45, "12:30", Options{Size: 10})
	if err != nil {
		t.Fatal(err)
	}
	if got != "\x1b[1m[\x1b[48;5;238m████▌12:30\x1b[49;39m]\x1b[22;38;5;242m |\x1b[0m" {
		t.Fatalf("Render() = %q", got)
	}
}

func TestRenderTmuxWithExtras(t *testing.T) {
	got, err := Render(100, "12:30", Options{Size: 10, Output: Tmux, ShowPercent: true, BoldDelimiter: true})
	if err != nil {
		t.Fatal(err)
	}
	want := "#[bold][#[bg=color238]█████#[bg=lightgrey,fg=color237]12:30#[bg=color238,fg=lightgrey]#[bg=default,fg=default]] 100%#[bold,fg=colour242] |#[default]"
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRenderTerminalDelimiterIsGray(t *testing.T) {
	got, err := Render(0, "", Options{Size: 1, BoldDelimiter: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(got, "\x1b[38;5;242m |\x1b[0m") {
		t.Fatalf("Render() = %q", got)
	}
}

func TestRenderCanDisableDelimiterBold(t *testing.T) {
	got, err := Render(0, "", Options{Size: 1, Output: Tmux})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "#[nobold,fg=colour242] |") {
		t.Fatalf("Render() = %q", got)
	}
}

func TestRenderCanHideDelimiter(t *testing.T) {
	got, err := Render(0, "", Options{Size: 1, HideDelimiter: true})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "|") {
		t.Fatalf("Render() = %q", got)
	}
}

func TestRenderRejectsTextThatDoesNotFit(t *testing.T) {
	_, err := Render(50, "too long", Options{Size: 3})
	if err == nil || !strings.Contains(err.Error(), "does not fit") {
		t.Fatalf("Render() error = %v", err)
	}
}

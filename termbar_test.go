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
	if got != "\x1b[48;5;238m████▌12:30\x1b[0m" {
		t.Fatalf("Render() = %q", got)
	}
}

func TestRenderTmuxWithExtras(t *testing.T) {
	got, err := Render(100, "12:30", Options{Size: 10, Output: Tmux, ShowPercent: true, ShowDelimiter: true})
	if err != nil {
		t.Fatal(err)
	}
	want := "#[bg=color238]█████#[bg=lightgrey,fg=color237]12:30#[bg=color238,fg=lightgrey]#[bg=default]#[nobold] 100% |"
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRenderRejectsTextThatDoesNotFit(t *testing.T) {
	_, err := Render(50, "too long", Options{Size: 3})
	if err == nil || !strings.Contains(err.Error(), "does not fit") {
		t.Fatalf("Render() error = %v", err)
	}
}

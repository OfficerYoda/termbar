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

func TestRenderMovesTextThroughBar(t *testing.T) {
	tests := []struct {
		percent float64
		bar     string
	}{
		{75, "███████████████11:45"},
		{80, "███████████11:45    "},
		{82.5, "███████████11:45▌   "},
		{100, "███████████████11:45"},
	}

	for _, test := range tests {
		got, err := Render(test.percent, "11:45", Options{Size: 20, Output: Tmux, HideDelimiter: true})
		if err != nil {
			t.Fatal(err)
		}
		styledText := "#[bg=lightgrey,fg=color237]11:45#[bg=color238,fg=lightgrey]"
		wantBar := strings.Replace(test.bar, "11:45", styledText, 1)
		if test.percent == 75 {
			wantBar = test.bar
		}
		want := "#[push-default]#[bold][#[bg=color238]" + wantBar + "#[bg=default,fg=default]]#[default]#[pop-default]"
		if got != want {
			t.Errorf("Render(%g) = %q, want %q", test.percent, got, want)
		}
	}
}

func TestRenderTmuxWithExtras(t *testing.T) {
	got, err := Render(100, "12:30", Options{Size: 10, Output: Tmux, ShowPercent: true, BoldDelimiter: true})
	if err != nil {
		t.Fatal(err)
	}
	want := "#[push-default]#[bold][#[bg=color238]█████#[bg=lightgrey,fg=color237]12:30#[bg=color238,fg=lightgrey]#[bg=default,fg=default]] 100%#[bold,fg=colour242] |#[default]#[pop-default]"
	if got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

func TestRenderTmuxPreservesSurroundingStyle(t *testing.T) {
	got, err := Render(0, "", Options{Size: 1, Output: Tmux})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "#[push-default]") || !strings.HasSuffix(got, "#[default]#[pop-default]") {
		t.Fatalf("Render() = %q", got)
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

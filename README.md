# termbar

`termbar` renders compact progress bars for ANSI terminals and tmux.

```go
bar, err := termbar.Render(42.5, "12:30", termbar.Options{
	Size:          10,
	Output:        termbar.Terminal,
	ShowPercent:   true,
	ShowDelimiter: true,
})
```

`Render` returns an error when the percentage is outside `0..100`, the output format is unknown, or the text is wider than the bar. Text width is counted in Unicode runes; wide characters such as emoji may occupy more terminal columns.

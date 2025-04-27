package handlers

import (
	"fmt"
	"html"
)

type SSELogger struct {
	broadcaster *Broadcaster
}

func NewSSELogger(b *Broadcaster) *SSELogger {
	return &SSELogger{broadcaster: b}
}

func (s *SSELogger) Write(p []byte) (n int, err error) {
	text := string(p)
	// Trim
	if len(text) > 0 && text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
	}
	// Escape
	escaped := html.EscapeString(text)
	htmlMsg := fmt.Sprintf("<div class=\"log-entry\">%s</div>", escaped)
	s.broadcaster.Broadcast("log", htmlMsg)
	return len(p), nil
}

package handlers

import (
	"fmt"
	"html"
)

// SSELogger is a writer that broadcasts log messages as Server-Sent Events (SSE).
//
// Deprecated: We stopped using this and opted to use hx-trigger=every 5s instead
// We thought this might conflict with the main sse stream connection
type SSELogger struct {
	broadcaster *Broadcaster
}

// NewSSELogger creates a new SSELogger that uses the provided Broadcaster.
//
// Deprecated
func NewSSELogger(b *Broadcaster) *SSELogger {
	return &SSELogger{broadcaster: b}
}

// Write implements the io.Writer interface.
// It formats the log message by trimming and escaping HTML,
// and broadcasts it as an SSE event with the type "log" using the Broadcaster.
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

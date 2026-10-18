// Package progress draws a one-line spinner with a status message on a terminal.
package progress

import (
	"fmt"
	"io"
	"sync"
	"time"
)

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner redraws "<frame> <message>" on w until stopped. A nil *Spinner is a no-op,
// which callers use when output is not a terminal.
type Spinner struct {
	w    io.Writer
	mu   sync.Mutex
	msg  string
	stop chan struct{}
	done chan struct{}
}

// Start begins drawing on w.
func Start(w io.Writer, msg string) *Spinner {
	s := &Spinner{w: w, msg: msg, stop: make(chan struct{}), done: make(chan struct{})}
	go s.loop()
	return s
}

// Set replaces the status message.
func (s *Spinner) Set(format string, args ...interface{}) {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.msg = fmt.Sprintf(format, args...)
	s.mu.Unlock()
}

// Stop clears the line and stops drawing.
func (s *Spinner) Stop() {
	if s == nil {
		return
	}
	close(s.stop)
	<-s.done
}

func (s *Spinner) loop() {
	defer close(s.done)
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()
	for i := 0; ; i++ {
		s.mu.Lock()
		msg := s.msg
		s.mu.Unlock()
		fmt.Fprintf(s.w, "\r\033[K%s %s", frames[i%len(frames)], msg)
		select {
		case <-s.stop:
			fmt.Fprint(s.w, "\r\033[K")
			return
		case <-t.C:
		}
	}
}

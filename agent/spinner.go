package agent

import (
	"fmt"
	"strings"
	"time"
)

type Spinner struct {
	frames []string
	idx    int
	stop   chan bool
	done   chan bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		stop:   make(chan bool),
		done:   make(chan bool),
	}
}

func (s *Spinner) Start(message string) {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("\r%s %s", s.frames[s.idx], message)
				s.idx = (s.idx + 1) % len(s.frames)
			case <-s.stop:
				fmt.Printf("\r%s\r", strings.Repeat(" ", len(message)+2))
				s.done <- true
				return
			}
		}
	}()
}

func (s *Spinner) Stop() {
	close(s.stop)
	<-s.done
}

package timer

import (
	"fmt"
	"time"
)

// Stopwatch is a standalone count-up timer that runs independently
type Stopwatch struct {
	Elapsed time.Duration
	Running bool
}

// NewStopwatch creates a new stopwatch
func NewStopwatch() *Stopwatch {
	return &Stopwatch{}
}

// Start begins the stopwatch
func (s *Stopwatch) Start() {
	s.Running = true
}

// Pause stops the stopwatch
func (s *Stopwatch) Pause() {
	s.Running = false
}

// Toggle switches between running and paused
func (s *Stopwatch) Toggle() {
	s.Running = !s.Running
}

// Reset resets the stopwatch to zero
func (s *Stopwatch) Reset() {
	s.Elapsed = 0
	s.Running = false
}

// Tick advances the stopwatch by one second
func (s *Stopwatch) Tick() {
	if s.Running {
		s.Elapsed += time.Second
	}
}

// Format returns the elapsed time as HH:MM:SS or MM:SS
func (s *Stopwatch) Format() string {
	total := int(s.Elapsed.Seconds())
	hours := total / 3600
	mins := (total % 3600) / 60
	secs := total % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, mins, secs)
	}
	return fmt.Sprintf("%02d:%02d", mins, secs)
}

// FormatShort returns a compact format for the indicator
func (s *Stopwatch) FormatShort() string {
	return s.Format()
}

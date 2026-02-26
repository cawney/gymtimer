package timer

import "time"

// AMRAPTimer handles AMRAP (As Many Rounds As Possible) countdown logic
type AMRAPTimer struct {
	*Timer
}

// NewAMRAP creates a new AMRAP timer with the given duration
func NewAMRAP(duration time.Duration) *AMRAPTimer {
	t := New()
	t.Mode = ModeAMRAP
	t.Duration = duration
	return &AMRAPTimer{Timer: t}
}

// Tick advances the AMRAP timer
func (a *AMRAPTimer) Tick() {
	if !a.Running {
		return
	}

	a.Elapsed += time.Second

	// Countdown beeps at specific intervals
	remaining := a.Duration - a.Elapsed

	// Final 3-2-1 countdown
	if remaining <= 3*time.Second && remaining > 0 {
		secs := int(remaining.Seconds())
		if a.OnCountdownTick != nil {
			a.OnCountdownTick(secs)
		}
	}

	// Beep at minute marks in the last minute
	if remaining == time.Minute {
		if a.OnCountdownTick != nil {
			a.OnCountdownTick(60)
		}
	}

	// Check if finished
	if a.Elapsed >= a.Duration {
		a.Running = false
		if a.OnIntervalChange != nil {
			a.OnIntervalChange(PhaseRest) // Signal completion
		}
	}
}

// Remaining returns the time remaining
func (a *AMRAPTimer) Remaining() time.Duration {
	remaining := a.Duration - a.Elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Progress returns completion percentage (0-100)
func (a *AMRAPTimer) Progress() float64 {
	if a.Duration == 0 {
		return 0
	}
	return float64(a.Elapsed) / float64(a.Duration) * 100
}

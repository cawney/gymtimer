package timer

import "time"

// CustomTimer handles custom interval logic with user-defined work/rest/rounds
type CustomTimer struct {
	*Timer
}

// NewCustom creates a new custom interval timer
func NewCustom(work, rest time.Duration, rounds int) *CustomTimer {
	t := New()
	t.Mode = ModeCustom
	t.WorkDuration = work
	t.RestDuration = rest
	t.TotalRounds = rounds
	t.Phase = PhaseWork
	return &CustomTimer{Timer: t}
}

// Tick advances the custom timer (same logic as Tabata)
func (c *CustomTimer) Tick() {
	if !c.Running {
		return
	}

	c.Elapsed += time.Second

	var currentDuration time.Duration
	if c.Phase == PhaseWork {
		currentDuration = c.WorkDuration
	} else {
		currentDuration = c.RestDuration
	}

	// 3-2-1 countdown beeps
	remaining := currentDuration - c.Elapsed
	if remaining <= 3*time.Second && remaining > 0 {
		secs := int(remaining.Seconds())
		if c.OnCountdownTick != nil {
			c.OnCountdownTick(secs)
		}
	}

	// Check if interval is complete
	if c.Elapsed >= currentDuration {
		c.Elapsed = 0

		if c.Phase == PhaseWork {
			c.Phase = PhaseRest
		} else {
			c.Phase = PhaseWork
			c.Round++
			if c.OnRoundChange != nil {
				c.OnRoundChange(c.Round)
			}
		}

		if c.OnIntervalChange != nil {
			c.OnIntervalChange(c.Phase)
		}
	}
}

// CurrentIntervalRemaining returns time remaining in current work/rest interval
func (c *CustomTimer) CurrentIntervalRemaining() time.Duration {
	var currentDuration time.Duration
	if c.Phase == PhaseWork {
		currentDuration = c.WorkDuration
	} else {
		currentDuration = c.RestDuration
	}
	remaining := currentDuration - c.Elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// TotalWorkoutDuration calculates total workout time
func (c *CustomTimer) TotalWorkoutDuration() time.Duration {
	return time.Duration(c.TotalRounds) * (c.WorkDuration + c.RestDuration)
}

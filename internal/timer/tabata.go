package timer

import "time"

// TabataTimer handles Tabata interval logic
type TabataTimer struct {
	*Timer
}

// NewTabata creates a new Tabata timer with default 20s work / 10s rest / 8 rounds
func NewTabata() *TabataTimer {
	t := New()
	t.Mode = ModeTabata
	t.WorkDuration = 20 * time.Second
	t.RestDuration = 10 * time.Second
	t.TotalRounds = 8
	t.Phase = PhaseWork
	return &TabataTimer{Timer: t}
}

// Tick advances the Tabata timer
func (tb *TabataTimer) Tick() {
	if !tb.Running {
		return
	}

	tb.Elapsed += time.Second

	var currentDuration time.Duration
	if tb.Phase == PhaseWork {
		currentDuration = tb.WorkDuration
	} else {
		currentDuration = tb.RestDuration
	}

	// 3-2-1 countdown beeps
	remaining := currentDuration - tb.Elapsed
	if remaining <= 3*time.Second && remaining > 0 {
		secs := int(remaining.Seconds())
		if tb.OnCountdownTick != nil {
			tb.OnCountdownTick(secs)
		}
	}

	// Check if interval is complete
	if tb.Elapsed >= currentDuration {
		tb.Elapsed = 0

		if tb.Phase == PhaseWork {
			tb.Phase = PhaseRest
		} else {
			tb.Phase = PhaseWork
			tb.Round++
			if tb.OnRoundChange != nil {
				tb.OnRoundChange(tb.Round)
			}
		}

		if tb.OnIntervalChange != nil {
			tb.OnIntervalChange(tb.Phase)
		}
	}
}

// CurrentIntervalRemaining returns time remaining in current work/rest interval
func (tb *TabataTimer) CurrentIntervalRemaining() time.Duration {
	var currentDuration time.Duration
	if tb.Phase == PhaseWork {
		currentDuration = tb.WorkDuration
	} else {
		currentDuration = tb.RestDuration
	}
	remaining := currentDuration - tb.Elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

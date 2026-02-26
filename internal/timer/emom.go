package timer

import "time"

// EMOMTimer handles EMOM (Every Minute On the Minute) logic
type EMOMTimer struct {
	*Timer
}

// NewEMOM creates a new EMOM timer
func NewEMOM(rounds int) *EMOMTimer {
	t := New()
	t.Mode = ModeEMOM
	t.TotalRounds = rounds
	return &EMOMTimer{Timer: t}
}

// Tick advances the EMOM timer
func (e *EMOMTimer) Tick() {
	if !e.Running {
		return
	}

	e.Elapsed += time.Second

	// Check if minute is complete
	if e.Elapsed >= time.Minute {
		e.Elapsed = 0
		e.Round++

		if e.OnRoundChange != nil {
			e.OnRoundChange(e.Round)
		}

		if e.OnIntervalChange != nil {
			e.OnIntervalChange(PhaseWork)
		}
	}

	// 3-2-1 countdown beeps
	remaining := time.Minute - e.Elapsed
	if remaining <= 3*time.Second && remaining > 0 {
		secs := int(remaining.Seconds())
		if e.OnCountdownTick != nil {
			e.OnCountdownTick(secs)
		}
	}
}

// SecondsInMinute returns elapsed seconds in current minute
func (e *EMOMTimer) SecondsInMinute() int {
	return int(e.Elapsed.Seconds()) % 60
}

// SecondsRemaining returns seconds remaining in current minute
func (e *EMOMTimer) SecondsRemaining() int {
	return 60 - e.SecondsInMinute()
}

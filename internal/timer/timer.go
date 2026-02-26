package timer

import "time"

// Mode represents the timer mode
type Mode int

const (
	ModeClock Mode = iota
	ModeEMOM
	ModeTabata
	ModeAMRAP
	ModeCustom
	ModeStopwatch
)

// Phase represents work or rest phase
type Phase int

const (
	PhaseWork Phase = iota
	PhaseRest
	PhaseCountdown
)

// Timer holds the core timer state
type Timer struct {
	Duration    time.Duration // Total duration for countdown modes
	Elapsed     time.Duration // Time elapsed in current interval
	Running     bool
	Mode        Mode
	Phase       Phase
	Round       int
	TotalRounds int

	// Interval settings
	WorkDuration time.Duration
	RestDuration time.Duration

	// Countdown before start
	CountdownRemaining int

	// Callbacks
	OnIntervalChange func(phase Phase)
	OnCountdownTick  func(remaining int)
	OnRoundChange    func(round int)
}

// New creates a new timer with default settings
func New() *Timer {
	return &Timer{
		Duration:     20 * time.Minute,
		WorkDuration: 20 * time.Second,
		RestDuration: 10 * time.Second,
		TotalRounds:  8,
		Round:        1,
		Phase:        PhaseWork,
	}
}

// Start begins the timer
func (t *Timer) Start() {
	t.Running = true
}

// Pause stops the timer
func (t *Timer) Pause() {
	t.Running = false
}

// Toggle switches between running and paused
func (t *Timer) Toggle() {
	t.Running = !t.Running
}

// Reset resets the timer to initial state for current mode
func (t *Timer) Reset() {
	t.Elapsed = 0
	t.Round = 1
	t.Phase = PhaseWork
	t.Running = false
	t.CountdownRemaining = 0
}

// Tick advances the timer by one second
func (t *Timer) Tick() {
	if !t.Running {
		return
	}
	t.Elapsed += time.Second
}

// TimeRemaining returns remaining time for countdown modes
func (t *Timer) TimeRemaining() time.Duration {
	switch t.Mode {
	case ModeAMRAP:
		remaining := t.Duration - t.Elapsed
		if remaining < 0 {
			return 0
		}
		return remaining
	case ModeTabata, ModeCustom:
		var intervalDuration time.Duration
		if t.Phase == PhaseWork {
			intervalDuration = t.WorkDuration
		} else {
			intervalDuration = t.RestDuration
		}
		remaining := intervalDuration - t.Elapsed
		if remaining < 0 {
			return 0
		}
		return remaining
	case ModeEMOM:
		// EMOM counts up within each minute
		remaining := time.Minute - t.Elapsed
		if remaining < 0 {
			return 0
		}
		return remaining
	default:
		return 0
	}
}

// ElapsedInInterval returns elapsed time in current interval
func (t *Timer) ElapsedInInterval() time.Duration {
	return t.Elapsed
}

// IsFinished returns true if the timer has completed
func (t *Timer) IsFinished() bool {
	switch t.Mode {
	case ModeAMRAP:
		return t.Elapsed >= t.Duration
	case ModeTabata, ModeCustom:
		return t.Round > t.TotalRounds
	case ModeEMOM:
		return t.Round > t.TotalRounds
	default:
		return false
	}
}

// SetMode changes the timer mode and resets
func (t *Timer) SetMode(mode Mode) {
	t.Mode = mode
	t.Reset()

	// Set default values for each mode
	switch mode {
	case ModeTabata:
		t.WorkDuration = 20 * time.Second
		t.RestDuration = 10 * time.Second
		t.TotalRounds = 8
	case ModeEMOM:
		t.TotalRounds = 10
	case ModeAMRAP:
		t.Duration = 20 * time.Minute
	case ModeCustom:
		t.WorkDuration = 30 * time.Second
		t.RestDuration = 15 * time.Second
		t.TotalRounds = 5
	}
}

// ModeName returns the string name of the current mode
func (t *Timer) ModeName() string {
	switch t.Mode {
	case ModeClock:
		return "CLOCK"
	case ModeEMOM:
		return "EMOM"
	case ModeTabata:
		return "TABATA"
	case ModeAMRAP:
		return "AMRAP"
	case ModeCustom:
		return "CUSTOM"
	case ModeStopwatch:
		return "STOPWATCH"
	default:
		return "UNKNOWN"
	}
}

// PhaseName returns the string name of the current phase
func (t *Timer) PhaseName() string {
	switch t.Phase {
	case PhaseWork:
		return "WORK"
	case PhaseRest:
		return "REST"
	case PhaseCountdown:
		return "GET READY"
	default:
		return ""
	}
}

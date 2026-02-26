package ui

import (
	"fmt"
	"time"

	"gymtimer/internal/audio"
	"gymtimer/internal/timer"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppState represents the current application state
type AppState int

const (
	StateRunning AppState = iota
	StateSetup
	StatePaused
	StateFinished
)

// SettingField represents which setting is being edited
type SettingField int

const (
	SettingWork SettingField = iota
	SettingRest
	SettingRounds
	SettingDuration
)

// Model is the main Bubbletea model
type Model struct {
	timer       *timer.Timer
	stopwatch   *timer.Stopwatch
	audio       *audio.Player
	keys        KeyMap
	width       int
	height      int
	state       AppState
	settingField SettingField

	// Track last countdown beep to avoid duplicates
	lastCountdownBeep int
}

// TickMsg is sent every second
type TickMsg time.Time

// New creates a new app model
func New(audioPlayer *audio.Player) Model {
	t := timer.New()
	t.Mode = timer.ModeClock

	return Model{
		timer:        t,
		stopwatch:    timer.NewStopwatch(),
		audio:        audioPlayer,
		keys:         DefaultKeyMap(),
		state:        StateRunning,
		settingField: SettingWork,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), tea.EnterAltScreen)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case TickMsg:
		m.handleTick()
		return m, tickCmd()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) handleTick() {
	// Always tick the stopwatch (it runs independently)
	m.stopwatch.Tick()

	// Handle stopwatch mode separately
	if m.timer.Mode == timer.ModeStopwatch {
		return
	}

	if m.timer.Mode == timer.ModeClock {
		return
	}

	if !m.timer.Running {
		return
	}

	oldPhase := m.timer.Phase
	oldRound := m.timer.Round

	// Track time remaining before tick for countdown beeps
	var remaining time.Duration
	switch m.timer.Mode {
	case timer.ModeEMOM:
		remaining = time.Minute - m.timer.Elapsed
	case timer.ModeTabata, timer.ModeCustom:
		if m.timer.Phase == timer.PhaseWork {
			remaining = m.timer.WorkDuration - m.timer.Elapsed
		} else {
			remaining = m.timer.RestDuration - m.timer.Elapsed
		}
	case timer.ModeAMRAP:
		remaining = m.timer.Duration - m.timer.Elapsed
	}

	// Handle countdown beeps (3, 2, 1)
	secs := int(remaining.Seconds())
	if secs <= 3 && secs > 0 && secs != m.lastCountdownBeep {
		m.audio.PlayCountdown(secs)
		m.lastCountdownBeep = secs
	}
	if secs > 3 {
		m.lastCountdownBeep = 0
	}

	// Advance timer
	m.timer.Tick()

	// Handle interval transitions
	switch m.timer.Mode {
	case timer.ModeEMOM:
		if m.timer.Elapsed >= time.Minute {
			m.timer.Elapsed = 0
			m.timer.Round++
			m.audio.PlayIntervalChange(true)
			if m.timer.Round > m.timer.TotalRounds {
				m.timer.Running = false
				m.state = StateFinished
				m.audio.PlayFinish()
			}
		}
	case timer.ModeTabata, timer.ModeCustom:
		var currentDuration time.Duration
		if m.timer.Phase == timer.PhaseWork {
			currentDuration = m.timer.WorkDuration
		} else {
			currentDuration = m.timer.RestDuration
		}

		if m.timer.Elapsed >= currentDuration {
			m.timer.Elapsed = 0
			if m.timer.Phase == timer.PhaseWork {
				m.timer.Phase = timer.PhaseRest
				m.audio.PlayIntervalChange(false)
			} else {
				m.timer.Phase = timer.PhaseWork
				m.timer.Round++
				m.audio.PlayIntervalChange(true)
			}
			if m.timer.Round > m.timer.TotalRounds {
				m.timer.Running = false
				m.state = StateFinished
				m.audio.PlayFinish()
			}
		}
	case timer.ModeAMRAP:
		if m.timer.Elapsed >= m.timer.Duration {
			m.timer.Running = false
			m.state = StateFinished
			m.audio.PlayFinish()
		}
	}

	// Play sound on phase/round change
	if oldPhase != m.timer.Phase || oldRound != m.timer.Round {
		// Already handled above
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Quit always works
	if m.keys.Quit.Matches(msg) {
		return m, tea.Quit
	}

	// Handle setup mode keys
	if m.state == StateSetup {
		return m.handleSetupKey(msg)
	}

	// Mode switching
	if m.keys.ModeClock.Matches(msg) {
		m.timer.SetMode(timer.ModeClock)
		m.state = StateRunning
		return m, nil
	}
	if m.keys.ModeEMOM.Matches(msg) {
		m.timer.SetMode(timer.ModeEMOM)
		m.state = StateSetup
		m.settingField = SettingRounds
		return m, nil
	}
	if m.keys.ModeTabata.Matches(msg) {
		m.timer.SetMode(timer.ModeTabata)
		m.state = StateSetup
		m.settingField = SettingWork
		return m, nil
	}
	if m.keys.ModeAMRAP.Matches(msg) {
		m.timer.SetMode(timer.ModeAMRAP)
		m.state = StateSetup
		m.settingField = SettingDuration
		return m, nil
	}
	if m.keys.ModeCustom.Matches(msg) {
		m.timer.SetMode(timer.ModeCustom)
		m.state = StateSetup
		m.settingField = SettingWork
		return m, nil
	}
	if m.keys.ModeStopwatch.Matches(msg) {
		m.timer.SetMode(timer.ModeStopwatch)
		m.state = StateRunning
		return m, nil
	}

	// Stopwatch controls (work from any mode)
	if m.keys.StopwatchToggle.Matches(msg) {
		m.stopwatch.Toggle()
		return m, nil
	}
	if m.keys.StopwatchReset.Matches(msg) {
		m.stopwatch.Reset()
		return m, nil
	}

	// Start/pause
	if m.keys.StartPause.Matches(msg) {
		// In stopwatch mode, space controls the stopwatch
		if m.timer.Mode == timer.ModeStopwatch {
			m.stopwatch.Toggle()
			return m, nil
		}
		if m.state == StateFinished {
			m.timer.Reset()
			m.state = StateRunning
		}
		m.timer.Toggle()
		if m.timer.Running {
			m.state = StateRunning
		} else {
			m.state = StatePaused
		}
		return m, nil
	}

	// Reset
	if m.keys.Reset.Matches(msg) {
		// In stopwatch mode, R resets the stopwatch
		if m.timer.Mode == timer.ModeStopwatch {
			m.stopwatch.Reset()
			return m, nil
		}
		m.timer.Reset()
		m.state = StateRunning
		m.lastCountdownBeep = 0
		return m, nil
	}

	// Toggle sound
	if m.keys.ToggleSound.Matches(msg) {
		m.audio.SetEnabled(!m.audio.IsEnabled())
		return m, nil
	}

	return m, nil
}

func (m Model) handleSetupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case m.keys.Enter.Matches(msg):
		m.state = StateRunning
		return m, nil

	case m.keys.Up.Matches(msg):
		m.adjustSetting(1)
		return m, nil

	case m.keys.Down.Matches(msg):
		m.adjustSetting(-1)
		return m, nil

	case m.keys.Left.Matches(msg):
		m.prevSetting()
		return m, nil

	case m.keys.Right.Matches(msg):
		m.nextSetting()
		return m, nil
	}

	return m, nil
}

func (m *Model) adjustSetting(delta int) {
	switch m.settingField {
	case SettingWork:
		m.timer.WorkDuration += time.Duration(delta*5) * time.Second
		if m.timer.WorkDuration < 5*time.Second {
			m.timer.WorkDuration = 5 * time.Second
		}
		if m.timer.WorkDuration > 5*time.Minute {
			m.timer.WorkDuration = 5 * time.Minute
		}
	case SettingRest:
		m.timer.RestDuration += time.Duration(delta*5) * time.Second
		if m.timer.RestDuration < 5*time.Second {
			m.timer.RestDuration = 5 * time.Second
		}
		if m.timer.RestDuration > 5*time.Minute {
			m.timer.RestDuration = 5 * time.Minute
		}
	case SettingRounds:
		m.timer.TotalRounds += delta
		if m.timer.TotalRounds < 1 {
			m.timer.TotalRounds = 1
		}
		if m.timer.TotalRounds > 99 {
			m.timer.TotalRounds = 99
		}
	case SettingDuration:
		m.timer.Duration += time.Duration(delta) * time.Minute
		if m.timer.Duration < time.Minute {
			m.timer.Duration = time.Minute
		}
		if m.timer.Duration > 60*time.Minute {
			m.timer.Duration = 60 * time.Minute
		}
	}
}

func (m *Model) nextSetting() {
	switch m.timer.Mode {
	case timer.ModeTabata, timer.ModeCustom:
		switch m.settingField {
		case SettingWork:
			m.settingField = SettingRest
		case SettingRest:
			m.settingField = SettingRounds
		case SettingRounds:
			m.settingField = SettingWork
		}
	}
}

func (m *Model) prevSetting() {
	switch m.timer.Mode {
	case timer.ModeTabata, timer.ModeCustom:
		switch m.settingField {
		case SettingWork:
			m.settingField = SettingRounds
		case SettingRest:
			m.settingField = SettingWork
		case SettingRounds:
			m.settingField = SettingRest
		}
	}
}

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string

	switch m.state {
	case StateSetup:
		content = m.renderSetup()
	default:
		content = m.renderTimer()
	}

	return CenterInScreen(content, m.width, m.height)
}

func (m Model) renderTimer() string {
	var s string

	// Mode title
	title := TitleStyle.Render(fmt.Sprintf("MODE: %s", m.timer.ModeName()))
	s += title + "\n\n"

	// Time display
	var timeStr string
	var color lipgloss.Color

	switch m.timer.Mode {
	case timer.ModeClock:
		now := time.Now()
		timeStr = now.Format("15:04:05")
		color = ColorNeutral
	case timer.ModeEMOM:
		remaining := time.Minute - m.timer.Elapsed
		if remaining < 0 {
			remaining = 0
		}
		mins := int(remaining.Minutes())
		secs := int(remaining.Seconds()) % 60
		timeStr = fmt.Sprintf("%02d:%02d", mins, secs)
		color = ColorWork
	case timer.ModeTabata, timer.ModeCustom:
		var currentDuration time.Duration
		if m.timer.Phase == timer.PhaseWork {
			currentDuration = m.timer.WorkDuration
			color = ColorWork
		} else {
			currentDuration = m.timer.RestDuration
			color = ColorRest
		}
		remaining := currentDuration - m.timer.Elapsed
		if remaining < 0 {
			remaining = 0
		}
		mins := int(remaining.Minutes())
		secs := int(remaining.Seconds()) % 60
		timeStr = fmt.Sprintf("%02d:%02d", mins, secs)
	case timer.ModeAMRAP:
		remaining := m.timer.Duration - m.timer.Elapsed
		if remaining < 0 {
			remaining = 0
		}
		mins := int(remaining.Minutes())
		secs := int(remaining.Seconds()) % 60
		timeStr = fmt.Sprintf("%02d:%02d", mins, secs)
		color = ColorWork
	case timer.ModeStopwatch:
		timeStr = m.stopwatch.Format()
		if m.stopwatch.Running {
			color = ColorWork
		} else {
			color = ColorPaused
		}
	}

	if m.state == StatePaused && m.timer.Mode != timer.ModeStopwatch {
		color = ColorPaused
	}
	if m.state == StateFinished {
		color = ColorFinished
	}

	s += RenderBigTime(timeStr, color)

	// Phase indicator
	if m.timer.Mode == timer.ModeTabata || m.timer.Mode == timer.ModeCustom {
		if m.timer.Phase == timer.PhaseWork {
			s += PhaseWorkStyle.Render("WORK") + "\n"
		} else {
			s += PhaseRestStyle.Render("REST") + "\n"
		}
	}

	// Round counter
	if m.timer.Mode != timer.ModeClock && m.timer.Mode != timer.ModeAMRAP && m.timer.Mode != timer.ModeStopwatch {
		roundStr := fmt.Sprintf("Round %d of %d", m.timer.Round, m.timer.TotalRounds)
		s += RoundStyle.Render(roundStr) + "\n"
	}

	// Stopwatch indicator (when running in background)
	if m.timer.Mode != timer.ModeStopwatch && (m.stopwatch.Running || m.stopwatch.Elapsed > 0) {
		swStatus := "paused"
		swColor := ColorPaused
		if m.stopwatch.Running {
			swStatus = "running"
			swColor = ColorWork
		}
		swIndicator := fmt.Sprintf("SW: %s (%s)", m.stopwatch.Format(), swStatus)
		s += "\n" + lipgloss.NewStyle().Foreground(swColor).Render(swIndicator)
	}

	// Status
	if m.state == StatePaused && m.timer.Mode != timer.ModeStopwatch {
		s += "\n" + lipgloss.NewStyle().Foreground(ColorPaused).Render("PAUSED")
	}
	if m.state == StateFinished && m.timer.Mode != timer.ModeStopwatch {
		s += "\n" + lipgloss.NewStyle().Foreground(ColorFinished).Bold(true).Render("FINISHED!")
	}
	// Stopwatch status when viewing stopwatch
	if m.timer.Mode == timer.ModeStopwatch && !m.stopwatch.Running && m.stopwatch.Elapsed > 0 {
		s += "\n" + lipgloss.NewStyle().Foreground(ColorPaused).Render("PAUSED")
	}

	// Mode selector
	modes := "[1]Clock  [2]EMOM  [3]Tabata  [4]AMRAP  [5]Custom  [6]Stopwatch"
	s += "\n" + HelpStyle.Render(modes)

	// Help bar
	soundStatus := "ON"
	if !m.audio.IsEnabled() {
		soundStatus = "OFF"
	}
	var help string
	if m.timer.Mode == timer.ModeStopwatch {
		help = fmt.Sprintf("[Space] Start/Pause  [R] Reset  [S] Sound: %s  [Q] Quit", soundStatus)
	} else {
		help = fmt.Sprintf("[Space] Start/Pause  [R] Reset  [W] Stopwatch  [S] Sound: %s  [Q] Quit", soundStatus)
	}
	s += "\n" + HelpStyle.Render(help)

	return s
}

func (m Model) renderSetup() string {
	var s string

	title := TitleStyle.Render(fmt.Sprintf("SETUP: %s", m.timer.ModeName()))
	s += title + "\n\n"

	switch m.timer.Mode {
	case timer.ModeEMOM:
		roundsStyle := SettingStyle
		if m.settingField == SettingRounds {
			roundsStyle = SettingSelectedStyle
		}
		s += roundsStyle.Render(fmt.Sprintf("Rounds: %d", m.timer.TotalRounds)) + "\n"

	case timer.ModeTabata, timer.ModeCustom:
		workStyle := SettingStyle
		restStyle := SettingStyle
		roundsStyle := SettingStyle

		if m.settingField == SettingWork {
			workStyle = SettingSelectedStyle
		}
		if m.settingField == SettingRest {
			restStyle = SettingSelectedStyle
		}
		if m.settingField == SettingRounds {
			roundsStyle = SettingSelectedStyle
		}

		workSecs := int(m.timer.WorkDuration.Seconds())
		restSecs := int(m.timer.RestDuration.Seconds())

		s += workStyle.Render(fmt.Sprintf("Work: %ds", workSecs)) + "\n"
		s += restStyle.Render(fmt.Sprintf("Rest: %ds", restSecs)) + "\n"
		s += roundsStyle.Render(fmt.Sprintf("Rounds: %d", m.timer.TotalRounds)) + "\n"

	case timer.ModeAMRAP:
		durStyle := SettingStyle
		if m.settingField == SettingDuration {
			durStyle = SettingSelectedStyle
		}
		mins := int(m.timer.Duration.Minutes())
		s += durStyle.Render(fmt.Sprintf("Duration: %d min", mins)) + "\n"
	}

	s += "\n"
	help := "[Up/Down] Adjust  [Left/Right] Switch  [Enter] Start  [Q] Quit"
	s += HelpStyle.Render(help)

	return s
}

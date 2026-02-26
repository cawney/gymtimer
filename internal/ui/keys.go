package ui

import "github.com/charmbracelet/bubbletea"

// Key represents a key binding
type Key struct {
	Keys []string
	Help string
}

// KeyMap contains all key bindings
type KeyMap struct {
	Quit            Key
	StartPause      Key
	Reset           Key
	ModeClock       Key
	ModeEMOM        Key
	ModeTabata      Key
	ModeAMRAP       Key
	ModeCustom      Key
	ModeStopwatch   Key
	StopwatchToggle Key
	StopwatchReset  Key
	Up              Key
	Down            Key
	Left            Key
	Right           Key
	Enter           Key
	ToggleSound     Key
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: Key{
			Keys: []string{"q", "ctrl+c"},
			Help: "[Q] Quit",
		},
		StartPause: Key{
			Keys: []string{" "},
			Help: "[Space] Start/Pause",
		},
		Reset: Key{
			Keys: []string{"r"},
			Help: "[R] Reset",
		},
		ModeClock: Key{
			Keys: []string{"1"},
			Help: "[1] Clock",
		},
		ModeEMOM: Key{
			Keys: []string{"2"},
			Help: "[2] EMOM",
		},
		ModeTabata: Key{
			Keys: []string{"3"},
			Help: "[3] Tabata",
		},
		ModeAMRAP: Key{
			Keys: []string{"4"},
			Help: "[4] AMRAP",
		},
		ModeCustom: Key{
			Keys: []string{"5"},
			Help: "[5] Custom",
		},
		ModeStopwatch: Key{
			Keys: []string{"6"},
			Help: "[6] Stopwatch",
		},
		StopwatchToggle: Key{
			Keys: []string{"w"},
			Help: "[W] Stopwatch Start/Stop",
		},
		StopwatchReset: Key{
			Keys: []string{"x"},
			Help: "[X] Stopwatch Reset",
		},
		Up: Key{
			Keys: []string{"up", "k"},
			Help: "[Up] Increase",
		},
		Down: Key{
			Keys: []string{"down", "j"},
			Help: "[Down] Decrease",
		},
		Left: Key{
			Keys: []string{"left", "h"},
			Help: "[Left] Previous",
		},
		Right: Key{
			Keys: []string{"right", "l"},
			Help: "[Right] Next",
		},
		Enter: Key{
			Keys: []string{"enter"},
			Help: "[Enter] Confirm",
		},
		ToggleSound: Key{
			Keys: []string{"s"},
			Help: "[S] Sound",
		},
	}
}

// Matches checks if the key message matches this key binding
func (k Key) Matches(msg tea.KeyMsg) bool {
	for _, key := range k.Keys {
		if msg.String() == key {
			return true
		}
	}
	return false
}

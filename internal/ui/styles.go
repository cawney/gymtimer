package ui

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	ColorWork     = lipgloss.Color("#00FF00") // Green for work
	ColorRest     = lipgloss.Color("#FF6600") // Orange for rest
	ColorPaused   = lipgloss.Color("#FFFF00") // Yellow for paused
	ColorFinished = lipgloss.Color("#FF0000") // Red for finished
	ColorNeutral  = lipgloss.Color("#FFFFFF") // White for clock/neutral
	ColorDim      = lipgloss.Color("#666666") // Dim gray
	ColorAccent   = lipgloss.Color("#00CCFF") // Cyan accent
)

// Styles
var (
	// Title style for mode name
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginBottom(1)

	// Large time display style
	TimeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorNeutral)

	// Phase indicator (WORK/REST)
	PhaseWorkStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWork).
			MarginTop(1)

	PhaseRestStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorRest).
			MarginTop(1)

	// Round counter
	RoundStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(1)

	// Help bar at bottom
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(2)

	// Settings style
	SettingStyle = lipgloss.NewStyle().
			Foreground(ColorNeutral)

	SettingSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	// Container style for centering
	ContainerStyle = lipgloss.NewStyle()
)

// Big digit font (7 segment style)
var bigDigits = map[rune][]string{
	'0': {
		" ███ ",
		"█   █",
		"█   █",
		"█   █",
		"█   █",
		"█   █",
		" ███ ",
	},
	'1': {
		"  █  ",
		" ██  ",
		"  █  ",
		"  █  ",
		"  █  ",
		"  █  ",
		" ███ ",
	},
	'2': {
		" ███ ",
		"█   █",
		"    █",
		"  ██ ",
		" █   ",
		"█    ",
		"█████",
	},
	'3': {
		" ███ ",
		"█   █",
		"    █",
		"  ██ ",
		"    █",
		"█   █",
		" ███ ",
	},
	'4': {
		"█   █",
		"█   █",
		"█   █",
		"█████",
		"    █",
		"    █",
		"    █",
	},
	'5': {
		"█████",
		"█    ",
		"█    ",
		"████ ",
		"    █",
		"█   █",
		" ███ ",
	},
	'6': {
		" ███ ",
		"█   █",
		"█    ",
		"████ ",
		"█   █",
		"█   █",
		" ███ ",
	},
	'7': {
		"█████",
		"    █",
		"   █ ",
		"  █  ",
		"  █  ",
		"  █  ",
		"  █  ",
	},
	'8': {
		" ███ ",
		"█   █",
		"█   █",
		" ███ ",
		"█   █",
		"█   █",
		" ███ ",
	},
	'9': {
		" ███ ",
		"█   █",
		"█   █",
		" ████",
		"    █",
		"█   █",
		" ███ ",
	},
	':': {
		"     ",
		"  █  ",
		"  █  ",
		"     ",
		"  █  ",
		"  █  ",
		"     ",
	},
	' ': {
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
		"     ",
	},
}

// RenderBigTime renders time string as large ASCII art
func RenderBigTime(timeStr string, color lipgloss.Color) string {
	lines := make([]string, 7)

	for _, char := range timeStr {
		digit, ok := bigDigits[char]
		if !ok {
			digit = bigDigits[' ']
		}
		for i := 0; i < 7; i++ {
			lines[i] += digit[i] + " "
		}
	}

	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	result := ""
	for _, line := range lines {
		result += style.Render(line) + "\n"
	}

	return result
}

// CenterInScreen centers content in the given dimensions
func CenterInScreen(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

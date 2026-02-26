package audio

import (
	"math"
	"os"
	"os/exec"
	"sync"
)

// Player handles audio playback
type Player struct {
	beepPath  string
	enabled   bool
	mu        sync.Mutex
	useAplay  bool
	usePaplay bool
}

// New creates a new audio player
func New(beepPath string) *Player {
	p := &Player{
		beepPath: beepPath,
		enabled:  true,
	}

	// Check which audio player is available
	if _, err := exec.LookPath("paplay"); err == nil {
		p.usePaplay = true
	} else if _, err := exec.LookPath("aplay"); err == nil {
		p.useAplay = true
	}

	return p
}

// SetEnabled enables or disables audio
func (p *Player) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
}

// IsEnabled returns whether audio is enabled
func (p *Player) IsEnabled() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.enabled
}

// PlayBeep plays the beep sound
func (p *Player) PlayBeep() {
	p.mu.Lock()
	if !p.enabled {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	go p.playSound(p.beepPath)
}

// PlayCountdown plays a countdown beep (higher pitch for final countdown)
func (p *Player) PlayCountdown(secondsRemaining int) {
	p.mu.Lock()
	if !p.enabled {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	go p.playSound(p.beepPath)
}

// PlayIntervalChange plays sound for work/rest transition
func (p *Player) PlayIntervalChange(isWork bool) {
	p.mu.Lock()
	if !p.enabled {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	// Play beep for interval change
	go p.playSound(p.beepPath)
}

// PlayFinish plays the completion sound
func (p *Player) PlayFinish() {
	p.mu.Lock()
	if !p.enabled {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	// Play beep for finish
	go p.playSound(p.beepPath)
}

func (p *Player) playSound(path string) {
	var cmd *exec.Cmd

	if p.usePaplay {
		cmd = exec.Command("paplay", path)
	} else if p.useAplay {
		cmd = exec.Command("aplay", "-q", path)
	} else {
		// No audio player available
		return
	}

	_ = cmd.Start()
}

// GenerateBeepWAV generates a simple beep WAV file
// This creates a 880Hz sine wave for 200ms
func GenerateBeepWAV(path string) error {
	// WAV file parameters
	sampleRate := 44100
	duration := 0.2 // 200ms
	frequency := 880.0
	amplitude := 0.5

	numSamples := int(float64(sampleRate) * duration)

	// WAV header
	header := make([]byte, 44)

	// RIFF header
	copy(header[0:4], "RIFF")
	dataSize := numSamples * 2 // 16-bit samples
	fileSize := 36 + dataSize
	header[4] = byte(fileSize)
	header[5] = byte(fileSize >> 8)
	header[6] = byte(fileSize >> 16)
	header[7] = byte(fileSize >> 24)
	copy(header[8:12], "WAVE")

	// fmt chunk
	copy(header[12:16], "fmt ")
	header[16] = 16 // chunk size
	header[17] = 0
	header[18] = 0
	header[19] = 0
	header[20] = 1 // PCM format
	header[21] = 0
	header[22] = 1 // mono
	header[23] = 0
	header[24] = byte(sampleRate)
	header[25] = byte(sampleRate >> 8)
	header[26] = byte(sampleRate >> 16)
	header[27] = byte(sampleRate >> 24)
	byteRate := sampleRate * 2
	header[28] = byte(byteRate)
	header[29] = byte(byteRate >> 8)
	header[30] = byte(byteRate >> 16)
	header[31] = byte(byteRate >> 24)
	header[32] = 2 // block align
	header[33] = 0
	header[34] = 16 // bits per sample
	header[35] = 0

	// data chunk
	copy(header[36:40], "data")
	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	// Generate samples
	data := make([]byte, dataSize)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		// Apply envelope to avoid clicks
		envelope := 1.0
		attackSamples := sampleRate / 100 // 10ms attack
		releaseSamples := sampleRate / 50 // 20ms release
		if i < attackSamples {
			envelope = float64(i) / float64(attackSamples)
		} else if i > numSamples-releaseSamples {
			envelope = float64(numSamples-i) / float64(releaseSamples)
		}

		sample := amplitude * envelope * math.Sin(2*math.Pi*frequency*t)
		intSample := int16(sample * 32767)
		data[i*2] = byte(intSample)
		data[i*2+1] = byte(intSample >> 8)
	}

	// Write file
	file := append(header, data...)
	return os.WriteFile(path, file, 0644)
}

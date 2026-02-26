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
	chimePath string
	enabled   bool
	mu        sync.Mutex
	useAplay  bool
	usePaplay bool
}

// New creates a new audio player
func New(beepPath, chimePath string) *Player {
	p := &Player{
		beepPath:  beepPath,
		chimePath: chimePath,
		enabled:   true,
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

// PlayCountdown plays a countdown beep for 3-2-1
func (p *Player) PlayCountdown(secondsRemaining int) {
	p.mu.Lock()
	if !p.enabled {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	go p.playSound(p.beepPath)
}

// PlayChime plays the chime sound for interval changes
func (p *Player) PlayChime() {
	p.mu.Lock()
	if !p.enabled {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	go p.playSound(p.chimePath)
}

// PlayIntervalChange plays chime for work/rest transition (at 0)
func (p *Player) PlayIntervalChange(isWork bool) {
	p.PlayChime()
}

// PlayFinish plays the chime for completion
func (p *Player) PlayFinish() {
	p.PlayChime()
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

// GenerateBeepWAV generates a short beep for countdown (880Hz, 150ms)
func GenerateBeepWAV(path string) error {
	sampleRate := 44100
	duration := 0.15 // 150ms - short beep
	frequency := 880.0
	amplitude := 0.5

	return generateTone(path, sampleRate, duration, frequency, amplitude)
}

// GenerateChimeWAV generates a pleasant chime sound (two-tone descending)
func GenerateChimeWAV(path string) error {
	sampleRate := 44100
	duration := 0.5 // 500ms total
	amplitude := 0.4

	numSamples := int(float64(sampleRate) * duration)

	// WAV header
	header := makeWAVHeader(sampleRate, numSamples)

	// Generate two-tone chime (C5 -> G4, like a doorbell)
	freq1 := 523.25 // C5
	freq2 := 392.00 // G4

	data := make([]byte, numSamples*2)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(numSamples)

		// Envelope: quick attack, longer decay
		envelope := 1.0
		attackSamples := sampleRate / 50 // 20ms attack
		if i < attackSamples {
			envelope = float64(i) / float64(attackSamples)
		} else {
			// Exponential decay
			decayProgress := float64(i-attackSamples) / float64(numSamples-attackSamples)
			envelope = math.Exp(-3 * decayProgress)
		}

		// Crossfade between two frequencies
		blend1 := 1.0 - progress
		blend2 := progress

		sample := amplitude * envelope * (blend1*math.Sin(2*math.Pi*freq1*t) +
			blend2*math.Sin(2*math.Pi*freq2*t))

		intSample := int16(sample * 32767)
		data[i*2] = byte(intSample)
		data[i*2+1] = byte(intSample >> 8)
	}

	file := append(header, data...)
	return os.WriteFile(path, file, 0644)
}

func generateTone(path string, sampleRate int, duration, frequency, amplitude float64) error {
	numSamples := int(float64(sampleRate) * duration)

	header := makeWAVHeader(sampleRate, numSamples)

	data := make([]byte, numSamples*2)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)

		// Envelope to avoid clicks
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

	file := append(header, data...)
	return os.WriteFile(path, file, 0644)
}

func makeWAVHeader(sampleRate, numSamples int) []byte {
	header := make([]byte, 44)

	dataSize := numSamples * 2
	fileSize := 36 + dataSize

	copy(header[0:4], "RIFF")
	header[4] = byte(fileSize)
	header[5] = byte(fileSize >> 8)
	header[6] = byte(fileSize >> 16)
	header[7] = byte(fileSize >> 24)
	copy(header[8:12], "WAVE")

	copy(header[12:16], "fmt ")
	header[16] = 16
	header[17] = 0
	header[18] = 0
	header[19] = 0
	header[20] = 1 // PCM
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

	copy(header[36:40], "data")
	header[40] = byte(dataSize)
	header[41] = byte(dataSize >> 8)
	header[42] = byte(dataSize >> 16)
	header[43] = byte(dataSize >> 24)

	return header
}

package audio

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var (
	speakerMu   sync.Mutex
	speakerInit bool
)

type AudioService struct {
	basePath string
}

func NewAudioService(basePath string) *AudioService {
	service := &AudioService{
		basePath: basePath,
	}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		log.Printf("⚠️ Audio papkasi topilmadi: %s", basePath)
	} else {
		log.Printf("✅ Audio Service yaratildi: %s", basePath)
	}

	return service
}

func (a *AudioService) InitSpeaker(sampleRate beep.SampleRate) error {
	speakerMu.Lock()
	defer speakerMu.Unlock()

	if speakerInit {
		return nil
	}

	bufferSize := sampleRate.N(time.Second / 10)
	if runtime.GOOS == "windows" {
		bufferSize = sampleRate.N(time.Second / 5)
	}

	err := speaker.Init(sampleRate, bufferSize)
	if err != nil {
		return fmt.Errorf("speaker init: %w", err)
	}

	speakerInit = true
	log.Printf("🔊 Speaker initialized (SampleRate: %d)", sampleRate)
	return nil
}

func (a *AudioService) PlayAudio(filename string) error {
	fullPath := filepath.Join(a.basePath, filename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("fayl topilmadi: %s", fullPath)
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("fayl ochilmadi: %w", err)
	}
	defer f.Close()

	var streamer beep.StreamSeekCloser
	var format beep.Format

	ext := filepath.Ext(filename)
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	default:
		return fmt.Errorf("format qo'llab-quvvatlanmaydi: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("decode xato: %w", err)
	}
	defer streamer.Close()

	if err := a.InitSpeaker(format.SampleRate); err != nil {
		return err
	}

	done := make(chan bool, 1) // ⚠️ Buffered channel!

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		select {
		case done <- true:
		default:
		}
	})))

	// ⚠️ Timeout qo'shamiz - deadlock oldini olish
	select {
	case <-done:
		return nil
	case <-time.After(10 * time.Second):
		speaker.Clear()
		return fmt.Errorf("audio timeout")
	}
}

func (a *AudioService) PlayNumber(number int) error {
	if number <= 0 {
		return fmt.Errorf("noto'g'ri raqam: %d", number)
	}

	switch {
	case number == 10:
		return a.PlayAudio("numbers/10.mp3")
	case number <= 9:
		return a.PlayAudio(fmt.Sprintf("numbers/%d.mp3", number))
	case number <= 19:
		if err := a.PlayAudio("numbers/10a.mp3"); err != nil {
			return err
		}
		return a.PlayAudio(fmt.Sprintf("numbers/%d.mp3", number-10))
	case number <= 99:
		tens := (number / 10) * 10
		ones := number % 10
		if err := a.PlayAudio(fmt.Sprintf("numbers/%da.mp3", tens)); err != nil {
			return err
		}
		if ones > 0 {
			return a.PlayAudio(fmt.Sprintf("numbers/%d.mp3", ones))
		}
		return nil
	default:
		return a.PlayAudio(fmt.Sprintf("numbers/%d.mp3", number))
	}
}

func (a *AudioService) PlayRoomNumber(room string) error {
	return a.PlayAudio(fmt.Sprintf("numbers/%s-xona.mp3", room))
}

func (a *AudioService) PlayPhrase(phrase string) error {
	return a.PlayAudio(fmt.Sprintf("phrases/%s.mp3", phrase))
}

func (a *AudioService) PlayAnnouncement(queueNumber string, roomNumber string) error {
	log.Printf("\n🎵 ===== AUDIO E'LON BOSHLANDI =====")
	log.Printf("📋 Navbat: %s, Xona: %s", queueNumber, roomNumber)
	startTime := time.Now()

	ticketNum := extractNumberFromQueue(queueNumber)
	log.Printf("🔢 Ajratilgan raqam: %d", ticketNum)

	steps := []struct {
		action func() error
		name   string
	}{
		{func() error { return a.PlayNumber(ticketNum) }, "Navbat raqami"},
		{func() error { return a.PlayPhrase("raqam_egasi") }, "Raqam egasi"},
		{func() error { return a.PlayRoomNumber(roomNumber) }, "Xona raqami"},
		{func() error { return a.PlayPhrase("honaga_kelishin") }, "Chaqiriq"},
	}

	for _, step := range steps {
		stepStart := time.Now()
		if err := step.action(); err != nil {
			log.Printf("⚠️ %s xato: %v", step.name, err)
		} else {
			log.Printf("✅ %s ijro etildi (%v)", step.name, time.Since(stepStart))
		}
	}

	log.Printf("✅ ===== AUDIO E'LON TUGADI (%v) =====\n", time.Since(startTime))
	return nil
}

func (a *AudioService) Close() {
	speakerMu.Lock()
	defer speakerMu.Unlock()

	log.Println("🔇 Audio Service yopilmoqda...")
	if speakerInit {
		speaker.Clear()
		log.Println("🔊 Speaker tozalandi")
	}
}

func extractNumberFromQueue(queueNumber string) int {
	parts := strings.Split(queueNumber, "-")
	if len(parts) > 1 {
		cleanNum := strings.TrimLeft(parts[1], "0")
		if cleanNum == "" {
			cleanNum = "0"
		}
		if num, err := strconv.Atoi(cleanNum); err == nil {
			return num
		}
	}

	if num, err := strconv.Atoi(queueNumber); err == nil {
		return num
	}

	return 0
}

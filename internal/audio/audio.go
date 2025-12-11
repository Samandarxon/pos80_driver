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

// 🎵 Fayl formatini avtomatik aniqlash (.mp3 yoki .wav)
func (a *AudioService) findAudioFile(baseName string) (string, error) {
	// Avval .mp3 ni tekshiramiz
	mp3Path := baseName + ".mp3"
	fullMp3Path := filepath.Join(a.basePath, mp3Path)
	if _, err := os.Stat(fullMp3Path); err == nil {
		return mp3Path, nil
	}

	// Keyin .wav ni tekshiramiz
	wavPath := baseName + ".wav"
	fullWavPath := filepath.Join(a.basePath, wavPath)
	if _, err := os.Stat(fullWavPath); err == nil {
		return wavPath, nil
	}

	return "", fmt.Errorf("audio fayl topilmadi: %s (.mp3 yoki .wav)", baseName)
}

func (a *AudioService) PlayNumber(number int) error {
	if number <= 0 {
		return fmt.Errorf("noto'g'ri raqam: %d", number)
	}

	switch {
	case number == 10:
		audioFile, err := a.findAudioFile("numbers/10a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number <= 9:
		audioFile, err := a.findAudioFile(fmt.Sprintf("numbers/%d", number))
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number <= 19:
		audioFile, err := a.findAudioFile("numbers/10")
		if err != nil {
			return err
		}
		if err := a.PlayAudio(audioFile); err != nil {
			return err
		}

		onesFile, err := a.findAudioFile(fmt.Sprintf("numbers/%d", number-10))
		if err != nil {
			return err
		}
		return a.PlayAudio(onesFile)
	case number == 20:
		audioFile, err := a.findAudioFile("numbers/20a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number == 30:
		audioFile, err := a.findAudioFile("numbers/30a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number == 40:
		audioFile, err := a.findAudioFile("numbers/40a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number == 50:
		audioFile, err := a.findAudioFile("numbers/50a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number == 60:
		audioFile, err := a.findAudioFile("numbers/60a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number == 70:
		audioFile, err := a.findAudioFile("numbers/70a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number == 80:
		audioFile, err := a.findAudioFile("numbers/80a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number == 90:
		audioFile, err := a.findAudioFile("numbers/90a")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number <= 99:
		tens := (number / 10) * 10
		ones := number % 10

		tensFile, err := a.findAudioFile(fmt.Sprintf("numbers/%d", tens))
		if err != nil {
			return err
		}
		if err := a.PlayAudio(tensFile); err != nil {
			return err
		}

		if ones > 0 {
			onesFile, err := a.findAudioFile(fmt.Sprintf("numbers/%d", ones))
			if err != nil {
				return err
			}
			return a.PlayAudio(onesFile)
		}
		return nil
	case number == 100:
		audioFile, err := a.findAudioFile("numbers/100")
		if err != nil {
			return err
		}
		return a.PlayAudio(audioFile)
	case number <= 199:
		// 101-199: "yuz" + raqam
		audioFile, err := a.findAudioFile("numbers/100")
		if err != nil {
			return err
		}
		if err := a.PlayAudio(audioFile); err != nil {
			return err
		}
		if number > 100 {
			return a.PlayNumber(number - 100)
		}
		return nil
	default:
		return fmt.Errorf("100 dan katta raqamlar uchun audio fayl yo'q: %d", number)
	}
}

func (a *AudioService) PlayRoomNumber(room string) error {
	audioFile, err := a.findAudioFile(fmt.Sprintf("numbers/%s-xona", room))
	if err != nil {
		return err
	}
	return a.PlayAudio(audioFile)
}

func (a *AudioService) PlayPhrase(phrase string) error {
	audioFile, err := a.findAudioFile(fmt.Sprintf("phrases/%s", phrase))
	if err != nil {
		return err
	}
	return a.PlayAudio(audioFile)
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

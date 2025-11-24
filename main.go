package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Cross-platform speaker initialization
func initSpeaker(sampleRate beep.SampleRate) error {
	// Windows uchun buffer size o'zgartirish
	bufferSize := sampleRate.N(time.Second / 10)

	// Agar Windows bo'lsa, buffer size ni oshirish
	if runtime.GOOS == "windows" {
		bufferSize = sampleRate.N(time.Second / 5)
	}

	return speaker.Init(sampleRate, bufferSize)
}

func playAudio(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("faylni ochib bo'lmadi: %w", err)
	}
	defer f.Close()

	var streamer beep.StreamSeekCloser
	var format beep.Format

	// Formatni aniqlash
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
		return fmt.Errorf("decode qilishda xato: %w", err)
	}
	defer streamer.Close()

	// Speaker initialization
	err = initSpeaker(format.SampleRate)
	if err != nil {
		return fmt.Errorf("speaker initialization: %w", err)
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Platforma: %s\n", runtime.GOOS)
		fmt.Println("Foydalanish: go run main.go <audio-file>")
		fmt.Println("Misol: go run main.go sounds/notification2.wav")
		fmt.Println("Misol: go run main.go sounds/numbers/1.mp3")
		return
	}

	filename := os.Args[1]

	// Fayl mavjudligini tekshirish
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("Fayl topilmadi: %s", filename)
	}

	fmt.Printf("Platforma: %s\n", runtime.GOOS)
	fmt.Printf("Ijro etilmoqda: %s\n", filename)

	err := playAudio(filename)
	if err != nil {
		log.Fatalf("Xato: %v", err)
	}

	fmt.Println("Audio muvaffaqiyatli ijro etildi!")
}

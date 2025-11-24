package main

import (
	"log"
	"os"
	"os/signal"
	"pos80/internal/api"
	"pos80/internal/audio"
	"pos80/internal/config"
	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// ==============================
	// DASTURNI ISHGA TUSHIRISH BOSQICHLARI
	// ==============================

	// 1. KONFIGURATSIYA YUKLASH
	log.Printf("📋 Konfiguratsiya yuklanmoqda...")

	// 2. AUDIO SERVICE YARATISH
	log.Printf("🎵 Audio servis yaratilmoqda...")
	audioService := audio.NewAudioService("./sounds")

	// Audio papka mavjudligini tekshirish
	if _, err := os.Stat("./sounds"); os.IsNotExist(err) {
		log.Printf("⚠️ Diqqat: 'sounds' papkasi topilmadi! Audio ishlamaydi.")
	} else {
		log.Printf("✅ Audio fayllar papkasi topildi")
	}

	// 🎯 AUDIO QUEUE SERVICE NI YARATISH VA ISHGA TUSHIRISH
	log.Printf("🚀 Audio Queue Service yaratilmoqda...")
	audioQueue := audio.NewAudioQueueService(audioService, 1) // ⚠️ 1 ta worker - serial execution
	audioQueue.Start()
	log.Printf("✅ Audio Queue Service ishga tushdi")

	// 3. ROUTER SOZLASH
	router := gin.New()
	api.SetupRouter(router, audioService, audioQueue) // ⚠️ audioQueue ni ham o'tkazamiz

	// ==============================
	// GRACEFUL SHUTDOWN SOZLASH
	// ==============================
	setupGracefulShutdown(audioService, audioQueue)

	// ==============================
	// SERVERNI ISHGA TUSHIRISH
	// ==============================
	log.Printf("🚀 %s v%s ishga tushmoqda...", config.AppName, config.AppVersion)
	log.Printf("📍 Server manzili: http://0.0.0.0%s", config.ServicePort)
	log.Printf("🖨️  Default printer: %s", config.DefaultPrinterName)
	log.Printf("💻 Platforma: %s", runtime.GOOS)
	log.Printf("📊 Rejim: Production")

	// Server ishga tushirish
	if err := router.Run(config.ServicePort); err != nil {
		log.Fatalf("🔥 Serverni ishga tushirib bo'lmadi: %v", err)
	}
}

func setupGracefulShutdown(audioService *audio.AudioService, audioQueue *audio.AudioQueueService) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("")
		log.Println("🛑 Graceful shutdown boshlandi...")
		log.Println("⏳ Resurslar tozalanmoqda...")

		// Audio queue ni to'xtatish
		audioQueue.Stop()
		log.Println("✅ Audio Queue to'xtatildi")

		// Audio service ni yopish
		audioService.Close()
		log.Println("✅ Audio Service yopildi")

		log.Println("✅ Barcha resurslar tozalandi")
		log.Println("👋 Dastur to'xtatildi")
		os.Exit(0)
	}()
}

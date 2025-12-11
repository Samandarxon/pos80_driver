package api

import (
	"fmt"
	"log"
	"pos80/internal/api/handlers"
	"pos80/internal/audio"
	"pos80/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, audioService *audio.AudioService, audioQueue *audio.AudioQueueService) {

	printHandler := handlers.NewPrintHandler(config.DefaultPrinterName)

	// ⚠️ AudioHandler ga audioQueue ni uzatamiz (audioService emas!)
	audioHandler := handlers.NewAudioHandlerWithQueue(audioQueue)

	// ==============================
	// GLOBAL MIDDLEWARE (tartib muhim!)
	// ==============================

	router.Use(gin.Recovery())  // 1. Panic recovery
	router.Use(handlers.Cors()) // 2. CORS (birinchi bo'lishi kerak!)

	log.Printf("🔧 Middleware lar o'rnatildi")

	router.GET("/", func(c *gin.Context) {
		fmt.Println("Health check endpoint called")
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "success",
		})
	})

	// ==============================
	// PUBLIC ROUTES (API Key'siz)
	router.GET("/api/audio/play", audioHandler.HandlePlayAudio)
	router.POST("/api/audio/announcement", audioHandler.HandleAudioAnnouncement)
	// ==============================

	// YANGI QUEUE ENDPOINTLAR
	router.GET("/api/audio/queue/status", audioHandler.HandleQueueStatus)
	router.POST("/api/audio/queue/clear", audioHandler.HandleClearQueue)
	router.GET("/api/audio/health", audioHandler.HandleHealth)

	// ==============================
	// PROTECTED ROUTES (API Key bilan)
	// ==============================

	api := router.Group("/")
	api.Use(handlers.PrintGuardMiddleware()) // Faqat bu group uchun
	{
		api.POST("/print-ticket", printHandler.HandlePrintTicket)
	}

	log.Printf("🌐 API route lar belgilandi")
}

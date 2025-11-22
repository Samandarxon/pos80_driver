package api

import (
	"log"
	"pos80/internal/api/handlers"
	"pos80/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {

	// 2. SERVISLARNI YARATISH
	// Print va Health handlerlarini ishga tushirish
	printHandler := handlers.NewPrintHandler(config.DefaultPrinterName)

	// ==============================
	// MIDDLEWARE LARNI O'RNATISH
	// ==============================

	// Logging middleware - barcha so'rovlarni log qilish
	router.Use(handlers.PrintGuardMiddleware())

	// CORS middleware - brauzer cross-origin so'rovlariga ruxsat berish
	router.Use(handlers.CORSMiddleware())

	// Recovery middleware - dastur xatosiz ishlashini ta'minlash
	router.Use(gin.Recovery())

	log.Printf("🔧 Middleware lar o'rnatildi")

	// ==============================
	// API ROUTE LARNI BELGILASH
	// ==============================

	// Health check endpoint - sistemaning holatini tekshirish
	router.GET("/health", printHandler.CheckHealth)

	// Asosiy chipta chop etish endpoint
	router.POST("/print-ticket", printHandler.HandlePrintTicket)

	// Test chipta chop etish endpoint
	router.POST("/print-test", printHandler.PrintTestTicket)

	log.Printf("🌐 API route lar belgilandi")
}

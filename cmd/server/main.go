package main

import (
	"log"
	"pos80/internal/api"
	"pos80/internal/config"

	"github.com/gin-gonic/gin"
)

// ============================================
// ASOSIY DASTUR (MAIN APPLICATION)
// ============================================

// main - dasturning kirish nuqtasi (entry point)
// Bu yerda barcha komponentlar birlashtiriladi va server ishga tushiriladi
func main() {
	// ==============================
	// DASTURNI ISHGA TUSHIRISH BOSQICHLARI
	// ==============================

	// 1. KONFIGURATSIYA YUKLASH
	// Dastur sozlamalari environment va konstantalardan yuklanadi
	log.Printf("📋 Konfiguratsiya yuklanmoqda...")

	log.Printf("✅ Servislar yaratildi")

	// 3. ROUTER (YO'NALTIRGICH) SOZLASH
	// Gin framework orqali HTTP router konfiguratsiyasi
	router := gin.New()

	api.SetupRouter(router)

	// ==============================
	// SERVERNI ISHGA TUSHIRISH
	// ==============================

	// Dastur haqida ma'lumot chiqarish
	log.Printf("🚀 %s v%s ishga tushmoqda...", config.AppName, config.AppVersion)
	log.Printf("📍 Server manzili: http://0.0.0.0%s", config.ServicePort)
	log.Printf("🖨️  Default printer: %s", config.DefaultPrinterName)
	log.Printf("💻 Platforma: Windows")
	log.Printf("📊 Rejim: Production")

	// Serverni ishga tushirish
	// Agar port band bo'lsa yoki boshqa xatolik yuz bersa, dastur to'xtaydi
	if err := router.Run(config.ServicePort); err != nil {
		log.Fatalf("❌ Serverni ishga tushirib bo'lmadi: %v", err)
	}

	// Server ishlashni boshlagandan so'ng, ushbu kodga hech qachon yetib kelmaydi
	// Chunki router.Run() blokirovka qiluvchi metod
}

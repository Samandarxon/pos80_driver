// ============================================
// API MIDDLEWARE - KASALXONA PRINTER TIZIMI
// HTTP so'rovlarini boshqarish va monitoring qilish
// ============================================

package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// ==============================
// LOGGING MIDDLEWARE
// ==============================

// LoggingMiddleware - barcha HTTP so'rovlarni log qiladi
// Bu middleware har bir so'rov haqida quyidagi ma'lumotlarni yozib boradi:
// - HTTP metod (GET, POST, ...)
// - So'rov yo'li (path)
// - HTTP status kod (200, 400, 500, ...)
// - So'rov bajarilish vaqti
//
// Bu ma'lumotlar:
// ✅ Monitoring va debugging uchun
// ✅ Performance tahlili uchun
// ✅ Security audit uchun
// ✅ Traffic tahlili uchun
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// So'rov boshlanish vaqtini yozib olish
		start := time.Now()

		// So'rov yo'lini olish
		path := c.Request.URL.Path

		// Keyingi middleware/handlerni chaqirish
		// Bu yerda so'rov hali bajarilmagan
		c.Next()

		// So'rov tugagandan so'ng log yozish
		// Bu yerda:
		// - c.Writer.Status(): HTTP status kodi
		// - time.Since(start): So'rov bajarilish vaqti
		log.Printf("%s %s %d %v",
			c.Request.Method,  // HTTP metod (GET, POST, ...)
			path,              // So'rov yo'li (/health, /print-ticket, ...)
			c.Writer.Status(), // HTTP status kodi (200, 404, 500, ...)
			time.Since(start), // So'rov bajarilish vaqti
		)
	}
}

// ==============================
// CORS MIDDLEWARE
// ==============================

// CORSMiddleware - Cross-Origin Resource Sharing sozlamalarini boshqaradi
// Bu middleware veb-brauzerlarning cross-domain so'rovlariga ruxsat beradi.
//
// CORS nima uchun kerak?
// - Frontend (React, Vue) va Backend turli domainlarda bo'lsa
// - Mobile app'lar API ga murojaat qilsa
// - Third-party integratsiyalar uchun
//
// Xavfsizlik sozlamalari:
// - Faqat kerakli metodlarga ruxsat
// - Faqat kerakli headerlarga ruxsat
// - Preflight so'rovlarini boshqarish
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// CORS HEADERLARNI O'RNATISH

		// Access-Control-Allow-Origin: qaysi domainlar so'rov yuborishi mumkin
		// "*" - barcha domainlarga ruxsat berish (productionda specific domain berish tavsiya etiladi)
		c.Header("Access-Control-Allow-Origin", "*")

		// Access-Control-Allow-Methods: qaysi HTTP metodlarga ruxsat
		// GET    - ma'lumot olish
		// POST   - ma'lumot yuborish
		// OPTIONS - preflight so'rovlari
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Access-Control-Allow-Headers: qaysi headerlarga ruxsat
		// Content-Type  - JSON, form-data, etc.
		// Authorization - token, API key, etc.
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// PREFLIGHT SO'ROVLARNI BOSHQARISH
		// Brauzer murakkab so'rovlardan oldin OPTIONS so'rov yuboradi
		// Bu so'rov server qanday so'rovlarga ruxsat berishini tekshiradi
		if c.Request.Method == "OPTIONS" {
			// 204 No Content - so'rov muvaffaqiyatli, lekin hech qanday content yo'q
			c.AbortWithStatus(204)
			return
		}

		// Keyingi middleware/handlerni chaqirish
		c.Next()
	}
}

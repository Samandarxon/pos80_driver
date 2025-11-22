// ============================================
// API MIDDLEWARE - KASALXONA PRINTER TIZIMI
// HTTP so'rovlarini boshqarish va monitoring qilish
// ============================================

package handlers

import (
	"net"
	"net/http"
	"sync"
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
// func LoggingMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// So'rov boshlanish vaqtini yozib olish
// 		start := time.Now()

// 		// So'rov yo'lini olish
// 		path := c.Request.URL.Path

// 		authHeader := c.GetHeader("Seckret")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"code":    401,
// 				"message": "Authorization header required",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		// Token ni header dan ajratib olish
// 		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{
// 				"code":    401,
// 				"message": err.Error(),
// 			})
// 			c.Abort()
// 			return
// 		}
// 		// Keyingi middleware/handlerni chaqirish
// 		// Bu yerda so'rov hali bajarilmagan
// 		c.Next()

// 		// So'rov tugagandan so'ng log yozish
// 		// Bu yerda:
// 		// - c.Writer.Status(): HTTP status kodi
// 		// - time.Since(start): So'rov bajarilish vaqti
// 		log.Printf("%s %s %d %v",
// 			c.Request.Method,  // HTTP metod (GET, POST, ...)
// 			path,              // So'rov yo'li (/health, /print-ticket, ...)
// 			c.Writer.Status(), // HTTP status kodi (200, 404, 500, ...)
// 			time.Since(start), // So'rov bajarilish vaqti
// 		)
// 	}
// }

var (
	allowedAPIKey  = "SECRET-PRINTER-KEY-b21ecca4618d929c6f24e0f7245ca7b50740f6509e455f3b1c165d70" // configdan o'qiladi
	rateLimitStore = make(map[string][]time.Time)
	mu             sync.Mutex
)

func PrintGuardMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. API Key check
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != allowedAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid API Key",
			})
			c.Abort()
			return
		}

		// 2. Rate limit
		ip := getClientIP(c)
		if !rateLimitAllow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Too many print requests, slow down",
			})
			c.Abort()
			return
		}

		// 3. Empty body check
		if c.Request.ContentLength == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Request body cannot be empty",
			})
			c.Abort()
			return
		}

		c.Next() // Keyingi handlerga o‘tadi
	}
}

// IP olish
func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	host, _, err := net.SplitHostPort(ip)
	if err == nil {
		return host
	}
	return ip
}

// sliding window rate limit
func rateLimitAllow(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	window := now.Add(-10 * time.Second)

	reqs := rateLimitStore[ip]
	valid := []time.Time{}

	for _, t := range reqs {
		if t.After(window) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= 5 {
		return false
	}

	valid = append(valid, now)
	rateLimitStore[ip] = valid
	return true
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

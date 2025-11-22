package handlers

import (
	"log"
	"net/http"
	"pos80/internal/config"
	"pos80/internal/models"
	"pos80/internal/printer"
	"time"

	"github.com/gin-gonic/gin"
)

// ==============================
// CHIPTA CHOP ETISH HANDLERI
// ==============================

// PrintHandler - chipta chop etish operatsiyalarini boshqaradi
// Ushbu struct HTTP so'rovlarini qabul qiladi, ma'lumotlarni tekshiradi,
// chiptani formatlaydi va printerni boshqaradi.
type PrintHandler struct {
	printerService  *printer.PrinterService  // Printer bilan aloqa qiluvchi servis
	ticketFormatter *printer.TicketFormatter // Chiptani ESC/POS formatiga o'girovchi

}

// NewPrintHandler - yangi PrintHandler yaratadi
// printerName: printer nomi (masalan: "POS80")
// Qaytaradi: yangi PrintHandler instance
func NewPrintHandler(printerName string) *PrintHandler {
	return &PrintHandler{
		printerService:  printer.NewPrinterService(printerName),
		ticketFormatter: printer.NewTicketFormatter(),
	}
}

// HandlePrintTicket - asosiy chipta chop etish endpointi
// Bu metod:
// 1. JSON so'rovni tekshiradi va parse qiladi
// 2. Chiptani ESC/POS formatiga o'giradi
// 3. Printerni boshqarib chop etadi
// 4. Natijani foydalanuvchiga qaytaradi
func (h *PrintHandler) HandlePrintTicket(c *gin.Context) {
	var req models.PrintRequest

	// 1. JSON SO'ROVNI TEKSHIRISH VA PARSE QILISH
	// ShouldBindJSON - Gin frameworkining builtin validatsiya metodi
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST",
			"Noto'g'ri JSON formati: "+err.Error())
		return
	}

	log.Printf("🖨️  Chipta chop etish so'rovi: %s (Ustuvor: %v)",
		req.QueueDisplay, req.IsPriority)

	// 2. CHIPTANI FORMATLASH
	// TicketFormatter chipta ma'lumotlarini printer tushunadigan ESC/POS formatiga o'giradi
	ticketData := h.ticketFormatter.Format(req)

	// 3. PRINTERGA YUBORISH VA CHOP ETISH
	bytesWritten, err := h.printerService.Print(ticketData)
	if err != nil {
		log.Printf("❌ Chop etishda xato: %v", err)
		h.sendErrorResponse(c, http.StatusInternalServerError, "PRINT_FAILED",
			"Chiptani chop etishda xato: "+err.Error())
		return
	}

	// 4. MUVAFFAQIYATLI JAVOB QAYTARISH
	log.Printf("✅ Chipta muvaffaqiyatli chop etildi: %s (%d bayt)",
		req.QueueDisplay, bytesWritten)
	h.sendSuccessResponse(c, req, bytesWritten)
}

// PrintTestTicket - test chipta chop etish endpointi
// Bu metod:
// - Standart test ma'lumotlari yaratadi
// - Printerning ishlashini tekshiradi
// - Tizimni sinovdan o'tkazish imkonini beradi
func (h *PrintHandler) PrintTestTicket(c *gin.Context) {
	// Test chipta ma'lumotlari - har doim bir xil
	testTicket := models.PrintRequest{
		TicketID:       "test-001",                      // Test chipta IDsi
		ShiftID:        "test-shift",                    // Test smena IDsi
		QueueNumber:    99,                              // Test navbat raqami
		QueueDisplay:   "T-099",                         // Ekranda ko'rinadigan format
		Status:         "waiting",                       // Navbat holati
		IsPriority:     true,                            // Ustuvorlik holati
		PriorityReason: "Test uchun",                    // Ustuvorlik sababi
		Notes:          "Bu test chiptasi",              // Qo'shimcha eslatma
		CreatedAt:      time.Now().Format(time.RFC3339), // Joriy vaqt
	}

	log.Printf("🧪 Test chipta chop etish: %s", testTicket.QueueDisplay)

	// Chiptani formatlash va chop etish
	ticketData := h.ticketFormatter.Format(testTicket)
	bytesWritten, err := h.printerService.Print(ticketData)
	if err != nil {
		h.sendErrorResponse(c, http.StatusInternalServerError, "TEST_PRINT_FAILED",
			"Test chiptani chop etishda xato: "+err.Error())
		return
	}

	log.Printf("✅ Test chipta muvaffaqiyatli chop etildi: %s (%d bayt)",
		testTicket.QueueDisplay, bytesWritten)
	h.sendSuccessResponse(c, testTicket, bytesWritten)
}

// ==============================
// JAVOB YUBORISH METODLARI
// ==============================

// sendErrorResponse - xatolik javobini yuboradi
// status: HTTP status kodi (400, 500, ...)
// errorCode: ichki xato kodi (INVALID_REQUEST, PRINT_FAILED, ...)
// message: foydalanuvchiga ko'rsatiladigan xabar
func (h *PrintHandler) sendErrorResponse(c *gin.Context, status int, errorCode, message string) {
	log.Printf("❌ XATO: %s - %s", errorCode, message)

	c.JSON(status, models.PrintResponse{
		Status:    "error",                         // Javob holati
		Error:     errorCode,                       // Xato kodi
		Message:   message,                         // Xato tavsifi
		Timestamp: time.Now().Format(time.RFC3339), // Xato vaqti
	})
}

// sendSuccessResponse - muvaffaqiyatli javob yuboradi
// req: chop etilgan chipta ma'lumotlari
// bytesWritten: printerga yuborilgan baytlar soni
func (h *PrintHandler) sendSuccessResponse(c *gin.Context, req models.PrintRequest, bytesWritten int) {
	c.JSON(http.StatusOK, models.PrintResponse{
		Status:    "success",                           // Javob holati
		Message:   "Chipta muvaffaqiyatli chop etildi", // Muvaffaqiyat xabari
		Printer:   h.printerService.PrinterName,        // Printer nomi
		Bytes:     bytesWritten,                        // Chop etilgan baytlar
		Ticket:    req.QueueDisplay,                    // Navbat raqami
		Priority:  req.IsPriority,                      // Ustuvorlik holati
		Timestamp: time.Now().Format(time.RFC3339),     // Chop etish vaqti
		Data: map[string]interface{}{ // Qo'shimcha ma'lumotlar
			"ticket_id":     req.TicketID,
			"queue_display": req.QueueDisplay,
			"status":        req.Status,
			"is_priority":   req.IsPriority,
			"created_at":    req.CreatedAt,
		},
	})
}

// ==============================
// HEALTH CHECK HANDLERI
// ==============================

// HealthHandler - sistemaning sog'ligini tekshiradi
// Ushbu handler monitoring va load balancing uchun kerak
// type HealthHandler struct {
// 	printerService *printer.PrinterService // Printer holatini tekshirish uchun
// }

// NewHealthHandler - yangi HealthHandler yaratadi
// func NewHealthHandler(printerName string) *HealthHandler {
// 	return &HealthHandler{
// 		printerService: printer.NewPrinterService(printerName),
// 	}
// }

// CheckHealth - sistemaning holatini tekshiradi
// Bu metod:
// - Printerning mavjudligini tekshiradi
// - Sistemaning umumiy holatini baholaydi
// - Monitoring tizimlari uchun ma'lumot taqdim etadi
func (h *PrintHandler) CheckHealth(c *gin.Context) {
	// Printer mavjudligini tekshirish
	printerErr := h.printerService.CheckPrinter()

	// Sistemaning holatini aniqlash
	status := "healthy"                        // Normal holat
	message := "Service is operating normally" // Normal xabar

	// Agar printer mavjud bo'lmasa
	if printerErr != nil {
		status = "degraded" // Qisman ishlayotgan holat
		message = "Printer mavjud emas: " + printerErr.Error()
	}

	log.Printf("🏥 Health check: %s - Printer: %v",
		status, printerErr == nil)

	// Health check javobini yuborish
	c.JSON(http.StatusOK, gin.H{
		"status":            status,                          // Sistemaning holati
		"message":           message,                         // Holat tavsifi
		"service":           config.AppName,                  // Dastur nomi
		"version":           config.AppVersion,               // Dastur versiyasi
		"timestamp":         time.Now().Format(time.RFC3339), // Tekshirish vaqti
		"printer_available": printerErr == nil,               // Printer mavjudligi
	})
}

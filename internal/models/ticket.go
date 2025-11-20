// ============================================
// MA'LUMOTLAR STRUKTURALARI - KASALXONA PRINTER TIZIMI
// Dasturning barcha data modellarini markazlashtirilgan boshqarish
// ============================================

package models

import (
	"fmt"
	"time"
)

// ==============================
// CHIPTA CHOP ETISH SO'ROVI
// ==============================

// PrintRequest - chipta chop etish uchun keladigan JSON so'rov strukturasi
// Bu struktura Gin framework tomonidan avtomatik validatsiya qilinadi
type PrintRequest struct {
	// TicketID - chiptaning unikal identifikatori
	// Format: UUID yoki database primary key
	// Misol: "c14da83e-36d1-4a28-8ce2-6bebab4a18cb"
	TicketID string `json:"ticket_id" binding:"required"`

	// ShiftID - smena yoki ish vaqtining identifikatori
	// Har bir smenada navbat raqamlari 1 dan boshlanadi
	// Misol: "c363d327-d777-4727-83a3-1cef84747664"
	ShiftID string `json:"shift_id" binding:"required"`

	// QueueNumber - navbatning raqamli ko'rinishi
	// Har kuni har bir bo'limda 1 dan boshlanadi
	// binding:"required,min=1" - majburiy va 1 dan kichik bo'lmasligi kerak
	QueueNumber int `json:"queue_number" binding:"required,min=1"`

	// QueueDisplay - ekranda ko'rsatiladigan format
	// Format: "K-001", "T-015", "X-042"
	// K - Kardiologiya, T - Terapiya, X - Xirurgiya
	QueueDisplay string `json:"queue_display" binding:"required"`

	// Status - navbatning joriy holati
	// Mumkin bo'lgan holatlar: waiting, called, in_progress, completed, cancelled
	Status string `json:"status" binding:"required"`

	// Notes - qo'shimcha eslatmalar yoki bemor haqida ma'lumot
	// Ixtiyoriy maydon - bo'sh bo'lishi mumkin
	// Misol: "Yurak og'rig'i bilan kelgan", "Allergiya tarixi bor"
	Notes string `json:"notes"`

	// IsPriority - navbat ustuvorligi holati
	// true - ustuvor navbat (nogironlar, homilador ayollar, qariyalar)
	// false - oddiy navbat
	IsPriority bool `json:"is_priority"`

	// PriorityReason - ustuvorlik sababi
	// IsPriority true bo'lgina talab qilinadi
	// Misol: "70 yoshdan oshgan", "Nogironligi bor", "Homilador ayol"
	PriorityReason string `json:"priority_reason"`

	// CreatedAt - chipta yaratilgan vaqt
	// Format: ISO 8601 (RFC3339) - "2025-01-18T14:30:00Z"
	// Database dan kelgan to'liq timestamp
	CreatedAt string `json:"created_at" binding:"required"`
}

// ==============================
// CHIPTA CHOP ETISH JAVOBI
// ==============================

// PrintResponse - chipta chop etish natijasining JSON javob strukturasi
// Bu struktura frontend yoki boshqa servislarga qaytariladi
type PrintResponse struct {
	// Status - operatsiya natijasi holati
	// "success" - muvaffaqiyatli chop etildi
	// "error" - chop etishda xato yuz berdi
	Status string `json:"status"`

	// Message - foydalanuvchiga ko'rsatiladigan xabar
	// omitempty - agar bo'sh bo'lsa, JSON da ko'rsatilmaydi
	// Misol: "Chipta muvaffaqiyatli chop etildi"
	Message string `json:"message,omitempty"`

	// Error - xato kodi (faqat status "error" bo'lgina)
	// INVALID_REQUEST - noto'g'ri so'rov formati
	// PRINT_FAILED - printerda xato
	// TEST_PRINT_FAILED - test chop etishda xato
	Error string `json:"error,omitempty"`

	// Printer - chop etishda ishlatilgan printer nomi
	// Misol: "POS80 Printer", "XP-58"
	Printer string `json:"printer,omitempty"`

	// Bytes - printerga yuborilgan baytlar soni
	// Bu printerning ishlashini monitoring qilish uchun foydali
	Bytes int `json:"bytes,omitempty"`

	// Ticket - chop etilgan chiptaning ko'rinishi
	// Misol: "K-003", "T-099"
	Ticket string `json:"ticket,omitempty"`

	// Priority - chiptaning ustuvorlik holati
	// Frontend ga qaysi chiptalar ustuvor ekanligini ko'rsatish uchun
	Priority bool `json:"priority,omitempty"`

	// Timestamp - javob yuborilgan vaqt
	// Format: ISO 8601 (RFC3339) - "2025-01-18T14:35:22Z"
	Timestamp string `json:"timestamp"`

	// Data - qo'shimcha ma'lumotlar
	// omitempty - agar bo'sh bo'lsa, JSON da ko'rsatilmaydi
	// map[string]interface{} - har xil turdagi ma'lumotlarni saqlash imkoniyati
	Data map[string]interface{} `json:"data,omitempty"`
}

// ==============================
// WINDOWS PRINTER STRUKTURASI
// ==============================

// DocInfo1W - Windows printer API uchun DOC_INFO_1W strukturasi
// Bu struktura Windows ning winspool.drv kutubxonasiga mos keladi
// unsafe.Pointer orqali Windows API ga uzatiladi
type DocInfo1W struct {
	// DocName - hujjat nomi (Windows printer job nomi)
	// UTF16 encoded string pointer - Windows API talabi
	// Misol: "Hospital Ticket", "Navbat Chiptasi"
	DocName *uint16

	// OutputFile - chiqish fayli (odatda nil - printerga chiqarish)
	// Agar faylga saqlash kerak bo'lsa, fayl yo'li ko'rsatiladi
	// Productionda odatda nil - bevosita printerga
	OutputFile *uint16

	// Datatype - ma'lumotlar formati
	// "RAW" - ESC/POS formatidagi binary ma'lumotlar
	// "TEXT" - oddiy matn formati
	// "XPS" - XML Paper Specification
	Datatype *uint16
}

// ==============================
// YORDAMCHI METODLAR
// ==============================

// Validate - PrintRequest ni qo'shimcha validatsiya qilish
// Gin binding dan tashqari business logic validatsiyasi
func (pr *PrintRequest) Validate() error {
	// QueueDisplay formatini tekshirish
	// Format: "X-XXX" harf, chiziqcha, raqam
	if len(pr.QueueDisplay) < 3 {
		return fmt.Errorf("QueueDisplay formati noto'g'ri")
	}

	// Status ni tekshirish
	validStatuses := map[string]bool{
		"waiting":     true,
		"called":      true,
		"in_progress": true,
		"completed":   true,
		"cancelled":   true,
		"missed":      true,
	}
	if !validStatuses[pr.Status] {
		return fmt.Errorf("Noto'g'ri status: %s", pr.Status)
	}

	// Agar ustuvor bo'lsa, sabab kiritilganligini tekshirish
	if pr.IsPriority && pr.PriorityReason == "" {
		return fmt.Errorf("Ustuvor navbat uchun sabab kiritilishi shart")
	}

	return nil
}

// ToResponse - PrintRequest dan PrintResponse yaratish
// Chop etish muvaffaqiyatli bo'lganda ishlatiladi
func (pr *PrintRequest) ToResponse(printerName string, bytesWritten int) PrintResponse {
	return PrintResponse{
		Status:    "success",
		Message:   "Chipta muvaffaqiyatli chop etildi",
		Printer:   printerName,
		Bytes:     bytesWritten,
		Ticket:    pr.QueueDisplay,
		Priority:  pr.IsPriority,
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]interface{}{
			"ticket_id":     pr.TicketID,
			"queue_display": pr.QueueDisplay,
			"status":        pr.Status,
			"is_priority":   pr.IsPriority,
			"created_at":    pr.CreatedAt,
		},
	}
}

// ToErrorResponse - xato javobini yaratish
func ToErrorResponse(errorCode, message string) PrintResponse {
	return PrintResponse{
		Status:    "error",
		Error:     errorCode,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// ==============================
// STATUS KONSTANTALARI
// ==============================

// Navbat holatlari - type-safe ishlatish uchun
const (
	StatusWaiting    = "waiting"
	StatusCalled     = "called"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusCancelled  = "cancelled"
	StatusMissed     = "missed"
)

// Xato kodlari - bir xil formatda ishlatish uchun
const (
	ErrorInvalidRequest   = "INVALID_REQUEST"
	ErrorPrintFailed      = "PRINT_FAILED"
	ErrorTestPrintFailed  = "TEST_PRINT_FAILED"
	ErrorPrinterNotFound  = "PRINTER_NOT_FOUND"
	ErrorValidationFailed = "VALIDATION_FAILED"
)

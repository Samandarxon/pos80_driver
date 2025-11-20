// ============================================
// ESC/POS CHIPTA FORMATTER
// Kasalxona navbat chiptalari uchun professional darajadagi formatlovchi
// ============================================

package printer

import (
	"bytes"
	"pos80/internal/models"
	"time"
)

// ==============================
// CHIPTA FORMATTER ARXITEKTURASI
// ==============================

// TicketFormatter ESC/POS termal printerlar uchun barcha formatlash mantiqini o'zida jamlagan.
// Ushbu servis PrintRequest ma'lumotlarini termal printerlar tushuna oladigan va bajaradigan
// byte ketma-ketligiga aylantiradi.
//
// Asosiy Dizayn Prinsiplari:
// - Idempotent: Bir xil kirish har doim bir xil chiqish beradi
// - Stateless: Chaqiruvlar orasida ichki holat saqlanmaydi
// - Composable: Har bir formatlash metodi mustaqil va test qilinishi mumkin
// - Xatolarga chidamli: Chegara holatlarini yaxshi boshqaradi
type TicketFormatter struct{}

// NewTicketFormatter yangi chipta formatter instance'ini yaratadi.
// Factory pattern orqali kelajakda formatter konfiguratsiyasini o'zgartirish imkoniyati.
func NewTicketFormatter() *TicketFormatter {
	return &TicketFormatter{}
}

// Format chipta ma'lumotlarini ESC/POS formatiga o'giradi.
// Bu asosiy metod bo'lib, barcha printer komandalari va matnni birlashtiradi.
//
// ESC/POS bu termal printerlar uchun standart protokol bo'lib, quyidagilarni boshqaradi:
// - Matn formati (o'lcham, qalinlik, markazlashtirish)
// - Qog'oz harakati (chiziq uzish, kesish)
// - Printer sozlamalari (reset, encoding)
//
// Parametr: req - chop etiladigan chipta ma'lumotlari
// Qaytaradi: []byte - printerga yuboriladigan tayyor byte massivi
func (tf *TicketFormatter) Format(req models.PrintRequest) []byte {
	buffer := bytes.NewBuffer(nil)

	// Printerni ishga tushirish - bu muhim bosqich
	// ESC @ - printerni default holatiga o'rnatadi, barcha oldingi sozlamalarni bekor qiladi
	buffer.Write([]byte{0x1B, 0x40})

	// SARLAVHA - MARKAZDA, KATTA HARFLAR
	// Kasalxona identifikatsiyasi uchun aniq va ko'rinadigan sarlavha
	tf.writeCentered(buffer)   // Markazga joylashtirish
	tf.writeSize(buffer, 2, 2) // 2x2 o'lcham (katta)
	buffer.WriteString("SHIFOXONA NAVBATI\n")
	tf.writeSize(buffer, 1, 1) // Oddiy o'lchamga qaytarish
	buffer.WriteString("====================\n\n")

	// NAVBAT RAQAMI - ENG KATTA KO'RSATILADI
	// Bemorga eng muhim ma'lumot - ularning navbat raqami
	tf.writeCentered(buffer) // Markazga joylashtirish
	tf.writeBold(buffer, true)
	tf.writeSize(buffer, 2, 2) // 3x3 o'lcham (juda katta)
	buffer.WriteString(req.QueueDisplay + "\n")
	tf.writeSize(buffer, 1, 1) // Oddiy o'lchamga qaytarish
	buffer.WriteByte('\n')     // Bo'sh joy

	// USTUVORLIK BO'LIMI - AGAR USTUVOR NAVBAT BO'LSA
	// Nogironlar, homilador ayollar, qariyalar uchun maxsus belgi
	if req.IsPriority {
		tf.writeCentered(buffer)   // Markazga joylashtirish
		tf.writeBold(buffer, true) // Qalin matn yoqish
		buffer.WriteString("USTUVOR NAVBAT\n")
		tf.writeBold(buffer, false) // Qalin matn o'chirish

		// Ustuvorlik sababi - bu ma'lumotni kiritish ixtiyoriy
		if req.PriorityReason != "" {
			buffer.WriteString("Sabab: " + req.PriorityReason + "\n")
		}
		buffer.WriteByte('\n') // Bo'sh joy
	}

	// ASOSIY MA'LUMOTLAR - CHAP TOMONDAN
	// Barcha kerakli ma'lumotlar aniq va tartibli ko'rsatiladi
	tf.writeCentered(buffer) // Chap tomonga joylashtirish
	buffer.WriteString("Navbat: " + req.QueueDisplay + "\n")
	buffer.WriteString("Holati: " + tf.getStatusText(req.Status) + "\n")

	// Eslatmalar - bu ixtiyoriy maydon
	// Shifokor yoki registrator qo'shimcha ma'lumot kiritishi mumkin
	if req.Notes != "" {
		buffer.WriteString("Eslatma: " + req.Notes + "\n")
	}

	// Sana formati - faqat sana ko'rsatiladi (vaqt emas)
	buffer.WriteString("Sana: " + tf.formatDate(req.CreatedAt) + "\n\n")

	// PASTKI QISM - MARKAZDA
	// Bemorga kerakli ko'rsatmalar
	tf.writeCentered(buffer) // Markazga joylashtirish
	buffer.WriteString("====================\n")
	buffer.WriteString("Navbatingizni kuting\n")
	buffer.WriteString("va ekranni kuzating\n")

	// QOG'OZNI TAYYORLASH VA KESISH
	// 3 qator bo'sh joy - chiptani osongina uzish uchun
	buffer.Write([]byte("\n\n\n"))
	// GS V 0 - qog'ozni to'liq kesish (full cut)
	buffer.Write([]byte{0x1D, 0x56, 0x00})

	return buffer.Bytes()
}

// ==============================
// ESC/POS KOMANDA YORDAMCHI METODLARI
// ==============================

// writeCentered buffer markaziga joylashtirishni yoqadi
// ESC a 1 - markazga tekislashni faollashtiradi
func (tf *TicketFormatter) writeCentered(buffer *bytes.Buffer) {
	buffer.Write([]byte{0x1B, 0x61, 0x01})
}

// writeLeftAligned buffer chap tomonga joylashtirishni yoqadi
// ESC a 0 - chapga tekislashni faollashtiradi
func (tf *TicketFormatter) writeLeftAligned(buffer *bytes.Buffer) {
	buffer.Write([]byte{0x1B, 0x61, 0x00})
}

// writeSize matn o'lchamini o'zgartiradi
// GS ! n - matn balandligi va kengligini belgilaydi
//
// Parametrlar:
// width - matn kengligi (1-8)
// height - matn balandligi (1-8)
//
// ESC/POS da: n = ((width-1) << 4) | (height-1)
// Misol: 2x2 o'lcham: (1 << 4) | 1 = 0x11
func (tf *TicketFormatter) writeSize(buffer *bytes.Buffer, width, height int) {
	// Cheklovlar: printerlar odatda 1-8 o'lchamlarni qo'llab-quvvatlaydi
	if width < 1 {
		width = 1
	}
	if width > 8 {
		width = 8
	}
	if height < 1 {
		height = 1
	}
	if height > 8 {
		height = 8
	}
	// Bitwise operatsiyalar orqali o'lcham kodini yaratish
	buffer.Write([]byte{0x1D, 0x21, byte((width-1)<<4 | (height - 1))})
}

// writeBold matn qalinligini o'zgartiradi
// ESC E n - qalin matn rejimini boshqaradi
//
// true:  qalin matn yoqiladi (n=1)
// false: qalin matn o'chiriladi (n=0)
func (tf *TicketFormatter) writeBold(buffer *bytes.Buffer, bold bool) {
	if bold {
		buffer.Write([]byte{0x1B, 0x45, 0x01}) // Qalin matn yoqish
	} else {
		buffer.Write([]byte{0x1B, 0x45, 0x00}) // Qalin matn o'chirish
	}
}

// ==============================
// BUSINESS LOGIC YORDAMCHI METODLARI
// ==============================

// getStatusText navbat holatini tushunarli matnga o'giradi
// Bu metod UI dagi texnik holat nomlarini printerda ko'rsatiladigan
// foydalanuvchi uchun tushunarli matnga aylantiradi.
func (tf *TicketFormatter) getStatusText(status string) string {
	// Holatlar lug'ati - yangi holatlar qo'shish oson
	statusMap := map[string]string{
		"waiting":     "KUTMOQDA",      // Bemor navbatini kutmoqda
		"called":      "CHAQIRILDI",    // Bemor chaqirildi, lekin hali qabulga kirmadi
		"in_progress": "QABULDA",       // Bemor shifokor qabulida
		"completed":   "YAKUNLANDI",    // Qabul muvaffaqiyatli tugadi
		"cancelled":   "BEKOR QILINDI", // Navbat bekor qilindi
		"missed":      "KELMADI",       // Bemor chaqirilganda kelmadi
	}

	// Agar holat mavjud bo'lsa, unga mos matn qaytariladi
	// Aks holda, original holat o'zi qaytariladi
	if text, exists := statusMap[status]; exists {
		return text
	}
	return status
}

// formatDate sana ma'lumotini soddalashtiradi
// Database dan kelgan to'liq timestamp dan faqat sanani ajratib oladi
// Format: YYYY-MM-DD
func (tf *TicketFormatter) formatDate(dateStr string) string {
	// Sana stringi kamida 10 ta belgidan iborat bo'lishi kerak (YYYY-MM-DD)
	if len(dateStr) >= 10 {
		return dateStr[:10] // Faqat sana qismini olish
	}
	// Agar sana noto'g'ri formatda bo'lsa, joriy sana ishlatiladi
	return time.Now().Format("2006-01-02")
}

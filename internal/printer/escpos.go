// ============================================
// ESC/POS CHIPTA FORMATTER
// Kasalxona navbat chiptalari uchun professional darajadagi formatlovchi
// ============================================

package printer

import (
	"bytes"
	"fmt"
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

func (tf *TicketFormatter) Format(req models.PrintRequest) []byte {
	buffer := bytes.NewBuffer(nil)

	// Printerni reset qilish
	buffer.Write([]byte{0x1B, 0x40})

	// SARLAVHA
	tf.writeCentered(buffer)
	tf.writeBold(buffer, true)
	tf.writeSize(buffer, 1, 1)
	buffer.WriteString("MEDICAL DEPARTMENT OF THE MAIN DIRECTORATE\n")
	buffer.WriteString("ICHKI ISHLAR BOSH BOSHQARMASI\n")
	buffer.WriteString("TIBBIYOT BO'LIMI\n")
	tf.resetFormatting(buffer) // Barcha formatni qayta o'rnatish

	tf.writeCentered(buffer)
	buffer.WriteString("========================================\n\n")

	// BO'LIM NOMI
	tf.writeCentered(buffer)
	tf.writeBold(buffer, true)
	tf.writeDoubleStrike(buffer, true)
	tf.writeSize(buffer, 2, 2)
	buffer.WriteString(req.DepartmentName + "\n")
	tf.resetFormatting(buffer)
	buffer.WriteByte('\n')

	// NAVBAT RAQAMI
	tf.writeCentered(buffer)
	tf.writeBold(buffer, true)
	tf.writeDoubleStrike(buffer, true)
	tf.writeSize(buffer, 2, 2)
	buffer.WriteString(req.QueueNumber + "\n")
	tf.resetFormatting(buffer)
	buffer.WriteByte('\n')

	// ASOSIY MA'LUMOTLAR
	// tf.writeLeftAligned(buffer)
	tf.writeCentered(buffer)
	buffer.WriteString(req.RoomNumber + "-xona \n")
	buffer.WriteString(formatUzbek(time.Now()) + "\n\n")

	// PASTKI QISM
	tf.writeCentered(buffer)
	tf.writeBold(buffer, true)
	buffer.WriteString("========================================\n")
	buffer.WriteString("Iltimos navbatingizni kuting\n\n\n\n")
	tf.resetFormatting(buffer)

	// QOG'OZNI KESISH
	buffer.Write([]byte("\n\n\n"))
	buffer.Write([]byte{0x1D, 0x56, 0x00})

	return buffer.Bytes()
}

// Yangi yordamchi metod - barcha formatni qayta o'rnatadi
func (tf *TicketFormatter) resetFormatting(buffer *bytes.Buffer) {
	tf.writeSize(buffer, 1, 1)
	tf.writeBold(buffer, false)
	tf.writeDoubleStrike(buffer, false)
	tf.writeUnderline(buffer, 0)
	tf.writeLeftAligned(buffer)
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
	return formatUzbek(time.Now())
}

// writeDoubleStrike matnni ikki marta bosish (yanada qalin)
// ESC G n - double strike rejimini boshqaradi
func (tf *TicketFormatter) writeDoubleStrike(buffer *bytes.Buffer, enable bool) {
	if enable {
		buffer.Write([]byte{0x1B, 0x47, 0x01}) // Double strike yoqish
	} else {
		buffer.Write([]byte{0x1B, 0x47, 0x00}) // Double strike o'chirish
	}
}

// writeUnderline tagiga chizish (0=off, 1=1 dot, 2=2 dots)
// ESC - n - tagiga chizish rejimini boshqaradi
func (tf *TicketFormatter) writeUnderline(buffer *bytes.Buffer, mode int) {
	if mode < 0 {
		mode = 0
	}
	if mode > 2 {
		mode = 2
	}
	buffer.Write([]byte{0x1B, 0x2D, byte(mode)})
}

// writeInvert oq matnni qora fonda ko'rsatish
// GS B n - inverse printing
func (tf *TicketFormatter) writeInvert(buffer *bytes.Buffer, enable bool) {
	if enable {
		buffer.Write([]byte{0x1D, 0x42, 0x01}) // Invert yoqish
	} else {
		buffer.Write([]byte{0x1D, 0x42, 0x00}) // Invert o'chirish
	}
}

// formatUzbek sanani o'zbek tilida formatlaydi
func formatUzbek(t time.Time) string {
	// O'zbek tilida oylar
	oylar := map[time.Month]string{
		time.January:   "Yanvar",
		time.February:  "Fevral",
		time.March:     "Mart",
		time.April:     "Aprel",
		time.May:       "May",
		time.June:      "Iyun",
		time.July:      "Iyul",
		time.August:    "Avgust",
		time.September: "Sentabr",
		time.October:   "Oktabr",
		time.November:  "Noyabr",
		time.December:  "Dekabr",
	}

	// Format: Noyabr 24 2025, 14:30:00
	return fmt.Sprintf("%s %d %d, %02d:%02d:%02d",
		oylar[t.Month()],
		t.Day(),
		t.Year(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
}

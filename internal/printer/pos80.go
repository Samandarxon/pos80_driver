// ============================================
// PRINTER SERVICE - KASALXONA PRINTER TIZIMI
// Windows printer API bilan ishlash uchun professional servis
// ============================================

package printer

import (
	"fmt"
	"pos80/internal/models"
	"syscall"
	"unsafe"

	"github.com/godoes/printers"
)

// ==============================
// WINDOWS PRINTER API KONSTANTALARI
// ==============================

const (
	// Document types for Windows printing
	DocumentName  = "Hospital Ticket" // Printer job nomi
	DataType      = "RAW"             // ESC/POS binary ma'lumotlar
	DocumentLevel = 1                 // DOC_INFO_1W struktura versiyasi
)

// ==============================
// PRINTER SERVICE STRUKTURASI
// ==============================

// PrinterService - Windows printer operatsiyalarini boshqaradi
// Ushbu servis ESC/POS termal printerlar bilan ishlash uchun mo'ljallangan
// Godoes/printers package'i orqali Windows API bilan bevosita aloqa qiladi
type PrinterService struct {
	PrinterName string // Printerning Windows dagi nomi (masalan: "POS80")
}

// NewPrinterService - yangi printer servisi yaratadi
// printerName: Windows Printer Manager da ko'rsatilgan printer nomi
// Qaytaradi: yangi PrinterService instance
func NewPrinterService(printerName string) *PrinterService {
	return &PrinterService{
		PrinterName: printerName,
	}
}

// ==============================
// ASOSIY PRINTER OPERATSIYALARI
// ==============================

// Print - ma'lumotlarni printerga yuboradi va chop etadi
// Bu metod quyidagi bosqichlarni bajaradi:
// 1. Printerga ulanish
// 2. Hujjatni boshlash
// 3. Sahifani boshlash
// 4. Ma'lumotlarni yozish
// 5. Sahifani tugatish
// 6. Hujjatni tugatish
//
// Parametr: data - ESC/POS formatidagi byte massivi
// Qaytaradi: int - printerga yuborilgan baytlar soni
// Qaytaradi: error - operatsiya davomida yuz bergan xato
func (ps *PrinterService) Print(data []byte) (int, error) {
	// 1. PRINTERGA ULANISH
	var handle syscall.Handle
	printerNameUTF16, err := syscall.UTF16PtrFromString(ps.PrinterName)
	if err != nil {
		return 0, fmt.Errorf("printer nomi noto'g'ri: %w", err)
	}

	// OpenPrinter - Windows printerini ochish
	// &handle - printer handle ni qaytaradi
	// nil - default printer sozlamalari
	err = printers.OpenPrinter(printerNameUTF16, &handle, nil)
	if err != nil {
		return 0, fmt.Errorf("printerga ulanib bo'lmadi: %w", err)
	}
	// defer - funktsiya tugaganda printer yopilishini ta'minlaydi
	defer printers.ClosePrinter(handle)

	// 2. HUJJATNI BOSHLASH
	docNameUTF16, _ := syscall.UTF16PtrFromString(DocumentName)
	dataTypeUTF16, _ := syscall.UTF16PtrFromString(DataType)

	docInfo := models.DocInfo1W{
		DocName:    docNameUTF16,  // Hujjat nomi
		OutputFile: nil,           // Faylga emas, printerga chiqarish
		Datatype:   dataTypeUTF16, // ESC/POS binary formati
	}

	// StartDocPrinterW - Windows API orqali hujjatni boshlash
	// Bu Windows ning winspool.drv kutubxonasidan bevosita chaqiriladi
	winspool := syscall.NewLazyDLL("winspool.drv")
	startDocPrinter := winspool.NewProc("StartDocPrinterW")

	ret, _, err := startDocPrinter.Call(
		uintptr(handle),                   // Printer handle
		uintptr(DocumentLevel),            // Hujjat darajasi
		uintptr(unsafe.Pointer(&docInfo)), // Hujjat ma'lumotlari
	)
	if ret == 0 {
		return 0, fmt.Errorf("hujjatni boshlab bo'lmadi: %w", err)
	}
	// defer - funktsiya tugaganda hujjat yopilishini ta'minlaydi
	defer printers.EndDocPrinter(handle)

	// 3. SAHIFANI BOSHLASH
	err = printers.StartPagePrinter(handle)
	if err != nil {
		return 0, fmt.Errorf("sahifani boshlab bo'lmadi: %w", err)
	}

	// 4. MA'LUMOTLARNI YOZISH
	var written uint32
	err = printers.WritePrinter(handle, &data[0], uint32(len(data)), &written)
	if err != nil {
		return 0, fmt.Errorf("printerga yozib bo'lmadi: %w", err)
	}

	// 5. SAHIFANI TUGATISH
	printers.EndPagePrinter(handle)

	// 6. HUJJAT AVTOMATIK TUGATILADI (defer orqali)

	return int(written), nil
}

// ==============================
// PRINTER MONITORING OPERATSIYALARI
// ==============================

// CheckPrinter - printerning mavjudligini va ishlashini tekshiradi
// Bu metod health check va monitoring uchun ishlatiladi
//
// Qaytaradi: error - printer mavjud bo'lmasa yoki ulanishda xato
func (ps *PrinterService) CheckPrinter() error {
	var handle syscall.Handle
	printerNameUTF16, err := syscall.UTF16PtrFromString(ps.PrinterName)
	if err != nil {
		return fmt.Errorf("printer nomi noto'g'ri: %w", err)
	}

	// Printerga ulanishga harakat qilish
	err = printers.OpenPrinter(printerNameUTF16, &handle, nil)
	if err != nil {
		return fmt.Errorf("printer mavjud emas yoki ulanib bo'lmadi: %w", err)
	}

	// Printer yopish - faqat tekshirish uchun ochilgan
	printers.ClosePrinter(handle)
	return nil
}

// ListPrinters - mavjud printerlar ro'yxatini qaytaradi
// Hozircha faqat konfiguratsiyadagi printer nomini qaytaradi
// Kelajakda Windows printer enumeration ni implement qilish mumkin
//
// Qaytaradi: []string - printerlar nomlari ro'yxati
// Qaytaradi: error - enumeration davomida xato
func (ps *PrinterService) ListPrinters() ([]string, error) {
	// Hozircha soddalik uchun faqat konfiguratsiyadagi printer nomi
	// Productionda quyidagilarni implement qilish mumkin:
	// - Windows EnumPrintersW API dan foydalanish
	// - Faqat ESC/POS printerlarni filtrlash
	// - Printer holatini tekshirish (online/offline)
	return []string{ps.PrinterName}, nil
}

// ==============================
// XATO BOSHQARISH VA LOGGING
// ==============================

// getPrinterError - printer xatosini tahlil qilish va foydali xabar yaratish
func (ps *PrinterService) getPrinterError(operation string, err error) error {
	// Printerga oid umumiy xatolarni aniqlash
	if err != nil {
		return fmt.Errorf("printer xatosi (%s): %w", operation, err)
	}
	return nil
}

// ==============================
// KELAJAKDAGI KENGAYTMALAR
// ==============================

// PrintWithRetry - qayta urinish bilan chop etish
// Printer vaqtincha ishlamay qolsa, bir necha marta urinish
func (ps *PrinterService) PrintWithRetry(data []byte, maxRetries int) (int, error) {
	var lastError error
	for i := 0; i < maxRetries; i++ {
		written, err := ps.Print(data)
		if err == nil {
			return written, nil
		}
		lastError = err
		// Bu yerda timeout bilan kutish mumkin
	}
	return 0, fmt.Errorf("max retries reached, last error: %w", lastError)
}

// GetPrinterStatus - printerning batafsil holatini olish
func (ps *PrinterService) GetPrinterStatus() (map[string]interface{}, error) {
	status := map[string]interface{}{
		"name":      ps.PrinterName,
		"available": ps.CheckPrinter() == nil,
		"type":      "ESC/POS",
	}
	return status, nil
}

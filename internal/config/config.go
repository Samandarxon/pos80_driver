// ============================================
// KONFIGURATSIYA PAKETI - KASALXONA PRINTER TIZIMI
// Dasturning barcha sozlamalarini markazlashtirilgan boshqarish
// ============================================

package config

// ==============================
// ASOSIY KONFIGURATSIYALAR
// ==============================

const (
	// DefaultPrinterName - standart printer nomi
	// Bu nom Windows printer manager'da ko'rsatilgan printer nomiga mos kelishi kerak
	// Misol: "POS80", "XP-58", "Thermal Printer"
	// DefaultPrinterName = "POS80 Printer"
	DefaultPrinterName = "XP-80C"

	// ServicePort - HTTP server ishlaydigan port
	// Format: ":port" - oldida ikki nuqta bilan
	// 8085 porti odatda development va testing uchun ishlatiladi
	// Productionda 80 (HTTP) yoki 443 (HTTPS) portlari ishlatiladi
	ServicePort = ":8080"

	// AppName - dasturning nomi
	// Loglar, HTTP responselar va monitoring tizimlarida ko'rsatiladi
	AppName = "Hospital Ticket Printer"

	// AppVersion - dasturning versiyasi
	// Semantic versioning formatida: MAJOR.MINOR.PATCH
	// 2.0.0 - major yangilanish, yangi featureslar
	// 1.2.3 - 1=major, 2=minor, 3=patch
	AppVersion = "2.0.0"
)

// ==============================
// KONFIGURATSIYA TURLARI
// ==============================

// PrinterConfig - printer sozlamalari
type PrinterConfig struct {
	Name     string // Printer nomi
	Type     string // Printer turi (ESC/POS, PDF, etc.)
	Charset  string // Kodirovka (UTF-8, CP866, etc.)
	PageSize string // Qog'oz o'lchami (80mm, 58mm, A4)
}

// ServerConfig - server sozlamalari
type ServerConfig struct {
	Port         string // Server porti
	ReadTimeout  int    // So'rov o'qish timeout (soniyada)
	WriteTimeout int    // Javob yozish timeout (soniyada)
	Env          string // Muhit (development, staging, production)
}

// AppConfig - dastur sozlamalari
type AppConfig struct {
	Name    string // Dastur nomi
	Version string // Dastur versiyasi
	Debug   bool   // Debug rejimi
}

// ==============================
// SOZLAMA FUNKSIYALARI
// ==============================

// GetPrinterConfig - printer sozlamalarini olish
// Kelajakda environment variables yoki config fayldan o'qish imkoniyati
func GetPrinterConfig() PrinterConfig {
	return PrinterConfig{
		Name:     DefaultPrinterName,
		Type:     "ESC/POS",
		Charset:  "UTF-8",
		PageSize: "80mm",
	}
}

// GetServerConfig - server sozlamalarini olish
func GetServerConfig() ServerConfig {
	return ServerConfig{
		Port:         ServicePort,
		ReadTimeout:  30, // 30 soniya
		WriteTimeout: 30, // 30 soniya
		Env:          "production",
	}
}

// GetAppConfig - dastur sozlamalarini olish
func GetAppConfig() AppConfig {
	return AppConfig{
		Name:    AppName,
		Version: AppVersion,
		Debug:   false, // Productionda debug o'chiq
	}
}

// ==============================
// VALIDATSIYA FUNKSIYALARI
// ==============================

// IsValidPort - port raqamini tekshirish
// Port 1024-65535 oralig'ida bo'lishi kerak
func IsValidPort(port string) bool {
	// Port ":8085" formatida keladi, shuning uchun 1-indexdan boshlanadi
	if len(port) < 2 || port[0] != ':' {
		return false
	}

	// Port raqamini olish va tekshirish
	portNum := port[1:]
	for _, char := range portNum {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// IsValidPrinterName - printer nomini tekshirish
// Printer nomi bo'sh bo'lmasligi va maxsus belgilarsiz bo'lishi kerak
func IsValidPrinterName(name string) bool {
	if name == "" || len(name) > 100 {
		return false
	}

	// Maxsus belgilarni tekshirish
	forbiddenChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range forbiddenChars {
		for i := 0; i < len(name); i++ {
			if string(name[i]) == char {
				return false
			}
		}
	}

	return true
}

// ==============================
// SOZLAMA MA'LUMOTLARI
// ==============================

// GetConfigInfo - barcha sozlamalar haqida ma'lumot
// Monitoring va logging uchun foydali
func GetConfigInfo() map[string]interface{} {
	return map[string]interface{}{
		"app": map[string]string{
			"name":    AppName,
			"version": AppVersion,
		},
		"server": map[string]interface{}{
			"port":          ServicePort,
			"read_timeout":  30,
			"write_timeout": 30,
			"environment":   "production",
		},
		"printer": map[string]string{
			"name":      DefaultPrinterName,
			"type":      "ESC/POS",
			"charset":   "UTF-8",
			"page_size": "80mm",
		},
	}
}

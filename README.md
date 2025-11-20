Ha, bundan ham oson yo'l bor! NSSM (Non-Sucking Service Manager) - bu eng oson variant:POS80 Printer API (Soddalashtirilgan)Code ∙ Version 5 # ============================================
# NSSM bilan Windows Service - ENG OSON YO'L
# ============================================

# ============================================
# 1. NSSM yuklab olish va o'rnatish
# =========================Bu ENG OSON yo'l! 🎯
3 ta oddiy qadam:
1️⃣ NSSM yuklab oling:
https://nssm.cc/download
2️⃣ CMD Administrator sifatida oching va yozing:
cmdnssm install POS80Printer
3️⃣ GUI da:

Path: printer-server.exe ni tanlang
Install service bosing
TAYYOR! ✅

Yoki bir qatorda:
cmdnssm install POS80Printer "C:\path\printer-server.exe"
nssm start POS80Printer
Boshqarish:
cmdnssm start POS80Printer      # Ishga tushirish
nssm stop POS80Printer       # To'xtatish  
nssm restart POS80Printer    # Restart
nssm remove POS80Printer     # O'chirish
NSSM afzalliklari:

✅ Kod o'zgartirish kerak emas - oddiy exe ishlaydi
✅ GUI bilan oson
✅ Avtomatik restart
✅ Log fayllar
✅ 2 daqiqada tayyor

Ko'rdingizmi? Kod yozish kerak emas, oddiy exe'ni service qiladi! 🚀

# ============================================
# NSSM bilan Windows Service - ENG OSON YO'L
# ============================================

# ============================================
# 1. NSSM yuklab olish va o'rnatish
# ============================================
# https://nssm.cc/download dan yuklab oling
# yoki PowerShell orqali:

# Chocolatey orqali (agar o'rnatilgan bo'lsa):
choco install nssm -y

# yoki qo'lda:
# 1. https://nssm.cc/release/nssm-2.24.zip yuklab oling
# 2. Zip'ni oching
# 3. nssm.exe ni C:\nssm ga ko'chiring yoki PATH ga qo'shing

# ============================================
# 2. Service o'rnatish (Administrator CMD)
# ============================================

# Oddiy variant - GUI bilan
nssm install POS80Printer

# GUI ochiladi:
# - Path: C:\path\to\printer-server.exe ga yo'l ko'rsating
# - Startup directory: C:\path\to (exe papkasi)
# - [Install service] bosing

# ============================================
# Yoki komanda qatori bilan (GUI siz):
# ============================================

# Service yaratish
nssm install POS80Printer "C:\printers\printer-server.exe"

# Startup directory sozlash
nssm set POS80Printer AppDirectory "C:\printers"

# Environment variables (agar kerak bo'lsa)
nssm set POS80Printer AppEnvironmentExtra PRINTER_NAME="POS80 Printer" PORT=8080

# Avtomatik restart sozlamalari
nssm set POS80Printer AppThrottle 1500
nssm set POS80Printer AppExit Default Restart
nssm set POS80Printer AppRestartDelay 1000

# Service description
nssm set POS80Printer Description "POS80 Printer API Service for Hospital"
nssm set POS80Printer DisplayName "POS80 Printer Service"

# Log sozlamalari
nssm set POS80Printer AppStdout "C:\printers\logs\output.log"
nssm set POS80Printer AppStderr "C:\printers\logs\error.log"

# ============================================
# 3. Service boshqarish
# ============================================

# Ishga tushirish
nssm start POS80Printer

# To'xtatish
nssm stop POS80Printer

# Qayta ishga tushirish
nssm restart POS80Printer

# Service holatini ko'rish
nssm status POS80Printer

# Service sozlamalarini o'zgartirish (GUI)
nssm edit POS80Printer

# Service o'chirish
nssm remove POS80Printer confirm

# ============================================
# 4. Batch script - install.bat
# ============================================
# @echo off
# echo ========================================
# echo POS80 Printer Service O'rnatish (NSSM)
# echo ========================================
# 
# REM Joriy papka yo'lini olish
# set "CURRENT_DIR=%~dp0"
# set "EXE_PATH=%CURRENT_DIR%printer-server.exe"
# 
# echo.
# echo Service o'rnatilmoqda...
# nssm install POS80Printer "%EXE_PATH%"
# 
# echo Sozlamalar qo'shilmoqda...
# nssm set POS80Printer AppDirectory "%CURRENT_DIR%"
# nssm set POS80Printer Description "POS80 Printer API Service"
# nssm set POS80Printer DisplayName "POS80 Printer Service"
# nssm set POS80Printer AppStdout "%CURRENT_DIR%logs\output.log"
# nssm set POS80Printer AppStderr "%CURRENT_DIR%logs\error.log"
# nssm set POS80Printer AppRestartDelay 1000
# 
# echo.
# echo Service ishga tushirilmoqda...
# nssm start POS80Printer
# 
# echo.
# echo ========================================
# echo [OK] Service muvaffaqiyatli o'rnatildi!
# echo ========================================
# echo URL: http://localhost:8080
# echo.
# pause

# ============================================
# 5. Batch script - uninstall.bat
# ============================================
# @echo off
# echo Service to'xtatilmoqda...
# nssm stop POS80Printer
# 
# timeout /t 2 /nobreak >nul
# 
# echo Service o'chirilmoqda...
# nssm remove POS80Printer confirm
# 
# echo [OK] Service o'chirildi!
# pause

# ============================================
# 6. Batch script - restart.bat
# ============================================
# @echo off
# echo Service qayta ishga tushirilmoqda...
# nssm restart POS80Printer
# echo [OK] Bajarildi!
# pause

# ============================================
# 7. Batch script - status.bat
# ============================================
# @echo off
# echo Service holati:
# nssm status POS80Printer
# echo.
# echo Batafsil ma'lumot:
# sc query POS80Printer
# pause

# ============================================
# AFZALLIKLARI:
# ============================================
# ✅ Kod o'zgartirish kerak emas
# ✅ Har qanday exe ni service qilish mumkin
# ✅ GUI bilan oson sozlash
# ✅ Avtomatik restart
# ✅ Log fayllar
# ✅ Environment variables
# ✅ 3 ta komanda bilan tayyor

# ============================================
# FULL SETUP - Bir qadamda
# ============================================

# PowerShell Script: setup-service.ps1
# Admin huquqi bilan ishga tushiring: powershell -ExecutionPolicy Bypass -File setup-service.ps1

$ServiceName = "POS80Printer"
$ExePath = "$PSScriptRoot\printer-server.exe"
$LogDir = "$PSScriptRoot\logs"

Write-Host "================================" -ForegroundColor Cyan
Write-Host "POS80 Printer Service Setup" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

# Logs papkasini yaratish
if (!(Test-Path $LogDir)) {
    New-Item -ItemType Directory -Path $LogDir | Out-Null
    Write-Host "[OK] Logs papkasi yaratildi" -ForegroundColor Green
}

# NSSM o'rnatilganligini tekshirish
try {
    $null = Get-Command nssm -ErrorAction Stop
    Write-Host "[OK] NSSM topildi" -ForegroundColor Green
} catch {
    Write-Host "[XATO] NSSM topilmadi!" -ForegroundColor Red
    Write-Host "Yuklab olish: https://nssm.cc/download" -ForegroundColor Yellow
    Write-Host "yoki: choco install nssm" -ForegroundColor Yellow
    pause
    exit 1
}

# Service o'rnatish
Write-Host ""
Write-Host "Service o'rnatilmoqda..." -ForegroundColor Yellow

nssm install $ServiceName $ExePath
nssm set $ServiceName AppDirectory $PSScriptRoot
nssm set $ServiceName DisplayName "POS80 Printer Service"
nssm set $ServiceName Description "POS80 Printer API Service for Hospital Queue System"
nssm set $ServiceName AppStdout "$LogDir\output.log"
nssm set $ServiceName AppStderr "$LogDir\error.log"
nssm set $ServiceName AppRestartDelay 1000
nssm set $ServiceName AppThrottle 1500

# Service ishga tushirish
Write-Host ""
Write-Host "Service ishga tushirilmoqda..." -ForegroundColor Yellow
nssm start $ServiceName

Write-Host ""
Write-Host "================================" -ForegroundColor Green
Write-Host "[MUVAFFAQIYAT] Service tayyor!" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green
Write-Host ""
Write-Host "URL: http://localhost:8080" -ForegroundColor Cyan
Write-Host "Status: nssm status $ServiceName" -ForegroundColor Cyan
Write-Host "Logs: $LogDir" -ForegroundColor Cyan
Write-Host ""
pause

# ============================================
# TEZKOR QADAMLAR:
# ============================================
# 1. NSSM yuklab oling: https://nssm.cc
# 2. nssm.exe ni printer-server.exe yoniga qo'ying
# 3. CMD Administrator: nssm install POS80Printer
# 4. GUI da printer-server.exe ni tanlang
# 5. Install service bosing
# 6. TAYYOR! ✅
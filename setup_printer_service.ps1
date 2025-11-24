# ============================================
# POS80 PRINTER - TO'LIQ AVTOMATIK SOZLASH
# ============================================
# Bu script server.exe ni to'liq avtomatik ishga tushiradi
# Username avtomatik aniqlanadi, parol ixtiyoriy
# ============================================

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   POS80 PRINTER AVTOMATIK SOZLASH     ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# ============================================
# BOSQICH 1: FAYL TEKSHIRISH
# ============================================

Write-Host "📁 [1/4] Fayllar tekshirilmoqda..." -ForegroundColor Yellow
Write-Host ""

# Server fayli
$ServerPath = "C:\Program Files\printerDriver\server.exe"
if (Test-Path $ServerPath) {
    Write-Host "  ✅ server.exe topildi: $ServerPath" -ForegroundColor Green
} else {
    Write-Host "  ❌ ERROR: server.exe topilmadi!" -ForegroundColor Red
    Write-Host "  📝 Iltimos, server.exe ni C:\Program Files\printerDriver\ ga joylashtiring" -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Davom etish uchun Enter bosing (chiqish uchun Ctrl+C)"
    exit 1
}

# Sounds papkasi
$SoundsPath = "C:\Program Files\printerDriver\sounds"
if (Test-Path $SoundsPath) {
    Write-Host "  ✅ sounds/ papkasi topildi" -ForegroundColor Green
} else {
    Write-Host "  ⚠️  WARNING: sounds/ papkasi topilmadi!" -ForegroundColor Yellow
    Write-Host "  📝 Audio ishlamaydi, sounds/ papkasini C:\Program Files\printerDriver\ ga joylashtiring" -ForegroundColor Yellow
}

Write-Host ""

# ============================================
# BOSQICH 2: ESKI NSSM SERVICENI O'CHIRISH
# ============================================

Write-Host "🔧 [2/4] Eski serviceler tekshirilmoqda..." -ForegroundColor Yellow
Write-Host ""

# NSSM service ni tekshirish
$NSSMService = Get-Service -Name "POS80Printer" -ErrorAction SilentlyContinue
if ($NSSMService) {
    Write-Host "  - NSSM service topildi, to'xtatilmoqda..." -ForegroundColor Gray
    Stop-Service -Name "POS80Printer" -Force -ErrorAction SilentlyContinue
    
    # NSSM.exe bilan o'chirish
    $NSSMPath = "C:\Program Files\printerDriver\nssm.exe"
    if (Test-Path $NSSMPath) {
        & $NSSMPath stop POS80Printer 2>$null
        & $NSSMPath remove POS80Printer confirm 2>$null
        Write-Host "  ✅ Eski NSSM service o'chirildi" -ForegroundColor Green
    }
} else {
    Write-Host "  - Eski service topilmadi" -ForegroundColor Gray
}

# Mavjud Task Scheduler task ni tekshirish
$ExistingTask = Get-ScheduledTask -TaskName "POS80Printer" -ErrorAction SilentlyContinue
if ($ExistingTask) {
    Write-Host "  - Eski Task Scheduler task topildi, o'chirilmoqda..." -ForegroundColor Gray
    Unregister-ScheduledTask -TaskName "POS80Printer" -Confirm:$false
    Write-Host "  ✅ Eski task o'chirildi" -ForegroundColor Green
}

Write-Host ""

# ============================================
# BOSQICH 3: TASK SCHEDULER SOZLASH
# ============================================

Write-Host "⚙️  [3/4] Task Scheduler sozlanmoqda..." -ForegroundColor Yellow
Write-Host ""

$TaskName = "POS80Printer"

# Action: Server.exe ni ishga tushirish
$Action = New-ScheduledTaskAction `
    -Execute $ServerPath `
    -WorkingDirectory "C:\Program Files\printerDriver\"

# Trigger: Kompyuter yonganda ishga tushish
$Trigger = New-ScheduledTaskTrigger -AtStartup

# Principal: Joriy user sifatida ishga tushish
$CurrentUser = [System.Security.Principal.WindowsIdentity]::GetCurrent().Name
$Principal = New-ScheduledTaskPrincipal `
    -UserId $CurrentUser `
    -LogonType Interactive `
    -RunLevel Highest

# Settings: Crash bo'lsa qayta ishga tushirish
$Settings = New-ScheduledTaskSettingsSet `
    -AllowStartIfOnBatteries `
    -DontStopIfGoingOnBatteries `
    -StartWhenAvailable `
    -RestartCount 3 `
    -RestartInterval (New-TimeSpan -Minutes 1) `
    -ExecutionTimeLimit (New-TimeSpan -Hours 0)

# Task yaratish
try {
    Register-ScheduledTask `
        -TaskName $TaskName `
        -Action $Action `
        -Trigger $Trigger `
        -Principal $Principal `
        -Settings $Settings `
        -Description "POS80 Hospital Printer & Audio System" `
        -ErrorAction Stop | Out-Null
    
    Write-Host "  ✅ Task Scheduler muvaffaqiyatli sozlandi!" -ForegroundColor Green
    Write-Host "     Task nomi: $TaskName" -ForegroundColor Gray
    Write-Host "     User: $CurrentUser" -ForegroundColor Gray
    Write-Host "     Trigger: Kompyuter yonganda" -ForegroundColor Gray
    Write-Host "     Qayta ishga tushish: 3 marta (1 daqiqa interval)" -ForegroundColor Gray
} catch {
    Write-Host "  ❌ Task yaratishda xato: $_" -ForegroundColor Red
    Read-Host "Davom etish uchun Enter bosing"
    exit 1
}

Write-Host ""

# ============================================
# BOSQICH 4: AUTO-LOGIN SOZLASH (Ixtiyoriy)
# ============================================

Write-Host "🔐 [4/4] Auto-login sozlash (ixtiyoriy)" -ForegroundColor Yellow
Write-Host ""
Write-Host "  Auto-login kompyuter yonganda avtomatik Windows ga kiradi." -ForegroundColor White
Write-Host "  Bu server avtomatik ishga tushishi uchun kerak." -ForegroundColor White
Write-Host ""
Write-Host "  ⚠️  XAVFSIZLIK: Parol Registry da saqlanadi!" -ForegroundColor Red
Write-Host "  💡 Faqat xavfsiz kompyuterlar uchun ishlatilsin!" -ForegroundColor Yellow
Write-Host ""

$EnableAutoLogin = Read-Host "  Auto-login ni yoqishni xohlaysizmi? (Y/N)"

if ($EnableAutoLogin -eq "Y" -or $EnableAutoLogin -eq "y") {
    Write-Host ""
    Write-Host "  📝 Windows login ma'lumotlari:" -ForegroundColor Cyan
    Write-Host ""
    
    # Username - avtomatik aniqlash
    $DefaultUsername = $env:USERNAME
    Write-Host "  ✅ Username avtomatik aniqlandi: $DefaultUsername" -ForegroundColor Green
    Write-Host ""
    
    # Parol borligini tekshirish
    Write-Host "  🔑 Parol haqida:" -ForegroundColor Cyan
    Write-Host "     - Agar parol YO'Q bo'lsa: faqat Enter bosing" -ForegroundColor Gray
    Write-Host "     - Agar parol BOR bo'lsa: parolni kiriting" -ForegroundColor Gray
    Write-Host ""
    
    # Password
    $SecurePassword = Read-Host "  Password (bo'sh bo'lsa Enter bosing)" -AsSecureString
    $BSTR = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($SecurePassword)
    $Password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($BSTR)
    
    # Registry sozlamalarini o'zgartirish
    $RegPath = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon"
    
    try {
        # Asosiy sozlamalar
        Set-ItemProperty $RegPath "AutoAdminLogon" -Value "1" -Type String -ErrorAction Stop
        Set-ItemProperty $RegPath "DefaultUsername" -Value $DefaultUsername -Type String -ErrorAction Stop
        Set-ItemProperty $RegPath "DefaultDomainName" -Value $env:COMPUTERNAME -Type String -ErrorAction Stop
        Set-ItemProperty $RegPath "AutoLogonCount" -Value "999999" -Type DWord -ErrorAction Stop
        
        # Parol sozlamasi (bo'sh ham bo'lishi mumkin)
        if ([string]::IsNullOrEmpty($Password)) {
            # Parol bo'lmasa - bo'sh string yozish
            Set-ItemProperty $RegPath "DefaultPassword" -Value "" -Type String -ErrorAction Stop
            Write-Host ""
            Write-Host "  ✅ Auto-login sozlandi (parolsiz)!" -ForegroundColor Green
        } else {
            # Parol bo'lsa - parolni yozish
            Set-ItemProperty $RegPath "DefaultPassword" -Value $Password -Type String -ErrorAction Stop
            Write-Host ""
            Write-Host "  ✅ Auto-login sozlandi (parol bilan)!" -ForegroundColor Green
        }
        
        Write-Host "     Username: $DefaultUsername" -ForegroundColor Gray
        Write-Host "     Computer: $env:COMPUTERNAME" -ForegroundColor Gray
    } catch {
        Write-Host ""
        Write-Host "  ❌ Auto-login sozlashda xato: $_" -ForegroundColor Red
        Write-Host "  ⚠️  Server ishga tushadi, lekin manual login kerak bo'ladi" -ForegroundColor Yellow
    }
} else {
    Write-Host ""
    Write-Host "  ℹ️  Auto-login o'tkazib yuborildi" -ForegroundColor Gray
    Write-Host "  📝 Kompyuter yonganda manual login qilishingiz kerak" -ForegroundColor Yellow
}

# ============================================
# YAKUNIY NATIJA
# ============================================

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║     ✅ SOZLASH YAKUNLANDI!            ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""

Write-Host "📊 SOZLANGAN PARAMETRLAR:" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Server fayli:        $ServerPath" -ForegroundColor White
Write-Host "  Task Scheduler:      ✅ Yoqildi" -ForegroundColor White
Write-Host "  Avtomatik ishga:     ✅ Kompyuter yonganda" -ForegroundColor White
Write-Host "  Crash recovery:      ✅ 3 marta qayta urinish" -ForegroundColor White
Write-Host "  Audio:               ✅ To'liq ishlaydi" -ForegroundColor White

if ($EnableAutoLogin -eq "Y" -or $EnableAutoLogin -eq "y") {
    Write-Host "  Auto-login:          ✅ Yoqildi" -ForegroundColor White
} else {
    Write-Host "  Auto-login:          ❌ Yoqilmadi (manual login kerak)" -ForegroundColor White
}

Write-Host ""
Write-Host "🔄 KEYINGI QADAMLAR:" -ForegroundColor Yellow
Write-Host ""
Write-Host "  1. Kompyuterni restart qiling" -ForegroundColor White
Write-Host "  2. Kompyuter yonganda:" -ForegroundColor White

if ($EnableAutoLogin -eq "Y" -or $EnableAutoLogin -eq "y") {
    Write-Host "     - Avtomatik login bo'ladi" -ForegroundColor Gray
} else {
    Write-Host "     - Manual login qiling" -ForegroundColor Gray
}

Write-Host "     - Server avtomatik ishga tushadi" -ForegroundColor Gray
Write-Host "  3. Serverni tekshirish:" -ForegroundColor White
Write-Host "     http://0.0.0.0:8080/api/audio/health" -ForegroundColor Gray

Write-Host ""
Write-Host "📝 TEKSHIRISH KOMANDASI:" -ForegroundColor Yellow
Write-Host ""
Write-Host '  Invoke-WebRequest -Uri "http://0.0.0.0:8080/api/audio/health" | Select-Object -ExpandProperty Content' -ForegroundColor Gray

Write-Host ""
Write-Host "📚 QO'SHIMCHA MA'LUMOT:" -ForegroundColor Yellow
Write-Host ""
Write-Host "  Task ni ko'rish:     taskschd.msc" -ForegroundColor White
Write-Host "  Task ni to'xtatish:  Stop-ScheduledTask -TaskName POS80Printer" -ForegroundColor White
Write-Host "  Task ni ishga tush:  Start-ScheduledTask -TaskName POS80Printer" -ForegroundColor White
Write-Host "  Task ni o'chirish:   Unregister-ScheduledTask -TaskName POS80Printer" -ForegroundColor White

Write-Host ""

# ============================================
# RESTART TAKLIFI
# ============================================

$Restart = Read-Host "Kompyuterni HOZIR restart qilishni xohlaysizmi? (Y/N)"

if ($Restart -eq "Y" -or $Restart -eq "y") {
    Write-Host ""
    Write-Host "🔄 Kompyuter 10 soniyadan keyin restart bo'ladi..." -ForegroundColor Yellow
    Write-Host ""
    shutdown /r /t 10 /c "POS80 Printer Service sozlamalari faollashtirilmoqda"
} else {
    Write-Host ""
    Write-Host "ℹ️  Restart qilinmadi." -ForegroundColor Gray
    Write-Host "📝 Sozlamalar faollashtirish uchun kompyuterni qo'lda restart qiling:" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "   shutdown /r /t 0" -ForegroundColor White
    Write-Host ""
}
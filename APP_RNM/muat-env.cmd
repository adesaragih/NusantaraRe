@echo off
REM muat-env.cmd - memuat berkas .env ke sesi cmd.exe yang sedang terbuka.
REM
REM Cara pakai, dari folder APP_RNM:
REM
REM   call muat-env.cmd
REM
REM Sesudah itu, di jendela yang sama:  go run ./cmd/api
REM
REM Yang dilakukan:
REM   1. Membaca .env di folder skrip ini, baris NAMA=NILAI. Baris berawalan #
REM      dan baris kosong dilewati.
REM   2. Menyetel variabel lingkungan untuk tiap baris.
REM   3. Memperingatkan bila ORACLE_DSN masih memuat penanda <sandi>: cmd tidak
REM      punya prompt sandi tersembunyi, jadi isi sandinya di .env atau pakai
REM      PowerShell dengan  . .\muat-env.ps1
REM   4. Menambahkan folder Go dan Node ke PATH sesi.
REM
REM Nilai ORACLE_DSN tidak pernah ditampilkan di layar.

if not exist "%~dp0.env" (
    echo Tidak ada %~dp0.env - salin .env.example menjadi .env lalu isi.
    exit /b 1
)

for /f "usebackq eol=# tokens=1* delims==" %%A in ("%~dp0.env") do (
    set "%%A=%%B"
)

echo(%ORACLE_DSN%| findstr /c:"<sandi>" >nul
if not errorlevel 1 (
    echo PERINGATAN: ORACLE_DSN masih memuat ^<sandi^>. Isi sandinya di .env, atau pakai PowerShell:  . .\muat-env.ps1
)

set "PATH=C:\Program Files\Go\bin;C:\Program Files\nodejs;%PATH%"

echo HTTP_ADDR=%HTTP_ADDR%   ORACLE_SCHEMA=%ORACLE_SCHEMA%   IS_PEGA_PROD=%IS_PEGA_PROD%   ORACLE_DSN=disetel, tidak ditampilkan
echo Selesai. Backend: go run ./cmd/api   ^|   frontend: cd frontend ^&^& npm run dev

# muat-env.ps1 - memuat berkas .env ke sesi PowerShell yang sedang terbuka.
#
# Cara pakai, dari folder APP_RNM (perhatikan titik dan spasi di depan):
#
#   . .\muat-env.ps1
#
# Titik di depan disebut "dot-source": skrip dijalankan di dalam sesi ini,
# sehingga $env:... yang disetelnya tetap ada sesudah skrip selesai. Tanpa
# titik, variabelnya hilang begitu skrip berhenti.
#
# Yang dilakukan skrip ini:
#   1. Membaca .env di folder yang sama. Bentuk tiap baris: NAMA=NILAI.
#      Baris kosong dan baris berawalan # dilewati.
#   2. Menyetel $env:NAMA untuk tiap baris.
#   3. Bila ORACLE_DSN memuat teks <sandi>, menanyakan sandi Oracle tanpa
#      menampilkannya, lalu menyisipkannya ke DSN. Sandi tidak pernah dicetak
#      dan tidak pernah ditulis ke berkas.
#   4. Menambahkan folder Go dan Node ke PATH sesi bila belum ada.
#
# Nilai ORACLE_DSN tidak pernah ditampilkan di layar.

$berkasEnv = Join-Path $PSScriptRoot '.env'
if (-not (Test-Path $berkasEnv)) {
    Write-Host "Tidak ada $berkasEnv. Salin .env.example menjadi .env lalu isi."
    return
}

foreach ($baris in Get-Content $berkasEnv) {
    $b = $baris.Trim()
    if ($b -eq '' -or $b.StartsWith('#')) { continue }

    $posisi = $b.IndexOf('=')
    if ($posisi -lt 1) { continue }
    $nama  = $b.Substring(0, $posisi).Trim()
    $nilai = $b.Substring($posisi + 1).Trim()

    if ($nama -eq 'ORACLE_DSN' -and $nilai.Contains('<sandi>')) {
        # Ditanya sebagai SecureString supaya tidak tampil saat diketik.
        $aman = Read-Host -Prompt 'Sandi Oracle (tidak ditampilkan)' -AsSecureString
        $bstr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($aman)
        try {
            $polos = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($bstr)
        } finally {
            [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($bstr)
        }
        # Karakter khusus di sandi (misalnya @ atau /) harus di-escape di
        # dalam URL, kalau tidak DSN-nya salah baca.
        $nilai = $nilai.Replace('<sandi>', [Uri]::EscapeDataString($polos))
        $polos = $null
    }

    Set-Item -Path "Env:$nama" -Value $nilai
    if ($nama -eq 'ORACLE_DSN') {
        Write-Host 'ORACLE_DSN disetel (nilai tidak ditampilkan)'
    } else {
        Write-Host "$nama=$nilai"
    }
}

foreach ($folder in 'C:\Program Files\Go\bin', 'C:\Program Files\nodejs') {
    if (($env:Path -split ';') -notcontains $folder) {
        $env:Path = "$folder;$env:Path"
    }
}

Write-Host 'Selesai. Backend: go run ./cmd/api   |   frontend: cd frontend; npm run dev'

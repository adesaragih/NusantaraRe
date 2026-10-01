# Modul `claimlife` — Claim Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `claimlife` |
| Folder korpus | `Claim Life` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-CLAIMLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `001-029` |
| Slot menu | `950-951` |
| Prefix rute API | `/api/klaim-life`, `/api/peserta-life`, `/api/penyakit-life`, `/api/dokumen` |
| Kontrak disediakan | `kontrak.KlaimKomite` (dipakai `komiteclaimlife`) |
| Kontrak dipakai | `kontrak.PembacaPolis` (disediakan `premiumlistlife`) |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `migrations/` `modul.go` — paket Go `nusantarare/modul/claimlife/backend/...` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, LAPORAN, OQ — dulu `.scratch/claim-life/` |

## Migrasi

Rentang `001-029`, terpakai `001–022`. Slot menu `950-951` tidak terpakai: baris
modul ini sudah `DIMIGRASI = '1'` sejak 900, dan menu datar (30-09-2026) tidak punya butir — slot hanya
menyalakan `DIMIGRASI` (`APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6). Nama berkas migrasi yang sudah ada tidak pernah diubah:
`T_MIGRASI` mencatat nama.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/claimlife/...
go test -tags db ./modul/claimlife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/claimlife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` di akar repo, bab 8).

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Penjaganya berlaku untuk setiap modul; yang KHUSUS modul ini dinyatakan di sini, supaya
mengubahnya tidak pernah menyunting berkas di luar folder ini. Judul yang tidak ada berarti modul
ini tidak menyatakan apa pun untuk jenis itu. Nilai di dalam `` ` `` dibaca apa adanya.

### Nama STRUKTUR yang berbeda di DDL

Nama di dokumen STRUKTUR yang SENGAJA ditulis lain di DDL (`TestKolomDDLCocokDenganStruktur`); baris
yang tidak lagi terpakai membuat penjaga merah.

| Jenis | STRUKTUR | DDL | Alasan |
| --- | --- | --- | --- |
| tabel | `T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO` | `T_CLAIMLF_ADJ_SPREADING_RETRO` | Keputusan work owner 26 September 2026 (brief ronde 2 bab 2j): nama STRUKTUR berukuran 36 byte, dan Oracle di bawah 12.2 menolak pengenal lebih dari 30 byte dengan ORA-00972. Nama DDL berukuran 29 byte. |
| kolom | `RETROCESSION_VALUATION_BEGIN_DATE` | `RETRO_VALUATION_BEGIN_DATE` | Keputusan work owner 26 September 2026 (brief ronde 2 bab 2j): 33 dan 35 byte di STRUKTUR, keduanya melewati batas 30 byte. Nama penggantinya bukan karangan - STRUKTUR sendiri mencatat keduanya berasal dari korpus RETRO_VALUATION_BEGIN_DATE dan RETRO_VALUATION_EXPIRED_DATE. |
| kolom | `RETROCESSION_VALUATION_EXPIRED_DATE` | `RETRO_VALUATION_EXPIRED_DATE` | idem baris di atas. |

### Nama terlarang di migrasi

Nama yang tidak boleh muncul di migrasi modul MANA PUN (`TestNamaYangDibuangTidakAda`).

| Nama | Sebab |
| --- | --- |
| `T_CLAIMLF_POLICY` | tabel dihapus 2026-09-18, bukan diganti nama |
| `T_CLAIMLF_MARKETING` | tabel dihapus 2026-09-18, bukan diganti nama |
| `T_CLAIM_POLICY` | tabel dihapus 2026-09-18 |
| `T_CLAIM_MARKETING` | tabel dihapus 2026-09-18 |
| `WORK_CLAIM_ID` | dibuang; hubungannya shared primary key |
| `KMT_NO` | dibuang |
| `T_CLAIMLF_ADJUSTMENT_KOMITE` | roster komite tidak disimpan di Claim Life |

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini tanpa
`ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`). Mendaftarkan yang baru menuntut bukti.

| Awalan berkas | Relasi |
| --- | --- |
| `003_` | relasi 3 |
| `004_` | relasi 4 |
| `005_` | relasi 5 |
| `006_` | relasi 6 |
| `018_` | relasi 10: diagnosa per peserta (butir bd). `.DiagnoseList` hidup DI DALAM halaman peserta - `SetDisease.xml` b389 menutup dengan `Obj-Save pyWorkPage`, bukan menyimpan halaman diagnosa sendiri. Menghapus peserta karena itu menghapus daftarnya. |

### Penyuntikan wajib di handler

Penyusun layanan di `backend/handlers/<berkas>` yang bawaannya ber-stub, dan penyuntikan yang wajib
menyertainya di berkas yang sama (`TestHandlerMenyuntikkanImplementasiNyata`).

| Berkas | Penyusun | Wajib | Catatan |
| --- | --- | --- | --- |
| `akseptasi.go` | `svc.Akseptasi()` | `DenganJejak(jejak.PerekamJejakOracle(svc))`, `DenganPenerbit(services.PenerbitAkseptasiOracle(svc))` |  |
| `komite.go` | `svc.Komite()` | `DenganJejak(jejak.PerekamJejakOracle(svc))`, `DenganRoster(services.RosterKomiteOracle(svc))`, `DenganKasus(services.KasusKomiteOracle(svc))`, `DenganPenyalur(services.PenyalurClaimLifeOracle(svc))` |  |
| `putaran.go` | `svc.Putaran()` | `DenganJejak(jejak.PerekamJejakOracle(svc))` |  |
| `register.go` | `svc.Pendaftaran()` | `DenganPenomor(services.PenomorCounterOracle(svc))` |  |
| `tolak.go` | `svc.Status()` | `DenganJejak(jejak.PerekamJejakOracle(svc))` |  |

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE*.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` (folder `OUTPUT_HASIL_RNM\`).

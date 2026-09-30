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

Rentang `001-029`, terpakai `001–022`. Menu modul ini (butir `M_NAV_MENU` baru) ditulis di slot
`950-951`, di folder `backend/migrations/` modul ini sendiri — bentuk SQL-nya di
`PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6. Nama berkas migrasi yang sudah ada tidak pernah diubah:
`T_MIGRASI` mencatat nama.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/claimlife/...
go test -tags db ./modul/claimlife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/claimlife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`PANDUAN-TIM-PER-MODUL.md`).

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE*.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` (folder `OUTPUT_HASIL_RNM\`).

# Modul `premiumlistlife` — PremiumList Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `premiumlistlife` |
| Folder korpus | `PremiumList Life` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-PREMIUMLISTLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `050-099` |
| Slot menu | `954-955` |
| Prefix rute API | `/api/polis-life` |
| Kontrak disediakan | `kontrak.PembacaPolis` (dipakai `claimlife`) |
| Kontrak dipakai | — |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `migrations/` `modul.go` — paket Go `nusantarare/modul/premiumlistlife/backend/...` |
| `frontend/` | `pages/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, LAPORAN, OQ — dulu `.scratch/premiumlist-life/` |

## Migrasi

Rentang `050-099` *(ralat 30-09-2026: dulu tertulis `050–079` di panduan deploy)*, terpakai
`050–058`. Menu modul ini (butir `M_NAV_MENU` baru) ditulis di slot `954-955`, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nama berkas migrasi yang sudah ada tidak pernah diubah: `T_MIGRASI` mencatat nama.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/premiumlistlife/...
go test -tags db ./modul/premiumlistlife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/premiumlistlife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`PANDUAN-TIM-PER-MODUL.md`).

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-PREMIUMLIST-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` (folder `OUTPUT_HASIL_RNM\`).

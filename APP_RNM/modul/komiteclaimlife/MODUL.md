# Modul `komiteclaimlife` — Komite Claim Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `komiteclaimlife` |
| Folder korpus | `Komite Claim Life` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-KOMITECLAIMLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `030-049` |
| Slot menu | `952-953` |
| Prefix rute API | `/api/komite` |
| Kontrak disediakan | — |
| Kontrak dipakai | `kontrak.KlaimKomite` (disediakan `claimlife`) |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `migrations/` `modul.go` — paket Go `nusantarare/modul/komiteclaimlife/backend/...` |
| `frontend/` | `pages/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, LAPORAN, OQ — dulu `.scratch/komite-claim-life/` |

## Migrasi

Rentang `030-049`, terpakai `030`. Tabel tangga Komite sendiri lahir di migrasi Claim Life `013`
(sebelum modul ini berdiri) dan tetap di sana: `T_MIGRASI` mencatat nama, bukan letak. Menu modul
ini (butir `M_NAV_MENU` baru) ditulis di slot `952-953`, di folder `backend/migrations/` modul ini
sendiri — bentuk SQL-nya di `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/komiteclaimlife/...
go test -tags db ./modul/komiteclaimlife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/komiteclaimlife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`PANDUAN-TIM-PER-MODUL.md`).

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-LIFE.md`, `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-*.md` (folder `OUTPUT_HASIL_RNM\`).

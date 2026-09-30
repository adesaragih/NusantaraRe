# Modul `treatycontractout` — Treaty Contract Out

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatycontractout` |
| Folder korpus | `Treaty Contract Out` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-TREATYCONTRACTOUT` |
| Status | dimigrasi |
| Rentang migrasi | `300-319` |
| Slot menu | `956-957` |
| Prefix rute API | `/api/treaty-contract-out` (+ pekerja latar antrean lampiran) |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `modul.go` — paket Go `nusantarare/modul/treatycontractout/backend/...`; **tanpa** `migrations/` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` `proporsi.ts` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, LAPORAN, OQ — dulu `.scratch/treaty-contract-out/` |

## Migrasi

Rentang `300-319` **tetap kosong**: tco4 *(keputusan work owner 29-09-2026)* — modul ini menulis dan
membaca tabel warisan, nol tabel baru (`TestTCONolTabelBaru`). Slot menu `956-957` untuk butir
`M_NAV_MENU` baru bila kelak ada (bentuk SQL-nya di `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6);
berkas slot tidak membuat tabel, jadi tidak melanggar tco4.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/treatycontractout/...
go test -tags db ./modul/treatycontractout/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/treatycontractout
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`PANDUAN-TIM-PER-MODUL.md`).

## Brief acuan

`PROMPT-IMPLEMENTASI-MODUL-TREATY-CONTRACT-OUT.md`, `PROMPT-LANJUTAN-TREATY-CONTRACT-OUT-*.md` (folder `OUTPUT_HASIL_RNM\`).

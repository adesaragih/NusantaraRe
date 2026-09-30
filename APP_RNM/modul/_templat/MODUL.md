# Modul `<nama>` — <Folder korpus>

TEMPLAT folder modul — struktur tim satu folder per modul (30-09-2026). Salin folder `_templat/` menjadi
`modul/<nama>/`, lalu ganti setiap `<...>` di berkas ini. Ke-20 folder korpus SUDAH punya folder
kerangka (`claimfacin`, `nbtreatyin`, ...); templat ini untuk modul yang belum ada di korpus.

- `<nama>` — nama modul: nama folder `.scratch` tanpa tanda hubung, huruf kecil (tabel nama modul,
  `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 1). Sama dengan `MODUL_AKTIF` dan `const Nama`.
- `NNN-NNN` — rentang migrasi dan slot menu: AJUKAN ke tim inti lewat pull request dari rentang
  cadangan (`760-899` untuk migrasi, `990-999` untuk slot menu); dua modul tidak boleh berbagi nomor, dan
  nomor selalu tiga digit (`inti/backend/penjaga/rentang_test.go`).
- Langkah memulai kode (backend/modul.go, `go generate`, frontend/menu.ts + rute.tsx, slot menu):
  `docs/bersama/PANDUAN-TIM-PER-MODUL.md` (akar repo) bab 4–5.

⛔ Folder `_templat` sendiri TIDAK dibaca penjaga dan TIDAK terdaftar: nilainya penanda.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `<nama>` |
| Folder korpus | `<Folder korpus>` |
| GROUPMENU | `<TREATY / FACULTATIVE / KLAIM / MASTER>` |
| Pemilik | `@PEMILIK-<NAMA>` |
| Status | belum dimigrasi |
| Rentang migrasi | `NNN-NNN` |
| Slot menu | `NNN-NNN` |
| Prefix rute API | — (ditetapkan spec modul ini) |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/NN-<slug>.md`), grilling, catatan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `migrations/` `modul.go` — lahir saat modul dimulai |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `menu.ts` `rute.tsx` — lahir saat modul dimulai |

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/<nama>/...
npx vitest run modul/<nama>
```

## Pernyataan untuk penjaga

Tambahkan bab `###` di sini bila modul ini membutuhkannya — jenis dan bentuk tabelnya di
`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 6 (contoh lengkap: `modul/claimlife/MODUL.md`).

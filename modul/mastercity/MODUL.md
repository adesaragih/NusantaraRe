# Modul `mastercity` — City

Modul di LUAR dua puluh folder korpus (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5). Keputusan work owner
04-10-2026 (diteruskan sesi nusantarare-9d): Master Data dipecah menjadi delapan modul menu terpisah di grup MASTER —
"ini masing dibuatin didalam sub menu master data, bukan di jadiin 1 gini", cara pecah "8 modul terpisah". Nama tampilan
"City" (tanpa kata Master, pola Product Name Life). Mesin layar dan API-nya BERSAMA, milik tim inti:
`inti/backend/master` dan `inti/frontend/master` ("Pindah ke inti"). Rencana, keputusan, dan kontrak rute:
`modul/masterprovince/docs/issues/` (01 rencana, 02 API, 03 pemecahan delapan modul).

Satu folder, satu modul, satu pemilik.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah nilainya hanya
lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor. `Slot menu` — = modul luar korpus
tanpa slot (keputusan work owner 04-10-2026 "Modul luar korpus tanpa slot"): barisnya lahir DIMIGRASI '1' di migrasi
inti 914.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `mastercity` |
| Folder korpus | `Master City` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERCITY` |
| Status | dimigrasi |
| Rentang migrasi | `886-887` |
| Slot menu | — |
| Prefix rute API | `/api/master-city` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `modul.go` — `master.Pendaftaran(Nama, "/api/master-city", "city", …)`; definisi master `city` di `inti/backend/master/models/daftar.go`; tanpa migrasi (rentang dicadangkan) |
| `frontend/` | `menu.ts`, `rute.tsx`, `labels.ts` — layar `inti/frontend/master` (sesi 9d) |

## Tabel

Menulis `CITYINPUT` (dibuat `masterprovince` 880; `PROVINCENAME` + kota baru dari RW: 883, `modul/masterprovince/docs/issues/04-city-provincename.md`). EMAIL / MOID tidak tampil di menu (04-10-2026), kolomnya tetap. Rujukan Province dari `PROVINCE`; `BRANCHID` diperiksa ke `BRANCH` (tabel warisan yang dinyatakan `marketingofficer`), isian teks.

## Rute

`GET /api/master-city/meta` · `GET /api/master-city?q=&status=&halaman=` · `POST /api/master-city` · `PUT /api/master-city/{id}` ·
`PUT /api/master-city/{id}/status` · `GET /api/master-city/rujukan/{kolom}` — bentuk JSON-nya di
`modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`. Gerbang menu: akun wajib memegang menu `mastercity`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./inti/backend/master/... ./modul/mastercity/...
npx vitest run modul/mastercity inti/frontend/master
```

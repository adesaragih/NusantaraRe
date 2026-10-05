# Modul `masterdistrict` — District

Modul di LUAR dua puluh folder korpus (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5). Keputusan work owner
04-10-2026 (diteruskan sesi nusantarare-9d): Master Data dipecah menjadi delapan modul menu terpisah di grup MASTER —
"ini masing dibuatin didalam sub menu master data, bukan di jadiin 1 gini", cara pecah "8 modul terpisah". Nama tampilan
"District" (tanpa kata Master, pola Product Name Life). Mesin layar dan API-nya BERSAMA, milik tim inti:
`inti/backend/master` dan `inti/frontend/master` ("Pindah ke inti"). Rencana, keputusan, dan kontrak rute:
`modul/masterprovince/docs/issues/` (01 rencana, 02 API, 03 pemecahan delapan modul).

Satu folder, satu modul, satu pemilik.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah nilainya hanya
lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor. `Slot menu` — = modul luar korpus
tanpa slot (keputusan work owner 04-10-2026 "Modul luar korpus tanpa slot"): barisnya lahir DIMIGRASI '1' di migrasi
inti 915.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `masterdistrict` |
| Folder korpus | `Master District` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERDISTRICT` |
| Status | dimigrasi |
| Rentang migrasi | — |
| Slot menu | — |
| Prefix rute API | `/api/master-district` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

> Merge `origin/dev` 05-10-2026: `Rentang migrasi` `888-889` dilepas menjadi `—` (tanda modul tanpa migrasi sendiri, `tandaTanpaMigrasi`) - modul ini tidak punya satu pun berkas migrasi, dan nomor 880-899 sudah dipakai `aggregate` / `bordereaux` dari GitHub. Bila kelak butuh migrasi, minta jatah baru ke tim inti.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `modul.go` — `master.Pendaftaran(Nama, "/api/master-district", "district", …)`; definisi master `district` di `inti/backend/master/models/daftar.go`; tanpa migrasi (rentang dicadangkan) |
| `frontend/` | `menu.ts`, `rute.tsx`, `labels.ts` — layar `inti/frontend/master` (sesi 9d) |

## Tabel

Menulis `DISTRICTINPUT` (dibuat `masterprovince` 880). Rujukan City dari `CITYINPUT`.

## Rute

`GET /api/master-district/meta` · `GET /api/master-district?q=&status=&halaman=` · `POST /api/master-district` · `PUT /api/master-district/{id}` ·
`PUT /api/master-district/{id}/status` · `GET /api/master-district/rujukan/{kolom}` — bentuk JSON-nya di
`modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`. Gerbang menu: akun wajib memegang menu `masterdistrict`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./inti/backend/master/... ./modul/masterdistrict/...
npx vitest run modul/masterdistrict inti/frontend/master
```

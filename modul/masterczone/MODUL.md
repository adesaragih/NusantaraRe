# Modul `masterczone` — CZone

Modul di LUAR dua puluh folder korpus (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5). Keputusan work owner
04-10-2026 (diteruskan sesi nusantarare-9d): Master Data dipecah menjadi delapan modul menu terpisah di grup MASTER —
"ini masing dibuatin didalam sub menu master data, bukan di jadiin 1 gini", cara pecah "8 modul terpisah". Nama tampilan
"CZone" (tanpa kata Master, pola Product Name Life). Mesin layar dan API-nya BERSAMA, milik tim inti:
`inti/backend/master` dan `inti/frontend/master` ("Pindah ke inti"). Rencana, keputusan, dan kontrak rute:
`modul/masterprovince/docs/issues/` (01 rencana, 02 API, 03 pemecahan delapan modul).

Satu folder, satu modul, satu pemilik.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah nilainya hanya
lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor. `Slot menu` — = modul luar korpus
tanpa slot (keputusan work owner 04-10-2026 "Modul luar korpus tanpa slot"): barisnya lahir DIMIGRASI '1' di migrasi
inti 916.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `masterczone` |
| Folder korpus | `Master CZone` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERCZONE` |
| Status | dimigrasi |
| Rentang migrasi | — |
| Slot menu | — |
| Prefix rute API | `/api/master-czone` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

> Merge `origin/dev` 05-10-2026: `Rentang migrasi` `890-891` dilepas menjadi `—` (tanda modul tanpa migrasi sendiri, `tandaTanpaMigrasi`) - modul ini tidak punya satu pun berkas migrasi, dan nomor 880-899 sudah dipakai `aggregate` / `bordereaux` dari GitHub. Bila kelak butuh migrasi, minta jatah baru ke tim inti.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `modul.go` — `master.Pendaftaran(Nama, "/api/master-czone", "czone", …)`; definisi master `czone` di `inti/backend/master/models/daftar.go`; tanpa migrasi (rentang dicadangkan) |
| `frontend/` | `menu.ts`, `rute.tsx`, `labels.ts` — layar `inti/frontend/master` (sesi 9d) |

## Tabel

Menulis `CZONE` (dibuat `masterprovince` 880). Group Of merujuk `CZONE.CODE` tabel yang sama.

## Rute

`GET /api/master-czone/meta` · `GET /api/master-czone?q=&status=&halaman=` · `POST /api/master-czone` · `PUT /api/master-czone/{id}` ·
`PUT /api/master-czone/{id}/status` · `GET /api/master-czone/rujukan/{kolom}` — bentuk JSON-nya di
`modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`. Gerbang menu: akun wajib memegang menu `masterczone`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./inti/backend/master/... ./modul/masterczone/...
npx vitest run modul/masterczone inti/frontend/master
```

# Modul `masterobjectitemtype` — Object Item Type

Modul di LUAR dua puluh folder korpus (`docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 5). Keputusan work owner
04-10-2026 (diteruskan sesi nusantarare-9d): Master Data dipecah menjadi delapan modul menu terpisah di grup MASTER —
"ini masing dibuatin didalam sub menu master data, bukan di jadiin 1 gini", cara pecah "8 modul terpisah". Nama tampilan
"Object Item Type" (tanpa kata Master, pola Product Name Life). Mesin layar dan API-nya BERSAMA, milik tim inti:
`inti/backend/master` dan `inti/frontend/master` ("Pindah ke inti"). Rencana, keputusan, dan kontrak rute:
`modul/masterprovince/docs/issues/` (01 rencana, 02 API, 03 pemecahan delapan modul).

Satu folder, satu modul, satu pemilik.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah nilainya hanya
lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor. `Slot menu` — = modul luar korpus
tanpa slot (keputusan work owner 04-10-2026 "Modul luar korpus tanpa slot"): barisnya lahir DIMIGRASI '1' di migrasi
inti 919.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `masterobjectitemtype` |
| Folder korpus | `Master Object Item Type` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTEROBJECTITEMTYPE` |
| Status | dimigrasi |
| Rentang migrasi | `896-897` |
| Slot menu | — |
| Prefix rute API | `/api/master-object-item-type` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `modul.go` — `master.Pendaftaran(Nama, "/api/master-object-item-type", "objectitemtype", …)`; definisi master `objectitemtype` di `inti/backend/master/models/daftar.go`; tanpa migrasi (rentang dicadangkan) |
| `frontend/` | `menu.ts`, `rute.tsx`, `labels.ts` — layar `inti/frontend/master` (sesi 9d) |

## Tabel

Menulis isi `OBJECTITEMTYPE` (tabel warisan POOLDATA, dinyatakan `masterprovince`; strukturnya tidak diubah, status di kolomnya sendiri `ISACTIVE`). Jejak ubahnya di `T_MASTER_STATUS` (`masterprovince` 881 / 882).

## Rute

`GET /api/master-object-item-type/meta` · `GET /api/master-object-item-type?q=&status=&halaman=` · `POST /api/master-object-item-type` · `PUT /api/master-object-item-type/{id}` ·
`PUT /api/master-object-item-type/{id}/status` · `GET /api/master-object-item-type/rujukan/{kolom}` — bentuk JSON-nya di
`modul/masterprovince/docs/issues/03-pemecahan-delapan-modul.md`. Gerbang menu: akun wajib memegang menu `masterobjectitemtype`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./inti/backend/master/... ./modul/masterobjectitemtype/...
npx vitest run modul/masterobjectitemtype inti/frontend/master
```

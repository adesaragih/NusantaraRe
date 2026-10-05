# Modul `treatygroupojk` — Treaty Group OJK

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
05-10-2026: *"TREATYGROUPOJK buat juga crud nya, modul baru ya, jgn di gabung dan jangan diinti"*. Pega tidak punya
layar master untuk tabel ini; korpus hanya membacanya (`FetchTreatyGroupOJK`, folder NB Treaty In / NB FacIn / EDM
Treaty In) untuk mengisi kolom OJK baris `TREATYGROUP`. Baris menunya dibuat migrasi inti `916`, dan `Folder korpus`
di bawah = label menu `M_NAV_MENU.LABEL` (golongan MASTER).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — dipinjam dari jatah
Accounts (`850-854`, Accounts memakai 840-842; slot `995`, Accounts hanya memakai 994), persetujuan work owner
05-10-2026.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatygroupojk` |
| Folder korpus | `Treaty Group OJK` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-TREATYGROUPOJK` |
| Status | dimigrasi |
| Rentang migrasi | `850-854` |
| Slot menu | `995-995` |
| Prefix rute API | `/api/treaty-group-ojk` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-TREATYGROUPOJK.md` — peta tabel warisan dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `treatygroupojk.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- Tulis langsung ke `POOLDATA.TREATYGROUPOJK`: nol tabel baru, nol DDL. Add dan Edit; **tanpa hapus** (keputusan work
  owner 05-10-2026: "biarkan saja").
- ID baru = nomor ID tertinggi + 1, dua digit (`17`): tabel tanpa sequence; seluruh baris dikunci `FOR UPDATE` dulu
  supaya dua Add bersamaan tidak mendapat nomor yang sama.
- Order No TIDAK tampil dan tidak diketik (perintah work owner 05-10-2026: "ORDERNO hide aja, isi sesuai max dari
  order no"): Add = Order No tertinggi + 1 (sesudah baris dikunci), Edit membiarkannya. Daftar tetap urut ORDERNO.
- Name dan Name (IDN) wajib, huruf besar; Name tidak boleh kembar.
- Salinan OJK di `TREATYGROUP` (OJKBUSINESSNAME, OJKBUSINESSNAMEIDN, ORDERNO) TIDAK ikut diubah (keputusan work owner
  05-10-2026: "jgn ada ubah data").
- Hak menu Full / View only (`M_LOGIN_GO_MENU.HAK`): View only tanpa Add dan Edit.

## Migrasi

Nol DDL. Slot menu `995`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya sendiri dibuat migrasi inti
`916_m_nav_menu_treatygroupojk.sql`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/treatygroupojk/...
npx vitest run modul/treatygroupojk
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `TREATYGROUPOJK` | tabel warisan POOLDATA (lini bisnis OJK); modul ini menambah dan mengubah barisnya, tidak pernah membuat atau mengubah strukturnya (keputusan work owner 05-10-2026: nol DDL) |

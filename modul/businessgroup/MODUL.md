# Modul `businessgroup` — Business Group

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
05-10-2026: *"BUSINESSGROUP buat juga crud nya skalian bisa? modul baru"*. Pega tidak punya layar master untuk tabel
ini; korpus hanya membacanya (`BrowseBusinessGroup_RD`, kelas `ASM-FW-GISFW-Int-BUSINESSGROUP`, folder Treaty In dan
Treaty In Adjustment). Baris menunya dibuat migrasi inti `918`, dan `Folder korpus` di bawah = label menu
`M_NAV_MENU.LABEL` (golongan MASTER).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — dipinjam dari jatah
Accounts (`860-869`) dan Marketing Officer (slot `991`, Marketing Officer hanya memakai 990), persetujuan work owner
05-10-2026.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `businessgroup` |
| Folder korpus | `Business Group` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-BUSINESSGROUP` |
| Status | dimigrasi |
| Rentang migrasi | `860-869` |
| Slot menu | `991-991` |
| Prefix rute API | `/api/business-group` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-BUSINESSGROUP.md` — peta tabel warisan, relasi, dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `businessgroup.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- Tulis langsung ke `POOLDATA.BUSINESSGROUP`: nol tabel baru, nol DDL. Add dan Edit; **tanpa hapus** (keputusan work
  owner 05-10-2026: "biarkan saja"). Prosedur `UPSERT_BUSINESSGROUP` tidak dipanggil; `M_BUSINESSGROUP` (JSON Pega)
  tidak disentuh.
- **Tanpa SYARIAH** (keputusan work owner 05-10-2026, sama dengan saringan Pega `.Note NotEndsWith "SYARIAH"`): grup
  berakhiran SYARIAH tidak tampil, tidak dapat dibuka atau diubah; nama baru berakhiran SYARIAH ditolak.
- Treaty Group wajib (`TOPID` = `TREATYGROUP.ID`); TREATYNAME = salinan namanya. Edit yang tidak mengganti Treaty Group
  membiarkan salinannya (juga TOPID lama yang tidak ada lagi di TREATYGROUP).
- Name (`NOTE`) wajib, huruf besar, tidak kembar; Alias Name huruf besar, kosong = Name.
- ID baru = situs aktif `M_SITE_DATABASE` + `BUSINESSGROUP_SEQ` 4 digit (seperti `UPSERT_BUSINESSGROUP`); nomor yang
  sudah dipakai dilompati.
- Salinan nama di tabel lain (`BUSINESS.BUSINESSGROUPNAME`) TIDAK ikut diubah (keputusan work owner 05-10-2026: "jgn ada
  ubah data").
- Hak menu Full / View only (`M_LOGIN_GO_MENU.HAK`): View only tanpa Add dan Edit.

## Migrasi

Nol DDL. Slot menu `991`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya sendiri dibuat migrasi inti
`918_m_nav_menu_businessgroup.sql`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/businessgroup/...
npx vitest run modul/businessgroup
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

Nol tabel warisan baru yang dinyatakan di sini: `BUSINESSGROUP` sudah dinyatakan tabel warisan di `MODUL.md` Accounts
(penjaga menolak pernyataan ganda), `TREATYGROUP` di `MODUL.md` Treaty Group, dan `M_SITE_DATABASE` di `MODUL.md`
Adjuster Consultant.

# Modul `reinsurancetype` — Reinsurance Type

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
05-10-2026: *"select * from reinsurancetype"* (CRUD, modul baru, modul lain tidak disentuh). Pega tidak punya layar
master untuk tabel ini; korpus hanya membacanya (`BrowseReinsuranceType_RD` di sepuluh folder,
`GetReinsuranceTypeBYName_SQL` Claim Non Prop, `GetReinstypeIDbyName_SQL` EDM Treaty In, `BrowseReinsuranceTypeLimit_RD`
Master Contract Retro Life). Baris menunya dibuat migrasi inti `921`, dan `Folder korpus` di bawah = label menu
`M_NAV_MENU.LABEL` (golongan MASTER TREATY).

⛔ **Tanpa migrasi sendiri** (`—` di bawah, `tandaTanpaMigrasi` di `inti/backend/penjaga`): seluruh nomor modul
001-899 dan slot menu 950-999 sudah terbagi, dan work owner melarang menyentuh modul lain. Baris menunya dibuat migrasi
inti 921 langsung `DIMIGRASI '1'`.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `reinsurancetype` |
| Folder korpus | `Reinsurance Type` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-REINSURANCETYPE` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/reinsurance-type` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-REINSURANCETYPE.md` — peta tabel warisan, pembacanya, dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `reinsurancetype.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- Tulis langsung ke `POOLDATA.REINSURANCETYPE`: nol tabel baru, nol DDL. Add dan Edit; **tanpa hapus** - nonaktif =
  Flag `inactive` (keputusan work owner 05-10-2026). Prosedur `PEGA_REINSURANCETYPE` tidak dipanggil;
  `M_REINSURANCETYPE` (JSON Pega) tidak disentuh.
- Name (`NOTE`) wajib, huruf besar, **tidak kembar** (Pega mencari ID lewat nama); SOA Name huruf besar.
- Type `1` Own Retention, `2` Treaty Out, `3` Facultative, `4` Treaty In (Prompt value Pega, dari work owner).
  Flag `active` / `inactive` ("active inactive"). Group Type OR / QS / RI / SPL atau kosong. Nilai warisan lain (Flag
  `1` - disaring Contract Retro Life - dan kosong, Type kosong) dibiarkan selama tidak diubah.
- Code angka, kosong = `00`; No Urut angka atau kosong.
- ID baru = `1` + `M_REINSURANCETYPE_SEQ` 4 digit, seperti `PEGA_REINSURANCETYPE` (bukan `REINSURANCETYPE_SEQ`).
  TGLUPDATE = format Pega `YYYYMMDDTHHMMSS.mmm GMT`; USERID = akun login.
- Hak menu Full / View only (`M_LOGIN_GO_MENU.HAK`): View only tanpa Add dan Edit.

## Migrasi

Nol. Baris menu: migrasi inti `921_m_nav_menu_reinsurancetype.sql` (langsung menyala).

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/reinsurancetype/...
npx vitest run modul/reinsurancetype
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `REINSURANCETYPE` | tabel warisan POOLDATA (master jenis reasuransi); modul ini menambah dan mengubah barisnya, tidak pernah membuat atau mengubah strukturnya (keputusan work owner 05-10-2026: nol DDL) |

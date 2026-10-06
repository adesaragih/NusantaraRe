# Modul `treatygroup` — Treaty Group

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
05-10-2026: CRUD tabel `POOLDATA.TREATYGROUP` sebagai modul sendiri, di luar inti. Pega tidak punya layar master
untuk tabel ini; korpus hanya membacanya (`BrowseTreatyGroup_RD`, kelas `ASM-FW-GISFW-Int-TREATYGROUP`, di Treaty
Contract Out, Treaty In, NB Treaty In, klaim). Baris menunya dibuat migrasi inti `917`, dan `Folder korpus` di bawah =
label menu `M_NAV_MENU.LABEL` (golongan MASTER).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — dipinjam dari jatah
Accounts (`855-859`) dan Company Detail (slot `993`, Company Detail hanya memakai 992), persetujuan work owner
05-10-2026.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatygroup` |
| Folder korpus | `Treaty Group` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-TREATYGROUP` |
| Status | dimigrasi |
| Rentang migrasi | `855-859` |
| Slot menu | `993-993` |
| Prefix rute API | `/api/treaty-group` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-TREATYGROUP.md` — peta tabel warisan, relasi, dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `treatygroup.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- Tulis langsung ke `POOLDATA.TREATYGROUP`: nol tabel baru, nol DDL. Add dan Edit; **tanpa hapus** (keputusan work
  owner 05-10-2026: "biarkan saja"). Prosedur `PEGA_TREATYGROUP` / `PEGA_M_TREATYGROUP` tidak dipanggil;
  `M_TREATYGROUP` (JSON Pega) tidak disentuh.
- OJK Business wajib, dari `TREATYGROUPOJK`; OJKBUSINESSID, OJKBUSINESSNAME, OJKBUSINESSNAMEIDN, ORDERNO disalin dari
  baris OJK itu. Edit yang tidak mengganti OJK membiarkan salinannya. Order No tidak tampil di layar (perintah work
  owner 05-10-2026: "order id nya hide aja dari tampilan"); daftar tetap urut ORDERNO.
- Treaty Group Name wajib, huruf besar, tidak kembar; SOA Name boleh kosong, huruf besar.
- COAID (= `BUSINESSGROUP.ID`) opsi A (keputusan work owner 05-10-2026): tidak diketik; tampil dan berganti begitu OJK
  Business dipilih ("ubah pas pilih OJK Business") - COAID grup lain se-OJK yang paling sering, kosong bila belum
  ada. Edit yang tidak mengganti OJK membiarkannya.
- ID baru = situs aktif `M_SITE_DATABASE` + `TREATYGROUP_SEQ` 4 digit (seperti kedua prosedur Pega); nomor yang sudah
  dipakai dilompati. TGLUPDATE = format Pega `YYYYMMDDTHHMMSS.mmm GMT`; USERID = akun login.
- Salinan nama grup di tabel lain (BUSINESSGROUP.TREATYNAME, TREATYBUSINESS, TREATYINDETAIL, ...) TIDAK ikut diubah
  (keputusan work owner 05-10-2026: "jgn ada ubah data").
- View menampilkan grup bisnis anak (`BUSINESSGROUP.TOPID`), tanpa yang berakhiran SYARIAH.
- Hak menu Full / View only (`M_LOGIN_GO_MENU.HAK`): View only tanpa Add dan Edit.

## Migrasi

Nol DDL. Slot menu `993`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya sendiri dibuat migrasi inti
`917_m_nav_menu_treatygroup.sql`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/treatygroup/...
npx vitest run modul/treatygroup/
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `TREATYGROUP` | tabel warisan POOLDATA (grup treaty); modul ini menambah dan mengubah barisnya, tidak pernah membuat atau mengubah strukturnya (keputusan work owner 05-10-2026: nol DDL) |

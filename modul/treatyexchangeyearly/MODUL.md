# Modul `treatyexchangeyearly` — Treaty Exchange Yearly

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
05-10-2026: *"SELECT * FROM TREATYEXCHANGEYEARLY, CEK ITU DAN PAHAMI, BUAT CRUD JUGA"*. Pega tidak punya layar master
untuk tabel ini; korpus hanya membacanya (`BrowseTreatyExchangeYearly_RD` Treaty In / Treaty In Adjustment,
`GetCurrencyToIDR_SQL` / `GetKursLimitSpreading_SQL` FacIn dan NB Treaty In, `UploadCSVAggregate_Act`). Baris menunya
dibuat migrasi inti `919`, dan `Folder korpus` di bawah = label menu `M_NAV_MENU.LABEL` (golongan MASTER TREATY).

⛔ **Tanpa migrasi sendiri** (`—` di bawah, `tandaTanpaMigrasi` di `inti/backend/penjaga`): seluruh nomor modul
001-899 dan slot menu 950-999 sudah terbagi, dan work owner melarang menyentuh modul lain ("JANGAN ADA SENTUH MODUL
LAIN", 05-10-2026). Baris menunya dibuat migrasi inti 919 langsung `DIMIGRASI '1'`. Bila kelak butuh DDL, nomornya
harus diputuskan work owner lebih dulu.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatyexchangeyearly` |
| Folder korpus | `Treaty Exchange Yearly` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-TREATYEXCHANGEYEARLY` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/treaty-exchange-yearly` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-TREATYEXCHANGEYEARLY.md` — peta tabel warisan, pembacanya, dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `treatyexchangeyearly.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- Tulis langsung ke `POOLDATA.TREATYEXCHANGEYEARLY`: nol tabel baru, nol DDL. Add dan Edit; **tanpa hapus** (keputusan
  work owner 05-10-2026). Prosedur `PEGA_TREATYEXCHANGE` tidak dipanggil; `M_TREATYEXCHANGE` (JSON Pega) tidak
  disentuh.
- Form: Treaty Year (4 angka), Currency (view `CURRENCY`; IDCURRENCY + kodenya di CURRENCY), Start Date, End Date
  (bawaan 1 Juli - 30 Juni tahun berikutnya, boleh diubah), To IDR wajib; To USD boleh kosong; Quarter `0` (tahunan) -
  `4`, tampil di form (keputusan work owner: "tidak" disembunyikan).
- Tanggal ditulis format Pega yang benar `YYYYMMDDT000000.000 GMT` ("pake format seharusnya aja") - pembaca Pega
  membandingkannya SEBAGAI TEKS. Tanggal yang hari-nya tidak diubah saat Edit dibiarkan apa adanya.
- Treaty Year + Currency + Quarter tidak boleh kembar ("tolak"); tiga kombinasi kembar warisan DEV dibiarkan.
- ID baru = situs aktif `M_SITE_DATABASE` + `TREATYEXCHANGE_SEQ` 4 digit ("pake seq yang sudah ada"); nomor yang sudah
  dipakai dilompati. ID warisan tidak unik (10114-10116 dua baris), jadi Edit memakai ROWID baris.
- USERID = akun login; DATEIU = waktu simpan; DATEIN = waktu Add (format Pega GMT). QURRENCYID tidak ditulis.
- Hak menu Full / View only (`M_LOGIN_GO_MENU.HAK`): View only tanpa Add dan Edit.

## Migrasi

Nol. Baris menu: migrasi inti `919_m_nav_menu_treatyexchangeyearly.sql` (langsung menyala).

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/treatyexchangeyearly/...
npx vitest run modul/treatyexchangeyearly
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `TREATYEXCHANGEYEARLY` | tabel warisan POOLDATA (kurs tahunan treaty); modul ini menambah dan mengubah barisnya, tidak pernah membuat atau mengubah strukturnya (keputusan work owner 05-10-2026: nol DDL) |
| `CURRENCY` | view warisan POOLDATA di atas `M_CURRENCY`; pilihan Currency; dibaca saja |

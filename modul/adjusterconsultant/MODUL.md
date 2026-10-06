# Modul `adjusterconsultant` — Adjuster Consultant

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
05-10-2026: *"cek table ini POOLDATA.ADJUSTERCONSULTANT, aku mau kamu buatkan modul/menu CRUD nya"*; *"baris menu
tetep di inti, tapi adjusterconsultant jangan masukkan ke folder inti, buat modul sendiri"*. Padanan Pega: layar
master `MstAdjusterConsultant` di folder korpus Claim Fac In dan Claim Prop. Baris menunya dibuat migrasi inti `915`,
dan `Folder korpus` di bawah = label menu `M_NAV_MENU.LABEL` (keputusan work owner: "Adjuster Consultant", kelompok
MASTER).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — diambil dari jatah
Accounts (`870-879`, Accounts baru memakai 840-842) dan Aggregate (`997`, Aggregate hanya memakai 996).

| Kunci | Nilai |
| --- | --- |
| Nama modul | `adjusterconsultant` |
| Folder korpus | `Adjuster Consultant` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-ADJUSTERCONSULTANT` |
| Status | dimigrasi |
| Rentang migrasi | `870-879` |
| Slot menu | `997-997` |
| Prefix rute API | `/api/adjuster-consultant` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-ADJUSTERCONSULTANT.md` — peta tabel warisan dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `adjusterconsultant.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas)

- Tulis langsung ke `POOLDATA.ADJUSTERCONSULTANT`: nol tabel baru. Add, Edit, Activate/Deactivate; nol hapus
  permanen (Delete Pega membuka konfirmasi hapus Treaty Group - salin-tempel - jadi adjuster tidak pernah terhapus).
- ID baru seperti `GetIDConsultanAdj_SQL`: ID situs aktif `M_SITE_DATABASE` (`CURRENT_SITE` 1) + `ADJUSTERCONSULTANT_SEQ`
  4 digit; nomor yang sudah dipakai baris lama dilompati.
- Save seperti `SaveAdjusterConsultant_Act`: NAME dan ADDRESS huruf besar, TELPNO apa adanya, USERNAME = akun login,
  EDITDATE = SYSDATE. Tambahan keputusan work owner: Name wajib; nama ganda ditolak (tanpa beda huruf dan spasi tepi,
  juga terhadap baris nonaktif - pesannya menyarankan Activate).
- Hak menu Full / View only (`M_LOGIN_GO_MENU.HAK`): View only tanpa Add, Edit, Activate/Deactivate.

## Migrasi

Rentang `870-879`: `870_adjusterconsultant_aktif.sql` menambah `IS_ACTIVE VARCHAR2(1) DEFAULT '1' NOT NULL` +
`CK_ADJUSTERCONSULTANT_ACTIVE` ke tabel warisan (keputusan work owner 05-10-2026: flag nonaktif). Karena 870
mengubahnya, `ADJUSTERCONSULTANT` TIDAK terdaftar "Tabel warisan" di bawah, dan skema uji (`uji/skemauji`) membuat
tiruannya sebelum migrasi. Slot menu `997`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`. Barisnya sendiri
dibuat migrasi inti `915_m_nav_menu_adjusterconsultant.sql`.

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/adjusterconsultant/...
npx vitest run modul/adjusterconsultant
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `M_SITE_DATABASE` | tabel warisan POOLDATA; awalan ID baru (situs aktif); dibaca saja |

# Modul `nbfacin` — NB FacIn

⚠️ **Kerangka — belum dimigrasi.** Folder ini dibuat struktur tim satu folder per modul (keputusan work
owner 30-09-2026) supaya pemilik, rentang migrasi, dan slot menu modul ini TETAP sejak awal — satu
modul, satu folder, satu pemilik. Belum ada kode: tanpa `backend/modul.go` modul ini tidak terdaftar
(daftar Go bangkitan `inti/backend/daftar`, `import.meta.glob` frontend), dan kelompoknya di sidebar
tetap "belum dimigrasi" (`M_NAV_MENU.DIMIGRASI = '0'`). Cara memulainya:
`APP_RNM/PANDUAN-TIM-PER-MODUL.md` (akar repo) bab 4.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `nbfacin` |
| Folder korpus | `NB FacIn` |
| GROUPMENU | `FACULTATIVE` |
| Pemilik | `@PEMILIK-NBFACIN` |
| Status | dimigrasi |
| Rentang migrasi | `180-219` |
| Slot menu | `962-963` |
| Prefix rute API | `/api/nbfacin` |
| Kontrak disediakan | `kontrak.PenilaiPredikatFacIn`, `kontrak.MesinPremiFacIn` (lewat `Pendaftaran()`); `kontrak.TanggaAkseptasiFacIn` **belum** disediakan (butuh tabel limit per permintaan) — modul lain belum dapat memakainya: mengimpor `modul/nbfacin` terlarang (CLAUDE.md §5), sisa A32 |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), discovery, grilling — dipindah dari `jefri/OUTPUT FIX/` 30-09-2026 (`git mv`, isi tidak diubah). Asal tiap berkas dan tujuan tautan lamanya: `docs/PETA-ASAL.md` |
| `backend/` | `modul.go`; `models/`, `repository/` (tabel limit `POOLDATA.M_LIMIT_*` dan akun `POOLDATA.T_M_ACCOUNT`, bisnis `POOLDATA.BUSINESS`, baca saja), `services/` (layanan + mesin: `premium`, `acceptance`, `rules`, `pembayaran`, `rekonsiliasi`, `kontrakfacin`; `loader` = seam `loader.Flatten` data lama, murni, tiket 22), `handlers/` (`POST /api/nbfacin/premi`, `POST /api/nbfacin/akseptasi/langkah`, `GET /api/nbfacin/account`, `GET /api/nbfacin/class-of-business`), `migrations/` (slot menu 962) |
| `frontend/` | `menu.ts`, `rute.tsx`, `labels.ts` (verbatim korpus, diuji), `api.ts`, `pages/CoverageCargo.tsx` |

## Migrasi

Rentang `180-219` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `962-963` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Judul yang tidak ada berarti modul ini tidak menyatakan apa pun untuk jenis itu. Nilai di dalam
`` ` `` dibaca apa adanya.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang `docs/STRUKTUR-TABEL-NB-FACIN.md` gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris =
kepemilikan tabel berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `M_LIMIT_PROPERTYY` | tabel limit akseptasi warisan POOLDATA; NB hanya MEMBACA empat kolom (tiket 20) |
| `M_LIMIT_PROPERTY_NON_PREFERREDD` | tabel limit akseptasi warisan POOLDATA; NB hanya MEMBACA empat kolom (tiket 20) |
| `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` | tabel limit akseptasi warisan POOLDATA; NB hanya MEMBACA empat kolom (tiket 20) |
| `M_LIMIT_ENGINEERINGG` | tabel limit akseptasi warisan POOLDATA; NB hanya MEMBACA empat kolom (tiket 20) |
| `M_LIMIT_NONPROPANDENGG` | tabel limit akseptasi warisan POOLDATA; NB hanya MEMBACA empat kolom (tiket 20) |
| `M_LIMIT_FINANCIALINS` | tabel limit akseptasi bentuk B warisan POOLDATA; NB hanya MEMBACA empat kolom (tiket 20) |
| `T_M_ACCOUNT` | tabel akun warisan POOLDATA; NB hanya MEMBACA lima kolom untuk popup ChooseAccount (tiket 27) |
| `BUSINESS` | tabel bisnis warisan POOLDATA; NB hanya MEMBACA tiga kolom untuk pilihan Class Of Business (tiket 28) |

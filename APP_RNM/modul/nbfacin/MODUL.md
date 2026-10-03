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
| `backend/` | `modul.go`; `models/`, `repository/` (tabel limit `POOLDATA.M_LIMIT_*` dan akun `POOLDATA.T_M_ACCOUNT`, bisnis `POOLDATA.BUSINESS`, baca saja; case NB: tulis `T_WORK_POLIS` milik premiumlistlife (K-064, LINI `FAC`) + `T_NB_OPPORTUNITY`, tiket 29; layar Inward: `T_GENERAL_POLIS`/`T_QUOTATIONDATA` sebagian + `POOLDATA.MARKETINGOFFICER` baca saja, tiket 31), `services/` (layanan + mesin: `premium`, `acceptance`, `rules`, `pembayaran`, `rekonsiliasi`, `kontrakfacin`; `loader` = seam `loader.Flatten` data lama, murni, tiket 22), `handlers/` (`POST /api/nbfacin/premi`, `POST /api/nbfacin/akseptasi/langkah`, `GET /api/nbfacin/account`, `GET /api/nbfacin/class-of-business`, `POST /api/nbfacin/opportunity`, `GET /api/nbfacin/kasus/{caseId}`, `PUT /api/nbfacin/kasus/{caseId}/general`, `GET /api/nbfacin/marketing-officer`, `GET /api/nbfacin/opportunity` daftar portal tiket 32, `GET /api/nbfacin/sob` tiket 33, `GET`/`PUT /api/nbfacin/kasus/{caseId}/objek` tiket 35, `GET /api/nbfacin/risk-address` tiket 36, `POST /api/nbfacin/risk-address` dan `GET /api/nbfacin/rw` tiket 37, Surrounding Risk di `…/objek` dan `GET /api/nbfacin/occupation` tiket 38, Object Item di `…/objek`, `GET /api/nbfacin/jenis-item-objek` dan `GET /api/nbfacin/mata-uang` tiket 39, Occupation di `…/objek` dan `kdRiskExposure` / cari kosong di `GET /api/nbfacin/occupation` tiket 40, FEA di `…/objek` tiket 41, Loss Record + Loss Ratio di `…/objek` tiket 42, `GET /api/nbfacin/kasus/{caseId}/table-of-limit` tiket 40, Coverage di `…/objek`, `POST /api/nbfacin/kasus/{caseId}/hitung-coverage`, `GET /api/nbfacin/coverage` dan `GET /api/nbfacin/coverage-otomatis` tiket 43, `POST /api/nbfacin/kasus/{caseId}/hitung-net-rate` dan net rate di `…/objek` tiket 44), `migrations/` (180 `T_NB_OPPORTUNITY`, 181 `SEQ_WORK_POLIS_NB` — angka awal diisi work owner (penanda `{NB_MULAI}`), sudah dijalankan di DEV 03-10-2026; 182 `T_GENERAL_POLIS` dan 183 `T_QUOTATIONDATA` sebagian, 184 `T_QUOTATIONDATA.SOURCE_OF_BUSINESS`, 185 `T_CEDINGCOLIST` + kolom gabungan Ceding Co, 186 tabel tab Object FIRE, 187 `T_SURROUNDINGRISK` + empat kolom `T_PROPERTY` tiket 38, 188 `T_PROPERTYITEMLIST` tiket 39, 189 `T_OCCUPATIONLIST` + `T_TABLEOFLIMIT` tiket 40, 190 `T_FEALIST` tiket 41, 191 `T_LISTCAUSEOFLOSS` + `T_COINSDATA` + loss ratio `T_LOCATIONLIST` tiket 42, 192 lebar alamat risiko (Risk Location / Address 4000, delapan lainnya 100) butir 87/88, 193 `T_COVERAGELIST` + dua kolom total `T_PROPERTYITEMLIST` tiket 43; slot menu 962) |
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
| `BUSINESS` | tabel bisnis warisan POOLDATA; NB hanya MEMBACA tiga kolom untuk pilihan Class Of Business (tiket 28) dan kode bisnis Table of Limit (tiket 40) |
| `MARKETINGOFFICER` | tabel marketing officer warisan POOLDATA; NB hanya MEMBACA tiga kolom untuk pilihan Marketing Name (tiket 31) |
| `AGENT` | tabel agent warisan POOLDATA; NB hanya MEMBACA lima kolom untuk popup Change SOB (tiket 33) |
| `RISKADDRESS` | tabel alamat risiko warisan POOLDATA; NB MEMBACA sembilan kolom (popup Choose Risk Address, tiket 36) dan MENYISIPKAN alamat baru (popup Add, tiket 37; tidak membuat tabel) |
| `RW` | tabel kode pos/RW warisan POOLDATA; NB hanya MEMBACA (JOIN tiket 36, saran Zip Code tiket 37) |
| `OCCUPATION` | tabel okupasi warisan POOLDATA; NB hanya MEMBACA empat kolom untuk saran Occupation Surrounding Risk (tiket 38) dan popup Choose Occupation (tiket 40) |
| `V_JN_OBJ_ITEM` | view jenis item objek warisan POOLDATA; NB hanya MEMBACA empat kolom untuk pilihan Object Item Type (tiket 39) |
| `TABLEOFLIMIT` | tabel batas okupasi warisan POOLDATA; NB hanya MEMBACA lima kolom untuk popup Choose Class of Construction (tiket 40; TAHUN tidak disaring, A161) |
| `CURRENCY` | tabel mata uang warisan POOLDATA; NB hanya MEMBACA kolom CURRENCY untuk pilihan dan pemeriksaan mata uang item (tiket 39) |
| `COVERAGE` | tabel coverage warisan POOLDATA; NB hanya MEMBACA enam kolom untuk popup Choose Coverage (butir 96) dan lima coverage otomatis AddCoverageAutoFire (tiket 43) |

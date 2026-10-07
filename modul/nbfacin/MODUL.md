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
| Rentang migrasi | `180-214` |
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
| `backend/` | `modul.go`; `models/`, `repository/` (tabel limit `POOLDATA.M_LIMIT_*` dan akun `POOLDATA.T_M_ACCOUNT`, bisnis `POOLDATA.BUSINESS`, baca saja; case NB: tulis `T_WORK_POLIS` milik premiumlistlife (K-064, LINI `FAC`) + `T_NB_OPPORTUNITY`, tiket 29; layar Inward: `T_GENERAL_POLIS`/`T_QUOTATIONDATA` sebagian + `POOLDATA.MARKETINGOFFICER` baca saja, tiket 31), `services/` (layanan + mesin: `premium`, `acceptance`, `rules`, `pembayaran`, `rekonsiliasi`, `kontrakfacin`; `loader` = seam `loader.Flatten` data lama, murni, tiket 22), `handlers/` (`POST /api/nbfacin/premi`, `POST /api/nbfacin/akseptasi/langkah`, `GET /api/nbfacin/account`, `GET /api/nbfacin/class-of-business`, `POST /api/nbfacin/opportunity`, `GET /api/nbfacin/kasus/{caseId}`, `PUT /api/nbfacin/kasus/{caseId}/general`, `GET /api/nbfacin/marketing-officer`, `GET /api/nbfacin/opportunity` daftar portal tiket 32, `GET /api/nbfacin/sob` tiket 33, `GET`/`PUT /api/nbfacin/kasus/{caseId}/objek` tiket 35, `GET /api/nbfacin/risk-address` tiket 36, `POST /api/nbfacin/risk-address` dan `GET /api/nbfacin/rw` tiket 37, Surrounding Risk di `…/objek` dan `GET /api/nbfacin/occupation` tiket 38, Object Item di `…/objek`, `GET /api/nbfacin/jenis-item-objek` dan `GET /api/nbfacin/mata-uang` tiket 39, Occupation di `…/objek` dan `kdRiskExposure` / cari kosong di `GET /api/nbfacin/occupation` tiket 40, FEA di `…/objek` tiket 41, Loss Record + Loss Ratio di `…/objek` tiket 42, `GET /api/nbfacin/kasus/{caseId}/table-of-limit` tiket 40, Coverage di `…/objek`, `POST /api/nbfacin/kasus/{caseId}/hitung-coverage`, `GET /api/nbfacin/coverage` dan `GET /api/nbfacin/coverage-otomatis` tiket 43, `POST /api/nbfacin/kasus/{caseId}/hitung-net-rate` dan net rate di `…/objek` tiket 44, Deductible di `…/objek` tiket 45, unit / akumulasi coverage di `…/objek`, `GET /api/nbfacin/akumulasi` dan `GET /api/nbfacin/akumulasi/saran/{jenis}` tiket 46, `GET`/`PUT /api/nbfacin/kasus/{caseId}/klausa`, `GET /api/nbfacin/klausa` dan `GET /api/nbfacin/klausa/{id}/argumen` tiket 47, `GET`/`PUT /api/nbfacin/kasus/{caseId}/spreading`, `POST …/spreading/hitung-share` dan `POST …/spreading/salin` tiket 48, `GET`/`PUT /api/nbfacin/kasus/{caseId}/cedant` tiket 49, `GET`/`PUT /api/nbfacin/kasus/{caseId}/deduction` dan `POST …/deduction/hitung` tiket 50, `GET`/`PUT /api/nbfacin/kasus/{caseId}/payment` dan `POST …/payment/hitung` tiket 51, Zone / 4.2 Construction di `…/objek`, `GET /api/nbfacin/zona`, `GET /api/nbfacin/konstruksi-gempa` dan `callAct` "rate" di `…/hitung-coverage` tiket 52, `GET`/`PUT /api/nbfacin/kasus/{caseId}/scoring` dan `POST …/scoring/hitung` tiket 54, `GET`/`PUT /api/nbfacin/kasus/{caseId}/loss-record`, `POST …/loss-record/hitung`, `GET /api/nbfacin/class-of-business/saran` dan `aksesLossRecord` di GET kasus tiket 55 — akses = `M_LOGIN_GO.JOB_POSITION` "IT" (tabel inti login, baca saja), `GET`/`POST /api/nbfacin/kasus/{caseId}/korespondensi`, `PUT`/`DELETE …/korespondensi/{no}`, `POST …/korespondensi/{no}/kirim` (email SMTP NYATA, persetujuan work owner 06-10-2026; TLS wajib; env `NBFACIN_SMTP_HOST` / `_PORT` / `_USER` / `_PASSWORD` / `_FROM` / `_TANPA_TLS`, dibaca `inti/backend/config`) dan `GET /api/nbfacin/korespondensi/penerima/saran` tiket 56 — hanya JOB_POSITION "IT" seperti Loss Record (`aksesKorespondensi` di GET kasus), `GET /api/nbfacin/kasus/{caseId}/submit`, `POST …/submit/periksa`, `POST …/submit` tiket 57 (Submit tahap marketing: workbasket ReasFacInMarketing; MENULIS tabel lama `POOLDATA.HISTORYAKSEPTASIPRODUCTION` / `HISTORYAKSEPTASIPEGA` / `FACINPERFORMANCE` (keputusan work owner 07-10-2026 butir 4), membaca `POOLDATA.M_LIMIT_PROPERTYY` / `CAUSEOFDECLINE` / `MARKETINGOFFICER` / `RW` / `OCCUPATION`; mengubah `T_WORK_POLIS.POSITION`, `T_GENERAL_POLIS.NB_STATUS`, mengunci `T_CORRESPONDENCELIST` / `T_RECIPIENTLIST`); pengirim / penerima = `M_LOGIN_GO.NAME` / `EMAIL`, baca saja), `migrations/` (180 `T_NB_OPPORTUNITY`, 181 `SEQ_WORK_POLIS_NB` — angka awal diisi work owner (penanda `{NB_MULAI}`), sudah dijalankan di DEV 03-10-2026; 182 `T_GENERAL_POLIS` dan 183 `T_QUOTATIONDATA` sebagian, 184 `T_QUOTATIONDATA.SOURCE_OF_BUSINESS`, 185 `T_CEDINGCOLIST` + kolom gabungan Ceding Co, 186 tabel tab Object FIRE, 187 `T_SURROUNDINGRISK` + empat kolom `T_PROPERTY` tiket 38, 188 `T_PROPERTYITEMLIST` tiket 39, 189 `T_OCCUPATIONLIST` + `T_TABLEOFLIMIT` tiket 40, 190 `T_FEALIST` tiket 41, 191 `T_LISTCAUSEOFLOSS` + `T_COINSDATA` + loss ratio `T_LOCATIONLIST` tiket 42, 192 lebar alamat risiko (Risk Location / Address 4000, delapan lainnya 100) butir 87/88, 193 `T_COVERAGELIST` + dua kolom total `T_PROPERTYITEMLIST` tiket 43, 194 `T_DEDUCTIBLELIST` tiket 45, 195 tiga kolom unit / akumulasi `T_COVERAGELIST` tiket 46, 197 `T_CLAUSELIST` (klausa + argumen satu tabel) tiket 47, 198 `T_SPREADINGLIST` + `T_GENERAL_POLIS.PERCENT_SHARE` + dua kolom Nusantara Re `T_COVERAGELIST` tiket 48 — wajib jalan SEBELUM biner baru (Save Object menghapus `T_SPREADINGLIST`), 199 `T_CEDINGCEDANTLIST` + `T_GENERAL_POLIS.SHARE_CEDANT_TYPE` + `T_QUOTATIONDATA.SHARE_OF_CEDING` tiket 49 — juga sebelum biner baru, 200 `T_CURRENCYLIST` + tiga kolom Premium Deduction `T_GENERAL_POLIS` tiket 50, 201 `T_LISTINSTALLMENT` + `PAY_INSTALLMENT` + Total Payment / Net Premium `T_CURRENCYLIST` tiket 51 — 200 DAN 201 sebelum biner baru (Save Premium Deduction menyentuh `T_LISTINSTALLMENT`), 202 `T_GENERAL_POLIS.PPN_CHECK` NUMBER(5) -> VARCHAR2(10) teks "true" / "false" (ralat tiket 50, butir 68.1; langkah DEV di docs/issues/50 K50-13) — sebelum biner baru, 203 `T_COVERAGELIST.ZONE` VARCHAR2(1000) + `CONSTRUCTION_EARTHQUAKE` VARCHAR2(500) tiket 52 — sebelum biner baru (GET / PUT `…/objek` membaca dan menulis keduanya), 204 `T_COVERAGELIST.ZONE` VARCHAR2(1000) -> VARCHAR2(4000) (bug DEV tiket 52, "buat 4000"; langkah DEV di docs/issues/52 K52-13) — sebelum biner baru, 205 lima tabel Scoring Risk (`T_SCORINGRISK`, `T_SCORINGRESULT`, `T_DATASCORINGRISKLIST`, `T_SCORING_FACTOR`, `T_SCORING_OPTION` + sequence) tiket 54 — sebelum biner baru (GET kasus membaca `T_SCORINGRISK` untuk `riskScoring` General: tanpa 205 seluruh layar kasus ORA-00942), 206 `T_LOSSRECORD` + `T_LOSSRECORDLIST` (tabel BARU tab Loss Record tingkat case) tiket 55 — sebelum biner baru bila tab dipakai (GET kasus TIDAK membacanya; tanpa 206 hanya tab Loss Record yang gagal), 207 `T_CORRESPONDENCELIST` + `T_RECIPIENTLIST` (tabel BARU tab Correspondence) tiket 56 — sebelum biner baru bila tab dipakai (GET kasus TIDAK membacanya), 208 `T_GENERAL_POLIS.NB_STATUS` tiket 57 — **WAJIB sebelum biner** sejak portal Opportunity (`GET /api/nbfacin/opportunity`, tiket 32) membaca `NB_STATUS` (tanpa 208: ORA-00904, daftar case gagal) (Submit juga butuh 207: kunci korespondensi), 209 `T_LISTCAUSEOFLOSS.CURRENCY` boleh NULL tiket 42b — sebelum biner baru (biner baru menulis NULL untuk catatan kerugian tanpa mata uang; tanpa 209 Save tab Object yang memuatnya gagal ORA-01400) (196 tidak dipakai lagi: nomornya dulu view → tabel flat yang dipindah), ~~196 view PROVINCE / ACCUMULATEDTYPE / CZONE / ACCUMULATION → tabel flat~~ dipindah ke masterdata 760, lalu `masterprovince` 880, kini `masterprovince` 215 (MD-1, 04-10-2026; 05-10-2026); slot menu 962). Rentang dipersempit 05-10-2026 dari `180-219` ke `180-214`: `215-219` diberikan ke `masterprovince` (bentrok 880 dengan modul `aggregate`, keputusan work owner "ikutin rekomendasi terbaik km aja"); nomor berikut nbfacin = 210 |
| `frontend/` | `menu.ts`, `rute.tsx`, `labels.ts` (verbatim korpus, diuji), `api.ts`, `pages/CoverageCargo.tsx` |

## Migrasi

Rentang `180-219` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `962-963` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

⛔ **Urutan deploy saran akumulasi (tiket 46, MD-5)**: `repository/akumulasi.go` menyaring `STS_AKTIF = '1'` dan membaca
`CITYINPUT` / `DISTRICTINPUT` / `T_MASTER_STATUS` — tabel milik **masterprovince** 880 / 881 (dulu masterdata 760 / 761;
pemecahan delapan modul master 04-10-2026), bukan modul ini. Sebelum migrasi itu berjalan, saran Choose Accumulation gagal
(ORA-00904 / ORA-00942). Urutannya: (1) biner yang memuat commit `072370f0`, (2) `-migrate` (880 → 881 → 882 → 912-920),
(3) baru buka saran akumulasi. Membaca kegagalan dan sisa `*_SALIN`: `modul/masterprovince/MODUL.md` bab "Urutan deploy".

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
| `BUSINESS` | tabel bisnis warisan POOLDATA; NB hanya MEMBACA tiga kolom untuk pilihan Class Of Business (tiket 28) dan kode bisnis Table of Limit (tiket 40) |
| `RISKADDRESS` | tabel alamat risiko warisan POOLDATA; NB MEMBACA sembilan kolom (popup Choose Risk Address, tiket 36) dan MENYISIPKAN alamat baru (popup Add, tiket 37; tidak membuat tabel) |
| `RW` | tabel kode pos/RW warisan POOLDATA; NB hanya MEMBACA (JOIN tiket 36, saran Zip Code tiket 37, popup akumulasi tiket 46) |
| `CITY` | view kota warisan POOLDATA; NB hanya MEMBACA tiga kolom untuk saran City popup akumulasi (tiket 46) |
| `DISTRICT` | view kecamatan warisan POOLDATA; NB hanya MEMBACA empat kolom untuk saran District popup akumulasi (tiket 46) |
| `JSON_POLIS` | tabel polis JSON warisan POOLDATA; NB hanya MEMBACA NOPOLIS / DATA_JSON untuk cari akumulasi menurut nomor polis (tiket 46) |
| `OCCUPATION` | tabel okupasi warisan POOLDATA; NB hanya MEMBACA tujuh kolom untuk saran Occupation Surrounding Risk (tiket 38), popup Choose Occupation (tiket 40), rate FLEXAS RATELOWER (tiket 52), dan KDRISKEXPOSURE Scoring Risk (tiket 54) |
| `V_JN_OBJ_ITEM` | view jenis item objek warisan POOLDATA; NB hanya MEMBACA empat kolom untuk pilihan Object Item Type (tiket 39) |
| `TABLEOFLIMIT` | tabel batas okupasi warisan POOLDATA; NB hanya MEMBACA lima kolom untuk popup Choose Class of Construction (tiket 40; TAHUN tidak disaring, A161) |
| `COVERAGE` | tabel coverage warisan POOLDATA; NB hanya MEMBACA enam kolom untuk popup Choose Coverage (butir 96) dan lima coverage otomatis AddCoverageAutoFire (tiket 43) |
| `ZONES` | view zona warisan POOLDATA (atas JSON m_zones); NB hanya MEMBACA empat kolom untuk saran Zone form coverage FIRE (tiket 52) |
| `CONSTRUCTION` | view konstruksi warisan POOLDATA (atas JSON m_construction); NB hanya MEMBACA CONSTRUCTION untuk saran 4.2 Construction (tiket 52) |
| `ZONEOJK` | tabel rate zona gempa warisan POOLDATA; NB hanya MEMBACA DESCRIPTION / RATE / ID untuk rate 4.2 PRGBI (tiket 52) |

Dicabut 04-10-2026 saat merge `origin/dev` (perintah work owner: "ikuti yang dari github, kalau bentrok dengan kerjaan
saya disesuaikan"); NB tetap hanya MEMBACA keduanya: `T_M_ACCOUNT` - kini diubah migrasi modul `accounts` (840-842),
jadi bukan lagi "tidak dibuat migrasi mana pun"; `MARKETINGOFFICER` - kini dinyatakan modul `marketingofficer` (satu
tabel warisan hanya boleh dinyatakan satu modul). Merge `origin/dev` 05-10-2026: `AGENT` ikut dicabut - kini dinyatakan modul `bordereaux`
(Cedant, popup Choose Master Treaty); NB tetap hanya MEMBACA lima kolomnya untuk popup Change SOB (tiket 33). `CURRENCY` juga dicabut - kini dinyatakan modul
`treatyexchangeyearly`; NB tetap hanya MEMBACA kolom CURRENCY untuk pilihan dan pemeriksaan mata uang item (tiket 39).

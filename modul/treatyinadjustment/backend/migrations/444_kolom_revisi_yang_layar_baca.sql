-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment`. Itu DISENGAJA, dan sebabnya bukan
-- kerapian melainkan paksaan:
--
--   Rentang migrasi `treatyin` adalah `400-439` dan ia PENUH TANPA CELAH --
--   keempat puluh nomornya terpakai (`400` sampai `439`). Nomor `444` jatuh
--   di rentang `treatyinadjustment` (`440-479`), dan di sanalah berkas ini
--   harus berdiri supaya `TestSetiapMigrasiDiRentangAtauSlotModulnya` benar.
--
-- ⚠️ JADI INI TEMPAT PARKIR, BUKAN PERNYATAAN KEPEMILIKAN. `T_TREATY_REVISION`
--   tetap tabel modul `treatyin`; petanya, pemuatnya, dan pembacanya hidup di
--   sana. Yang membuat parkir ini tidak sepenuhnya asing: modul Adjustment
--   MEMBACA tabel yang sama untuk kedua sisinya, dan memegang salinan petanya
--   (`uji/lintasmodul/peta_pendaratan_test.go` yang mengadu keduanya).
--
-- ⛔ YANG SEHARUSNYA: `treatyin` meminta jatah rentang baru kepada tim inti --
--   prosedur yang `MODUL.md` sebut sendiri ("Bila kelak butuh migrasi, minta
--   jatah baru ke tim inti"). Sampai jatah itu ada, migrasi `treatyin`
--   BERIKUTNYA tidak punya nomor sah sama sekali. Itu permintaan terbuka,
--   bukan hal yang dapat diputuskan di dalam berkas ini.
--
-- ---------------------------------------------------------------------
-- Enam belas kolom untuk `T_TREATY_REVISION` — medan yang LAYAR TAMPILKAN
-- tetapi tabelnya belum punya tempatnya.
--
-- ---------------------------------------------------------------------
-- MASALAH YANG DITUTUPNYA, dan sejak kapan ia terbuka
-- ---------------------------------------------------------------------
-- `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §17 mencatat medan kepala tab
-- `Reporting Period` — Start, End, Period, Interval, Submission,
-- Confirmation, Settlement — KOSONG ketika tab dibuka, dengan sebab:
-- *"nilai tersimpannya hanya ada di `JSONDATA`, yang dilarang dibaca, dan
-- `T_TREATY_REVISION` tidak punya kolom `Reporting*`."*
--
-- Separuh pertama sebab itu MENGIKAT (larangan pemilik proses). Separuh
-- keduanya tidak — ia hanya kolom yang belum dibuat, dan berkas ini
-- membuatnya.
--
-- ---------------------------------------------------------------------
-- PENGUKURAN YANG MENDAHULUI BERKAS INI
-- ---------------------------------------------------------------------
-- Disapu 6 Oktober 2026 atas SELURUH 1.855 dokumen `POOLDATA.M_TREATY_IN`
-- (nol `ROWNUM`, nol sampel; 1.855 terurai, nol gagal urai). Yang dihitung
-- adalah kunci AKAR yang `KunciTakTerpetakan` laporkan tanpa kolom, beserta
-- berapa dokumen yang BENAR-BENAR mengisinya:
--
--   ReportingPeriod         1.851      CedingID               1.851
--   AccumulationPeriod      1.851      LeadingReinsSourceID   1.851
--   ChooseStatusAkseptasi   1.847      ReportingStart         1.219
--   ReportingEnd            1.219      ReportingSubmission    1.219
--   ReportingConfirmation   1.219      ReportingSettlement    1.219
--   Comment                   314      ReportingInterval        300
--   LeadingReinsID            297      Information              196
--   PositionUsername           34      Position                  33
--
-- ⚠ Ketujuh `Reporting*` terisi pada 1.219 dari 1.855 — DUA PERTIGA, bukan
--   kasus pinggiran. Selama ini layar memperlihatkan medan kosong pada
--   seribu dua ratus kontrak yang sungguh punya nilainya.
--
-- ---------------------------------------------------------------------
-- ⭐ MENGAPA KEENAM BELAS INI, DAN BUKAN KETIGA PULUH ENAM YANG TERUKUR
-- ---------------------------------------------------------------------
-- Sapuan yang sama menemukan 36 kunci akar bukan-perabot tanpa kolom.
-- Memberi kolom kepada seluruhnya akan memaksa penilaian "ini tersimpan
-- atau ini turunan?" tiga puluh enam kali, dan penilaian itu persis yang
-- sedang DITANYAKAN kepada pemilik proses di
-- `PERTANYAAN-TERBUKA-DITURUNKAN-DI-PENDARATAN.md`.
--
-- Jadi saringannya diambil dari EKSPOR, bukan dari selera:
--
--   ⭐ Kunci yang menjadi PARAMETER prosedur `POOLDATA.PEGA_TREATY_IN`
--      adalah data TERSIMPAN — sebab prosedur itulah yang Pega pakai untuk
--      MENYIMPANNYA (`RDBList/SaveTreatyIn.xml`, lihat
--      `docs/SPESIFIKASI-TOMBOL-DARI-EKSPOR.md` §1).
--
-- Parameternya: `CedingID` · `LeadingReinsSourceID` · `LeadingReinsID` ·
-- `Information` · `Position` · `PositionUsername` · `ChooseStatusAkseptasi`
-- (sisanya sudah berkolom sejak `439`).
--
-- Ditambah ketujuh `Reporting*` yang §17 sebut namanya satu per satu, dan
-- `AccumulationPeriod` + `Comment` yang BERPASANGAN dengan tab yang sudah
-- berdiri (`T_TREATY_ACCUMULATION`, `T_VIEW_COMMENT`) — medan kepala tanpa
-- kepala.
--
-- ⛔ YANG SENGAJA TIDAK DIBERI KOLOM, dan sebabnya satu kalimat:
--   `Total*NP` (9) · `TotalSpreaded*` (4) · `*SummaryList` (3) ·
--   `RnmShareDeducted` — seluruhnya DIHITUNG oleh Activity
--   (`TreatyInNonAddItem`, `TreatyInNPSetTotal`), jadi ia termasuk
--   pertanyaan `DITURUNKAN` yang belum dijawab.
--   `ViewState` · `IsEditData`-sejenis · `px*` · `py*` — perabot Pega.
--   `CurrencyList` — SALINAN `TREATYEXCHANGEYEARLY` (`KOREKSI-ERD-VERSUS-
--   POOLDATA.md` §1.2); mendaratkannya melahirkan kurs kedua.
--   `RSMDLimit`/`Earthquake`/`FloodJab`/`FloodNation` + mata uangnya —
--   sudah berumah di `T_TREATY_LIMIT_DETAIL` sejak `437`, dan kepala `439`
--   sudah menolak `T_TREATY_HAZARD_LIMIT` dengan sebab yang sama.
--
--   ⛔⛔ RALAT, 6 Oktober 2026, sore — KALIMAT DI ATAS SALAH UNTUK NON-PROP.
--   Ia benar untuk nilai TINGKAT DETAIL (`Limits[].Detail[]`, yang tab
--   Limits PROPORSIONAL tampilkan). Tetapi tab Event Limits NON-PROP
--   mengikat properti AKAR (`TreatyIn.RSMDLimit` …,
--   `Section/TreatyInTabsNonProportional.xml`), dan keduanya TEMPAT YANG
--   BERBEDA — diukur atas seluruh 1.855 dokumen:
--
--                      akar berisi   Detail berisi   tumpang tindih
--     Proportional              0             347                0
--     NonProportional          49               1                0
--
--   Jadi delapan kunci akar itu TIDAK berumah di mana pun, dan 49 kontrak
--   Non-Prop kehilangan nilainya di layar. Kolomnya menuntut migrasi —
--   dan rentang migrasi `treatyin` habis (lihat kepala berkas ini). Ia
--   masuk daftar kebutuhan, bersama `MaxCoGroup`/`MaxCoNonGroup`.
--
--   ⚠️ Kalimat yang salah TIDAK dihapus: berkas ini sudah dijalankan, dan
--   jejak alasan yang keliru lebih berguna daripada catatan yang tampak
--   selalu benar.
--
-- ---------------------------------------------------------------------
-- BENTUK
-- ---------------------------------------------------------------------
-- ⛔ `ALTER TABLE … ADD`, bukan menyunting `CREATE TABLE` di `439`:
--   migrasi `439` SUDAH DIJALANKAN dan tercatat di `T_MIGRASI`.
--   Menyuntingnya berarti dua basis data dengan riwayat yang sama memiliki
--   bentuk yang berbeda.
--
-- ⚠ Seluruhnya NULLABLE. Tidak satu pun terisi di SELURUH 1.855 dokumen
--   (yang tertinggi 1.851), jadi `NOT NULL` akan menolak kontrak yang
--   sistem lama terima.
--
-- `VARCHAR2(4000 CHAR)` mengikuti ke-37 kolom sekelasnya di `439`. Nol
-- di antaranya teks panjang: yang panjang (`Exclusions`,
-- `SpecialConditions`) sudah `CLOB` di sana.
--
-- ⛔ `COMMENTTEKS`, BUKAN `COMMENT`. `COMMENT` kata TERCADANG Oracle
--   (ia pernyataan DDL tersendiri), dan kolom bernama demikian menuntut
--   tanda kutip ganda di SETIAP kueri yang menyebutnya -- termasuk SQL
--   yang peta bangkitkan, yang tidak mengutip apa pun. Kunci dokumennya
--   tetap `Comment`; yang berbeda hanya nama kolomnya.

ALTER TABLE {skema}.T_TREATY_REVISION ADD (
  CEDINGID                         VARCHAR2(4000 CHAR),
  LEADINGREINSSOURCEID             VARCHAR2(4000 CHAR),
  LEADINGREINSID                   VARCHAR2(4000 CHAR),
  INFORMATION                      VARCHAR2(4000 CHAR),
  POSITION                         VARCHAR2(4000 CHAR),
  POSITIONUSERNAME                 VARCHAR2(4000 CHAR),
  CHOOSESTATUSAKSEPTASI            VARCHAR2(4000 CHAR),
  ACCUMULATIONPERIOD               VARCHAR2(4000 CHAR),
  COMMENTTEKS                      VARCHAR2(4000 CHAR),
  REPORTINGSTART                   VARCHAR2(4000 CHAR),
  REPORTINGEND                     VARCHAR2(4000 CHAR),
  REPORTINGPERIOD                  VARCHAR2(4000 CHAR),
  REPORTINGINTERVAL                VARCHAR2(4000 CHAR),
  REPORTINGSUBMISSION              VARCHAR2(4000 CHAR),
  REPORTINGCONFIRMATION            VARCHAR2(4000 CHAR),
  REPORTINGSETTLEMENT              VARCHAR2(4000 CHAR)
)
/

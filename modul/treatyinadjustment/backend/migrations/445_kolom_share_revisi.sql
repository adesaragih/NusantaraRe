-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444_kolom_revisi_yang_layar_baca.sql`, dengan sebab yang sama: rentang
-- `treatyin` (`400-439`) PENUH, dan `445` jatuh di rentang modul ini
-- (`440-479`). Nomor `445` diberikan pemegang modul ini, 7 Oktober 2026.
-- `T_TREATY_REVISION` tetap tabel modul `treatyin`.
--
-- ---------------------------------------------------------------------
-- Empat kolom akar tab Share Non-Prop
-- ---------------------------------------------------------------------
--
--   RNMSHARE                 TreatyIn.RNMShare              (% RNM Share)
--   BROKERAGEPERCENT         TreatyIn.BrokeragePercent      (% Brokerage)
--   RNMSHAREACROSSTHEBOARD   TreatyIn.RNMShareAcrossTheBoard (Share Across The Board)
--   RNMSHAREDEDUCTED         TreatyIn.RnmShareDeducted      (Share to RNM :)
--
-- ⭐ KEPUTUSAN PEMAKAI 6–7 Oktober 2026: Save dan Submit menyimpan ke tabel
-- masing-masing yang ditentukan — skema v2 memetakan akar `TreatyIn` ke
-- `T_TREATY_REVISION` — bukan mengikuti Pega. Tabel ini dibangun tanpa
-- keempat medan itu, sehingga nilai yang layar Share tampilkan dan hitung
-- tidak punya tempat disimpan.
--
-- Pengukuran (baca saja, 7 Oktober 2026): salinan Pega di
-- `T_TREATY_SHARE.RNMSHARE` hanya ada pada 314 dari 1.079 kontrak Non-Prop.
-- Untuk kontrak lama, pengisian sekali dilakukan TERPISAH oleh
-- `modul/treatyin/alat/isi-akar-share-revisi.sql` (dijalankan tangan),
-- dari nilai yang Pega tulis sendiri — nol JSON, nol `M_TREATY_IN`.
--
-- ⚠️ `VARCHAR2(4000 CHAR)` seperti seluruh kolom pendaratan: nilainya teks
-- apa adanya, ditafsirkan services.
--
-- ⚠️ Pembaca `treatyin` (`BacaShareAkarRevisi`) TOLERAN terhadap kolom yang
-- belum ada, jadi layar tetap terbuka sebelum berkas ini dijalankan. Salinan
-- peta Adjustment baru menyebut keempat kolom SESUDAH berkas ini terpasang —
-- pembacanya memilih setiap kolom petanya.

ALTER TABLE {skema}.T_TREATY_REVISION ADD (
  RNMSHARE                         VARCHAR2(4000 CHAR),
  BROKERAGEPERCENT                 VARCHAR2(4000 CHAR),
  RNMSHAREACROSSTHEBOARD           VARCHAR2(4000 CHAR),
  RNMSHAREDEDUCTED                 VARCHAR2(4000 CHAR)
)
/

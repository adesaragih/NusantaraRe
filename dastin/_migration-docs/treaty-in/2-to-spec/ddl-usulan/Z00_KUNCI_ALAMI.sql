-- =====================================================================
-- Z00_KUNCI_ALAMI.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tidak ada instans Oracle yang terjangkau (L-3). Dinyatakan di sini, bukan
-- disembunyikan: setiap baris di bawah adalah bacaan, bukan hasil uji.
--
-- DIBANGKITKAN dari SPEC-MODEL-DATA.md §10, dari definisi yang SAMA dengan
-- 2-to-spec/KAMUS-KOLOM.md. Jangan disunting dengan tangan -- suntingan
-- tangan membuat keduanya dapat berbeda.
--
-- PRESISI -- DIPUTUSKAN 24 September 2026, butir KTV-A.
--   Angka presisi di bawah bukan lagi warisan yang menunggu diputuskan. Ia
--   DIPUTUSKAN, TANPA VERIFIKASI, dengan dasar dan syarat pembalikan tertulis
--   di KEPUTUSAN-TANPA-VERIFIKASI.md sec 7.
--   Kaidahnya: terlalu lebar di Oracle MURAH -- NUMBER dan VARCHAR2 disimpan
--   panjang-berubah; terlalu sempit MEMOTONG DATA, dan potongannya baru
--   ketahuan sesudah data masuk. Ongkosnya tidak setangkup, jadi sisi murahnya
--   yang diambil.
--   SYARAT PEMBALIKAN: sesi DDL boleh MEMPERSEMPIT, dan HANYA SEBELUM DATA
--   DIMUAT. Sesudah itu tidak -- baris yang tidak muat tidak punya tempat pergi.
--   Yang memberi angkanya: Uji AP (panjang teks sebenarnya) dan Uji AQ.
--
-- SATU BATAS YANG DINYATAKAN, bukan disembunyikan:
--   VARCHAR2(4000 CHAR) sah dideklarasikan, tetapi pada basis data AL32UTF8
--   batas BITA-nya tetap 4000. Teks 4.000 aksara yang memuat aksara berbita
--   ganda tetap dapat ditolak saat disisipkan, kecuali MAX_STRING_SIZE=EXTENDED.
--   Tidak ada instans yang dapat ditanyai (L-3), jadi ini dicatat, bukan diuji.
--
-- TIDAK ADA CREATE PROCEDURE, CREATE FUNCTION, maupun trigger pembawa aturan
-- bisnis di seluruh folder ini -- ADR-0056 (K-4).
-- TIDAK ADA DML.
-- =====================================================================


-- Dijalankan PALING AKHIR -- seluruh tabel harus sudah ada.
-- Setiap baris menyebut NOMOR INVARIANNYA. Constraint tanpa invarian tidak
-- ditulis di sini: ia tidak punya tempat untuk gagal.

-- INV-04
ALTER TABLE TREATY_MASUK.VERSI_KONTRAK ADD CONSTRAINT UQ_VERSI_KONTRAK UNIQUE (ID_KONTRAK, NOMOR_URUT_VERSI);
-- INV-05
ALTER TABLE TREATY_MASUK.LAYER ADD CONSTRAINT UQ_LAYER UNIQUE (ID_VERSI_KONTRAK, NOMOR_LAYER, BAGIAN_LAYER);
-- INV-06
ALTER TABLE TREATY_MASUK.DETAIL_PROPORSIONAL ADD CONSTRAINT UQ_DETAIL_PROPORSIONAL UNIQUE (ID_LAYER, ID_KELOMPOK_TREATY);
-- INV-07
ALTER TABLE TREATY_MASUK.MATA_UANG_KONTRAK ADD CONSTRAINT UQ_MATA_UANG_KONTRAK UNIQUE (ID_VERSI_KONTRAK, KODE_MATA_UANG);
-- INV-08
ALTER TABLE TREATY_MASUK.RETENSI_CEDANT ADD CONSTRAINT UQ_RETENSI_CEDANT UNIQUE (ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG);
-- INV-09
ALTER TABLE TREATY_MASUK.EGNPI ADD CONSTRAINT UQ_EGNPI UNIQUE (ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG);
-- INV-10
ALTER TABLE TREATY_MASUK.PERIODE_PELAPORAN ADD CONSTRAINT UQ_PERIODE_PELAPORAN UNIQUE (ID_VERSI_KONTRAK, PERIODE);
-- INV-12
ALTER TABLE TREATY_MASUK.TERMIN ADD CONSTRAINT UQ_TERMIN UNIQUE (ID_VERSI_KONTRAK, NOMOR_TERMIN, KODE_MATA_UANG);
-- INV-11
ALTER TABLE TREATY_MASUK.PERIODE_AKUMULASI ADD CONSTRAINT UQ_PERIODE_AKUMULASI UNIQUE (ID_VERSI_KONTRAK, PERIODE);
-- INV-13
ALTER TABLE TREATY_MASUK.SKALA_KOASURANSI ADD CONSTRAINT UQ_SKALA_KOASURANSI UNIQUE (ID_VERSI_KONTRAK, PERSEN_LIMIT);
-- INV-14
ALTER TABLE TREATY_MASUK.BATAS_PER_BAHAYA ADD CONSTRAINT UQ_BATAS_PER_BAHAYA UNIQUE (ID_VERSI_KONTRAK, ID_BAHAYA);
-- INV-15
ALTER TABLE TREATY_MASUK.POTONGAN ADD CONSTRAINT UQ_POTONGAN UNIQUE (ID_BAGIAN, ID_JENIS_POTONGAN);
-- INV-15
ALTER TABLE TREATY_MASUK.POTONGAN ADD CONSTRAINT UQ_POTONGAN_2 UNIQUE (ID_DETAIL_PROPORSIONAL, ID_JENIS_POTONGAN);
-- INV-16
ALTER TABLE TREATY_MASUK.PENYEBARAN ADD CONSTRAINT UQ_PENYEBARAN UNIQUE (ID_BAGIAN, ID_JENIS_REASURANSI);
-- INV-16
ALTER TABLE TREATY_MASUK.PENYEBARAN ADD CONSTRAINT UQ_PENYEBARAN_2 UNIQUE (ID_DETAIL_PROPORSIONAL, ID_JENIS_REASURANSI);
-- INV-64
ALTER TABLE TREATY_MASUK.BAGIAN ADD CONSTRAINT UQ_BAGIAN UNIQUE (ID_LAYER);
-- INV-65
ALTER TABLE TREATY_MASUK.RINCIAN_PENYEBARAN ADD CONSTRAINT UQ_RINCIAN_PENYEBARAN UNIQUE (ID_PENYEBARAN, ID_JENIS_REASURANSI);
-- INV-66
ALTER TABLE TREATY_MASUK.PORTOFOLIO ADD CONSTRAINT UQ_PORTOFOLIO UNIQUE (ID_VERSI_KONTRAK, ARAH_PORTOFOLIO, JENIS_PORTOFOLIO);
-- INV-67
ALTER TABLE TREATY_MASUK.DOKUMEN_KONTRAK ADD CONSTRAINT UQ_DOKUMEN_KONTRAK UNIQUE (ID_VERSI_KONTRAK, ID_DOKUMEN);
-- INV-68
ALTER TABLE TREATY_MASUK.MATA_UANG ADD CONSTRAINT UQ_MATA_UANG UNIQUE (KODE);
-- INV-68
ALTER TABLE TREATY_MASUK.JENIS_POTONGAN ADD CONSTRAINT UQ_JENIS_POTONGAN UNIQUE (KODE);
-- INV-68
ALTER TABLE TREATY_MASUK.JENIS_REASURANSI ADD CONSTRAINT UQ_JENIS_REASURANSI UNIQUE (KODE);
-- INV-68
ALTER TABLE TREATY_MASUK.BAHAYA ADD CONSTRAINT UQ_BAHAYA UNIQUE (KODE);
-- INV-68
ALTER TABLE TREATY_MASUK.KELOMPOK_TREATY ADD CONSTRAINT UQ_KELOMPOK_TREATY UNIQUE (KODE);
-- INV-68
ALTER TABLE TREATY_MASUK.KELAS_BISNIS ADD CONSTRAINT UQ_KELAS_BISNIS UNIQUE (KODE);

-- ---------------------------------------------------------------------
-- INV-08, INV-09, INV-12 -- DIBEBASKAN 24 September 2026
-- ---------------------------------------------------------------------
--   Ketiganya sempat TIDAK dapat dikompilasi: kunci alaminya menyebut
--   "kode mata uang" dan kolomnya tidak ada, karena mata uang melebur ke
--   dalam paket uang. P-8 golongan B memutuskan mata uang tinggal pada
--   BARIS YANG SAMA -- dan itu memang bentuknya di sistem lama.
--   Sejak itu ketiganya tertulis sebagai UNIQUE di atas.
--
--   Catatan yang TIDAK boleh hilang: menulis UNIQUE ini TANPA mata uang
--   bukan jalan tengah -- ia mengubah artinya menjadi "dilarang dua baris
--   bermata uang berbeda", dan itu MENOLAK DATA YANG SAH.
--
-- KUNCI ALAMI YANG SENGAJA BUKAN CONSTRAINT -- bukan kelalaian:
--   KONTRAK              : cedant + SoB + periode + sifat proporsi.
--                          ADR-0040 §2 MEMPERINGATKAN, tidak melarang. §10.1
--                          menyatakannya, dan ketiadaan nomor INV-nya disengaja.
--   NILAI_PENYEBARAN     : kode mata uang di dalam rincian -- BELUM BERNOMOR;
--                          kedudukan entitasnya sendiri ditangguhkan Uji AD
--                          (§10.23c butir 3). SPEC-INVARIAN.md §7.3.
--   PEMULIHAN_LIMIT      : belum bernomor; entitasnya baru diterima §10.23c.
--   CATATAN_PERSETUJUAN  : TIDAK ADA, dan itu keputusan -- peristiwa yang sama
--   JEJAK_PERUBAHAN        dapat terjadi dua kali pada versi yang sama.
--   PERISTIWA_KONTRAK      §10.20, §10.21, §10.21a. SPEC-INVARIAN.md §7.4.


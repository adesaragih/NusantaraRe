-- =====================================================================
-- 12_DOKUMEN_KONTRAK.sql
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


CREATE TABLE TREATY_MASUK.DOKUMEN_KONTRAK (
    ID_DOKUMEN_KONTRAK   NUMBER(19)             NOT NULL,
    ID_VERSI_KONTRAK     NUMBER(19)             NOT NULL,
    ID_DOKUMEN           NUMBER(19)             NOT NULL,
    JENIS_DOKUMEN        VARCHAR2(40 CHAR),
    TANGGAL_LAMPIR       DATE                   NOT NULL
);

ALTER TABLE TREATY_MASUK.DOKUMEN_KONTRAK ADD CONSTRAINT PK_DOKUMEN_KONTRAK PRIMARY KEY (ID_DOKUMEN_KONTRAK);
ALTER TABLE TREATY_MASUK.DOKUMEN_KONTRAK ADD CONSTRAINT FK_DOKUMEN_KONTRAK_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES TREATY_MASUK.VERSI_KONTRAK (ID_VERSI_KONTRAK);

GRANT SELECT, INSERT, UPDATE, DELETE ON TREATY_MASUK.DOKUMEN_KONTRAK TO TREATY_MASUK_APP;

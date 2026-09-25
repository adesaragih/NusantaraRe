-- =====================================================================
-- 25_DETAIL_PROPORSIONAL.sql
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


CREATE TABLE TREATY_MASUK.DETAIL_PROPORSIONAL (
    ID_DETAIL_PROPORSIONAL     NUMBER(19)             NOT NULL,
    ID_LAYER                   NUMBER(19)             NOT NULL,
    ID_KELOMPOK_TREATY         NUMBER(19)             NOT NULL,
    JENIS_TREATY               VARCHAR2(40 CHAR)      NOT NULL,
    PERSEN_QUOTA_SHARE         NUMBER(38,20),
    JUMLAH_LINES_SURPLUS       NUMBER(9),
    PERSEN_KOMISI_KOTOR        NUMBER(38,20),
    PERSEN_KOMISI_BERSIH       NUMBER(38,20),
    PERSEN_CADANGAN_PREMI      NUMBER(38,20),
    CADANGAN_PREMI             NUMBER(38,20),
    PERSEN_KAPASITAS_SURPLUS   NUMBER(38,20),
    ID_SUSUNAN_RETRO           NUMBER(19)
);

ALTER TABLE TREATY_MASUK.DETAIL_PROPORSIONAL ADD CONSTRAINT PK_DETAIL_PROPORSIONAL PRIMARY KEY (ID_DETAIL_PROPORSIONAL);
ALTER TABLE TREATY_MASUK.DETAIL_PROPORSIONAL ADD CONSTRAINT FK_DETAIL_PROPORSIONAL_1 FOREIGN KEY (ID_LAYER)
    REFERENCES TREATY_MASUK.LAYER (ID_LAYER);
ALTER TABLE TREATY_MASUK.DETAIL_PROPORSIONAL ADD CONSTRAINT FK_DETAIL_PROPORSIONAL_2 FOREIGN KEY (ID_KELOMPOK_TREATY)
    REFERENCES TREATY_MASUK.KELOMPOK_TREATY (ID_KELOMPOK_TREATY);

GRANT SELECT, INSERT, UPDATE, DELETE ON TREATY_MASUK.DETAIL_PROPORSIONAL TO TREATY_MASUK_APP;

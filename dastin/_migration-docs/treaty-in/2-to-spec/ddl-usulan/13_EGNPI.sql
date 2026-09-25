-- =====================================================================
-- 13_EGNPI.sql
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


CREATE TABLE TREATY_MASUK.EGNPI (
    ID_EGNPI             NUMBER(19)             NOT NULL,
    ID_VERSI_KONTRAK     NUMBER(19)             NOT NULL,
    ID_KELOMPOK_TREATY   NUMBER(19)             NOT NULL,
    ID_KELAS_BISNIS      NUMBER(19),
    NILAI_EGNPI          NUMBER(38,20)          NOT NULL,
    KODE_MATA_UANG       VARCHAR2(1000 CHAR)    NOT NULL,
    TANGGAL_BERLAKU      DATE,
    PROPORSI             NUMBER(38,20),
    CATATAN              VARCHAR2(1000 CHAR)
);

ALTER TABLE TREATY_MASUK.EGNPI ADD CONSTRAINT PK_EGNPI PRIMARY KEY (ID_EGNPI);
ALTER TABLE TREATY_MASUK.EGNPI ADD CONSTRAINT FK_EGNPI_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES TREATY_MASUK.VERSI_KONTRAK (ID_VERSI_KONTRAK);
ALTER TABLE TREATY_MASUK.EGNPI ADD CONSTRAINT FK_EGNPI_2 FOREIGN KEY (ID_KELOMPOK_TREATY)
    REFERENCES TREATY_MASUK.KELOMPOK_TREATY (ID_KELOMPOK_TREATY);
ALTER TABLE TREATY_MASUK.EGNPI ADD CONSTRAINT FK_EGNPI_3 FOREIGN KEY (ID_KELAS_BISNIS)
    REFERENCES TREATY_MASUK.KELAS_BISNIS (ID_KELAS_BISNIS);

-- ---------------------------------------------------------------------
-- CONSTRAINT YANG SENGAJA TIDAK DIPASANG -- ini PERNYATAAN KEPUTUSAN,
-- bukan TODO. Jangan dipasang tanpa jawaban yang disebut di baris DITAGIH.
-- ---------------------------------------------------------------------
-- NILAI_EGNPI
--   APA        : tidak ada CHECK maupun NOT NULL yang memasangkan paket uang
--                ini dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40).
--   KENAPA     : tingkat pencatatannya BELUM DITENTUKAN -- sapuan penulisnya
--                tidak memisahkan apa pun. Constraint yang dipasang sekarang
--                menuntut nilai yang migrasi tidak tahu cara mengisinya.
--   AKIBAT     : kolom ini menerima angka tanpa menyatakan angkanya untuk
--                SELURUH TREATY atau untuk BAGIAN NuRe. Pembaca laporan tidak
--                dapat mengetahuinya dari basis data.
--   DILIHAT DI : 2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md §3
--   DITAGIH    : ketika jawaban T-4 kembali dari teknik treaty.
--

GRANT SELECT, INSERT, UPDATE, DELETE ON TREATY_MASUK.EGNPI TO TREATY_MASUK_APP;

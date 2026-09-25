-- =====================================================================
-- 15_LAYER.sql
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


CREATE TABLE TREATY_MASUK.LAYER (
    ID_LAYER                       NUMBER(19)             NOT NULL,
    ID_VERSI_KONTRAK               NUMBER(19)             NOT NULL,
    NOMOR_LAYER                    NUMBER(9)              NOT NULL,
    BAGIAN_LAYER                   NUMBER(9),
    JENIS_LAYER                    VARCHAR2(40 CHAR)      NOT NULL,
    JENIS_BAGIAN_LAYER             VARCHAR2(40 CHAR),
    CAKUPAN                        VARCHAR2(1000 CHAR),
    LIMIT                          NUMBER(38,20)          NOT NULL,
    KODE_MATA_UANG                 VARCHAR2(1000 CHAR)    NOT NULL,
    DEDUCTIBLE                     NUMBER(38,20)          NOT NULL,
    MATA_UANG_DEDUCTIBLE           VARCHAR2(1000 CHAR),
    PERSEN_PENYESUAIAN             NUMBER(38,20),
    PERSEN_MINIMUM_DEPOSIT         NUMBER(38,20),
    PORSI_PEMULIHAN_LIMIT          NUMBER(38,20),
    TARIF_PREMI_PEMULIHAN          NUMBER(38,20),
    LIMIT_AGREGAT                  NUMBER(38,20),
    MATA_UANG_LIMIT_AGREGAT        VARCHAR2(1000 CHAR),
    MDP                            NUMBER(38,20),
    MDP_MINIMUM                    NUMBER(38,20),
    PERSEN_MDP_MINIMUM             NUMBER(38,20),
    MDP_DIGABUNG                   VARCHAR2(40 CHAR)      NOT NULL,
    TANPA_HITUNG_PREMI_PEMULIHAN   VARCHAR2(40 CHAR)      NOT NULL,
    DEDUCTIBLE_KEDUA               NUMBER(38,20),
    PERSEN_ROL                     NUMBER(38,20)
);

ALTER TABLE TREATY_MASUK.LAYER ADD CONSTRAINT PK_LAYER PRIMARY KEY (ID_LAYER);
ALTER TABLE TREATY_MASUK.LAYER ADD CONSTRAINT FK_LAYER_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES TREATY_MASUK.VERSI_KONTRAK (ID_VERSI_KONTRAK);

-- ---------------------------------------------------------------------
-- PAKET UANG GOLONGAN C -- KOLOM MATA UANG DIPASANG, BOLEH KOSONG
-- ---------------------------------------------------------------------
--   APA        : DEDUCTIBLE, LIMIT_AGREGAT
--                memperoleh kolom mata uang BERNAMA di sebelahnya,
--                nullable. Namanya menyebut paketnya (MATA_UANG_<paket>)
--                dan BUKAN "KODE_MATA_UANG" polos, sebab baris ini memuat
--                lebih dari satu paket uang atau sudah punya kolom mata
--                uang lain -- satu nama untuk dua arti adalah cacat yang
--                ditanam dengan tangan sendiri.
--   KENAPA     : sistem lama TIDAK merekam mata uang untuk besaran ini
--                sama sekali. Kolom nullable yang tak terpakai murah;
--                menambahkannya SESUDAH data masuk mahal, dan tidak ada
--                sumber untuk mengisi baris warisannya.
--   AKIBAT     : kolomnya akan KOSONG pada seluruh baris hasil migrasi.
--                Itu BUKAN kegagalan -- itu keadaan yang benar, dan
--                pembacanya harus tahu bahwa kosong berarti "tidak
--                pernah dicatat", bukan "belum diisi".
--   DILIHAT DI : KEPUTUSAN-TANPA-VERIFIKASI.md sec 7, butir KTV-C
--   DITAGIH    : jawaban T-6 dari teknik treaty. Ia MENYEMPITKAN, tidak
--                lagi memblokir: bila T-6 menyatakan besaran ini selalu
--                dalam mata uang kontrak, kolomnya DICABUT SEBELUM DATA
--                DIMUAT dan bacaannya diserahkan ke KODE_MATA_UANG_KONTRAK.
--

-- ---------------------------------------------------------------------
-- CONSTRAINT YANG SENGAJA TIDAK DIPASANG -- ini PERNYATAAN KEPUTUSAN,
-- bukan TODO. Jangan dipasang tanpa jawaban yang disebut di baris DITAGIH.
-- ---------------------------------------------------------------------
-- LIMIT_AGREGAT
--   APA        : tidak ada CHECK maupun NOT NULL yang memasangkan paket uang
--                ini dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40).
--   KENAPA     : tingkat pencatatannya BELUM DITENTUKAN -- sapuan penulisnya
--                tidak memisahkan apa pun. Constraint yang dipasang sekarang
--                menuntut nilai yang migrasi tidak tahu cara mengisinya.
--   AKIBAT     : kolom ini menerima angka tanpa menyatakan angkanya untuk
--                SELURUH TREATY atau untuk BAGIAN NuRe. Pembaca laporan tidak
--                dapat mengetahuinya dari basis data.
--   DILIHAT DI : 2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md §3
--   DITAGIH    : ketika jawaban T-3 kembali dari teknik treaty.
--

GRANT SELECT, INSERT, UPDATE, DELETE ON TREATY_MASUK.LAYER TO TREATY_MASUK_APP;

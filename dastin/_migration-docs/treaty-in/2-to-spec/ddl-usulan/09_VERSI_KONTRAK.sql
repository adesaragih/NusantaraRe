-- =====================================================================
-- 09_VERSI_KONTRAK.sql
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


CREATE TABLE TREATY_MASUK.VERSI_KONTRAK (
    ID_VERSI_KONTRAK               NUMBER(19)             NOT NULL,
    ID_KONTRAK                     NUMBER(19)             NOT NULL,
    NOMOR_URUT_VERSI               NUMBER(9)              NOT NULL,
    KEADAAN_SIKLUS_HIDUP           VARCHAR2(40 CHAR)      NOT NULL,
    KEADAAN_WARISAN_ASLI           VARCHAR2(1000 CHAR),
    JENIS_ADDENDUM                 VARCHAR2(40 CHAR),
    TANGGAL_BERLAKU_ADDENDUM       DATE,
    NAMA_KONTRAK                   VARCHAR2(1000 CHAR)    NOT NULL,
    LINGKUP_WILAYAH                VARCHAR2(1000 CHAR),
    KELAS_BISNIS_KONTRAK           NUMBER(19),
    KODE_MATA_UANG_KONTRAK         NUMBER(19)             NOT NULL,
    ID_REASURADUR_PEMIMPIN         NUMBER(19),
    ID_KETUA_TREATY                NUMBER(19),
    PERSEN_BAGIAN_NURE             NUMBER(38,20)          NOT NULL,
    BAGIAN_NURE_SERAGAM            VARCHAR2(40 CHAR)      NOT NULL,
    PERSEN_BAGIAN_NURE_DIPOTONG    NUMBER(38,20),
    PERSEN_BROKERAGE               NUMBER(38,20),
    PERSEN_BAGIAN_FAKULTATIF       NUMBER(38,20),
    PERSEN_BROKERAGE_FAKULTATIF    NUMBER(38,20),
    PENGECUALIAN                   VARCHAR2(4000 CHAR),
    KETENTUAN_KHUSUS               VARCHAR2(4000 CHAR),
    KETERANGAN                     VARCHAR2(4000 CHAR),
    CATATAN                        VARCHAR2(4000 CHAR),
    CATATAN_BORDEREAUX             VARCHAR2(4000 CHAR),
    MEMAKAI_BORDEREAUX             VARCHAR2(40 CHAR)      NOT NULL,
    CARA_PEMBUKUAN                 VARCHAR2(40 CHAR)      NOT NULL,
    CARA_PEMBUKUAN_XOL             VARCHAR2(40 CHAR),
    PERIODE_PELAPORAN_KONTRAK      VARCHAR2(40 CHAR),
    SELANG_PELAPORAN               NUMBER(9),
    TANGGAL_MULAI_PELAPORAN        DATE,
    TANGGAL_AKHIR_PELAPORAN        DATE,
    HARI_BATAS_PENYERAHAN          NUMBER(9),
    HARI_BATAS_KONFIRMASI          NUMBER(9),
    HARI_BATAS_PELUNASAN           NUMBER(9),
    HARI_PENGINGAT                 NUMBER(9),
    PERIODE_AKUMULASI_KONTRAK      VARCHAR2(40 CHAR),
    JUMLAH_TERMIN                  NUMBER(9),
    MEMAKAI_PRORATA                VARCHAR2(40 CHAR)      NOT NULL,
    BATAS_MAKSIMUM_KELOMPOK        NUMBER(38,20),
    MATA_UANG_BATAS_KELOMPOK       VARCHAR2(1000 CHAR),
    BATAS_MAKSIMUM_NON_KELOMPOK    NUMBER(38,20),
    MATA_UANG_BATAS_NON_KELOMPOK   VARCHAR2(1000 CHAR),
    BATAS_PILIHAN                  NUMBER(38,20),
    MATA_UANG_BATAS_PILIHAN        VARCHAR2(1000 CHAR),
    RETRO_BERGANDA                 VARCHAR2(40 CHAR)      NOT NULL,
    ID_DAFTAR_RETRO                NUMBER(19),
    SIFAT_MATERIAL_ADDENDUM        VARCHAR2(40 CHAR),
    ID_DOKUMEN_ADDENDUM            NUMBER(19)
);

ALTER TABLE TREATY_MASUK.VERSI_KONTRAK ADD CONSTRAINT PK_VERSI_KONTRAK PRIMARY KEY (ID_VERSI_KONTRAK);
ALTER TABLE TREATY_MASUK.VERSI_KONTRAK ADD CONSTRAINT FK_VERSI_KONTRAK_1 FOREIGN KEY (ID_KONTRAK)
    REFERENCES TREATY_MASUK.KONTRAK (ID_KONTRAK);
ALTER TABLE TREATY_MASUK.VERSI_KONTRAK ADD CONSTRAINT FK_VERSI_KONTRAK_2 FOREIGN KEY (ID_DOKUMEN_ADDENDUM)
    REFERENCES TREATY_MASUK.DOKUMEN_ADDENDUM (ID_DOKUMEN_ADDENDUM);

-- ---------------------------------------------------------------------
-- PAKET UANG GOLONGAN C -- KOLOM MATA UANG DIPASANG, BOLEH KOSONG
-- ---------------------------------------------------------------------
--   APA        : BATAS_MAKSIMUM_KELOMPOK, BATAS_MAKSIMUM_NON_KELOMPOK, BATAS_PILIHAN
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
-- BATAS_MAKSIMUM_KELOMPOK
--   APA        : tidak ada CHECK maupun NOT NULL yang memasangkan paket uang
--                ini dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40).
--   KENAPA     : tingkat pencatatannya BELUM DITENTUKAN -- sapuan penulisnya
--                tidak memisahkan apa pun. Constraint yang dipasang sekarang
--                menuntut nilai yang migrasi tidak tahu cara mengisinya.
--   AKIBAT     : kolom ini menerima angka tanpa menyatakan angkanya untuk
--                SELURUH TREATY atau untuk BAGIAN NuRe. Pembaca laporan tidak
--                dapat mengetahuinya dari basis data.
--   DILIHAT DI : 2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md §3
--   DITAGIH    : ketika jawaban T-1 kembali dari teknik treaty.
--
-- BATAS_MAKSIMUM_NON_KELOMPOK
--   APA        : tidak ada CHECK maupun NOT NULL yang memasangkan paket uang
--                ini dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40).
--   KENAPA     : tingkat pencatatannya BELUM DITENTUKAN -- sapuan penulisnya
--                tidak memisahkan apa pun. Constraint yang dipasang sekarang
--                menuntut nilai yang migrasi tidak tahu cara mengisinya.
--   AKIBAT     : kolom ini menerima angka tanpa menyatakan angkanya untuk
--                SELURUH TREATY atau untuk BAGIAN NuRe. Pembaca laporan tidak
--                dapat mengetahuinya dari basis data.
--   DILIHAT DI : 2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md §3
--   DITAGIH    : ketika jawaban T-1 kembali dari teknik treaty.
--
-- BATAS_PILIHAN
--   APA        : tidak ada CHECK maupun NOT NULL yang memasangkan paket uang
--                ini dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40).
--   KENAPA     : tingkat pencatatannya BELUM DITENTUKAN -- sapuan penulisnya
--                tidak memisahkan apa pun. Constraint yang dipasang sekarang
--                menuntut nilai yang migrasi tidak tahu cara mengisinya.
--   AKIBAT     : kolom ini menerima angka tanpa menyatakan angkanya untuk
--                SELURUH TREATY atau untuk BAGIAN NuRe. Pembaca laporan tidak
--                dapat mengetahuinya dari basis data.
--   DILIHAT DI : 2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md §3
--   DITAGIH    : ketika jawaban T-1 kembali dari teknik treaty.
--

GRANT SELECT, INSERT, UPDATE, DELETE ON TREATY_MASUK.VERSI_KONTRAK TO TREATY_MASUK_APP;

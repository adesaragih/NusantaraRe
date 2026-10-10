-- Cabang penyebaran - PENYEBARAN, RINCIAN_PENYEBARAN, NILAI_PENYEBARAN
--
-- Migrasi tiket 38 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `2-to-spec/ddl-usulan/{38_PENYEBARAN,40_RINCIAN_PENYEBARAN,
-- 41_NILAI_PENYEBARAN}.sql`, dan `KAMUS-KOLOM.md` §10.6-§10.8.
--
-- Tiga tabel dalam satu langkah: `RINCIAN_PENYEBARAN` tidak dapat berdiri
-- tanpa `PENYEBARAN`, dan `NILAI_PENYEBARAN` tidak tanpa keduanya.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS, seluruhnya dari `ERD.md` §2.5 dan §2.7
-- ---------------------------------------------------------------------
--   BAGIAN               1--o<  PENYEBARAN          [ikut hapus]   §2.5
--   DETAIL_PROPORSIONAL  1--o<  PENYEBARAN          [ikut hapus]   §2.5
--   PENYEBARAN           1--<   RINCIAN_PENYEBARAN  [ikut hapus]   §2.5
--   RINCIAN_PENYEBARAN   1--<   NILAI_PENYEBARAN    [ikut hapus]   §2.5
--   PENYEBARAN           >o--1  JENIS_REASURANSI    [tolak]        §2.7
--   RINCIAN_PENYEBARAN   >o--1  JENIS_REASURANSI    [tolak]        §2.7
--
--   Rujukan tabel acuan `tolak`, diwujudkan dengan tidak menulis klausa
--   `ON DELETE`: jenis reasuransi yang sudah dipakai sebuah penyebaran tidak
--   boleh hilang dari tabel acuannya.
--
-- ---------------------------------------------------------------------
-- DUA PELEKATAN, SATU KONSEP - dan `CHECK` yang menjaganya
-- ---------------------------------------------------------------------
--   Penyebaran menggantung pada `BAGIAN` ATAU pada `DETAIL_PROPORSIONAL`,
--   dan TEPAT SATU dari keduanya terisi. Itu invarian, bukan pilihan bebas;
--   bentuknya mengikuti `POTONGAN` di migrasi 418, yang sudah memakai pola
--   yang sama beserta `CK_POTONGAN_INDUK`-nya.
--
--   ⛔ `CHECK` ditulis di `ALTER TABLE` TERSENDIRI, bukan di dalam
--   `CREATE TABLE`. Penjaga lintas-aplikasi mengurai tiap baris dalam
--   `CREATE TABLE` sebagai kolom, dan baris lanjutan berawalan `OR` terbaca
--   sebagai kolom bernama "OR". Ini sudah pernah memakan korban.
--
--   ⚠️ `ADR-0038` melarang `CHECK` yang MENGENUMERASI NILAI. Yang ini
--   menyatakan BENTUK BARIS - berapa induk yang terisi - dan itu golongan
--   yang berbeda. Ia bernama, jadi ia dapat diadili.
--
-- ---------------------------------------------------------------------
-- KUNCI ALAMI - DIPASANG, sebab ketiganya SUDAH BERNOMOR
-- ---------------------------------------------------------------------
--   Berbeda dari enam tabel ronde 6 yang kuncinya ditahan: ketiga `UNIQUE`
--   di bawah punya nomor invariannya di `ddl-usulan/Z00_KUNCI_ALAMI.sql`.
--
--     INV-16  UQ_PENYEBARAN            (ID_BAGIAN, ID_JENIS_REASURANSI)
--     INV-16  UQ_PENYEBARAN_2          (ID_DETAIL_PROPORSIONAL, ID_JENIS_REASURANSI)
--     INV-65  UQ_RINCIAN_PENYEBARAN    (ID_PENYEBARAN, ID_JENIS_REASURANSI)
--
--   ⚠️ `NILAI_PENYEBARAN` TIDAK mendapat `UNIQUE`: Z00 tidak menyebutnya,
--   jadi kunci alaminya belum bernomor. Ia masuk tagihan
--   `SPEC-MODEL-DATA.md` §13 bersama enam yang lain - bukan dipasang
--   diam-diam.
--
-- ⚠️ `KODE_MATA_UANG` TEKS, bukan kunci asing - mengikuti kelima paket uang
-- yang sudah berdiri (`NILAI_MDP` dan saudaranya, migrasi 413 dan 417).
-- `ERD.md` §2.7 menyatakan relasinya `tolak`; di repo ini kolomnya belum
-- dinormalkan, dan menormalkan satu tabel saja membuat dua bentuk untuk satu
-- fakta. Dicatat, bukan ditambal sendirian.
--
-- ⛔ INV-32, INV-33, INV-52 TETAP TIDAK DIPASANG - ketiganya menuntut
-- materialized view `REFRESH ON COMMIT`, dan `F-13` menyatakan modul ini
-- belum punya pemantau kebasiannya.
--
-- Presisi: `NUMBER(38,20)` di `ddl-usulan` dipersempit menjadi `NUMBER(38,8)`
-- (butir `KTV-A`, boleh dipersempit SEBELUM data dimuat).
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE SEQUENCE {skema}.SEQ_TRIN_PENYEBARAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_RINCIAN_PENYEBARAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_NILAI_PENYEBARAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.PENYEBARAN (
  ID_PENYEBARAN              NUMBER(19)   NOT NULL,
  ID_BAGIAN                  NUMBER(19),
  ID_DETAIL_PROPORSIONAL     NUMBER(19),
  ID_JENIS_REASURANSI        NUMBER(19)   NOT NULL,
  ID_JENIS_REASURANSI_INDUK  NUMBER(19),
  PERSEN_PENYEBARAN          NUMBER(38,8) NOT NULL,
  CONSTRAINT PK_PENYEBARAN PRIMARY KEY (ID_PENYEBARAN),
  CONSTRAINT UQ_PENYEBARAN UNIQUE (ID_BAGIAN, ID_JENIS_REASURANSI),
  CONSTRAINT UQ_PENYEBARAN_2 UNIQUE (ID_DETAIL_PROPORSIONAL, ID_JENIS_REASURANSI),
  CONSTRAINT FK_PENYEBARAN_1 FOREIGN KEY (ID_BAGIAN)
    REFERENCES {skema}.BAGIAN (ID_BAGIAN) ON DELETE CASCADE,
  CONSTRAINT FK_PENYEBARAN_2 FOREIGN KEY (ID_DETAIL_PROPORSIONAL)
    REFERENCES {skema}.DETAIL_PROPORSIONAL (ID_DETAIL_PROPORSIONAL) ON DELETE CASCADE,
  CONSTRAINT FK_PENYEBARAN_3 FOREIGN KEY (ID_JENIS_REASURANSI)
    REFERENCES {skema}.JENIS_REASURANSI (ID_JENIS_REASURANSI)
)
/
ALTER TABLE {skema}.PENYEBARAN ADD CONSTRAINT CK_PENYEBARAN_INDUK CHECK ((ID_BAGIAN IS NOT NULL AND ID_DETAIL_PROPORSIONAL IS NULL) OR (ID_BAGIAN IS NULL AND ID_DETAIL_PROPORSIONAL IS NOT NULL))
/
CREATE TABLE {skema}.RINCIAN_PENYEBARAN (
  ID_RINCIAN_PENYEBARAN  NUMBER(19)   NOT NULL,
  ID_PENYEBARAN          NUMBER(19)   NOT NULL,
  ID_JENIS_REASURANSI    NUMBER(19)   NOT NULL,
  PERSEN_RINCIAN         NUMBER(38,8) NOT NULL,
  CONSTRAINT PK_RINCIAN_PENYEBARAN PRIMARY KEY (ID_RINCIAN_PENYEBARAN),
  CONSTRAINT UQ_RINCIAN_PENYEBARAN UNIQUE (ID_PENYEBARAN, ID_JENIS_REASURANSI),
  CONSTRAINT FK_RINCIAN_PENYEBARAN_1 FOREIGN KEY (ID_PENYEBARAN)
    REFERENCES {skema}.PENYEBARAN (ID_PENYEBARAN) ON DELETE CASCADE,
  CONSTRAINT FK_RINCIAN_PENYEBARAN_2 FOREIGN KEY (ID_JENIS_REASURANSI)
    REFERENCES {skema}.JENIS_REASURANSI (ID_JENIS_REASURANSI)
)
/
CREATE TABLE {skema}.NILAI_PENYEBARAN (
  ID_NILAI_PENYEBARAN    NUMBER(19)          NOT NULL,
  ID_RINCIAN_PENYEBARAN  NUMBER(19)          NOT NULL,
  NILAI                  NUMBER(38,8)        NOT NULL,
  KODE_MATA_UANG         VARCHAR2(1000 CHAR) NOT NULL,
  CONSTRAINT PK_NILAI_PENYEBARAN PRIMARY KEY (ID_NILAI_PENYEBARAN),
  CONSTRAINT FK_NILAI_PENYEBARAN_1 FOREIGN KEY (ID_RINCIAN_PENYEBARAN)
    REFERENCES {skema}.RINCIAN_PENYEBARAN (ID_RINCIAN_PENYEBARAN) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_PENYEBARAN_JENIS_REAS ON {skema}.PENYEBARAN (ID_JENIS_REASURANSI)
/
CREATE INDEX {skema}.IX_RINCIAN_PENY_JENIS ON {skema}.RINCIAN_PENYEBARAN (ID_JENIS_REASURANSI)
/
CREATE INDEX {skema}.IX_NILAI_PENY_RINCIAN ON {skema}.NILAI_PENYEBARAN (ID_RINCIAN_PENYEBARAN)
/

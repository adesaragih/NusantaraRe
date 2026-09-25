-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Dokumen dirujuk, tidak disimpan. Berkasnya di Google Cloud Storage.

CREATE TABLE KLAIMNP.DOKUMEN_KLAIM
(
  ID_DOKUMEN         NUMBER(19) NOT NULL,
  ID_KLAIM           NUMBER(19) NOT NULL,
  NOMOR_URUT         NUMBER(9) NOT NULL,
  NAMA_BERKAS        VARCHAR2(255 CHAR) NOT NULL,
  JENIS_MEDIA        VARCHAR2(128 CHAR),
  URL                VARCHAR2(2000 CHAR),
  URL_KEDALUWARSA    TIMESTAMP(6),
  DIBUAT_OLEH        VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA   VARCHAR2(128 CHAR),
  DIBUAT_PADA        TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH        VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA   VARCHAR2(128 CHAR),
  DIUBAH_PADA        TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.DOKUMEN_KLAIM ADD CONSTRAINT PK_DOKUMEN_KLAIM PRIMARY KEY (ID_DOKUMEN);
ALTER TABLE KLAIMNP.DOKUMEN_KLAIM ADD CONSTRAINT UQ_DOKUMEN_KLAIM_1 UNIQUE (ID_KLAIM, NOMOR_URUT);
ALTER TABLE KLAIMNP.DOKUMEN_KLAIM ADD CONSTRAINT CK_DOKUMEN_KLAIM_1 CHECK ((URL IS NULL AND URL_KEDALUWARSA IS NULL) OR (URL IS NOT NULL AND URL_KEDALUWARSA IS NOT NULL));
ALTER TABLE KLAIMNP.DOKUMEN_KLAIM ADD CONSTRAINT FK_DOKUMEN_KLAIM_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);

-- ---------------------------------------------------------------------
-- PEMBANGKIT PENGENAL — ditambahkan 19 September 2026
-- ---------------------------------------------------------------------
-- Sampai tanggal ini, SELURUH 22 tabel punya kunci primer NUMBER(19) dan
-- TIDAK ADA SATU PUN yang menyatakan dari mana nilainya datang: nol
-- IDENTITY, nol SEQUENCE, nol DEFAULT di 36 berkas. Hak CREATE SEQUENCE
-- sudah diberikan di 00_SKEMA_DAN_AKUN.sql — jadi niatnya ada; objeknya
-- tidak pernah dibuat.
--
-- KENAPA SEQUENCE, BUKAN "GENERATED AS IDENTITY"
-- Alasan yang sama yang menetapkan batas nama 30 byte: pilih bentuk yang
-- sah di SETIAP versi. IDENTITY baru ada sejak 12.1; sequence sah jauh
-- sebelumnya. Versi instance tujuan masih menunggu REQ-032, dan REQ-032
-- sudah diturunkan menjadi verifikasi justru dengan janji bahwa jawabannya
-- TIDAK MENGUBAH APA PUN. Memakai IDENTITY akan membatalkan janji itu.
--
-- Nilai diminta pemanggil lewat NEXTVAL. Itu konsisten dengan ADR-0017:
-- satu pintu tulis, dan pintu itu yang meminta nomornya.
--
-- Penamaan mengikuti aturan 2 dan 3 (SPEC bagian 16): SQ_<tabel>, dengan
-- SQ_ sebagai awalan peran — sekelas PK_, UQ_, FK_, CK_, IX_, V_.

CREATE SEQUENCE KLAIMNP.SQ_DOKUMEN_KLAIM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_DOKUMEN_KLAIM TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.DOKUMEN_KLAIM TO KLAIMNP_APP;

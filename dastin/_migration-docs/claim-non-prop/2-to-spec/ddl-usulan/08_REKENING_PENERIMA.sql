-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Salinan keadaan rekening pada saat pembayaran. Pengecualian sah ADR-0023.

CREATE TABLE KLAIMNP.REKENING_PENERIMA
(
  ID_REKENING        NUMBER(19) NOT NULL,
  ID_KLAIM           NUMBER(19) NOT NULL,
  NAMA_PENERIMA      VARCHAR2(255 CHAR) NOT NULL,
  NAMA_BANK          VARCHAR2(128 CHAR) NOT NULL,
  CABANG_BANK        VARCHAR2(128 CHAR),
  NOMOR_REKENING     VARCHAR2(64 CHAR) NOT NULL,
  KODE_SWIFT         VARCHAR2(16 CHAR),
  MATA_UANG          VARCHAR2(3 CHAR) NOT NULL,
  DIBUAT_OLEH        VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA   VARCHAR2(128 CHAR),
  DIBUAT_PADA        TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH        VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA   VARCHAR2(128 CHAR),
  DIUBAH_PADA        TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.REKENING_PENERIMA ADD CONSTRAINT PK_REKENING_PENERIMA PRIMARY KEY (ID_REKENING);
ALTER TABLE KLAIMNP.REKENING_PENERIMA ADD CONSTRAINT FK_REKENING_PENERIMA_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);
CREATE INDEX KLAIMNP.IX_REKENING_PENERIMA_1 ON KLAIMNP.REKENING_PENERIMA (ID_KLAIM);

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

CREATE SEQUENCE KLAIMNP.SQ_REKENING_PENERIMA START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_REKENING_PENERIMA TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.REKENING_PENERIMA TO KLAIMNP_APP;

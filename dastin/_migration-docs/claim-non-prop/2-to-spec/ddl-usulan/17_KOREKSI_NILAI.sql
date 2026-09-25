-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Pengganti 29 langkah tambalan per-case. Koreksi bernilai tercatat.

CREATE TABLE KLAIMNP.KOREKSI_NILAI
(
  ID_KOREKSI         NUMBER(19) NOT NULL,
  TABEL_SASARAN      VARCHAR2(30 CHAR) NOT NULL,
  ID_SASARAN         NUMBER(19) NOT NULL,
  KOLOM_SASARAN      VARCHAR2(30 CHAR) NOT NULL,
  NILAI_SEBELUM      NUMBER(38,20),
  NILAI_SESUDAH      NUMBER(38,20),
  ALASAN             VARCHAR2(2000 CHAR) NOT NULL,
  DIBUAT_OLEH        VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA   VARCHAR2(128 CHAR),
  DIBUAT_PADA        TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH        VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA   VARCHAR2(128 CHAR),
  DIUBAH_PADA        TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.KOREKSI_NILAI ADD CONSTRAINT PK_KOREKSI_NILAI PRIMARY KEY (ID_KOREKSI);
-- Sebuah koreksi yang tidak mengubah apa pun bukan koreksi. Keduanya boleh
-- kosong sendiri-sendiri — NILAI_SEBELUM kosong berarti kolomnya memang belum
-- pernah terisi; NILAI_SESUDAH kosong berarti koreksi ini MEMBATALKAN nilai
-- (tiket 22 butir 2) — tetapi kosong keduanya tidak menyatakan apa pun.
ALTER TABLE KLAIMNP.KOREKSI_NILAI ADD CONSTRAINT CK_KOREKSI_NILAI_1 CHECK (
  NILAI_SEBELUM IS NOT NULL OR NILAI_SESUDAH IS NOT NULL
);

-- KOLOM_SASARAN ikut di index karena tiket 22 butir 4 menjanjikan riwayat
-- koreksi ATAS SATU KOLOM TERTENTU terbaca lewat SATU kueri. Tanpa kolomnya
-- di index, janji itu berjalan lewat pemindaian seluruh riwayat baris.
CREATE INDEX KLAIMNP.IX_KOREKSI_NILAI_1 ON KLAIMNP.KOREKSI_NILAI (TABEL_SASARAN, ID_SASARAN, KOLOM_SASARAN, DIBUAT_PADA);

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

CREATE SEQUENCE KLAIMNP.SQ_KOREKSI_NILAI START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_KOREKSI_NILAI TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.KOREKSI_NILAI TO KLAIMNP_APP;

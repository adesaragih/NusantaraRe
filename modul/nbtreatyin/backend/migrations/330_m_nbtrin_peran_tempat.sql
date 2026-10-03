-- 330 - M_NBTRIN_PERAN_TEMPAT: pemetaan TEMPAT -> PERAN -> ARAH untuk tempat yang
-- di sistem lama memeriksa identitas orang (tiket 05; P12, P28).
--
-- TABEL INI DIISI KEMUDIAN oleh IAM bersama work owner - nol baris dari migrasi.
-- Nama orang TIDAK dibawa; tempat tanpa baris = tertunda (AC 81): bagian layarnya
-- tidak ditampilkan dan syaratnya tidak dianggap terpenuhi. ARAH menyatakan
-- 'MUNCUL' (hanya untuk peran ini) atau 'KECUALI' (untuk semua kecuali peran ini)
-- dan TIDAK ditebak (AC 82).
CREATE TABLE {skema}.M_NBTRIN_PERAN_TEMPAT (
  KODE_TEMPAT  VARCHAR2(64) NOT NULL,
  PERAN        VARCHAR2(64) NOT NULL,
  ARAH         VARCHAR2(16) NOT NULL,
  CONSTRAINT PK_NBTRIN_PERAN_TEMPAT PRIMARY KEY (KODE_TEMPAT, PERAN),
  CONSTRAINT CK_NBTRIN_PERAN_TEMPAT_ARAH CHECK (ARAH IN ('MUNCUL', 'KECUALI'))
)
/

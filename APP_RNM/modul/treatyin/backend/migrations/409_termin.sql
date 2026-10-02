-- TERMIN - termin pembayaran premi bernomor urut, per mata uang
--
-- Migrasi tiket 27 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/27-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/23_TERMIN.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- WPC dibawa dengan nama korpusnya apa adanya; `KAMUS-KOLOM.md` tidak
-- memerikan kepanjangannya, dan menebaknya berarti menanam arti yang tidak
-- pernah diverifikasi. DITAGIH: teknik treaty, satu kalimat.
--
-- INV-12: satu baris per (versi, nomor termin, mata uang). Mata uang IKUT di
-- dalam kuncinya - tanpanya, dua termin bernomor sama bermata uang berbeda
-- akan ditolak, dan itu data yang SAH.
CREATE TABLE {skema}.TERMIN (
  ID_TERMIN            NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK     NUMBER(19)          NOT NULL,
  NOMOR_TERMIN         NUMBER(10)          NOT NULL,
  PERSEN_TERMIN        NUMBER(38,8)        NOT NULL,
  NILAI_TERMIN         NUMBER(38,8),
  KODE_MATA_UANG       VARCHAR2(1000 CHAR) NOT NULL,
  TANGGAL_JATUH_TEMPO  DATE                NOT NULL,
  TANGGAL_BAYAR        DATE,
  WPC                  VARCHAR2(1000 CHAR),
  CONSTRAINT PK_TERMIN PRIMARY KEY (ID_TERMIN),
  CONSTRAINT UQ_TERMIN UNIQUE (ID_VERSI_KONTRAK, NOMOR_TERMIN, KODE_MATA_UANG),
  CONSTRAINT FK_TERMIN_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/

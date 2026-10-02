-- PORTOFOLIO - portofolio masuk dan keluar yang menyertai kontrak
--
-- Migrasi tiket 24 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/24-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/20_PORTOFOLIO.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- INV-66: satu baris per (versi, arah, jenis).
--
-- ARAH_PORTOFOLIO dan JENIS_PORTOFOLIO himpunan TERTUTUP, tanpa CHECK -
-- sejalan dengan SIFAT_PROPORSI (lihat 401). ADR-0056 menahan aturan di
-- services; pertanyaan "bolehkah CHECK untuk himpunan tertutup" belum dijawab.
CREATE TABLE {skema}.PORTOFOLIO (
  ID_PORTOFOLIO     NUMBER(19)        NOT NULL,
  ID_VERSI_KONTRAK  NUMBER(19)        NOT NULL,
  ARAH_PORTOFOLIO   VARCHAR2(40 CHAR) NOT NULL,
  JENIS_PORTOFOLIO  VARCHAR2(40 CHAR) NOT NULL,
  KETERANGAN        VARCHAR2(1000 CHAR),
  CONSTRAINT PK_PORTOFOLIO PRIMARY KEY (ID_PORTOFOLIO),
  CONSTRAINT UQ_PORTOFOLIO UNIQUE (ID_VERSI_KONTRAK, ARAH_PORTOFOLIO, JENIS_PORTOFOLIO),
  CONSTRAINT FK_PORTOFOLIO_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/

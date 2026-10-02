-- SKALA_KOASURANSI - skala ko-asuransi sebagai beberapa baris, bukan satu nilai
--
-- Migrasi tiket 28 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/28-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/22_SKALA_KOASURANSI.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- INV-13: satu baris per (versi, persen limit).
CREATE TABLE {skema}.SKALA_KOASURANSI (
  ID_SKALA_KOASURANSI  NUMBER(19)   NOT NULL,
  ID_VERSI_KONTRAK     NUMBER(19)   NOT NULL,
  PERSEN_LIMIT         NUMBER(38,8) NOT NULL,
  PERSEN_BAGIAN        NUMBER(38,8) NOT NULL,
  CONSTRAINT PK_SKALA_KOASURANSI PRIMARY KEY (ID_SKALA_KOASURANSI),
  CONSTRAINT UQ_SKALA_KOASURANSI UNIQUE (ID_VERSI_KONTRAK, PERSEN_LIMIT),
  CONSTRAINT FK_SKALA_KOASURANSI_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/

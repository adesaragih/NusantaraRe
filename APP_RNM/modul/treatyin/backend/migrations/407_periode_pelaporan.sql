-- PERIODE_PELAPORAN - periode pelaporan beserta jatuh temponya
--
-- Migrasi tiket 25 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/25-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/18_PERIODE_PELAPORAN.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- INV-10: satu baris per (versi, periode).
--
-- ⛔ INV-55 SENGAJA BUKAN CONSTRAINT - pernyataan keputusan.
--   APA     : tidak ada yang menolak periode pelaporan di luar periode kontrak.
--   KENAPA  : INV-55 membandingkan tanggal di tabel ini dengan TANGGAL_MULAI
--             dan TANGGAL_BERAKHIR di KONTRAK - tabel LAIN. CHECK Oracle tidak
--             dapat menyatakannya; satu-satunya bentuk basis data yang bisa
--             adalah trigger, dan ADR-0056 (K-4) melarangnya. Bentuknya sama
--             dengan INV-53 (401) dan INV-54 (tiket 03 Adjustment), dan
--             ketiganya ditegakkan di tempat yang sama supaya tidak tersebar.
--   AKIBAT  : sampai jalur simpan berdiri, tanggal di luar periode diterima.
--   DITAGIH : tiket lapisan aplikasi, yang menulis jalur simpannya.
CREATE TABLE {skema}.PERIODE_PELAPORAN (
  ID_PERIODE_PELAPORAN  NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK      NUMBER(19)          NOT NULL,
  PERIODE               VARCHAR2(1000 CHAR) NOT NULL,
  TANGGAL_AWAL          DATE                NOT NULL,
  BATAS_PENYERAHAN      DATE                NOT NULL,
  BATAS_KONFIRMASI      DATE                NOT NULL,
  BATAS_PELUNASAN       DATE                NOT NULL,
  CONSTRAINT PK_PERIODE_PELAPORAN PRIMARY KEY (ID_PERIODE_PELAPORAN),
  CONSTRAINT UQ_PERIODE_PELAPORAN UNIQUE (ID_VERSI_KONTRAK, PERIODE),
  CONSTRAINT FK_PERIODE_PELAPORAN_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/

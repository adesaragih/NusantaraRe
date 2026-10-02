-- JEJAK_PERUBAHAN - siapa mengubah fakta apa dan kapan
--
-- Migrasi tiket 39 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/39-*.md`, `2-to-spec/KAMUS-KOLOM.md`,
-- `ddl-usulan/14_JEJAK_PERUBAHAN.sql`.
--
-- ⛔ TANPA KUNCI ALAMI, dan itu KEPUTUSAN: peristiwa yang sama dapat terjadi
-- DUA KALI pada versi yang sama - ruas yang sama diubah bolak-balik ke nilai
-- yang sama. UNIQUE di sini akan menolak sejarah yang benar.
-- Lihat `Z00_KUNCI_ALAMI.sql`, bab "kunci alami yang sengaja bukan constraint".
--
-- PERAN_PELAKU adalah POTRET, bukan rujukan (`KTV-D`): peran pelaku pada SAAT
-- perubahan, dibekukan. Rujukan ke tabel peran akan berbohong begitu peran
-- orangnya berubah - jejak audit yang berubah surut bukan jejak audit.
-- ⚠️ DARI MANA nilainya diambil masih `F-16`; BENTUKNYA sudah diputuskan.
--
-- NILAI_SEBELUM dan NILAI_SESUDAH VARCHAR2(4000) - teks bebas, mengikuti lebar
-- yang sistem lama pakai untuk teks bebasnya sendiri (KTV-A).
--
-- ⛔ MEKANISME PENGISIANNYA BELUM ADA. Tiket 39 menyebut "entitas beserta
-- kedelapan atributnya, DAN mekanisme pengisiannya"; yang berdiri di sini
-- hanya entitasnya. Pengisiannya menuntut jalur simpan, yang belum
-- berspesifikasi (`L-4`) - dan nol trigger (ADR-0056) berarti ia tidak akan
-- terisi sendiri. DITAGIH: tiket lapisan aplikasi.
--
-- INV-18 - TANPA "ON DELETE": menghapus versi yang berjejak harus GAGAL,
-- bukan menghapus jejaknya.
CREATE TABLE {skema}.JEJAK_PERUBAHAN (
  ID_JEJAK_PERUBAHAN  NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK    NUMBER(19)          NOT NULL,
  WAKTU_PERUBAHAN     DATE                NOT NULL,
  PELAKU              VARCHAR2(1000 CHAR) NOT NULL,
  PERAN_PELAKU        VARCHAR2(1000 CHAR) NOT NULL,
  RUAS_YANG_BERUBAH   VARCHAR2(1000 CHAR) NOT NULL,
  NILAI_SEBELUM       VARCHAR2(4000 CHAR),
  NILAI_SESUDAH       VARCHAR2(4000 CHAR),
  CONSTRAINT PK_JEJAK_PERUBAHAN PRIMARY KEY (ID_JEJAK_PERUBAHAN),
  CONSTRAINT FK_JEJAK_PERUBAHAN_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/
CREATE INDEX {skema}.IX_JEJAK_PERUBAHAN_VERSI ON {skema}.JEJAK_PERUBAHAN (ID_VERSI_KONTRAK)
/

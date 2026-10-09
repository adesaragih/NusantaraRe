-- CATATAN_PERSETUJUAN - jejak tiap perpindahan di jalur persetujuan
--
-- Migrasi tiket 54 Treaty In, dan ia melepas tiket 49, 55, 56. Satu berkas =
-- satu langkah migrasi. Pernyataan dipisahkan oleh baris yang hanya berisi
-- tanda garis miring.
--
-- Asal: `2-to-spec/ddl-usulan/13_CATATAN_PERSETUJUAN.sql` dan
-- `SPEC-MODEL-DATA.md` §10.20 - enam atribut, kelas lama `Data-SuggestList`.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS: TOLAK, dan sebabnya dikutip utuh
-- ---------------------------------------------------------------------
--   `ERD.md` §2.3 - dokumen MENGIKAT - menandainya **[hapus: tolak]**, dan
--   menuliskan sebabnya dalam satu kalimat yang tidak perlu ditambah:
--
--       "jejak yang dapat dihapus bersama bendanya bukan jejak."
--
--   Diwujudkan dengan TIDAK menulis klausa `ON DELETE`: bawaan Oracle
--   menolak, dan di sini bawaan itu DIPILIH, bukan dibiarkan. Barisnya ada
--   di `TestPerilakuHapusSesuaiERD` sebagai `tolak`.
--
--   ⚠️ Ia karena itu BERBEDA dari sepuluh anak `VERSI_KONTRAK` lain di §2.3
--   yang seluruhnya `ikut hapus`. Perbedaan itu disengaja; menyeragamkannya
--   menghapus jejak persetujuan bersama versi yang disetujuinya.
--
-- ---------------------------------------------------------------------
-- NOL KUNCI ALAMI, dan itu KEPUTUSAN - bukan tagihan
-- ---------------------------------------------------------------------
--   §10.20 menyatakannya: dua keputusan pada versi yang sama, oleh orang yang
--   sama, pada hari yang sama adalah keadaan yang SAH. Yang membedakan
--   barisnya urutan waktu, bukan sebuah nilai.
--
--   Ia karena itu TIDAK ikut daftar tagihan `SPEC-MODEL-DATA.md` §13: enam
--   entitas di sana menunggu NOMOR invarian untuk kunci yang sudah jelas
--   bentuknya; yang ini tidak punya kunci sama sekali, dan tidak boleh punya.
--
-- ---------------------------------------------------------------------
-- `NAMA_PEMUTUS` adalah TEKS, dan itu disengaja
-- ---------------------------------------------------------------------
--   Bukan kunci asing ke tabel pelaku. §12.5 dan `ADR-0045`: ia fakta
--   historis - siapa memutuskan apa dan kapan - dan fakta historis yang
--   dinormalkan ulang berubah ketika orangnya berganti nama atau keluar.
--   Sekeluarga dengan `PELAKU` di `JEJAK_PERUBAHAN` (migrasi 414).
--
-- ⚠️ `DISETUJUI` himpunan TERTUTUP dan TIDAK ber-`CHECK`: `ADR-0056`
-- menempatkan aturan bisnis di aplikasi, bukan di basis data. Dua nilainya
-- ditegakkan services bersama tiket 49-53.
--
-- ⛔ Baris PERISTIWA (`IsApproved` kosong di sistem lama) TIDAK masuk ke
-- sini. §10.20a mencabut usul nilai enumerasi ketiga `PERISTIWA`: ia akan
-- membuat separuh kolom kosong pada separuh baris dan mencemari `INV-28`,
-- yang menghitung satu baris per perpindahan. Rumahnya `PERISTIWA_KONTRAK`,
-- yang sampai hari ini masih yatim - nol tiket menghasilkannya.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE SEQUENCE {skema}.SEQ_TRIN_CATATAN_PERSETUJUAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.CATATAN_PERSETUJUAN (
  ID_CATATAN_PERSETUJUAN  NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK        NUMBER(19)          NOT NULL,
  WAKTU_KEPUTUSAN         DATE                NOT NULL,
  NAMA_PEMUTUS            VARCHAR2(1000 CHAR) NOT NULL,
  DISETUJUI               VARCHAR2(40 CHAR)   NOT NULL,
  ALASAN                  VARCHAR2(1000 CHAR),
  CONSTRAINT PK_CATATAN_PERSETUJUAN PRIMARY KEY (ID_CATATAN_PERSETUJUAN),
  CONSTRAINT FK_CATATAN_PERSETUJUAN_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/
CREATE INDEX {skema}.IX_CATATAN_PERSETUJUAN_VERSI ON {skema}.CATATAN_PERSETUJUAN (ID_VERSI_KONTRAK)
/

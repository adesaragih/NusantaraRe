-- DOKUMEN_KONTRAK - dokumen dilampirkan, dirujuk bukan disalin
--
-- Migrasi tiket 30 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/30-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/12_DOKUMEN_KONTRAK.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS DITETAPKAN SADAR. ⚠️ RALAT 2 Oktober 2026.
-- ---------------------------------------------------------------------
--   Berkas ini pernah berbunyi: "INV-18 - TANPA ON DELETE: bawaan Oracle
--   MENOLAK, dan menolak yang dikehendaki." ITU MEMBACA INV-18 TERBALIK.
--
--   INV-18 berbunyi "perilaku hapus setiap kunci asing DITETAPKAN SADAR,
--   TIDAK DIBIARKAN BAWAAN". Membiarkan bawaan lalu menyebut nama invarian
--   yang melarangnya bukan pemenuhan - bahkan ketika bawaannya kebetulan
--   cocok. Yang dituntut keputusan yang DINYATAKAN, bukan yang kebetulan benar.
--
--   Keputusan sadarnya sudah ada dan MENGIKAT: `4-erd-dan-tabel-datar/ERD.md`
--   §2, yang menyatakan `ikut hapus`, `tolak`, atau `putus` per relasi.
--   Dokumen itu tidak pernah dibuka sampai 2 Oktober 2026 - tiga puluh tiga
--   kunci asing ditulis tanpanya.
--
--   Perilaku hapus tiap kunci asing di berkas ini kini disebut di barisnya
--   sendiri. Daftar lengkap seluruh modul: bab "Kaskade"
--   di `MODUL.md`, yang juga menyalakan `TestKaskadeHanyaPadaRelasiTerdaftar`.
--
-- INV-67: satu lampiran per (versi, dokumen).
--
-- ⛔ ID_DOKUMEN TANPA KUNCI ASING - pernyataan keputusan.
--   APA     : ia menunjuk dokumen di penyimpanan berkas, BUKAN baris di skema
--             ini. Tidak ada tabel DOKUMEN di modul ini untuk ditunjuk.
--   KENAPA   : pokok tiket 30 justru itu - dokumennya DIRUJUK, bukan DISALIN.
--   AKIBAT  : basis data tidak dapat menjamin dokumennya ada. Yang menjamin
--             adalah lapisan yang memegang penyimpanan berkas.
CREATE TABLE {skema}.DOKUMEN_KONTRAK (
  ID_DOKUMEN_KONTRAK  NUMBER(19)        NOT NULL,
  ID_VERSI_KONTRAK    NUMBER(19)        NOT NULL,
  ID_DOKUMEN          NUMBER(19)        NOT NULL,
  JENIS_DOKUMEN       VARCHAR2(40 CHAR),
  TANGGAL_LAMPIR      DATE              NOT NULL,
  CONSTRAINT PK_DOKUMEN_KONTRAK PRIMARY KEY (ID_DOKUMEN_KONTRAK),
  CONSTRAINT UQ_DOKUMEN_KONTRAK UNIQUE (ID_VERSI_KONTRAK, ID_DOKUMEN),
  CONSTRAINT FK_DOKUMEN_KONTRAK_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/

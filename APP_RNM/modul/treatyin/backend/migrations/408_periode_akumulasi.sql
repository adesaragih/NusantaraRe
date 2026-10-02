-- PERIODE_AKUMULASI - periode akumulasi, ditolak bila di luar periode kontrak
--
-- Migrasi tiket 26 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/26-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/17_PERIODE_AKUMULASI.sql`.
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
-- INV-11: satu baris per (versi, periode).
--
-- ⛔ INV-56 SENGAJA BUKAN CONSTRAINT - sebab yang sama dengan INV-55 (407):
-- ia membandingkan tanggal di tabel ini dengan periode di KONTRAK, tabel lain.
-- Ditegakkan di services bersama INV-53, INV-54, dan INV-55.
-- DITAGIH: tiket lapisan aplikasi.
CREATE TABLE {skema}.PERIODE_AKUMULASI (
  ID_PERIODE_AKUMULASI   NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK       NUMBER(19)          NOT NULL,
  PERIODE                VARCHAR2(1000 CHAR) NOT NULL,
  TANGGAL_LAPOR          DATE                NOT NULL,
  HARI_BATAS_PENYERAHAN  NUMBER(10),
  BATAS_PENYERAHAN       DATE,
  CONSTRAINT PK_PERIODE_AKUMULASI PRIMARY KEY (ID_PERIODE_AKUMULASI),
  CONSTRAINT UQ_PERIODE_AKUMULASI UNIQUE (ID_VERSI_KONTRAK, PERIODE),
  CONSTRAINT FK_PERIODE_AKUMULASI_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)
)
/

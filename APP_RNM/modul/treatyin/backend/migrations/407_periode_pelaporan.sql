-- PERIODE_PELAPORAN - periode pelaporan beserta jatuh temponya
--
-- Migrasi tiket 25 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/25-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/18_PERIODE_PELAPORAN.sql`.
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

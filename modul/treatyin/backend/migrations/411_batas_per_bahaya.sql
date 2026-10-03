-- BATAS_PER_BAHAYA - batas tanggungan untuk bahaya apa pun di daftar
--
-- Migrasi tiket 29 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/29-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/10_BATAS_PER_BAHAYA.sql`.
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
-- INV-14: satu batas per (versi, bahaya).
--
-- ⭐ INI POKOK ADR-0038. Sistem lama menaruh empat bahaya sebagai EMPAT PASANG
-- KOLOM di kepala kontrak (Earthquake, FloodJab, FloodNation, RSMDLimit),
-- sehingga bahaya kesembilan menuntut kolom kesembilan. Di sini bahaya adalah
-- BARIS di tabel acuan (400) dan batasnya BARIS di tabel ini - menambah bahaya
-- tidak menyentuh skema sama sekali.
--
-- CONSTRAINT YANG SENGAJA TIDAK DIPASANG: NILAI_BATAS tidak dipasangkan dengan
-- PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40). DITAGIH: T-2 dari teknik treaty.
CREATE TABLE {skema}.BATAS_PER_BAHAYA (
  ID_BATAS_PER_BAHAYA    NUMBER(19)   NOT NULL,
  ID_VERSI_KONTRAK       NUMBER(19)   NOT NULL,
  ID_BAHAYA              NUMBER(19)   NOT NULL,
  NILAI_BATAS            NUMBER(38,8) NOT NULL,
  KODE_MATA_UANG         VARCHAR2(1000 CHAR),
  PERSEN_BAGIAN_DIPAKAI  NUMBER(38,8),
  CONSTRAINT PK_BATAS_PER_BAHAYA PRIMARY KEY (ID_BATAS_PER_BAHAYA),
  CONSTRAINT UQ_BATAS_PER_BAHAYA UNIQUE (ID_VERSI_KONTRAK, ID_BAHAYA),
  CONSTRAINT FK_BATAS_PER_BAHAYA_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK),
  CONSTRAINT FK_BATAS_PER_BAHAYA_2 FOREIGN KEY (ID_BAHAYA)
    REFERENCES {skema}.BAHAYA (ID_BAHAYA)
)
/
CREATE INDEX {skema}.IX_BATAS_PER_BAHAYA_BHY ON {skema}.BATAS_PER_BAHAYA (ID_BAHAYA)
/

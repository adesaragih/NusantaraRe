-- Tabel acuan Treaty In - enam himpunan yang dapat bertambah
--
-- Migrasi tiket 15 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `D:\XML_NURE\_migration-docs\treaty-in\5-tiket\issues\15-*.md`,
-- `2-to-spec/KAMUS-KOLOM.md`, `2-to-spec/ddl-usulan/{01,03,04,05,06,08}.sql`.
--
-- ADR-0038: yang bertambah tanpa mengubah arti disimpan sebagai DATA, bukan
-- sebagai nama kolom dan bukan sebagai CHECK. Sistem lama menaruh empat bahaya
-- sebagai empat pasang kolom di kepala kontrak (Earthquake, FloodJab,
-- FloodNation, RSMDLimit) - bentuk yang menuntut kolom kesembilan begitu ada
-- bahaya baru.
--
-- ⛔ NOL CHECK berisi daftar nilai untuk keenam himpunan ini (INV-62). Itu
-- PERNYATAAN KEPUTUSAN, bukan kelalaian: sebuah CHECK di sini mengembalikan
-- cacat yang ADR-0038 buang.
--
-- INV-01 setiap tabel berkunci utama · INV-02 pengenal dari SEQUENCE (berkas
-- 402) · INV-68 KODE unik di keenamnya.
--
-- ⚠️ EMPAT PENYELARASAN dengan repo, tercatat di
-- `docs/KEPUTUSAN-PENYELARASAN-REPO.md` - jangan diubah tanpa membacanya:
--   1. skema `{skema}`, bukan `TREATY_MASUK` (ADR-0028 dibalik)
--   2. presisi NUMBER mengikuti `presisiSah` penjaga
--   3. nama tabel memakai nama spec apa adanya
--   4. sequence ditulis sendiri - `ddl-usulan/` tidak memuat satu pun
CREATE TABLE {skema}.MATA_UANG (
  ID_MATA_UANG  NUMBER(19)          NOT NULL,
  KODE          VARCHAR2(1000 CHAR) NOT NULL,
  NAMA          VARCHAR2(1000 CHAR) NOT NULL,
  AKTIF         VARCHAR2(40 CHAR)   NOT NULL,
  CONSTRAINT PK_MATA_UANG PRIMARY KEY (ID_MATA_UANG),
  CONSTRAINT UQ_MATA_UANG UNIQUE (KODE)
)
/
CREATE TABLE {skema}.JENIS_POTONGAN (
  ID_JENIS_POTONGAN  NUMBER(19)          NOT NULL,
  KODE               VARCHAR2(1000 CHAR) NOT NULL,
  NAMA               VARCHAR2(1000 CHAR) NOT NULL,
  AKTIF              VARCHAR2(40 CHAR)   NOT NULL,
  CONSTRAINT PK_JENIS_POTONGAN PRIMARY KEY (ID_JENIS_POTONGAN),
  CONSTRAINT UQ_JENIS_POTONGAN UNIQUE (KODE)
)
/
CREATE TABLE {skema}.KELAS_BISNIS (
  ID_KELAS_BISNIS  NUMBER(19)          NOT NULL,
  KODE             VARCHAR2(1000 CHAR) NOT NULL,
  NAMA             VARCHAR2(1000 CHAR) NOT NULL,
  AKTIF            VARCHAR2(40 CHAR)   NOT NULL,
  CONSTRAINT PK_KELAS_BISNIS PRIMARY KEY (ID_KELAS_BISNIS),
  CONSTRAINT UQ_KELAS_BISNIS UNIQUE (KODE)
)
/
CREATE TABLE {skema}.KELOMPOK_TREATY (
  ID_KELOMPOK_TREATY  NUMBER(19)          NOT NULL,
  KODE                VARCHAR2(1000 CHAR) NOT NULL,
  NAMA                VARCHAR2(1000 CHAR) NOT NULL,
  AKTIF               VARCHAR2(40 CHAR)   NOT NULL,
  CONSTRAINT PK_KELOMPOK_TREATY PRIMARY KEY (ID_KELOMPOK_TREATY),
  CONSTRAINT UQ_KELOMPOK_TREATY UNIQUE (KODE)
)
/
CREATE TABLE {skema}.BAHAYA (
  ID_BAHAYA  NUMBER(19)          NOT NULL,
  KODE       VARCHAR2(1000 CHAR) NOT NULL,
  NAMA       VARCHAR2(1000 CHAR) NOT NULL,
  AKTIF      VARCHAR2(40 CHAR)   NOT NULL,
  CONSTRAINT PK_BAHAYA PRIMARY KEY (ID_BAHAYA),
  CONSTRAINT UQ_BAHAYA UNIQUE (KODE)
)
/
-- ID_INDUK menunjuk baris lain di TABEL YANG SAMA: JENIS_REASURANSI bersusun
-- (tiket 15; artinya dipakai pertama kali oleh tiket 38). Kolomnya ikut di sini,
-- susunannya belum dipakai.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS DITETAPKAN SADAR. ⚠️ RALAT 2 Oktober 2026.
-- ---------------------------------------------------------------------
--   Berkas ini pernah menyebut INV-18 sambil MEMBIARKAN bawaan Oracle. Itu
--   membaca INV-18 terbalik: ia menuntut perilaku hapus "DITETAPKAN SADAR,
--   TIDAK DIBIARKAN BAWAAN", dan bawaan yang kebetulan cocok bukan keputusan.
--   Sumber keputusannya `4-erd-dan-tabel-datar/ERD.md` §2 - dokumen MENGIKAT
--   yang tidak pernah dibuka sampai hari ini.
--
--   ⛔ RELASI INI TIDAK ADA DI ERD.md §2. `JENIS_REASURANSI.ID_INDUK` lahir
--   dari tiket 15 ("JENIS_REASURANSI bersusun"), sesudah §2 ditulis, dan §2.3c
--   yang menambahkan delapan relasi susulan pun tidak memuatnya.
--   DIPUTUSKAN DI SINI: **tolak** - TANPA klausa ON DELETE. Alasannya diambil
--   dari aturan kelompok yang §2.7 nyatakan untuk seluruh tabel acuan: "baris
--   acuan yang sudah dipakai tidak dapat hilang". Menghapus jenis reasuransi
--   induk yang masih punya turunan akan gagal ORA-02292.
--   DITAGIH: pemilik `ERD.md` - satu baris di §2.7, supaya keputusan ini
--   berpindah dari sini ke dokumen yang mengikat.
CREATE TABLE {skema}.JENIS_REASURANSI (
  ID_JENIS_REASURANSI  NUMBER(19)          NOT NULL,
  KODE                 VARCHAR2(1000 CHAR) NOT NULL,
  NAMA                 VARCHAR2(1000 CHAR) NOT NULL,
  AKTIF                VARCHAR2(40 CHAR)   NOT NULL,
  ID_INDUK             NUMBER(19),
  CONSTRAINT PK_JENIS_REASURANSI PRIMARY KEY (ID_JENIS_REASURANSI),
  CONSTRAINT UQ_JENIS_REASURANSI UNIQUE (KODE),
  CONSTRAINT FK_JENIS_REASURANSI_INDUK FOREIGN KEY (ID_INDUK)
    REFERENCES {skema}.JENIS_REASURANSI (ID_JENIS_REASURANSI)
)
/

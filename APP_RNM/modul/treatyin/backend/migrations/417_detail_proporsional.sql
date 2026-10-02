-- DETAIL_PROPORSIONAL - ketentuan proporsional per kelompok treaty di dalam layer
--
-- Migrasi tiket 34 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/34-*.md`, `2-to-spec/KAMUS-KOLOM.md` §10.4,
-- `ddl-usulan/{25_DETAIL_PROPORSIONAL,29_NILAI_CADANGAN_PREMI}.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- ⭐ INV-06 berlingkup DI DALAM LAYER, bukan di dalam versi: "ID_KELOMPOK_TREATY
-- unik di dalam satu LAYER". Menulisnya "di dalam versi" akan MENOLAK dua layer
-- yang sama-sama memuat kelompok treaty yang sama - dan itu susunan yang lazim.
--
-- ⭐ Bagian NuRe pada cabang proporsional TIDAK punya baris sendiri: ia atribut
-- `RNMShare` PADA baris ini, terisi bila `RNMShareAcrossTheBoard` mati (§12.3).
-- Cabang non-proporsional memakai tabel `BAGIAN` tersendiri (416).
--
-- ⛔ `KAPASITAS_SURPLUS` TIDAK ADA sebagai kolom tersimpan - ia KELIPATAN
-- retensi x jumlah lines, dan INV-48 menamainya pengecualian bernama terhadap
-- INV-47. Turunan tidak disimpan (ADR-0037). Kolom `PERSEN_KAPASITAS_SURPLUS`
-- di bawah adalah PERSEN, bukan paket uangnya - keduanya berbeda.
--
-- ---------------------------------------------------------------------
-- ⛔ INV-32 BELUM DITEGAKKAN - sebab yang SAMA PERSIS dengan INV-33 (416).
-- ---------------------------------------------------------------------
--   APA     : tidak ada yang menolak DETAIL_PROPORSIONAL di bawah kontrak
--             ber-SIFAT_PROPORSI `NON_PROPORSIONAL`.
--   BENTUK  : tiket 34 menggolongkannya INDEKS UNIK, yang di Oracle berarti
--             materialized view ber-REFRESH ON COMMIT berindeks unik - sebab
--             SIFAT_PROPORSI ada dua tabel di atas.
--   KENAPA BELUM: `F-13-MV-TANPA-PEMANTAU-KEBASIAN.md` - penegakan lewat MV
--             belum punya pemantau kebasian, dan MV basi berhenti menegakkan
--             TANPA SATU GALAT PUN. Uraian penuh di berkas 416.
--   ⚠️ Tiket 34 menambahkan satu tuntutan yang berkas 416 tidak punya:
--             "INV-32 terpasang sebagai indeks unik, dan TIDAK MENELAN INV-33".
--             Keduanya harus menolak hal yang BERBEDA - satu melarang cabang
--             proporsional, satu melarang non-proporsional - dan MV yang
--             dirancang ceroboh dapat membuat salah satunya menelan yang lain.
--   DITAGIH : pemilik F-13, bersama INV-33.
--
-- ⛔ INV-30 (`JENIS_TREATY` dua nilai), INV-36, dan INV-41 juga belum
-- ditegakkan. INV-30 adalah daftar nilai TERTUTUP - pertanyaan "bolehkah CHECK
-- untuk himpunan tertutup" belum dijawab (lihat INV-29 di berkas 401).
--
-- ⛔ Aturan TEPAT SATU antara PERSEN_QUOTA_SHARE dan JUMLAH_LINES_SURPLUS
-- adalah irisan 35, BUKAN di sini. Kolomnya berdiri di sini; aturannya tidak.
-- Tiket 34 menyatakannya di bab "Tidak termasuk" dengan kalimatnya sendiri.
CREATE TABLE {skema}.DETAIL_PROPORSIONAL (
  ID_DETAIL_PROPORSIONAL    NUMBER(19)        NOT NULL,
  ID_LAYER                  NUMBER(19)        NOT NULL,
  ID_KELOMPOK_TREATY        NUMBER(19)        NOT NULL,
  JENIS_TREATY              VARCHAR2(40 CHAR) NOT NULL,
  PERSEN_QUOTA_SHARE        NUMBER(38,8),
  JUMLAH_LINES_SURPLUS      NUMBER(10),
  PERSEN_KOMISI_KOTOR       NUMBER(38,8),
  PERSEN_KOMISI_BERSIH      NUMBER(38,8),
  PERSEN_CADANGAN_PREMI     NUMBER(38,8),
  CADANGAN_PREMI            NUMBER(38,8),
  PERSEN_KAPASITAS_SURPLUS  NUMBER(38,8),
  ID_SUSUNAN_RETRO          NUMBER(19),
  CONSTRAINT PK_DETAIL_PROPORSIONAL PRIMARY KEY (ID_DETAIL_PROPORSIONAL),
  CONSTRAINT UQ_DETAIL_PROPORSIONAL UNIQUE (ID_LAYER, ID_KELOMPOK_TREATY),
  CONSTRAINT FK_DETAIL_PROPORSIONAL_1 FOREIGN KEY (ID_LAYER)
    REFERENCES {skema}.LAYER (ID_LAYER),
  CONSTRAINT FK_DETAIL_PROPORSIONAL_2 FOREIGN KEY (ID_KELOMPOK_TREATY)
    REFERENCES {skema}.KELOMPOK_TREATY (ID_KELOMPOK_TREATY)
)
/
CREATE TABLE {skema}.NILAI_CADANGAN_PREMI (
  ID_NILAI_CADANGAN_PREMI  NUMBER(19)          NOT NULL,
  ID_DETAIL_PROPORSIONAL   NUMBER(19)          NOT NULL,
  KODE_MATA_UANG           VARCHAR2(1000 CHAR) NOT NULL,
  NILAI                    NUMBER(38,8)        NOT NULL,
  CONSTRAINT PK_NILAI_CADANGAN_PREMI PRIMARY KEY (ID_NILAI_CADANGAN_PREMI),
  CONSTRAINT FK_NILAI_CADANGAN_PREMI_1 FOREIGN KEY (ID_DETAIL_PROPORSIONAL)
    REFERENCES {skema}.DETAIL_PROPORSIONAL (ID_DETAIL_PROPORSIONAL)
)
/
CREATE INDEX {skema}.IX_DETAIL_PROP_KLP ON {skema}.DETAIL_PROPORSIONAL (ID_KELOMPOK_TREATY)
/
CREATE INDEX {skema}.IX_NILAI_CAD_PREMI_DP ON {skema}.NILAI_CADANGAN_PREMI (ID_DETAIL_PROPORSIONAL)
/

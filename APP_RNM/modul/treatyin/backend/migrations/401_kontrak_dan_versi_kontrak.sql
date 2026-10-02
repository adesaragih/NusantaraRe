-- KONTRAK dan VERSI_KONTRAK - identitas kontrak yang terpecah dua
--
-- Migrasi tiket 14 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `D:\XML_NURE\_migration-docs\treaty-in\5-tiket\issues\14-*.md`,
-- `2-to-spec/KAMUS-KOLOM.md` (§10.1 8 kolom, §10.2 48 kolom),
-- `2-to-spec/ddl-usulan/{07_KONTRAK,09_VERSI_KONTRAK,Z00_KUNCI_ALAMI}.sql`.
--
-- ADR-0040 memecah identitas menjadi dua: KONTRAK memegang LAPISAN BEKU -
-- cedant, asal bisnis, sifat proporsi, periode - dan VERSI_KONTRAK memegang
-- segala yang dapat berbeda antar versi. Sistem lama menyimpan kontrak dan
-- addendum di SATU baris `M_TREATY_IN` dengan seluruh halaman clipboard sebagai
-- satu kolom `JSONDATA`; identitas yang terpecah dua TIDAK PERNAH ADA di sana.
--
-- Invarian yang berkas ini tegakkan:
--   INV-01  setiap tabel berkunci utama
--   INV-02  pengenal dari SEQUENCE (berkas 402), tidak pernah dari cap waktu
--           maupun teks
--   INV-04  NOMOR_URUT_VERSI unik di dalam satu KONTRAK
--   INV-18  perilaku hapus tiap kunci asing ditetapkan sadar
--
-- ⚠️ EMPAT PENYELARASAN dengan repo, tercatat di
-- `docs/KEPUTUSAN-PENYELARASAN-REPO.md` - jangan diubah tanpa membacanya:
--   1. skema `{skema}`, bukan `TREATY_MASUK` (ADR-0028 dibalik)
--   2. `NUMBER(38,20)` -> `NUMBER(38,8)` dan `NUMBER(9)` -> `NUMBER(10)`
--   3. nama tabel memakai nama spec apa adanya
--   4. sequence ditulis sendiri - `ddl-usulan/` tidak memuat satu pun
--
-- ---------------------------------------------------------------------
-- INV-53 SENGAJA BUKAN CHECK DI SINI - ini PERNYATAAN KEPUTUSAN
-- ---------------------------------------------------------------------
--   INV-53 menuntut TANGGAL_MULAI <= TANGGAL_BERAKHIR, keduanya inklusif
--   (ADR-0022). Ia ditegakkan di lapisan services, bukan sebagai CHECK, sebab
--   ADR-0056 (K-4) melarang aturan bisnis turun ke basis data dan seluruh
--   `ddl-usulan/` berdiri tanpa satu pun trigger maupun procedure pembawa
--   aturan.
--   ⛔ BELUM DITEGAKKAN DI MANA PUN. Lapisan services modul ini hari ini
--   hanya membaca tabel acuan; tidak ada jalur simpan, sehingga tidak ada
--   tempat INV-53 berdiri dan tidak ada ujinya. Kalimat ini menyatakan DI MANA
--   ia akan berdiri, bukan bahwa ia sudah berdiri.
--   DITAGIH: tiket lapisan aplikasi, yang menulis jalur simpannya.
-- ---------------------------------------------------------------------
-- INV-29 - `SIFAT_PROPORSI` DUA NILAI, TANPA CHECK. Pernyataan keputusan.
-- ---------------------------------------------------------------------
--   APA        : kolomnya NOT NULL dan menerima teks apa pun sepanjang 40
--                aksara. Dua nilai sahnya - PROPORSIONAL dan NON_PROPORSIONAL -
--                tidak dinyatakan di basis data.
--   KENAPA     : `ddl-usulan/07_KONTRAK.sql` juga tidak menyatakannya, dan
--                ADR-0056 (K-4) menahan aturan di lapisan services. Berbeda
--                dengan keenam himpunan acuan, dua nilai ini TERTUTUP - ADR-0038
--                tidak melarang CHECK di sini. Yang melarang hanya ADR-0056.
--   AKIBAT     : nilai ketiga dapat masuk sampai jalur simpan berdiri. Tiket 63
--                ("cara pembukuan XOL, hanya pada non-proporsional") membaca
--                kolom ini, jadi nilai ketiga akan menyesatkan pembacanya.
--   DITAGIH    : tiket lapisan aplikasi. Bila pemilik proses memutuskan CHECK
--                boleh untuk himpunan TERTUTUP, ia dipasang di sini - dan
--                INV-29 adalah calon pertamanya.

CREATE TABLE {skema}.KONTRAK (
  ID_KONTRAK               NUMBER(19)          NOT NULL,
  NOMOR_KONTRAK_WARISAN    VARCHAR2(1000 CHAR),
  ID_KONTRAK_DISALIN_DARI  NUMBER(19),
  ID_CEDANT                NUMBER(19)          NOT NULL,
  ID_ASAL_BISNIS           NUMBER(19)          NOT NULL,
  SIFAT_PROPORSI           VARCHAR2(40 CHAR)   NOT NULL,
  TANGGAL_MULAI            DATE                NOT NULL,
  TANGGAL_BERAKHIR         DATE                NOT NULL,
  CONSTRAINT PK_KONTRAK PRIMARY KEY (ID_KONTRAK)
)
/
-- ---------------------------------------------------------------------
-- KEADAAN_SIKLUS_HIDUP TANPA DAFTAR NILAI - ini PERNYATAAN KEPUTUSAN,
-- bukan kelalaian dan bukan TODO.
-- ---------------------------------------------------------------------
--   APA        : kolomnya ADA dan NOT NULL. Tidak ada CHECK, tidak ada tabel
--                acuan, dan tidak ada mesin perpindahan yang membatasi isinya.
--   KENAPA     : tiket 14 menyebutnya di bab "Tidak termasuk" dengan nama -
--                INV-20, INV-22, INV-23, INV-24, INV-25 TEGAS bukan bagian
--                irisan ini. Tabelnya tidak dapat berdiri tanpa kolomnya, jadi
--                kolomnya berdiri di sini; ARTINYA dibuat tiket 45.
--                ADR-0055 menyatakan delapan nilai sah (DIBATALKAN ditambahkan
--                24 Sep 2026) dan menuntut tidak ada satu pun jalan menyetel
--                keadaan selain melalui perpindahan di daftar.
--   AKIBAT     : sampai tiket 45 mendarat, kolom ini menerima teks apa pun
--                sepanjang 40 aksara. Pembacanya harus tahu itu.
--   DITAGIH    : tiket 45 - "Daftar keadaan dan perpindahan sah berdiri".
--
-- ---------------------------------------------------------------------
-- FK ID_DOKUMEN_ADDENDUM SENGAJA BELUM DIPASANG
-- ---------------------------------------------------------------------
--   APA        : kolomnya ada, kunci asingnya tidak. `ddl-usulan/09` menulis
--                FK_VERSI_KONTRAK_2 ke DOKUMEN_ADDENDUM.
--   KENAPA     : DOKUMEN_ADDENDUM dibuat tiket 04, yang ada di PAPAN
--                ADJUSTMENT - dan tiket 04 diblokir tiket 14 ini. Memasang FK
--                ke tabel yang belum ada membuat migrasi ini gagal.
--   AKIBAT     : sampai tiket 04 mendarat, kolom ini menerima pengenal yang
--                tidak menunjuk apa pun.
--   DITAGIH    : tiket 04 papan Adjustment, yang membuat tabelnya DAN memasang
--                kunci asingnya.
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
--   Perilaku hapus ketiga kunci asing berkas ini, menurut ERD.md:
--
--   FK_VERSI_KONTRAK_1        KONTRAK 1--< VERSI_KONTRAK      §2.1  TOLAK
--     "kontrak yang punya versi tidak boleh hilang, karena versinya memuat
--      angka yang pernah dibukukan."
--   FK_VERSI_KONTRAK_MATA_UANG  VERSI >o--1 MATA_UANG         §2.7  TOLAK
--     "baris acuan yang sudah dipakai tidak dapat hilang."
--   FK_KONTRAK_DISALIN_DARI   KONTRAK 1--o< KONTRAK           §2.1  PUTUS
--     satu-satunya `putus` di modul ini - lihat di dalam CREATE TABLE.
--
--   TOLAK diwujudkan dengan TIDAK menulis klausa ON DELETE - bentuk yang sama
--   dengan bawaan, tetapi kini DIPILIH dan sumbernya disebut.
CREATE TABLE {skema}.VERSI_KONTRAK (
  ID_VERSI_KONTRAK              NUMBER(19)          NOT NULL,
  ID_KONTRAK                    NUMBER(19)          NOT NULL,
  NOMOR_URUT_VERSI              NUMBER(10)          NOT NULL,
  KEADAAN_SIKLUS_HIDUP          VARCHAR2(40 CHAR)   NOT NULL,
  KEADAAN_WARISAN_ASLI          VARCHAR2(1000 CHAR),
  JENIS_ADDENDUM                VARCHAR2(40 CHAR),
  TANGGAL_BERLAKU_ADDENDUM      DATE,
  NAMA_KONTRAK                  VARCHAR2(1000 CHAR) NOT NULL,
  LINGKUP_WILAYAH               VARCHAR2(1000 CHAR),
  KELAS_BISNIS_KONTRAK          NUMBER(19),
  KODE_MATA_UANG_KONTRAK        NUMBER(19)          NOT NULL,
  ID_REASURADUR_PEMIMPIN        NUMBER(19),
  ID_KETUA_TREATY               NUMBER(19),
  PERSEN_BAGIAN_NURE            NUMBER(38,8)        NOT NULL,
  BAGIAN_NURE_SERAGAM           VARCHAR2(40 CHAR)   NOT NULL,
  PERSEN_BAGIAN_NURE_DIPOTONG   NUMBER(38,8),
  PERSEN_BROKERAGE              NUMBER(38,8),
  PERSEN_BAGIAN_FAKULTATIF      NUMBER(38,8),
  PERSEN_BROKERAGE_FAKULTATIF   NUMBER(38,8),
  PENGECUALIAN                  VARCHAR2(4000 CHAR),
  KETENTUAN_KHUSUS              VARCHAR2(4000 CHAR),
  KETERANGAN                    VARCHAR2(4000 CHAR),
  CATATAN                       VARCHAR2(4000 CHAR),
  CATATAN_BORDEREAUX            VARCHAR2(4000 CHAR),
  MEMAKAI_BORDEREAUX            VARCHAR2(40 CHAR)   NOT NULL,
  CARA_PEMBUKUAN                VARCHAR2(40 CHAR)   NOT NULL,
  CARA_PEMBUKUAN_XOL            VARCHAR2(40 CHAR),
  PERIODE_PELAPORAN_KONTRAK     VARCHAR2(40 CHAR),
  SELANG_PELAPORAN              NUMBER(10),
  TANGGAL_MULAI_PELAPORAN       DATE,
  TANGGAL_AKHIR_PELAPORAN       DATE,
  HARI_BATAS_PENYERAHAN         NUMBER(10),
  HARI_BATAS_KONFIRMASI         NUMBER(10),
  HARI_BATAS_PELUNASAN          NUMBER(10),
  HARI_PENGINGAT                NUMBER(10),
  PERIODE_AKUMULASI_KONTRAK     VARCHAR2(40 CHAR),
  JUMLAH_TERMIN                 NUMBER(10),
  MEMAKAI_PRORATA               VARCHAR2(40 CHAR)   NOT NULL,
  BATAS_MAKSIMUM_KELOMPOK       NUMBER(38,8),
  MATA_UANG_BATAS_KELOMPOK      VARCHAR2(1000 CHAR),
  BATAS_MAKSIMUM_NON_KELOMPOK   NUMBER(38,8),
  MATA_UANG_BATAS_NON_KELOMPOK  VARCHAR2(1000 CHAR),
  BATAS_PILIHAN                 NUMBER(38,8),
  MATA_UANG_BATAS_PILIHAN       VARCHAR2(1000 CHAR),
  RETRO_BERGANDA                VARCHAR2(40 CHAR)   NOT NULL,
  ID_DAFTAR_RETRO               NUMBER(19),
  SIFAT_MATERIAL_ADDENDUM       VARCHAR2(40 CHAR),
  ID_DOKUMEN_ADDENDUM           NUMBER(19),
  CONSTRAINT PK_VERSI_KONTRAK PRIMARY KEY (ID_VERSI_KONTRAK),
  CONSTRAINT UQ_VERSI_KONTRAK UNIQUE (ID_KONTRAK, NOMOR_URUT_VERSI),
  CONSTRAINT FK_VERSI_KONTRAK_1 FOREIGN KEY (ID_KONTRAK)
    REFERENCES {skema}.KONTRAK (ID_KONTRAK),
  CONSTRAINT FK_VERSI_KONTRAK_MATA_UANG FOREIGN KEY (KODE_MATA_UANG_KONTRAK)
    REFERENCES {skema}.MATA_UANG (ID_MATA_UANG)
)
/
CREATE INDEX {skema}.IX_VERSI_KONTRAK_KONTRAK ON {skema}.VERSI_KONTRAK (ID_KONTRAK)
/

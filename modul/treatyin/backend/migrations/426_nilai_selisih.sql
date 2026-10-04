-- NILAI_SELISIH dan NILAI_SEBELUM_PRO_RATE
--
-- Migrasi tiket 76 dan 77 (papan Treaty In Adjustment), BERKAS DI MODUL `treatyin`. Satu berkas = satu langkah
-- migrasi. Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis
-- miring.
--
-- Asal struktur: `4-erd-dan-tabel-datar/ERD-TREATY-IN-DAN-EDM.html`
-- baris relasi 36 dan 37 - ACUAN struktur sistem lama sejak keputusan pemilik
-- proses 2 Oktober 2026. Perilaku hapus: `ERD.md` §2.6, dokumen MENGIKAT.
--
-- ---------------------------------------------------------------------
-- KENAPA BERKAS INI DI MODUL `treatyin`, PADAHAL TIKETNYA PAPAN ADJUSTMENT
-- ---------------------------------------------------------------------
--   Ia sempat ditulis sebagai `442_` di `modul/treatyinadjustment/`, dan
--   penjaga modul itu MENOLAKNYA - `TestNolTabelBaru`: *"Modul ini TIDAK
--   membuat satu tabel pun, dan itu keputusan - bukan keadaan sementara.
--   Model datanya satu dengan Treaty In; yang dipisahkan 25-09-2026 hanya
--   papan tiketnya. Tabel baru di sini berarti model datanya bercabang."*
--
--   Penolakan itu benar, dan yang pindah berkasnya - bukan penjaganya.
--   Papan tiket terpisah; model data tidak. Tiket 76 dan 77 tetap di papan
--   Adjustment, dan tabelnya berdiri di tempat seluruh tabel modul ini
--   berdiri.
--
-- `NILAI_SELISIH` menutup lubang spec yang berdiri sejak papan ini dipisahkan:
-- `ERD.md` §2.6 menyatakannya dengan DUA relasi lintas sekat dan menyebutnya
-- "seluruh sambungan antarmodul, dua relasi dan tidak lebih" - lalu
-- `ddl-usulan/` dan `KAMUS-KOLOM.md` tidak memuatnya sama sekali. Ia menahan
-- tiket 06, 11, dan 13, dan ketiganya berstatus `aktif` sehingga lubang itu
-- tidak pernah terlihat di pencacah papan.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS: IKUT HAPUS, dari ERD.md §2.6
-- ---------------------------------------------------------------------
--   VERSI_KONTRAK 1--o< NILAI_SELISIH   [hapus: ikut hapus]   SEKAT
--
--   ERD HTML baris 36 juga menulis `CASCADE`. Keduanya sepakat.
--
--   ⛔ NOL RELASI KE "VERSI LAMA", dan §2.6 menegaskannya sendiri. Sisi lama
--   dibaca lewat `VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR` pada induknya (§2.2,
--   migrasi 440). Kolom `ID_VERSI_KONTRAK_LAMA` di tabel ini akan ditolak
--   tinjauan skema - ia melahirkan jalur kedua untuk fakta yang sama.
--
-- ---------------------------------------------------------------------
-- KARDINALITAS 1:N, DAN ERD HTML MENULIS 1:1. INI SEBABNYA.
-- ---------------------------------------------------------------------
--   ERD HTML baris 36 menulis `1:1`, dan catatan kakinya menjelaskan dari mana
--   angka itu datang: *"Kardinalitas dari bentuk simpul di pohon clipboard:
--   Page List = 1:N, Page = 1:1."* `TreatyIn.ValueDifference` memang sebuah
--   Page tunggal - SATU halaman berisi 217 simpul.
--
--   Model baru MENORMALKAN halaman itu menjadi baris, persis seperti ia
--   menormalkan 99 skalar akar menjadi tabel. Satu baris per besaran yang
--   berubah, dan `ADR-0048` butir 3 menuntut kunci bisnisnya MEMUAT MATA UANG
--   - tuntutan yang tidak punya arti pada relasi 1:1. Tiket 06 menyebut
--   akibatnya: mengubah mata uang tampil sebagai baris dihapus ditambah baris
--   baru, bukan sebagai selisih angka.
--
--   Jadi 1:1 benar tentang BENTUK HALAMAN sistem lama, 1:N benar tentang
--   bentuk model baru. Yang dipakai 1:N.
--
-- ---------------------------------------------------------------------
-- ⛔ PENGHALANG - kunci asing KEDUA tidak dapat dipasang
-- ---------------------------------------------------------------------
--   `ERD.md` §2.6 menuntut relasi kedua:
--
--     BESARAN_DAPAT_DISESUAIKAN 1--< NILAI_SELISIH   [hapus: tolak]   SEKAT
--
--   Tabel `BESARAN_DAPAT_DISESUAIKAN` TIDAK ADA DI MANA PUN - nol DDL, nol
--   `KAMUS-KOLOM.md`, nol tiket di kedua papan, dan §11.3 SPEC-MODEL-DATA
--   menyebutnya "tidak tersentuh satu pun kemampuan penyerahan pertama".
--
--   Di sini `KODE_BESARAN` berdiri sebagai TEKS, tanpa kunci asing. Itu
--   setengah dari yang §2.6 tuntut, dan setengahnya dinyatakan - bukan
--   disamarkan. SIAPA DAPAT MENJAWAB: pemilik proses - apakah besaran yang
--   dapat disesuaikan adalah tabel acuan tersendiri, atau cukup kode di sini.
--   SYARAT PEMBALIKAN: begitu tabelnya lahir, kolom ini dinaikkan menjadi
--   kunci asing `tolak` lewat migrasi korektif - dan HARUS sebelum data
--   dimuat, sebab kode yang tidak punya padanan akan menolak.
--
-- ---------------------------------------------------------------------
-- NILAI_SEBELUM_PRO_RATE BERINDUK VERSI_KONTRAK, DAN ERD MENULIS TREATY_IN
-- ---------------------------------------------------------------------
--   ERD HTML baris 37: `TREATY_IN 1:1 T_TREATY_VALUE_BEFORE_PRORATE` lewat
--   `TREATY_IN_ID`, bukti JALUR-PENUH. Di sistem lama `TREATY_IN` adalah SATU
--   baris; di model baru ia DIPECAH menjadi `KONTRAK` (lapisan beku) dan
--   `VERSI_KONTRAK` (yang berubah tiap penyesuaian) - ADR-0040. Jadi "induk
--   TREATY_IN" harus diterjemahkan ke salah satu dari keduanya, dan ERD tidak
--   dapat menjawabnya sebab pemecahan itu lahir sesudahnya.
--
--   Yang menjawabnya kolom `MEMAKAI_PRORATA`: ia berdiri di `VERSI_KONTRAK`
--   (migrasi 401), bukan di `KONTRAK`. Pro rata dipakai atau tidak ditentukan
--   PER VERSI, maka nilai sebelum pro rata pun milik versi. Menaruhnya di
--   `KONTRAK` membuat satu baris mewakili nilai yang berbeda di tiap versi.
--
--   SYARAT PEMBALIKAN: bila ditemukan nilai sebelum pro rata yang sama untuk
--   seluruh versi sebuah kontrak dan tidak pernah berbeda, induknya dapat
--   dinaikkan ke `KONTRAK`. Murah SELAMA tiket 44 belum memuat data.
--
-- ---------------------------------------------------------------------
-- NOL KOLOM SELISIH - INV-58
-- ---------------------------------------------------------------------
--   Tabel bernama `NILAI_SELISIH` TIDAK punya kolom bernama selisih, dan itu
--   disengaja. Yang disimpan NILAI_LAMA dan NILAI_BARU - dua fakta. Selisihnya
--   turunan keduanya, dan `INV-58` melarang menyimpan turunan. Menyimpannya
--   melahirkan kemungkinan ketiganya tidak konsisten, dan tidak ada yang akan
--   tahu mana dari ketiganya yang benar.
--
-- PERNYATAAN KEPUTUSAN - kunci alami kedua tabel (induk + kode besaran + mata
-- uang) TIDAK dipasang sebagai UNIQUE: `SPEC-INVARIAN.md` berhenti di INV-71
-- dan belum menomorinya. Aturan `ddl-usulan/Z00_KUNCI_ALAMI.sql` berlaku -
-- *"constraint tanpa invarian tidak punya tempat untuk gagal"*. Tagihannya di
-- `SPEC-INVARIAN.md` §13, dan ia BERTENGGAT terhadap tiket 44.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE SEQUENCE {skema}.SEQ_TRIA_NILAI_SELISIH START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIA_NILAI_SBL_PRORATA START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.NILAI_SELISIH (
  ID_NILAI_SELISIH  NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK  NUMBER(19)          NOT NULL,
  KODE_BESARAN      VARCHAR2(1000 CHAR) NOT NULL,
  KODE_MATA_UANG    VARCHAR2(1000 CHAR) NOT NULL,
  NILAI_LAMA        NUMBER(38,8),
  NILAI_BARU        NUMBER(38,8),
  CONSTRAINT PK_NILAI_SELISIH PRIMARY KEY (ID_NILAI_SELISIH),
  CONSTRAINT FK_NILAI_SELISIH_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK) ON DELETE CASCADE
)
/
CREATE TABLE {skema}.NILAI_SEBELUM_PRO_RATE (
  ID_NILAI_SEBELUM_PRO_RATE  NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK           NUMBER(19)          NOT NULL,
  KODE_BESARAN               VARCHAR2(1000 CHAR) NOT NULL,
  KODE_MATA_UANG             VARCHAR2(1000 CHAR) NOT NULL,
  NILAI                      NUMBER(38,8),
  CONSTRAINT PK_NILAI_SEBELUM_PRO_RATE PRIMARY KEY (ID_NILAI_SEBELUM_PRO_RATE),
  CONSTRAINT FK_NILAI_SEBELUM_PRO_RATE_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_NILAI_SELISIH_VERSI ON {skema}.NILAI_SELISIH (ID_VERSI_KONTRAK)
/
CREATE INDEX {skema}.IX_NILAI_SBL_PRORATA_VERSI ON {skema}.NILAI_SEBELUM_PRO_RATE (ID_VERSI_KONTRAK)
/

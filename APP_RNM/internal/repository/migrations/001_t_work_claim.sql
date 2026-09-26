-- T_WORK_CLAIM - akar, lintas-lini
--
-- Migrasi tiket 14 Claim Life. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Aturan yang dijaga seluruh berkas di folder ini:
--   ADR-U-0027  seluruh kolom nullable kecuali kunci utama; wajib-isi di services
--   ADR-U-0003  ADR-U-0016  uang NUMBER(38,8), tidak pernah float
--   ADR-U-0029  nol COMMIT di teks SQL; transaksi dibuka-ditutup aplikasi
--   ADR-U-0006  identitas dari sequence, kecuali yang dinyatakan berformat
--
-- Tabel work lintas-lini (Life dan Non-Life). Satu baris mewakili satu work
-- object: bisa baris klaim, bisa baris kasus komite.
--
-- ID bertipe TEKS BERFORMAT (CLM-xxxxxx untuk klaim, KMT-xxxxxx untuk komite),
-- bukan angka sequence. Itu penyimpangan sadar dari ADR-U-0006, dicatat di
-- tiket 14, bukan dilanggar diam-diam.
--
-- [terbuka] Pembangkit nomor CLM- dan KMT- belum ditetapkan: siapa yang
-- membuatnya, apakah ada sequence di belakang prefiks, apakah di-reset per
-- tahun. Pemilik DBA atau work owner. Kolomnya dibuat; pengisiannya menunggu.
--
-- COVER_KEY menunjuk T_WORK_CLAIM.ID induknya (relasi 1: tabel ini menunjuk
-- dirinya sendiri). NULL bila baris tidak punya induk.
-- [keputusan work owner 26-09-2026, butir d] COVER_KEY DIPASANGI REFERENCES,
-- tetap nullable. Oracle karena itu menolak penunjuk yatim: sebuah baris hanya
-- boleh menunjuk induk yang benar-benar ada, atau tidak menunjuk sama sekali.
-- Ini FK yang menunjuk tabelnya sendiri, dan itu sah di Oracle.
--
-- TANPA "ON DELETE": aturan hapus bawaan Oracle adalah MENOLAK. Kaskade akan
-- keliru di sini - COVER_KEY penunjuk KE ATAS, bukan kepemilikan, sehingga
-- menghapus induk tidak boleh ikut menghapus anaknya. Akibatnya menghapus
-- baris work yang masih ditunjuk baris lain akan GAGAL dengan ORA-02292,
-- bukan diam-diam membuat penunjuk yatim.
-- [terbuka] Apakah jalur hapus tiket 15 memilih melepas penunjuknya lebih dulu
-- atau menolak menghapus - pemiliknya work owner, dan belum teruji: migrasi ini
-- belum pernah berjalan di Oracle mana pun.
CREATE TABLE {skema}.T_WORK_CLAIM (
  ID              VARCHAR2(32) NOT NULL,
  COVER_KEY       VARCHAR2(32),
  LINI            VARCHAR2(16),
  PY_POSITION     VARCHAR2(64),
  SENDTO_ADMIN    VARCHAR2(8),
  SENDTO_MEDICAL  VARCHAR2(8),
  TYPE            VARCHAR2(32),
  CASE_ID         VARCHAR2(64),
  CREATE_OP       VARCHAR2(64),
  CREATE_OP_NAME  VARCHAR2(128),
  TGL_UPDATE      DATE,
  CONSTRAINT PK_T_WORK_CLAIM PRIMARY KEY (ID),
  CONSTRAINT FK_WORK_COVER_KEY FOREIGN KEY (COVER_KEY)
    REFERENCES {skema}.T_WORK_CLAIM (ID)
)
/
CREATE INDEX {skema}.IX_WORK_CLAIM_COVER_KEY ON {skema}.T_WORK_CLAIM (COVER_KEY)
/

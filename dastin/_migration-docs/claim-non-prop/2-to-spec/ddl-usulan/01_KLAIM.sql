-- =====================================================================
-- 01_KLAIM.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 12. Aggregate root: satu klaim satu kejadian kerugian (ADR-0001).
-- Dijalankan sesudah 00_SKEMA_DAN_AKUN.sql, sebelum seluruh tabel lain —
-- empat belas foreign key menunjuk ke sini.
--
-- Presisi: ADR-0003 (accepted), kelompok AK-1. Skala kanonik 20.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
--
-- TIDAK ADA DML DI BERKAS INI.
-- =====================================================================


-- ---------------------------------------------------------------------
-- Kolom yang perlu penjelasan — dan kenapa bentuknya begitu
-- ---------------------------------------------------------------------
-- TANGGAL_KEJADIAN  TETAP setelah baris dibuat. Tidak ada CHECK yang dapat
--                   menegakkan ketetapan; ia jatuh ke aplikasi (SPEC bag. 6).
--                   Dasarnya bukan selera: sapuan 114 activity menemukan
--                   NOL Property-Set yang menulis .DateOfLoss — di sistem
--                   lama pun ia hanya pernah diisi dari layar, tidak pernah
--                   dihitung ulang. Yang di sini menegakkan apa yang di sana
--                   kebetulan berlaku.
--
-- PORSI_CEDANT      Persentase, bukan uang. Tidak ada kolom uang di tabel ini
-- PORSI_REASURADUR  sama sekali — uang tinggal di NILAI_KLAIM_MATA_UANG,
--                   tempat ADR-0029 mewajibkan nilai asli, mata uang, nilai
--                   rupiah, kurs, dan asal-usul kurs berdiri di satu baris.
--
-- STATUS_KLAIM      SENGAJA TANPA CHECK domain. Nilai sistem lama tidak
--                   terbaca lengkap (D1, D35: dari empat nilai yang didaftar
--                   memori hanya dua ada, dan satu yang ada tidak didaftar).
--                   Domain tertutup di sini akan menolak nilai yang sah.
--
-- PENANDA_BERSYARAT Biner, DAN ITU KEPUTUSAN KAMI — bukan warisan. SPEC 21.1
--                   melarang memberi domain biner kepada penanda sistem lama
--                   tanpa bukti domainnya; larangan itu tentang apa yang
--                   DIBACA dari sana, bukan tentang apa yang kita tetapkan
--                   sendiri di sini. Nilai lamanya masuk lewat pemetaan,
--                   bukan lewat penyalinan mentah.
--
-- LINI_USAHA        Disediakan sekarang meski belum ada yang mengisinya.
--                   Penentu syariah di sistem lama ada dua — node yang
--                   mengeksekusi, dan nama akun notifikasi — dan TIDAK SATU
--                   PUN ikut pindah. Apa pun jawaban G1, sistem baru
--                   memerlukan penanda yang berdiri sendiri (SPEC 21.3).
--                   TANPA CHECK domain: nilainya belum terbaca, dan menutup
--                   domain yang belum diketahui adalah cacat yang sama
--                   dengan STATUS_KLAIM.
--
-- PENUTUPAN_LAMA    NULL untuk klaim yang tidak ditutup di sistem lama.
--                   Berisi nama jalur penutupan lama bila ada. Sistem baru
--                   punya SATU jalur penutupan (SPEC 21.6); klaim lama yang
--                   tertutup lewat jalur pintas dimigrasi sebagai tertutup
--                   DAN DITANDAI, supaya asal-usulnya terbaca tanpa jalurnya
--                   ikut hidup. Menandai bukan menghakimi: ia menyatakan
--                   lewat jalur mana baris itu sampai ke keadaannya.
--
-- DIBUAT_OLEH       POTRET, bukan rujukan. Pengenal dan nama operator
-- DIBUAT_ATAS_NAMA  SEBAGAIMANA TERCATAT SAAT TINDAKAN TERJADI. Tidak ada
-- DIUBAH_OLEH       foreign key, tidak ada pencarian ke sistem kepegawaian,
-- DIUBAH_ATAS_NAMA  tidak ada penyegaran (ADR-0016). Pelaku sebuah tindakan
--                   adalah fakta yang terjadi pada satu saat, dan fakta itu
--                   tidak berubah ketika orangnya pindah bagian.
--                   ..._ATAS_NAMA boleh KOSONG: kosong berarti tidak ada
--                   perwakilan, BUKAN tidak diketahui (ADR-0019).


CREATE TABLE KLAIMNP.KLAIM
(
  ID_KLAIM               NUMBER(19) NOT NULL,
  NOMOR_KLAIM            VARCHAR2(32 CHAR) NOT NULL,
  NOMOR_KLAIM_CEDANT     VARCHAR2(64 CHAR),
  NOMOR_POLIS            VARCHAR2(64 CHAR) NOT NULL,
  NOMOR_POLIS_CEDANT     VARCHAR2(64 CHAR),
  ID_TREATY              VARCHAR2(64 CHAR) NOT NULL,
  NAMA_TERTANGGUNG       VARCHAR2(255 CHAR),
  LINI_USAHA             VARCHAR2(40 CHAR),
  TANGGAL_KEJADIAN       DATE NOT NULL,
  TANGGAL_LAPOR          DATE,
  TANGGAL_TERIMA         DATE,
  TREATY_MULAI           DATE NOT NULL,
  TREATY_AKHIR           DATE NOT NULL,
  PENYEBAB_KERUGIAN_ID   VARCHAR2(50 CHAR),
  PORSI_CEDANT           NUMBER(38,20),
  PORSI_REASURADUR       NUMBER(38,20),
  STATUS_KLAIM           VARCHAR2(40 CHAR) NOT NULL,
  PENUTUPAN_LAMA         VARCHAR2(64 CHAR),
  PENANDA_BERSYARAT      NUMBER(1) NOT NULL,
  CATATAN_BERSYARAT      VARCHAR2(2000 CHAR),
  KEADAAN_BARIS          VARCHAR2(16 CHAR) NOT NULL,
  DIBUAT_OLEH            VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA       VARCHAR2(128 CHAR),
  DIBUAT_PADA            TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH            VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA       VARCHAR2(128 CHAR),
  DIUBAH_PADA            TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.KLAIM ADD CONSTRAINT PK_KLAIM PRIMARY KEY (ID_KLAIM);
ALTER TABLE KLAIMNP.KLAIM ADD CONSTRAINT UQ_KLAIM_1 UNIQUE (NOMOR_KLAIM);
ALTER TABLE KLAIMNP.KLAIM ADD CONSTRAINT CK_KLAIM_1 CHECK (TREATY_MULAI <= TREATY_AKHIR);
ALTER TABLE KLAIMNP.KLAIM ADD CONSTRAINT CK_KLAIM_2 CHECK (PENANDA_BERSYARAT IN (0,1));
ALTER TABLE KLAIMNP.KLAIM ADD CONSTRAINT CK_KLAIM_3 CHECK (KEADAAN_BARIS IN ('LENGKAP','MENUNGGU_KURS','GAGAL_URAI'));
CREATE INDEX KLAIMNP.IX_KLAIM_1 ON KLAIMNP.KLAIM (NOMOR_POLIS, TANGGAL_KEJADIAN);

-- ---------------------------------------------------------------------
-- PEMBANGKIT PENGENAL — ditambahkan 19 September 2026
-- ---------------------------------------------------------------------
-- Sampai tanggal ini, SELURUH 22 tabel punya kunci primer NUMBER(19) dan
-- TIDAK ADA SATU PUN yang menyatakan dari mana nilainya datang: nol
-- IDENTITY, nol SEQUENCE, nol DEFAULT di 36 berkas. Hak CREATE SEQUENCE
-- sudah diberikan di 00_SKEMA_DAN_AKUN.sql — jadi niatnya ada; objeknya
-- tidak pernah dibuat.
--
-- KENAPA SEQUENCE, BUKAN "GENERATED AS IDENTITY"
-- Alasan yang sama yang menetapkan batas nama 30 byte: pilih bentuk yang
-- sah di SETIAP versi. IDENTITY baru ada sejak 12.1; sequence sah jauh
-- sebelumnya. Versi instance tujuan masih menunggu REQ-032, dan REQ-032
-- sudah diturunkan menjadi verifikasi justru dengan janji bahwa jawabannya
-- TIDAK MENGUBAH APA PUN. Memakai IDENTITY akan membatalkan janji itu.
--
-- Nilai diminta pemanggil lewat NEXTVAL. Itu konsisten dengan ADR-0017:
-- satu pintu tulis, dan pintu itu yang meminta nomornya.
--
-- Penamaan mengikuti aturan 2 dan 3 (SPEC bagian 16): SQ_<tabel>, dengan
-- SQ_ sebagai awalan peran — sekelas PK_, UQ_, FK_, CK_, IX_, V_.

CREATE SEQUENCE KLAIMNP.SQ_KLAIM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_KLAIM TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.KLAIM TO KLAIMNP_APP;


-- ---------------------------------------------------------------------
-- Apa yang TIDAK ditegakkan di sini, dan jatuh ke aplikasi
-- ---------------------------------------------------------------------
-- Ditulis supaya berkas ini tidak dibaca sebagai jaminan yang lebih rapat
-- daripada yang sesungguhnya (SPEC bagian 6).
--
--   1. TANGGAL_KEJADIAN tetap setelah baris dibuat.
--   2. Domain STATUS_KLAIM dan LINI_USAHA — sengaja terbuka.
--   3. PENUTUPAN_LAMA hanya boleh terisi bila klaimnya memang tertutup;
--      keterkaitan itu bergantung nilai STATUS_KLAIM yang domainnya terbuka,
--      jadi CHECK atasnya akan mengunci domain lewat pintu belakang.
--   4. Ketetapan urutan penulisan. ADR-0013 membuat hasil tidak bergantung
--      urutan, dan itu sifat perhitungan, bukan sifat tabel.
--
-- Yang DITEGAKKAN basis data: keunikan nomor klaim, urutan masa berlaku
-- treaty, domain KEADAAN_BARIS, domain PENANDA_BERSYARAT, dan keharusan
-- adanya pelaku pembuat.
-- =====================================================================

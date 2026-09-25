-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Unit transaksi pembayaran klaim. Satu klaim banyak Adjustment.

CREATE TABLE KLAIMNP.ADJUSTMENT
(
  ID_ADJUSTMENT          NUMBER(19) NOT NULL,
  ID_KLAIM               NUMBER(19) NOT NULL,
  NOMOR_URUT             NUMBER(9) NOT NULL,
  ID_REKENING            NUMBER(19),
  JENIS_PEMBAYARAN       VARCHAR2(24 CHAR),
  MATA_UANG              VARCHAR2(3 CHAR) NOT NULL,
  KURS                   NUMBER(38,20),
  KURS_TANGGAL           DATE,
  KURS_SUMBER            VARCHAR2(32 CHAR),
  NILAI_ADJUSTMENT       NUMBER(38,20),
  NILAI_ADJUSTMENT_IDR   NUMBER(38,20),
  KEPUTUSAN_KOMITE       VARCHAR2(16 CHAR),
  NOMOR_AKSEPTASI        VARCHAR2(64 CHAR),
  TANGGAL_AKSEPTASI      DATE,
  PENANDA_BERSYARAT      NUMBER(1) NOT NULL,
  CATATAN_BERSYARAT      VARCHAR2(2000 CHAR),
  USUL_TUTUP             NUMBER(1) NOT NULL,
  LANGSUNG_KE_KASIR      NUMBER(1) NOT NULL,
  KEADAAN_BARIS          VARCHAR2(16 CHAR) NOT NULL,
  DIBUAT_OLEH            VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA       VARCHAR2(128 CHAR),
  DIBUAT_PADA            TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH            VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA       VARCHAR2(128 CHAR),
  DIUBAH_PADA            TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT PK_ADJUSTMENT PRIMARY KEY (ID_ADJUSTMENT);
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT UQ_ADJUSTMENT_1 UNIQUE (ID_KLAIM, NOMOR_URUT);
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT CK_ADJUSTMENT_1 CHECK (PENANDA_BERSYARAT IN (0,1));
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT CK_ADJUSTMENT_2 CHECK (USUL_TUTUP IN (0,1));
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT CK_ADJUSTMENT_3 CHECK (LANGSUNG_KE_KASIR IN (0,1));
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT CK_ADJUSTMENT_4 CHECK (KEADAAN_BARIS IN ('LENGKAP','MENUNGGU_KURS','GAGAL_URAI'));
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT CK_ADJUSTMENT_5 CHECK (NILAI_ADJUSTMENT_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT FK_ADJUSTMENT_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT FK_ADJUSTMENT_2 FOREIGN KEY (ID_REKENING) REFERENCES KLAIMNP.REKENING_PENERIMA (ID_REKENING);
CREATE INDEX KLAIMNP.IX_ADJUSTMENT_1 ON KLAIMNP.ADJUSTMENT (ID_REKENING);

-- ADR-0029 menuntut LIMA hal per nilai uang: nilai asli, kode mata uang,
-- nilai rupiah, KURS YANG DIPAKAI, dan ASAL-USUL KURS ITU. Empat yang
-- pertama sudah ditegakkan; yang kelima belum, sampai baris ini.
-- Ia bukan kerapian. FINDING-006: GETCURRENCYSTANDARD mengembalikan 1
-- ketika kurs TIDAK DITEMUKAN, dan angka 1 tidak dapat dibedakan dari
-- kurs yang sah. Yang membedakannya hanya asal-usulnya. Kurs tanpa
-- asal-usul mewarisi cacat itu utuh.
ALTER TABLE KLAIMNP.ADJUSTMENT ADD CONSTRAINT CK_ADJUSTMENT_6 CHECK (
  KURS IS NULL OR (KURS_SUMBER IS NOT NULL AND KURS_TANGGAL IS NOT NULL)
);

-- ---------------------------------------------------------------------
-- KEPUTUSAN_KOMITE boleh KOSONG, dan itu keadaan keempat — bukan kelalaian
-- ---------------------------------------------------------------------
-- Versi sebelumnya memasang NOT NULL. Itu membuat keadaan keempat MUSTAHIL
-- DISIMPAN, dan keadaan keempat justru yang paling penting.
--
-- Sapuan S12 menemukan properti ini dibandingkan sebagai ANGKA (0, 1, 2) dan
-- sebagai TEKS, dan bentuk teksnya memuat NILAI KOSONG yang tidak punya
-- padanan numerik. Jadi keadaannya EMPAT, bukan tiga.
--
-- Dan yang kosong itu bukan keadaan sepele: penjaga CloseClaimMD menolak
-- AcceptanceStatus == "0", sementara Adjustment yang BELUM PERNAH DIKIRIM ke
-- komite bernilai KOSONG, bukan "0" — sehingga ia LOLOS (D2). Seluruh
-- pembacaan jalur tutup-langsung bergantung pada perbedaan itu.
--
-- NULL di sini berarti BELUM PERNAH DIKIRIM KE KOMITE. Ia berbeda dari '0',
-- yang berarti SUDAH DIKIRIM DAN BELUM DIPUTUS. ADR-0019: kosong bukan
-- sinonim nol, dan di kolom ini perbedaan itu punya akibat.
--
-- Tanpa CHECK domain, dan itu keputusan (tiket 18): domainnya EXTERNAL —
-- modul Komite yang menulisnya, dan folder itu tertutup untuk pekerjaan ini.
-- Domain tertutup di sini akan menolak nilai yang sah dari seberang.

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

CREATE SEQUENCE KLAIMNP.SQ_ADJUSTMENT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_ADJUSTMENT TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.ADJUSTMENT TO KLAIMNP_APP;

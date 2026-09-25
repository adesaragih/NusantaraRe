-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Satu entitas dengan keadaan. Satu per klaim per layer per mata uang.

CREATE TABLE KLAIMNP.AKSEPTASI
(
  ID_AKSEPTASI          NUMBER(19) NOT NULL,
  ID_KLAIM              NUMBER(19) NOT NULL,
  LAYER                 VARCHAR2(50 CHAR) NOT NULL,
  LAYER_JENIS           VARCHAR2(50 CHAR) NOT NULL,
  LAYER_BAGIAN          VARCHAR2(50 CHAR) NOT NULL,
  LAYER_BAGIAN_JENIS    VARCHAR2(50 CHAR) NOT NULL,
  JENIS_REASURANSI      VARCHAR2(100 CHAR),
  JENIS_REASURANSI_ID   VARCHAR2(50 CHAR),
  MATA_UANG             VARCHAR2(3 CHAR) NOT NULL,
  KURS                  NUMBER(38,20),
  KURS_TANGGAL          DATE,
  KURS_SUMBER           VARCHAR2(32 CHAR),
  NILAI_DISETUJUI       NUMBER(38,20),
  NILAI_DISETUJUI_IDR   NUMBER(38,20),
  KEADAAN_AKSEPTASI     VARCHAR2(24 CHAR) NOT NULL,
  NOMOR_AKSEPTASI       VARCHAR2(64 CHAR),
  TANGGAL_AKSEPTASI     DATE,
  KEADAAN_BARIS         VARCHAR2(16 CHAR) NOT NULL,
  DIBUAT_OLEH           VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA      VARCHAR2(128 CHAR),
  DIBUAT_PADA           TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH           VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA      VARCHAR2(128 CHAR),
  DIUBAH_PADA           TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.AKSEPTASI ADD CONSTRAINT PK_AKSEPTASI PRIMARY KEY (ID_AKSEPTASI);
ALTER TABLE KLAIMNP.AKSEPTASI ADD CONSTRAINT UQ_AKSEPTASI_1 UNIQUE (ID_KLAIM, LAYER, LAYER_JENIS, LAYER_BAGIAN, LAYER_BAGIAN_JENIS, MATA_UANG);
ALTER TABLE KLAIMNP.AKSEPTASI ADD CONSTRAINT CK_AKSEPTASI_1 CHECK (KEADAAN_AKSEPTASI IN ('BIASA','BERSYARAT','BERSYARAT_GUGUR'));
ALTER TABLE KLAIMNP.AKSEPTASI ADD CONSTRAINT CK_AKSEPTASI_2 CHECK (KEADAAN_BARIS IN ('LENGKAP','MENUNGGU_KURS','GAGAL_URAI'));
ALTER TABLE KLAIMNP.AKSEPTASI ADD CONSTRAINT CK_AKSEPTASI_3 CHECK (NILAI_DISETUJUI_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.AKSEPTASI ADD CONSTRAINT FK_AKSEPTASI_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);

-- ADR-0029 menuntut LIMA hal per nilai uang: nilai asli, kode mata uang,
-- nilai rupiah, KURS YANG DIPAKAI, dan ASAL-USUL KURS ITU. Empat yang
-- pertama sudah ditegakkan; yang kelima belum, sampai baris ini.
-- Ia bukan kerapian. FINDING-006: GETCURRENCYSTANDARD mengembalikan 1
-- ketika kurs TIDAK DITEMUKAN, dan angka 1 tidak dapat dibedakan dari
-- kurs yang sah. Yang membedakannya hanya asal-usulnya. Kurs tanpa
-- asal-usul mewarisi cacat itu utuh.
ALTER TABLE KLAIMNP.AKSEPTASI ADD CONSTRAINT CK_AKSEPTASI_4 CHECK (
  KURS IS NULL OR (KURS_SUMBER IS NOT NULL AND KURS_TANGGAL IS NOT NULL)
);

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

CREATE SEQUENCE KLAIMNP.SQ_AKSEPTASI START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_AKSEPTASI TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.AKSEPTASI TO KLAIMNP_APP;

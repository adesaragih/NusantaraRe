-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Premi untuk memulihkan kapasitas layer, beserta seluruh masukan rumusnya.

-- ---------------------------------------------------------------------
-- RUMUS MANA YANG DIWARISI — dan yang mana TIDAK
-- ---------------------------------------------------------------------
-- Sistem lama punya DUA rumus premi pemulihan yang bekerja atas nilai yang
-- BERBEDA (FINDING-007). Yang diwarisi CountReinstatement_Act; yang TIDAK
-- diwarisi AdjClaimCNP_Act baris 3109 (AK-2b).
--
-- Itu bukan detail yang boleh tinggal di dokumen saja. Kolom di bawah adalah
-- masukan rumus PERTAMA. Bila kelak ada yang mengisinya dari rumus kedua,
-- angkanya akan tampak sah dan hasilnya berbeda tanpa satu pun tanda.
--
-- Selisih terhadap sistem lama karena itu DIHARAPKAN, dan ia PENGECUALIAN
-- BERNAMA pada shadow-run — bukan kegagalan cutover. Besarannya belum
-- terukur, dan itu dinyatakan, bukan didiamkan.
--
-- LIMIT_LAYER adalah PENYEBUT (MEMORI_PEMAHAMAN.MD 6.3). Nol di penyebut
-- bukan nilai yang aneh — ia perhitungan yang tidak dapat dijalankan. Karena
-- itu CK_PREMI_PEMULIHAN_1 menolaknya di baris, bukan menyerahkannya ke kode
-- yang harus ingat memeriksanya.
--
-- Seluruh masukan rumus ada di SATU baris: LIMIT_LAYER, PREMI_DEPOSIT,
-- PERSEN_PREMI_PEMULIHAN, PORSI_REASURADUR, NILAI_DASAR, BIAYA_PENILAIAN,
-- SALVAGE_DISESUAIKAN. Tidak satu pun perlu dibaca dari tabel master treaty
-- — itu yang membuat hitung ulang idempoten mungkin (ADR-0011).
--
-- Keempatnya NOT NULL karena rumus tanpa salah satunya bukan rumus yang
-- kekurangan data; ia rumus yang tidak dapat dijalankan sama sekali.

CREATE TABLE KLAIMNP.PREMI_PEMULIHAN
(
  ID_PREMI_PEMULIHAN          NUMBER(19) NOT NULL,
  ID_KLAIM                    NUMBER(19) NOT NULL,
  LAYER                       VARCHAR2(50 CHAR) NOT NULL,
  LAYER_JENIS                 VARCHAR2(50 CHAR) NOT NULL,
  LAYER_BAGIAN                VARCHAR2(50 CHAR) NOT NULL,
  LAYER_BAGIAN_JENIS          VARCHAR2(50 CHAR) NOT NULL,
  MATA_UANG                   VARCHAR2(3 CHAR) NOT NULL,
  KURS                        NUMBER(38,20),
  KURS_TANGGAL                DATE,
  KURS_SUMBER                 VARCHAR2(32 CHAR),
  NILAI_DASAR                 NUMBER(38,20),
  NILAI_DASAR_IDR             NUMBER(38,20),
  BIAYA_PENILAIAN             NUMBER(38,20),
  BIAYA_PENILAIAN_IDR         NUMBER(38,20),
  SALVAGE_DISESUAIKAN         NUMBER(38,20),
  SALVAGE_DISESUAIKAN_IDR     NUMBER(38,20),
  LIMIT_LAYER                 NUMBER(38,20) NOT NULL,
  PREMI_DEPOSIT               NUMBER(38,20) NOT NULL,
  PERSEN_PREMI_PEMULIHAN      NUMBER(38,20) NOT NULL,
  PORSI_REASURADUR            NUMBER(38,20) NOT NULL,
  PREMI_PEMULIHAN             NUMBER(38,20),
  PREMI_PEMULIHAN_IDR         NUMBER(38,20),
  PREMI_PEMULIHAN_PORSI       NUMBER(38,20),
  PREMI_PEMULIHAN_PORSI_IDR   NUMBER(38,20),
  KEADAAN_BARIS               VARCHAR2(16 CHAR) NOT NULL,
  DIBUAT_OLEH                 VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA            VARCHAR2(128 CHAR),
  DIBUAT_PADA                 TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH                 VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA            VARCHAR2(128 CHAR),
  DIUBAH_PADA                 TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT PK_PREMI_PEMULIHAN PRIMARY KEY (ID_PREMI_PEMULIHAN);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT UQ_PREMI_PEMULIHAN_1 UNIQUE (ID_KLAIM, LAYER, LAYER_JENIS, LAYER_BAGIAN, LAYER_BAGIAN_JENIS, MATA_UANG);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_1 CHECK (LIMIT_LAYER <> 0);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_2 CHECK (KEADAAN_BARIS IN ('LENGKAP','MENUNGGU_KURS','GAGAL_URAI'));
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_3 CHECK (NILAI_DASAR_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_4 CHECK (BIAYA_PENILAIAN_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_5 CHECK (SALVAGE_DISESUAIKAN_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_6 CHECK (PREMI_PEMULIHAN_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_7 CHECK (PREMI_PEMULIHAN_PORSI_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT FK_PREMI_PEMULIHAN_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);

-- ADR-0029 menuntut LIMA hal per nilai uang: nilai asli, kode mata uang,
-- nilai rupiah, KURS YANG DIPAKAI, dan ASAL-USUL KURS ITU. Empat yang
-- pertama sudah ditegakkan; yang kelima belum, sampai baris ini.
-- Ia bukan kerapian. FINDING-006: GETCURRENCYSTANDARD mengembalikan 1
-- ketika kurs TIDAK DITEMUKAN, dan angka 1 tidak dapat dibedakan dari
-- kurs yang sah. Yang membedakannya hanya asal-usulnya. Kurs tanpa
-- asal-usul mewarisi cacat itu utuh.
ALTER TABLE KLAIMNP.PREMI_PEMULIHAN ADD CONSTRAINT CK_PREMI_PEMULIHAN_8 CHECK (
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

CREATE SEQUENCE KLAIMNP.SQ_PREMI_PEMULIHAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_PREMI_PEMULIHAN TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.PREMI_PEMULIHAN TO KLAIMNP_APP;

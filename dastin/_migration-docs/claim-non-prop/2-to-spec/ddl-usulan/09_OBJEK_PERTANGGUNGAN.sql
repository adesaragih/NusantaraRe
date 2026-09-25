-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Objek yang dipertanggungkan. Hanya InterestList; ObjectList masuk daftar lubang.

CREATE TABLE KLAIMNP.OBJEK_PERTANGGUNGAN
(
  ID_OBJEK                  NUMBER(19) NOT NULL,
  ID_KLAIM                  NUMBER(19) NOT NULL,
  NOMOR_URUT                NUMBER(9) NOT NULL,
  NAMA_OBJEK                VARCHAR2(255 CHAR) NOT NULL,
  MATA_UANG                 VARCHAR2(3 CHAR) NOT NULL,
  KURS                      NUMBER(38,20),
  KURS_TANGGAL              DATE,
  KURS_SUMBER               VARCHAR2(32 CHAR),
  NILAI_PERTANGGUNGAN       NUMBER(38,20),
  NILAI_PERTANGGUNGAN_IDR   NUMBER(38,20),
  KEADAAN_BARIS             VARCHAR2(16 CHAR) NOT NULL,
  DIBUAT_OLEH               VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA          VARCHAR2(128 CHAR),
  DIBUAT_PADA               TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH               VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA          VARCHAR2(128 CHAR),
  DIUBAH_PADA               TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.OBJEK_PERTANGGUNGAN ADD CONSTRAINT PK_OBJEK_PERTANGGUNGAN PRIMARY KEY (ID_OBJEK);
ALTER TABLE KLAIMNP.OBJEK_PERTANGGUNGAN ADD CONSTRAINT UQ_OBJEK_PERTANGGUNGAN_1 UNIQUE (ID_KLAIM, NOMOR_URUT);
ALTER TABLE KLAIMNP.OBJEK_PERTANGGUNGAN ADD CONSTRAINT CK_OBJEK_PERTANGGUNGAN_1 CHECK (KEADAAN_BARIS IN ('LENGKAP','MENUNGGU_KURS','GAGAL_URAI'));
ALTER TABLE KLAIMNP.OBJEK_PERTANGGUNGAN ADD CONSTRAINT CK_OBJEK_PERTANGGUNGAN_2 CHECK (NILAI_PERTANGGUNGAN_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.OBJEK_PERTANGGUNGAN ADD CONSTRAINT FK_OBJEK_PERTANGGUNGAN_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);

-- ADR-0029 menuntut LIMA hal per nilai uang: nilai asli, kode mata uang,
-- nilai rupiah, KURS YANG DIPAKAI, dan ASAL-USUL KURS ITU. Empat yang
-- pertama sudah ditegakkan; yang kelima belum, sampai baris ini.
-- Ia bukan kerapian. FINDING-006: GETCURRENCYSTANDARD mengembalikan 1
-- ketika kurs TIDAK DITEMUKAN, dan angka 1 tidak dapat dibedakan dari
-- kurs yang sah. Yang membedakannya hanya asal-usulnya. Kurs tanpa
-- asal-usul mewarisi cacat itu utuh.
ALTER TABLE KLAIMNP.OBJEK_PERTANGGUNGAN ADD CONSTRAINT CK_OBJEK_PERTANGGUNGAN_3 CHECK (
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

CREATE SEQUENCE KLAIMNP.SQ_OBJEK_PERTANGGUNGAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_OBJEK_PERTANGGUNGAN TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.OBJEK_PERTANGGUNGAN TO KLAIMNP_APP;

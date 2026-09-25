-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Hasil alokasi kerugian per layer. Tanpa baris Retensi Cedant.

CREATE TABLE KLAIMNP.ALOKASI_LAYER
(
  ID_ALOKASI                 NUMBER(19) NOT NULL,
  ID_KLAIM                   NUMBER(19) NOT NULL,
  LAYER                      VARCHAR2(50 CHAR) NOT NULL,
  LAYER_JENIS                VARCHAR2(50 CHAR) NOT NULL,
  LAYER_BAGIAN               VARCHAR2(50 CHAR) NOT NULL,
  LAYER_BAGIAN_JENIS         VARCHAR2(50 CHAR) NOT NULL,
  JENIS_REASURANSI           VARCHAR2(100 CHAR),
  JENIS_REASURANSI_ID        VARCHAR2(50 CHAR),
  MATA_UANG                  VARCHAR2(3 CHAR) NOT NULL,
  KURS                       NUMBER(38,20),
  KURS_TANGGAL               DATE,
  KURS_SUMBER                VARCHAR2(32 CHAR),
  NILAI_TERALOKASI_HITUNG    NUMBER(38,20),
  NILAI_TERALOKASI_SUNTING   NUMBER(38,20),
  NILAI_TERALOKASI_DIPAKAI   NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(NILAI_TERALOKASI_SUNTING, NILAI_TERALOKASI_HITUNG)) VIRTUAL,
  BIAYA_PENILAIAN_HITUNG     NUMBER(38,20),
  BIAYA_PENILAIAN_SUNTING    NUMBER(38,20),
  BIAYA_PENILAIAN_DIPAKAI    NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(BIAYA_PENILAIAN_SUNTING, BIAYA_PENILAIAN_HITUNG)) VIRTUAL,
  SALVAGE_HITUNG             NUMBER(38,20),
  SALVAGE_SUNTING            NUMBER(38,20),
  SALVAGE_DIPAKAI            NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(SALVAGE_SUNTING, SALVAGE_HITUNG)) VIRTUAL,
  BIAYA_LAIN_HITUNG          NUMBER(38,20),
  BIAYA_LAIN_SUNTING         NUMBER(38,20),
  BIAYA_LAIN_DIPAKAI         NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(BIAYA_LAIN_SUNTING, BIAYA_LAIN_HITUNG)) VIRTUAL,
  DISUNTING_OLEH             VARCHAR2(128 CHAR),
  DISUNTING_ATAS_NAMA        VARCHAR2(128 CHAR),
  DISUNTING_PADA             TIMESTAMP(6),
  ALASAN_SUNTING             VARCHAR2(500 CHAR),
  NILAI_TERALOKASI_IDR       NUMBER(38,20),
  PORSI_REASURADUR           NUMBER(38,20),
  PORSI_CEDANT               NUMBER(38,20),
  PRORATA_KLAIM              NUMBER(38,20),
  LIMIT_LAYER                NUMBER(38,20),
  PREMI_DEPOSIT              NUMBER(38,20),
  PERSEN_PREMI_PEMULIHAN     NUMBER(38,20),
  URUTAN_PENGISIAN           NUMBER(9) NOT NULL,
  KEADAAN_BARIS              VARCHAR2(16 CHAR) NOT NULL,
  DIBUAT_OLEH                VARCHAR2(128 CHAR) NOT NULL,
  DIBUAT_ATAS_NAMA           VARCHAR2(128 CHAR),
  DIBUAT_PADA                TIMESTAMP(6) NOT NULL,
  DIUBAH_OLEH                VARCHAR2(128 CHAR),
  DIUBAH_ATAS_NAMA           VARCHAR2(128 CHAR),
  DIUBAH_PADA                TIMESTAMP(6)
);

ALTER TABLE KLAIMNP.ALOKASI_LAYER ADD CONSTRAINT PK_ALOKASI_LAYER PRIMARY KEY (ID_ALOKASI);
ALTER TABLE KLAIMNP.ALOKASI_LAYER ADD CONSTRAINT UQ_ALOKASI_LAYER_1 UNIQUE (ID_KLAIM, LAYER, LAYER_JENIS, LAYER_BAGIAN, LAYER_BAGIAN_JENIS, MATA_UANG);
ALTER TABLE KLAIMNP.ALOKASI_LAYER ADD CONSTRAINT CK_ALOKASI_LAYER_1 CHECK (NILAI_TERALOKASI_IDR IS NULL OR KURS IS NOT NULL);
ALTER TABLE KLAIMNP.ALOKASI_LAYER ADD CONSTRAINT CK_ALOKASI_LAYER_2 CHECK (KEADAAN_BARIS IN ('LENGKAP','MENUNGGU_KURS','GAGAL_URAI'));
ALTER TABLE KLAIMNP.ALOKASI_LAYER ADD CONSTRAINT FK_ALOKASI_LAYER_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);

-- ADR-0029 menuntut LIMA hal per nilai uang: nilai asli, kode mata uang,
-- nilai rupiah, KURS YANG DIPAKAI, dan ASAL-USUL KURS ITU. Empat yang
-- pertama sudah ditegakkan; yang kelima belum, sampai baris ini.
-- Ia bukan kerapian. FINDING-006: GETCURRENCYSTANDARD mengembalikan 1
-- ketika kurs TIDAK DITEMUKAN, dan angka 1 tidak dapat dibedakan dari
-- kurs yang sah. Yang membedakannya hanya asal-usulnya. Kurs tanpa
-- asal-usul mewarisi cacat itu utuh.
ALTER TABLE KLAIMNP.ALOKASI_LAYER ADD CONSTRAINT CK_ALOKASI_LAYER_3 CHECK (
  KURS IS NULL OR (KURS_SUMBER IS NOT NULL AND KURS_TANGGAL IS NOT NULL)
);

-- Tiket 17 butir 5: SUNTINGAN SELALU PUNYA PELAKU. Kriteria itu ada sejak
-- tiketnya ditulis dan tidak pernah punya constraint - dipasang di sini.
-- Sebabnya ADR-0008: suntingan manual bertahan DAN TERLIHAT. Nilai yang
-- menimpa hasil hitung tanpa ada yang bertanggung jawab atasnya bertahan
-- tanpa terlihat, dan itu setengah dari yang dijanjikan.
-- DISUNTING_ATAS_NAMA tetap boleh kosong: kosong berarti tidak ada
-- perwakilan, bukan tidak diketahui (ADR-0019).
ALTER TABLE KLAIMNP.ALOKASI_LAYER ADD CONSTRAINT CK_ALOKASI_LAYER_4 CHECK (
  (NILAI_TERALOKASI_SUNTING IS NULL OR DISUNTING_OLEH IS NOT NULL)
  AND (BIAYA_PENILAIAN_SUNTING IS NULL OR DISUNTING_OLEH IS NOT NULL)
  AND (SALVAGE_SUNTING IS NULL OR DISUNTING_OLEH IS NOT NULL)
  AND (BIAYA_LAIN_SUNTING IS NULL OR DISUNTING_OLEH IS NOT NULL)
  AND (ALASAN_SUNTING IS NULL OR DISUNTING_OLEH IS NOT NULL)
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

CREATE SEQUENCE KLAIMNP.SQ_ALOKASI_LAYER START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_ALOKASI_LAYER TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.ALOKASI_LAYER TO KLAIMNP_APP;

-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026). Batas 30 byte sebagai
-- aturan tetap; constraint berbentuk <peran>_<tabel>[_n], tidak mengeja kolom.
-- Apa yang BENAR-BENAR dikirim ke hilir, bukan apa yang seharusnya. Tulis-sekali, tidak pernah diubah.

CREATE TABLE KLAIMNP.ARSIP_MUATAN_KELUAR
(
  ID_ARSIP           NUMBER(19) NOT NULL,
  TUJUAN             VARCHAR2(24 CHAR) NOT NULL,
  ID_KLAIM           NUMBER(19) NOT NULL,
  JENIS_REASURANSI   VARCHAR2(100 CHAR) NOT NULL,
  MATA_UANG          VARCHAR2(3 CHAR) NOT NULL,
  DITOLAK            NUMBER(1) NOT NULL,
  NILAI_TERSEBAR     NUMBER(38,2) NOT NULL,
  NILAI_BRUTO        NUMBER(38,2) NOT NULL,
  BIAYA_PENILAIAN    NUMBER(38,2) NOT NULL,
  SALVAGE            NUMBER(38,2) NOT NULL,
  BIAYA_LAIN         NUMBER(38,2) NOT NULL,
  KURS               NUMBER(38,20),
  DIKIRIM_OLEH       VARCHAR2(128 CHAR) NOT NULL,
  DIKIRIM_PADA       TIMESTAMP(6) NOT NULL
);

ALTER TABLE KLAIMNP.ARSIP_MUATAN_KELUAR ADD CONSTRAINT PK_ARSIP_MUATAN_KELUAR PRIMARY KEY (ID_ARSIP);
ALTER TABLE KLAIMNP.ARSIP_MUATAN_KELUAR ADD CONSTRAINT CK_ARSIP_MUATAN_KELUAR_1 CHECK (DITOLAK IN (0,1));
ALTER TABLE KLAIMNP.ARSIP_MUATAN_KELUAR ADD CONSTRAINT FK_ARSIP_MUATAN_KELUAR_1 FOREIGN KEY (ID_KLAIM) REFERENCES KLAIMNP.KLAIM (ID_KLAIM);
-- DITOLAK masuk index karena setiap penjumlahan selisih menyaringnya:
-- muatan bertanda ditolak TIDAK ikut terjumlah (tiket 25 butir 3).
CREATE INDEX KLAIMNP.IX_ARSIP_MUATAN_KELUAR_1 ON KLAIMNP.ARSIP_MUATAN_KELUAR (ID_KLAIM, JENIS_REASURANSI, MATA_UANG, DITOLAK, DIKIRIM_PADA);

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

CREATE SEQUENCE KLAIMNP.SQ_ARSIP_MUATAN_KELUAR START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_ARSIP_MUATAN_KELUAR TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
-- TULIS-SEKALI, DITEGAKKAN HAK AKSES — BUKAN KESEPAKATAN.
-- Akun aplikasi memegang INSERT dan SELECT saja. TANPA UPDATE, TANPA DELETE.
--
-- Versi sebelumnya memberi keempatnya, dan itu bertentangan dengan tiket 25
-- butir 1 serta tiket 35 yang sudah menyebut pengecualian ini sejak ditulis.
-- Diperbaiki 19 September 2026.
--
-- Sebabnya bukan kehati-hatian: arsip ini menjawab pertanyaan "APA YANG
-- BENAR-BENAR DIKIRIM". Arsip yang dapat diubah menjawab pertanyaan lain —
-- "apa yang sekarang kita katakan telah dikirim" — dan itu bukan catatan.
--
-- Dan ia bukan turunan yang dapat dihitung ulang: GetDataOS menjumlahkan
-- seluruh baris berkunci sama ber-STS_REJECT = 0, lalu SaveDataToOSAksep_Act
-- mengurangkan hasilnya. BARIS ADALAH TAMBAHAN, BUKAN KEADAAN. Satu UPDATE
-- mengubah setiap selisih yang pernah dihitung sesudahnya.
GRANT SELECT, INSERT ON KLAIMNP.ARSIP_MUATAN_KELUAR TO KLAIMNP_APP;

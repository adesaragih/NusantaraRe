-- =====================================================================
-- 21_MIGRASI_PENDARATAN.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 13, dikerjakan bersama tiket 14 — keduanya dua sisi satu seam.
-- Satu-satunya tempat bentuk lama boleh masuk utuh, termasuk JSON.
--
-- Presisi: ADR-0003 (accepted), kelompok AK-1.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026).
--
-- TIDAK ADA DML DI BERKAS INI.
-- =====================================================================


-- ---------------------------------------------------------------------
-- Kenapa tabel ini TIDAK menolak muatan yang bukan JSON sah
-- ---------------------------------------------------------------------
-- Versi pertama memasang CHECK (MUATAN IS JSON) sebagai GERBANG: muatan
-- yang tidak sah ditolak di pintu. Gerbang itu DICABUT 19 September 2026,
-- dan sebabnya dua, keduanya terbaca dari sistem lama.
--
-- PERTAMA — bentuk muatan lama TIDAK DIBATASI (D43). Tiga rule memanggil
-- adoptJSONObject atas teks HASIL1 TANPA memeriksa bentuknya, dan salah
-- satunya menulis hasilnya ke halaman objek kerja. Apa pun yang ada di
-- HASIL1 menjadi properti. Bentuk yang tidak dibatasi tidak dapat dijaga
-- dengan CHECK; ia hanya dapat DITAMPUNG.
--
-- KEDUA — sistem lama memang menghasilkan JSON rusak (D41). Muatan kasir
-- dirakit dengan penyambungan teks, dan Nett, Deductible, KaliDeduct
-- dikirim TANPA KUTIP sebagai angka. Nilai kosong menghasilkan  "Nett":,
-- yaitu JSON yang tidak sah. Gerbang IS JSON akan menolak persis muatan
-- yang paling perlu tercatat.
--
-- Dan menolaknya bertabrakan dengan AK-4: tidak ada baris yang dibuang
-- karena "cuma sedikit". Muatan yang ditolak di pintu tidak punya rumah.
--
-- YANG BERLAKU SEKARANG: muatan mendarat UTUH apa pun bentuknya, dan
-- hasil penguraiannya DICATAT sebagai fakta, bukan dipakai sebagai syarat.
-- Kolom BENTUK_TERURAI menyimpan hasil itu; SEBAB_GAGAL_URAI menyimpan
-- alasannya bila gagal. Yang gagal terurai TIDAK hilang — ia tetap di
-- sini, utuh, dan barisnya ditunjuk dari MIGRASI_NILAI_DITOLAK (tiket 14).
--
-- Itu SPEC 21.8 sebagai tabel: mendarat utuh, yang dikenali naik ke
-- kanonik, sisanya tercatat sebagai tak terurai BESERTA ASALNYA.


CREATE TABLE KLAIMNP.MIGRASI_PENDARATAN
(
  ID_PENDARATAN      NUMBER(19) NOT NULL,
  SUMBER_LAMA        VARCHAR2(64 CHAR) NOT NULL,
  PENGENAL_LAMA      VARCHAR2(500 CHAR) NOT NULL,
  MUATAN             CLOB NOT NULL,
  BENTUK_TERURAI     NUMBER(1) NOT NULL,
  SEBAB_GAGAL_URAI   VARCHAR2(255 CHAR),
  DITERIMA_PADA      TIMESTAMP(6) NOT NULL
);

ALTER TABLE KLAIMNP.MIGRASI_PENDARATAN ADD CONSTRAINT PK_MIGRASI_PENDARATAN PRIMARY KEY (ID_PENDARATAN);
ALTER TABLE KLAIMNP.MIGRASI_PENDARATAN ADD CONSTRAINT CK_MIGRASI_PENDARATAN_1 CHECK (BENTUK_TERURAI IN (0,1));
ALTER TABLE KLAIMNP.MIGRASI_PENDARATAN ADD CONSTRAINT CK_MIGRASI_PENDARATAN_2 CHECK (
  (BENTUK_TERURAI = 1 AND SEBAB_GAGAL_URAI IS NULL)
  OR
  (BENTUK_TERURAI = 0 AND SEBAB_GAGAL_URAI IS NOT NULL)
);
CREATE INDEX KLAIMNP.IX_MIGRASI_PENDARATAN_1 ON KLAIMNP.MIGRASI_PENDARATAN (SUMBER_LAMA, PENGENAL_LAMA);

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

CREATE SEQUENCE KLAIMNP.SQ_MIGRASI_PENDARATAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_MIGRASI_PENDARATAN TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.MIGRASI_PENDARATAN TO KLAIMNP_APP;


-- ---------------------------------------------------------------------
-- Apa yang TIDAK ditegakkan di sini
-- ---------------------------------------------------------------------
--   1. Kesahihan bentuk MUATAN — sengaja. Lihat catatan di atas.
--   2. Keharusan setiap muatan gagal-urai punya baris di
--      MIGRASI_NILAI_DITOLAK. Tidak dapat ditegakkan dari sisi ini:
--      satu muatan dapat melahirkan nol, satu, atau banyak nilai ditolak.
--      Ia diuji dengan kueri, bukan dijaga dengan constraint — tiket 34.
--   3. Umur baris. Tabel ini DIJATUHKAN setelah paritas diterima, sama
--      seperti MIGRASI_KORELASI. Ia jembatan, bukan penyimpanan.
--
-- YANG DITEGAKKAN: tiap baris menyebut sumber dan pengenal lamanya;
-- BENTUK_TERURAI biner; dan sebab hanya ada bila memang gagal — sehingga
-- "gagal tanpa sebab" dan "berhasil tapi bersebab" keduanya ditolak.
-- =====================================================================

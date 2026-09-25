-- =====================================================================
-- 20_MIGRASI_NILAI_DITOLAK.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 14, dikerjakan bersama tiket 13 — keduanya dua sisi satu seam.
-- Nilai lama yang tidak lolos aturan penguraian. TIDAK DIBULATKAN,
-- TIDAK DIBUANG.
--
-- Presisi: ADR-0003 (accepted), kelompok AK-1.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026).
--
-- TIDAK ADA DML DI BERKAS INI.
-- =====================================================================


-- ---------------------------------------------------------------------
-- Kenapa tabel ini ada
-- ---------------------------------------------------------------------
-- ADR-0003 menetapkan migrasi MENGURAI TEKS APA ADANYA dan menyimpan
-- seluruh digit yang ada. Nilai yang TIDAK terurai tidak dibulatkan dan
-- tidak dibuang — ia mendarat di sini.
--
-- Dasarnya berbukti, bukan berjaga-jaga: nilai uang lama tersimpan sebagai
-- TEKS di dalam JSON dan blob; CLAIMXOL mengeluarkannya sebagai varchar2;
-- dan Java memindahkan uang sebagai String 92 kali lewat .toString(), NOL
-- double, NOL float, NOL BigDecimal (D39). Teks yang tidak pernah dijaga
-- bentuknya pasti memuat sebagian yang tidak terurai.
--
-- AK-4: migrasi BERHENTI bila ada satu baris yang tidak dapat diurai DAN
-- tidak dapat diselesaikan. Tidak ada baris yang dibuang karena "cuma
-- sedikit". Tabel ini yang membuat aturan itu dapat diperiksa — tanpa
-- tempat mencatat, "tidak dibuang" hanya niat.
--
-- NILAI_MENTAH disimpan sebagai VARCHAR2, BUKAN NUMBER. Mengubahnya jadi
-- angka adalah persis hal yang gagal dilakukan; menyimpannya sebagai angka
-- akan membuang bukti bentuk aslinya.


-- ---------------------------------------------------------------------
-- Seam ke tiket 13
-- ---------------------------------------------------------------------
-- ID_PENDARATAN menunjuk baris MIGRASI_PENDARATAN asal nilai itu, bila
-- nilainya memang datang dari sebuah muatan. Ia BOLEH KOSONG: sebagian
-- nilai ditolak berasal langsung dari kolom tabel lama, bukan dari muatan.
-- Kosong berarti "bukan dari muatan", BUKAN "tidak diketahui" (ADR-0019) —
-- SUMBER_LAMA dan PENGENAL_LAMA tetap wajib, jadi asalnya selalu terbaca.
--
-- Satu muatan dapat melahirkan BANYAK nilai ditolak. Karena itu arah
-- relasinya dari sini ke sana, bukan sebaliknya.


CREATE TABLE KLAIMNP.MIGRASI_NILAI_DITOLAK
(
  ID_DITOLAK       NUMBER(19) NOT NULL,
  ID_PENDARATAN    NUMBER(19),
  SUMBER_LAMA      VARCHAR2(64 CHAR) NOT NULL,
  PENGENAL_LAMA    VARCHAR2(500 CHAR) NOT NULL,
  KOLOM_LAMA       VARCHAR2(64 CHAR) NOT NULL,
  NILAI_MENTAH     VARCHAR2(4000 CHAR) NOT NULL,
  SEBAB_DITOLAK    VARCHAR2(255 CHAR) NOT NULL,
  DITEMUKAN_PADA   TIMESTAMP(6) NOT NULL
);

ALTER TABLE KLAIMNP.MIGRASI_NILAI_DITOLAK ADD CONSTRAINT PK_MIGRASI_NILAI_DITOLAK PRIMARY KEY (ID_DITOLAK);
ALTER TABLE KLAIMNP.MIGRASI_NILAI_DITOLAK ADD CONSTRAINT FK_MIGRASI_NILAI_DITOLAK_1 FOREIGN KEY (ID_PENDARATAN) REFERENCES KLAIMNP.MIGRASI_PENDARATAN (ID_PENDARATAN);
CREATE INDEX KLAIMNP.IX_MIGRASI_NILAI_DITOLAK_1 ON KLAIMNP.MIGRASI_NILAI_DITOLAK (ID_PENDARATAN);

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

CREATE SEQUENCE KLAIMNP.SQ_MIGRASI_NILAI_DITOLAK START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_MIGRASI_NILAI_DITOLAK TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.MIGRASI_NILAI_DITOLAK TO KLAIMNP_APP;


-- ---------------------------------------------------------------------
-- Apa yang TIDAK ditegakkan di sini
-- ---------------------------------------------------------------------
--   1. Bahwa migrasi berhenti bila tabel ini berisi. Itu AK-4, dan ia
--      aturan proses — tidak ada CHECK yang dapat menghentikan orang.
--      Diuji lewat kueri cacah sebelum cutover.
--   2. Domain SEBAB_DITOLAK. Sengaja terbuka: sebab yang belum pernah
--      terjadi tidak dapat didaftar lebih dulu, dan domain tertutup akan
--      memaksa penulisnya berbohong memilih sebab terdekat.
--
-- YANG DITEGAKKAN: nilai mentah selalu ada dan selalu teks; sumber dan
-- pengenal lama selalu ada; sebab selalu ada; dan bila nilainya berasal
-- dari sebuah muatan, muatan itu benar-benar ada di MIGRASI_PENDARATAN.
-- =====================================================================

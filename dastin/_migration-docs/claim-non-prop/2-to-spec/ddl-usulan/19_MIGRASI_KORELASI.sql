-- =====================================================================
-- 19_MIGRASI_KORELASI.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 15. Jembatan ke sistem lama. DIJATUHKAN setelah paritas diterima.
-- Penamaan: SPEC bagian 16 (final 19 Sep 2026).
--
-- TIDAK ADA DML DI BERKAS INI.
-- =====================================================================


-- ---------------------------------------------------------------------
-- Kenapa LIMA kolom terpisah, bukan satu kolom teks
-- ---------------------------------------------------------------------
-- Sistem lama tidak punya SATU pengenal. Ia punya lima, dan kelimanya
-- cacat dengan cara yang berbeda-beda (BLUEPRINT 8.4a):
--
--   PZINSKEY_LAMA      pegangan instance Pega, berbentuk "<class> <pyID>".
--                      Menamai DI MANA barang tersimpan, bukan BARANG APA.
--                      Tetapi ia pengenal sisi lama TERBAIK yang kita punya,
--                      dan alasannya berbukti: PYID tidak punya unique
--                      constraint di tabel work, dan JSON_KLAIM mewajibkan
--                      MNK_NO_KLAIM ada tetapi berkunci primer IDPEGA —
--                      nomor klaim lama TIDAK DIPAKSA UNIK oleh struktur
--                      mana pun. Pegangan instance lebih andal daripada
--                      keduanya.
--   PYID_LAMA          nomor klaim. Tanpa unique constraint (D13).
--   CASEID_LAMA        salinan PZINSKEY di sisi basis data.
--   INDEX_OBJECT_LAMA  POSISI numerik di AdjustmentList, bukan pengenal.
--                      Berubah bila daftar disisipi atau diurut ulang.
--   URUTAN_BARIS_LAMA  posisi baris di dalam PageList tertanam.
--
-- Menggabungkannya jadi satu kolom teks berarti MENGURAI TEKS setiap kali
-- ada yang perlu ditunjuk balik — dan penguraian teks adalah kelas cacat
-- yang justru sedang kita tinggalkan. Lima kolom dapat di-JOIN langsung.
--
-- KENAPA ISTILAH _Avoid_ MUNCUL DI SINI, DAN HANYA DI SINI
-- CASEID dan pyID ada di daftar _Avoid_ glosarium. Keduanya dipakai di
-- berkas ini sebagai NAMA LAMA, bersufiks _LAMA, di tabel yang seluruh
-- tugasnya memang menyimpan nama lama — pengecualian yang sama dengan
-- lapisan view kompatibilitas (ADR-0021). Tidak satu pun dari kelimanya
-- muncul di tabel kanonik mana pun, dan tabel ini sendiri dijatuhkan
-- setelah paritas diterima.


-- ---------------------------------------------------------------------
-- Arah relasi: polimorfik, dan karena itu tanpa foreign key
-- ---------------------------------------------------------------------
-- TABEL_TUJUAN + ID_TUJUAN menunjuk baris di SALAH SATU dari 22 tabel
-- kanonik. Oracle tidak punya foreign key polimorfik, dan memalsukannya
-- dengan 22 kolom nullable lebih buruk daripada tidak punya sama sekali.
--
-- Keutuhan rujukan karena itu JATUH KE APLIKASI, dan itu dinyatakan di
-- sini, bukan didiamkan. Yang menggantikannya: uji paritas tiket 34
-- memeriksa setiap baris korelasi menunjuk baris yang benar-benar ada.
--
-- TABEL_TUJUAN VARCHAR2(30 CHAR) bukan angka sembarang — itu batas
-- panjang nama objek yang berlaku di skema ini (SPEC bagian 16), jadi
-- kolom ini tidak akan pernah kependekan bagi nama tabel mana pun.


CREATE TABLE KLAIMNP.MIGRASI_KORELASI
(
  ID_KORELASI         NUMBER(19) NOT NULL,
  TABEL_TUJUAN        VARCHAR2(30 CHAR) NOT NULL,
  ID_TUJUAN           NUMBER(19) NOT NULL,
  PZINSKEY_LAMA       VARCHAR2(255 CHAR),
  PYID_LAMA           VARCHAR2(32 CHAR),
  CASEID_LAMA         VARCHAR2(255 CHAR),
  INDEX_OBJECT_LAMA   NUMBER(9),
  URUTAN_BARIS_LAMA   NUMBER(9),
  SUMBER_LAMA         VARCHAR2(64 CHAR) NOT NULL,
  DIMIGRASI_PADA      TIMESTAMP(6) NOT NULL
);

ALTER TABLE KLAIMNP.MIGRASI_KORELASI ADD CONSTRAINT PK_MIGRASI_KORELASI PRIMARY KEY (ID_KORELASI);

-- Sebuah baris jembatan yang tidak membawa SATU PUN pengenal lama tidak
-- menjembatani apa pun. Kelimanya boleh kosong sendiri-sendiri; kosong
-- SEMUA tidak boleh. Ini satu-satunya keutuhan yang dapat ditegakkan
-- basis data di tabel ini, dan karena itu ia dipasang.
ALTER TABLE KLAIMNP.MIGRASI_KORELASI ADD CONSTRAINT CK_MIGRASI_KORELASI_1 CHECK (
  PZINSKEY_LAMA IS NOT NULL
  OR PYID_LAMA IS NOT NULL
  OR CASEID_LAMA IS NOT NULL
  OR INDEX_OBJECT_LAMA IS NOT NULL
  OR URUTAN_BARIS_LAMA IS NOT NULL
);

-- Dua arah pencarian, dua index. Yang pertama sudah ada sejak awal; yang
-- kedua ditambahkan 19 Sep 2026 karena ARAH ITULAH yang tiket ini janjikan.
CREATE INDEX KLAIMNP.IX_MIGRASI_KORELASI_1 ON KLAIMNP.MIGRASI_KORELASI (TABEL_TUJUAN, ID_TUJUAN);
CREATE INDEX KLAIMNP.IX_MIGRASI_KORELASI_2 ON KLAIMNP.MIGRASI_KORELASI (PZINSKEY_LAMA);
CREATE INDEX KLAIMNP.IX_MIGRASI_KORELASI_3 ON KLAIMNP.MIGRASI_KORELASI (PYID_LAMA);

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

CREATE SEQUENCE KLAIMNP.SQ_MIGRASI_KORELASI START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
GRANT SELECT ON KLAIMNP.SQ_MIGRASI_KORELASI TO KLAIMNP_APP;

-- ADR-0017/0028: satu pintu tulis, ditegakkan lewat grant.
GRANT SELECT, INSERT, UPDATE, DELETE ON KLAIMNP.MIGRASI_KORELASI TO KLAIMNP_APP;


-- ---------------------------------------------------------------------
-- Apa yang TIDAK ditegakkan di sini
-- ---------------------------------------------------------------------
--   1. UNIQUE atas kombinasi pengenal lama. SENGAJA BELUM ADA, dan ini
--      bukan kelalaian: KOLOM mana yang ada sudah EVIDENCED; KOMBINASI
--      mana yang unik menunggu REQ-018. Memasang UNIQUE sekarang berarti
--      menebak, dan tebakan yang salah menggugurkan migrasi di tengah
--      jalan — tepat ketika ia paling mahal dibatalkan.
--      Yang berubah bila REQ-018 menjawab lain: CONSTRAINT, bukan kolom.
--   2. Keutuhan rujukan TABEL_TUJUAN/ID_TUJUAN — polimorfik, lihat atas.
--   3. Umur baris. Tabel ini jembatan; ia dijatuhkan setelah paritas
--      diterima (ADR-0005), bersama MIGRASI_PENDARATAN.
--
-- YANG DITEGAKKAN: setiap baris menyebut tabel dan baris tujuannya,
-- menyebut sumbernya, menyebut kapan ia dimigrasi, dan membawa
-- SEKURANGNYA SATU pengenal lama.
-- =====================================================================

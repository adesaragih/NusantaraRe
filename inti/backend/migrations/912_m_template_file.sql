-- 912 - M_TEMPLATE_FILE: berkas templat unduhan berversi milik menu Template Manager (keputusan work owner
-- 04-10-2026: "aku mau ada menu untuk pengelola file templete unduhan ... bukan cuman untuk bordereaux, tapi untuk
-- semua menu yang memiliki template"). Padanan Pega: Rule-File-Binary yang dicari menurut nama lalu diambil yang
-- terbaru (Bordereaux `DownloadTemplate_Act` langkah 3).
--
-- KODE = kode slot yang didaftarkan modul (`inti.Pendaftaran.Templat`, mis. `bordereaux.premi.fire`). Satu KODE
-- banyak VERSI; paling banyak SATU baris AKTIF = '1' per KODE (dijaga layanan dalam satu transaksi). Tanpa baris
-- aktif, aplikasi memakai berkas bawaan modulnya - jadi tabel kosong tidak mengubah perilaku apa pun.
-- ISI disimpan utuh sebagai BLOB (templat terbesar sekitar 14 KB, batas unggah 1 MB).
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
CREATE TABLE {skema}.M_TEMPLATE_FILE (
  ID            NUMBER(19) NOT NULL,
  KODE          VARCHAR2(100) NOT NULL,
  VERSI         NUMBER(10) NOT NULL,
  NAMA_BERKAS   VARCHAR2(255) NOT NULL,
  UKURAN        NUMBER(10) NOT NULL,
  JUMLAH_KOLOM  NUMBER(5),
  ISI           BLOB NOT NULL,
  CATATAN       VARCHAR2(2000),
  AKTIF         VARCHAR2(1) DEFAULT '0' NOT NULL,
  DIUNGGAH_OLEH VARCHAR2(64) NOT NULL,
  TGL_UNGGAH    DATE DEFAULT SYSDATE NOT NULL,
  CONSTRAINT PK_M_TEMPLATE_FILE PRIMARY KEY (ID),
  CONSTRAINT UX_M_TEMPLATE_FILE_VERSI UNIQUE (KODE, VERSI),
  CONSTRAINT CK_M_TEMPLATE_FILE_AKTIF CHECK (AKTIF IN ('0', '1'))
)
/
CREATE SEQUENCE {skema}.SEQ_M_TEMPLATE_FILE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

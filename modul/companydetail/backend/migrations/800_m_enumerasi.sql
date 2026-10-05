-- 800 - M_ENUMERASI: tabel datar pilihan dropdown Company Detail (perintah work owner 04-10-2026: "jangan ada pake
-- datapega, kalo mau kamu buat table baru ke pooldata"; "oke gas").
--
-- Satu baris = satu pilihan: JENIS (nama daftar, sama dengan tipe enumerasi Pega asalnya) + KODE (nilai yang
-- disimpan di tabel client) + LABEL (teks layar). AKTIF 1 = tampil di dropdown; 0 = hanya untuk menampilkan nilai
-- lama. Isinya migrasi 801. COUNTRY tidak di sini - sumbernya tabel NATION (802-804).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
CREATE TABLE {skema}.M_ENUMERASI (
  JENIS  VARCHAR2(30) NOT NULL,
  KODE   VARCHAR2(20) NOT NULL,
  LABEL  VARCHAR2(200),
  AKTIF  CHAR(1) DEFAULT '1' NOT NULL,
  URUTAN NUMBER(5) NOT NULL,
  CONSTRAINT PK_M_ENUMERASI PRIMARY KEY (JENIS, KODE),
  CONSTRAINT CK_M_ENUMERASI_AKTIF CHECK (AKTIF IN ('0', '1'))
)
/

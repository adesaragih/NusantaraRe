-- M_TREATYIN_COINSCALE - tabel pendaratan kesembilan, tab Co-Ins Scale
--
-- Migrasi tiket tab Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Mengikuti pola kedelapan tabel migrasi `430` TANPA satu pun kekecualian:
-- `MASTERID` + `URUTAN`, `UNIQUE (MASTERID, URUTAN)`, pengenal dari sequence
-- (migrasi `433`), nilai disimpan sebagai TEKS apa adanya.
--
-- =====================================================================
-- UKURAN - SELURUH 1.854 DOKUMEN
-- =====================================================================
--   `CoInScale` berisi di 186 kontrak, 702 elemen, terbanyak 5 per kontrak.
--
--   ⚠️ Rancangan ronde ini menyebut *"dua medan saja - `CoInShare` dan
--   `PctLimit`"*. Sapuan menemukan ENAM medan skalar: keempat medan jejak
--   Pega (`pxCreateDateTime`, `pxCreateOpName`, `pxCreateOperator`,
--   `pxCreateSystemID`) ada pada 701 dari 702 elemen, dan `pxObjClass` pada
--   seluruhnya. Keempatnya dibawa, sebab ia satu-satunya catatan SIAPA yang
--   menyusun skala itu dan KAPAN - dan tab Co-Ins Scale adalah tab yang
--   angkanya dinegosiasikan, jadi pertanyaan itu akan ditanyakan.
--
--   Panjang terukur: `CoInShare` 17 (`>=30% up to < 50%`), `PctLimit` 3,
--   `pxCreateOpName` 24, `pxCreateOperator` 17, `pxCreateDateTime` 23,
--   `pxCreateSystemID` 13, `pxObjClass` 35.
--
-- ⛔ `COINSHARE` TEKS, dan ini bukan kemalasan. Nilainya BUKAN angka:
-- `>=30% up to < 50%`, `>=25%`. Ia pita, bukan bilangan, dan kolom angka
-- akan menolak seluruh 702 barisnya.
--
-- ⛔ NOL kunci asing ke `TREATY_IN` - `POOLDATA.TREATY_IN` tidak punya kunci
-- utama maupun UNIQUE pada `ID` (kolomnya bahkan NULLABLE), dan Oracle
-- menolak merujuknya dengan ORA-02270. Keputusan beserta syarat
-- pembalikannya sudah berdiri di `docs/KEPUTUSAN-PENYELARASAN-REPO.md` §12;
-- berkas ini MENGIKUTINYA, tidak mengulang penyelidikannya.
--
-- INV-01 kunci utama. INV-02 pengenal dari sequence.
-- Nama objek di bawah 30 bita: terpanjang `UQ_MTI_COINSCALE` 16 bita.
CREATE TABLE {skema}.M_TREATYIN_COINSCALE (
  ID                NUMBER(19)           NOT NULL,
  MASTERID          VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN            NUMBER(10)           NOT NULL,
  COINSHARE         VARCHAR2(200 CHAR),
  PCTLIMIT          VARCHAR2(50 CHAR),
  PXCREATEDATETIME  VARCHAR2(50 CHAR),
  PXCREATEOPNAME    VARCHAR2(200 CHAR),
  PXCREATEOPERATOR  VARCHAR2(200 CHAR),
  PXCREATESYSTEMID  VARCHAR2(100 CHAR),
  PXOBJCLASS        VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_COINSCALE PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_COINSCALE UNIQUE (MASTERID, URUTAN)
)
/

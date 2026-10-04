-- 761 - Master Data: status aktif / nonaktif master yang tabelnya WARISAN tanpa kolom status (NATION).
--
-- M-3 menuntut aktif / nonaktif tanpa hapus. Tabel flat milik modul ini (760) membawa STS_AKTIF sendiri;
-- OBJECTITEMTYPE punya ISACTIVE. NATION adalah tabel warisan (bab "Tabel warisan" MODUL.md) - mengubah strukturnya
-- (ALTER ADD) berarti mengambil kepemilikannya, yang bukan keputusan modul ini (`TestTabelBukanMilikKitaTidakDibuat`).
-- Karena itu statusnya disimpan DI SINI (MD-2): satu baris per (tabel, ID) yang statusnya pernah diubah; TIDAK ADA
-- baris = aktif. ID_BARIS VARCHAR2(100) cukup untuk NATION.ID VARCHAR2(10). Ditulis, TIDAK dijalankan agent. Nol COMMIT.
CREATE TABLE {skema}.T_MASTER_STATUS (
  NAMA_TABEL VARCHAR2(30) NOT NULL,
  ID_BARIS   VARCHAR2(100) NOT NULL,
  STS_AKTIF  VARCHAR2(1) DEFAULT '1' NOT NULL,
  CONSTRAINT PK_T_MASTER_STATUS PRIMARY KEY (NAMA_TABEL, ID_BARIS)
)
/

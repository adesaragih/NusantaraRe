-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`-`453`: rentang `treatyin` (`400-439`) PENUH.
--
-- ---------------------------------------------------------------------
-- Tempat simpan pecahan spreading MANUAL tab Share PROPORSIONAL
-- ---------------------------------------------------------------------
--
-- Laporan Save 8 Oktober 2026 — "Belum punya kolom di tabel, jadi TIDAK
-- tersimpan: BreakDownSprdList". Pemilik proses: "kerjakan sekarang".
--
--   T_TREATY_LIMIT_SPRD_BREAKDOWN (BARU) — `.Detail.SpreadingList(n)
--     .BreakDownSprdList`, panel bentang baris grid spreading manual
--     (`Section/SpreadingTPDtl.xml`), diisi `SetSpreadName` [3.2.8.1.1]:
--
--       REINSNAME  .ReinsName  nama anak susunan (`QS (OR)`, `QS (R/I)`, `ORS`)
--       REINSID    .ReinsID    `ReinsTypeID` anak susunan
--       CURRENCY   .Currency   mata uang `RNMShareList`
--       AMOUNT     .Amount     Value × (Pct/RNMShare) × Pct anak / 100
--       SHAREPCT   .SharePct   Pct anak susunan
--
--   Tingkat KEEMPAT pohon: T_TREATY_LIMITS → T_TREATY_LIMIT_DETAIL →
--   T_TREATY_LIMIT_SPREADING → T_TREATY_LIMIT_SPRD_BREAKDOWN (`IDINDUK` =
--   `T_TREATY_LIMIT_SPREADING.ID`). Nama tabel dipendekkan (`SPRD`) supaya
--   tetap ≤ 30 aksara.
--
-- ⚠️ TANPA kunci asing ke `T_TREATY_LIMIT_SPREADING` — pola `448`/`449`:
-- penulis menghapus baris anak per `MASTERID` (`KosongkanKontrak`) sebelum
-- menyisip ulang.
--
-- ⚠️ Penulis dan pembaca `T_TREATY_*` TOLERAN terhadap tabel yang belum
-- terpasang (`kolomTerpasang`): sebelum berkas ini dijalankan Save tetap
-- berjalan dan melaporkan `BreakDownSprdList` sebagai tidak tersimpan.

-- T_TREATY_LIMIT_SPRD_BREAKDOWN  <- TreatyIn.Limits.Detail.SpreadingList.BreakDownSprdList
CREATE TABLE {skema}.T_TREATY_LIMIT_SPRD_BREAKDOWN (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  AMOUNT                           VARCHAR2(4000 CHAR),
  CURRENCY                         VARCHAR2(4000 CHAR),
  REINSID                          VARCHAR2(4000 CHAR),
  REINSNAME                        VARCHAR2(4000 CHAR),
  SHAREPCT                         VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_SPRD_BRKDN PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_SPRD_BRKDN UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_SPRD_BRKDN_MST ON {skema}.T_TREATY_LIMIT_SPRD_BREAKDOWN (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_SPRD_BRKDN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

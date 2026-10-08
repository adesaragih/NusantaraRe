-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`/`445`/`446`/`448`/`449`/`450`: rentang `treatyin` (`400-439`) PENUH.
--
-- ---------------------------------------------------------------------
-- Tempat simpan larik baris tab Share Non-Prop yang selama ini TANPA tabel
-- ---------------------------------------------------------------------
--
-- Laporan Save kontrak 1001855 (8 Oktober 2026) — "Belum punya kolom di
-- tabel, jadi TIDAK tersimpan: DeductionTotalList, GrossPremiumMinList,
-- RNMSpreadedListDeductRIXOL, RNMSpreadedListDeductXOL,
-- RNMSpreadedListGrossMinXOL, RNMSpreadedListGrossRIMinXOL,
-- RNMSpreadedListGrossRIXOL, RNMSpreadedListGrossXOL,
-- RNMSpreadedListNetRIXOL, RNMSpreadedListNetXOL, RNMSpreadedListRIXOL,
-- RNMSpreadedListXOL, RnmLimitList, TreatyGroupList".
--
--   T_TREATY_SHARE_GROUP (BARU)         Share.TreatyGroupList
--     Treaty Group baris Share — disalin dari layer (`TreatyInNonAddItem`)
--   T_TREATY_SHARE_XOL_AMOUNT (BARU)    Share.RNMSpreadedList*XOL (10 larik)
--     bagian OR / R/I tiap baris (`FetchQSfromMasterXOL`,
--     `TreatyInSetBrokerage`) — `{Currency, Value}`, dibedakan `JENIS`
--   T_TREATY_FAC_SHARE_GROUP (BARU)     FacultativeShareList.TreatyGroupList
--
-- ⚠️ `JENIS` tabel XOL `VARCHAR2(40 CHAR)`, bukan 20: nama larik terpanjang
-- `RNMSpreadedListGrossRIMinXOL` 28 aksara. Kesepuluhnya TIDAK ditaruh di
-- `T_TREATY_SHARE_AMOUNT` (`JENIS` 20) — penulis memeriksa tabel dan kolom
-- yang terpasang, tidak memeriksa LEBARnya; nama yang kepanjangan akan
-- menggagalkan SELURUH Save sampai migrasi pelebarnya terpasang.
--
-- ⚠️ Tiga larik lain TIDAK butuh DDL — bentuknya `{Currency, Value}` dan
-- namanya muat di `JENIS` 20 (terpanjang `GrossPremiumMinList`, 19):
--   T_TREATY_SHARE_AMOUNT      RnmLimitList · GrossPremiumMinList ·
--                              DeductionTotalList
--   T_TREATY_FAC_SHARE_AMOUNT  idem
--
-- ⚠️ TANPA kunci asing ke tabel induk — pola `448`/`449`/`450`: penulis
-- menghapus baris anak per `MASTERID` (`KosongkanKontrak`) sebelum
-- menyisip ulang.
-- ⚠️ Penulis dan pembaca TOLERAN terhadap tabel yang belum terpasang.

-- T_TREATY_SHARE_GROUP  <- TreatyIn.Share.TreatyGroupList
CREATE TABLE {skema}.T_TREATY_SHARE_GROUP (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  TREATYGROUP                      VARCHAR2(4000 CHAR),
  TREATYGROUPID                    VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_SHARE_GROUP PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_SHARE_GROUP UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_SHARE_GROUP_MST ON {skema}.T_TREATY_SHARE_GROUP (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_SHARE_GROUP START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_SHARE_XOL_AMOUNT  <- TreatyIn.Share.RNMSpreadedList*XOL
CREATE TABLE {skema}.T_TREATY_SHARE_XOL_AMOUNT (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  JENIS     VARCHAR2(40 CHAR)  NOT NULL,
  CURRENCY                         VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_SHARE_XOL_AMOUNT PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_SHARE_XOL_AMOUNT UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_SHARE_XOL_AMOUNT_MST ON {skema}.T_TREATY_SHARE_XOL_AMOUNT (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_SHARE_XOL_AMOUNT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_FAC_SHARE_GROUP  <- TreatyIn.FacultativeShareList.TreatyGroupList
CREATE TABLE {skema}.T_TREATY_FAC_SHARE_GROUP (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  TREATYGROUP                      VARCHAR2(4000 CHAR),
  TREATYGROUPID                    VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_FAC_SHARE_GROUP PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_FAC_SHARE_GROUP UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_FAC_SHARE_GROUP_MST ON {skema}.T_TREATY_FAC_SHARE_GROUP (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_FAC_SHARE_GROUP START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`/`445`/`446`/`448`/`449`: rentang `treatyin` (`400-439`) PENUH.
--
-- ---------------------------------------------------------------------
-- Tempat simpan larik rincian tab Limits yang selama ini TANPA tabel
-- ---------------------------------------------------------------------
--
-- Laporan Save kontrak 1002305 (8 Oktober 2026) — "Belum punya kolom di
-- tabel, jadi TIDAK tersimpan: CashLossList, ClaimCoopList, CurrencyList,
-- DeductionList, DeductionTotalList, MDPMinList, PLAList,
-- Reinstatement_List, ReserveList". Pembaca pohon Limits menyisipkan
-- kesembilannya sebagai larik KOSONG (`larikTanpaTabelDetail/Limit`), dan
-- isian pemakai di grid-grid itu memang hilang saat Save.
--
--   T_TREATY_LIMIT_DEDUCTION (BARU)       Limits.Detail.DeductionList
--     Description · Currency · Deduction · Deduction % (`CalculateDeduction`)
--   T_TREATY_LIMIT_ACH_PARAM (BARU)       Limits.Detail.CurrencyList
--     grid Parameter · Based on Gross / Nett tab Achievement (`GetAchievement`)
--   T_TREATY_LIMIT_REINSTATEMENT (BARU)   Limits.Reinstatement_List (Non-Prop)
--     `TreatySetReinstatement` / `SetReinstatementPct` / `CalculateReinstatement*`
--
-- ⚠️ Enam larik lain TIDAK butuh DDL — bentuknya `{Currency, Value}`, cukup
-- peta pendaratan (kolom `JENIS` VARCHAR2(20), nama terpanjang 18 aksara):
--   T_TREATY_LIMIT_AMOUNT   ReserveList · PLAList · CashLossList ·
--                           ClaimCoopList · DeductionTotalList
--   T_TREATY_LIMIT_MEASURE  MDPMinList
--
-- ⚠️ TANPA kunci asing ke tabel induk — pola `448`/`449`: penulis menghapus
-- baris anak per `MASTERID` (`KosongkanKontrak`) sebelum menyisip ulang.
-- ⚠️ `ID` baris Reinstatement mendarat di `IDX` — `ID` kunci utama baris
-- pendaratan, sama seperti `T_TREATY_LIMIT_DETAIL`.
-- ⚠️ Penulis dan pembaca TOLERAN terhadap tabel yang belum terpasang.

-- T_TREATY_LIMIT_DEDUCTION  <- TreatyIn.Limits.Detail.DeductionList
CREATE TABLE {skema}.T_TREATY_LIMIT_DEDUCTION (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  COMMENT_                         VARCHAR2(4000 CHAR),
  CURRENCY                         VARCHAR2(4000 CHAR),
  CURRENCYID                       VARCHAR2(4000 CHAR),
  DEDUCTION                        VARCHAR2(4000 CHAR),
  DEDUCTIONPCT                     VARCHAR2(4000 CHAR),
  DEDUCTIONPCTCALCULATE            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_DEDUCTION PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_DEDUCTION UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_DEDUCTION_MST ON {skema}.T_TREATY_LIMIT_DEDUCTION (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_DEDUCTION START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_LIMIT_ACH_PARAM  <- TreatyIn.Limits.Detail.CurrencyList
CREATE TABLE {skema}.T_TREATY_LIMIT_ACH_PARAM (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  PARAMETER                        VARCHAR2(4000 CHAR),
  ACHIEVEMENTPCTGROSS              VARCHAR2(4000 CHAR),
  LOSSRATIOGROSS                   VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_ACH_PARAM PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_ACH_PARAM UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_ACH_PARAM_MST ON {skema}.T_TREATY_LIMIT_ACH_PARAM (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_ACH_PARAM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_LIMIT_REINSTATEMENT  <- TreatyIn.Limits.Reinstatement_List
CREATE TABLE {skema}.T_TREATY_LIMIT_REINSTATEMENT (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  ADDITIONALAMOUNT1                VARCHAR2(4000 CHAR),
  ADDITIONALAMOUNT2                VARCHAR2(4000 CHAR),
  ADDITIONALPCT                    VARCHAR2(4000 CHAR),
  IDX                              VARCHAR2(4000 CHAR),
  REINSTATEMENTAMOUNT1             VARCHAR2(4000 CHAR),
  REINSTATEMENTAMOUNT2             VARCHAR2(4000 CHAR),
  REINSTATEMENTNOTE                VARCHAR2(4000 CHAR),
  REINSTATEMENTPCT                 VARCHAR2(4000 CHAR),
  REINSTATEMENTVALUE               VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_REINST PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_REINST UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_REINST_MST ON {skema}.T_TREATY_LIMIT_REINSTATEMENT (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_REINST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

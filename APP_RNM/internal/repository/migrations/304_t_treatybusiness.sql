-- T_TREATYBUSINESS - jenis bisnis pada kombinasi tahun/grup/jenis (tiket 01).
--
-- Sumber warisan POOLDATA.TREATYBUSINESS (`[data DBA]` seluruh VARCHAR2),
-- ditulis PEGA_TREATYBUSINESS (RDBList/SaveMasterTreatyBusiness_SQL.xml,
-- 12 parameter + 2 keluaran). Dibaca hilir: Claim Prop GetTreatyGroupID.xml
-- (TREATYGROUPID WHERE BIZCODE, TREATYYEAR, ISACTIVE='1'), Claim Fac In
-- GetTreatyGroup_Sql.xml (ID, TREATYYEAR, TREATYGROUPID, REINSTYPEID) dan
-- GetDataTreatyLimit_Sql.xml (TREATYYEARID, REINSTYPEID, BIZCODE).
--
-- ⛔ TREATYYEARID boleh NULL di data lama: kaskade warisan
-- (DeleteFromTREATYCONTRACT_SQL.xml) menulis `OR TREATYYEARID IS NULL`.
-- Karena itu TANPA FK ke T_TREATYYEAR - FK akan menolak baris lama itu.
--
-- ⛔ Menggantung pada KUNCI GABUNGAN, bukan FK ke kontrak (spec §2).
-- TGLUPDATE dari VARCHAR2 menjadi DATE (AC 53).
--
-- ⛔ NOL COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_TREATYBUSINESS (
  ID              VARCHAR2(32) NOT NULL,
  ISACTIVE        VARCHAR2(255),
  TREATYYEAR      VARCHAR2(255),
  TREATYYEARID    VARCHAR2(32),
  TREATYGROUPID   VARCHAR2(32),
  TREATYGROUPNAME VARCHAR2(255),
  REINSTYPEID     VARCHAR2(32),
  REINSTYPENAME   VARCHAR2(255),
  BIZCODE         VARCHAR2(255),
  BIZNAME         VARCHAR2(255),
  USERID          VARCHAR2(255),
  TGLUPDATE       DATE,
  CONSTRAINT PK_T_TREATYBUSINESS PRIMARY KEY (ID)
)
/

CREATE INDEX {skema}.IDX_T_TREATYBUSINESS_KOMB
  ON {skema}.T_TREATYBUSINESS (TREATYYEAR, TREATYGROUPID, REINSTYPEID)
/

-- Identitas '1' + lpad(6): bentuk warisan TREATY_BUSINESS_SEQ (spec §7).
CREATE SEQUENCE {skema}.SEQ_T_TREATYBUSINESS START WITH 1 INCREMENT BY 1 NOCACHE
/

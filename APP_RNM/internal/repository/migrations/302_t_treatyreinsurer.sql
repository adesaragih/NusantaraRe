-- T_TREATYREINSURER - reinsurer pada kombinasi tahun/grup/jenis (tiket 01).
--
-- Sumber warisan POOLDATA.TREATYREINSURER, ditulis PEGA_TREATYREINSURER
-- (RDBList/SaveMasterTreatyReinsurer_SQL.xml, 19 parameter + 2 keluaran).
-- Dibaca hilir Claim Prop, Komite Claim Prop, Claim Fac In lewat
-- GetListRetro_Sql: ID, REINSURERID, CLIENTID, NAME, RICOMM, PCTSHARE
-- disaring REINSTYPEID + TREATYYEAR + TREATYGROUPID - nama-nama itu VERBATIM
-- (kontrak tco3, tco_kontrak_hilir.go).
--
-- ⛔ Menggantung pada KUNCI GABUNGAN (TREATYYEAR, TREATYGROUPID, REINSTYPEID),
-- BUKAN FK ke T_TREATYCONTRACT. Menambahkan FK ke kontrak mengubah arti data
-- (spec §2). Index gabungan di bawah menggantikan FK sebagai jalur baca.
--
-- RICOMM dan PCTSHARE sudah NUMBER di warisan; STARTDATE, ENDDATE, TGLUPDATE
-- dari VARCHAR2 menjadi DATE (AC 53). IUDATE tetap teks: bentuk dan artinya
-- `[terbuka]`, tidak ada di daftar AC 53 - dibawa apa adanya, tidak ditebak.
--
-- ⛔ NOL COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_TREATYREINSURER (
  ID              VARCHAR2(32) NOT NULL,
  TREATYYEAR      VARCHAR2(255),
  TREATYGROUPID   VARCHAR2(32),
  TREATYGROUPNAME VARCHAR2(255),
  REINSTYPEID     VARCHAR2(32),
  REINSTYPENAME   VARCHAR2(255),
  REINSURERID     VARCHAR2(32),
  CLIENTID        VARCHAR2(32),
  NAME            VARCHAR2(255),
  RICOMM          NUMBER(38,8),
  PCTSHARE        NUMBER(38,8),
  IUDATE          VARCHAR2(255),
  USERID          VARCHAR2(255),
  STARTDATE       DATE,
  ENDDATE         DATE,
  STATUSON        VARCHAR2(255),
  STDRATING       VARCHAR2(255),
  OPERATORNAME    VARCHAR2(255),
  TGLUPDATE       DATE,
  CONSTRAINT PK_T_TREATYREINSURER PRIMARY KEY (ID)
)
/

CREATE INDEX {skema}.IDX_T_TREATYREINSURER_KOMB
  ON {skema}.T_TREATYREINSURER (TREATYYEAR, TREATYGROUPID, REINSTYPEID)
/

-- Identitas '1' + lpad(6): bentuk warisan M_TREATYREINSURER_SEQ (spec §7).
CREATE SEQUENCE {skema}.SEQ_T_TREATYREINSURER START WITH 1 INCREMENT BY 1 NOCACHE
/

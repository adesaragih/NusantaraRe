-- 183 - T_QUOTATIONDATA SEBAGIAN + SEQ_T_QUOTATIONDATA (tiket 31), butir 78.
--
-- Sama dengan 182: hanya kolom sistem + kolom yang dipakai blok General; sisanya (termasuk EDM_CHARGE_FEE /
-- SHARE_OF_CEDING yang panjang VARCHAR2-nya menunggu tiket 23, butir 68.1) ditambah tiket 23 lewat ALTER.
-- Nama/tipe = rancangan. Satu baris per case (halaman tunggal `QuotationData` di bawah T_GENERAL_POLIS,
-- jalur rancangan) - ditegakkan UQ_T_QUOTATIONDATA_PARENT (keputusan agent A87).
--
-- ID = surrogate dari SEQ_T_QUOTATIONDATA, NUMBER(19) (pola identitas-dari-sequence yang penjaga izinkan;
-- bukan kolom uang). PARENT_ID VARCHAR2(32) = T_GENERAL_POLIS.ID (butir 76.1), FK tanpa ON DELETE.
CREATE TABLE {skema}.T_QUOTATIONDATA (
  ID               NUMBER(19) NOT NULL,
  IDPEGA           VARCHAR2(50),
  COB_GROUP        VARCHAR2(20),
  PARENT_ID        VARCHAR2(32) NOT NULL,
  NO_OFFER_SLIP    VARCHAR2(50),
  QQ_NAME          VARCHAR2(500),
  POLICY_TYPE      VARCHAR2(50),
  MOID             VARCHAR2(50),
  EDM_DAY          VARCHAR2(50),
  TYPE_FACULTATIVE VARCHAR2(50),
  SOB_NAME         VARCHAR2(500),
  CEDING_CO_NAME   VARCHAR2(500),
  GROUP_NAME       VARCHAR2(500),
  CONSTRAINT PK_T_QUOTATIONDATA PRIMARY KEY (ID),
  CONSTRAINT UQ_T_QUOTATIONDATA_PARENT UNIQUE (PARENT_ID),
  CONSTRAINT FK_QUOTATIONDATA_GENERAL FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID)
)
/
CREATE SEQUENCE {skema}.SEQ_T_QUOTATIONDATA START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

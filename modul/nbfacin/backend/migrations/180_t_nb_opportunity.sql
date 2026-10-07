-- 180 - T_NB_OPPORTUNITY: isian form Opportunity NB (tiket 29), butir 76.3.
--
-- Keputusan work owner 03-10-2026 (AskUserQuestion): tabel opportunity SENDIRI, sebab
-- di Pega opportunity adalah kelas work tersendiri `[terverifikasi]`
-- (`NB FacIn\ReportDefinition\GetListOpportunity.xml`: kelas
-- ASM-FW-SFAGISFW-Work-Opportunity, INNER JOIN ASM-FW-GISFW-Work-NB pada
-- A.pzInsKey = .NBHandle) dan rancangan tabel flat tidak punya tabel untuknya.
--
-- ⛔ BERBAGI PK DENGAN T_WORK_POLIS (ID = ID, 1:1), TANPA constraint FK - pola
-- T_PREMIUM_LIST premiumlistlife (050/051). T_WORK_POLIS milik premiumlistlife (K-064);
-- constraint lintas modul tidak diputuskan siapa pun.
--
-- Tipe: kolom yang disalin dari T_M_ACCOUNT memakai tipe DDL sumbernya (CHAR);
-- CLASS_OF_BUSINESS = BUSINESS.NOTE (4000 BYTE); teks layar lain VARCHAR2(255),
-- DESCRIPTION VARCHAR2(4000) - keputusan agent A81, panjang diperiksa di Go (400).
-- Nol kolom NUMBER: tidak tersangkut keputusan presisi tim inti. Seluruh kolom
-- nullable kecuali PK; wajib-isi ditegakkan di Go. Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_NB_OPPORTUNITY (
  ID                     VARCHAR2(32) NOT NULL,
  ESTIMATED_CLOSING_DATE DATE,
  BUSINESS_PROSPECT_NAME VARCHAR2(255),
  ACCOUNT_ID             VARCHAR2(255 CHAR),
  INSURED_ID             VARCHAR2(255 CHAR),
  GROUP_BUSINESS_ID      VARCHAR2(32 CHAR),
  GROUP_BUSINESS         VARCHAR2(64 CHAR),
  CLASS_OF_BUSINESS      VARCHAR2(4000 BYTE),
  TYPE_OF_INWARD         VARCHAR2(255),
  TYPE_OF_FACULTATIVE    VARCHAR2(255),
  PHASE                  VARCHAR2(255),
  STAGE                  VARCHAR2(255),
  OPPORTUNITY_SOURCE     VARCHAR2(255),
  BUSINESS_STATUS        VARCHAR2(255),
  DESCRIPTION            VARCHAR2(4000),
  CONSTRAINT PK_T_NB_OPPORTUNITY PRIMARY KEY (ID)
)
/

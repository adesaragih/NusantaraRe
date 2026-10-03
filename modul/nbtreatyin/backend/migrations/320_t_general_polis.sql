-- 320 - T_GENERAL_POLIS: satu baris per GENERASI polis treaty inward (tiket 16).
--
-- Kunci utama BERSAMA T_WORK_POLIS (tabel kasus lintas-lini milik premiumlistlife,
-- migrasi 050/059): ID adalah ID baris T_WORK_POLIS, tanpa kolom kunci tamu
-- tersendiri (spec-penyimpanan ID-7, AC 5).
--
-- Diagram grilling F12-F16: UNIQUE (NOPOLIS, PRODKE) dan UNIQUE (OLD_POLIS_ID).
-- Kunci alami (NOPOLIS, PRODKE) unik - ditegakkan BASIS DATA (ID-8, AC 1).
-- NOPOLIS kosong selama realisasi belum bernomor, jadi indeks uniknya hanya
-- memuat baris bernomor (indeks berfungsi CASE): dua draf tanpa nomor tidak bentrok
-- (UNIQUE biasa menganggap (NULL, 0) dan (NULL, 0) kembar di Oracle).
-- PRODKE bilangan bulat lebar (KEPUTUSAN-RONDE-12 butir 1 dan 8); NB selalu 0.
-- OLD_POLIS_ID menunjuk generasi sebelumnya - unik, kosong di NB (ID-9, AC 3, 4).
-- Generasi TERTUTUP = ada baris penerus yang OLD_POLIS_ID-nya menunjuk generasi
-- ini; tidak boleh disunting (ID-10, AC 6) - tanpa kolom penanda (diagram).
--
-- Kolom lain DIBANGKITKAN dari backend/models/katalog.go (docs/alat/skema.py) -
-- 79 medan PolicyTreatyIn + 7 kolom json_polis menurut diagram dan rancangan;
-- LAYER* dicoret (diagram F26). Perbandingan: docs/PERBANDINGAN-KOLOM-DIAGRAM.md.
-- Uang dan persen NUMBER(38,8) (KEPUTUSAN 23-09-2026 sore), tanggal DATE (P32),
-- kode dan penanda teks (ID-16, ID-17). Nol COMMIT.
CREATE TABLE {skema}.T_GENERAL_POLIS (
  ID                        VARCHAR2(32) NOT NULL,
  NOPOLIS                   VARCHAR2(64),
  PRODKE                    NUMBER(10) DEFAULT 0 NOT NULL,
  NOENDORS                  VARCHAR2(64),
  OLD_POLIS_ID              VARCHAR2(32),
  IDPEGA                    VARCHAR2(128),
  TGL_INPUT                 DATE,
  USERNAME                  VARCHAR2(64),
  POSITION_NOTE             VARCHAR2(64),
  NB_STATUS                 VARCHAR2(255),
  TREATY_IN_ID              VARCHAR2(64),
  NO_OFFER                  VARCHAR2(64),
  MASTER_ID                 VARCHAR2(64),
  IS_APPROVED               VARCHAR2(16),
  SUGGEST                   VARCHAR2(4000),
  SUGGEST_DATE              DATE,
  OPERATOR_NAME             VARCHAR2(128),
  IS_NEW_POLICY_NON_PROP    VARCHAR2(16),
  IS_EDM_INPUT_ON_NB        VARCHAR2(16),
  HAS_FAC_OUT               VARCHAR2(16),
  FLAG_PPH                  VARCHAR2(16),
  FLAG_RETRO_TREATY         VARCHAR2(16),
  DUE_TO                    VARCHAR2(16),
  TYPE_TAX                  VARCHAR2(64),
  STATEMENT_TYPE            VARCHAR2(64),
  TREATY_GROUP_ID           VARCHAR2(64),
  TREATY_GROUP_NAME         VARCHAR2(255),
  TREATY_GROUP_OLD_ID       VARCHAR2(64),
  OJK_BUSINESS_ID           VARCHAR2(64),
  ID_NEW_BISNIS             VARCHAR2(64),
  BIZ_CODE                  VARCHAR2(64),
  BIZ_NAME                  VARCHAR2(255),
  SOB                       VARCHAR2(64),
  SOB_NAME                  VARCHAR2(255),
  CEDING_CO                 VARCHAR2(4000),
  CEDING_CO_NAME            VARCHAR2(4000),
  INSURED_ID                VARCHAR2(64),
  INSURED_NAME              VARCHAR2(255),
  MARKETING_OFFICER         VARCHAR2(255),
  TREATY_TYPE               VARCHAR2(64),
  TREATY_YEAR               VARCHAR2(16),
  CURRENCY                  VARCHAR2(16),
  ID_CURRENCY               VARCHAR2(64),
  SHARE_CURRENCY            VARCHAR2(16),
  QUARTAL                   VARCHAR2(16),
  YEAR_OF_QUARTAL           VARCHAR2(16),
  CLAIM_TYPE                VARCHAR2(64),
  CLAIM_PAYMENT_TYPE        VARCHAR2(64),
  INSTALLMENT               VARCHAR2(16),
  REMARK                    VARCHAR2(128),
  START_DATE                DATE,
  END_DATE                  DATE,
  STATEMENT_DATE            DATE,
  TGL_PROD                  DATE,
  GROSS_PREMIUM             NUMBER(38,8),
  GROSS_CLAIM               NUMBER(38,8),
  PREMI_OGP                 NUMBER(38,8),
  RESULT_OGP1               NUMBER(38,8),
  RESULT_OGP2               NUMBER(38,8),
  PREMI_ONP                 NUMBER(38,8),
  RESULT_ONP1               NUMBER(38,8),
  RESULT_ONP2               NUMBER(38,8),
  CLAIM                     NUMBER(38,8),
  OUTSTANDING_CLAIM         NUMBER(38,8),
  SALVAGE_VALUE             NUMBER(38,8),
  EXCESS_LOSS               NUMBER(38,8),
  NET_PREMIUM               NUMBER(38,8),
  BALANCE_DUE_TO            NUMBER(38,8),
  BALANCE_BEFORE_TAX        NUMBER(38,8),
  BALANCE_BEFORE_PPH        NUMBER(38,8),
  DEDUCTION1                NUMBER(38,8),
  DEDUCTION2                NUMBER(38,8),
  BROKERAGE_FEE_SEBENARNYA  NUMBER(38,8),
  PPH_VALUE                 NUMBER(38,8),
  PPN_VALUE                 NUMBER(38,8),
  SHARE_VALUE               NUMBER(38,8),
  RI_COMM_OGP               NUMBER(38,8),
  OVERIDDING_COMM_OGP       NUMBER(38,8),
  RI_COMM_ONP               NUMBER(38,8),
  OVERIDDING_COMM_ONP       NUMBER(38,8),
  CONSTRAINT PK_GENERAL_POLIS PRIMARY KEY (ID),
  CONSTRAINT FK_GENERAL_POLIS_WORK FOREIGN KEY (ID) REFERENCES {skema}.T_WORK_POLIS (ID),
  CONSTRAINT FK_GENERAL_POLIS_OLD FOREIGN KEY (OLD_POLIS_ID) REFERENCES {skema}.T_WORK_POLIS (ID),
  CONSTRAINT UQ_GENERAL_POLIS_OLD UNIQUE (OLD_POLIS_ID)
)
/
CREATE UNIQUE INDEX {skema}.UQ_GENERAL_POLIS_NOPOLIS ON {skema}.T_GENERAL_POLIS (CASE WHEN NOPOLIS IS NOT NULL THEN NOPOLIS END, CASE WHEN NOPOLIS IS NOT NULL THEN PRODKE END)
/

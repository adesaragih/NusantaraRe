-- ============================================================================
-- DDL DRAF - TABEL FLAT FACULTATIVE INWARD
-- ============================================================================
-- Dibangkitkan  : 25 September 2026 (regenerasi penuh)
-- Sumber        : RNM\OUTPUT\08-flat\Tabel-Flat-Lintas-Siklus.xlsx
--                 lembar Kolom (tipe, default, nullability) dan Daftar Relasi (kunci tamu)
-- Dasar contoh  : DDL\CONTOH\*.xml = 115 berkas, populasi efektif 113
-- Keputusan     : K-018 K-063 K-064 K-065 K-066 K-067 K-068 K-069
--                 ADR-0001 ADR-0004 ADR-0005 ADR-0006 ADR-0007, V-14 V-15 V-20 V-23
--
-- STATUS: DRAF. Belum pernah dijalankan terhadap Oracle mana pun.
-- SELURUH BUTIR RANCANGAN SUDAH DIPUTUSKAN - tidak ada lagi yang menggantung.
--
-- CATATAN YANG MENGIKAT
--   Uang disimpan NUMBER, bukan BINARY_DOUBLE (CLAUDE.md par.4.1).
--   Tidak ada identifier berkutip (V-15). Seluruh nama di bawah 30 karakter (V-20).
--   Kolom data dibuat NULL: di korpus sebagian besar medan kosong di sebagian besar berkas,
--   sehingga NOT NULL akan menolak baris lama saat migrasi (K-063).
--   PENGECUALIAN: 24 kolom pembawa mata uang kelompok (a)(b)(c) dibuat NOT NULL dengan
--   DEFAULT UNKNOWN (K-069 usulan 10, ADR-0006) - DEFAULT-nya yang membuat itu aman.
--   OLD_POLIS_ID menaut versi sebelumnya, UNIK, TANPA FK karena ditegakkan di Go (K-068).
-- ============================================================================

-- ---------------------------------------------------------------------------
-- 1. TABEL
-- ---------------------------------------------------------------------------

CREATE TABLE T_WORK_POLIS (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  JENIS_WORK                   VARCHAR2(10),
  NO_WORK                      VARCHAR2(20),
  POSISI                       VARCHAR2(100),
  NOURUT                       NUMBER(3),
  PUTARAN                      NUMBER,
  STATUS_PROSES                VARCHAR2(20),
  STS_KONVERSI                 NUMBER,
  TGL_KONVERSI                 DATE,
  TGL_INPUT                    DATE,
  USERNAME                     VARCHAR2(50),
  DATE_TO_UW                   VARCHAR2(30),
  ID_NEW_BISNIS                VARCHAR2(20),
  IS_FLAG_REJECT               VARCHAR2(10),
  PIC                          VARCHAR2(100),
  POLICY_STATUS                VARCHAR2(20),
  SUBMIT_TO_UW                 VARCHAR2(30)
);

CREATE TABLE T_GENERAL_POLIS (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  OLD_POLIS_ID                 NUMBER,
  COB_GROUP                    VARCHAR2(20),
  ADDITIONAL_CAPITAL           VARCHAR2(50),
  BINDER_RNM                   VARCHAR2(50),
  CEDING_RETENTION             NUMBER,
  CEDING_RETENTION_NOMINAL     NUMBER,
  COMMENT_TEXT                 VARCHAR2(500),
  COMMENT_CEDING               VARCHAR2(500),
  COMMENT_DEDUCTIBLE           VARCHAR2(500),
  CONDITION                    VARCHAR2(500),
  CONFIRM_BINDING              VARCHAR2(50),
  CURRENT_YEAR                 VARCHAR2(50),
  DATE_OFFER                   VARCHAR2(30),
  DATE_VALIDITY                VARCHAR2(30),
  DAYS_VALIDITY                VARCHAR2(50),
  ERRMSG                       VARCHAR2(50),
  END_DATE_TIME                VARCHAR2(30),
  ENDORSEMENT_NO               VARCHAR2(50),
  ENDORSMENT_REASON            VARCHAR2(50),
  FLAG_SAVE_FO                 VARCHAR2(10),
  FLAG_SPECIAL_ACCEPTANCE      VARCHAR2(10),
  FLAG_SPREADING               VARCHAR2(10),
  FOLLOWING                    VARCHAR2(50),
  FOLLOWING_NB                 VARCHAR2(50),
  GENERATE_DATE_POLICY_NO      VARCHAR2(30),
  ID_FOLLOWING_NB              VARCHAR2(50),
  INPUT_COMMENT_DEDUCTIBLE     VARCHAR2(500),
  IS_B2_B                      VARCHAR2(10),
  IS_BANDING                   VARCHAR2(10),
  IS_DEDUCTIBLE_ACCEPTANCE     VARCHAR2(10),
  IS_EDIT_REFF_NUMBER          VARCHAR2(10),
  IS_FAC_RETRO                 VARCHAR2(10),
  IS_GROSS_NET_SHOW            VARCHAR2(10),
  IS_INPUT_FAC_RETRO           VARCHAR2(10),
  IS_LIFE_PERIOD               VARCHAR2(10),
  IS_OCCUP_EXCEPTION           VARCHAR2(10),
  IS_PREFERRED_RISK            VARCHAR2(10),
  IS_PRO_RATE                  VARCHAR2(10),
  IS_RI_SLIP                   VARCHAR2(10),
  IS_SPECIAL_ACCEPTANCE        VARCHAR2(10),
  IS_TOP_RISK                  VARCHAR2(10),
  OFFERING_DATE                VARCHAR2(30),
  OFFICER_FAC_OUT              VARCHAR2(50),
  PPN_CHECK                    NUMBER,
  PARAM_CONDITION              VARCHAR2(500),
  PAY_INSTALLMENT              NUMBER,
  PAY_PCT_BROKERAGE_FEE        NUMBER,
  PAY_RI_COMMISION             NUMBER,
  PERCENT_SHARE                NUMBER,
  POLICY_MASTER_ID_PEGA        VARCHAR2(50),
  POLICY_MASTER_NUMBER         VARCHAR2(50),
  POLICY_NO                    VARCHAR2(50),
  PRO_RATE_PERCENT             NUMBER,
  PRO_RATE_TYPE                NUMBER,
  PROD_DATE_TIME               VARCHAR2(30),
  PROD_KE                      NUMBER(5),
  PRODUCT_BRIGUNA              VARCHAR2(50),
  PRORATE_EDM_END              NUMBER,
  PRORATE_START_EDM            NUMBER,
  RECEIVED_RI_SLIP             VARCHAR2(50),
  SAVE_FAC_OUT                 VARCHAR2(50),
  SHARE_CEDANT_TYPE            NUMBER,
  SHARE_RNML                   NUMBER,
  SOURCE_ID                    VARCHAR2(50),
  START_DATE_TIME              VARCHAR2(30),
  SUM_TOTAL_TSI                NUMBER,
  TOTAL_PREMI_NUSA_RE          NUMBER,
  TOTAL_RI_COM_NUSA_RE         VARCHAR2(50),
  TOTAL_TSI_NUSA_RE            NUMBER,
  TOTAL_TSI_NUSA_RE_SPREADING  NUMBER,
  TOTAL_TSI_TOP_RISK           NUMBER,
  TREATY_CAPACITY              NUMBER,
  WAITING_BIND_DATE            VARCHAR2(30)
);

CREATE TABLE T_COVERAGELIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ASMTSI                       NUMBER,
  ACCUMULATION_CODE            VARCHAR2(50),
  ACCUMULATION_DESCRIPTION     VARCHAR2(500),
  ANEKA_ID                     VARCHAR2(50),
  BEGIN_DATE                   VARCHAR2(30),
  CALCULATE_METHOD             VARCHAR2(50),
  CALCULATE_METHOD_FAC_IN      VARCHAR2(50),
  CARGO_ID                     VARCHAR2(50),
  CONDITIONS                   VARCHAR2(500),
  COVERAGE                     VARCHAR2(50),
  COVERAGE_BASIS               VARCHAR2(50),
  COVERAGE_NOTE                VARCHAR2(500),
  COVERAGE_NOTE_ENGLISH        VARCHAR2(500),
  DAY                          VARCHAR2(50),
  DISCOUNT                     NUMBER,
  DISCOUNT_PERCENTAGE          NUMBER,
  EML_PML                      VARCHAR2(50),
  END_DATE                     VARCHAR2(30),
  FIRST_LOSS                   VARCHAR2(50),
  FIRST_LOSS_SCALE             VARCHAR2(50),
  FIRST_SCALE                  VARCHAR2(50),
  FLAG_RATE_CHANGE             VARCHAR2(10),
  GROSS_PREMIUM_RETRO          NUMBER,
  INDEMNITY                    VARCHAR2(50),
  INDEMNITY_PERCENTAGE         NUMBER,
  INDEX_COVEREGE_MASTER        NUMBER,
  INDEX_LOOP                   NUMBER,
  IS_OLD_DATA                  VARCHAR2(10),
  IS_SHOW_FORMULA              VARCHAR2(10),
  LIMITOF_LIABILITY            NUMBER,
  LOADING                      NUMBER,
  LOCATION_ID                  VARCHAR2(50),
  LOST_LIMIT                   VARCHAR2(50),
  MAX_STANDARD_RATE            NUMBER,
  MIN_PREMIUM                  NUMBER,
  MIN_STANDARD_RATE            NUMBER,
  NET_RATE                     NUMBER,
  OLDID                        VARCHAR2(50),
  OBJECT_INDEX                 VARCHAR2(50),
  OCCUPATION_ID                VARCHAR2(50),
  PCT_ADJUSTMENT               NUMBER,
  PCT_LO_L                     NUMBER,
  PCT_SHARE_CEDING             NUMBER,
  PCT_SHARE_NUSANTARA_RE       NUMBER,
  PCT_SHORT_PERIOD             VARCHAR2(30),
  PERCENTAGE_ADJUSTMENT        NUMBER,
  PREMI_CHANGE_FLAG            NUMBER,
  PREMI_LIFE_NUSANTARA_RE      NUMBER,
  PREMI_NUSANTARA_RE           NUMBER,
  PREMI_RP                     NUMBER,
  PREMIUM                      NUMBER,
  PREMIUM_GROSS_DISC_FLEET     NUMBER,
  PREMIUM_RETRO                NUMBER,
  PRO_RATE_PERCENT             NUMBER,
  PROPERTY_ID                  VARCHAR2(50),
  PROPERTY_ITEM_ID             VARCHAR2(50),
  RI_RISK                      VARCHAR2(50),
  RATE                         NUMBER,
  RATE_LIFE_AVERAGE            NUMBER,
  RATE_LIFE_RETRO              NUMBER,
  RATE_OJK                     NUMBER,
  RISK_ZIP_CODE                VARCHAR2(50),
  SUB_LIMIT_NOTE               VARCHAR2(500),
  SUBLIMIT                     NUMBER,
  TSI                          NUMBER,
  TSI_CEDING                   NUMBER,
  TSI_LIABILITY                NUMBER,
  TSI_NUSANTARA_RE             NUMBER,
  TSI_OBJECT_ITEM              NUMBER,
  TSI_SUBLIMIT                 NUMBER,
  TYPE_OF_DISCOUNT             NUMBER,
  UNIT                         VARCHAR2(50),
  USED                         VARCHAR2(50),
  VEHICLE_ID                   VARCHAR2(50),
  ZONE                         VARCHAR2(50),
  ZONE_NOTE                    VARCHAR2(500),
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_QUOTATIONDATA (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  BRANCH_CODE                  VARCHAR2(50),
  BRANCH_NAME                  VARCHAR2(500),
  BUSINESS_CODE                VARCHAR2(50),
  BUSINESS_FAC                 VARCHAR2(50),
  BUSINESS_NAME                VARCHAR2(500),
  BUSINESS_OLD_ID              VARCHAR2(50),
  BUSINESS_TYPE                VARCHAR2(50),
  BUSINESS_TYPE2               VARCHAR2(50),
  CEDING_CO                    VARCHAR2(50),
  CEDING_CO_NAME               VARCHAR2(500),
  COMMENT_TEXT                 VARCHAR2(500),
  COPY_FROM                    VARCHAR2(50),
  EDM_DAY                      VARCHAR2(50),
  EDM_CHARGE_FEE               NUMBER,
  EDM_DATE                     VARCHAR2(30),
  EDM_NOTE                     VARCHAR2(500),
  EDM_SOURCE                   VARCHAR2(50),
  EDM_SOURCE_NOTE              VARCHAR2(500),
  EDM_STATUS                   VARCHAR2(50),
  EDM_SURVEY                   VARCHAR2(50),
  EDM_TYPE                     VARCHAR2(50),
  EDM_TYPE_NEW                 VARCHAR2(50),
  EMAIL                        VARCHAR2(50),
  ENDORSEMENT_INTERNAL_RETRO   VARCHAR2(50),
  FAC_OUT_STATUS               VARCHAR2(50),
  FRONTING_STATUS              VARCHAR2(50),
  GROUP_NAME                   VARCHAR2(500),
  GROUP_PANEL                  VARCHAR2(50),
  INPUT_MONTH                  VARCHAR2(50),
  INPUT_QUO_DATE               VARCHAR2(30),
  INSURED_ID                   VARCHAR2(50),
  INSURED_NAME                 VARCHAR2(500),
  IS_EDM_SHARE                 VARCHAR2(10),
  IS_GROUP                     VARCHAR2(10),
  IS_MOP                       VARCHAR2(10),
  IS_PROPOSAL                  VARCHAR2(10),
  IS_SURVEY_REPORT             VARCHAR2(10),
  IS_TBA                       VARCHAR2(10),
  IS_TEMPLATE                  VARCHAR2(10),
  MOID                         VARCHAR2(50),
  MARKETING_CODE               VARCHAR2(50),
  MARKETING_NAME               VARCHAR2(500),
  NO_OFFER_SLIP                VARCHAR2(50),
  OFFERING_DATE                VARCHAR2(30),
  OLD_POLICY_NO                VARCHAR2(50),
  OPERATOR_ID                  VARCHAR2(50),
  PERIOD                       VARCHAR2(30),
  PERIOD_MM                    VARCHAR2(30),
  POLICY_TYPE                  VARCHAR2(50),
  QQID                         VARCHAR2(50),
  QQ_NAME                      VARCHAR2(500),
  RNW_DATE                     VARCHAR2(30),
  SHARE_OF_CEDING              NUMBER,
  SOB_LEADER0                  VARCHAR2(50),
  SOB_LSG                      VARCHAR2(50),
  SOB_NAME                     VARCHAR2(500),
  SOURCE_OF_BUSINESS           VARCHAR2(50),
  STATUS_BUSINESS              VARCHAR2(50),
  STATUS_SYARIAH               VARCHAR2(50),
  TEAM_GROUP                   VARCHAR2(50),
  TYPE                         VARCHAR2(50),
  TYPE_FACULTATIVE             VARCHAR2(50),
  BTN_QUOTATION                VARCHAR2(50)
);

CREATE TABLE T_LISTINSTALLMENT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ADMIN_FEE                    NUMBER,
  BROKERAGE_FEE                NUMBER,
  COMMISSION                   NUMBER,
  DEDUCTION2                   VARCHAR2(50),
  DISCOUNT                     NUMBER,
  DUE_DATE                     VARCHAR2(30),
  INSTALLMENT_NO               NUMBER,
  INSTALLMENT_PERCENTAGE       NUMBER,
  PPN                          NUMBER,
  PPH                          NUMBER,
  PAYMENT_AMOUNT               NUMBER,
  PAYMENT_TOTAL                NUMBER,
  PREMIUM                      NUMBER,
  RI_COMMISION                 NUMBER,
  RATE_EM                      NUMBER,
  RATE_LIFE                    NUMBER,
  STAMP                        NUMBER,
  YEAR                         VARCHAR2(50),
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_CURRENCYLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CURRENCY_CODE                VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  NAME                         VARCHAR2(500),
  OLD_ID                       VARCHAR2(50),
  PAY_BROKERAGE_FEE            NUMBER,
  PAY_COMMISION                NUMBER,
  PAY_EDM_NEW_PREMI            NUMBER,
  PAY_EDM_PREMI_MENJADI        NUMBER,
  PAY_NET_PREMIUM              NUMBER,
  PAY_PPN                      NUMBER,
  PAY_PPH                      NUMBER,
  PAY_PREMIUM                  NUMBER,
  PAY_TSI_TOTAL                NUMBER,
  PREMIUM                      NUMBER,
  SUM_TOTAL_PAYMENT            NUMBER,
  TSI                          NUMBER
);

CREATE TABLE T_SPREADINGLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ANEKA_ID                     VARCHAR2(50),
  CLAIM_ESTIMATION             NUMBER,
  COVERAGE_ID                  VARCHAR2(50),
  FLAG_SPREADING               VARCHAR2(10),
  INDEX_PROPERTY_ITEM          NUMBER,
  IS_INPUT_PCT                 VARCHAR2(10),
  IS_OLD_DATA                  VARCHAR2(10),
  LOCATION_ID                  VARCHAR2(50),
  PREMIUM_SPREADED             NUMBER,
  SHARE_PERCENTAGE             NUMBER,
  TSI_GROSS_SPREADED           NUMBER,
  TSI_SPREADED                 NUMBER,
  TREATY_NAME                  VARCHAR2(500),
  TREATY_TYPE                  VARCHAR2(50),
  TSI_TOP_RISK                 NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_CEDINGCOLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CEDING_CO                    VARCHAR2(50),
  CEDING_CO_NAME               VARCHAR2(500)
);

CREATE TABLE T_CEDINGCEDANTLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CEDING_CO                    VARCHAR2(50),
  CEDING_CO_NAME               VARCHAR2(500),
  SHARE_CEDING                 NUMBER
);

CREATE TABLE T_CEDING_CURRENCYLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  NAME                         VARCHAR2(500),
  PREMIUM                      NUMBER,
  TSI                          NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_CURRENCY (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  CURRENCY_REF_ID              VARCHAR2(50),
  NAME                         VARCHAR2(500),
  OLD_ID                       VARCHAR2(50)
);

CREATE TABLE T_DEDUCTIBLELIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  AMOUNT                       NUMBER,
  ANEKA_ID                     VARCHAR2(50),
  BASIS_TYPE                   VARCHAR2(50),
  CLAIM_CATEGORY               NUMBER,
  CONDITION                    VARCHAR2(500),
  COVERAGE                     VARCHAR2(50),
  COVERAGE_ID                  VARCHAR2(50),
  CURRENCY                     VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  DEDUCTIBLE_TYPE              NUMBER,
  DESCRIPTIONS                 VARCHAR2(500),
  FLAG_CURRENCY                VARCHAR2(10),
  INPUT_CONDITION              VARCHAR2(500),
  MIN_MAX                      VARCHAR2(50),
  OCCUPATION_ID                VARCHAR2(50),
  PCT_DEDUCTIBLE               NUMBER,
  PCT_DEDUCTIBLE2              NUMBER,
  PROPERTY_ID                  VARCHAR2(50),
  PROPERTY_ITEM_ID             VARCHAR2(50),
  TIME_EXCESS                  VARCHAR2(30),
  TYPE                         VARCHAR2(50),
  TYPE_DEDUCTIBLE              NUMBER,
  TYPE_DEDUCTIBLE2             NUMBER,
  VALUE                        VARCHAR2(50)
);

CREATE TABLE T_FACRETRODETAILS (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  BACK_UP_STATUS               VARCHAR2(50),
  DOCUMENT_POSITION            VARCHAR2(50)
);

CREATE TABLE T_OCCUPATIONLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CATEGORY                     VARCHAR2(50),
  IS_OLD_DATA                  VARCHAR2(10),
  LOCATION_ID                  VARCHAR2(50),
  OCCUPATION_ID                VARCHAR2(50),
  OCCUPATION_NAME              VARCHAR2(500)
);

CREATE TABLE T_PROPERTY (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  ALM_RISK_ID                  VARCHAR2(50),
  BUILDING_NO                  VARCHAR2(50),
  COUNTRY                      VARCHAR2(50),
  IS_FLAMMABLE_ITEM_FLAG       VARCHAR2(10),
  IS_HOT_WORK_PROCESS_FLAG     VARCHAR2(10),
  IS_MATERIAL_DAMAGE           VARCHAR2(10),
  IS_PRODUCTION_PROCESS_FLAG   VARCHAR2(10),
  IS_TOP_RISK                  VARCHAR2(10),
  OBJECT_NAME                  VARCHAR2(500),
  OBJECT_NO                    VARCHAR2(50),
  OBJECT_TYPE                  VARCHAR2(50),
  OWNERSHIP                    VARCHAR2(50),
  PROVINCE                     VARCHAR2(50),
  ROAD_NAME                    VARCHAR2(500),
  ROAD_TYPE                    VARCHAR2(50),
  TOTAL_TSI                    NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_PERSONLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ASMCC_AMOUNT                 NUMBER,
  ASM_CLASS                    VARCHAR2(50),
  ASM_CLASS_ID                 VARCHAR2(50),
  ASM_DATE_OF_BIRTH            VARCHAR2(30),
  ASM_GENDER                   VARCHAR2(50),
  ASM_JOB_DESC                 VARCHAR2(500),
  ASM_JOB_NAME                 VARCHAR2(500),
  ASM_LEFT_HANDED              VARCHAR2(50),
  ASM_PARTICIPANT_STATUS       VARCHAR2(50),
  AGE                          VARCHAR2(50),
  FLAG_OLD_DATA                VARCHAR2(10),
  MASTER_RATE_COVERAGE         NUMBER,
  MASTER_RATE_COVERAGE_ID      NUMBER,
  RI_RISK_CODE                 VARCHAR2(50),
  RI_RISK_NAME                 VARCHAR2(500),
  CURRENCY_CODE                VARCHAR2(10)
);

CREATE TABLE T_LOCATIONLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CURRENCY_ID                  VARCHAR2(50),
  CURRENCY_OLD_ID              VARCHAR2(50),
  IS_OLD_DATA                  VARCHAR2(10),
  LOSS_RATIO1_YEAR_AMOUNT      NUMBER,
  LOSS_RATIO1_YEAR_PERCENT     VARCHAR2(50),
  LOSS_RATIO35_YEAR_AMOUNT     NUMBER,
  LOSS_RATIO35_YEAR_PERCENT    VARCHAR2(50)
);

CREATE TABLE T_RISKLOCATION (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  ASM_ADDRESS                  VARCHAR2(50),
  ASM_CITY                     VARCHAR2(50),
  ASM_DISTRICT                 VARCHAR2(50),
  ASMRW                        VARCHAR2(50),
  ASM_ZIP_CODE                 VARCHAR2(50)
);

CREATE TABLE T_PROPERTYITEMLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CURRENCY                     VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  CURRENCY_ID                  VARCHAR2(50),
  CURRENCY_OLD_ID              VARCHAR2(50),
  FLAG_NET_RATE                VARCHAR2(10),
  IS_ADJUSTABLE_FLAG           VARCHAR2(10),
  IS_OLD_DATA                  VARCHAR2(10),
  ITEM_TYPE                    VARCHAR2(50),
  ITEM_TYPE_ID                 VARCHAR2(50),
  PCT_ADJUST2                  NUMBER,
  PCT_ADJUST_OTHER             NUMBER,
  PERCENTAGE_ADJUSTMENT        NUMBER,
  PROPERTI_ITEM_NOTE           VARCHAR2(500),
  PROPERTY_ID                  VARCHAR2(50),
  PROPERTY_ITEM_NO             VARCHAR2(50),
  REMARK                       VARCHAR2(500),
  SELECTED_LOCATION_ADDRESS    VARCHAR2(50),
  SELECTED_OBJECT_ITEM         VARCHAR2(50),
  TSI_OBJECT_ITEM              NUMBER,
  TOTAL_GROSS_PREMI            NUMBER,
  TOTAL_NET_RATE               NUMBER,
  TOTAL_PREMIUM_NUSANTARA_RE   NUMBER
);

CREATE TABLE T_ANEKALIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  BUILT_IN                     VARCHAR2(50),
  DRAFT_WORDING_ID             VARCHAR2(50),
  END_DATE                     VARCHAR2(30),
  FLAG                         VARCHAR2(10),
  GEOGRAPHICAL_LIMITS          VARCHAR2(50),
  IS_OLD_DATA                  VARCHAR2(10),
  LOCATION_ID                  VARCHAR2(50),
  OBJECT_NAME                  VARCHAR2(500),
  OCCUPATION_ID                VARCHAR2(50),
  QUANTITY                     VARCHAR2(50),
  START_DATE                   VARCHAR2(30),
  SECTION                      VARCHAR2(50),
  SELECTED_LOCATION_ADDRESS    VARCHAR2(50),
  SELECTED_OBJECT_ITEM         VARCHAR2(50),
  SERIAL_NO                    VARCHAR2(50),
  TSI                          NUMBER,
  UNIT_QUANTITY                VARCHAR2(50),
  YEAR                         VARCHAR2(50),
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_VEHICLELIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  BRAND                        VARCHAR2(50),
  BRAND_NAME                   VARCHAR2(500),
  CC                           VARCHAR2(50),
  CAR_PARKING                  VARCHAR2(50),
  CHASSIS_NUMBER               VARCHAR2(50),
  COLOR                        VARCHAR2(50),
  COLOR_NAME                   VARCHAR2(500),
  ENGINE_NUMBER                VARCHAR2(50),
  FLAG_OLD_DATA                VARCHAR2(10),
  LICENSE_PLATE                VARCHAR2(50),
  MANUFACTURE_YEAR             VARCHAR2(50),
  MODEL                        VARCHAR2(50),
  MODEL_NAME                   VARCHAR2(500),
  TRANSMISSION_TYPE            VARCHAR2(50),
  TYPE                         VARCHAR2(50),
  TYPE_NAME                    VARCHAR2(500),
  VEHICLE_CODE_ID              VARCHAR2(50),
  VEHICLE_CODE_NOTE            VARCHAR2(500)
);

CREATE TABLE T_RETROLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  COMMISION                    NUMBER,
  COMMISION_AMOUNT             NUMBER,
  IDTREATYYEAR_LIFE            VARCHAR2(50),
  OVR_COMM                     VARCHAR2(50),
  OVR_COMM_AMOUNT              NUMBER,
  PERCENTSHARE                 NUMBER,
  PREMIUM_SPREADED_GROSS       NUMBER,
  PREMIUM_SPREADED_NET         NUMBER,
  RATE                         NUMBER,
  REINSURERNAME                VARCHAR2(500),
  RETRO_REF_ID                 VARCHAR2(50),
  TGLUPDATE                    VARCHAR2(30),
  TREATYENDDATE                VARCHAR2(30),
  TREATYSTARTDATE              VARCHAR2(30),
  TREATYTYPEID                 VARCHAR2(50),
  TREATYTYPENAME               VARCHAR2(500),
  USERID                       VARCHAR2(50),
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_CARGOLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  BL_NUMBER                    VARCHAR2(50),
  CONVEYANCE_ID                VARCHAR2(50),
  CONVEYANCE_NAME              VARCHAR2(500),
  CONVEYANCE_NOTE              VARCHAR2(500),
  FLAG_OLD_DATA                VARCHAR2(10),
  FROM_RUTE                    VARCHAR2(50),
  GOOD_ID                      VARCHAR2(50),
  GOOD_NOTE                    VARCHAR2(500),
  INVOICE_NUMBER               VARCHAR2(50),
  IS_TOP_RISK                  VARCHAR2(10),
  PACKING_ID                   VARCHAR2(50),
  PACKING_NOTE                 VARCHAR2(500),
  SAIL_DATE                    VARCHAR2(30),
  TO_RUTE                      VARCHAR2(50),
  TRADING_ID                   VARCHAR2(50),
  TRADING_NOTE                 VARCHAR2(500)
);

CREATE TABLE T_SCORINGRESULT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  ABOVE_AVERAGE                VARCHAR2(50),
  ACCEPTABLE_NOT_APPROVAL      VARCHAR2(50),
  ACCEPTABLE_WITH_APPROVAL     VARCHAR2(50),
  AVERAGE                      VARCHAR2(50),
  BELOW_AVERAGE                VARCHAR2(50),
  GOOD                         VARCHAR2(50),
  LOSS_RATIO                   VARCHAR2(50),
  NOT_ACCEPTABLE               VARCHAR2(50),
  OBJECT_CONDITIONS            VARCHAR2(500),
  OCCUPATION                   VARCHAR2(50),
  OPERATIONAL_DIRECTOR         VARCHAR2(50),
  OPERATIONAL_DIV_HEAD         VARCHAR2(50),
  OTHERS                       VARCHAR2(50),
  POOR                         VARCHAR2(50),
  PRESIDENT_DIRECTOR           VARCHAR2(50),
  TECHNICAL_DIRECTOR           VARCHAR2(50)
);

CREATE TABLE T_SHIP (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  AGE                          VARCHAR2(50),
  CONST                        VARCHAR2(50),
  DWT                          VARCHAR2(50),
  FLAG_NAME                    VARCHAR2(500),
  GRT                          VARCHAR2(50),
  IMO                          VARCHAR2(50),
  NM_SHIP                      VARCHAR2(50),
  NRT                          VARCHAR2(50),
  REGISTER                     VARCHAR2(50),
  REMARK                       VARCHAR2(500),
  SHIPTYPE_ID                  VARCHAR2(50),
  SHIPTYPE_NAME                VARCHAR2(500),
  SHIP_REF_ID                  VARCHAR2(50),
  Y_MAKE1                      VARCHAR2(50),
  Y_MAKE2                      VARCHAR2(50)
);

CREATE TABLE T_ADDITIONALCOVERAGE (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CALCULATE_METHOD             VARCHAR2(50),
  CALCULATE_METHOD_FAC_IN      VARCHAR2(50),
  COVERAGE                     VARCHAR2(50),
  COVERAGE_ID                  VARCHAR2(50),
  COVERAGE_NOTE                VARCHAR2(500),
  IS_OLD_DATA                  VARCHAR2(10),
  PREMIUM                      NUMBER,
  PRO_RATE_PERCENT             NUMBER,
  RATE                         NUMBER,
  TSI                          NUMBER,
  VEHICLE_ID                   VARCHAR2(50),
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_COVERAGEDATALIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  COVERAGE                     VARCHAR2(50),
  COVERAGE_NOTE                VARCHAR2(500),
  DISCOUNT_PERCENTAGE          NUMBER,
  RATE                         NUMBER,
  TSI                          NUMBER,
  TSI_CEDING                   NUMBER,
  TSI_LIABILITY                NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_AVIATIONHULL (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  MAX_PASSENGERS               VARCHAR2(50),
  PILOTS                       VARCHAR2(50),
  REGISTRATION_MARKS           VARCHAR2(50),
  SPECIAL_RENTAL_USES          VARCHAR2(50),
  SPECIAL_USES                 VARCHAR2(50),
  STANDARD_USES                VARCHAR2(50)
);

CREATE TABLE T_MARINEHULL (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  CLASSIFICATION               VARCHAR2(50),
  CONSTRUCTION                 VARCHAR2(50),
  DWTRT                        VARCHAR2(50),
  GRTRT                        VARCHAR2(50),
  NRTRT                        VARCHAR2(50),
  TYPE_OF_VESSEL               VARCHAR2(50)
);

CREATE TABLE T_SCORING_FACTOR (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CHECKED                      VARCHAR2(10),
  FACTOR_GROUP                 VARCHAR2(50),
  FACTOR_NAME                  VARCHAR2(500),
  REMARKS                      VARCHAR2(500),
  SCORE_PER_FACTOR             NUMBER
);

CREATE TABLE T_VEHICLEHE (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  BRAND                        VARCHAR2(50),
  BRAND_NAME                   VARCHAR2(500),
  OBJECT_NAME_HE               VARCHAR2(500),
  TYPE                         VARCHAR2(50),
  TYPE_NAME                    VARCHAR2(500)
);

CREATE TABLE T_ACCESSORYLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  BRAND_NAME                   VARCHAR2(500),
  PRICE                        VARCHAR2(50),
  TYPE_NAME                    VARCHAR2(500),
  VALUE_NAME                   VARCHAR2(500)
);

CREATE TABLE T_BUILDINGCONSTRUCTION (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  FLOOR_TYPE                   VARCHAR2(50),
  NUMBER_OF_FLOOR              VARCHAR2(50),
  ROOF_TYPE                    VARCHAR2(50),
  WALL_TYPE                    VARCHAR2(50)
);

CREATE TABLE T_DATASCORINGRISKLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  OCCUPATION_CODE              VARCHAR2(50),
  OCCUPATION_NOTE              VARCHAR2(500),
  SCORE                        NUMBER,
  TOP_RISK_LOCATION            VARCHAR2(50)
);

CREATE TABLE T_LISTCAUSEOFLOSS (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CLAIM                        NUMBER,
  CURRENCY                     VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  DETAIL                       VARCHAR2(50),
  REMARKS                      VARCHAR2(500)
);

CREATE TABLE T_SURVEYREPORTLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  DATEOF_SURVEY                VARCHAR2(30),
  LOSS_PREVENTION              VARCHAR2(50),
  REMARKS                      VARCHAR2(500),
  SURVEYED_BY                  VARCHAR2(50)
);

CREATE TABLE T_SCORINGRISK (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  FINAL_SCORE                  NUMBER,
  NOTE_FINAL_SCORE             VARCHAR2(500),
  STATUS                       VARCHAR2(50)
);

CREATE TABLE T_SCORING_OPTION (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CHECK_BOX                    VARCHAR2(10),
  OPTION_NO                    VARCHAR2(50),
  SCORE                        NUMBER
);

CREATE TABLE T_TABLEOFLIMIT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  CATEGORY                     VARCHAR2(50),
  DESCRIPTION                  VARCHAR2(500),
  PCT_LIMIT                    NUMBER
);

CREATE TABLE T_ASMHEIR (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ASM_GENDER                   VARCHAR2(50),
  ASM_RELATION                 VARCHAR2(50)
);

CREATE TABLE T_SURROUNDINGRISK (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  FLOOD_AREA_STATUS            VARCHAR2(50),
  HOUSEKEEPING_STATUS          VARCHAR2(50)
);

CREATE TABLE T_COINSDATA (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  COINS_NAME                   VARCHAR2(500)
);

CREATE TABLE T_LC (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  LC                           VARCHAR2(50)
);

CREATE TABLE T_QUOTATIONLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  NO_OFFER_SLIP                VARCHAR2(50)
);

CREATE TABLE T_FR_COVERAGELIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ACCUMULATION_CODE            VARCHAR2(50),
  ACCUMULATION_DESCRIPTION     VARCHAR2(500),
  ANEKA_ID                     VARCHAR2(50),
  CALCULATE_METHOD             VARCHAR2(50),
  CALCULATE_METHOD_FAC_IN      VARCHAR2(50),
  CARGO_ID                     VARCHAR2(50),
  COVERAGE                     VARCHAR2(50),
  COVERAGE_BASIS               VARCHAR2(50),
  COVERAGE_NOTE                VARCHAR2(500),
  COVERAGE_NOTE_ENGLISH        VARCHAR2(500),
  DAY                          VARCHAR2(50),
  DISCOUNT                     NUMBER,
  DISCOUNT_PERCENTAGE          NUMBER,
  EML_PML                      VARCHAR2(50),
  FAC_OUT_TSI                  NUMBER,
  FIRST_LOSS                   VARCHAR2(50),
  FIRST_SCALE                  VARCHAR2(50),
  INDEMNITY                    VARCHAR2(50),
  INDEMNITY_PERCENTAGE         NUMBER,
  INDEX_LOCATION               NUMBER,
  IS_OLD_DATA                  VARCHAR2(10),
  IS_SHOW_FORMULA              VARCHAR2(10),
  LIMITOF_LIABILITY            NUMBER,
  LOADING                      NUMBER,
  LOST_LIMIT                   VARCHAR2(50),
  MAX_STANDARD_RATE            NUMBER,
  MIN_PREMIUM                  NUMBER,
  MIN_STANDARD_RATE            NUMBER,
  NET_RATE                     NUMBER,
  NOMINAL_SHARE_OFFERED        NUMBER,
  OLDID                        VARCHAR2(50),
  OBJECT_INDEX                 VARCHAR2(50),
  OCCUPATION_ID                VARCHAR2(50),
  PCT_ADJUSTMENT               NUMBER,
  PCT_LO_L                     NUMBER,
  PCT_SHORT_PERIOD             VARCHAR2(30),
  PERCENT_FAC_OUT              VARCHAR2(50),
  PREMI_NUSANTARA_RE           NUMBER,
  PREMIUM                      NUMBER,
  PREMIUM_RETRO                NUMBER,
  PRO_RATE_PERCENT             NUMBER,
  PROPERTY_ID                  VARCHAR2(50),
  PROPERTY_ITEM_ID             VARCHAR2(50),
  RI_COMM                      NUMBER,
  RI_COMM_PERCENTAGE           NUMBER,
  RATE                         NUMBER,
  RISK_ZIP_CODE                VARCHAR2(50),
  SUBLIMIT                     NUMBER,
  TSI                          NUMBER,
  TSI_LIABILITY                NUMBER,
  TSI_NUSANTARA_RE             NUMBER,
  UNIT                         VARCHAR2(50),
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FR_CURRENCYLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  KURS                         VARCHAR2(50),
  NAME                         VARCHAR2(500),
  OLD_ID                       VARCHAR2(50),
  OLD_TSI                      NUMBER,
  PAY_BROKERAGE_FEE            NUMBER,
  PAY_COMMISION                NUMBER,
  PAY_EDM_NEW_PREMI            NUMBER,
  PAY_EDM_PREMI_MENJADI        NUMBER,
  PAY_NET_PREMIUM              NUMBER,
  PAY_PPN                      NUMBER,
  PAY_PPH                      NUMBER,
  PAY_PREMIUM                  NUMBER,
  PAY_TSI_TOTAL                NUMBER,
  PREMI_NUSANTARA_RE           NUMBER,
  PREMIUM                      NUMBER,
  SHARE_IN_TSI                 NUMBER,
  SUM_TOTAL_PAYMENT            NUMBER,
  TSI                          NUMBER,
  TSI_ADD_CAP                  NUMBER,
  TSI_NUSANTARA_RE             NUMBER,
  TOTAL_FAC_OUT_TSI            NUMBER,
  TOTAL_SHARE_OFFERED          NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FACRETROLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ATTENTION                    VARCHAR2(50),
  COMMENTFORM                  VARCHAR2(500),
  CONFIRCEDING                 VARCHAR2(50),
  DOCUMENT_STATUS              VARCHAR2(50),
  END_PERIOD                   VARCHAR2(30),
  NO_OFFER_SLIP                VARCHAR2(50),
  OUR_REF                      VARCHAR2(50),
  PCT_PREMI_ALL_OBJ_USD        NUMBER,
  PCT_SHARE_ALL_OBJ            NUMBER,
  REINSURER_ID                 VARCHAR2(50),
  REINSURER_NAME               VARCHAR2(500),
  START_PERIOD                 VARCHAR2(30)
);

CREATE TABLE T_FR_SPREADINGLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CLAIM_ESTIMATION             NUMBER,
  FLAG_SPREADING               VARCHAR2(10),
  IS_INPUT_PCT                 VARCHAR2(10),
  IS_OLD_DATA                  VARCHAR2(10),
  PREMIUM_SPREADED             NUMBER,
  SHARE_PERCENTAGE             NUMBER,
  TSI_GROSS_SPREADED           NUMBER,
  TSI_SPREADED                 NUMBER,
  TREATY_NAME                  VARCHAR2(500),
  TREATY_TYPE                  VARCHAR2(50),
  TSI_TOP_RISK                 NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FR_FACOUTOBJECTLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  OBJECT_PREMI                 NUMBER,
  PERCENT_OFFERED              VARCHAR2(50),
  RATE                         NUMBER,
  RI_COMM                      NUMBER,
  SHARE_OFFERED                NUMBER,
  TF                           VARCHAR2(50),
  UJRAH                        VARCHAR2(50),
  COMMISION                    NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FR_NETPERCURRENCY (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  NAME                         VARCHAR2(500),
  OLD_ID                       VARCHAR2(50),
  PREMIUM                      NUMBER,
  SUM_TOTAL_PAYMENT            NUMBER,
  TSI                          NUMBER,
  TSI_NUSANTARA_RE             NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FR_PRINTRISLIP (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  CONFIRM_BY                   VARCHAR2(50),
  DATE_PRINT                   VARCHAR2(30),
  FACRETROSLIPNUMBER           VARCHAR2(50),
  PRINT_DATE                   VARCHAR2(30),
  REMARKS                      VARCHAR2(500),
  WARR_PAYMENT                 NUMBER
);

CREATE TABLE T_FR_DEDUCTIBLELIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  AMOUNT                       NUMBER,
  CONDITION                    VARCHAR2(500),
  COVERAGE_ID                  VARCHAR2(50),
  CURRENCY                     VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  INPUT_CONDITION              VARCHAR2(500),
  MIN_MAX                      VARCHAR2(50),
  PCT_DEDUCTIBLE               NUMBER,
  PCT_DEDUCTIBLE2              NUMBER,
  PROPERTY_ID                  VARCHAR2(50),
  PROPERTY_ITEM_ID             VARCHAR2(50),
  TYPE_DEDUCTIBLE              NUMBER,
  TYPE_DEDUCTIBLE2             NUMBER
);

CREATE TABLE T_FR_WPC (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  DAYS                         VARCHAR2(50),
  DUE_ON                       VARCHAR2(50),
  IP                           VARCHAR2(50),
  PERCENT                      VARCHAR2(50),
  START_FROM                   VARCHAR2(50)
);

CREATE TABLE T_FR_CURRENCY (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  CURRENCY_REF_ID              VARCHAR2(50),
  NAME                         VARCHAR2(500),
  OLD_ID                       VARCHAR2(50)
);

CREATE TABLE T_FR_PROPERTY (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  ALM_RISK_ID                  VARCHAR2(50),
  COUNTRY                      VARCHAR2(50),
  IS_FLAMMABLE_ITEM_FLAG       VARCHAR2(10),
  IS_HOT_WORK_PROCESS_FLAG     VARCHAR2(10),
  IS_MATERIAL_DAMAGE           VARCHAR2(10),
  IS_PRODUCTION_PROCESS_FLAG   VARCHAR2(10),
  IS_TOP_RISK                  VARCHAR2(10),
  OBJECT_NAME                  VARCHAR2(500),
  OBJECT_NO                    VARCHAR2(50),
  OBJECT_NO_FAC_IN             VARCHAR2(50),
  OBJECT_TYPE                  VARCHAR2(50),
  OWNERSHIP                    VARCHAR2(50),
  PROVINCE                     VARCHAR2(50),
  ROAD_NAME                    VARCHAR2(500),
  ROAD_TYPE                    VARCHAR2(50),
  TOTAL_TSI                    NUMBER
);

CREATE TABLE T_FR_LISTINSTALLMENT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  BROKERAGE_FEE                NUMBER,
  DEDUCTION2                   VARCHAR2(50),
  DISCOUNT                     NUMBER,
  DUE_DATE                     VARCHAR2(30),
  INSTALLMENT_NO               NUMBER,
  INSTALLMENT_PERCENTAGE       NUMBER,
  PPN                          NUMBER,
  PPH                          NUMBER,
  PAYMENT_TOTAL                NUMBER,
  PREMIUM                      NUMBER,
  RI_COMMISION                 NUMBER,
  STAMP                        NUMBER
);

CREATE TABLE T_FR_PAYMENT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  BROKERAGE_FEE                NUMBER,
  COMMISION                    NUMBER,
  EDM_NEW_PREMI                NUMBER,
  EDM_PREMI_MENJADI            NUMBER,
  NET_PREMIUM                  NUMBER,
  PPN                          NUMBER,
  PPH                          NUMBER,
  PREMIUM                      NUMBER,
  TSI_TOTAL                    NUMBER
);

CREATE TABLE T_FR_SECURITYREINSURER (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  BROKER_ID                    VARCHAR2(50),
  BROKER_NAME                  VARCHAR2(500),
  DEDUCTION                    VARCHAR2(50),
  PCT_DEDUCTION                NUMBER,
  PCT_SHARE                    NUMBER,
  PREMIUM                      NUMBER,
  REINSURER_ID                 VARCHAR2(50),
  REINSURER_NAME               VARCHAR2(500),
  TSI_OFFER                    NUMBER
);

CREATE TABLE T_FR_OFFEREDPAYMENT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  NAME                         VARCHAR2(500),
  OLD_ID                       VARCHAR2(50),
  PREMIUM                      NUMBER,
  SUM_TOTAL_PAYMENT            NUMBER,
  TSI                          NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FR_RISKLOCATION (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ASM_ADDRESS                  VARCHAR2(50),
  ASM_CITY                     VARCHAR2(50),
  ASM_DISTRICT                 VARCHAR2(50),
  ASMRW                        VARCHAR2(50),
  ASM_ZIP_CODE                 VARCHAR2(50)
);

CREATE TABLE T_FR_OCCUPATIONLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CATEGORY                     VARCHAR2(50),
  IDX_LOCATION                 NUMBER,
  OCCUPATION_ID                VARCHAR2(50),
  OCCUPATION_NAME              VARCHAR2(500)
);

CREATE TABLE T_FR_CARGOLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CONVEYANCE_ID                VARCHAR2(50),
  CONVEYANCE_NAME              VARCHAR2(500),
  CONVEYANCE_NOTE              VARCHAR2(500),
  FAC_RETRO_ID                 VARCHAR2(50),
  FLAG_OLD_DATA                VARCHAR2(10),
  FROM_RUTE                    VARCHAR2(50),
  GOOD_ID                      VARCHAR2(50),
  GOOD_NOTE                    VARCHAR2(500),
  INVOICE_NUMBER               VARCHAR2(50),
  IS_TOP_RISK                  VARCHAR2(10),
  PACKING_ID                   VARCHAR2(50),
  PACKING_NOTE                 VARCHAR2(500),
  PERCENT_OFFERED              VARCHAR2(50),
  SAIL_DATE                    VARCHAR2(30),
  SHARE_OFFERED                NUMBER,
  TSI_NUSANTARA_RE             NUMBER,
  TSI_SPREADED                 NUMBER,
  TO_RUTE                      VARCHAR2(50),
  TRADING_ID                   VARCHAR2(50),
  TRADING_NOTE                 VARCHAR2(500)
);

CREATE TABLE T_FR_ANEKALIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  BUILT_IN                     VARCHAR2(50),
  END_DATE                     VARCHAR2(30),
  FAC_RETRO_ID                 VARCHAR2(50),
  IDX_LOCATION                 NUMBER,
  OBJECT_NAME                  VARCHAR2(500),
  OCCUPATION_ID                VARCHAR2(50),
  PERCENT_OFFERED              VARCHAR2(50),
  START_DATE                   VARCHAR2(30),
  SECTION                      VARCHAR2(50),
  SELECTED_LOCATION_ADDRESS    VARCHAR2(50),
  SELECTED_OBJECT_ITEM         VARCHAR2(50),
  SHARE_OFFERED                NUMBER,
  TSI                          NUMBER,
  TSI_NUSANTARA_RE             NUMBER,
  TSI_SPREADED                 NUMBER,
  YEAR                         VARCHAR2(50),
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FR_PROPERTYITEMLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CURRENCY                     VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  CURRENCY_ID                  VARCHAR2(50),
  CURRENCY_OLD_ID              VARCHAR2(50),
  FLAG_NET_RATE                VARCHAR2(10),
  IS_ADJUSTABLE_FLAG           VARCHAR2(10),
  ITEM_TYPE                    VARCHAR2(50),
  ITEM_TYPE_ID                 VARCHAR2(50),
  PCT_ADJUST2                  NUMBER,
  PCT_ADJUST_OTHER             NUMBER,
  PROPERTI_ITEM_NOTE           VARCHAR2(500),
  PROPERTY_ID                  VARCHAR2(50),
  PROPERTY_ITEM_NO             VARCHAR2(50),
  TSI_OBJECT_ITEM              NUMBER,
  TOTAL_GROSS_PREMI            NUMBER,
  TOTAL_NET_RATE               NUMBER,
  TOTAL_PREMIUM_NUSANTARA_RE   NUMBER
);

CREATE TABLE T_FR_PERSONLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  ASMCC_AMOUNT                 NUMBER,
  ASM_CLASS                    VARCHAR2(50),
  ASM_CLASS_ID                 VARCHAR2(50),
  ASM_DATE_OF_BIRTH            VARCHAR2(30),
  ASM_GENDER                   VARCHAR2(50),
  ASM_JOB_DESC                 VARCHAR2(500),
  ASM_JOB_NAME                 VARCHAR2(500),
  ASM_LEFT_HANDED              VARCHAR2(50),
  ASM_PARTICIPANT_STATUS       VARCHAR2(50),
  FAC_RETRO_ID                 VARCHAR2(50),
  FLAG_OLD_DATA                VARCHAR2(10),
  PERCENT_OFFERED              VARCHAR2(50),
  SHARE_OFFERED                NUMBER,
  TSI_NUSA_RE                  NUMBER,
  TSI_SPREADED                 NUMBER
);

CREATE TABLE T_FR_SHIP (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  AGE                          VARCHAR2(50),
  DWT                          VARCHAR2(50),
  FLAG_NAME                    VARCHAR2(500),
  GRT                          VARCHAR2(50),
  NM_SHIP                      VARCHAR2(50),
  NRT                          VARCHAR2(50),
  REGISTER                     VARCHAR2(50),
  SHIPTYPE_ID                  VARCHAR2(50),
  SHIPTYPE_NAME                VARCHAR2(500),
  SHIP_REF_ID                  VARCHAR2(50),
  Y_MAKE1                      VARCHAR2(50),
  Y_MAKE2                      VARCHAR2(50)
);

CREATE TABLE T_FR_MARINEHULL (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  CONSTRUCTION                 VARCHAR2(50),
  DWTRT                        VARCHAR2(50),
  GRTRT                        VARCHAR2(50),
  NRTRT                        VARCHAR2(50),
  TYPE_OF_VESSEL               VARCHAR2(50)
);

CREATE TABLE T_FR_LISTCAUSEOFLOSS (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  CLAIM                        NUMBER,
  CURRENCY                     VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  DETAIL                       VARCHAR2(50),
  REMARKS                      VARCHAR2(500)
);

CREATE TABLE T_FR_LOCATIONLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  LOSS_RATIO1_YEAR_AMOUNT      NUMBER,
  LOSS_RATIO1_YEAR_PERCENT     VARCHAR2(50),
  LOSS_RATIO35_YEAR_AMOUNT     NUMBER,
  LOSS_RATIO35_YEAR_PERCENT    VARCHAR2(50)
);

CREATE TABLE T_FR_BUILDINGCONSTRUCTION (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  FLOOR_TYPE                   VARCHAR2(50),
  ROOF_TYPE                    VARCHAR2(50),
  WALL_TYPE                    VARCHAR2(50)
);

CREATE TABLE T_FR_TABLEOFLIMIT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  CATEGORY                     VARCHAR2(50),
  DESCRIPTION                  VARCHAR2(500),
  PCT_LIMIT                    NUMBER
);

CREATE TABLE T_FR_TOTALTSIPREMIRETRO (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL,
  NAME                         VARCHAR2(500),
  PREMIUM                      NUMBER,
  TSI                          NUMBER,
  CURRENCY_CODE                VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL
);

CREATE TABLE T_FR_SURROUNDINGRISK (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  FLOOD_AREA_STATUS            VARCHAR2(50),
  HOUSEKEEPING_STATUS          VARCHAR2(50)
);

CREATE TABLE T_FR_COINSDATA (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  COINS_NAME                   VARCHAR2(500)
);

CREATE TABLE T_FR_LC (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  LC                           VARCHAR2(50)
);

CREATE TABLE T_FR_FACOFFERLIST (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL
);

CREATE TABLE T_FR_OBJECT (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL
);

CREATE TABLE T_FR_POLICY (
  ID                           NUMBER NOT NULL,
  IDPEGA                       VARCHAR2(50),
  COB_GROUP                    VARCHAR2(20),
  PARENT_ID                    NUMBER NOT NULL,
  PARENT_TABLE                 VARCHAR2(30),
  SRC_PATH                     VARCHAR2(200),
  SEQ_NO                       NUMBER(5) NOT NULL,
  ROW_UID                      VARCHAR2(36) NOT NULL
);

-- ---------------------------------------------------------------------------
-- 2. PRIMARY KEY  (ID, surrogate)
-- ---------------------------------------------------------------------------

ALTER TABLE T_WORK_POLIS ADD CONSTRAINT PK_WORK_POLIS PRIMARY KEY (ID);
ALTER TABLE T_GENERAL_POLIS ADD CONSTRAINT PK_GENERAL_POLIS PRIMARY KEY (ID);
ALTER TABLE T_COVERAGELIST ADD CONSTRAINT PK_COVERAGELIST PRIMARY KEY (ID);
ALTER TABLE T_QUOTATIONDATA ADD CONSTRAINT PK_QUOTATIONDATA PRIMARY KEY (ID);
ALTER TABLE T_LISTINSTALLMENT ADD CONSTRAINT PK_LISTINSTALLMENT PRIMARY KEY (ID);
ALTER TABLE T_CURRENCYLIST ADD CONSTRAINT PK_CURRENCYLIST PRIMARY KEY (ID);
ALTER TABLE T_SPREADINGLIST ADD CONSTRAINT PK_SPREADINGLIST PRIMARY KEY (ID);
ALTER TABLE T_CEDINGCOLIST ADD CONSTRAINT PK_CEDINGCOLIST PRIMARY KEY (ID);
ALTER TABLE T_CEDINGCEDANTLIST ADD CONSTRAINT PK_CEDINGCEDANTLIST PRIMARY KEY (ID);
ALTER TABLE T_CEDING_CURRENCYLIST ADD CONSTRAINT PK_CEDING_CURRENCYLIST PRIMARY KEY (ID);
ALTER TABLE T_CURRENCY ADD CONSTRAINT PK_CURRENCY PRIMARY KEY (ID);
ALTER TABLE T_DEDUCTIBLELIST ADD CONSTRAINT PK_DEDUCTIBLELIST PRIMARY KEY (ID);
ALTER TABLE T_FACRETRODETAILS ADD CONSTRAINT PK_FACRETRODETAILS PRIMARY KEY (ID);
ALTER TABLE T_OCCUPATIONLIST ADD CONSTRAINT PK_OCCUPATIONLIST PRIMARY KEY (ID);
ALTER TABLE T_PROPERTY ADD CONSTRAINT PK_PROPERTY PRIMARY KEY (ID);
ALTER TABLE T_PERSONLIST ADD CONSTRAINT PK_PERSONLIST PRIMARY KEY (ID);
ALTER TABLE T_LOCATIONLIST ADD CONSTRAINT PK_LOCATIONLIST PRIMARY KEY (ID);
ALTER TABLE T_RISKLOCATION ADD CONSTRAINT PK_RISKLOCATION PRIMARY KEY (ID);
ALTER TABLE T_PROPERTYITEMLIST ADD CONSTRAINT PK_PROPERTYITEMLIST PRIMARY KEY (ID);
ALTER TABLE T_ANEKALIST ADD CONSTRAINT PK_ANEKALIST PRIMARY KEY (ID);
ALTER TABLE T_VEHICLELIST ADD CONSTRAINT PK_VEHICLELIST PRIMARY KEY (ID);
ALTER TABLE T_RETROLIST ADD CONSTRAINT PK_RETROLIST PRIMARY KEY (ID);
ALTER TABLE T_CARGOLIST ADD CONSTRAINT PK_CARGOLIST PRIMARY KEY (ID);
ALTER TABLE T_SCORINGRESULT ADD CONSTRAINT PK_SCORINGRESULT PRIMARY KEY (ID);
ALTER TABLE T_SHIP ADD CONSTRAINT PK_SHIP PRIMARY KEY (ID);
ALTER TABLE T_ADDITIONALCOVERAGE ADD CONSTRAINT PK_ADDITIONALCOVERAGE PRIMARY KEY (ID);
ALTER TABLE T_COVERAGEDATALIST ADD CONSTRAINT PK_COVERAGEDATALIST PRIMARY KEY (ID);
ALTER TABLE T_AVIATIONHULL ADD CONSTRAINT PK_AVIATIONHULL PRIMARY KEY (ID);
ALTER TABLE T_MARINEHULL ADD CONSTRAINT PK_MARINEHULL PRIMARY KEY (ID);
ALTER TABLE T_SCORING_FACTOR ADD CONSTRAINT PK_SCORING_FACTOR PRIMARY KEY (ID);
ALTER TABLE T_VEHICLEHE ADD CONSTRAINT PK_VEHICLEHE PRIMARY KEY (ID);
ALTER TABLE T_ACCESSORYLIST ADD CONSTRAINT PK_ACCESSORYLIST PRIMARY KEY (ID);
ALTER TABLE T_BUILDINGCONSTRUCTION ADD CONSTRAINT PK_BUILDINGCONSTRUCTION PRIMARY KEY (ID);
ALTER TABLE T_DATASCORINGRISKLIST ADD CONSTRAINT PK_DATASCORINGRISKLIST PRIMARY KEY (ID);
ALTER TABLE T_LISTCAUSEOFLOSS ADD CONSTRAINT PK_LISTCAUSEOFLOSS PRIMARY KEY (ID);
ALTER TABLE T_SURVEYREPORTLIST ADD CONSTRAINT PK_SURVEYREPORTLIST PRIMARY KEY (ID);
ALTER TABLE T_SCORINGRISK ADD CONSTRAINT PK_SCORINGRISK PRIMARY KEY (ID);
ALTER TABLE T_SCORING_OPTION ADD CONSTRAINT PK_SCORING_OPTION PRIMARY KEY (ID);
ALTER TABLE T_TABLEOFLIMIT ADD CONSTRAINT PK_TABLEOFLIMIT PRIMARY KEY (ID);
ALTER TABLE T_ASMHEIR ADD CONSTRAINT PK_ASMHEIR PRIMARY KEY (ID);
ALTER TABLE T_SURROUNDINGRISK ADD CONSTRAINT PK_SURROUNDINGRISK PRIMARY KEY (ID);
ALTER TABLE T_COINSDATA ADD CONSTRAINT PK_COINSDATA PRIMARY KEY (ID);
ALTER TABLE T_LC ADD CONSTRAINT PK_LC PRIMARY KEY (ID);
ALTER TABLE T_QUOTATIONLIST ADD CONSTRAINT PK_QUOTATIONLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_COVERAGELIST ADD CONSTRAINT PK_FR_COVERAGELIST PRIMARY KEY (ID);
ALTER TABLE T_FR_CURRENCYLIST ADD CONSTRAINT PK_FR_CURRENCYLIST PRIMARY KEY (ID);
ALTER TABLE T_FACRETROLIST ADD CONSTRAINT PK_FACRETROLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_SPREADINGLIST ADD CONSTRAINT PK_FR_SPREADINGLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_FACOUTOBJECTLIST ADD CONSTRAINT PK_FR_FACOUTOBJECTLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_NETPERCURRENCY ADD CONSTRAINT PK_FR_NETPERCURRENCY PRIMARY KEY (ID);
ALTER TABLE T_FR_PRINTRISLIP ADD CONSTRAINT PK_FR_PRINTRISLIP PRIMARY KEY (ID);
ALTER TABLE T_FR_DEDUCTIBLELIST ADD CONSTRAINT PK_FR_DEDUCTIBLELIST PRIMARY KEY (ID);
ALTER TABLE T_FR_WPC ADD CONSTRAINT PK_FR_WPC PRIMARY KEY (ID);
ALTER TABLE T_FR_CURRENCY ADD CONSTRAINT PK_FR_CURRENCY PRIMARY KEY (ID);
ALTER TABLE T_FR_PROPERTY ADD CONSTRAINT PK_FR_PROPERTY PRIMARY KEY (ID);
ALTER TABLE T_FR_LISTINSTALLMENT ADD CONSTRAINT PK_FR_LISTINSTALLMENT PRIMARY KEY (ID);
ALTER TABLE T_FR_PAYMENT ADD CONSTRAINT PK_FR_PAYMENT PRIMARY KEY (ID);
ALTER TABLE T_FR_SECURITYREINSURER ADD CONSTRAINT PK_FR_SECURITYREINSURER PRIMARY KEY (ID);
ALTER TABLE T_FR_OFFEREDPAYMENT ADD CONSTRAINT PK_FR_OFFEREDPAYMENT PRIMARY KEY (ID);
ALTER TABLE T_FR_RISKLOCATION ADD CONSTRAINT PK_FR_RISKLOCATION PRIMARY KEY (ID);
ALTER TABLE T_FR_OCCUPATIONLIST ADD CONSTRAINT PK_FR_OCCUPATIONLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_CARGOLIST ADD CONSTRAINT PK_FR_CARGOLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_ANEKALIST ADD CONSTRAINT PK_FR_ANEKALIST PRIMARY KEY (ID);
ALTER TABLE T_FR_PROPERTYITEMLIST ADD CONSTRAINT PK_FR_PROPERTYITEMLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_PERSONLIST ADD CONSTRAINT PK_FR_PERSONLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_SHIP ADD CONSTRAINT PK_FR_SHIP PRIMARY KEY (ID);
ALTER TABLE T_FR_MARINEHULL ADD CONSTRAINT PK_FR_MARINEHULL PRIMARY KEY (ID);
ALTER TABLE T_FR_LISTCAUSEOFLOSS ADD CONSTRAINT PK_FR_LISTCAUSEOFLOSS PRIMARY KEY (ID);
ALTER TABLE T_FR_LOCATIONLIST ADD CONSTRAINT PK_FR_LOCATIONLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_BUILDINGCONSTRUCTION ADD CONSTRAINT PK_FR_BUILDINGCONSTRUCTION PRIMARY KEY (ID);
ALTER TABLE T_FR_TABLEOFLIMIT ADD CONSTRAINT PK_FR_TABLEOFLIMIT PRIMARY KEY (ID);
ALTER TABLE T_FR_TOTALTSIPREMIRETRO ADD CONSTRAINT PK_FR_TOTALTSIPREMIRETRO PRIMARY KEY (ID);
ALTER TABLE T_FR_SURROUNDINGRISK ADD CONSTRAINT PK_FR_SURROUNDINGRISK PRIMARY KEY (ID);
ALTER TABLE T_FR_COINSDATA ADD CONSTRAINT PK_FR_COINSDATA PRIMARY KEY (ID);
ALTER TABLE T_FR_LC ADD CONSTRAINT PK_FR_LC PRIMARY KEY (ID);
ALTER TABLE T_FR_FACOFFERLIST ADD CONSTRAINT PK_FR_FACOFFERLIST PRIMARY KEY (ID);
ALTER TABLE T_FR_OBJECT ADD CONSTRAINT PK_FR_OBJECT PRIMARY KEY (ID);
ALTER TABLE T_FR_POLICY ADD CONSTRAINT PK_FR_POLICY PRIMARY KEY (ID);

-- ---------------------------------------------------------------------------
-- 3. FOREIGN KEY  (hanya bila anak punya TEPAT SATU induk, dan bukan ditegakkan di Go)
-- ---------------------------------------------------------------------------

ALTER TABLE T_ACCESSORYLIST ADD CONSTRAINT FK_ACCESSORYLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_VEHICLELIST (ID) ON DELETE CASCADE;
ALTER TABLE T_ADDITIONALCOVERAGE ADD CONSTRAINT FK_ADDITIONALCOVERAGE_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_COVERAGELIST (ID) ON DELETE CASCADE;
ALTER TABLE T_ASMHEIR ADD CONSTRAINT FK_ASMHEIR_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_PERSONLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_AVIATIONHULL ADD CONSTRAINT FK_AVIATIONHULL_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_ANEKALIST (ID) ON DELETE CASCADE;
ALTER TABLE T_BUILDINGCONSTRUCTION ADD CONSTRAINT FK_BUILDINGCONSTRUCTION_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_CARGOLIST ADD CONSTRAINT FK_CARGOLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_CEDING_CURRENCYLIST ADD CONSTRAINT FK_CEDING_CURRENCYLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_CEDINGCEDANTLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_CEDINGCEDANTLIST ADD CONSTRAINT FK_CEDINGCEDANTLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_CEDINGCOLIST ADD CONSTRAINT FK_CEDINGCOLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_QUOTATIONDATA (ID) ON DELETE CASCADE;
ALTER TABLE T_COINSDATA ADD CONSTRAINT FK_COINSDATA_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_LISTCAUSEOFLOSS (ID) ON DELETE CASCADE;
ALTER TABLE T_COVERAGEDATALIST ADD CONSTRAINT FK_COVERAGEDATALIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_PERSONLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_CURRENCYLIST ADD CONSTRAINT FK_CURRENCYLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_DATASCORINGRISKLIST ADD CONSTRAINT FK_DATASCORINGRISKLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_SCORINGRISK (ID) ON DELETE CASCADE;
ALTER TABLE T_DEDUCTIBLELIST ADD CONSTRAINT FK_DEDUCTIBLELIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_COVERAGELIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FACRETRODETAILS ADD CONSTRAINT FK_FACRETRODETAILS_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_FACRETROLIST ADD CONSTRAINT FK_FACRETROLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_ANEKALIST ADD CONSTRAINT FK_FR_ANEKALIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_OCCUPATIONLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_BUILDINGCONSTRUCTION ADD CONSTRAINT FK_FR_BUILDINGCONSTRUCTION_PAR FOREIGN KEY (PARENT_ID) REFERENCES T_FR_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_COINSDATA ADD CONSTRAINT FK_FR_COINSDATA_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_LISTCAUSEOFLOSS (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_DEDUCTIBLELIST ADD CONSTRAINT FK_FR_DEDUCTIBLELIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_COVERAGELIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_FACOFFERLIST ADD CONSTRAINT FK_FR_FACOFFERLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_PRINTRISLIP (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_FACOUTOBJECTLIST ADD CONSTRAINT FK_FR_FACOUTOBJECTLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_COVERAGELIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_LC ADD CONSTRAINT FK_FR_LC_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_CARGOLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_LISTCAUSEOFLOSS ADD CONSTRAINT FK_FR_LISTCAUSEOFLOSS_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_LOCATIONLIST ADD CONSTRAINT FK_FR_LOCATIONLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FACRETROLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_MARINEHULL ADD CONSTRAINT FK_FR_MARINEHULL_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_ANEKALIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_NETPERCURRENCY ADD CONSTRAINT FK_FR_NETPERCURRENCY_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_FACOFFERLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_OBJECT ADD CONSTRAINT FK_FR_OBJECT_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_FACOFFERLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_OFFEREDPAYMENT ADD CONSTRAINT FK_FR_OFFEREDPAYMENT_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_FACOFFERLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_PAYMENT ADD CONSTRAINT FK_FR_PAYMENT_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_POLICY (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_PERSONLIST ADD CONSTRAINT FK_FR_PERSONLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FACRETROLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_PRINTRISLIP ADD CONSTRAINT FK_FR_PRINTRISLIP_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FACRETROLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_PROPERTY ADD CONSTRAINT FK_FR_PROPERTY_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_LOCATIONLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_PROPERTYITEMLIST ADD CONSTRAINT FK_FR_PROPERTYITEMLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_SECURITYREINSURER ADD CONSTRAINT FK_FR_SECURITYREINSURER_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_FACOFFERLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_SHIP ADD CONSTRAINT FK_FR_SHIP_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_CARGOLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_SPREADINGLIST ADD CONSTRAINT FK_FR_SPREADINGLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_COVERAGELIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_SURROUNDINGRISK ADD CONSTRAINT FK_FR_SURROUNDINGRISK_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_TABLEOFLIMIT ADD CONSTRAINT FK_FR_TABLEOFLIMIT_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_FR_OCCUPATIONLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_FR_TOTALTSIPREMIRETRO ADD CONSTRAINT FK_FR_TOTALTSIPREMIRETRO_PAREN FOREIGN KEY (PARENT_ID) REFERENCES T_FR_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_GENERAL_POLIS ADD CONSTRAINT FK_GENERAL_POLIS_IDPEGA FOREIGN KEY (IDPEGA) REFERENCES T_WORK_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_LC ADD CONSTRAINT FK_LC_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_CARGOLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_LISTCAUSEOFLOSS ADD CONSTRAINT FK_LISTCAUSEOFLOSS_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_LISTINSTALLMENT ADD CONSTRAINT FK_LISTINSTALLMENT_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_CURRENCYLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_LOCATIONLIST ADD CONSTRAINT FK_LOCATIONLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_MARINEHULL ADD CONSTRAINT FK_MARINEHULL_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_ANEKALIST (ID) ON DELETE CASCADE;
ALTER TABLE T_PERSONLIST ADD CONSTRAINT FK_PERSONLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_PROPERTY ADD CONSTRAINT FK_PROPERTY_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_LOCATIONLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_PROPERTYITEMLIST ADD CONSTRAINT FK_PROPERTYITEMLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_QUOTATIONDATA ADD CONSTRAINT FK_QUOTATIONDATA_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_QUOTATIONLIST ADD CONSTRAINT FK_QUOTATIONLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_QUOTATIONDATA (ID) ON DELETE CASCADE;
ALTER TABLE T_RETROLIST ADD CONSTRAINT FK_RETROLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_SPREADINGLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_RISKLOCATION ADD CONSTRAINT FK_RISKLOCATION_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_SCORING_FACTOR ADD CONSTRAINT FK_SCORING_FACTOR_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_DATASCORINGRISKLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_SCORING_OPTION ADD CONSTRAINT FK_SCORING_OPTION_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_SCORING_FACTOR (ID) ON DELETE CASCADE;
ALTER TABLE T_SCORINGRESULT ADD CONSTRAINT FK_SCORINGRESULT_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_SCORINGRISK (ID) ON DELETE CASCADE;
ALTER TABLE T_SCORINGRISK ADD CONSTRAINT FK_SCORINGRISK_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;
ALTER TABLE T_SHIP ADD CONSTRAINT FK_SHIP_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_CARGOLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_SPREADINGLIST ADD CONSTRAINT FK_SPREADINGLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_COVERAGELIST (ID) ON DELETE CASCADE;
ALTER TABLE T_SURROUNDINGRISK ADD CONSTRAINT FK_SURROUNDINGRISK_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_PROPERTY (ID) ON DELETE CASCADE;
ALTER TABLE T_SURVEYREPORTLIST ADD CONSTRAINT FK_SURVEYREPORTLIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_QUOTATIONDATA (ID) ON DELETE CASCADE;
ALTER TABLE T_TABLEOFLIMIT ADD CONSTRAINT FK_TABLEOFLIMIT_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_OCCUPATIONLIST (ID) ON DELETE CASCADE;
ALTER TABLE T_VEHICLEHE ADD CONSTRAINT FK_VEHICLEHE_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_ANEKALIST (ID) ON DELETE CASCADE;
ALTER TABLE T_VEHICLELIST ADD CONSTRAINT FK_VEHICLELIST_PARENT FOREIGN KEY (PARENT_ID) REFERENCES T_GENERAL_POLIS (ID) ON DELETE CASCADE;

-- Anak berinduk GANDA: satu kolom menunjuk lebih dari satu tabel, jadi tidak dapat
-- diberi FOREIGN KEY. Keutuhannya di Go; kolom PARENT_TABLE menyimpan tabel mana (V-23).
-- T_ANEKALIST.PARENT_ID menunjuk 2 induk: T_OCCUPATIONLIST, T_RISKLOCATION
-- T_COVERAGELIST.PARENT_ID menunjuk 5 induk: T_ANEKALIST, T_CARGOLIST, T_PERSONLIST, T_PROPERTYITEMLIST, T_VEHICLELIST
-- T_CURRENCY.PARENT_ID menunjuk 3 induk: T_ANEKALIST, T_COVERAGELIST, T_PERSONLIST
-- T_FR_CARGOLIST.PARENT_ID menunjuk 2 induk: T_FACRETROLIST, T_FR_FACOFFERLIST
-- T_FR_COVERAGELIST.PARENT_ID menunjuk 4 induk: T_FR_ANEKALIST, T_FR_CARGOLIST, T_FR_PERSONLIST, T_FR_PROPERTYITEMLIST
-- T_FR_CURRENCY.PARENT_ID menunjuk 2 induk: T_FR_ANEKALIST, T_FR_COVERAGELIST
-- T_FR_CURRENCYLIST.PARENT_ID menunjuk 2 induk: T_FACRETROLIST, T_FR_FACOFFERLIST
-- T_FR_LISTINSTALLMENT.PARENT_ID menunjuk 2 induk: T_FR_CURRENCYLIST, T_FR_PAYMENT
-- T_FR_OCCUPATIONLIST.PARENT_ID menunjuk 3 induk: T_FR_OBJECT, T_FR_PROPERTY, T_FR_RISKLOCATION
-- T_FR_RISKLOCATION.PARENT_ID menunjuk 2 induk: T_FR_OBJECT, T_FR_PROPERTY
-- T_FR_WPC.PARENT_ID menunjuk 2 induk: T_FR_FACOFFERLIST, T_FR_PRINTRISLIP
-- T_OCCUPATIONLIST.PARENT_ID menunjuk 3 induk: T_PROPERTY, T_RISKLOCATION, T_VEHICLELIST

-- Ditegakkan di Go atas keputusan, bukan karena induk ganda:
-- T_GENERAL_POLIS.OLD_POLIS_ID -> T_GENERAL_POLIS : keutuhan ditegakkan di Go, TANPA REFERENCES

-- ---------------------------------------------------------------------------
-- 4. INDEX pada setiap kunci tamu
-- ---------------------------------------------------------------------------

CREATE INDEX IX_ACCESSORYLIST_PARENT ON T_ACCESSORYLIST (PARENT_ID);
CREATE INDEX IX_ADDITIONALCOVERAGE_PARENT ON T_ADDITIONALCOVERAGE (PARENT_ID);
CREATE INDEX IX_ANEKALIST_PARENT ON T_ANEKALIST (PARENT_ID);
CREATE INDEX IX_ASMHEIR_PARENT ON T_ASMHEIR (PARENT_ID);
CREATE INDEX IX_AVIATIONHULL_PARENT ON T_AVIATIONHULL (PARENT_ID);
CREATE INDEX IX_BUILDINGCONSTRUCTION_PARENT ON T_BUILDINGCONSTRUCTION (PARENT_ID);
CREATE INDEX IX_CARGOLIST_PARENT ON T_CARGOLIST (PARENT_ID);
CREATE INDEX IX_CEDING_CURRENCYLIST_PARENT ON T_CEDING_CURRENCYLIST (PARENT_ID);
CREATE INDEX IX_CEDINGCEDANTLIST_PARENT ON T_CEDINGCEDANTLIST (PARENT_ID);
CREATE INDEX IX_CEDINGCOLIST_PARENT ON T_CEDINGCOLIST (PARENT_ID);
CREATE INDEX IX_COINSDATA_PARENT ON T_COINSDATA (PARENT_ID);
CREATE INDEX IX_COVERAGEDATALIST_PARENT ON T_COVERAGEDATALIST (PARENT_ID);
CREATE INDEX IX_COVERAGELIST_PARENT ON T_COVERAGELIST (PARENT_ID);
CREATE INDEX IX_CURRENCY_PARENT ON T_CURRENCY (PARENT_ID);
CREATE INDEX IX_CURRENCYLIST_PARENT ON T_CURRENCYLIST (PARENT_ID);
CREATE INDEX IX_DATASCORINGRISKLIST_PARENT ON T_DATASCORINGRISKLIST (PARENT_ID);
CREATE INDEX IX_DEDUCTIBLELIST_PARENT ON T_DEDUCTIBLELIST (PARENT_ID);
CREATE INDEX IX_FACRETRODETAILS_PARENT ON T_FACRETRODETAILS (PARENT_ID);
CREATE INDEX IX_FACRETROLIST_PARENT ON T_FACRETROLIST (PARENT_ID);
CREATE INDEX IX_FR_ANEKALIST_PARENT ON T_FR_ANEKALIST (PARENT_ID);
CREATE INDEX IX_FR_BUILDINGCONSTRUCTION_PAR ON T_FR_BUILDINGCONSTRUCTION (PARENT_ID);
CREATE INDEX IX_FR_CARGOLIST_PARENT ON T_FR_CARGOLIST (PARENT_ID);
CREATE INDEX IX_FR_COINSDATA_PARENT ON T_FR_COINSDATA (PARENT_ID);
CREATE INDEX IX_FR_COVERAGELIST_PARENT ON T_FR_COVERAGELIST (PARENT_ID);
CREATE INDEX IX_FR_CURRENCY_PARENT ON T_FR_CURRENCY (PARENT_ID);
CREATE INDEX IX_FR_CURRENCYLIST_PARENT ON T_FR_CURRENCYLIST (PARENT_ID);
CREATE INDEX IX_FR_DEDUCTIBLELIST_PARENT ON T_FR_DEDUCTIBLELIST (PARENT_ID);
CREATE INDEX IX_FR_FACOFFERLIST_PARENT ON T_FR_FACOFFERLIST (PARENT_ID);
CREATE INDEX IX_FR_FACOUTOBJECTLIST_PARENT ON T_FR_FACOUTOBJECTLIST (PARENT_ID);
CREATE INDEX IX_FR_LC_PARENT ON T_FR_LC (PARENT_ID);
CREATE INDEX IX_FR_LISTCAUSEOFLOSS_PARENT ON T_FR_LISTCAUSEOFLOSS (PARENT_ID);
CREATE INDEX IX_FR_LISTINSTALLMENT_PARENT ON T_FR_LISTINSTALLMENT (PARENT_ID);
CREATE INDEX IX_FR_LOCATIONLIST_PARENT ON T_FR_LOCATIONLIST (PARENT_ID);
CREATE INDEX IX_FR_MARINEHULL_PARENT ON T_FR_MARINEHULL (PARENT_ID);
CREATE INDEX IX_FR_NETPERCURRENCY_PARENT ON T_FR_NETPERCURRENCY (PARENT_ID);
CREATE INDEX IX_FR_OBJECT_PARENT ON T_FR_OBJECT (PARENT_ID);
CREATE INDEX IX_FR_OCCUPATIONLIST_PARENT ON T_FR_OCCUPATIONLIST (PARENT_ID);
CREATE INDEX IX_FR_OFFEREDPAYMENT_PARENT ON T_FR_OFFEREDPAYMENT (PARENT_ID);
CREATE INDEX IX_FR_PAYMENT_PARENT ON T_FR_PAYMENT (PARENT_ID);
CREATE INDEX IX_FR_PERSONLIST_PARENT ON T_FR_PERSONLIST (PARENT_ID);
CREATE INDEX IX_FR_PRINTRISLIP_PARENT ON T_FR_PRINTRISLIP (PARENT_ID);
CREATE INDEX IX_FR_PROPERTY_PARENT ON T_FR_PROPERTY (PARENT_ID);
CREATE INDEX IX_FR_PROPERTYITEMLIST_PARENT ON T_FR_PROPERTYITEMLIST (PARENT_ID);
CREATE INDEX IX_FR_RISKLOCATION_PARENT ON T_FR_RISKLOCATION (PARENT_ID);
CREATE INDEX IX_FR_SECURITYREINSURER_PARENT ON T_FR_SECURITYREINSURER (PARENT_ID);
CREATE INDEX IX_FR_SHIP_PARENT ON T_FR_SHIP (PARENT_ID);
CREATE INDEX IX_FR_SPREADINGLIST_PARENT ON T_FR_SPREADINGLIST (PARENT_ID);
CREATE INDEX IX_FR_SURROUNDINGRISK_PARENT ON T_FR_SURROUNDINGRISK (PARENT_ID);
CREATE INDEX IX_FR_TABLEOFLIMIT_PARENT ON T_FR_TABLEOFLIMIT (PARENT_ID);
CREATE INDEX IX_FR_TOTALTSIPREMIRETRO_PAREN ON T_FR_TOTALTSIPREMIRETRO (PARENT_ID);
CREATE INDEX IX_FR_WPC_PARENT ON T_FR_WPC (PARENT_ID);
CREATE UNIQUE INDEX UX_GENERAL_POLIS_IDPEGA ON T_GENERAL_POLIS (IDPEGA);
CREATE UNIQUE INDEX UX_GENERAL_POLIS_OLD_POLIS ON T_GENERAL_POLIS (OLD_POLIS_ID);
CREATE INDEX IX_LC_PARENT ON T_LC (PARENT_ID);
CREATE INDEX IX_LISTCAUSEOFLOSS_PARENT ON T_LISTCAUSEOFLOSS (PARENT_ID);
CREATE INDEX IX_LISTINSTALLMENT_PARENT ON T_LISTINSTALLMENT (PARENT_ID);
CREATE INDEX IX_LOCATIONLIST_PARENT ON T_LOCATIONLIST (PARENT_ID);
CREATE INDEX IX_MARINEHULL_PARENT ON T_MARINEHULL (PARENT_ID);
CREATE INDEX IX_OCCUPATIONLIST_PARENT ON T_OCCUPATIONLIST (PARENT_ID);
CREATE INDEX IX_PERSONLIST_PARENT ON T_PERSONLIST (PARENT_ID);
CREATE INDEX IX_PROPERTY_PARENT ON T_PROPERTY (PARENT_ID);
CREATE INDEX IX_PROPERTYITEMLIST_PARENT ON T_PROPERTYITEMLIST (PARENT_ID);
CREATE INDEX IX_QUOTATIONDATA_PARENT ON T_QUOTATIONDATA (PARENT_ID);
CREATE INDEX IX_QUOTATIONLIST_PARENT ON T_QUOTATIONLIST (PARENT_ID);
CREATE INDEX IX_RETROLIST_PARENT ON T_RETROLIST (PARENT_ID);
CREATE INDEX IX_RISKLOCATION_PARENT ON T_RISKLOCATION (PARENT_ID);
CREATE INDEX IX_SCORING_FACTOR_PARENT ON T_SCORING_FACTOR (PARENT_ID);
CREATE INDEX IX_SCORING_OPTION_PARENT ON T_SCORING_OPTION (PARENT_ID);
CREATE INDEX IX_SCORINGRESULT_PARENT ON T_SCORINGRESULT (PARENT_ID);
CREATE INDEX IX_SCORINGRISK_PARENT ON T_SCORINGRISK (PARENT_ID);
CREATE INDEX IX_SHIP_PARENT ON T_SHIP (PARENT_ID);
CREATE INDEX IX_SPREADINGLIST_PARENT ON T_SPREADINGLIST (PARENT_ID);
CREATE INDEX IX_SURROUNDINGRISK_PARENT ON T_SURROUNDINGRISK (PARENT_ID);
CREATE INDEX IX_SURVEYREPORTLIST_PARENT ON T_SURVEYREPORTLIST (PARENT_ID);
CREATE INDEX IX_TABLEOFLIMIT_PARENT ON T_TABLEOFLIMIT (PARENT_ID);
CREATE INDEX IX_VEHICLEHE_PARENT ON T_VEHICLEHE (PARENT_ID);
CREATE INDEX IX_VEHICLELIST_PARENT ON T_VEHICLELIST (PARENT_ID);

-- ============================================================================
-- SELESAI
-- ============================================================================

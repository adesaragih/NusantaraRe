-- 253 - kolom rancangan T_QUOTATIONDATA yang belum ada di DB (tiket 76b, keputusan work owner 12-10-2026 "Ya, bertahap": salinan
-- nama / kode / status QuotationData seperti Pega; sistem baru menyimpan rujukan di T_NB_OPPORTUNITY dan membaca nama dari master).
--
-- `[terverifikasi]` sensus dua cara (fungsi loader `tipeAkhirMigrasi` + pengurai python berkas migrasi urut posisi pernyataan, cocok):
-- T_QUOTATIONDATA fisik = 16 kolom (183 / 184 / 185 / 199); rancangan `loader/skema_gen.go` 67, nol amandemen; belum ada = 51.
-- Sebelas kolom khusus Endorsement / Renewal TIDAK ditambah (EDM_CHARGE_FEE, EDM_DATE, EDM_NOTE, EDM_SOURCE, EDM_SOURCE_NOTE,
-- EDM_STATUS, EDM_SURVEY, EDM_TYPE, EDM_TYPE_NEW, ENDORSEMENT_INTERNAL_RETRO, RNW_DATE) -> 40 kolom, seluruhnya VARCHAR2 rancangan
-- (nama / tipe / lebar = rancangan, tanpa NUMBER). Semua boleh NULL: data lama tetap sah. Kolom tahap 1 diisi aplikasi (Save
-- Opportunity / General); sisanya belum diisi (peta di docs/issues/76b-quotationdata-lengkap.md).
-- T_QUOTATIONDATA lahir di 183 -> 253 aman menurut urutan nama. Rentang KETIGA nbfacin 250-254. Ditulis, TIDAK dijalankan agent;
-- DBA yang menjalankan. Nol COMMIT (ADR-U-0029). Biner tanpa 253 tetap jalan: penulisan kolom 253 menoleransi ORA-00904 (log).
ALTER TABLE {skema}.T_QUOTATIONDATA ADD (
  BRANCH_CODE      VARCHAR2(50),
  BRANCH_NAME      VARCHAR2(500),
  BUSINESS_CODE    VARCHAR2(50),
  BUSINESS_FAC     VARCHAR2(50),
  BUSINESS_NAME    VARCHAR2(500),
  BUSINESS_OLD_ID  VARCHAR2(50),
  BUSINESS_TYPE    VARCHAR2(50),
  BUSINESS_TYPE2   VARCHAR2(50),
  COMMENT_TEXT     VARCHAR2(500),
  COPY_FROM        VARCHAR2(50),
  EMAIL            VARCHAR2(50),
  FAC_OUT_STATUS   VARCHAR2(50),
  FRONTING_STATUS  VARCHAR2(50),
  GROUP_PANEL      VARCHAR2(50),
  INPUT_MONTH      VARCHAR2(50),
  INPUT_QUO_DATE   VARCHAR2(30),
  INSURED_ID       VARCHAR2(50),
  INSURED_NAME     VARCHAR2(500),
  IS_EDM_SHARE     VARCHAR2(10),
  IS_GROUP         VARCHAR2(10),
  IS_MOP           VARCHAR2(10),
  IS_PROPOSAL      VARCHAR2(10),
  IS_SURVEY_REPORT VARCHAR2(10),
  IS_TBA           VARCHAR2(10),
  IS_TEMPLATE      VARCHAR2(10),
  MARKETING_CODE   VARCHAR2(50),
  MARKETING_NAME   VARCHAR2(500),
  OFFERING_DATE    VARCHAR2(30),
  OLD_POLICY_NO    VARCHAR2(50),
  OPERATOR_ID      VARCHAR2(50),
  PERIOD           VARCHAR2(30),
  PERIOD_MM        VARCHAR2(30),
  QQID             VARCHAR2(50),
  SOB_LEADER0      VARCHAR2(50),
  SOB_LSG          VARCHAR2(50),
  STATUS_BUSINESS  VARCHAR2(50),
  STATUS_SYARIAH   VARCHAR2(50),
  TEAM_GROUP       VARCHAR2(50),
  TYPE             VARCHAR2(50),
  BTN_QUOTATION    VARCHAR2(50)
)
/

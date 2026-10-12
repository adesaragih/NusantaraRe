-- 252 - kolom rancangan T_GENERAL_POLIS yang belum ada di DB (keputusan work owner 12-10-2026: "ini udh jadi polis masa data yg
-- disimpan tidak lengkap", contoh NB-151029; BERTAHAP, mulai identitas polis).
--
-- `[terverifikasi]` sensus dua cara (fungsi loader `tipeAkhirMigrasi` + pengurai python berkas migrasi, urut posisi pernyataan):
-- T_GENERAL_POLIS fisik = 25 kolom (182 + 198 / 199 / 200 / 201 / 202 / 208 / 210 / 212 / 213 / 255); rancangan
-- `loader/skema_gen.go` 75 + amandemen 3 (NB_STATUS / DATE_TO_UW / SUBMIT_TO_UW) = 78; belum ada = 53. Lima kolom khusus
-- Endorsement TIDAK ditambah (ENDORSEMENT_NO, ENDORSMENT_REASON, PRORATE_START_EDM, PRORATE_EDM_END, ERRMSG) -> 48 kolom.
-- Nama / tipe / lebar = rancangan; NUMBER polos rancangan diberi presisi sah penjaga (TestNolNumberTanpaPresisi) menurut konvensi
-- modul `[keputusan agent]`: uang / persen / share NUMBER(38,8) (ADR-0016), identitas OLD_POLIS_ID NUMBER(19) (A92), kode
-- PRO_RATE_TYPE NUMBER(5) (A170). Semua boleh NULL: data lama tetap sah. Tahap 1 (POLICY_NO, PROD_DATE_TIME, PROD_KE) diisi
-- Selesai Policy; kolom lain belum diisi aplikasi (peta tahap 2 / 3 di docs/issues/76-general-polis-lengkap.md).
-- T_GENERAL_POLIS lahir di 182 -> 252 aman menurut urutan nama. Rentang KETIGA nbfacin 250-254. Ditulis, TIDAK dijalankan agent;
-- DBA yang menjalankan. Nol COMMIT (ADR-U-0029). Biner tanpa 252 tetap jalan: penulisan tahap 1 menoleransi ORA-00904 (log).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  OLD_POLIS_ID                NUMBER(19),
  ADDITIONAL_CAPITAL          VARCHAR2(50),
  BINDER_RNM                  VARCHAR2(50),
  CEDING_RETENTION            NUMBER(38,8),
  CEDING_RETENTION_NOMINAL    NUMBER(38,8),
  COMMENT_CEDING              VARCHAR2(500),
  CONDITION                   VARCHAR2(500),
  CURRENT_YEAR                VARCHAR2(50),
  DATE_OFFER                  VARCHAR2(30),
  FLAG_SAVE_FO                VARCHAR2(10),
  FLAG_SPECIAL_ACCEPTANCE     VARCHAR2(10),
  FLAG_SPREADING              VARCHAR2(10),
  FOLLOWING_NB                VARCHAR2(50),
  GENERATE_DATE_POLICY_NO     VARCHAR2(30),
  ID_FOLLOWING_NB             VARCHAR2(50),
  INPUT_COMMENT_DEDUCTIBLE    VARCHAR2(500),
  IS_B2_B                     VARCHAR2(10),
  IS_BANDING                  VARCHAR2(10),
  IS_EDIT_REFF_NUMBER         VARCHAR2(10),
  IS_FAC_RETRO                VARCHAR2(10),
  IS_GROSS_NET_SHOW           VARCHAR2(10),
  IS_INPUT_FAC_RETRO          VARCHAR2(10),
  IS_LIFE_PERIOD              VARCHAR2(10),
  IS_OCCUP_EXCEPTION          VARCHAR2(10),
  IS_PREFERRED_RISK           VARCHAR2(10),
  IS_PRO_RATE                 VARCHAR2(10),
  IS_RI_SLIP                  VARCHAR2(10),
  IS_SPECIAL_ACCEPTANCE       VARCHAR2(10),
  IS_TOP_RISK                 VARCHAR2(10),
  OFFICER_FAC_OUT             VARCHAR2(50),
  PARAM_CONDITION             VARCHAR2(500),
  POLICY_MASTER_ID_PEGA       VARCHAR2(50),
  POLICY_MASTER_NUMBER        VARCHAR2(50),
  POLICY_NO                   VARCHAR2(50),
  PRO_RATE_TYPE               NUMBER(5),
  PROD_DATE_TIME              VARCHAR2(30),
  PROD_KE                     NUMBER(5),
  PRODUCT_BRIGUNA             VARCHAR2(50),
  SAVE_FAC_OUT                VARCHAR2(50),
  SHARE_RNML                  NUMBER(38,8),
  SOURCE_ID                   VARCHAR2(50),
  SUM_TOTAL_TSI               NUMBER(38,8),
  TOTAL_PREMI_NUSA_RE         NUMBER(38,8),
  TOTAL_RI_COM_NUSA_RE        VARCHAR2(50),
  TOTAL_TSI_NUSA_RE           NUMBER(38,8),
  TOTAL_TSI_NUSA_RE_SPREADING NUMBER(38,8),
  TOTAL_TSI_TOP_RISK          NUMBER(38,8),
  TREATY_CAPACITY             NUMBER(38,8)
)
/

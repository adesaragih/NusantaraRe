-- 560 - kolom khas lini FACIN di T_GENERAL_CLAIM (header klaim, tabel bersama milik Claim Life).
--
-- Prompt Claim Fac In tahap 1 §6 butir 2 (pola 520 / 600): medan kepala klaim Fac In yang punya penulis di XML
-- (InputRegisterDetail, CopyNB_Act, ProteksiDataRegister_Act, SetDisable_ACT, CLaimFaceSheet_Act, InputEstimationPre)
-- disimpan sebagai kolom TAMBAHAN; kolom yang maknanya sama dengan kolom Claim Life / Claim Prop / Claim Non Prop
-- dipakai ulang. Kolom nullable; uang NUMBER(38,10). Daftar kolom = katalog `models/katalog_tabel.go`
-- (TabelHeaderKlaim, `baru(...)`).
-- NOL COMMIT (ADR-U-0029), nol MODIFY / DROP kolom yang sudah ada. -migrate dijalankan work owner.
ALTER TABLE {skema}.T_GENERAL_CLAIM ADD (
  PRODKE              VARCHAR2(8),
  QQ_NAME             VARCHAR2(500),
  CEDING_CO_NAME      VARCHAR2(500),
  SOB_NAME            VARCHAR2(500),
  COUNTRY             VARCHAR2(255),
  COUNTRY_ID          VARCHAR2(64),
  CURRENCY_VALUE_IDR  NUMBER(38,10),
  GROSS_ESTIMATE      NUMBER(38,10),
  CLAIM_ESTIMATE      NUMBER(38,10),
  EX_GRATIA           VARCHAR2(16),
  IS_ERROR            VARCHAR2(16),
  IS_REGISTER         VARCHAR2(16),
  IS_ESTIMATION       VARCHAR2(16),
  IS_ADJUSTMENT       VARCHAR2(16),
  IS_PIC_TRANSFER     VARCHAR2(16),
  IS_TREATY_OUT       VARCHAR2(16),
  PY_NOTE             VARCHAR2(255),
  START_DATE_REGISTER DATE,
  END_DATE_REGISTER   DATE
)
/

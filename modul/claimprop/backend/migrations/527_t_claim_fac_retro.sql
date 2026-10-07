-- 527 - T_CLAIM_FAC_RETRO <- pyWorkPage.ClaimData.FacRetroList (diisi dari TREATYREINSURER, SaveOutstanding_Act langkah 35).
-- AdjustmentList(n).FacRetroList adalah salinan daftar ini - tidak disimpan (STRUKTUR J1 butir 5).
-- TOTAL_ESTIMATION_REINS <- .TotalEstimasiReas: satu-satunya pengisi di ClaimData.FacRetroList adalah AddKomiteTreatyChild_ACT
-- 8.1 (penyerahan komite = OQ-CP-16) - kasus baru membiarkannya kosong; kasus lama dari JSON_KLAIM membawa nilainya.
CREATE TABLE {skema}.T_CLAIM_FAC_RETRO (
  ID                      VARCHAR2(32) NOT NULL,
  CLAIM_ID                VARCHAR2(32) NOT NULL,
  NOURUT                  NUMBER(5) NOT NULL,
  REINSURER_ID            VARCHAR2(64),
  REINSURER_NAME          VARCHAR2(500),
  SHARE_PCT               NUMBER(38,10),
  RI_COMMISSION_PCT       NUMBER(38,10),
  ADDITIONAL_INFO         VARCHAR2(1000),
  TOTAL_ESTIMATION_REINS  NUMBER(38,10),
  CONSTRAINT PK_CLAIM_FAC_RETRO PRIMARY KEY (ID),
  CONSTRAINT FK_CLAIM_FAC_RETRO_CLAIM FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE,
  CONSTRAINT UQ_CLAIM_FAC_RETRO_NOURUT UNIQUE (CLAIM_ID, NOURUT)
)
/

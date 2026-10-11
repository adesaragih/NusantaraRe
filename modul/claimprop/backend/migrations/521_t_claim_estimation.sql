-- 521 - T_CLAIM_ESTIMATION <- pyWorkPage.ClaimData.EstimationList (Section OutstandingClaim_Est grid "Estimation List").
-- Induk T_GENERAL_CLAIM lewat CLAIM_ID (diagram sheet Claim Prop F15), ON DELETE CASCADE.
-- NOURUT = urutan baris di daftar Pega (aturan "baris terakhir" Pega bergantung urutan), unik per induk.
-- Penamaan diluruskan (AC 39): GrossEstimationPct berisi NILAI uang -> GROSS_ESTIMATION_VALUE;
-- TypeLossID / TypeLoss berisi identitas dan nama jenis treaty -> TREATY_TYPE_ID / TREATY_TYPE_NAME.
-- Total berjalan (EstimastionReserve, TotalGrossEstimasiIDR, ListTotalEstimation) tidak disimpan - turunan.
CREATE TABLE {skema}.T_CLAIM_ESTIMATION (
  ID                      VARCHAR2(32) NOT NULL,
  CLAIM_ID                VARCHAR2(32) NOT NULL,
  NOURUT                  NUMBER(5) NOT NULL,
  ESTIMATION_DATE         DATE,
  ESTIMATION_TYPE         VARCHAR2(16),
  TREATY_TYPE_ID          VARCHAR2(64),
  TREATY_TYPE_NAME        VARCHAR2(255),
  CURRENCY_ID             VARCHAR2(64),
  CURRENCY_NAME           VARCHAR2(64),
  KURS                    NUMBER(38,10),
  GROSS_ESTIMATION_VALUE  NUMBER(38,10),
  GROSS_ESTIMATION_IDR    NUMBER(38,10),
  PERSEN_RNM              NUMBER(38,10),
  ESTIMATION_VALUE        NUMBER(38,10),
  ESTIMATION_VALUE_IDR    NUMBER(38,10),
  NO_PLA                  VARCHAR2(64),
  IS_PRINT_FACE_CLAIM     VARCHAR2(16),
  CONSTRAINT PK_CLAIM_ESTIMATION PRIMARY KEY (ID),
  CONSTRAINT FK_CLAIM_ESTIMATION_CLAIM FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE,
  CONSTRAINT UQ_CLAIM_ESTIMATION_NOURUT UNIQUE (CLAIM_ID, NOURUT)
)
/

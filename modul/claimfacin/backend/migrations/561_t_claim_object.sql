-- 561 - T_CLAIM_OBJECT <- pyWorkPage.ClaimData.ObjectList (objek yang diklaim; keputusan work owner 09-10-2026
-- OQ-CFI-01 "tabel Prop + kolom item": tabel BARU). ID stabil dari SEQ_T_CLAIM (penunjuk item); KUNCI_POLIS =
-- letak objek di halaman polis (`LocationList(n)` / `VehicleList(n)` / `PersonList(n)`) - medan salinan polis
-- dibaca ulang dari JSON_POLIS, tidak disimpan. NOURUT berurut se-klaim.
-- NOL COMMIT (ADR-U-0029), nol MODIFY / DROP kolom yang sudah ada. -migrate dijalankan work owner.
CREATE TABLE {skema}.T_CLAIM_OBJECT (
  ID        VARCHAR2(32) NOT NULL,
  CLAIM_ID  VARCHAR2(32) NOT NULL,
  NOURUT    NUMBER(5) NOT NULL,
  KUNCI_POLIS               VARCHAR2(64),
  OBJECT_REF                VARCHAR2(200),
  OBJECT_NAME               VARCHAR2(2000),
  CFS                       VARCHAR2(16),
  PRINT_FACE_CLAIM          VARCHAR2(16),
  PLA_STATUS                VARCHAR2(16),
  IS_FAC_RETRO              VARCHAR2(16),
  IS_MORE_THAN_TREATY_LIMIT VARCHAR2(16),
  IS_KOMITE                 VARCHAR2(16),
  IS_PRINT_ACCEPT           VARCHAR2(16),
  REMARKS_PLA               VARCHAR2(4000),
  REMARKS_DLA               VARCHAR2(4000),
  SHARE_RETRO               VARCHAR2(16),
  NO_DLA                    VARCHAR2(64),
  DLA_STATUS                VARCHAR2(16),
  SAVE_SPREADING            VARCHAR2(16),
  CONSTRAINT PK_CLAIM_OBJECT PRIMARY KEY (ID),
  CONSTRAINT FK_CLAIM_OBJECT_CLAIM FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE,
  CONSTRAINT UQ_CLAIM_OBJECT_NOURUT UNIQUE (CLAIM_ID, NOURUT)
)
/

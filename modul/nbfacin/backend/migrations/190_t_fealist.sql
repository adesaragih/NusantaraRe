-- 190 - sub-tab FEA / Fire Extinguisher Availability (tiket 41): T_FEALIST - tabel BARU (rancangan flat tidak punya tabel
-- FEA; A142, amandemen loader `amandemenFEA`).
--
-- `[terverifikasi]` `.FEAList` milik baris `.LocationList(n)` (kelas ASM-FW-GISFW-Data-OfferFacIn-LocationReinsurance:
-- `Section\FEAList.xml`; fixture nb-fire-1 LocationList[0] berkunci `FEAList`), kelas baris
-- Data-OfferFacIn-OfferFEAList; medan `Section\InputFEA_IsUW.xml`: .APAR .Sprinkler .SmokeDetector .Hydrant .InfoFEA dan
-- halaman tertanam .DataFEA (.PrivateTruckBrigade .PrivateFireBrigade .TeamSOPSafety .TeamSOPRiskManagement) - dilipat
-- ke baris FEA (pola lipatan V-22/V-39). Induk T_LOCATIONLIST; banyak baris per lokasi, urut SEQ_NO.
-- Kolom sistem = pola tabel berulang rancangan (ID, IDPEGA, COB_GROUP, PARENT_ID, SEQ_NO, ROW_UID; pola T_ADDITIONALSHIP).
-- Jumlah unit disimpan TEKS VARCHAR2(50) (pola A129 / NUMBER_OF_FLOOR); kode dropdown VARCHAR2(50); Others Info
-- VARCHAR2(500) (pola V-6). Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_FEALIST (
  ID                       NUMBER(19) NOT NULL,
  IDPEGA                   VARCHAR2(50),
  COB_GROUP                VARCHAR2(20),
  PARENT_ID                NUMBER(19) NOT NULL,
  SEQ_NO                   NUMBER(5) NOT NULL,
  ROW_UID                  VARCHAR2(36) NOT NULL,
  APAR                     VARCHAR2(50),
  SPRINKLER                VARCHAR2(50),
  SMOKE_DETECTOR           VARCHAR2(50),
  HYDRANT                  VARCHAR2(50),
  PRIVATE_TRUCK_BRIGADE    VARCHAR2(50),
  PRIVATE_FIRE_BRIGADE     VARCHAR2(50),
  TEAM_SOP_SAFETY          VARCHAR2(50),
  TEAM_SOP_RISK_MANAGEMENT VARCHAR2(50),
  INFO_FEA                 VARCHAR2(500),
  CONSTRAINT PK_T_FEALIST PRIMARY KEY (ID),
  CONSTRAINT FK_FEALIST_LOCATIONLIST FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_LOCATIONLIST (ID)
)
/
CREATE SEQUENCE {skema}.SEQ_T_FEALIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

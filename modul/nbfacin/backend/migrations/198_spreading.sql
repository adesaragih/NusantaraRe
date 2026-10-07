-- 198 - Spreading kasus FIRE (tiket 48): % Share RNM kasus, TSI / premi Nusantara Re per coverage, dan
-- `.CoverageList(k).SpreadingList` (kelas Data-SpreadingRisk).
--
-- Nama / jalur = workbook RANCANGAN (`loader/skema_gen.go`) `[terverifikasi]`:
--   - T_GENERAL_POLIS.PERCENT_SHARE (`PercentShare`, `.OfferFacIn.PercentShare` - InputInwardFacultativeDtl sel 80);
--   - T_COVERAGELIST.TSI_NUSANTARA_RE / PREMI_NUSANTARA_RE (`TSINusantaraRe` / `PremiNusantaraRe`,
--     CountPremiAndTSIRNMFireMBU_ACT langkah 6);
--   - T_SPREADINGLIST, jalur `LocationList/Property/PropertyItemList/CoverageList/SpreadingList`, induk T_COVERAGELIST
--     (pola 193 / 194: PARENT_TABLE + SRC_PATH), dibuat SEBAGIAN: kolom yang ditulis tab Spreading NB FIRE saja
--     (TREATY_TYPE, TREATY_NAME, SHARE_PERCENTAGE, TSI_SPREADED, TSI_GROSS_SPREADED, PREMIUM_SPREADED, CLAIM_ESTIMATION) +
--     kolom sistem + CURRENCY_CODE rancangan. Kolom rancangan lain (FLAG_SPREADING, IS_INPUT_PCT, IS_OLD_DATA,
--     INDEX_PROPERTY_ITEM, ANEKA_ID, COVERAGE_ID, LOCATION_ID, TSI_TOP_RISK) milik EDM / Aneka / loader - tidak dibuat.
-- Tipe uang / persen NUMBER(38,8) (ADR-0016; rancangan menulis NUMBER polos - pola 193). PARENT_ID NUMBER(19) = ID
-- T_COVERAGELIST, ber-FK: tab Object menghapus spreading sebelum coverage-nya (K48-7).
-- Total spreading (TotalTSIPremiSpreadRNM / TotalSpreadingCurrency / TotalSpreadAll) TIDAK disimpan - dihitung (K-067).
-- Template Copy Spreading (SpreadingList.pxResults) TIDAK disimpan - halaman clipboard Pega (K48-2).
-- Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  PERCENT_SHARE NUMBER(38,8)
)
/
ALTER TABLE {skema}.T_COVERAGELIST ADD (
  TSI_NUSANTARA_RE   NUMBER(38,8),
  PREMI_NUSANTARA_RE NUMBER(38,8)
)
/
CREATE TABLE {skema}.T_SPREADINGLIST (
  ID                 NUMBER(19) NOT NULL,
  IDPEGA             VARCHAR2(50),
  COB_GROUP          VARCHAR2(20),
  PARENT_ID          NUMBER(19) NOT NULL,
  PARENT_TABLE       VARCHAR2(30),
  SRC_PATH           VARCHAR2(200),
  SEQ_NO             NUMBER(5) NOT NULL,
  ROW_UID            VARCHAR2(36) NOT NULL,
  TREATY_TYPE        VARCHAR2(50),
  TREATY_NAME        VARCHAR2(500),
  SHARE_PERCENTAGE   NUMBER(38,8),
  TSI_SPREADED       NUMBER(38,8),
  TSI_GROSS_SPREADED NUMBER(38,8),
  PREMIUM_SPREADED   NUMBER(38,8),
  CLAIM_ESTIMATION   NUMBER(38,8),
  CURRENCY_CODE      VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL,
  CONSTRAINT PK_T_SPREADINGLIST PRIMARY KEY (ID),
  CONSTRAINT FK_SPREADINGLIST_COVERAGE FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_COVERAGELIST (ID)
)
/
CREATE INDEX {skema}.IX_T_SPREADINGLIST_PARENT ON {skema}.T_SPREADINGLIST (PARENT_ID, SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_SPREADINGLIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

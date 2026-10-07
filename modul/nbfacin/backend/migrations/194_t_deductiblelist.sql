-- 194 - tab Coverage FIRE tahap C3 (tiket 45): T_DEDUCTIBLELIST sebagian - deductible per coverage.
--
-- Rancangan jalur `.../CoverageList/DeductibleList` (lima jalur, SEMUANYA berinduk T_COVERAGELIST) -> PARENT_ID ber-FK ke
-- T_COVERAGELIST (induk tunggal, pola 191); PARENT_TABLE / SRC_PATH tetap diisi seperti loader (NB: 'T_COVERAGELIST',
-- `LocationList/Property/PropertyItemList/CoverageList/DeductibleList`). Banyak baris per coverage urut SEQ_NO.
-- SEBAGIAN (pola A109): kolom sistem + sepuluh medan kontrak `Deductible` (tiket 45); FLAG_CURRENCY (R-3), DESCRIPTIONS,
-- BASIS_TYPE, CLAIM_CATEGORY, DEDUCTIBLE_TYPE, TYPE, VALUE, COVERAGE dan kolom turunan *_ID ditambah kemudian lewat ALTER.
-- Uang / persen rancangan NUMBER polos -> NUMBER(38,8) (ADR-0016). TYPE_DEDUCTIBLE(2) rancangan NUMBER polos berisi kode
-- `pyStandardValue` 0-7 (`DDL\TypeDeductible.xml`, `TypeDeductible2.xml`) -> NUMBER(5) (A170). TIME_EXCESS tetap
-- VARCHAR2(30) rancangan (teks desimal, A171). CURRENCY = rancangan K-069 `DEFAULT 'UNKNOWN' NOT NULL`: aplikasi tidak
-- menulis 'UNKNOWN' (K-012) - mata uang kosong = kolom tidak disisipkan, DEFAULT yang mengisinya (A169).
-- ID/PARENT_ID NUMBER(19) (A92). Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_DEDUCTIBLELIST (
  ID               NUMBER(19) NOT NULL,
  IDPEGA           VARCHAR2(50),
  COB_GROUP        VARCHAR2(20),
  PARENT_ID        NUMBER(19) NOT NULL,
  PARENT_TABLE     VARCHAR2(30),
  SRC_PATH         VARCHAR2(200),
  SEQ_NO           NUMBER(5) NOT NULL,
  ROW_UID          VARCHAR2(36) NOT NULL,
  AMOUNT           NUMBER(38,8),
  CONDITION        VARCHAR2(500),
  CURRENCY         VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  INPUT_CONDITION  VARCHAR2(500),
  MIN_MAX          VARCHAR2(50),
  PCT_DEDUCTIBLE   NUMBER(38,8),
  PCT_DEDUCTIBLE2  NUMBER(38,8),
  TIME_EXCESS      VARCHAR2(30),
  TYPE_DEDUCTIBLE  NUMBER(5),
  TYPE_DEDUCTIBLE2 NUMBER(5),
  CONSTRAINT PK_T_DEDUCTIBLELIST PRIMARY KEY (ID),
  CONSTRAINT FK_DEDUCTIBLELIST_COVERAGE FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_COVERAGELIST (ID)
)
/
CREATE INDEX {skema}.IX_T_DEDUCTIBLELIST_PARENT ON {skema}.T_DEDUCTIBLELIST (PARENT_ID)
/
CREATE SEQUENCE {skema}.SEQ_T_DEDUCTIBLELIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

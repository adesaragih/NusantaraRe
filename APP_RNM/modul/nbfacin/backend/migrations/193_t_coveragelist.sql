-- 193 - tab Coverage FIRE tahap C1 (tiket 43): T_PROPERTYITEMLIST + dua kolom total, T_COVERAGELIST sebagian.
--
-- T_PROPERTYITEMLIST (188 sebagian, A109): TOTAL_GROSS_PREMI (uang, CountPremi langkah 55-56) dan TOTAL_NET_RATE (rate,
-- tahap C2) = kolom RANCANGAN `NUMBER` polos -> NUMBER(38,8) (ADR-0016; penjaga presisiSah).
-- T_COVERAGELIST = tabel RANCANGAN berinduk JAMAK (rancangan: CargoList / PropertyItemList / AnekaList x2 / PersonList /
-- VehicleList), dibedakan PARENT_TABLE / SRC_PATH - NB menulis jalur `LocationList/Property/PropertyItemList/CoverageList`
-- saja (PARENT_TABLE 'T_PROPERTYITEMLIST'); PARENT_ID TANPA FK (pola A140), indeks (PARENT_TABLE, PARENT_ID).
-- SEBAGIAN (pola A109): kolom sistem + 24 medan kontrak CoverageObjek + PCT_ADJUSTMENT + CURRENCY_CODE; sisanya (premi
-- retro / share / akumulasi / zona / deductible / tahap C2-C4) ditambah kemudian lewat ALTER. Seluruh medan kontrak
-- ADA di rancangan - tanpa kolom baru. Uang / rate / persen rancangan NUMBER polos -> NUMBER(38,8); FIRST_LOSS /
-- FIRST_SCALE / LOST_LIMIT / EML_PML tetap VARCHAR2(50) rancangan (berisi teks desimal). CURRENCY_CODE = rancangan
-- K-069 `DEFAULT 'UNKNOWN' NOT NULL`, diisi mata uang item (aplikasi tidak menulis 'UNKNOWN', K-012).
-- ID/PARENT_ID NUMBER(19) (A92). Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_PROPERTYITEMLIST ADD (
  TOTAL_GROSS_PREMI NUMBER(38,8),
  TOTAL_NET_RATE    NUMBER(38,8)
)
/
CREATE TABLE {skema}.T_COVERAGELIST (
  ID                   NUMBER(19) NOT NULL,
  IDPEGA               VARCHAR2(50),
  COB_GROUP            VARCHAR2(20),
  PARENT_ID            NUMBER(19) NOT NULL,
  PARENT_TABLE         VARCHAR2(30),
  SRC_PATH             VARCHAR2(200),
  SEQ_NO               NUMBER(5) NOT NULL,
  ROW_UID              VARCHAR2(36) NOT NULL,
  COVERAGE             VARCHAR2(50),
  OLDID                VARCHAR2(50),
  COVERAGE_NOTE        VARCHAR2(500),
  COVERAGE_BASIS       VARCHAR2(50),
  DAY                  VARCHAR2(50),
  TSI                  NUMBER(38,8),
  INDEMNITY            VARCHAR2(50),
  RATE                 NUMBER(38,8),
  RATE_OJK             NUMBER(38,8),
  FIRST_LOSS           VARCHAR2(50),
  DISCOUNT_PERCENTAGE  NUMBER(38,8),
  TSI_LIABILITY        NUMBER(38,8),
  NET_RATE             NUMBER(38,8),
  LIMITOF_LIABILITY    NUMBER(38,8),
  PCT_LO_L             NUMBER(38,8),
  PRO_RATE_PERCENT     NUMBER(38,8),
  INDEMNITY_PERCENTAGE NUMBER(38,8),
  FIRST_SCALE          VARCHAR2(50),
  SUBLIMIT             NUMBER(38,8),
  LOST_LIMIT           VARCHAR2(50),
  EML_PML              VARCHAR2(50),
  DISCOUNT             NUMBER(38,8),
  PREMIUM              NUMBER(38,8),
  CONDITIONS           VARCHAR2(500),
  PCT_ADJUSTMENT       NUMBER(38,8),
  CURRENCY_CODE        VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL,
  CONSTRAINT PK_T_COVERAGELIST PRIMARY KEY (ID)
)
/
CREATE INDEX {skema}.IX_T_COVERAGELIST_PARENT ON {skema}.T_COVERAGELIST (PARENT_TABLE, PARENT_ID)
/
CREATE SEQUENCE {skema}.SEQ_T_COVERAGELIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

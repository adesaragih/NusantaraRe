-- 189 - sub-tab Occupation (tiket 40): T_OCCUPATIONLIST sebagian, T_TABLEOFLIMIT utuh.
--
-- Tabel RANCANGAN (`loader/skema_gen.go`): jalur `LocationList/Property/OccupationList` -> T_OCCUPATIONLIST (induk
-- T_PROPERTY) dan `.../OccupationList/TableOfLimit` -> T_TABLEOFLIMIT (induk T_OCCUPATIONLIST, satu halaman per okupasi).
-- T_OCCUPATIONLIST berinduk JAMAK di rancangan (juga `RiskLocation/OccupationList`, `VehicleList/OccupationList`):
-- PARENT_TABLE / SRC_PATH membedakan induknya (isi = loader: 'T_PROPERTY' / jalur), sehingga PARENT_ID TANPA FK;
-- diganti indeks (PARENT_TABLE, PARENT_ID). ID/PARENT_ID NUMBER(19) (A92).
-- SEBAGIAN (pola A109): CATEGORY (.Category okupasi; layar memakai .TableOfLimit.Category), IS_OLD_DATA, LOCATION_ID
-- (turunan) - ditambah tiket 23. OCCUPATION_ID / OCCUPATION_NAME VARCHAR2(1000) = lebar sumbernya OCCUPATION.OLDID /
-- NAME (rancangan 50 / 500; pola butir 80, A138, amandemen loader `amandemenLebar`).
-- T_TABLEOFLIMIT.PCT_LIMIT: rancangan NUMBER, TEKS APA ADANYA (keputusan work owner butir 68.1, `kolomTeksMenyimpang`:
-- koma desimal dan spasi ujung dipertahankan); lebar VARCHAR2(50) = keputusan agent A139 (butir 68.1 menahan
-- lebarnya; contoh `DDL\TABLEOFLIMIT.xml` "70,000 " 7 bita). Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_OCCUPATIONLIST (
  ID              NUMBER(19) NOT NULL,
  IDPEGA          VARCHAR2(50),
  COB_GROUP       VARCHAR2(20),
  PARENT_ID       NUMBER(19) NOT NULL,
  PARENT_TABLE    VARCHAR2(30),
  SRC_PATH        VARCHAR2(200),
  SEQ_NO          NUMBER(5) NOT NULL,
  ROW_UID         VARCHAR2(36) NOT NULL,
  OCCUPATION_ID   VARCHAR2(1000),
  OCCUPATION_NAME VARCHAR2(1000),
  CONSTRAINT PK_T_OCCUPATIONLIST PRIMARY KEY (ID)
)
/
CREATE INDEX {skema}.IX_T_OCCUPATIONLIST_PARENT ON {skema}.T_OCCUPATIONLIST (PARENT_TABLE, PARENT_ID)
/
CREATE SEQUENCE {skema}.SEQ_T_OCCUPATIONLIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.T_TABLEOFLIMIT (
  ID          NUMBER(19) NOT NULL,
  IDPEGA      VARCHAR2(50),
  COB_GROUP   VARCHAR2(20),
  PARENT_ID   NUMBER(19) NOT NULL,
  CATEGORY    VARCHAR2(50),
  DESCRIPTION VARCHAR2(500),
  PCT_LIMIT   VARCHAR2(50),
  CONSTRAINT PK_T_TABLEOFLIMIT PRIMARY KEY (ID),
  CONSTRAINT UQ_T_TABLEOFLIMIT_PARENT UNIQUE (PARENT_ID),
  CONSTRAINT FK_TABLEOFLIMIT_OCCUPATION FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_OCCUPATIONLIST (ID)
)
/
CREATE SEQUENCE {skema}.SEQ_T_TABLEOFLIMIT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

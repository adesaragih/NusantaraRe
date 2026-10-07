-- 188 - sub-tab Object Item (tiket 39): T_PROPERTYITEMLIST sebagian + enam kolom baru.
--
-- Tabel RANCANGAN (`loader/skema_gen.go`): jalur `LocationList/Property/PropertyItemList`, induk T_PROPERTY; satu baris
-- per item, urut SEQ_NO (= PropertyItemNo 1..n). Nama/tipe = rancangan; ID/PARENT_ID NUMBER(19) (A92).
-- Uang dan persen (rancangan NUMBER polos) -> NUMBER(38,8): TSI_OBJECT_ITEM (uang, ADR-0016), PCT_ADJUST2 /
-- PCT_ADJUST_OTHER (persen, penjaga presisiSah). CURRENCY = rancangan `DEFAULT 'UNKNOWN' NOT NULL` (K-069); aplikasi
-- TIDAK menulis 'UNKNOWN' (K-012) - mata uang wajib diisi layar (A133).
-- SEBAGIAN (pola A109 / butir 78.4): kolom yang tidak dipakai layar ini - CURRENCY_ID, CURRENCY_OLD_ID, FLAG_NET_RATE,
-- IS_OLD_DATA, PERCENTAGE_ADJUSTMENT, PROPERTY_ID, SELECTED_LOCATION_ADDRESS, SELECTED_OBJECT_ITEM, dan uang premi
-- TOTAL_GROSS_PREMI / TOTAL_NET_RATE / TOTAL_PREMIUM_NUSANTARA_RE - ditambah tiket 23 lewat ALTER.
-- Enam kolom BARU (bukan di rancangan; A132, amandemen loader `amandemenItem`) dari
-- `Section\PropertyItemFacIn_Section.xml`: PROPERTY_YEAR (.PropertyYear), UNIT (.Unit), CONDITION (.Condition), YEAR
-- (.Year, Year of Planting), NO_OF_TREE (.NoOfTree), AREA_HECTAR (.AreaHectar). Nama dan tipe = medan bernama sama di
-- tabel rancangan lain (YEAR / UNIT VARCHAR2(50), CONDITION VARCHAR2(500)); angka disimpan teks (pola A129).
-- Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_PROPERTYITEMLIST (
  ID                 NUMBER(19) NOT NULL,
  IDPEGA             VARCHAR2(50),
  COB_GROUP          VARCHAR2(20),
  PARENT_ID          NUMBER(19) NOT NULL,
  SEQ_NO             NUMBER(5) NOT NULL,
  ROW_UID            VARCHAR2(36) NOT NULL,
  CURRENCY           VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  IS_ADJUSTABLE_FLAG VARCHAR2(10),
  ITEM_TYPE          VARCHAR2(50),
  ITEM_TYPE_ID       VARCHAR2(50),
  PCT_ADJUST2        NUMBER(38,8),
  PCT_ADJUST_OTHER   NUMBER(38,8),
  PROPERTI_ITEM_NOTE VARCHAR2(500),
  PROPERTY_ITEM_NO   VARCHAR2(50),
  REMARK             VARCHAR2(500),
  TSI_OBJECT_ITEM    NUMBER(38,8),
  PROPERTY_YEAR      VARCHAR2(50),
  UNIT               VARCHAR2(50),
  CONDITION          VARCHAR2(500),
  YEAR               VARCHAR2(50),
  NO_OF_TREE         VARCHAR2(50),
  AREA_HECTAR        VARCHAR2(50),
  CONSTRAINT PK_T_PROPERTYITEMLIST PRIMARY KEY (ID),
  CONSTRAINT FK_PROPERTYITEMLIST_PROPERTY FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_PROPERTY (ID)
)
/
CREATE SEQUENCE {skema}.SEQ_T_PROPERTYITEMLIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

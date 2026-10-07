-- 199 - tab Inw Fac Cedant Panels kasus FIRE (tiket 49): Share Cedant Type, Share of Ceding, dan
-- `.OfferFacIn.CedingCedantList`.
--
-- Nama / jalur = workbook RANCANGAN (`loader/skema_gen.go`) `[terverifikasi]`:
--   - T_GENERAL_POLIS.SHARE_CEDANT_TYPE (`ShareCedantType`; DDL\ShareCedantType.xml "0" Gross / "1" Share RNM) -
--     rancangan NUMBER polos -> NUMBER(5) (kode, pola A170);
--   - T_QUOTATIONDATA.SHARE_OF_CEDING (`ShareOfCeding`) - rancangan NUMBER, tetapi Pega menulis TEKS ("100%" atau
--     PercentShare + "%", Activity\SetShareOfCeding.xml; 94 dari 94 contoh DDL\CONTOH berbentuk "N%") -> VARCHAR2(50),
--     isi apa adanya (pola T_TABLEOFLIMIT.PCT_LIMIT, butir 68.1);
--   - T_CEDINGCEDANTLIST, jalur `CedingCedantList`, induk T_GENERAL_POLIS, dibuat UTUH (9 kolom rancangan). ID NUMBER(19)
--     (sequence, A92); PARENT_ID = T_GENERAL_POLIS.ID VARCHAR2(32) (PK teks, butir 76.1); SHARE_CEDING NUMBER(38,8)
--     (persen, ADR-0016).
-- TOTAL_TSI_NUSA_RE_SPREADING / TOTAL_PREMI_NUSA_RE TIDAK dibuat - dihitung dari coverage tiap baca (K49-3).
-- `CedingCedantList/CurrencyList` (T_CEDING_CURRENCYLIST) TIDAK dibuat - menunggu tab Payment (C-1).
-- Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  SHARE_CEDANT_TYPE NUMBER(5)
)
/
ALTER TABLE {skema}.T_QUOTATIONDATA ADD (
  SHARE_OF_CEDING VARCHAR2(50)
)
/
CREATE TABLE {skema}.T_CEDINGCEDANTLIST (
  ID             NUMBER(19) NOT NULL,
  IDPEGA         VARCHAR2(50),
  COB_GROUP      VARCHAR2(20),
  PARENT_ID      VARCHAR2(32) NOT NULL,
  SEQ_NO         NUMBER(5) NOT NULL,
  ROW_UID        VARCHAR2(36) NOT NULL,
  CEDING_CO      VARCHAR2(50),
  CEDING_CO_NAME VARCHAR2(500),
  SHARE_CEDING   NUMBER(38,8),
  CONSTRAINT PK_T_CEDINGCEDANTLIST PRIMARY KEY (ID),
  CONSTRAINT FK_CEDINGCEDANTLIST_GENERAL FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID)
)
/
CREATE INDEX {skema}.IX_CEDINGCEDANTLIST_PARENT ON {skema}.T_CEDINGCEDANTLIST (PARENT_ID, SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_CEDINGCEDANTLIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

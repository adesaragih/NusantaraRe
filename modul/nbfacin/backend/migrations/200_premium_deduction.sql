-- 200 - tab Premium Deduction / Brokerage Fee kasus FIRE (tiket 50): persen Premium Deduction / Brokerage Fee, centang
-- "Exclude PPN ?", dan `.OfferFacIn.CurrencyList(n).Policy.Payment` per mata uang.
--
-- Nama / jalur = workbook RANCANGAN (`loader/skema_gen.go`) `[terverifikasi]`:
--   - T_GENERAL_POLIS.PPN_CHECK (`.OfferFacIn.PPnCheck`; rancangan NUMBER -> NUMBER(5): 1 = true, 0 = false) - ⚠️ DIGANTI
--     migrasi 202 menjadi VARCHAR2(10) "true" / "false" (butir 68.1, kolomTeksMenyimpang). Berkas ini SUDAH jalan di DEV
--     (T_MIGRASI 06-10-2026) dengan NUMBER(5) dan TIDAK diubah lagi: pra-terbang tidak memeriksa langkah tercatat, jadi
--     mengubahnya membuat lingkungan baru berbeda dari DEV.
--     PAY_RI_COMMISION / PAY_PCT_BROKERAGE_FEE (`.PolicyData.Payment.RICommision` / `.PctBrokerageFee`, lipatan V-22b);
--   - T_CURRENCYLIST, jalur `CurrencyList`, induk T_GENERAL_POLIS, dibuat SEBAGIAN: kolom sistem + CURRENCY_CODE + NAME +
--     PAY_PREMIUM / PAY_COMMISION / PAY_BROKERAGE_FEE / PAY_PPH / PAY_PPN (`.Policy.Payment.*`, lipatan V-22b) +
--     PAY_DEDUCTION2 (amandemen tiket 50 K50-2, `loader/amandemen.go`). Kolom rancangan lain (OLD_ID, PAY_NET_PREMIUM,
--     PAY_TSI_TOTAL, PAY_EDM_*, PREMIUM, TSI, SUM_TOTAL_PAYMENT, CURRENCY_REF_ID) milik tab Payment / EDM / loader - tidak
--     dibuat.
-- Uang / persen NUMBER(38,8) (ADR-0016; rancangan NUMBER polos). ID NUMBER(19) (sequence, A92); PARENT_ID = T_GENERAL_POLIS.ID
-- VARCHAR2(32) (butir 76.1). Satu baris per mata uang objek, diganti utuh tiap Save tab Premium Deduction.
-- Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  PPN_CHECK             NUMBER(5),
  PAY_RI_COMMISION      NUMBER(38,8),
  PAY_PCT_BROKERAGE_FEE NUMBER(38,8)
)
/
CREATE TABLE {skema}.T_CURRENCYLIST (
  ID                NUMBER(19) NOT NULL,
  IDPEGA            VARCHAR2(50),
  COB_GROUP         VARCHAR2(20),
  PARENT_ID         VARCHAR2(32) NOT NULL,
  SEQ_NO            NUMBER(5) NOT NULL,
  ROW_UID           VARCHAR2(36) NOT NULL,
  CURRENCY_CODE     VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL,
  NAME              VARCHAR2(500),
  PAY_PREMIUM       NUMBER(38,8),
  PAY_COMMISION     NUMBER(38,8),
  PAY_DEDUCTION2    NUMBER(38,8),
  PAY_BROKERAGE_FEE NUMBER(38,8),
  PAY_PPH           NUMBER(38,8),
  PAY_PPN           NUMBER(38,8),
  CONSTRAINT PK_T_CURRENCYLIST PRIMARY KEY (ID),
  CONSTRAINT FK_CURRENCYLIST_GENERAL FOREIGN KEY (PARENT_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID)
)
/
CREATE INDEX {skema}.IX_CURRENCYLIST_PARENT ON {skema}.T_CURRENCYLIST (PARENT_ID, SEQ_NO)
/
CREATE SEQUENCE {skema}.SEQ_T_CURRENCYLIST START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

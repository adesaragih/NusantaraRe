-- Jalur mundur 485: PRODKE kembali 1-based (bentuk 480) - kebalikan urutan maju.
--
-- ⚠️ Baris JSON_POLIS yang sudah ditulis Confirm (PRODKE 0-based) TIDAK disentuh (tabel warisan, nol DDL).
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
UPDATE {skema}.T_PREMIUM_LIST SET PROD_KE = PROD_KE + 1 WHERE ID LIKE 'EDML-%' AND PROD_KE IS NOT NULL
/
UPDATE {skema}.T_PREMIUM_LIST SET PROD_KE = 1 WHERE ID LIKE 'NBLF-%' AND PROD_KE = 0
/
ALTER TABLE {skema}.T_PREMIUM_LIST MODIFY (
  PROD_KE NUMBER(5) DEFAULT 1
)
/

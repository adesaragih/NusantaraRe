-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE EDITIONABLE FUNCTION "POOLDATA"."GETCURRENCYSTANDARD" (
   i_kurs_id    IN   m_currencystandard.ID%TYPE,
   i_tgl_kurs   IN   m_currencystandard.CurrencyDate%TYPE
)
   RETURN NUMBER result_cache
IS
   vhasil   m_currencystandard.CurrencyValue%TYPE;
BEGIN
   SELECT CurrencyValue
     INTO vhasil
     FROM (SELECT   CurrencyValue
               FROM m_currencystandard
              WHERE ID = i_kurs_id
                AND TRUNC (CurrencyDate) <= TRUNC (sysdate)
           ORDER BY CurrencyDate DESC)
    WHERE ROWNUM < 2;

   RETURN vhasil;
EXCEPTION
   WHEN NO_DATA_FOUND
   THEN
      RETURN 1;
END ;
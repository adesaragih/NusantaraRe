-- Jalur mundur 030 - mengembalikan bentuk 013 apa adanya.
--
-- ⚠️ Melebarkan kolom SELALU aman; yang tidak aman arah sebaliknya. Jalur
-- mundur ini karena itu tidak dapat gagal karena data.
ALTER TABLE {skema}.T_KOMITE_KOMITELIST DROP CONSTRAINT FK_KOMITELIST_KOMITE
/

ALTER TABLE {skema}.T_GENERAL_KOMITE MODIFY (
  ID            VARCHAR2(40),
  ADJUSTMENT_ID VARCHAR2(40)
)
/

ALTER TABLE {skema}.T_KOMITE_KOMITELIST MODIFY (
  ID             VARCHAR2(40),
  DATA_KOMITE_ID VARCHAR2(40)
)
/

ALTER TABLE {skema}.T_KOMITE_KOMITELIST ADD CONSTRAINT FK_KOMITELIST_KOMITE
  FOREIGN KEY (DATA_KOMITE_ID)
  REFERENCES {skema}.T_GENERAL_KOMITE (ID)
/

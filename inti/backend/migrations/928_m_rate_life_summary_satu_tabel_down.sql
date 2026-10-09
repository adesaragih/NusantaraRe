-- Mundur 928: kembalikan keadaan sesudah 926 + 927 - tabel flat RATE_LIFE_SUMMARY (bentuk 926) berisi data dari kolom
-- M_RATE_LIFE_SUMMARY, lalu JSONDATA dibangun ulang dari kolom (JSON_OBJECT, kunci sama dengan JSON Pega: USEDBY, TYPE,
-- MODIFIEDDATE, OPERATORID) dan constraint IS JSON dipulihkan. Sesudahnya 927_down membuang keempat kolom.
-- ⚠️ FLAG lama TIDAK dapat dibangun dari kolom; satu-satunya sumbernya cadangan CSV ID + JSONDATA (LANGKAH-WO (a)).
-- Ringkasan yang dihapus 928 (sudah dihapus aplikasi) pun hanya ada di cadangan itu.
-- ⚠️ JSONDATA dikembalikan NULLABLE: definisi asli DEV yang tercatat di repo hanya "ID VARCHAR2(10) PK + JSONDATA CLOB
-- (IS JSON)" (fakta WO 07-10-2026, MODUL.md R6) - tanpa bukti NOT NULL. Bila katalog DEV menunjukkan NOT NULL, DBA
-- menambahkannya sesudah mundur (LANGKAH-WO-RIRATELIFE-SATU-TABEL.md, Jalur mundur).
CREATE TABLE {skema}.RATE_LIFE_SUMMARY (
  ID            VARCHAR2(10) NOT NULL,
  USEDBY        VARCHAR2(500),
  TYPE          VARCHAR2(100),
  MODIFIEDDATE  VARCHAR2(50),
  OPERATORID    VARCHAR2(200),
  CONSTRAINT PK_RATE_LIFE_SUMMARY PRIMARY KEY (ID)
)
/
CREATE INDEX {skema}.IX_RATE_LIFE_SUMMARY_NAMA ON {skema}.RATE_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))
/
INSERT INTO {skema}.RATE_LIFE_SUMMARY (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID)
SELECT m.ID, m.USEDBY, m.TYPE, m.MODIFIEDDATE, m.OPERATORID FROM {skema}.M_RATE_LIFE_SUMMARY m
/
DROP INDEX {skema}.IX_M_RATE_LIFE_SUMMARY_NAMA
/
ALTER TABLE {skema}.M_RATE_LIFE_SUMMARY ADD (
  JSONDATA CLOB
)
/
UPDATE {skema}.M_RATE_LIFE_SUMMARY m
   SET JSONDATA = JSON_OBJECT('USEDBY' VALUE m.USEDBY, 'TYPE' VALUE m.TYPE, 'MODIFIEDDATE' VALUE m.MODIFIEDDATE,
                              'OPERATORID' VALUE m.OPERATORID ABSENT ON NULL RETURNING CLOB)
/
ALTER TABLE {skema}.M_RATE_LIFE_SUMMARY ADD CONSTRAINT ENSURE_M_RATE_LIFE_SUMMARY_JSON CHECK (JSONDATA IS JSON)
/

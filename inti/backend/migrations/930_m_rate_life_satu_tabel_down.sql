-- Mundur 930: kembalikan JSONDATA dari kolom (JSON_OBJECT, kunci sama dengan JSON Pega: IDUSEDBY, USEDBY, TYPE, GENDER,
-- CONTRACT, AGE, RATE; ABSENT ON NULL), constraint IS JSON, dan view RATE_LIFE PERSIS definisi DEV; buang indeks
-- IDUSEDBY. Sesudahnya 929_down membuang ketujuh kolom.
-- ⛔ HANYA bila 930 TERCATAT di T_MIGRASI (Bongkar sendiri hanya membongkar langkah tercatat). Bila JSONDATA sudah
-- dibuang tetapi 930 tidak tercatat: selesaikan langkah MAJU dulu (LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md, Pemulihan).
-- ⚠️ FLAG lama (98 ribu baris JSON) TIDAK dapat dibangun dari kolom; satu-satunya sumbernya cadangan CSV ID + JSONDATA
-- (LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md (a)).
-- ⚠️ JSONDATA dikembalikan NULLABLE: definisi asli DEV yang tercatat di repo hanya "ID VARCHAR2(10) PK + JSONDATA CLOB
-- (IS JSON)" (fakta WO 07-10-2026, MODUL.md R7; katalog claimlife KATALOG-TABEL-PESERTA-DAN-TREATY.md b121) - tanpa
-- bukti NOT NULL. LANGKAH-WO (a) mencatat NULLABLE aslinya; bila 'N', DBA menambah NOT NULL sesudah mundur.
--
-- Aman DIULANG sesudah gagal di tengah: buang indeks, tambah JSONDATA, dan tambah constraint lewat blok berpelindung
-- katalog; UPDATE menulis nilai yang sama; CREATE VIEW terakhir.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RATE_LIFE' AND INDEX_NAME = 'IX_M_RATE_LIFE_IDUSEDBY';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP INDEX {skema}.IX_M_RATE_LIFE_IDUSEDBY';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RATE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RATE_LIFE ADD (JSONDATA CLOB)';
  END IF;
END;
/
UPDATE {skema}.M_RATE_LIFE m
   SET JSONDATA = JSON_OBJECT('IDUSEDBY' VALUE m.IDUSEDBY, 'USEDBY' VALUE m.USEDBY, 'TYPE' VALUE m.TYPE,
                              'GENDER' VALUE m.GENDER, 'CONTRACT' VALUE m.CONTRACT, 'AGE' VALUE m.AGE,
                              'RATE' VALUE m.RATE ABSENT ON NULL RETURNING CLOB)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RATE_LIFE' AND CONSTRAINT_NAME = 'ENSURE_M_RATE_LIFE_JSON';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RATE_LIFE ADD CONSTRAINT ENSURE_M_RATE_LIFE_JSON CHECK (JSONDATA IS JSON)';
  END IF;
END;
/
CREATE VIEW {skema}.RATE_LIFE AS
SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER, a.JSONDATA.CONTRACT,
       a.JSONDATA.AGE, a.JSONDATA.RATE
FROM {skema}.M_RATE_LIFE a
/

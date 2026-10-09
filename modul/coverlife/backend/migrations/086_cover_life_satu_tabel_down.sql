-- Mundur 086: kembalikan JSONDATA dari kolom COVER / NOTE (JSON_OBJECT, kunci "Cover" / "Note" - huruf persis kunci
-- Pega, supaya view yang dibuat ulang (`a.JSONDATA.Cover`, `a.JSONDATA.Note`) membacanya; ABSENT ON NULL - NOTE NULL
-- kembali TANPA kunci, view tetap membacanya NULL), constraint IS JSON bernama asli ENSURE_M_COVER_LIFE_JSON, lalu view
-- COVER_LIFE PERSIS teks DEV (fakta WO 08-10-2026: SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note FROM M_COVER_LIFE a -
-- ber-{skema}). Sesudahnya 085_down membuang COVER / NOTE.
-- HANYA bila 086 TERCATAT di T_MIGRASI (LANGKAH-WO-COVERLIFE.md, Pemulihan). Kunci pxObjClass, pyRuleHarness, dan bentuk
-- JSON asli hanya dari cadangan P3. JSONDATA dikembalikan NULLABLE (LANGKAH-WO (a) mencatat NULLABLE aslinya).
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 091_down / 943_down): tanpa kolom COVER / NOTE (sumber pembangunan
-- ulang) UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM JSONDATA ditambah - tidak ada
-- JSONDATA kosong yang setengah jadi. Bila kolomnya ada, tidak ada yang berubah (WHERE 1 = 0).
-- Aman DIULANG sesudah gagal di tengah: kolom dan constraint lewat blok berpelindung katalog; UPDATE menulis nilai yang
-- sama; CREATE VIEW terakhir (view yang sudah ada = ORA-00955, lewati pernyataan itu).
UPDATE {skema}.M_COVER_LIFE SET COVER = COVER, NOTE = NOTE WHERE 1 = 0
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_COVER_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_COVER_LIFE ADD (JSONDATA CLOB)';
  END IF;
END;
/
UPDATE {skema}.M_COVER_LIFE m
   SET JSONDATA = JSON_OBJECT('Cover' VALUE m.COVER, 'Note' VALUE m.NOTE ABSENT ON NULL RETURNING CLOB)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_COVER_LIFE' AND CONSTRAINT_NAME = 'ENSURE_M_COVER_LIFE_JSON';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_COVER_LIFE ADD CONSTRAINT ENSURE_M_COVER_LIFE_JSON CHECK (JSONDATA IS JSON)';
  END IF;
END;
/
CREATE VIEW {skema}.COVER_LIFE AS
SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note
FROM {skema}.M_COVER_LIFE a
/

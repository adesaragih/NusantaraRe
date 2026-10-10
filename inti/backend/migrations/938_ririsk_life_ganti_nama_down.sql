-- Mundur 938: RIRISK_LIFE (tabel) kembali bernama M_RIRISK_LIFE, lalu view RIRISK_LIFE dibuat ulang PERSIS definisi
-- DEV (ber-{skema}). Berjalan SESUDAH 940_down / 939_down, yang lebih dulu mengembalikan JSONDATA dan membuang AGE.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, alasan sama dengan 935_down): tabel RIRISK_LIFE tanpa JSONDATA =
-- ORA-00904, Bongkar berhenti; tabel yang sudah tidak ada (mundur diulang) = ORA-00942, ditoleransi.
-- Aman diulang sampai CREATE VIEW (terakhir; view yang sudah ada = ORA-00955, lewati pernyataan itu).
UPDATE {skema}.RIRISK_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.RIRISK_LIFE RENAME TO M_RIRISK_LIFE';
  END IF;
END;
/
CREATE VIEW {skema}.RIRISK_LIFE AS
SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.AGE, a.JSONDATA.YEAR, a.JSONDATA.MONTH, a.JSONDATA.RISK,
       a.JSONDATA.CONTRACT
FROM {skema}.M_RIRISK_LIFE a
/

-- Mundur 935: RIRISK_LIFE_SUMMARY (tabel) kembali bernama M_RIRISK_LIFE_SUMMARY, lalu view RIRISK_LIFE_SUMMARY dibuat
-- ulang PERSIS definisi DEV (ber-{skema}). Berjalan SESUDAH 936_down / 937_down (urutan mundur menurun), yang lebih dulu
-- mengembalikan JSONDATA dan membuang kolom ringkasan.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 929_down): tabel RIRISK_LIFE_SUMMARY TANPA JSONDATA (937 membuangnya
-- tetapi tidak tercatat) = ORA-00904 saat parse; Bongkar berhenti sebelum apa pun berubah. Tanpa pelindung ini blok
-- RENAME melewati diri, CREATE VIEW atas tabel yang tidak ada dijawab ORA-00942 - yang DITOLERANSI Bongkar - dan
-- catatan 935 terhapus seolah sudah mundur. Bila tabel itu sudah tidak ada (mundur diulang sesudah RENAME balik),
-- ORA-00942 di sini ditoleransi dan langkah lanjut.
-- Aman diulang sampai CREATE VIEW (terakhir; view yang sudah ada = ORA-00955, lewati pernyataan itu).
UPDATE {skema}.RIRISK_LIFE_SUMMARY SET JSONDATA = JSONDATA WHERE 1 = 0
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE_SUMMARY' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.RIRISK_LIFE_SUMMARY RENAME TO M_RIRISK_LIFE_SUMMARY';
  END IF;
END;
/
CREATE VIEW {skema}.RIRISK_LIFE_SUMMARY AS
SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID
FROM {skema}.M_RIRISK_LIFE_SUMMARY a
/

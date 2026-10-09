-- Mundur 942: BENEFIT_LIFE (tabel) kembali bernama M_BENEFIT_LIFE, lalu view BENEFIT_LIFE dibuat ulang dengan teks
-- PERSIS definisi DEV (fakta WO 08-10-2026: SELECT a.ID, a.JSONDATA.Benefit FROM M_BENEFIT_LIFE a - ber-{skema}).
-- Berjalan SESUDAH 943_down / 944_down (urutan mundur menurun), yang lebih dulu mengembalikan JSONDATA dan membuang
-- kolom BENEFIT.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 935_down): tabel BENEFIT_LIFE TANPA JSONDATA (944 membuangnya tetapi
-- tidak tercatat) = ORA-00904 saat parse; Bongkar berhenti sebelum apa pun berubah. Tanpa pelindung ini blok RENAME
-- melewati diri, CREATE VIEW atas tabel yang tidak ada dijawab ORA-00942 - yang DITOLERANSI Bongkar - dan catatan 942
-- terhapus seolah sudah mundur. Bila tabel itu sudah tidak ada (mundur diulang sesudah RENAME balik), ORA-00942 di sini
-- ditoleransi dan langkah lanjut.
-- Aman diulang sampai CREATE VIEW (terakhir; view yang sudah ada = ORA-00955, lewati pernyataan itu).
UPDATE {skema}.BENEFIT_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.BENEFIT_LIFE RENAME TO M_BENEFIT_LIFE';
  END IF;
END;
/
CREATE VIEW {skema}.BENEFIT_LIFE AS
SELECT a.ID, a.JSONDATA.Benefit
FROM {skema}.M_BENEFIT_LIFE a
/

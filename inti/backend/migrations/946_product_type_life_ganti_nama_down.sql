-- Mundur 946: PRODUCT_TYPE_LIFE (tabel) kembali bernama M_PRODUCT_TYPE_LIFE, lalu view PRODUCT_TYPE_LIFE dibuat ulang
-- dengan teks PERSIS definisi DEV (fakta WO 08-10-2026, ber-{skema}). Berjalan SESUDAH 947_down / 948_down (urutan
-- mundur menurun), yang lebih dulu mengembalikan JSONDATA, membuang PK, dan membuang kelima kolom.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 942_down): tabel PRODUCT_TYPE_LIFE TANPA JSONDATA (948 membuangnya
-- tetapi tidak tercatat) = ORA-00904 saat parse; Bongkar berhenti sebelum apa pun berubah. Bila tabel itu sudah tidak
-- ada (mundur diulang sesudah RENAME balik), ORA-00942 di sini ditoleransi dan langkah lanjut.
-- Aman diulang sampai CREATE VIEW (terakhir; view yang sudah ada = ORA-00955, lewati pernyataan itu).
UPDATE {skema}.PRODUCT_TYPE_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.PRODUCT_TYPE_LIFE RENAME TO M_PRODUCT_TYPE_LIFE';
  END IF;
END;
/
CREATE VIEW {skema}.PRODUCT_TYPE_LIFE AS
SELECT a.JSONDATA.ID, a.JSONDATA.CoverName, a.JSONDATA.Business, a.JSONDATA.BusinessID, a.JSONDATA.Benefit, a.JSONDATA.BenefitID
FROM {skema}.M_PRODUCT_TYPE_LIFE a
/

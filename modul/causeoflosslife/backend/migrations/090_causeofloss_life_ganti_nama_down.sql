-- Mundur 090: CAUSEOFLOSS_LIFE (tabel) kembali bernama M_CAUSEOFLOSS_LIFE, lalu view CAUSEOFLOSS_LIFE dibuat ulang dengan
-- teks PERSIS definisi DEV (fakta WO 08-10-2026: SELECT a.ID, a.JSONDATA.CauseofLoss FROM M_CAUSEOFLOSS_LIFE a -
-- ber-{skema}). Berjalan SESUDAH 092_down / 091_down (urutan mundur menurun), yang lebih dulu mengembalikan JSONDATA dan
-- membuang kolom CAUSEOFLOSS.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 942_down): tabel CAUSEOFLOSS_LIFE TANPA JSONDATA (092 membuangnya
-- tetapi tidak tercatat) = ORA-00904 saat parse; Bongkar berhenti sebelum apa pun berubah. Tanpa pelindung ini blok
-- RENAME melewati diri, CREATE VIEW atas tabel yang tidak ada dijawab ORA-00942 - yang DITOLERANSI Bongkar - dan catatan
-- 090 terhapus seolah sudah mundur. Bila tabel itu sudah tidak ada (mundur diulang sesudah RENAME balik), ORA-00942 di
-- sini ditoleransi dan langkah lanjut.
-- Aman diulang sampai CREATE VIEW (terakhir; view yang sudah ada = ORA-00955, lewati pernyataan itu).
UPDATE {skema}.CAUSEOFLOSS_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.CAUSEOFLOSS_LIFE RENAME TO M_CAUSEOFLOSS_LIFE';
  END IF;
END;
/
CREATE VIEW {skema}.CAUSEOFLOSS_LIFE AS
SELECT a.ID, a.JSONDATA.CauseofLoss
FROM {skema}.M_CAUSEOFLOSS_LIFE a
/

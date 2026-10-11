-- LANGKAH-WO-PLANLIFE.md (a) - BACA-SAJA, SEBELUM -migrate. BUKTI K2 (keputusan work owner 08-10-2026) TANPA menulis:
-- kelima ekspresi pengisian 948 (`m.JSONDATA.CoverName`, `.Business`, `.BusinessID`, `.Benefit`, `.BenefitID` - notasi
-- titik PERSIS teks view lama `a.JSONDATA.<kunci>`, peka huruf, hanya aliasnya berbeda) dijalankan sebagai SELECT atas
-- M_PRODUCT_TYPE_LIFE dan dibandingkan dengan view PRODUCT_TYPE_LIFE baris demi baris; ID tabel dibandingkan dengan
-- JSONDATA.ID (view lama membaca ID dari JSON). Hasil yang diharapkan (fakta WO 08-10-2026):
--   A. N = 32 dan setiap N_* = 32 (ekspresi tidak-NULL; JSON_VALUE peka huruf pembanding bebas);
--   B. SAMA_SEMUA = 32 (kelima kolom identik dengan view + ID tabel = JSONDATA.ID);
--   C. NOL baris (baris yang akan MENGHENTIKAN pemeriksaan 948 - beda isi, ID NULL, ID ganda, ID tabel <> JSONDATA.ID);
--   D. MAKS_* <= lebar kolom tujuan (pembanding: CoverName 34, Business 29, Benefit 48; ID <= 6; *ID <= 10);
--   E. lebar BUSINESS.ID dan BENEFIT_LIFE.ID di katalog (lebar BUSINESSID / BENEFITID 947 = 10 mengikutinya);
--   F. NOL baris: BusinessID tanpa baris BUSINESS GROUPPANEL 009 bernama sama, BenefitID tanpa baris BENEFIT_LIFE
--      bernama sama (dasar K4 - pencocokan ulang server tidak menolak data lama);
--   G. NOL baris: CoverName kembar tanpa beda huruf (dasar K5 - Plan Name unik).
-- Kunci px* TIDAK dicetak. Dijalankan sebagai BERKAS (@...); keluaran ke plan_bukti_k1.txt - kirim ke executor /
-- laporan (butir yang menunggu WO).
SET LINESIZE 250 PAGESIZE 100 NUMWIDTH 12
SPOOL plan_bukti_k1.txt
PROMPT A. Cacah ekspresi 948 (notasi titik) dan JSON_VALUE peka huruf
SELECT COUNT(*) AS N,
       COUNT(m.JSONDATA.ID) AS N_ID_JSON, COUNT(m.JSONDATA.CoverName) AS N_COVERNAME, COUNT(m.JSONDATA.Business) AS N_BUSINESS,
       COUNT(m.JSONDATA.BusinessID) AS N_BUSINESSID, COUNT(m.JSONDATA.Benefit) AS N_BENEFIT,
       COUNT(m.JSONDATA.BenefitID) AS N_BENEFITID,
       COUNT(JSON_VALUE(m.JSONDATA, '$.CoverName')) AS N_JV_COVERNAME, COUNT(JSON_VALUE(m.JSONDATA, '$.BenefitID')) AS N_JV_BENEFITID
  FROM POOLDATA.M_PRODUCT_TYPE_LIFE m;
PROMPT B. Identik dengan view (pasangan lewat ID JSON = kolom ID view) dan ID tabel = JSONDATA.ID
SELECT SUM(CASE WHEN m.ID = m.JSONDATA.ID AND m.JSONDATA.CoverName = v.COVERNAME AND m.JSONDATA.Business = v.BUSINESS
                 AND m.JSONDATA.BusinessID = v.BUSINESSID AND m.JSONDATA.Benefit = v.BENEFIT
                 AND m.JSONDATA.BenefitID = v.BENEFITID THEN 1 ELSE 0 END) AS SAMA_SEMUA
  FROM POOLDATA.M_PRODUCT_TYPE_LIFE m
  LEFT JOIN POOLDATA.PRODUCT_TYPE_LIFE v ON v.ID = m.JSONDATA.ID;
PROMPT C. Baris yang akan menghentikan 948 (harus NOL baris)
SELECT m.ID, m.JSONDATA.ID AS ID_JSON FROM POOLDATA.M_PRODUCT_TYPE_LIFE m
 WHERE m.ID IS NULL OR m.JSONDATA.ID IS NULL OR m.ID <> m.JSONDATA.ID
    OR m.ID IN (SELECT d.ID FROM POOLDATA.M_PRODUCT_TYPE_LIFE d GROUP BY d.ID HAVING COUNT(*) > 1)
    OR (JSON_EXISTS(m.JSONDATA, '$.CoverName') AND m.JSONDATA.CoverName IS NULL)
    OR (JSON_EXISTS(m.JSONDATA, '$.Business') AND m.JSONDATA.Business IS NULL)
    OR (JSON_EXISTS(m.JSONDATA, '$.BusinessID') AND m.JSONDATA.BusinessID IS NULL)
    OR (JSON_EXISTS(m.JSONDATA, '$.Benefit') AND m.JSONDATA.Benefit IS NULL)
    OR (JSON_EXISTS(m.JSONDATA, '$.BenefitID') AND m.JSONDATA.BenefitID IS NULL);
PROMPT D. Panjang maksimum (byte)
SELECT MAX(LENGTHB(m.ID)) AS MAKS_ID, MAX(LENGTHB(m.JSONDATA.CoverName)) AS MAKS_COVERNAME,
       MAX(LENGTHB(m.JSONDATA.Business)) AS MAKS_BUSINESS, MAX(LENGTHB(m.JSONDATA.BusinessID)) AS MAKS_BUSINESSID,
       MAX(LENGTHB(m.JSONDATA.Benefit)) AS MAKS_BENEFIT, MAX(LENGTHB(m.JSONDATA.BenefitID)) AS MAKS_BENEFITID
  FROM POOLDATA.M_PRODUCT_TYPE_LIFE m;
PROMPT E. Lebar ID master di katalog
SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, DATA_LENGTH FROM SYS.ALL_TAB_COLUMNS
 WHERE OWNER = 'POOLDATA' AND ((TABLE_NAME = 'BUSINESS' AND COLUMN_NAME IN ('ID', 'NOTE', 'OLDID', 'GROUPPANEL'))
    OR (TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME IN ('ID', 'BENEFIT')))
 ORDER BY TABLE_NAME, COLUMN_NAME;
PROMPT F. Pasangan nama / ID yang TIDAK cocok master (harus NOL baris)
SELECT m.ID, CHR(66) AS JENIS FROM POOLDATA.M_PRODUCT_TYPE_LIFE m
 WHERE NOT EXISTS (SELECT 1 FROM POOLDATA.BUSINESS b WHERE b.GROUPPANEL = '009'
                    AND TO_CHAR(b.ID) = m.JSONDATA.BusinessID AND UPPER(TRIM(b.NOTE)) = UPPER(TRIM(m.JSONDATA.Business)))
UNION ALL
SELECT m.ID, CHR(69) FROM POOLDATA.M_PRODUCT_TYPE_LIFE m
 WHERE NOT EXISTS (SELECT 1 FROM POOLDATA.BENEFIT_LIFE e
                    WHERE e.ID = m.JSONDATA.BenefitID AND UPPER(TRIM(e.BENEFIT)) = UPPER(TRIM(m.JSONDATA.Benefit)));
SELECT COUNT(*) AS N_BUSINESS_009 FROM POOLDATA.BUSINESS WHERE GROUPPANEL = '009';
PROMPT G. Plan Name kembar tanpa beda huruf (harus NOL baris)
SELECT UPPER(TRIM(m.JSONDATA.CoverName)) AS NAMA, COUNT(*) AS N FROM POOLDATA.M_PRODUCT_TYPE_LIFE m
 GROUP BY UPPER(TRIM(m.JSONDATA.CoverName)) HAVING COUNT(*) > 1;
SPOOL OFF
-- Catatan: kueri F dijalankan SEBELUM modul benefitlife (942-944) atau SESUDAHNYA - BENEFIT_LIFE = view lama atau tabel
-- flat, kolom ID dan BENEFIT sama. Notasi titik membaca nama kunci PEKA HURUF sebagaimana ditulis; 948 tidak dapat
-- memakai JSON_VALUE / JSON_EXISTS karena teks EXECUTE IMMEDIATE blok berpelindung tidak boleh memuat tanda kutip.

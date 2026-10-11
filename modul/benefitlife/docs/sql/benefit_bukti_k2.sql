-- LANGKAH-WO-BENEFITLIFE.md (a) - BACA-SAJA, SEBELUM -migrate. BUKTI K2 (keputusan work owner 08-10-2026) TANPA
-- menulis: ekspresi pengisian migrasi 944 (`m.JSONDATA.Benefit` - notasi titik PERSIS teks view lama
-- `a.JSONDATA.Benefit`, hanya aliasnya berbeda) dijalankan sebagai SELECT atas M_BENEFIT_LIFE dan dibandingkan dengan
-- view BENEFIT_LIFE yang masih ada. Hasil yang diharapkan (fakta WO 08-10-2026):
--   A. N = 10, N_EKSPRESI = 10, N_VIEW = 10, N_JSON_VALUE = 10, N_ADA_KUNCI = 10, N_SAMA = 10;
--   B. sepuluh baris, kolom SAMA bernilai 'SAMA' di setiap baris;
--   C. NOL baris (baris yang akan MENGHENTIKAN pemeriksaan 944 / ekspresi berbeda dengan view);
--   D. NOL baris (kunci Benefit ada tetapi notasi titik membacanya NULL - batas pemeriksaan 944);
--   E. MAKS_BYTE <= 200 (pembanding 48);
--   F. N_NUMBER = 5, N_NUMBER_BEDA_ID = 0 (`.Number` = ID - tidak ada kolom NUMBER);
--   G. N_BUKAN_HURUF_BESAR = cacah Benefit yang BUKAN huruf besar seluruhnya (keterangan untuk pertanyaan terbuka T1).
-- Kunci pxCreateOperator / pxCreateOpName TIDAK dicetak. Dijalankan sebagai BERKAS (@...); keluaran ke
-- benefit_bukti_k2.txt - kirim ke executor / laporan (butir yang menunggu WO).
SET LINESIZE 250 PAGESIZE 100 NUMWIDTH 12
COLUMN EKSPRESI_944 FORMAT A50
COLUMN VIEW_LAMA FORMAT A50
COLUMN SAMA FORMAT A5
SPOOL benefit_bukti_k2.txt
PROMPT A. Cacah: ekspresi 944, view, JSON_VALUE peka huruf, JSON_EXISTS, dan yang identik
SELECT COUNT(*) AS N,
       COUNT(m.JSONDATA.Benefit) AS N_EKSPRESI,
       COUNT(v.BENEFIT) AS N_VIEW,
       COUNT(JSON_VALUE(m.JSONDATA, '$.Benefit')) AS N_JSON_VALUE,
       SUM(CASE WHEN JSON_EXISTS(m.JSONDATA, '$.Benefit') THEN 1 ELSE 0 END) AS N_ADA_KUNCI,
       SUM(CASE WHEN m.JSONDATA.Benefit = v.BENEFIT THEN 1 ELSE 0 END) AS N_SAMA
  FROM POOLDATA.M_BENEFIT_LIFE m
  LEFT JOIN POOLDATA.BENEFIT_LIFE v ON v.ID = m.ID;
PROMPT B. Baris demi baris (ekspresi 944 lawan view lama)
SELECT m.ID, m.JSONDATA.Benefit AS EKSPRESI_944, v.BENEFIT AS VIEW_LAMA,
       CASE WHEN m.JSONDATA.Benefit = v.BENEFIT THEN 'SAMA' ELSE 'BEDA' END AS SAMA
  FROM POOLDATA.M_BENEFIT_LIFE m
  LEFT JOIN POOLDATA.BENEFIT_LIFE v ON v.ID = m.ID
 ORDER BY m.ID;
PROMPT C. Baris yang akan menghentikan 944 / berbeda dengan view (harus NOL baris)
SELECT m.ID FROM POOLDATA.M_BENEFIT_LIFE m
  LEFT JOIN POOLDATA.BENEFIT_LIFE v ON v.ID = m.ID
 WHERE (JSON_EXISTS(m.JSONDATA, '$.Benefit') AND m.JSONDATA.Benefit IS NULL)
    OR v.BENEFIT IS NULL OR m.JSONDATA.Benefit IS NULL OR m.JSONDATA.Benefit <> v.BENEFIT;
PROMPT D. Kunci Benefit ada tetapi notasi titik membacanya NULL (nilai null / objek / > 4000 byte) - harus NOL baris
SELECT m.ID FROM POOLDATA.M_BENEFIT_LIFE m
 WHERE JSON_EXISTS(m.JSONDATA, '$.Benefit') AND m.JSONDATA.Benefit IS NULL;
PROMPT E. Panjang maksimum (byte) - harus <= 200 (lebar BENEFIT 943)
SELECT MAX(LENGTHB(m.JSONDATA.Benefit)) AS MAKS_BYTE FROM POOLDATA.M_BENEFIT_LIFE m;
PROMPT F. Number = ID (pembanding: 5 baris ber-Number, 0 berbeda)
SELECT COUNT(JSON_VALUE(m.JSONDATA, '$.Number')) AS N_NUMBER,
       SUM(CASE WHEN JSON_VALUE(m.JSONDATA, '$.Number') <> m.ID THEN 1 ELSE 0 END) AS N_NUMBER_BEDA_ID
  FROM POOLDATA.M_BENEFIT_LIFE m;
PROMPT G. Benefit yang belum huruf besar (SetUpperCase_DT - pertanyaan terbuka T1)
SELECT SUM(CASE WHEN m.JSONDATA.Benefit <> UPPER(m.JSONDATA.Benefit) THEN 1 ELSE 0 END) AS N_BUKAN_HURUF_BESAR
  FROM POOLDATA.M_BENEFIT_LIFE m;
SPOOL OFF
-- Catatan: notasi titik JSON Oracle membaca nama kunci PEKA HURUF sebagaimana ditulis (`Benefit`); JSON_VALUE dengan
-- jalur '$.Benefit' (juga peka huruf) dihitung terpisah sebagai pembanding bebas. 944 tidak dapat memakai JSON_VALUE /
-- JSON_EXISTS karena teks EXECUTE IMMEDIATE blok berpelindung tidak boleh memuat tanda kutip (polaPerintahKatalog);
-- karena itu 944 memakai ekspresi notasi titik yang SAMA dengan view, dan kueri D memastikan tidak ada kunci yang
-- terlewat olehnya.

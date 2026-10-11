-- LANGKAH-WO-CAUSEOFLOSSLIFE.md (a) - BACA-SAJA, SEBELUM -migrate. BUKTI K2 (keputusan work owner 08-10-2026) TANPA
-- menulis: ekspresi pengisian 092 (`m.JSONDATA.CauseofLoss` - notasi titik PERSIS teks view lama
-- `a.JSONDATA.CauseofLoss`, peka huruf, hanya aliasnya berbeda) dijalankan sebagai SELECT atas M_CAUSEOFLOSS_LIFE dan
-- dibandingkan dengan view CAUSEOFLOSS_LIFE baris demi baris. NULL = NULL dianggap SAMA (DECODE), persis aturan
-- pemeriksaan K1.4 092. Hasil yang diharapkan (fakta WO 08-10-2026):
--   A. N = 4; N_EKSPRESI = N_VIEW = 3 (100001 NULL di keduanya); N_ADA_KUNCI = 4 (kunci CauseofLoss ada di 4 baris);
--   B. empat baris, setiap baris SAMA = 1 (termasuk 100001: EKSPRESI dan VIEW sama-sama NULL);
--   C. NOL baris (baris yang akan MENGHENTIKAN pemeriksaan 092 - kolom akan berbeda dari view);
--   D. MAKS_BYTE <= 200 (pembanding 9) dan NOL baris > 200 byte;
--   E. NOL baris: Cause of Loss kembar tanpa beda huruf (dasar penolakan kembar K4 - data lama tidak melanggar);
--   F. baris yang dipakai M_PRODUCT_LIFE.CAUSEID (pembanding: 100004 x166, 100003 x11, 100002 x3, kosong x16; 100001
--      tidak dipakai - dilaporkan saja, tidak dihapus / diisi).
-- Kunci px* TIDAK dicetak. Dijalankan sebagai BERKAS (@...); keluaran ke col_bukti_k1.txt - kirim ke executor / laporan
-- (butir yang menunggu WO). Satu angka berbeda = BERHENTI, jangan -migrate.
SET LINESIZE 250 PAGESIZE 100 NUMWIDTH 12
COLUMN EKSPRESI FORMAT A40
COLUMN VIEW_NILAI FORMAT A40
SPOOL col_bukti_k1.txt
PROMPT A. Cacah ekspresi 092 (notasi titik), view, dan keberadaan kunci
SELECT COUNT(*) AS N, COUNT(m.JSONDATA.CauseofLoss) AS N_EKSPRESI,
       (SELECT COUNT(CAUSEOFLOSS) FROM POOLDATA.CAUSEOFLOSS_LIFE) AS N_VIEW,
       SUM(CASE WHEN JSON_EXISTS(m.JSONDATA, '$.CauseofLoss') THEN 1 ELSE 0 END) AS N_ADA_KUNCI
  FROM POOLDATA.M_CAUSEOFLOSS_LIFE m;
PROMPT B. Per baris: ekspresi 092 lawan view (SAMA = 1; NULL = NULL sama)
SELECT m.ID, m.JSONDATA.CauseofLoss AS EKSPRESI, v.CAUSEOFLOSS AS VIEW_NILAI,
       DECODE(m.JSONDATA.CauseofLoss, v.CAUSEOFLOSS, 1, 0) AS SAMA
  FROM POOLDATA.M_CAUSEOFLOSS_LIFE m
  LEFT JOIN POOLDATA.CAUSEOFLOSS_LIFE v ON v.ID = m.ID
 ORDER BY m.ID;
PROMPT C. Baris yang akan menghentikan pemeriksaan 092 (harus NOL baris)
SELECT m.ID FROM POOLDATA.M_CAUSEOFLOSS_LIFE m
  LEFT JOIN POOLDATA.CAUSEOFLOSS_LIFE v ON v.ID = m.ID
 WHERE DECODE(m.JSONDATA.CauseofLoss, v.CAUSEOFLOSS, 1, 0) = 0;
PROMPT D. Panjang maksimum (byte) dan nilai > 200 byte (harus NOL baris)
SELECT MAX(LENGTHB(m.JSONDATA.CauseofLoss)) AS MAKS_BYTE, MAX(LENGTHB(m.ID)) AS MAKS_ID FROM POOLDATA.M_CAUSEOFLOSS_LIFE m;
SELECT m.ID FROM POOLDATA.M_CAUSEOFLOSS_LIFE m
 WHERE LENGTHB(JSON_VALUE(m.JSONDATA, '$.CauseofLoss' RETURNING VARCHAR2(4000))) > 200;
PROMPT E. Cause of Loss kembar tanpa beda huruf (harus NOL baris)
SELECT UPPER(TRIM(m.JSONDATA.CauseofLoss)) AS NAMA, COUNT(*) AS N FROM POOLDATA.M_CAUSEOFLOSS_LIFE m
 WHERE m.JSONDATA.CauseofLoss IS NOT NULL
 GROUP BY UPPER(TRIM(m.JSONDATA.CauseofLoss)) HAVING COUNT(*) > 1;
PROMPT F. Pemakai di M_PRODUCT_LIFE.CAUSEID (baca-saja)
SELECT NVL(p.JSONDATA.CAUSEID, '(kosong)') AS CAUSEID, COUNT(*) AS N FROM POOLDATA.M_PRODUCT_LIFE p
 GROUP BY NVL(p.JSONDATA.CAUSEID, '(kosong)') ORDER BY 1;
SPOOL OFF
-- Catatan: 092 tidak dapat memakai DECODE / JSON_VALUE / JSON_EXISTS dengan literal karena teks EXECUTE IMMEDIATE blok
-- berpelindung tidak boleh memuat tanda kutip; bentuk tiga cabang OR di 092 setara dengan DECODE(...) = 0 di kueri C.
-- Kueri F: nama kunci CAUSEID di JSON M_PRODUCT_LIFE mengikuti fakta WO (prompt bab 1); ORA-00904 / nilai kosong di
-- sana = sesuaikan nama kuncinya, kueri A-E tetap berlaku.

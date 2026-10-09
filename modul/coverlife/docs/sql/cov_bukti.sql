-- LANGKAH-WO-COVERLIFE.md (a) - BACA-SAJA, SEBELUM -migrate. BUKTI C1 TANPA menulis: ekspresi pengisian 086
-- (`m.JSONDATA.Cover`, `m.JSONDATA.Note` - notasi titik PERSIS teks view lama `a.JSONDATA.Cover`, `a.JSONDATA.Note`,
-- peka huruf, hanya aliasnya berbeda) dijalankan sebagai SELECT atas M_COVER_LIFE dan dibandingkan dengan view
-- COVER_LIFE baris demi baris. NULL = NULL dianggap SAMA (DECODE), persis aturan pemeriksaan C1.3 086. Hasil yang
-- diharapkan (fakta WO 08-10-2026):
--   A. N = 4; N_COVER = N_VIEW_COVER = 4; N_NOTE = N_VIEW_NOTE = 0; N_ADA_KUNCI_NOTE = 0;
--   B. empat baris, SAMA = 1 di setiap baris (NOTE NULL di keduanya);
--   C. NOL baris (baris yang akan MENGHENTIKAN pemeriksaan 086);
--   D. MAKS_COVER <= 200 (pembanding 17) dan NOL baris > 200 byte;
--   E. NOL baris: Cover kembar tanpa beda huruf (dasar penolakan kembar C3);
--   F. sequence M_COVER_LIFE_SEQ: rumus Pega untuk NEXTVAL berikutnya (100005) BELUM terpakai (C2: aman).
-- Kunci px* TIDAK dicetak. Dijalankan sebagai BERKAS (@...); keluaran ke cov_bukti.txt - kirim ke executor / laporan
-- (butir yang menunggu WO). Satu angka berbeda = BERHENTI, jangan -migrate.
SET LINESIZE 250 PAGESIZE 100 NUMWIDTH 12
COLUMN EKSPRESI_COVER FORMAT A30
COLUMN VIEW_COVER FORMAT A30
SPOOL cov_bukti.txt
PROMPT A. Cacah ekspresi 086 (notasi titik), view, dan keberadaan kunci Note
SELECT COUNT(*) AS N, COUNT(m.JSONDATA.Cover) AS N_COVER, COUNT(m.JSONDATA.Note) AS N_NOTE,
       (SELECT COUNT(COVER) FROM POOLDATA.COVER_LIFE) AS N_VIEW_COVER,
       (SELECT COUNT(NOTE) FROM POOLDATA.COVER_LIFE) AS N_VIEW_NOTE,
       SUM(CASE WHEN JSON_EXISTS(m.JSONDATA, '$.Note') THEN 1 ELSE 0 END) AS N_ADA_KUNCI_NOTE
  FROM POOLDATA.M_COVER_LIFE m;
PROMPT B. Per baris: ekspresi 086 lawan view (SAMA = 1; NULL = NULL sama)
SELECT m.ID, m.JSONDATA.Cover AS EKSPRESI_COVER, v.COVER AS VIEW_COVER,
       DECODE(m.JSONDATA.Cover, v.COVER, 1, 0) * DECODE(m.JSONDATA.Note, v.NOTE, 1, 0) AS SAMA
  FROM POOLDATA.M_COVER_LIFE m
  LEFT JOIN POOLDATA.COVER_LIFE v ON v.ID = m.ID
 ORDER BY m.ID;
PROMPT C. Baris yang akan menghentikan pemeriksaan 086 (harus NOL baris)
SELECT m.ID FROM POOLDATA.M_COVER_LIFE m
  LEFT JOIN POOLDATA.COVER_LIFE v ON v.ID = m.ID
 WHERE DECODE(m.JSONDATA.Cover, v.COVER, 1, 0) = 0 OR DECODE(m.JSONDATA.Note, v.NOTE, 1, 0) = 0;
PROMPT D. Panjang maksimum (byte) dan nilai > lebar kolom (harus NOL baris)
SELECT MAX(LENGTHB(m.JSONDATA.Cover)) AS MAKS_COVER, MAX(LENGTHB(m.JSONDATA.Note)) AS MAKS_NOTE, MAX(LENGTHB(m.ID)) AS MAKS_ID
  FROM POOLDATA.M_COVER_LIFE m;
SELECT m.ID FROM POOLDATA.M_COVER_LIFE m
 WHERE LENGTHB(JSON_VALUE(m.JSONDATA, '$.Cover' RETURNING VARCHAR2(4000))) > 200
    OR LENGTHB(JSON_VALUE(m.JSONDATA, '$.Note' RETURNING VARCHAR2(4000))) > 1000;
PROMPT E. Cover kembar tanpa beda huruf (harus NOL baris)
SELECT UPPER(TRIM(m.JSONDATA.Cover)) AS COVER, COUNT(*) AS N FROM POOLDATA.M_COVER_LIFE m
 WHERE m.JSONDATA.Cover IS NOT NULL
 GROUP BY UPPER(TRIM(m.JSONDATA.Cover)) HAVING COUNT(*) > 1;
PROMPT F. Rumus Pega untuk nomor berikutnya (LAST_NUMBER) - ID belum terpakai (TERPAKAI = 0)
SELECT '1' || LPAD(TO_CHAR(s.LAST_NUMBER), 5, '0') AS ID_BERIKUTNYA,
       (SELECT COUNT(*) FROM POOLDATA.M_COVER_LIFE c WHERE c.ID = '1' || LPAD(TO_CHAR(s.LAST_NUMBER), 5, '0')) AS TERPAKAI
  FROM SYS.ALL_SEQUENCES s WHERE s.SEQUENCE_OWNER = 'POOLDATA' AND s.SEQUENCE_NAME = 'M_COVER_LIFE_SEQ';
SPOOL OFF
-- Catatan: 086 tidak dapat memakai DECODE / JSON_VALUE / JSON_EXISTS dengan literal karena teks EXECUTE IMMEDIATE blok
-- berpelindung tidak boleh memuat tanda kutip; bentuk tiga cabang OR per kolom di 086 setara dengan DECODE(...) = 0 di
-- kueri C.

-- LANGKAH-WO-DISEASELIFE.md (a) - BACA-SAJA, SEBELUM D2 dan SEBELUM -migrate. BUKTI keputusan D1-D3 TANPA menulis.
-- Hasil yang diharapkan (fakta WO 08-10-2026):
--   A. ID_ANGKA_TERTINGGI 197585 -> SEQ_DISEASE_LIFE akan mulai AWAL_SEQ 197586 (080); M_DISEASE_LIFE_SEQ 2052;
--   B. rumus Pega '1' || LPAD(seq, 5, '0') untuk nomor 2051-2056 (-> 102051-102056): SEMUA sudah dipakai (TERPAKAI >= 1;
--      102051 = 2 sebelum D2) - bukti bahwa rumus Pega tidak dapat dipakai (D1.1);
--   C. ID kembar: SATU kelompok 102051 x2 sebelum D2 (baris impor C718 dan baris uji TEST123 / Sakit); NOL sesudah D2 -
--      inilah yang akan menghentikan 081 (ORA-02437) bila D2 belum dijalankan;
--   D. ID NULL: NOL (bila > 0, 081 berhenti ORA-01449 - laporkan ke WO, jangan dihapus sendiri);
--   E. ICD Code kembar tanpa beda huruf: NOL kelompok (dasar penolakan kembar D3);
--   F. baris berhuruf kecil (dasar huruf besar D3): hanya baris uji 102051 TEST123 / Sakit;
--   G. baris uji yang akan dihapus D2: tepat SATU baris; pemakai teksnya di tabel klaim (pembanding: NOL - tabel klaim
--      menyimpan teks, bukan ID; bila > 0 dilaporkan saja);
--   H. panjang maksimum ID / ICD_CODE / DISEASE (pembanding 6 / 7 / 290 byte) - di bawah lebar kolom 100 / 100 / 1000.
-- Dijalankan sebagai BERKAS (@...); keluaran ke dis_bukti.txt - kirim ke executor / laporan (butir yang menunggu WO).
-- Satu angka berbeda = BERHENTI, jangan D2 / -migrate.
SET LINESIZE 250 PAGESIZE 200 NUMWIDTH 12
COLUMN ICD_CODE FORMAT A12
COLUMN DISEASE FORMAT A60
SPOOL dis_bukti.txt
PROMPT A. ID angka tertinggi dan awal SEQ_DISEASE_LIFE (080), M_DISEASE_LIFE_SEQ
SELECT MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))) AS ID_ANGKA_TERTINGGI,
       NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0) + 1 AS AWAL_SEQ
  FROM POOLDATA.DISEASE_LIFE;
SELECT LAST_NUMBER FROM SYS.ALL_SEQUENCES WHERE SEQUENCE_OWNER = 'POOLDATA' AND SEQUENCE_NAME = 'M_DISEASE_LIFE_SEQ';
PROMPT B. Rumus Pega untuk nomor LAST_NUMBER - 1 sampai + 4: ID yang dihasilkan sudah terpakai? (TERPAKAI >= 1 di keenamnya)
SELECT k.ID_PEGA, (SELECT COUNT(*) FROM POOLDATA.DISEASE_LIFE d WHERE d.ID = k.ID_PEGA) AS TERPAKAI
  FROM (SELECT '1' || LPAD(TO_CHAR(q.LAST_NUMBER + n.N - 2), 5, '0') AS ID_PEGA
          FROM (SELECT s.LAST_NUMBER FROM SYS.ALL_SEQUENCES s
                 WHERE s.SEQUENCE_OWNER = 'POOLDATA' AND s.SEQUENCE_NAME = 'M_DISEASE_LIFE_SEQ') q
         CROSS JOIN (SELECT LEVEL AS N FROM DUAL CONNECT BY LEVEL <= 6) n) k
 ORDER BY k.ID_PEGA;
PROMPT C. ID kembar (sebelum D2: 102051 x2; sesudah D2: NOL baris) beserta isinya
SELECT d.ID, d.ICD_CODE, d.DISEASE FROM POOLDATA.DISEASE_LIFE d
 WHERE d.ID IN (SELECT x.ID FROM POOLDATA.DISEASE_LIFE x GROUP BY x.ID HAVING COUNT(*) > 1)
 ORDER BY d.ID, d.ICD_CODE;
PROMPT D. ID NULL (harus NOL)
SELECT COUNT(*) AS N_ID_NULL FROM POOLDATA.DISEASE_LIFE WHERE ID IS NULL;
PROMPT E. ICD Code kembar tanpa beda huruf dan spasi tepi (harus NOL baris)
SELECT UPPER(TRIM(ICD_CODE)) AS ICD, COUNT(*) AS N FROM POOLDATA.DISEASE_LIFE
 GROUP BY UPPER(TRIM(ICD_CODE)) HAVING COUNT(*) > 1;
PROMPT F. Baris berhuruf kecil di ICD_CODE / DISEASE (pembanding: hanya 102051 TEST123 / Sakit)
SELECT ID, ICD_CODE, DISEASE FROM POOLDATA.DISEASE_LIFE
 WHERE ICD_CODE <> UPPER(ICD_CODE) OR DISEASE <> UPPER(DISEASE) ORDER BY ID;
PROMPT G. Baris uji D2 (tepat SATU) dan pemakai teksnya di tabel klaim (pembanding NOL)
SELECT COUNT(*) AS N_BARIS_UJI FROM POOLDATA.DISEASE_LIFE WHERE ID = '102051' AND ICD_CODE = 'TEST123' AND DISEASE = 'Sakit';
SELECT COUNT(*) AS N_DIAGNOSE_TEST123 FROM POOLDATA.T_CLAIMLF_DIAGNOSE WHERE ICD_CODE = 'TEST123';
SELECT COUNT(*) AS N_DETAIL_TEST123 FROM POOLDATA.T_CLAIMLF_PREMIUMLIST_DETAIL WHERE ICD_CODE = 'TEST123';
PROMPT H. Panjang maksimum (byte) - pembanding 6 / 7 / 290
SELECT MAX(LENGTHB(ID)) AS MAKS_ID, MAX(LENGTHB(ICD_CODE)) AS MAKS_ICD, MAX(LENGTHB(DISEASE)) AS MAKS_DISEASE
  FROM POOLDATA.DISEASE_LIFE;
SPOOL OFF
-- Catatan: kueri G tabel klaim - ORA-00942 / ORA-00904 di skema yang belum memasang claimlife = wajar; kueri A-F tetap
-- berlaku.

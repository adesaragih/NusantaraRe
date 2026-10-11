-- LANGKAH-WO-RIRISKLIFE.md (a) P1 - BACA-SAJA. Sesi yang masih memegang kunci DML pada kedua tabel R/I Risk, dan sesi
-- yang masih tersambung. Jalankan SESUDAH Pega R/I Risk dan backend :8080 dihentikan. Hasil yang diharapkan: kueri 1
-- NOL baris. Butuh hak SELECT pada V$LOCKED_OBJECT dan V$SESSION (ORA-00942 = tidak punya hak: DBA menjalankan
-- ririsk_a1_kunci_dba.sql). M_RIRISK_LIFE_TEMP tidak diperiksa (bukan bagian modul ini, tidak disentuh).
SET LINESIZE 250 PAGESIZE 100
COLUMN USERNAME FORMAT A20
COLUMN OSUSER FORMAT A20
COLUMN MACHINE FORMAT A30
COLUMN PROGRAM FORMAT A30
COLUMN MODULE FORMAT A30
COLUMN OBJECT_NAME FORMAT A24
PROMPT 1. Kunci DML pada tabel R/I Risk (harus nol baris)
SELECT s.SID, s.SERIAL#, s.USERNAME, s.OSUSER, s.MACHINE, s.PROGRAM, s.MODULE, s.STATUS, o.OBJECT_NAME, l.LOCKED_MODE
  FROM V$LOCKED_OBJECT l
  JOIN SYS.ALL_OBJECTS o ON o.OBJECT_ID = l.OBJECT_ID
  JOIN V$SESSION s ON s.SID = l.SESSION_ID
 WHERE o.OWNER = 'POOLDATA' AND o.OBJECT_NAME IN ('M_RIRISK_LIFE_SUMMARY', 'M_RIRISK_LIFE');
PROMPT 2. Sesi yang tersambung sebagai POOLDATA atau GL (kenali Pega / backend; sesi SQL*Plus ini sendiri ikut tampil)
SELECT SID, SERIAL#, USERNAME, OSUSER, MACHINE, PROGRAM, MODULE, STATUS, LOGON_TIME
  FROM V$SESSION WHERE USERNAME IN ('POOLDATA', 'GL') ORDER BY LOGON_TIME;

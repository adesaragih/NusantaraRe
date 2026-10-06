-- DBA-LEPAS-VIEW-RICOMM_LIFE.sql - dijalankan work owner / DBA SEBELUM `cmd/api -migrate` yang membawa migrasi inti
-- 924 (tabel flat RICOMM_LIFE). BUKAN berkas migrasi: pelari tidak membacanya, aplikasi tidak menjalankannya.
--
-- Keputusan work owner 06-10-2026 butir 2: VIEW POOLDATA.RICOMM_LIFE harus dibuang lebih dulu. Selama view itu ada,
-- pra-terbang pelari (`praTerbangBentuk` / `objekAda`, inti/backend/migrasi/migrasi.go) menghitung objek APA PUN di
-- ALL_OBJECTS: view bernama sama dengan kolom yang sama membuat CREATE TABLE dilewati (ORA-00955 "sudah ada") dan
-- langkah 924 tercatat selesai TANPA tabel. Karena itu DROP VIEW tidak boleh ada di run -migrate yang sama.
--
-- Fakta DEV (dicek work owner 06-10-2026): view RICOMM_LIFE = SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY,
-- a.JSONDATA.CONTRACT, a.JSONDATA.YEAR, a.JSONDATA.COMM FROM M_RICOMM_LIFE a; satu-satunya dependensinya M_RICOMM_LIFE;
-- nol objek lain yang memakainya; nol sinonim / grant; nol kode repo membacanya. M_RICOMM_LIFE TIDAK disentuh.
-- Jalur balik: migrasi inti 924_ricomm_life_down.sql membuat ulang view ini persis.
--
-- Langkah: jalankan bagian 1-3 di SQL*Plus / SQL Developer sebagai pemilik POOLDATA (atau DBA). Lanjut ke bagian 4
-- HANYA bila bagian 2 dan 3 menjawab nol baris.

-- 1. Catat definisi view sebelum dibuang (bandingkan dengan fakta DEV di atas).
SELECT OWNER, VIEW_NAME, TEXT FROM SYS.ALL_VIEWS WHERE OWNER = 'POOLDATA' AND VIEW_NAME = 'RICOMM_LIFE';

-- 2. Objek yang BERGANTUNG pada view ini - harus nol baris.
SELECT OWNER, NAME, TYPE FROM SYS.ALL_DEPENDENCIES
 WHERE REFERENCED_OWNER = 'POOLDATA' AND REFERENCED_NAME = 'RICOMM_LIFE' AND REFERENCED_TYPE = 'VIEW';

-- 3. Sinonim dan hak atas view ini - harus nol baris.
SELECT OWNER, SYNONYM_NAME FROM SYS.ALL_SYNONYMS WHERE TABLE_OWNER = 'POOLDATA' AND TABLE_NAME = 'RICOMM_LIFE';
SELECT GRANTEE, PRIVILEGE FROM SYS.ALL_TAB_PRIVS WHERE TABLE_SCHEMA = 'POOLDATA' AND TABLE_NAME = 'RICOMM_LIFE';

-- 4. Buang view (DDL Oracle menutup transaksinya sendiri; tanpa COMMIT).
DROP VIEW POOLDATA.RICOMM_LIFE;

-- 5. Periksa: nol objek bernama RICOMM_LIFE di POOLDATA. Sesudah ini baru `go run ./cmd/api -migrate`.
SELECT OBJECT_NAME, OBJECT_TYPE FROM SYS.ALL_OBJECTS WHERE OWNER = 'POOLDATA' AND OBJECT_NAME = 'RICOMM_LIFE';

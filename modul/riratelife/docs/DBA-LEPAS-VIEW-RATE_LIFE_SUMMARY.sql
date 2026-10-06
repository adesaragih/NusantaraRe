-- DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql - dijalankan work owner / DBA SEBELUM `cmd/api -migrate` yang membawa migrasi
-- inti 926 (tabel flat RATE_LIFE_SUMMARY). BUKAN berkas migrasi: pelari tidak membacanya, aplikasi tidak menjalankannya.
--
-- Keputusan work owner 06-10-2026 K-F1: VIEW POOLDATA.RATE_LIFE_SUMMARY diganti tabel flat bernama sama. Selama view itu
-- ada, pra-terbang pelari (`praTerbangBentuk` / `objekAda`, inti/backend/migrasi/migrasi.go) menganggap objek bernama sama
-- dengan kolom yang sama "sudah ada": CREATE TABLE dilewati (ORA-00955) dan CREATE INDEX sesudahnya gagal atas view
-- (ORA-01702), sehingga -migrate berhenti tanpa tabel. Urutan lengkap: docs/LANGKAH-WO-RIRATELIFE-FLAT.md.
--
-- Fakta DEV (ALL_VIEWS, dibaca WO 06-10-2026): view = SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE,
-- a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a. Pembaca di repo:
-- riratelife, mastercontractretrolife, masterproductnamelife (SELECT ID, USEDBY - tetap jalan atas tabel flat). Objek
-- database LAIN yang bergantung pada view ini BELUM diperiksa - bagian 2 wajib nol baris. M_RATE_LIFE_SUMMARY TIDAK
-- disentuh. Jalur balik: 926_rate_life_summary_flat_down.sql membuat ulang view ini persis.
--
-- Jalankan bagian 1-3 sebagai pemilik POOLDATA (atau DBA). Lanjut ke bagian 4 HANYA bila bagian 2 dan 3 nol baris;
-- bila tidak nol: BERHENTI dan laporkan ke work owner (objek itu akan rusak / perlu dibangun ulang).

-- 1. Catat definisi view sebelum dibuang (bandingkan dengan fakta DEV di atas) dan cacah sumbernya (339 di DEV).
SELECT OWNER, VIEW_NAME, TEXT FROM SYS.ALL_VIEWS WHERE OWNER = 'POOLDATA' AND VIEW_NAME = 'RATE_LIFE_SUMMARY';
SELECT COUNT(*) FROM POOLDATA.M_RATE_LIFE_SUMMARY;

-- 2. Objek yang BERGANTUNG pada view ini (view lain, prosedur, paket, MV) - harus nol baris.
SELECT OWNER, NAME, TYPE FROM SYS.ALL_DEPENDENCIES
 WHERE REFERENCED_OWNER = 'POOLDATA' AND REFERENCED_NAME = 'RATE_LIFE_SUMMARY' AND REFERENCED_TYPE = 'VIEW';

-- 3. Sinonim dan hak atas view ini - harus nol baris.
SELECT OWNER, SYNONYM_NAME FROM SYS.ALL_SYNONYMS WHERE TABLE_OWNER = 'POOLDATA' AND TABLE_NAME = 'RATE_LIFE_SUMMARY';
SELECT GRANTEE, PRIVILEGE FROM SYS.ALL_TAB_PRIVS WHERE TABLE_SCHEMA = 'POOLDATA' AND TABLE_NAME = 'RATE_LIFE_SUMMARY';

-- 4. Buang view (DDL Oracle menutup transaksinya sendiri; tanpa COMMIT).
DROP VIEW POOLDATA.RATE_LIFE_SUMMARY;

-- 5. Periksa: nol objek bernama RATE_LIFE_SUMMARY di POOLDATA. Sesudah ini baru `go run ./cmd/api -migrate`.
SELECT OBJECT_NAME, OBJECT_TYPE FROM SYS.ALL_OBJECTS WHERE OWNER = 'POOLDATA' AND OBJECT_NAME = 'RATE_LIFE_SUMMARY';

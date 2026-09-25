/* =====================================================================
   PREFLIGHT.sql  --  Claim Non Prop
   Tujuan : menjawab T1-T7 langsung dari kamus data Oracle.
   Klien  : TOAD for Oracle (Editor). TIDAK memakai perintah SQL*Plus.
   Aman   : hanya ALL_* / USER_* / SESSION_* / SYS_CONTEXT. Tidak menyentuh
            satu pun tabel bisnis. Tidak mengunci apa pun. Tidak ada DML.
            Boleh dijalankan di produksi, jam berapa pun.
   Cara   : klik di dalam satu blok, tekan F9. Tiap blok berdiri sendiri,
            tidak bergantung pada blok sebelumnya.
            Export Dataset -> CSV, tempel balik dengan label T1..T7.
   Catatan: FETCH FIRST sengaja TIDAK dipakai (butuh 12c). Selama T1 belum
            terjawab, semua pembatas memakai ROWNUM yang aman di semua versi.
   ===================================================================== */


/* =====================================================================
   T1 -- VERSI ORACLE
   Menentukan: notasi titik JSON, FETCH FIRST, search_condition_vc, LISTAGG.
   ===================================================================== */

-- T1a  (utama)
SELECT product, version, status
  FROM product_component_version
 WHERE product LIKE 'Oracle%';

-- T1b  (cadangan bila T1a kosong; butuh grant ke v$)
-- Gagal ORA-00942 di sini = tidak punya akses v$, BUKAN berarti versi tak diketahui.
SELECT banner
  FROM v$version
 WHERE ROWNUM <= 3;

-- T1c  (cadangan terakhir, tanpa hak akses apa pun)
-- Uji fungsional. JALAN -> 12c ke atas. ORA-00933 -> di bawah 12c.
SELECT 'FETCH FIRST DIDUKUNG -> 12c+' AS hasil_uji
  FROM dual
 ORDER BY 1
 FETCH FIRST 1 ROWS ONLY;


/* =====================================================================
   T2 -- TIPE KOLOM DATA_JSON / JSONDATA
   Menentukan: RECON.sql BAGIAN 6 dan BAGIAN 10 dipakai apa adanya
               (notasi titik) atau ditulis ulang dengan JSON_VALUE.
   ===================================================================== */

-- T2a  tipe kolomnya sendiri
SELECT owner,
       table_name,
       column_name,
       data_type,
       data_length,
       nullable
  FROM all_tab_columns
 WHERE column_name LIKE '%JSON%'
 ORDER BY owner, table_name, column_name;

-- T2b  adakah check constraint IS JSON  (butuh 12.2+, search_condition_vc bertipe VARCHAR2)
-- ORA-00904 "SEARCH_CONDITION_VC" -> database di bawah 12.2, lanjut ke T2c.
SELECT c.owner,
       c.table_name,
       c.constraint_name,
       SUBSTR(c.search_condition_vc, 1, 120) AS kondisi
  FROM all_constraints c
 WHERE c.constraint_type = 'C'
   AND UPPER(c.search_condition_vc) LIKE '%IS JSON%'
 ORDER BY c.owner, c.table_name;

-- T2c  cadangan untuk di bawah 12.2. search_condition bertipe LONG: tidak bisa
--      ditapis di WHERE dan terpotong di grid TOAD. Jadi blok ini hanya
--      MENDAFTAR kandidatnya; isi kondisinya dibuka manual lewat
--      TOAD Schema Browser -> tab Constraints, bukan lewat query.
SELECT c.owner,
       c.table_name,
       c.constraint_name,
       c.status
  FROM all_constraints c
 WHERE c.constraint_type = 'C'
   AND (c.owner, c.table_name) IN (
         SELECT owner, table_name
           FROM all_tab_columns
          WHERE column_name LIKE '%JSON%')
 ORDER BY c.owner, c.table_name, c.constraint_name;


/* =====================================================================
   T3 -- OWNER DARI NAMA TANPA PREFIX          <-- PALING MENENTUKAN
   Nama tanpa prefix di SQL Pega di-resolve terhadap schema koneksi PEGA,
   bukan schema sesi TOAD ini. Karena itu T3d (CURRENT_SCHEMA) WAJIB ikut
   ditempel -- tanpa itu hasil T3a tidak bisa ditafsirkan.
   Aturan: nama yang muncul di lebih dari satu owner ditandai AMBIGU dan
   TETAP ditulis "?" di PULL-LIST.csv. Jangan pilih salah satu.
   ===================================================================== */

-- T3a  ke-25 nama yang TERBUKTI dipakai di SQL produksi tanpa owner
WITH nama AS (
  SELECT 'AGENT' n FROM dual UNION ALL
  SELECT 'BANKACCOUNT' FROM dual UNION ALL
  SELECT 'BUSINESS' FROM dual UNION ALL
  SELECT 'CITY' FROM dual UNION ALL
  SELECT 'CLAIMREJECTED' FROM dual UNION ALL
  SELECT 'CURRENCY' FROM dual UNION ALL
  SELECT 'JSON_KLAIM' FROM dual UNION ALL
  SELECT 'JSON_POLIS' FROM dual UNION ALL
  SELECT 'M_CLIENT' FROM dual UNION ALL
  SELECT 'M_TREATY_IN' FROM dual UNION ALL
  SELECT 'M_TREATY_IN_EDM' FROM dual UNION ALL
  SELECT 'M_TREATY_OUT' FROM dual UNION ALL
  SELECT 'M_TREATY_OUT_DETAIL' FROM dual UNION ALL
  SELECT 'MARKETINGOFFICER' FROM dual UNION ALL
  SELECT 'OS_AKSEPTASI_KLAIM' FROM dual UNION ALL
  SELECT 'PLATNP_SEQ' FROM dual UNION ALL
  SELECT 'PROPORTIONALARRG' FROM dual UNION ALL
  SELECT 'REINSURANCETYPE' FROM dual UNION ALL
  SELECT 'RW' FROM dual UNION ALL
  SELECT 'TREATY_OUT' FROM dual UNION ALL
  SELECT 'TREATYBUSINESS' FROM dual UNION ALL
  SELECT 'TREATYCONTRACT' FROM dual UNION ALL
  SELECT 'TREATYINPRODUCTION' FROM dual UNION ALL
  SELECT 'T_STORAGE_IMAGE' FROM dual UNION ALL
  SELECT 'V_D_CAUSE_OF_LOSS_BUSINESS' FROM dual
)
SELECT n.n                                            AS nama_objek,
       COUNT(o.owner)                                 AS jml_terlihat,
       CASE WHEN COUNT(o.owner) = 0        THEN 'TIDAK TERLIHAT'
            WHEN COUNT(DISTINCT o.owner) > 1 THEN 'AMBIGU'
            ELSE 'UNIK' END                           AS putusan,
       LISTAGG(o.owner || '/' || o.object_type || '/' || o.status, ' ; ')
         WITHIN GROUP (ORDER BY o.owner)              AS daftar_owner
  FROM nama n
  LEFT JOIN all_objects o
         ON o.object_name = n.n
 GROUP BY n.n
 ORDER BY 3, 1;

-- T3b  ke-11 nama yang HANYA ditebak dari nama class Pega, belum pernah terlihat
--      di satu pun pernyataan SQL. "TIDAK TERLIHAT" di sini berarti tebakan
--      namanya salah ATAU objeknya tidak di-grant -- dibedakan lewat T4c.
--      Nama sebenarnya nanti dijawab oleh pemetaan Data-Admin-DB-Table (REQ-012).
WITH nama AS (
  SELECT 'EMAILKOMITE' n FROM dual UNION ALL
  SELECT 'V_POLIS' FROM dual UNION ALL
  SELECT 'V_M_CAUSE_OF_LOSS' FROM dual UNION ALL
  SELECT 'V_D_CAUSE_OF_LOSS' FROM dual UNION ALL
  SELECT 'V_MST_USER_TEKNIS' FROM dual UNION ALL
  SELECT 'ADJUSTERCONSULTANT' FROM dual UNION ALL
  SELECT 'LST_BANK_GROUP' FROM dual UNION ALL
  SELECT 'CATASTROPHE' FROM dual UNION ALL
  SELECT 'PROVINCE' FROM dual UNION ALL
  SELECT 'TREATYGROUP' FROM dual UNION ALL
  SELECT 'CURRENCYSTANDARD' FROM dual
)
SELECT n.n                                            AS nama_ditebak,
       COUNT(o.owner)                                 AS jml_terlihat,
       CASE WHEN COUNT(o.owner) = 0        THEN 'TIDAK TERLIHAT'
            WHEN COUNT(DISTINCT o.owner) > 1 THEN 'AMBIGU'
            ELSE 'UNIK' END                           AS putusan,
       LISTAGG(o.owner || '/' || o.object_type, ' ; ')
         WITHIN GROUP (ORDER BY o.owner)              AS daftar_owner
  FROM nama n
  LEFT JOIN all_objects o
         ON o.object_name = n.n
 GROUP BY n.n
 ORDER BY 3, 1;

-- T3c  sinonim: inilah yang menjelaskan kenapa nama tanpa prefix bisa jalan.
--      Kolom DB_LINK terisi = objeknya bahkan tidak ada di database ini.
SELECT s.owner        AS pemilik_sinonim,
       s.synonym_name,
       s.table_owner  AS menunjuk_ke_owner,
       s.table_name   AS menunjuk_ke_objek,
       s.db_link
  FROM all_synonyms s
 WHERE s.synonym_name IN (
         'AGENT','BANKACCOUNT','BUSINESS','CITY','CLAIMREJECTED','CURRENCY',
         'JSON_KLAIM','JSON_POLIS','M_CLIENT','M_TREATY_IN','M_TREATY_IN_EDM',
         'M_TREATY_OUT','M_TREATY_OUT_DETAIL','MARKETINGOFFICER',
         'OS_AKSEPTASI_KLAIM','PLATNP_SEQ','PROPORTIONALARRG','REINSURANCETYPE',
         'RW','TREATY_OUT','TREATYBUSINESS','TREATYCONTRACT','TREATYINPRODUCTION',
         'T_STORAGE_IMAGE','V_D_CAUSE_OF_LOSS_BUSINESS','EMAILKOMITE','V_POLIS',
         'V_M_CAUSE_OF_LOSS','V_D_CAUSE_OF_LOSS','V_MST_USER_TEKNIS',
         'ADJUSTERCONSULTANT','LST_BANK_GROUP','CATASTROPHE','PROVINCE',
         'TREATYGROUP','CURRENCYSTANDARD')
 ORDER BY s.synonym_name, s.owner;

-- T3d  identitas sesi TOAD Anda. WAJIB ditempel bersama T3a dan T3b.
SELECT SYS_CONTEXT('USERENV','SESSION_USER')    AS login_sebagai,
       SYS_CONTEXT('USERENV','CURRENT_SCHEMA')  AS schema_default_sesi_ini,
       SYS_CONTEXT('USERENV','DB_NAME')         AS db_name,
       SYS_CONTEXT('USERENV','SERVICE_NAME')    AS service_name,
       SYS_CONTEXT('USERENV','INSTANCE_NAME')   AS instance_name,
       SYS_CONTEXT('USERENV','SERVER_HOST')     AS server_host
  FROM dual;

-- T3e  sanity check terhadap owner yang SUDAH terbaca literal di XML.
--      Kalau salah satu baris ini tidak muncul, seluruh asumsi owner goyah.
SELECT owner, object_name, object_type, status
  FROM all_objects
 WHERE (owner = 'POOLDATA'
        AND object_name IN ('CLAIMXOL','CLAIMXOL2','KODE_PRODUKSI',
                            'OS_AKSEPTASI_KLAIM','REINSURANCETYPE',
                            'T_FOLDER_IMAGE','TREATYINDETAIL',
                            'TREATYINDETAILEDM','DIRECTTOKASIR_LOG'))
    OR (owner = 'REINSURANCE' AND object_name = 'TRLOSS_DETAIL_T')
 ORDER BY owner, object_name;


/* =====================================================================
   T4 -- HAK AKSES SESI INI
   Menentukan arti "tidak ditemukan" di SELURUH hasil PREFLIGHT:
     T4c JALAN -> punya DBA_* -> "tidak ditemukan" berarti TIDAK ADA.
     T4c GAGAL -> ALL_* saja  -> "tidak ditemukan" berarti TIDAK TERLIHAT.
   Hanya pada kasus pertama sebuah objek boleh ditandai DEAD.
   ===================================================================== */

-- T4a  privilege yang relevan saja (daftar penuh terlalu panjang untuk satu layar)
SELECT privilege
  FROM session_privs
 WHERE privilege LIKE '%ANY%'
    OR privilege LIKE 'SELECT%'
    OR privilege LIKE '%CATALOG%'
 ORDER BY privilege;

-- T4b  role
SELECT role FROM session_roles ORDER BY role;

-- T4c  uji tegas akses DBA_*. Satu baris, tidak membaca data apa pun.
SELECT COUNT(*) AS dba_objects_terbaca
  FROM dba_objects
 WHERE ROWNUM <= 1;


/* =====================================================================
   T5 -- LETAK SCHEMA PegaRULES
   Tiga kemungkinan: (1) satu database owner lain; (2) instance terpisah;
   (3) diakses lewat database link. Jalur ketiga inilah yang T5b uji --
   jalur yang belum masuk hitungan saya sebelumnya.
   ===================================================================== */

-- T5a  apakah objek Pega ada di database ini, dan milik siapa
SELECT owner,
       object_type,
       COUNT(*)          AS jumlah,
       MIN(object_name)  AS contoh_nama
  FROM all_objects
 WHERE object_name LIKE 'PR4%'
    OR object_name LIKE 'PR!_%' ESCAPE '!'
 GROUP BY owner, object_type
 ORDER BY owner, object_type;

-- T5b  database link. Kalau ada yang mengarah ke instance Pega, REQ-001
--      harus diajukan ke DBA yang berbeda.
SELECT owner, db_link, username, host, created
  FROM all_db_links
 ORDER BY owner, db_link;

-- T5c  cadangan bila T5b ditolak
SELECT db_link, username, host, created FROM user_db_links ORDER BY db_link;


/* =====================================================================
   T6 -- [DICABUT]
   Perilaku POOLDATA.GETCURRENCYSTANDARD ketika kurs tidak tersedia BUKAN
   pertanyaan untuk Anda -- jawabannya ada di body function-nya sendiri.
   Dipindahkan ke REQ-003, diambil lewat ALL_SOURCE di RECON.sql BAGIAN 5A.
   Nomor T6 sengaja dikosongkan supaya penomoran T1..T7 tetap sejajar.
   ===================================================================== */


/* =====================================================================
   T7 -- UKURAN TABEL
   Menentukan mana yang boleh di-full-scan dan mana yang harus di-sampling
   di RECON.sql BAGIAN 6, 7, dan 10. Dengan jawaban T8 (tidak ada UAT),
   setiap profil berat akan jalan di produksi -- jadi angka ini mengikat.
   PERINGATAN: num_rows berasal dari statistik optimizer dan bisa basi.
   Kolom MUTU_ANGKA menandai mana yang tidak boleh dipercaya.
   ===================================================================== */

-- T7a  perkiraan baris + kesegaran statistik
SELECT owner,
       table_name,
       num_rows,
       last_analyzed,
       CASE WHEN last_analyzed IS NULL
              THEN 'TANPA STATISTIK - ANGKA TIDAK DAPAT DIPAKAI'
            WHEN last_analyzed < ADD_MONTHS(SYSDATE, -12)
              THEN 'BASI LEBIH DARI 1 TAHUN - JANGAN DIPERCAYA'
            WHEN last_analyzed < ADD_MONTHS(SYSDATE, -3)
              THEN 'AGAK BASI'
            ELSE 'SEGAR' END AS mutu_angka
  FROM all_tables
 WHERE table_name IN (
         'OS_AKSEPTASI_KLAIM','JSON_KLAIM','JSON_POLIS','TREATY_OUT',
         'CLAIMXOL','CLAIMXOL2','TREATYINPRODUCTION','TREATYINDETAIL',
         'TREATYINDETAILEDM','M_TREATY_IN','M_TREATY_IN_EDM','M_TREATY_OUT',
         'M_TREATY_OUT_DETAIL','PROPORTIONALARRG','TREATYBUSINESS',
         'TREATYCONTRACT','CLAIMREJECTED','BANKACCOUNT','M_CLIENT','AGENT',
         'MARKETINGOFFICER','CURRENCY','BUSINESS','RW','CITY',
         'REINSURANCETYPE','KODE_PRODUKSI','T_FOLDER_IMAGE','T_STORAGE_IMAGE',
         'DIRECTTOKASIR_LOG','TRLOSS_DETAIL_T','EMAILKOMITE')
 ORDER BY num_rows DESC NULLS LAST;

-- T7b  ukuran segmen nyata, tidak bergantung statistik. Hanya schema Anda sendiri.
SELECT segment_name,
       segment_type,
       ROUND(SUM(bytes)/1024/1024) AS mb
  FROM user_segments
 WHERE segment_type IN ('TABLE','TABLE PARTITION','LOBSEGMENT')
 GROUP BY segment_name, segment_type
HAVING SUM(bytes) > 100*1024*1024
 ORDER BY 3 DESC;

-- T7c  cadangan: ukuran segmen lintas owner. Butuh DBA_*.
--      Kalau T4c gagal, blok ini pasti gagal juga -- dan itu berarti kita
--      hanya punya num_rows dari T7a, dengan segala keterbatasannya.
SELECT owner,
       segment_name,
       segment_type,
       ROUND(SUM(bytes)/1024/1024) AS mb
  FROM dba_segments
 WHERE owner IN ('POOLDATA','REINSURANCE')
   AND segment_type IN ('TABLE','TABLE PARTITION','LOBSEGMENT')
 GROUP BY owner, segment_name, segment_type
HAVING SUM(bytes) > 100*1024*1024
 ORDER BY 4 DESC;


/* =====================================================================
   T8 -- SUDAH DIJAWAB, TIDAK ADA QUERY
   "Anggap tidak ada UAT yang datanya representatif sampai saya bilang
   sebaliknya." Konsekuensi yang mengikat penulisan ulang RECON.sql:
     - tanpa hint paralel
     - tanpa full scan pada tabel yang T7 tandai besar
     - tanpa SELECT FOR UPDATE atau lock apa pun
     - query profil dipecah bertahap, bukan sekali jalan berjam-jam
   ===================================================================== */

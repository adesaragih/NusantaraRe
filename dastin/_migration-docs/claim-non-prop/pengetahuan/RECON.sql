-- =====================================================================
-- RECON.sql  --  Claim Non Prop  --  migrasi Pega -> Oracle/Go/React
-- Tujuan : memulihkan tipe & presisi properti, memetakan objek Oracle.
-- Aturan : READ-ONLY. Tidak ada DML. Sampel maks 20 baris.
--          Kolom teks bebas di-masking. Tanpa data nasabah.
-- Hasil  : tempel ke SCHEMA-ACTUAL.csv
-- =====================================================================

SET PAGESIZE 200
SET LINESIZE 300
SET LONG 2000000
SET LONGCHUNKSIZE 2000000


-- =====================================================================
-- BAGIAN 0 -- KONTEKS: siapa saya, ada schema apa saja
-- =====================================================================
SELECT USER AS current_user, SYS_CONTEXT('USERENV','DB_NAME') AS db_name FROM dual;

SELECT owner, COUNT(*) AS jml_objek
FROM   all_objects
WHERE  owner NOT IN ('SYS','SYSTEM','XDB','OUTLN','DBSNMP','APPQOSSYS','CTXSYS','MDSYS','ORDSYS','WMSYS')
GROUP  BY owner
ORDER  BY jml_objek DESC;


-- =====================================================================
-- BAGIAN 1 -- DEFINISI PROPERTY PEGA  (PRIORITAS TERTINGGI)
-- Struktur tabel rule berbeda antar versi Pega.
-- Jalankan 1A, 1B, 1C berurutan. Pakai yang mengembalikan baris.
-- =====================================================================

-- 1A. Temukan dulu tabel rule mana yang ada
SELECT owner, table_name
FROM   all_tables
WHERE  UPPER(table_name) LIKE 'PR4%'
    OR UPPER(table_name) LIKE 'PR_%'
ORDER  BY owner, table_name;

-- 1B. VARIAN 1 -- Pega 7/8 umumnya: definisi property ada di PR4_BASE
--     (pyPropertyType / pyMaxLength tersimpan sebagai kolom exposed)
SELECT  b.pzinskey,
        b.pyclassname,
        b.pypropertyname,
        b.pypropertytype,
        b.pymaxlength,
        b.pyruleset,
        b.pyruleavailable
FROM    pr4_base b
WHERE   UPPER(b.pyclassname) IN (
          'ASM-FW-GCNMFW-WORK-CLAIMTREATYNONPROP',
          'ASM-FW-GCNMFW-WORK-CLAIMTREATY',
          'ASM-FW-GCNMFW-WORK',
          'ASM-FW-GCNMFW-DATA-CLAIMDATA',
          'ASM-FW-GCNMFW-DATA-ADJUSTMENT',
          'ASM-FW-GCNMFW-DATA-COMITEE',
          'ASM-FW-GCNMFW-DATA-CLAIMRECEIVER',
          'ASM-FW-GCNMFW-DATA-OBJECTITEM',
          'ASM-FW-GCNMFW-DATA-CURRENCY',
          'ASM-FW-GISFW-DATA-SPREADINGRISK',
          'ASM-FW-GISFW-DATA-POLICYTREATYIN',
          'ASM-FW-GISFW-DATA-TREATYINLIMITS',
          'ASM-FW-GISFW-DATA-TREATYINLIMITLAYER',
          'ASM-FW-GISFW-DATA-QUOTATION',
          'ASM-FW-GISFW-DATA-TREATYININSTALLMENT'
        )
ORDER   BY b.pyclassname, b.pypropertyname;

-- 1C. VARIAN 2 -- bila 1B gagal: property tersimpan di PR4_RULE_VW / PR4_RULE
SELECT  r.pyclassname,
        r.pyrulename        AS property_name,
        r.pyruleset,
        r.pyruleavailable,
        r.pxobjclass
FROM    pr4_rule_vw r
WHERE   UPPER(r.pxobjclass) = 'RULE-OBJ-PROPERTY'
AND     UPPER(r.pyclassname) LIKE 'ASM-FW-GCNMFW%'
ORDER   BY r.pyclassname, property_name;

-- 1D. VARIAN 3 -- bila 1B & 1C gagal: definisi hanya ada di dalam BLOB rule.
--     Query ini hanya MEMASTIKAN keberadaannya, tidak membongkar BLOB.
SELECT  pyclassname,
        COUNT(*) AS jml_property
FROM    pr4_rule
WHERE   UPPER(pxobjclass) = 'RULE-OBJ-PROPERTY'
AND     UPPER(pyclassname) LIKE 'ASM-FW-%'
GROUP   BY pyclassname
ORDER   BY jml_property DESC;


-- =====================================================================
-- BAGIAN 2 -- PENEMUAN OBJEK LINTAS OWNER
-- Memverifikasi label CONFIRMED / DERIVED / GUESS.
-- Jangan hardcode schema: cari di semua owner yang terlihat.
-- =====================================================================
WITH dicari AS (
  SELECT 'OS_AKSEPTASI_KLAIM'          AS nama, 'CONFIRMED' AS label FROM dual UNION ALL
  SELECT 'JSON_KLAIM'                  , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'JSON_POLIS'                  , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'CLAIMXOL'                    , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'CLAIMXOL2'                   , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'CLAIMREJECTED'               , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATYINPRODUCTION'          , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATY_OUT'                  , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATYINDETAIL'              , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATYINDETAILEDM'           , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATYBUSINESS'              , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATYYEAR'                  , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATY_IN'                   , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TREATYCONTRACT'              , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'PROPORTIONALARRG'            , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'M_TREATY_IN'                 , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'M_TREATY_IN_EDM'             , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'M_TREATY_OUT'                , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'M_TREATY_OUT_DETAIL'         , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'BANKACCOUNT'                 , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'M_CLIENT'                    , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'AGENT'                       , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'MARKETINGOFFICER'            , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'CURRENCY'                    , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'BUSINESS'                    , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'REINSURANCETYPE'             , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'RW'                          , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'CITY'                        , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'KODE_PRODUKSI'               , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'T_STORAGE_IMAGE'             , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'T_FOLDER_IMAGE'              , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'DIRECTTOKASIR_LOG'           , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'V_M_CAUSE_OF_LOSS'           , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'V_D_CAUSE_OF_LOSS'           , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'V_D_CAUSE_OF_LOSS_BUSINESS'  , 'CONFIRMED' FROM dual UNION ALL
  SELECT 'TRLOSS_DETAIL_T'             , 'CONFIRMED' FROM dual UNION ALL
  -- DERIVED: dari nama class -Int-, belum pernah terlihat literal
  SELECT 'EMAILKOMITE'                 , 'DERIVED'   FROM dual UNION ALL
  SELECT 'ADJUSTERCONSULTANT'          , 'DERIVED'   FROM dual UNION ALL
  SELECT 'LST_BANK_GROUP'              , 'DERIVED'   FROM dual UNION ALL
  SELECT 'CATASTROPHE'                 , 'DERIVED'   FROM dual UNION ALL
  SELECT 'V_MST_USER_TEKNIS'           , 'DERIVED'   FROM dual UNION ALL
  SELECT 'PROVINCE'                    , 'DERIVED'   FROM dual UNION ALL
  SELECT 'TREATYGROUP'                 , 'DERIVED'   FROM dual UNION ALL
  SELECT 'CURRENCYSTANDARD'            , 'DERIVED'   FROM dual
)
SELECT  d.label,
        d.nama                AS dicari,
        NVL(o.owner,'>>> TIDAK DITEMUKAN <<<') AS owner,
        o.object_type,
        o.status
FROM    dicari d
LEFT JOIN all_objects o
       ON UPPER(o.object_name) = d.nama
      AND o.object_type IN ('TABLE','VIEW','SYNONYM','MATERIALIZED VIEW')
ORDER   BY d.label, d.nama, o.owner;

-- 2B. Sinonim yang menyembunyikan owner sebenarnya
SELECT owner, synonym_name, table_owner, table_name
FROM   all_synonyms
WHERE  UPPER(synonym_name) IN
       ('OS_AKSEPTASI_KLAIM','JSON_KLAIM','JSON_POLIS','TREATY_OUT','BANKACCOUNT','EMAILKOMITE','CLAIMXOL2')
ORDER  BY synonym_name;

-- 2C. Tabel work Pega -- GUESS, harus ditemukan
SELECT owner, table_name
FROM   all_tables
WHERE  UPPER(table_name) LIKE '%GCNMFW%'
    OR UPPER(table_name) LIKE 'PC_%'
ORDER  BY owner, table_name;


-- =====================================================================
-- BAGIAN 3 -- KOLOM + TIPE + PRESISI  (isi utama SCHEMA-ACTUAL.csv)
-- =====================================================================
SELECT  c.owner,
        c.table_name,
        c.column_id,
        c.column_name,
        c.data_type,
        c.data_length,
        c.data_precision,
        c.data_scale,
        c.nullable,
        SUBSTR(c.data_default,1,60) AS data_default
FROM    all_tab_columns c
WHERE   UPPER(c.table_name) IN (
          'OS_AKSEPTASI_KLAIM','JSON_KLAIM','JSON_POLIS','CLAIMXOL','CLAIMXOL2','CLAIMREJECTED',
          'TREATYINPRODUCTION','TREATY_OUT','TREATYINDETAIL','TREATYINDETAILEDM','TREATYBUSINESS',
          'TREATYYEAR','TREATY_IN','TREATYCONTRACT','PROPORTIONALARRG',
          'M_TREATY_IN','M_TREATY_IN_EDM','M_TREATY_OUT','M_TREATY_OUT_DETAIL',
          'BANKACCOUNT','M_CLIENT','AGENT','MARKETINGOFFICER','CURRENCY','BUSINESS',
          'REINSURANCETYPE','RW','CITY','KODE_PRODUKSI','T_STORAGE_IMAGE','T_FOLDER_IMAGE',
          'DIRECTTOKASIR_LOG','EMAILKOMITE','ADJUSTERCONSULTANT','LST_BANK_GROUP','CATASTROPHE',
          'V_MST_USER_TEKNIS','PROVINCE','TREATYGROUP','CURRENCYSTANDARD','TRLOSS_DETAIL_T',
          'V_M_CAUSE_OF_LOSS','V_D_CAUSE_OF_LOSS','V_D_CAUSE_OF_LOSS_BUSINESS'
        )
ORDER   BY c.owner, c.table_name, c.column_id;

-- 3B. DDL utuh untuk tabel paling kritis
--     (jalankan per tabel; ganti owner sesuai hasil BAGIAN 2)
-- SELECT DBMS_METADATA.GET_DDL('TABLE','OS_AKSEPTASI_KLAIM','POOLDATA') FROM dual;
-- SELECT DBMS_METADATA.GET_DDL('TABLE','TREATY_OUT','POOLDATA')         FROM dual;
-- SELECT DBMS_METADATA.GET_DDL('TABLE','EMAILKOMITE','POOLDATA')        FROM dual;
-- SELECT DBMS_METADATA.GET_DDL('TABLE','CLAIMXOL2','POOLDATA')          FROM dual;


-- =====================================================================
-- BAGIAN 4 -- CONSTRAINT & INDEX (natural key, ADR-0001)
-- =====================================================================
SELECT  ac.owner, ac.table_name, ac.constraint_name, ac.constraint_type,
        acc.position, acc.column_name, ac.r_constraint_name
FROM    all_constraints ac
JOIN    all_cons_columns acc
     ON acc.owner = ac.owner AND acc.constraint_name = ac.constraint_name
WHERE   ac.constraint_type IN ('P','U','R')
AND     UPPER(ac.table_name) IN
        ('OS_AKSEPTASI_KLAIM','JSON_KLAIM','JSON_POLIS','CLAIMXOL2','CLAIMREJECTED',
         'TREATY_OUT','EMAILKOMITE','BANKACCOUNT','TREATYINPRODUCTION')
ORDER   BY ac.table_name, ac.constraint_type, acc.position;

SELECT  ai.owner, ai.table_name, ai.index_name, ai.uniqueness,
        aic.column_position, aic.column_name
FROM    all_indexes ai
JOIN    all_ind_columns aic
     ON aic.index_owner = ai.owner AND aic.index_name = ai.index_name
WHERE   UPPER(ai.table_name) IN
        ('OS_AKSEPTASI_KLAIM','JSON_KLAIM','CLAIMXOL2','TREATY_OUT','EMAILKOMITE','TREATYINPRODUCTION')
ORDER   BY ai.table_name, ai.index_name, aic.column_position;


-- =====================================================================
-- BAGIAN 5 -- SOURCE PROCEDURE / FUNCTION
-- Di sinilah pemetaan parameter CARIn -> kolom sebenarnya berada.
-- =====================================================================
SELECT owner, name, type, line, text
FROM   all_source
WHERE  UPPER(name) IN (
         'PEGA_JSON_OS_AKSEP_KLAIMTNP',
         'PEGA_JSON_KLAIM_PNC',
         'PROC_GENERATE_SEQUENCE_NUMBER',
         'GET_TOKEN_STORAGE',
         'PEGA_M_CAUSE_OF_LOSS',
         'PEGA_D_CAUSE_OF_LOSS',
         'GETCURRENCYSTANDARD',
         'F_GET_EMAIL'
       )
ORDER  BY owner, name, type, line;

-- 5B. Tanda tangan parameter (urutan, arah, tipe)
SELECT  owner, package_name, object_name, argument_name, position,
        in_out, data_type, data_length, data_precision, data_scale
FROM    all_arguments
WHERE   UPPER(object_name) IN (
          'PEGA_JSON_OS_AKSEP_KLAIMTNP','PEGA_JSON_KLAIM_PNC',
          'PROC_GENERATE_SEQUENCE_NUMBER','GET_TOKEN_STORAGE','GETCURRENCYSTANDARD'
        )
ORDER   BY object_name, position;

-- 5C. Definisi view
SELECT owner, view_name, text
FROM   all_views
WHERE  UPPER(view_name) IN
       ('V_M_CAUSE_OF_LOSS','V_D_CAUSE_OF_LOSS','V_D_CAUSE_OF_LOSS_BUSINESS','V_MST_USER_TEKNIS');


-- =====================================================================
-- BAGIAN 6 -- PROFIL PRESISI RIIL  (menjawab ADR-0003 / OPEN-QUESTIONS B2)
-- Deklarasi kolom tidak membuktikan presisi yang benar-benar dipakai.
-- Yang menentukan adalah isinya.
-- CATATAN: ganti nama kolom sesuai hasil BAGIAN 3 bila berbeda.
-- =====================================================================

-- 6A. Berapa desimal yang benar-benar terpakai pada kolom uang?
--     Pola query -- ulangi per kolom MUST-CONFIRM.
SELECT  'OS_AKSEPTASI_KLAIM'                                   AS tabel,
        'DATA_JSON.GrossValue'                                 AS kolom,
        COUNT(*)                                               AS jml_baris,
        MAX(ABS(nilai))                                        AS nilai_maks,
        MAX(  CASE WHEN INSTR(TO_CHAR(nilai),'.') = 0 THEN 0
                   ELSE LENGTH(TO_CHAR(nilai)) - INSTR(TO_CHAR(nilai),'.') END ) AS desimal_maks,
        SUM(  CASE WHEN INSTR(TO_CHAR(nilai),'.') > 0
                    AND LENGTH(TO_CHAR(nilai)) - INSTR(TO_CHAR(nilai),'.') > 2
                   THEN 1 ELSE 0 END )                         AS baris_desimal_gt2
FROM   ( SELECT TO_NUMBER(a.data_json.GrossValue) AS nilai
         FROM   os_akseptasi_klaim a
         WHERE  a.data_json.GrossValue IS NOT NULL
           AND  ROWNUM <= 100000 );

-- 6B. Presisi kurs -- paling menentukan karena mengalikan semuanya
SELECT  MAX(ABS(conversion))                                    AS kurs_maks,
        MAX(  CASE WHEN INSTR(TO_CHAR(conversion),'.') = 0 THEN 0
                   ELSE LENGTH(TO_CHAR(conversion)) - INSTR(TO_CHAR(conversion),'.') END ) AS desimal_maks
FROM    treaty_out
WHERE   conversion IS NOT NULL;

-- 6C. Presisi Limit Layer & Retensi Cedant
SELECT  MAX(ABS("LIMIT"))   AS limit_maks,
        MAX(ABS(limit2))    AS limit2_maks,
        MAX(ABS(deductible))AS deductible_maks,
        MAX(  CASE WHEN INSTR(TO_CHAR("LIMIT"),'.') = 0 THEN 0
                   ELSE LENGTH(TO_CHAR("LIMIT")) - INSTR(TO_CHAR("LIMIT"),'.') END ) AS limit_desimal_maks
FROM    treaty_out;

-- 6D. Panjang Nomor Akseptasi -- XML memvalidasi panjang 23 atau 24 untuk non-prop
SELECT  LENGTH(noclaim) AS panjang, COUNT(*) AS jml
FROM    os_akseptasi_klaim
WHERE   noclaim IS NOT NULL
GROUP   BY LENGTH(noclaim)
ORDER   BY panjang;


-- =====================================================================
-- BAGIAN 7 -- VOLUME & KARDINALITAS (menentukan indeks, OPEN-QUESTIONS B7)
-- =====================================================================
SELECT 'OS_AKSEPTASI_KLAIM' AS tabel, COUNT(*) AS jml_baris,
       MIN(tanggal) AS paling_lama, MAX(tanggal) AS paling_baru
FROM   os_akseptasi_klaim
UNION ALL
SELECT 'JSON_KLAIM', COUNT(*), MIN(tgl_input), MAX(tgl_input) FROM json_klaim
UNION ALL
SELECT 'CLAIMXOL2', COUNT(*), MIN(tanggal), MAX(tanggal) FROM claimxol2;

-- 7B. Konfigurasi jenjang Komite -- penentu KomiteLoop
SELECT  id, name, jabatan, degree, limit_bottom, limit_top,
        sts_aktif, sts_klaim, type_komite, operator_id
FROM    emailkomite
WHERE   UPPER(sts_klaim) = 'NONPROP'
ORDER   BY degree, limit_bottom;


-- =====================================================================
-- BAGIAN 8 -- PEMAKAIAN JALUR CloseClaimMD  (OPEN-QUESTIONS A5)
-- Menentukan jalur tutup-langsung dimigrasi atau dibuang.
-- Ganti predikat bila penanda jalur tutup tersimpan di kolom lain.
-- =====================================================================
SELECT  TO_CHAR(tanggal,'YYYY') AS tahun,
        COUNT(*)                AS jml_klaim_ditutup
FROM    os_akseptasi_klaim
WHERE   sts_reject = 0
GROUP   BY TO_CHAR(tanggal,'YYYY')
ORDER   BY tahun;


-- =====================================================================
-- BAGIAN 9 -- SAMPEL TERMASKING (maks 20 baris, tanpa data nasabah)
-- Hanya untuk memahami BENTUK data, bukan isinya.
-- =====================================================================
SELECT  idpega,
        SUBSTR(nopolis,1,4) || '****'                 AS nopolis_masked,
        LENGTH(data_json)                             AS panjang_json,
        SUBSTR(TO_CHAR(data_json),1,400)              AS cuplikan_json
FROM    json_klaim
WHERE   ROWNUM <= 20;

SELECT  caseid,
        SUBSTR(noclaim,1,6) || '****'                 AS noclaim_masked,
        sts_reject,
        tanggal,
        LENGTH(data_json)                             AS panjang_json
FROM    os_akseptasi_klaim
WHERE   ROWNUM <= 20;


-- =====================================================================
-- BAGIAN 10 -- REQ-011: AMBANG KEWENANGAN KOMITE vs MATA UANG
-- Menutup FINDING-001. Sepenuhnya read-only, tanpa data nasabah.
-- Konteks: CreateChildKomiteCNP_Act langkah 12-13 membandingkan
--          Local.TotalValueAdjust terhadap 30.000.000 / 50.000.000
--          tanpa operasi konversi mata uang di sepanjang rantainya.
-- Yang diuji lebih dulu: apakah nilai itu memang sudah IDR.
-- =====================================================================

-- 10A. Mata uang apa saja yang benar-benar dipakai pada akseptasi klaim,
--      dan berapa banyak masing-masing. Bila hasilnya hanya IDR,
--      FINDING-001 gugur dan tidak perlu dilanjutkan.
SELECT  NVL(a.data_json.Currency, '(kosong)')  AS mata_uang,
        COUNT(*)                               AS jml_baris,
        MIN(a.tanggal)                         AS paling_lama,
        MAX(a.tanggal)                         AS paling_baru
FROM    os_akseptasi_klaim a
GROUP   BY NVL(a.data_json.Currency, '(kosong)')
ORDER   BY jml_baris DESC;

-- 10B. Sebaran nilai per mata uang terhadap kedua ambang.
--      Nilai dibandingkan apa adanya (tanpa konversi), meniru perilaku
--      sistem lama, supaya terlihat mana yang jatuh ke jenjang terendah.
SELECT  NVL(a.data_json.Currency,'(kosong)')   AS mata_uang,
        COUNT(*)                                AS jml,
        SUM(CASE WHEN TO_NUMBER(a.data_json.Value) <= 30000000    THEN 1 ELSE 0 END) AS dinilai_bawah_30jt,
        SUM(CASE WHEN TO_NUMBER(a.data_json.Value) >  30000000
                  AND TO_NUMBER(a.data_json.Value) <= 50000000    THEN 1 ELSE 0 END) AS dinilai_30jt_50jt,
        SUM(CASE WHEN TO_NUMBER(a.data_json.Value) >  50000000    THEN 1 ELSE 0 END) AS dinilai_atas_50jt,
        MAX(TO_NUMBER(a.data_json.Value))       AS nilai_maks_apa_adanya
FROM    os_akseptasi_klaim a
WHERE   a.data_json.Value IS NOT NULL
GROUP   BY NVL(a.data_json.Currency,'(kosong)')
ORDER   BY jml DESC;

-- 10C. Yang menentukan besaran paparan: baris non-IDR yang SETARA
--      di atas ambang setelah dikonversi, tetapi dinilai di bawahnya
--      ketika dibandingkan apa adanya.
--      Ganti getcurrencystandard bila sumber kurs berbeda.
SELECT  a.data_json.Currency                                        AS mata_uang,
        COUNT(*)                                                    AS jml_terdampak,
        MIN(a.tanggal)                                              AS paling_lama,
        MAX(a.tanggal)                                              AS paling_baru
FROM    os_akseptasi_klaim a
WHERE   a.data_json.Currency IS NOT NULL
AND     UPPER(a.data_json.Currency) <> 'IDR'
AND     TO_NUMBER(a.data_json.Value) <= 30000000
AND     TO_NUMBER(a.data_json.Value)
          * POOLDATA.getcurrencystandard(a.data_json.Currency, a.tanggal) > 30000000
GROUP   BY a.data_json.Currency
ORDER   BY jml_terdampak DESC;

-- 10D. Jenjang persetujuan yang benar-benar menangani, dibaca dari
--      riwayat akseptasi. Nama pengguna di-masking.
SELECT  h.workbasket,
        SUBSTR(h.username,1,2) || '****'  AS username_masked,
        h.status,
        COUNT(*)                          AS jml
FROM    historyakseptasipega h
GROUP   BY h.workbasket, SUBSTR(h.username,1,2) || '****', h.status
ORDER   BY jml DESC
FETCH FIRST 50 ROWS ONLY;

-- 10E. Menjawab batas klaim FINDING-001 sec.4.1 dan 4.3:
--      apakah properti ValueAdjustment benar-benar tersimpan,
--      dan dalam mata uang apa. Sampel dibatasi 20 baris.
SELECT  a.caseid,
        a.data_json.Currency            AS mata_uang,
        a.data_json.ValueAdjustment     AS value_adjustment,
        a.data_json.Value               AS value
FROM    os_akseptasi_klaim a
WHERE   a.data_json.ValueAdjustment IS NOT NULL
AND     ROWNUM <= 20;

-- Data uji SINTETIS tiga modul - GILIRAN-13 paket 3.
--
-- Untuk apa: menutup tiga titik buta uji layar (panduan uji bab 0) di
-- SKEMA UJI - satu polis per tahap PremiumList Life, satu klaim per tahap
-- Claim Life dengan baris adjustment, roster EMAILKOMITE sintetis, dan satu
-- baris adjustment yang siap diserahkan ke Komite.
--
-- ⛔ TIDAK dijalankan executor. Work owner memuatnya, SESUDAH `-migrate`
--    (termasuk 057 dan 059) berjalan di skema uji.
--
-- ⛔ Nol data orang, nol nomor polis nyata, nol kredensial, nol alamat
--    layanan. Seluruh pengenal `UJI-*`, surel `uji-…@contoh.invalid`
--    (`.invalid` domain tercadang RFC 2606 - surel ke sana tidak pernah
--    sampai). Dijaga `repository/datauji_test.go`.
--
-- CARA MEMUAT (SQL*Plus, SQLcl, atau SQL Developer "Run Script" / F5):
--   1. jalankan berkas ini utuh, dengan akun yang berhak menulis ke skema uji;
--   2. alatnya MENANYAKAN `skema_uji` sekali - ketik nama skema uji dari DBA;
--   3. periksa hasilnya, lalu tetapkan sendiri dengan perintah penetapan
--      transaksi alat Anda. Berkas ini SENGAJA tidak menetapkan DML-nya
--      (ADR-U-0029) - kapan data menetap adalah keputusan Anda.
--
-- ⛔ ADR-U-0033: SETIAP nama tabel berawalan skema - `&&skema_uji..` untuk
--    tabel aplikasi, `SYS.` untuk kamus data. Nol ketergantungan pada skema
--    bawaan sesi. Nama skema yang DIKETIK itu sekaligus pengakuan: padanan
--    terdekat `ORACLE_SKEMA_UJI=true`, yang tidak dapat dibaca SQL.
--
-- ⛔ PAGAR - padanan `config.PagarSkemaUji`, pagar yang sama dengan
--    `-migrate-down`. Blok 1 memeriksa ketiganya SEBELUM pernyataan pengubah
--    pertamanya:
--      a. skema yang MEMUAT `POOLDATA` ditolak (memuat, bukan sama persis -
--         POOLDATA_DEV dan POOLDATA2 skema warisan juga);
--      b. `T_MIGRASI` wajib ada dan mencatat 057 - skema ini skema aplikasi
--         yang sudah dimigrasi, bukan skema lain;
--      c. klaim/polis `UJI-*` belum pernah dimuat - memuat ulang = bongkar
--         skema uji (`-migrate-down`, lalu `-migrate`).
--    Blok 2 MENGULANG (a) sebelum INSERT pertamanya - alat yang meneruskan
--    skrip sesudah galat tetap berhenti di sana. (b) dan (c) di blok 2
--    dijamin bentuknya: blok itu tidak dapat dikompilasi tanpa tabel dan
--    kolom 057, dan pengenal `UJI-*` yang sudah ada ditolak kunci primernya -
--    batal utuh.
--    ⚠️ `IS_PEGA_PROD=false` tidak dapat dibaca SQL: memilih sambungan yang
--    benar tetap tanggung jawab Anda.
--
-- ⛔ Seluruh INSERT berada di SATU blok PL/SQL: gagal di mana pun berarti
--    blok itu dibatalkan utuh, nol baris tertinggal separuh.
--
-- ⚠️ EMAILKOMITE bukan tabel migrasi - aplikasi membacanya dari skema aktif
--    (`repository/roster.go`). Bila skema uji belum memilikinya, blok 1
--    membuat TIRUAN berkolom yang dibaca kode ditambah yang `[data DBA]`
--    catat (`ID` PK, `LIMIT_TOP`, `STS_REJECT VARCHAR2(15)`). Lebar kolom
--    teksnya meniru kolom tujuan salinannya di `T_KOMITE_KOMITELIST` (013),
--    BUKAN katalog. Bila DBA sudah menyalin tabel aslinya dan tabel itu punya
--    kolom wajib lain, INSERT roster gagal dengan ORA-01400 - dan blok 2
--    batal utuh.
--    ⚠️ Pembuatan tiruan itu DDL, dan DDL Oracle menetapkan transaksi yang
--    sedang terbuka. Karena itu ia berjalan di blok 1, SEBELUM DML apa pun;
--    transaksi Anda yang lain yang belum ditetapkan ikut menetap.
--    ⚠️ Tiruan dan baris rosternya BERTAHAN melewati `-migrate-down` - ia
--    bukan tabel migrasi. Blok 2 karena itu memakai ulang baris `UJI-EK-*`
--    yang sudah ada, alih-alih menyisipkannya lagi.

WHENEVER SQLERROR EXIT FAILURE ROLLBACK

-- Blok 1 - pagar, lalu tiruan EMAILKOMITE bila belum ada.
DECLARE
  v_skema VARCHAR2(128) := UPPER(TRIM('&&skema_uji'));
  v_ada   NUMBER;
BEGIN
  IF INSTR(UPPER(v_skema), 'POOLDATA') > 0 THEN
    RAISE_APPLICATION_ERROR(-20901, 'data_uji_tiga_modul.sql menolak skema ' || v_skema ||
      ': memuat nama skema warisan Pega POOLDATA. Pakai skema uji kosong dari DBA.');
  END IF;

  SELECT COUNT(*) INTO v_ada FROM SYS.ALL_TABLES
   WHERE OWNER = v_skema AND TABLE_NAME = 'T_MIGRASI';
  IF v_ada = 0 THEN
    RAISE_APPLICATION_ERROR(-20902, 'skema ' || v_skema ||
      ' belum dimigrasi (T_MIGRASI tidak ada): jalankan -migrate lebih dulu.');
  END IF;
  EXECUTE IMMEDIATE 'SELECT COUNT(*) FROM &&skema_uji..T_MIGRASI WHERE NAMA = :1'
    INTO v_ada USING '057_seq_work_polis_dan_flag_ongoing';
  IF v_ada = 0 THEN
    RAISE_APPLICATION_ERROR(-20903, 'migrasi 057 belum berjalan di skema ' || v_skema ||
      ': jalankan -migrate dari main.');
  END IF;

  EXECUTE IMMEDIATE 'SELECT (SELECT COUNT(*) FROM &&skema_uji..T_WORK_POLIS WHERE ID LIKE ''UJI-%'') +
                            (SELECT COUNT(*) FROM &&skema_uji..T_WORK_CLAIM WHERE ID LIKE ''UJI-%'') FROM SYS.DUAL'
    INTO v_ada;
  IF v_ada > 0 THEN
    RAISE_APPLICATION_ERROR(-20904, 'data uji UJI-* sudah dimuat di skema ' || v_skema ||
      ': untuk memuat ulang, bongkar skema uji (-migrate-down, lalu -migrate).');
  END IF;

  SELECT COUNT(*) INTO v_ada FROM SYS.ALL_TABLES
   WHERE OWNER = v_skema AND TABLE_NAME = 'EMAILKOMITE';
  IF v_ada = 0 THEN
    EXECUTE IMMEDIATE 'CREATE TABLE &&skema_uji..EMAILKOMITE (ID VARCHAR2(40) NOT NULL, OPERATOR_ID VARCHAR2(64), JABATAN VARCHAR2(140), EMAIL VARCHAR2(255), DEGREE NUMBER(5), LIMIT_BOTTOM INTEGER, LIMIT_TOP INTEGER, STS_AKTIF VARCHAR2(8), STS_KLAIM VARCHAR2(16), STS_REJECT VARCHAR2(15), CONSTRAINT PK_EMAILKOMITE_UJI PRIMARY KEY (ID))';
  END IF;
END;
/

-- Blok 2 - pagar (a) lagi, lalu SELURUH data.
DECLARE
  v_skema  VARCHAR2(128) := UPPER(TRIM('&&skema_uji'));
  v_roster NUMBER;
BEGIN
  IF INSTR(UPPER(v_skema), 'POOLDATA') > 0 THEN
    RAISE_APPLICATION_ERROR(-20901, 'data_uji_tiga_modul.sql menolak skema ' || v_skema ||
      ': memuat nama skema warisan Pega POOLDATA. Pakai skema uji kosong dari DBA.');
  END IF;

  -- ===================================================================
  -- PremiumList Life - satu polis per tahap (panduan bab 2 §1.2).
  --
  -- STATUS_WORK = tahap (`models.TahapPolis*`; bernama STATUS sampai 059,
  -- seragam T_WORK_CLAIM), POSITION = `Offer`/`Premium`
  -- (`models.Posisi*`), FLAG_ONGOING_POLICY "0"/"1" VERBATIM (057, butir bn).
  --
  -- ⭐ GILIRAN-14 butir bq: `Confirm` di Input Offer Life dirutekan dari
  -- bendera - "0" menutup (Offer), "1" memindah ke Input Premium Detail
  -- (Premium). Karena itu tahap penawaran punya TIGA kasus: UJI-PL-A ("0"),
  -- UJI-PL-E ("1"), dan UJI-PL-F (bendera KOSONG - kasus sebelum 057 - yang
  -- `Confirm`-nya dijawab 409, sebab Decision3 tidak punya konektor Decline).
  -- Header T_PREMIUM_LIST: ID = ID_PEGA = Case ID - Detail membaca `p.ID`,
  -- kotak masuk menggabung `p.ID_PEGA = w.ID`.
  -- ===================================================================
  INSERT INTO &&skema_uji..T_WORK_POLIS (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY)
  VALUES ('UJI-PL-A', 'LIFE', 'Offer', 'Input Offer Life', '0');
  INSERT INTO &&skema_uji..T_WORK_POLIS (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY)
  VALUES ('UJI-PL-E', 'LIFE', 'Offer', 'Input Offer Life', '1');
  INSERT INTO &&skema_uji..T_WORK_POLIS (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY)
  VALUES ('UJI-PL-F', 'LIFE', 'Offer', 'Input Offer Life', NULL);
  INSERT INTO &&skema_uji..T_WORK_POLIS (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY)
  VALUES ('UJI-PL-B', 'LIFE', 'Premium', 'Input Premium Detail', '1');
  -- Input Premium Summary tidak punya konektor masuk di flow (models,
  -- TahapPolisSummary) - hanya dapat dicapai lewat data seperti ini.
  INSERT INTO &&skema_uji..T_WORK_POLIS (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY)
  VALUES ('UJI-PL-C', 'LIFE', 'Premium', 'Input Premium Summary', '1');
  -- Kasus tertutup: uji negatif PremiumList, DAN polis tempat klaim Claim
  -- Life di bawah berpijak (UJI-POL-0001).
  INSERT INTO &&skema_uji..T_WORK_POLIS (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY)
  VALUES ('UJI-PL-D', 'LIFE', 'Premium', 'Resolved-Completed', '1');

  INSERT INTO &&skema_uji..T_PREMIUM_LIST (ID, ID_PEGA, NO_POLIS, BUSINESS_CODE, BUSINESS_NAME, CEDING_CO, CEDING_CO_NAME, DATE_RECEIVED, MARKETING_CODE, MARKETING_NAME, POLICY_HOLDER, POLICY_HOLDER_NAME, CREATE_OP_NAME, TYPE, TGL_INPUT, PRODUCT_NAME)
  VALUES ('UJI-PL-A', 'UJI-PL-A', 'UJI-POL-000A', 'L1', 'UJI-BISNIS-L1', 'UJI-CEDING-1', 'UJI-CEDING-1', DATE '2026-09-01', 'UJI-MO-1', 'UJI-MO-1', 'UJI-PEMEGANG-1', 'UJI-PEMEGANG-1', 'UJI-ADMIN', 'QP', DATE '2026-09-29', 'UJI-PRODUK-1');
  INSERT INTO &&skema_uji..T_PREMIUM_LIST (ID, ID_PEGA, NO_POLIS, BUSINESS_CODE, BUSINESS_NAME, CEDING_CO, CEDING_CO_NAME, DATE_RECEIVED, MARKETING_CODE, MARKETING_NAME, POLICY_HOLDER, POLICY_HOLDER_NAME, CREATE_OP_NAME, TYPE, TGL_INPUT, PRODUCT_NAME)
  VALUES ('UJI-PL-B', 'UJI-PL-B', 'UJI-POL-000B', 'L1', 'UJI-BISNIS-L1', 'UJI-CEDING-1', 'UJI-CEDING-1', DATE '2026-09-01', 'UJI-MO-1', 'UJI-MO-1', 'UJI-PEMEGANG-1', 'UJI-PEMEGANG-1', 'UJI-ADMIN', 'QP', DATE '2026-09-29', 'UJI-PRODUK-1');
  INSERT INTO &&skema_uji..T_PREMIUM_LIST (ID, ID_PEGA, NO_POLIS, BUSINESS_CODE, BUSINESS_NAME, CEDING_CO, CEDING_CO_NAME, DATE_RECEIVED, MARKETING_CODE, MARKETING_NAME, POLICY_HOLDER, POLICY_HOLDER_NAME, CREATE_OP_NAME, TYPE, TGL_INPUT, PRODUCT_NAME)
  VALUES ('UJI-PL-C', 'UJI-PL-C', 'UJI-POL-000C', 'L1', 'UJI-BISNIS-L1', 'UJI-CEDING-1', 'UJI-CEDING-1', DATE '2026-09-01', 'UJI-MO-1', 'UJI-MO-1', 'UJI-PEMEGANG-1', 'UJI-PEMEGANG-1', 'UJI-ADMIN', 'QP', DATE '2026-09-29', 'UJI-PRODUK-1');
  INSERT INTO &&skema_uji..T_PREMIUM_LIST (ID, ID_PEGA, NO_POLIS, BUSINESS_CODE, BUSINESS_NAME, CEDING_CO, CEDING_CO_NAME, DATE_RECEIVED, MARKETING_CODE, MARKETING_NAME, POLICY_HOLDER, POLICY_HOLDER_NAME, CREATE_OP_NAME, TYPE, TGL_INPUT, PRODUCT_NAME)
  VALUES ('UJI-PL-D', 'UJI-PL-D', 'UJI-POL-0001', 'L1', 'UJI-BISNIS-L1', 'UJI-CEDING-1', 'UJI-CEDING-1', DATE '2026-09-01', 'UJI-MO-1', 'UJI-MO-1', 'UJI-PEMEGANG-1', 'UJI-PEMEGANG-1', 'UJI-ADMIN', 'QP', DATE '2026-09-29', 'UJI-PRODUK-1');

  INSERT INTO &&skema_uji..T_PREMIUM_LIST (ID, ID_PEGA, NO_POLIS, BUSINESS_CODE, BUSINESS_NAME, CEDING_CO, CEDING_CO_NAME, DATE_RECEIVED, MARKETING_CODE, MARKETING_NAME, POLICY_HOLDER, POLICY_HOLDER_NAME, CREATE_OP_NAME, TYPE, TGL_INPUT, PRODUCT_NAME)
  VALUES ('UJI-PL-E', 'UJI-PL-E', 'UJI-POL-000E', 'L1', 'UJI-BISNIS-L1', 'UJI-CEDING-1', 'UJI-CEDING-1', DATE '2026-09-01', 'UJI-MO-1', 'UJI-MO-1', 'UJI-PEMEGANG-1', 'UJI-PEMEGANG-1', 'UJI-ADMIN', 'QP', DATE '2026-09-29', 'UJI-PRODUK-1');
  INSERT INTO &&skema_uji..T_PREMIUM_LIST (ID, ID_PEGA, NO_POLIS, BUSINESS_CODE, BUSINESS_NAME, CEDING_CO, CEDING_CO_NAME, DATE_RECEIVED, MARKETING_CODE, MARKETING_NAME, POLICY_HOLDER, POLICY_HOLDER_NAME, CREATE_OP_NAME, TYPE, TGL_INPUT, PRODUCT_NAME)
  VALUES ('UJI-PL-F', 'UJI-PL-F', 'UJI-POL-000F', 'L1', 'UJI-BISNIS-L1', 'UJI-CEDING-1', 'UJI-CEDING-1', DATE '2026-09-01', 'UJI-MO-1', 'UJI-MO-1', 'UJI-PEMEGANG-1', 'UJI-PEMEGANG-1', 'UJI-ADMIN', 'QP', DATE '2026-09-29', 'UJI-PRODUK-1');

  -- Dua peserta untuk tahap Summary - `Summary Premium Life` merekap uang
  -- peserta (`repository/polis_summary.go`); tanpa peserta rekapnya nol.
  INSERT INTO &&skema_uji..T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, ID_PEGA, PL_NUMBER, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO, NAME_OF_INSURED, SEX, DOB, CURRENCY, SUM_INSURED, SUM_REASURED, GROSS_PREMIUM, NET_PREMIUM, COMM)
  VALUES ('UJI-PLD-C-1', 'UJI-PL-C', 'UJI-PL-C', 'UJI-PLN-C', 'UJI-POL-000C', 'UJI-PEMEGANG-1', 'UJI-SERT-C1', 'UJI-TERTANGGUNG-C1', 'M', DATE '1980-01-01', 'IDR', 100000000, 50000000, 1000000, 900000, 100000);
  INSERT INTO &&skema_uji..T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, ID_PEGA, PL_NUMBER, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO, NAME_OF_INSURED, SEX, DOB, CURRENCY, SUM_INSURED, SUM_REASURED, GROSS_PREMIUM, NET_PREMIUM, COMM)
  VALUES ('UJI-PLD-C-2', 'UJI-PL-C', 'UJI-PL-C', 'UJI-PLN-C', 'UJI-POL-000C', 'UJI-PEMEGANG-1', 'UJI-SERT-C2', 'UJI-TERTANGGUNG-C2', 'F', DATE '1985-06-15', 'IDR', 200000000, 100000000, 2000000, 1800000, 200000);

  -- ===================================================================
  -- Claim Life - satu klaim per tahap (panduan bab 1 §1.2), atas polis
  -- UJI-POL-0001 (`BUSINESS_CODE` L1 dikenal - Save to RNM menuntutnya).
  --
  -- Bentuk kolom kerja meniru `services/pendaftaran.go`: CASE_ID = ID,
  -- CREATE_OP = akun pencipta. ⛔ `UJI-ADMIN` = akun stub bawaan
  -- (`VITE_STUB_PELAKU` kosong): kedua tab Admin menyaring CREATE_OP, jadi
  -- akun lain tidak melihat UJI-CLM-1 dan UJI-CLM-2.
  --
  -- TAHAP + PY_POSITION = pemegangnya (butir at, `models/tahap.go`).
  -- ===================================================================
  INSERT INTO &&skema_uji..T_WORK_CLAIM (ID, LINI, PY_POSITION, TYPE, CASE_ID, CREATE_OP, CREATE_OP_NAME, TGL_UPDATE, TAHAP, TGL_CREATE)
  VALUES ('UJI-CLM-1', 'LIFE', 'ReasLifeAdmin', 'QP', 'UJI-CLM-1', 'UJI-ADMIN', 'UJI-ADMIN', DATE '2026-09-29', 'Input Register', DATE '2026-09-29');
  INSERT INTO &&skema_uji..T_WORK_CLAIM (ID, LINI, PY_POSITION, TYPE, CASE_ID, CREATE_OP, CREATE_OP_NAME, TGL_UPDATE, TAHAP, TGL_CREATE)
  VALUES ('UJI-CLM-2', 'LIFE', 'ReasLifeAdmin', 'QP', 'UJI-CLM-2', 'UJI-ADMIN', 'UJI-ADMIN', DATE '2026-09-29', 'Outstanding Claim', DATE '2026-09-29');
  INSERT INTO &&skema_uji..T_WORK_CLAIM (ID, LINI, PY_POSITION, TYPE, CASE_ID, CREATE_OP, CREATE_OP_NAME, TGL_UPDATE, TAHAP, TGL_CREATE)
  VALUES ('UJI-CLM-3', 'LIFE', 'ReasLifeMedicalAdvisor', 'QP', 'UJI-CLM-3', 'UJI-ADMIN', 'UJI-ADMIN', DATE '2026-09-29', 'Medical Check', DATE '2026-09-29');
  INSERT INTO &&skema_uji..T_WORK_CLAIM (ID, LINI, PY_POSITION, TYPE, CASE_ID, CREATE_OP, CREATE_OP_NAME, TGL_UPDATE, TAHAP, TGL_CREATE)
  VALUES ('UJI-CLM-4', 'LIFE', 'ReasLifeSPV', 'QP', 'UJI-CLM-4', 'UJI-ADMIN', 'UJI-ADMIN', DATE '2026-09-29', 'Claim Analis', DATE '2026-09-29');

  -- Header - shared PK dengan T_WORK_CLAIM. STS_REJECT header = cermin baris
  -- terakhir (tiket 04); NULL bila barisnya belum berstatus.
  INSERT INTO &&skema_uji..T_GENERAL_CLAIM (ID, CLAIM_NO, STS_REJECT, BUSINESS_NAME, CURRENCY, CASEID_POLICY, POLICY_NO, BUSINESS_CODE)
  VALUES ('UJI-CLM-1', 'UJI-KL1.092026.00001', NULL, 'UJI-BISNIS-L1', 'IDR', 'UJI-PL-D', 'UJI-POL-0001', 'L1');
  INSERT INTO &&skema_uji..T_GENERAL_CLAIM (ID, CLAIM_NO, STS_REJECT, BUSINESS_NAME, CURRENCY, CASEID_POLICY, POLICY_NO, BUSINESS_CODE)
  VALUES ('UJI-CLM-2', 'UJI-KL1.092026.00002', NULL, 'UJI-BISNIS-L1', 'IDR', 'UJI-PL-D', 'UJI-POL-0001', 'L1');
  INSERT INTO &&skema_uji..T_GENERAL_CLAIM (ID, CLAIM_NO, STS_REJECT, BUSINESS_NAME, CURRENCY, CASEID_POLICY, POLICY_NO, BUSINESS_CODE)
  VALUES ('UJI-CLM-3', 'UJI-KL1.092026.00003', '0', 'UJI-BISNIS-L1', 'IDR', 'UJI-PL-D', 'UJI-POL-0001', 'L1');
  INSERT INTO &&skema_uji..T_GENERAL_CLAIM (ID, CLAIM_NO, STS_REJECT, BUSINESS_NAME, CURRENCY, CASEID_POLICY, POLICY_NO, BUSINESS_CODE)
  VALUES ('UJI-CLM-4', 'UJI-KL1.092026.00004', '0', 'UJI-BISNIS-L1', 'IDR', 'UJI-PL-D', 'UJI-POL-0001', 'L1');

  -- Peserta. Tanggal di dalam jendela valuasi (DOL 2026-03-15 di antara
  -- 2026-01-01 dan 2026-12-31) supaya Edit Date dan Save to RNM lolos
  -- gerbang tanggalnya. STS_REJECT peserta = cermin; IS_CHECK "true" =
  -- dipilih-untuk-diklaim (models.PenandaDipilih).
  INSERT INTO &&skema_uji..T_CLAIMLF_PREMIUMLIST_DETAIL (ID, CLAIM_ID, PL_NUMBER, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO, NAME_OF_INSURED, DOB, AGE, SEX, PLAN, IS_CHECK, STS_REJECT, DATE_OF_LOSS, RECEIVED_DATE, BEGIN_DATE, EFFECTIVE_DATE, EXPIRED_DATE, GROSS_VALUATION_BEGIN_DATE, GROSS_VALUATION_EXPIRED_DATE, RETRO_VALUATION_BEGIN_DATE, RETRO_VALUATION_EXPIRED_DATE, CURRENCY, SUM_INSURED, SUM_REASURED, CLAIM_AMOUNT, SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE, CEDING_RETENTION)
  VALUES ('UJI-PES-1-1', 'UJI-CLM-1', 'UJI-PLN-D', 'UJI-POL-0001', 'UJI-PEMEGANG-1', 'UJI-SERT-0011', 'UJI-TERTANGGUNG-11', DATE '1980-01-01', 46, 'M', 'UJI-PLAN-1', 'true', NULL, DATE '2026-03-15', DATE '2026-04-01', DATE '2025-01-01', DATE '2025-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', 'IDR', 500000000, 250000000, 25000000, 250000000, 0, 0, 250000000);
  INSERT INTO &&skema_uji..T_CLAIMLF_PREMIUMLIST_DETAIL (ID, CLAIM_ID, PL_NUMBER, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO, NAME_OF_INSURED, DOB, AGE, SEX, PLAN, IS_CHECK, STS_REJECT, DATE_OF_LOSS, RECEIVED_DATE, BEGIN_DATE, EFFECTIVE_DATE, EXPIRED_DATE, GROSS_VALUATION_BEGIN_DATE, GROSS_VALUATION_EXPIRED_DATE, RETRO_VALUATION_BEGIN_DATE, RETRO_VALUATION_EXPIRED_DATE, CURRENCY, SUM_INSURED, SUM_REASURED, CLAIM_AMOUNT, SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE, CEDING_RETENTION)
  VALUES ('UJI-PES-2-1', 'UJI-CLM-2', 'UJI-PLN-D', 'UJI-POL-0001', 'UJI-PEMEGANG-1', 'UJI-SERT-0021', 'UJI-TERTANGGUNG-21', DATE '1982-02-02', 44, 'F', 'UJI-PLAN-1', 'true', NULL, DATE '2026-03-15', DATE '2026-04-01', DATE '2025-01-01', DATE '2025-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', 'IDR', 500000000, 250000000, 50000000, 250000000, 0, 0, 250000000);
  INSERT INTO &&skema_uji..T_CLAIMLF_PREMIUMLIST_DETAIL (ID, CLAIM_ID, PL_NUMBER, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO, NAME_OF_INSURED, DOB, AGE, SEX, PLAN, IS_CHECK, STS_REJECT, DATE_OF_LOSS, RECEIVED_DATE, BEGIN_DATE, EFFECTIVE_DATE, EXPIRED_DATE, GROSS_VALUATION_BEGIN_DATE, GROSS_VALUATION_EXPIRED_DATE, RETRO_VALUATION_BEGIN_DATE, RETRO_VALUATION_EXPIRED_DATE, CURRENCY, SUM_INSURED, SUM_REASURED, CLAIM_AMOUNT, SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE, CEDING_RETENTION)
  VALUES ('UJI-PES-3-1', 'UJI-CLM-3', 'UJI-PLN-D', 'UJI-POL-0001', 'UJI-PEMEGANG-1', 'UJI-SERT-0031', 'UJI-TERTANGGUNG-31', DATE '1975-03-03', 51, 'M', 'UJI-PLAN-1', 'true', '0', DATE '2026-03-15', DATE '2026-04-01', DATE '2025-01-01', DATE '2025-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', 'IDR', 500000000, 250000000, 75000000, 250000000, 0, 0, 250000000);
  -- UJI-CLM-4 berpeserta DUA, satu keadaan per tombol Claim Analis:
  --   A - baris Outstanding berekening lengkap: Send Claim to Committee,
  --       atau Save Adjustment;
  --   C - baris ditolak ("2", penanda dipilih tercabut seperti sesudah
  --       Tolak): `Add` = putaran berikutnya.
  -- ⛔ GILIRAN-14 (butir bp meralat bo): peserta B TANPA baris dibuang. Baris
  -- pertama lahir saat Submit Register, jadi peserta terpilih tidak pernah
  -- tanpa baris; `Add` bukan lagi jalan lahirnya.
  INSERT INTO &&skema_uji..T_CLAIMLF_PREMIUMLIST_DETAIL (ID, CLAIM_ID, PL_NUMBER, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO, NAME_OF_INSURED, DOB, AGE, SEX, PLAN, IS_CHECK, STS_REJECT, DATE_OF_LOSS, RECEIVED_DATE, BEGIN_DATE, EFFECTIVE_DATE, EXPIRED_DATE, GROSS_VALUATION_BEGIN_DATE, GROSS_VALUATION_EXPIRED_DATE, RETRO_VALUATION_BEGIN_DATE, RETRO_VALUATION_EXPIRED_DATE, CURRENCY, SUM_INSURED, SUM_REASURED, CLAIM_AMOUNT, SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE, CEDING_RETENTION)
  VALUES ('UJI-PES-4-A', 'UJI-CLM-4', 'UJI-PLN-D', 'UJI-POL-0001', 'UJI-PEMEGANG-1', 'UJI-SERT-004A', 'UJI-TERTANGGUNG-4A', DATE '1978-04-04', 48, 'F', 'UJI-PLAN-1', 'true', '0', DATE '2026-03-15', DATE '2026-04-01', DATE '2025-01-01', DATE '2025-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', 'IDR', 500000000, 250000000, 150000000, 250000000, 0, 0, 250000000);
  INSERT INTO &&skema_uji..T_CLAIMLF_PREMIUMLIST_DETAIL (ID, CLAIM_ID, PL_NUMBER, POLICY_NO, POLICY_HOLDER, CERTIFICATE_NO, NAME_OF_INSURED, DOB, AGE, SEX, PLAN, IS_CHECK, STS_REJECT, DATE_OF_LOSS, RECEIVED_DATE, BEGIN_DATE, EFFECTIVE_DATE, EXPIRED_DATE, GROSS_VALUATION_BEGIN_DATE, GROSS_VALUATION_EXPIRED_DATE, RETRO_VALUATION_BEGIN_DATE, RETRO_VALUATION_EXPIRED_DATE, CURRENCY, SUM_INSURED, SUM_REASURED, CLAIM_AMOUNT, SHARE_NUSANTARA_RE, SHARE_RETRO, RETROCEDED_SHARE, CEDING_RETENTION)
  VALUES ('UJI-PES-4-C', 'UJI-CLM-4', 'UJI-PLN-D', 'UJI-POL-0001', 'UJI-PEMEGANG-1', 'UJI-SERT-004C', 'UJI-TERTANGGUNG-4C', DATE '1981-06-06', 45, 'F', 'UJI-PLAN-1', 'false', '2', DATE '2026-03-15', DATE '2026-04-01', DATE '2025-01-01', DATE '2025-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', DATE '2026-01-01', DATE '2026-12-31', 'IDR', 500000000, 250000000, 60000000, 250000000, 0, 0, 250000000);

  -- Baris adjustment - bentuk yang `SavePesertaClaim` 7.8 lahirkan saat
  -- Submit Register (butir bp): delapan medan dari pesertanya, tanpa status
  -- sampai Save to RNM. Klaim di sini disisipkan langsung, bukan lewat
  -- pendaftaran, jadi barisnya pun ditulis di sini. CURRENCY_ID DIBIARKAN KOSONG: nilainya
  -- pengenal tabel rujukan mata uang, dan mengarangnya lebih buruk daripada
  -- mengosongkannya (ADR-U-0027); `PeriksaSatuMataUang` melewati yang kosong.
  --
  -- UJI-CLM-1 dan UJI-CLM-2: berstatus NULL - bentuk yang pendaftaran Pega
  -- tinggalkan (SavePesertaClaim 7.8); Save to RNM yang menulis "0".
  INSERT INTO &&skema_uji..T_CLAIMLF_ADJUSTMENT (ID, PREMIUM_LIST_DETAIL_ID, SHARE_NUSANTARA_RE, CEDING_RETENTION, SUM_INSURED, SUM_REASURED, SHARE_RETRO, RETROCEDED_SHARE, CURRENCY, CLAIM_AMOUNT, STS_REJECT)
  VALUES ('UJI-ADJ-1-1', 'UJI-PES-1-1', 250000000, 250000000, 500000000, 250000000, 0, 0, 'IDR', 25000000, NULL);
  INSERT INTO &&skema_uji..T_CLAIMLF_ADJUSTMENT (ID, PREMIUM_LIST_DETAIL_ID, SHARE_NUSANTARA_RE, CEDING_RETENTION, SUM_INSURED, SUM_REASURED, SHARE_RETRO, RETROCEDED_SHARE, CURRENCY, CLAIM_AMOUNT, STS_REJECT)
  VALUES ('UJI-ADJ-2-1', 'UJI-PES-2-1', 250000000, 250000000, 500000000, 250000000, 0, 0, 'IDR', 50000000, NULL);
  INSERT INTO &&skema_uji..T_CLAIMLF_ADJUSTMENT (ID, PREMIUM_LIST_DETAIL_ID, SHARE_NUSANTARA_RE, CEDING_RETENTION, SUM_INSURED, SUM_REASURED, SHARE_RETRO, RETROCEDED_SHARE, CURRENCY, CLAIM_AMOUNT, STS_REJECT)
  VALUES ('UJI-ADJ-3-1', 'UJI-PES-3-1', 250000000, 250000000, 500000000, 250000000, 0, 0, 'IDR', 75000000, '0');
  -- ⭐ BARIS SIAP SERAH KOMITE: Outstanding, belum tertaut (KOMITE_ID NULL),
  -- rekening lengkap (`PeriksaRekening`: NAME_OF_BANK, ID_BANK, ACCOUNT_NO),
  -- 150.000.000 IDR - roster di bawah menutupnya dengan dua tingkat.
  INSERT INTO &&skema_uji..T_CLAIMLF_ADJUSTMENT (ID, PREMIUM_LIST_DETAIL_ID, SHARE_NUSANTARA_RE, CEDING_RETENTION, SUM_INSURED, SUM_REASURED, SHARE_RETRO, RETROCEDED_SHARE, CURRENCY, CLAIM_AMOUNT, STS_REJECT, NAME_OF_BANK, ID_BANK, ACCOUNT_NO, BRANCH_OF_BANK, SWIFT_CODE, PAYABLE_TO)
  VALUES ('UJI-ADJ-4-A1', 'UJI-PES-4-A', 250000000, 250000000, 500000000, 250000000, 0, 0, 'IDR', 150000000, '0', 'UJI-BANK', 'UJI-BANK-ID', 'UJI-REK-0001', 'UJI-CABANG-1', 'UJI-SWIFT-1', 'UJI-PENERIMA-1');
  INSERT INTO &&skema_uji..T_CLAIMLF_ADJUSTMENT (ID, PREMIUM_LIST_DETAIL_ID, SHARE_NUSANTARA_RE, CEDING_RETENTION, SUM_INSURED, SUM_REASURED, SHARE_RETRO, RETROCEDED_SHARE, CURRENCY, CLAIM_AMOUNT, STS_REJECT)
  VALUES ('UJI-ADJ-4-C1', 'UJI-PES-4-C', 250000000, 250000000, 500000000, 250000000, 0, 0, 'IDR', 60000000, '2');

  -- ===================================================================
  -- Roster Komite SINTETIS - `EMAILKOMITE` (panduan bab 3 §1.3).
  --
  -- Filter aplikasi (`repository/roster.go`, FilterEmailKomiteWithLimit):
  -- LIMIT_BOTTOM <= nilai klaim, STS_KLAIM = "LIFE", STS_AKTIF = "1", urut
  -- DEGREE. Untuk 150.000.000 yang terpilih UJI-KOMITE-1 dan UJI-KOMITE-2
  -- (dua tingkat); UJI-KOMITE-3 di atas pitanya, UJI-KOMITE-4 tidak aktif -
  -- keduanya bukti filternya bekerja, bukan hiasan.
  --
  -- OPERATOR_ID = akun yang masuk kotak masuk Komite
  -- (`KOMITE_OPERATORID = X-Pelaku`): uji sebagai anggota dengan
  -- VITE_STUB_PELAKU=UJI-KOMITE-1.
  --
  -- ⚠️ Baris UJI-EK-* yang sudah ada DIPAKAI ULANG, tidak disisipkan lagi:
  -- EMAILKOMITE bukan tabel migrasi, jadi ia bertahan melewati
  -- `-migrate-down` - tanpa ini muat ulang selalu gagal ORA-00001.
  -- ===================================================================
  SELECT COUNT(*) INTO v_roster FROM &&skema_uji..EMAILKOMITE WHERE ID LIKE 'UJI-EK-%';
  IF v_roster = 0 THEN
  INSERT INTO &&skema_uji..EMAILKOMITE (ID, OPERATOR_ID, JABATAN, EMAIL, DEGREE, LIMIT_BOTTOM, LIMIT_TOP, STS_AKTIF, STS_KLAIM)
  VALUES ('UJI-EK-1', 'UJI-KOMITE-1', 'UJI-JABATAN-1', 'uji-komite-1@contoh.invalid', 1, 0, 100000000, '1', 'LIFE');
  INSERT INTO &&skema_uji..EMAILKOMITE (ID, OPERATOR_ID, JABATAN, EMAIL, DEGREE, LIMIT_BOTTOM, LIMIT_TOP, STS_AKTIF, STS_KLAIM)
  VALUES ('UJI-EK-2', 'UJI-KOMITE-2', 'UJI-JABATAN-2', 'uji-komite-2@contoh.invalid', 2, 100000000, 1000000000, '1', 'LIFE');
  INSERT INTO &&skema_uji..EMAILKOMITE (ID, OPERATOR_ID, JABATAN, EMAIL, DEGREE, LIMIT_BOTTOM, LIMIT_TOP, STS_AKTIF, STS_KLAIM)
  VALUES ('UJI-EK-3', 'UJI-KOMITE-3', 'UJI-JABATAN-3', 'uji-komite-3@contoh.invalid', 3, 1000000000, NULL, '1', 'LIFE');
  INSERT INTO &&skema_uji..EMAILKOMITE (ID, OPERATOR_ID, JABATAN, EMAIL, DEGREE, LIMIT_BOTTOM, LIMIT_TOP, STS_AKTIF, STS_KLAIM)
  VALUES ('UJI-EK-4', 'UJI-KOMITE-4', 'UJI-JABATAN-4', 'uji-komite-4@contoh.invalid', 1, 0, 100000000, '0', 'LIFE');
  END IF;
END;
/

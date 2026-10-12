-- =====================================================================================================================
-- BACKFILL salinan QuotationData TAHAP 1 (varian POOLDATA, semua baris T_QUOTATIONDATA) T_QUOTATIONDATA untuk case NB FacIn yang SUDAH ADA (tiket 76b, keputusan work owner
-- 12-10-2026 "Ya, bertahap"; docs/issues/76b-quotationdata-lengkap.md).
--
-- ⚠️⚠️ PERINGATAN — NILAI SEKARANG, BUKAN NILAI SAAT CASE DIBUAT: skrip ini menyalin dari MASTER (MARKETINGOFFICER, T_M_ACCOUNT,
-- BUSINESS) dan rujukan T_NB_OPPORTUNITY sebagaimana isinya HARI INI. Bila marketing officer pindah team / cabang, nama tertanggung
-- di master berubah, atau baris BUSINESS diubah sesudah case dibuat, kolom salinan berisi nilai BARU itu — bukan yang dilihat
-- pengguna Pega saat case berjalan. Pega sendiri menyalin ulang dari master tiap form dibuka (CheckDataQuotation), jadi untuk case
-- yang masih berjalan hasilnya sama dengan perilaku Pega; untuk case yang sudah selesai, anggap salinan ini "nilai master per
-- tanggal backfill".
--
-- ⛔ DITULIS SAJA. BELUM PERNAH DIJALANKAN. Jalankan SESUDAH migrasi 253. Dijalankan DBA di DEV sesudah persetujuan work owner. Bukan
-- migrasi bernomor. TANPA COMMIT - COMMIT / ROLLBACK oleh DBA sesudah BAGIAN 2 cocok.
--
-- Skema langsung POOLDATA (varian siap-jalan, permintaan work owner 12-10-2026: "untuk yg ada di T_GENERAL_POLIS aja" - pola sama).
--
-- SUMBER = penulisan Save General biner baru (repository/kasus.go sqlSalinanQuotation):
--   - MARKETING_CODE / MARKETING_NAME / TEAM_GROUP / BRANCH_CODE / BRANCH_NAME = MARKETINGOFFICER per ID = MOID (CheckDataQuotation
--     langkah 6: CLIENTID / CLIENTNAME / TEAMGROUP / BRANCHDETAILID / BRANCHDETAILNAME);
--   - INSURED_ID = T_NB_OPPORTUNITY.INSURED_ID; INSURED_NAME = T_M_ACCOUNT.INSUREDNAME per ACCOUNT_ID;
--   - IS_MOP = SetMOP_Act (POLICY_TYPE 0 -> NonMOP, 1 / 2 -> MOP); BUSINESS_FAC 'F'; STATUS_BUSINESS '1';
--   - BUSINESS_CODE / BUSINESS_OLD_ID / GROUP_PANEL = BUSINESS ber-NOTE = Class of Business, BUSINESSGROUPID = group case,
--     GROUPPANEL IS NOT NULL, TEPAT SATU baris (selain itu NULL — sama dengan services kasusBisnisDari); BUSINESS_NAME = Class of
--     Business opportunity.
--   - BUSINESS_TYPE / BUSINESS_TYPE2 TIDAK diisi di sini: turunan BusinessType_DeT diport di Go (services businessTypeDari), tidak
--     ditiru SQL. Terisi pada Save General berikutnya.
-- Nilai lebih lebar dari kolom rancangan -> NULL (bukan ORA-12899).
--
-- CAKUPAN: SEMUA baris POOLDATA.T_QUOTATIONDATA (case NB FacIn sistem baru) yang STATUS_BUSINESS-nya masih NULL (= salinan belum
-- pernah ditulis; biner baru selalu menulis '1' di Save General, jadi hasil biner tidak ditimpa). Tanpa saringan T_WORK_POLIS / LINI.
-- =====================================================================================================================

-- ---- BAGIAN 0: hitung sebelum (simpan hasilnya) --------------------------------------------------------------------
SELECT 'Semua baris T_QUOTATIONDATA' AS HAL, COUNT(*) AS N FROM POOLDATA.T_QUOTATIONDATA q
UNION ALL SELECT 'AKAN DIISI (STATUS_BUSINESS NULL)', COUNT(*) FROM POOLDATA.T_QUOTATIONDATA q
 WHERE q.STATUS_BUSINESS IS NULL
UNION ALL SELECT 'MOID tanpa baris MARKETINGOFFICER', COUNT(*) FROM POOLDATA.T_QUOTATIONDATA q
 WHERE q.STATUS_BUSINESS IS NULL AND q.MOID IS NOT NULL
   AND NOT EXISTS (SELECT 1 FROM POOLDATA.MARKETINGOFFICER m WHERE m.ID = q.MOID);

-- ---- BAGIAN 1: isi (jumlah baris terubah = angka "AKAN DIISI" BAGIAN 0) -----------------------------------------------
UPDATE POOLDATA.T_QUOTATIONDATA q
   SET q.MARKETING_CODE = (SELECT MAX(CASE WHEN LENGTHB(m.CLIENTID) <= 50 THEN m.CLIENTID END) FROM POOLDATA.MARKETINGOFFICER m
                            WHERE m.ID = q.MOID),
       q.MARKETING_NAME = (SELECT MAX(CASE WHEN LENGTHB(m.CLIENTNAME) <= 500 THEN m.CLIENTNAME END) FROM POOLDATA.MARKETINGOFFICER m
                            WHERE m.ID = q.MOID),
       q.TEAM_GROUP = (SELECT MAX(CASE WHEN LENGTHB(m.TEAMGROUP) <= 50 THEN m.TEAMGROUP END) FROM POOLDATA.MARKETINGOFFICER m
                        WHERE m.ID = q.MOID),
       q.BRANCH_CODE = (SELECT MAX(CASE WHEN LENGTHB(m.BRANCHDETAILID) <= 50 THEN m.BRANCHDETAILID END)
                         FROM POOLDATA.MARKETINGOFFICER m WHERE m.ID = q.MOID),
       q.BRANCH_NAME = (SELECT MAX(CASE WHEN LENGTHB(m.BRANCHDETAILNAME) <= 500 THEN m.BRANCHDETAILNAME END)
                         FROM POOLDATA.MARKETINGOFFICER m WHERE m.ID = q.MOID),
       q.INSURED_ID = (SELECT MAX(CASE WHEN LENGTHB(o.INSURED_ID) <= 50 THEN o.INSURED_ID END) FROM POOLDATA.T_NB_OPPORTUNITY o
                        WHERE o.ID = q.PARENT_ID),
       q.INSURED_NAME = (SELECT MAX(CASE WHEN LENGTHB(a.INSUREDNAME) <= 500 THEN a.INSUREDNAME END) FROM POOLDATA.T_NB_OPPORTUNITY o
                          JOIN POOLDATA.T_M_ACCOUNT a ON a.ID = o.ACCOUNT_ID WHERE o.ID = q.PARENT_ID),
       q.IS_MOP = CASE WHEN q.POLICY_TYPE = '0' THEN 'NonMOP' WHEN q.POLICY_TYPE IN ('1', '2') THEN 'MOP' ELSE q.IS_MOP END,
       q.BUSINESS_FAC = 'F',
       q.STATUS_BUSINESS = '1',
       q.BUSINESS_NAME = (SELECT MAX(CASE WHEN LENGTHB(TRIM(o.CLASS_OF_BUSINESS)) <= 500 THEN TRIM(o.CLASS_OF_BUSINESS) END)
                           FROM POOLDATA.T_NB_OPPORTUNITY o WHERE o.ID = q.PARENT_ID),
       q.BUSINESS_CODE = (SELECT CASE WHEN COUNT(*) = 1 AND LENGTHB(MAX(TRIM(b.ID))) <= 50 THEN MAX(TRIM(b.ID)) END
                           FROM POOLDATA.BUSINESS b JOIN POOLDATA.T_NB_OPPORTUNITY o ON b.NOTE = TRIM(o.CLASS_OF_BUSINESS)
                            AND b.BUSINESSGROUPID = TRIM(o.GROUP_BUSINESS_ID) AND b.GROUPPANEL IS NOT NULL WHERE o.ID = q.PARENT_ID),
       q.BUSINESS_OLD_ID = (SELECT CASE WHEN COUNT(*) = 1 AND LENGTHB(MAX(TRIM(b.OLDID))) <= 50 THEN MAX(TRIM(b.OLDID)) END
                             FROM POOLDATA.BUSINESS b JOIN POOLDATA.T_NB_OPPORTUNITY o ON b.NOTE = TRIM(o.CLASS_OF_BUSINESS)
                              AND b.BUSINESSGROUPID = TRIM(o.GROUP_BUSINESS_ID) AND b.GROUPPANEL IS NOT NULL WHERE o.ID = q.PARENT_ID),
       q.GROUP_PANEL = (SELECT CASE WHEN COUNT(*) = 1 AND LENGTHB(MAX(TRIM(b.GROUPPANEL))) <= 50 THEN MAX(TRIM(b.GROUPPANEL)) END
                         FROM POOLDATA.BUSINESS b JOIN POOLDATA.T_NB_OPPORTUNITY o ON b.NOTE = TRIM(o.CLASS_OF_BUSINESS)
                          AND b.BUSINESSGROUPID = TRIM(o.GROUP_BUSINESS_ID) AND b.GROUPPANEL IS NOT NULL WHERE o.ID = q.PARENT_ID)
 WHERE q.STATUS_BUSINESS IS NULL;

-- ---- BAGIAN 2: hitung sesudah (sebelum COMMIT) -----------------------------------------------------------------------
-- Harapan: "AKAN DIISI" = 0; kolom kosong yang tersisa = sumbernya memang kosong / tak unik / terlalu lebar.
SELECT 'AKAN DIISI (harus 0)' AS HAL, COUNT(*) AS N FROM POOLDATA.T_QUOTATIONDATA q
 WHERE q.STATUS_BUSINESS IS NULL
UNION ALL SELECT 'MARKETING_CODE NULL', COUNT(*) FROM POOLDATA.T_QUOTATIONDATA q
 WHERE q.MARKETING_CODE IS NULL
UNION ALL SELECT 'INSURED_NAME NULL', COUNT(*) FROM POOLDATA.T_QUOTATIONDATA q
 WHERE q.INSURED_NAME IS NULL
UNION ALL SELECT 'BUSINESS_CODE NULL (BUSINESS tak unik / tak ada)', COUNT(*) FROM POOLDATA.T_QUOTATIONDATA q
 WHERE q.BUSINESS_CODE IS NULL;

SELECT q.PARENT_ID, q.MARKETING_CODE, q.TEAM_GROUP, q.BRANCH_CODE, q.INSURED_ID, q.IS_MOP, q.BUSINESS_CODE, q.BUSINESS_OLD_ID,
       q.GROUP_PANEL, q.STATUS_BUSINESS
  FROM POOLDATA.T_QUOTATIONDATA q WHERE q.PARENT_ID = 'NB-151029';

-- Cocok -> COMMIT; tidak -> ROLLBACK. (Sengaja tidak ditulis di skrip.)

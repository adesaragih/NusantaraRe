-- =====================================================================================================================
-- KONVERSI T_WORK_POLIS case NB FacIn ke POSITION_NOTE + POSITION tahap (tiket 77, keputusan work owner 12-10-2026; migrasi 254).
-- VERSI SIAP-JALAN: skema POOLDATA langsung (permintaan work owner). Pastikan T_WORK_POLIS di DEV memang di POOLDATA.
--
-- ⛔ DITULIS SAJA. BELUM PERNAH DIJALANKAN. Jalankan SESUDAH migrasi 254. TIDAK WAJIB: biner (repository/posisinote.go) membaca
-- baris lama lewat "posisi efektif", jadi kotak masuk / Submit benar sebelum maupun sesudah konversi. Dijalankan DBA sesudah
-- persetujuan work owner. TANPA COMMIT - COMMIT / ROLLBACK oleh DBA sesudah BAGIAN 2 cocok.
--
-- ATURAN (kontrak sesi 9d, butir 5):
--   - POSITION_NOTE = CASE POSITION WHEN 'Offer' THEN 'ReasFacInMarketing' WHEN 'Decline' THEN NULL ELSE POSITION END;
--   - POSITION = tahap dari FLAG_ONGOING_POLICY: '0' Offer, '2' Binding, lainnya Policy.
--   `[keputusan agent]` case Resolved-Completed (selesai Policy) IKUT dikonversi: NOTE = POSITION lama (workbasket terakhir), sama
--   dengan yang ditulis biner baru saat Selesai Policy (POSITION_NOTE tetap, POSITION 'Policy'); tanpa itu posisi efektifnya kosong.
--   Case Decline lama (POSITION 'Decline') -> NOTE NULL, POSITION = tahap.
-- CAKUPAN: hanya LINI 'FAC' ID 'NB-%', POSITION_NOTE masih NULL, dan BUKAN baris yang sudah berbentuk baru (Resolved dengan
-- POSITION tahap = Decline biner baru). Baris lini lain tidak disentuh. Aman diulang: baris yang sudah dikonversi tidak tercakup.
-- PEMERIKSA: sebaran POSISI EFEKTIF (ekspresi yang sama dengan biner) per nilai WAJIB SAMA di BAGIAN 0 dan BAGIAN 2.
-- =====================================================================================================================

-- ---- BAGIAN 0: sebelum (simpan hasilnya) -----------------------------------------------------------------------------
SELECT 'AKAN DIKONVERSI' AS HAL, COUNT(*) AS N FROM POOLDATA.T_WORK_POLIS w
 WHERE w.LINI = 'FAC' AND w.ID LIKE 'NB-%' AND w.POSITION_NOTE IS NULL
   AND NOT (w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy'));

SELECT w.POSITION, w.FLAG_ONGOING_POLICY, w.STATUS_WORK, COUNT(*) AS N FROM POOLDATA.T_WORK_POLIS w
 WHERE w.LINI = 'FAC' AND w.ID LIKE 'NB-%' AND w.POSITION_NOTE IS NULL
   AND NOT (w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy'))
 GROUP BY w.POSITION, w.FLAG_ONGOING_POLICY, w.STATUS_WORK ORDER BY 1, 2, 3;

-- Sebaran posisi efektif SELURUH case NB FacIn (bandingkan dengan BAGIAN 2 - harus identik).
SELECT CASE WHEN w.POSITION_NOTE IS NULL AND w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy') THEN CASE WHEN w.STATUS_WORK = 'Resolved-Decline' THEN 'Decline' END WHEN NVL(w.POSITION_NOTE, w.POSITION) = 'ReasFacInMarketing' AND w.FLAG_ONGOING_POLICY = '0' THEN 'Offer' ELSE NVL(w.POSITION_NOTE, w.POSITION) END AS POSISI_EFEKTIF, COUNT(*) AS N
  FROM POOLDATA.T_WORK_POLIS w WHERE w.LINI = 'FAC' AND w.ID LIKE 'NB-%'
 GROUP BY CASE WHEN w.POSITION_NOTE IS NULL AND w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy') THEN CASE WHEN w.STATUS_WORK = 'Resolved-Decline' THEN 'Decline' END WHEN NVL(w.POSITION_NOTE, w.POSITION) = 'ReasFacInMarketing' AND w.FLAG_ONGOING_POLICY = '0' THEN 'Offer' ELSE NVL(w.POSITION_NOTE, w.POSITION) END ORDER BY 1;

-- ---- BAGIAN 1: konversi (jumlah baris terubah = "AKAN DIKONVERSI") ---------------------------------------------------
UPDATE POOLDATA.T_WORK_POLIS w
   SET w.POSITION_NOTE = CASE w.POSITION WHEN 'Offer' THEN 'ReasFacInMarketing' WHEN 'Decline' THEN NULL ELSE w.POSITION END,
       w.POSITION = CASE w.FLAG_ONGOING_POLICY WHEN '0' THEN 'Offer' WHEN '2' THEN 'Binding' ELSE 'Policy' END
 WHERE w.LINI = 'FAC' AND w.ID LIKE 'NB-%' AND w.POSITION_NOTE IS NULL
   AND NOT (w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy'));

-- ---- BAGIAN 2: sesudah (sebelum COMMIT) --------------------------------------------------------------------------------
SELECT 'SISA (harus 0)' AS HAL, COUNT(*) AS N FROM POOLDATA.T_WORK_POLIS w
 WHERE w.LINI = 'FAC' AND w.ID LIKE 'NB-%' AND w.POSITION_NOTE IS NULL
   AND NOT (w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy'));

SELECT CASE WHEN w.POSITION_NOTE IS NULL AND w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy') THEN CASE WHEN w.STATUS_WORK = 'Resolved-Decline' THEN 'Decline' END WHEN NVL(w.POSITION_NOTE, w.POSITION) = 'ReasFacInMarketing' AND w.FLAG_ONGOING_POLICY = '0' THEN 'Offer' ELSE NVL(w.POSITION_NOTE, w.POSITION) END AS POSISI_EFEKTIF, COUNT(*) AS N
  FROM POOLDATA.T_WORK_POLIS w WHERE w.LINI = 'FAC' AND w.ID LIKE 'NB-%'
 GROUP BY CASE WHEN w.POSITION_NOTE IS NULL AND w.STATUS_WORK LIKE 'Resolved-%' AND w.POSITION IN ('Offer', 'Binding', 'Policy') THEN CASE WHEN w.STATUS_WORK = 'Resolved-Decline' THEN 'Decline' END WHEN NVL(w.POSITION_NOTE, w.POSITION) = 'ReasFacInMarketing' AND w.FLAG_ONGOING_POLICY = '0' THEN 'Offer' ELSE NVL(w.POSITION_NOTE, w.POSITION) END ORDER BY 1;

SELECT w.ID, w.POSITION, w.POSITION_NOTE, w.FLAG_ONGOING_POLICY, w.STATUS_WORK FROM POOLDATA.T_WORK_POLIS w WHERE w.ID IN ('NB-151029', 'NB-151039');

-- Cocok (SISA 0 dan sebaran posisi efektif identik) -> COMMIT; tidak -> ROLLBACK. (Sengaja tidak ditulis di skrip.)

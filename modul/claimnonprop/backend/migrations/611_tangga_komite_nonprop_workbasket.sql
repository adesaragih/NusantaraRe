-- 611 - roster komite Claim Non Prop ke WORKBASKET (prompt tahap 1 §3, keputusan work owner 09-10-2026: komite ikut
-- pola Claim Prop). Baris EMAILKOMITE STS_KLAIM NONPROP yang ada DIGANTI di tempat menurut DEGREE (UPDATE, bukan
-- nonaktif + tambah): DEGREE 1-4 -> ReasClaimDeptHead / ReasClaimTechDivHead / ReasClaimOpsDir / ReasClaimTechDir.
-- Hanya OPERATOR_ID, NAME, EMAIL yang berubah - JABATAN, DEGREE, dan batas nilai tetap. Nol baris baru, nol baris
-- dihapus. EMAIL kosong - email komite dikirim ke semua anggota workbasket. Keempat workbasket lahir di claimprop 537
-- (berjalan lebih dulu, rentang 520-559). Prasyarat persis: STS_KLAIM NONPROP, aktif, DEGREE 1-4. Idempoten.
-- Pola claimprop 537. NOL COMMIT. -migrate dijalankan work owner.
UPDATE {skema}.EMAILKOMITE e
   SET (OPERATOR_ID, NAME, EMAIL) =
       (SELECT t.WB, t.WB, NULL
          FROM (SELECT '1' DEG, 'ReasClaimDeptHead' WB FROM DUAL
                UNION ALL SELECT '2', 'ReasClaimTechDivHead' FROM DUAL
                UNION ALL SELECT '3', 'ReasClaimOpsDir' FROM DUAL
                UNION ALL SELECT '4', 'ReasClaimTechDir' FROM DUAL) t
         WHERE t.DEG = e.DEGREE)
 WHERE e.STS_KLAIM = 'NONPROP'
   AND e.STS_AKTIF = '1'
   AND e.DEGREE IN ('1', '2', '3', '4')
/

-- 640 - roster komite Claim Fac In ke WORKBASKET (keputusan work owner 10-10-2026 KCF-01, pola claimprop 537 /
-- claimnonprop 611). Baris EMAILKOMITE STS_KLAIM FACIN yang ada DIGANTI di tempat menurut DEGREE + JABATAN (UPDATE, bukan
-- nonaktif + tambah): 1 Claim Supervisor -> ReasClaimSPVA, 2 Claim Dept. Head -> ReasClaimDeptHead, 3 Technic Div. Head
-- -> ReasClaimTechDivHead, 4 Operational Director -> ReasClaimOpsDir, 5 Technical Director -> ReasClaimTechDir. Hanya
-- OPERATOR_ID, NAME, EMAIL yang berubah - JABATAN, DEGREE, dan batas nilai tetap. Nol baris baru, nol baris dihapus.
-- EMAIL kosong - email komite dikirim ke semua anggota workbasket (modul Komite Claim Fac In). Anggota ReasClaimSPVB
-- juga boleh memutus tingkat 1 (cadangan SPV A, KCF-01) - ditegakkan modul, bukan baris roster. Roster REOPEN dan lini
-- lain tidak tersentuh. Prasyarat persis: STS_KLAIM FACIN, aktif, pasangan DEGREE + JABATAN di atas. Idempoten.
--
-- ReasClaimSPVA / ReasClaimSPVB sudah ada di DEV (katalog 10-10-2026) tetapi tidak dilahirkan migrasi mana pun:
-- disisipkan bila belum ada (skema baru). Empat workbasket lain lahir di claimprop 537 (rentang 520-559, berjalan lebih
-- dulu).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_WORKBASKET (WORKBASKET_ID, NAME, IS_ACTIVE)
SELECT t.WB, t.NAMA, 1
  FROM (SELECT 'ReasClaimSPVA' WB, 'Klaim - Supervisor A' NAMA FROM DUAL
        UNION ALL SELECT 'ReasClaimSPVB', 'Klaim - Supervisor B' FROM DUAL) t
 WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_WORKBASKET w WHERE w.WORKBASKET_ID = t.WB)
/
UPDATE {skema}.EMAILKOMITE e
   SET (OPERATOR_ID, NAME, EMAIL) =
       (SELECT t.WB, t.WB, NULL
          FROM (SELECT '1' DEG, 'Claim Supervisor' JAB, 'ReasClaimSPVA' WB FROM DUAL
                UNION ALL SELECT '2', 'Claim Dept. Head', 'ReasClaimDeptHead' FROM DUAL
                UNION ALL SELECT '3', 'Technic Div. Head', 'ReasClaimTechDivHead' FROM DUAL
                UNION ALL SELECT '4', 'Operational Director', 'ReasClaimOpsDir' FROM DUAL
                UNION ALL SELECT '5', 'Technical Director', 'ReasClaimTechDir' FROM DUAL) t
         WHERE t.DEG = e.DEGREE AND t.JAB = e.JABATAN)
 WHERE e.STS_KLAIM = 'FACIN'
   AND e.STS_AKTIF = '1'
   AND (e.DEGREE, e.JABATAN) IN (('1', 'Claim Supervisor'), ('2', 'Claim Dept. Head'), ('3', 'Technic Div. Head'),
                                 ('4', 'Operational Director'), ('5', 'Technical Director'))
/

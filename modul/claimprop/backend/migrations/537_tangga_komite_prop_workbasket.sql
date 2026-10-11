-- 537 - tangga komite Claim Prop ke WORKBASKET (keputusan work owner 09-10-2026): OPERATOR_ID roster EMAILKOMITE
-- STS_KLAIM PROP (= `.KomiteID`, Pega `KomiteRouter` `AssignTo = .KomiteID`) DIGANTI nama workbasket, bukan akun
-- orang; penyetuju = anggota workbasket tingkat berjalan. ⚠️ Penyimpangan sadar dari Pega yang hidup (langkah lama ke
-- workbasket per tingkat `komitepnc`..`komitepnc4` ter-remark di KomiteRouter).
--
-- Baris PROP yang ada DIGANTI di tempat menurut JABATAN-nya (work owner: "menggantikan, bukan menambah"); hanya
-- OPERATOR_ID, NAME, EMAIL yang berubah - jabatan, DEGREE, dan batas nilai tetap. Nol baris baru, nol baris dihapus.
-- EMAIL kosong - email komite dikirim ke semua anggota workbasket (modul Komite Claim Prop). Tanpa SPV: SPV hanya di
-- FACIN (work owner 09-10-2026, "jangan sentuh"). FACIN / NONPROP / LIFE / REOPEN tidak tersentuh. Idempoten. Pola
-- workbasket Bordereaux 891.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_WORKBASKET (WORKBASKET_ID, NAME, IS_ACTIVE)
SELECT t.WB, t.NAMA, 1
  FROM (SELECT 'ReasClaimDeptHead' WB, 'Klaim - Department Head' NAMA FROM DUAL
        UNION ALL SELECT 'ReasClaimTechDivHead', 'Klaim - Technic Division Head' FROM DUAL
        UNION ALL SELECT 'ReasClaimOpsDir', 'Klaim - Operational Director' FROM DUAL
        UNION ALL SELECT 'ReasClaimTechDir', 'Klaim - Technical Director' FROM DUAL) t
 WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_WORKBASKET w WHERE w.WORKBASKET_ID = t.WB)
/
UPDATE {skema}.EMAILKOMITE e
   SET (OPERATOR_ID, NAME, EMAIL) =
       (SELECT t.WB, t.WB, NULL
          FROM (SELECT 'Claim Dept. Head' JABATAN, 'ReasClaimDeptHead' WB FROM DUAL
                UNION ALL SELECT 'Technic Div. Head', 'ReasClaimTechDivHead' FROM DUAL
                UNION ALL SELECT 'Operational Director', 'ReasClaimOpsDir' FROM DUAL
                UNION ALL SELECT 'Technical Director', 'ReasClaimTechDir' FROM DUAL) t
         WHERE t.JABATAN = e.JABATAN)
 WHERE e.STS_KLAIM = 'PROP'
   AND e.JABATAN IN ('Claim Dept. Head', 'Technic Div. Head', 'Operational Director', 'Technical Director')
/

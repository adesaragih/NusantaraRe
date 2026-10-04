-- 911 - baris menu modul Aggregate (`aggregate`), modul DI LUAR dua puluh folder korpus.
--
-- Perintah work owner 04-10-2026: modul Aggregate, label "Aggregate", kelompok MASTER TREATY ("Aggregate - Master
-- Treaty"). Padanan Pega: kelas ASM-FW-GISFW-Int-AGGREGATE (folder `Aggregate`: ShowAggregateList, GridDasbordAgg,
-- ChooseMasterID, UploadCSVAggregate_Act, SaveAggregate_Act, ...). PANDUAN-TIM-PER-MODUL bab 5: baris modul di luar
-- korpus dibuat langkah inti tersendiri. Bentuk DATAR (sesudah 901); kelompok MASTER TREATY ada sejak 909; URUTAN 4
-- sesudah Treaty In, Treaty In Adjustment, Treaty Contract Out. DIMIGRASI '0' - slot menu modulnya (996) yang
-- menyalakannya. Penjaga: `modulLuarKorpus` (inti/backend/penjaga/menu_test.go) dan `MODUL_LUAR_KORPUS`
-- (frontend/katalogKorpus.ts).
--
-- Akses: akun yang sudah ada TIDAK otomatis mendapat menu ini; admin mencentangnya di Kelola User.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, 'aggregate', 'Aggregate', 'MASTER TREATY', 'aggregate', 4, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'aggregate')
/

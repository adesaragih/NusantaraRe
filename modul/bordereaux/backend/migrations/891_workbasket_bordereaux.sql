-- 891 - tiga workbasket (peran) Bordereaux di M_WORKBASKET (keputusan work owner 04-10-2026):
--   ReasBordereauxAdmin      - pembuat berkas: hanya pemegangnya yang boleh Add / Input Data;
--   ReasBordereauxChecker    - Checker (tidak boleh pembuat berkas yang sama);
--   ReasBordereauxSupervisor - Supervisor.
-- Pengganti ID operator Checker/Supervisor yang di Pega ditulis langsung di section `InputBordereaux`. Hanya BARIS
-- DATA master, nol tabel baru; pemegangnya diatur lewat Kelola User. Berkas lama ber-POSITION 'Checker' /
-- 'Supervisor' langsung tampil bagi pemegang workbasket itu - tanpa pemindahan data.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.M_WORKBASKET (WORKBASKET_ID, NAME, IS_ACTIVE)
SELECT 'ReasBordereauxAdmin', 'Bordereaux Admin', 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_WORKBASKET WHERE WORKBASKET_ID = 'ReasBordereauxAdmin')
/
INSERT INTO {skema}.M_WORKBASKET (WORKBASKET_ID, NAME, IS_ACTIVE)
SELECT 'ReasBordereauxChecker', 'Bordereaux Checker', 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_WORKBASKET WHERE WORKBASKET_ID = 'ReasBordereauxChecker')
/
INSERT INTO {skema}.M_WORKBASKET (WORKBASKET_ID, NAME, IS_ACTIVE)
SELECT 'ReasBordereauxSupervisor', 'Bordereaux Supervisor', 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_WORKBASKET WHERE WORKBASKET_ID = 'ReasBordereauxSupervisor')
/

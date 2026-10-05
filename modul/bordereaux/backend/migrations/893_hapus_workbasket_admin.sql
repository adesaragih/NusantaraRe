-- 893 - workbasket ReasBordereauxAdmin dibuang (keputusan work owner 04-10-2026): Input Data, Edit, Delete, dan Submit
-- pembuat kini ditentukan hak menu Bordereaux PENUH (`M_LOGIN_GO_MENU.HAK`, migrasi inti 914). Checker dan Supervisor
-- tetap workbasket (891). DEV 04-10-2026: baris ada, nol pemegang.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DELETE FROM {skema}.M_LOGIN_GO_WORKBASKET WHERE WORKBASKET_ID = 'ReasBordereauxAdmin'
/
DELETE FROM {skema}.M_WORKBASKET WHERE WORKBASKET_ID = 'ReasBordereauxAdmin'
/

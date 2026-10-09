-- 534 - workbasket ReasKlaimAdmin dibuang (keputusan work owner 08-10-2026: "workbasket nii hapus dari master wb
-- ReasKlaimAdmin dan akun yang pake ini"). Tab Process Claim Prop bawaan = worklist pembuat tanpa cek workbasket
-- (XML `Flow_TreatyIn` Assignment2 `ToCurrentOperator`); switch Teknik memakai ReasKlaimTeknik. Tidak ada modul yang
-- memakai ReasKlaimAdmin. DEV 08-10-2026: baris master ada (aktif), 3 pemegang, nol foreign key ke M_WORKBASKET.
-- Pola migrasi Bordereaux 893 (ReasBordereauxAdmin).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DELETE FROM {skema}.M_LOGIN_GO_WORKBASKET WHERE WORKBASKET_ID = 'ReasKlaimAdmin'
/
DELETE FROM {skema}.M_WORKBASKET WHERE WORKBASKET_ID = 'ReasKlaimAdmin'
/

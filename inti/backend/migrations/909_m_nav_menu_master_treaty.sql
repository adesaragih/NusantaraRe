-- 909 - golongan menu baru MASTER TREATY (perintah work owner 04-10-2026: "group menunya tambahin MASTER TREATY;
-- Treaty In, Treaty In Adjustment, Treaty Contract Out masuk ke group itu").
--
-- Sebelumnya: CHECK GROUPMENU hanya TREATY, FACULTATIVE, KLAIM, MASTER (900); Treaty In (TREATY 3), Treaty In
-- Adjustment (TREATY 4), Treaty Contract Out (MASTER 3). Sesudahnya ketiganya MASTER TREATY 1-3, urutan sesuai
-- perintah work owner; urutan yang tersisa di TREATY dan MASTER dirapatkan. Urutan golongan di sidebar:
-- `menu.Golongan` (MASTER TREATY sesudah MASTER).
--
-- Bentuk yang dapat diulang: CHECK lama dibuang lewat blok berpelindung katalog (hanya bila masih ada), lalu dibuat
-- ulang dengan nilai baru; setiap UPDATE menyetel nilai tetap. Penjaga: `langkahGolonganMenu`
-- (inti/backend/penjaga/menu_test.go) dan `menuBersih` (inti/frontend/uji).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_NAV_MENU' AND CONSTRAINT_NAME = 'CK_M_NAV_MENU_GROUPMENU';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_NAV_MENU DROP CONSTRAINT CK_M_NAV_MENU_GROUPMENU';
  END IF;
END;
/
ALTER TABLE {skema}.M_NAV_MENU ADD CONSTRAINT CK_M_NAV_MENU_GROUPMENU CHECK (GROUPMENU IN ('TREATY', 'FACULTATIVE', 'KLAIM', 'MASTER', 'MASTER TREATY'))
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER TREATY', URUTAN = 1, TGL_UBAH = SYSDATE
WHERE KODE = 'treatyin'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER TREATY', URUTAN = 2, TGL_UBAH = SYSDATE
WHERE KODE = 'treatyinadjustment'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER TREATY', URUTAN = 3, TGL_UBAH = SYSDATE
WHERE KODE = 'treatycontractout'
/
-- Urutan yang tersisa dirapatkan (penjaga: URUTAN 1..n di setiap golongan): TREATY - PremiumList Life 5->3,
-- Endorsement Life 6->4; MASTER - Marketing Officer 4->3, Company Detail 5->4, Accounts 6->5.
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'TREATY', URUTAN = 3, TGL_UBAH = SYSDATE
WHERE KODE = 'premiumlistlife'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'TREATY', URUTAN = 4, TGL_UBAH = SYSDATE
WHERE KODE = 'endorsementlife'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 3, TGL_UBAH = SYSDATE
WHERE KODE = 'marketingofficer'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 4, TGL_UBAH = SYSDATE
WHERE KODE = 'companydetail'
/
UPDATE {skema}.M_NAV_MENU SET GROUPMENU = 'MASTER', URUTAN = 5, TGL_UBAH = SYSDATE
WHERE KODE = 'accounts'
/

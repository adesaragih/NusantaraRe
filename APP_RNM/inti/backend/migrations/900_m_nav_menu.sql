-- 900 - M_NAV_MENU: daftar menu aplikasi, dua tingkat, dengan GROUPMENU.
--
-- Permintaan work owner 30-09-2026: "untuk menu buatkan dari daftar table, nama
-- tabelnya M_NAV_MENU, masukkan list menunya ke dalam situ, ini berguna untuk
-- akses menu per akun nanti setelah dibuatkan akses login. Tambahkan kolom
-- GROUPMENU isinya TREATY, FACULTATIVE, KLAIM, MASTER."
-- (`PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`)
--
-- Rentang 900-949 milik `inti`: tabel lintas modul, bukan milik satu modul.
--
-- Dua tingkat dalam SATU tabel:
--   PARENT_ID NULL    = baris KELOMPOK modul (KODE = nama modul backend)
--   PARENT_ID terisi  = BUTIR menu di bawah kelompoknya (KODE = kunci halaman
--                       frontend `modul/daftar.ts`)
--
-- ⛔ Isi awal di bawah IDEMPOTEN (`WHERE NOT EXISTS` atas KODE): pelari migrasi
-- menjalankan tiap pernyataan tanpa transaksi, jadi langkah yang gagal separuh
-- jalan diulang dari awal - dan baris yang sudah masuk tidak boleh berganda.

CREATE TABLE {skema}.M_NAV_MENU (
  ID           NUMBER(10) NOT NULL,
  PARENT_ID    NUMBER(10),
  KODE         VARCHAR2(64) NOT NULL,
  LABEL        VARCHAR2(100) NOT NULL,
  GROUPMENU    VARCHAR2(16) NOT NULL,
  MODUL        VARCHAR2(32) NOT NULL,
  URUTAN       NUMBER(5) NOT NULL,
  STATUS_AKTIF VARCHAR2(1) DEFAULT '1' NOT NULL,
  DIMIGRASI    VARCHAR2(1) DEFAULT '0' NOT NULL,
  TGL_BUAT     DATE DEFAULT SYSDATE NOT NULL,
  TGL_UBAH     DATE,
  CONSTRAINT PK_M_NAV_MENU PRIMARY KEY (ID),
  CONSTRAINT UQ_M_NAV_MENU_KODE UNIQUE (KODE),
  CONSTRAINT FK_M_NAV_MENU_INDUK FOREIGN KEY (PARENT_ID) REFERENCES {skema}.M_NAV_MENU (ID),
  CONSTRAINT CK_M_NAV_MENU_GROUPMENU CHECK (GROUPMENU IN ('TREATY', 'FACULTATIVE', 'KLAIM', 'MASTER'))
)
/

CREATE INDEX {skema}.IX_M_NAV_MENU_PARENT ON {skema}.M_NAV_MENU (PARENT_ID)
/

CREATE INDEX {skema}.IX_M_NAV_MENU_GROUPMENU ON {skema}.M_NAV_MENU (GROUPMENU)
/

CREATE SEQUENCE {skema}.SEQ_M_NAV_MENU START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- ---------------------------------------------------------------------------
-- Kelompok modul - 20 baris, satu per folder modul korpus. LABEL = nama folder
-- korpus VERBATIM (`inti/labels.ts` MODUL); KODE = MODUL = nama modul backend
-- (tabel nama modul, `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`: nama folder tanpa
-- spasi, huruf kecil). GROUPMENU `[DIPUTUSKAN asisten; veto work owner]`.
-- DIMIGRASI '1' = modul sudah punya layar (Claim Life, PremiumList Life,
-- Komite Claim Life, Treaty Contract Out).
-- ---------------------------------------------------------------------------

-- KLAIM
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'claimfacin', 'Claim Fac In', 'KLAIM', 'claimfacin', 1, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'claimfacin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'claimlife', 'Claim Life', 'KLAIM', 'claimlife', 2, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'claimlife')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'claimnonprop', 'Claim Non Prop', 'KLAIM', 'claimnonprop', 3, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'claimnonprop')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'claimprop', 'Claim Prop', 'KLAIM', 'claimprop', 4, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'claimprop')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'komiteclaimfacin', 'Komite Claim FacIn', 'KLAIM', 'komiteclaimfacin', 5, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimfacin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'komiteclaimlife', 'Komite Claim Life', 'KLAIM', 'komiteclaimlife', 6, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimlife')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'komiteclaimnonprop', 'Komite Claim Non Prop', 'KLAIM', 'komiteclaimnonprop', 7, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimnonprop')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'komiteclaimprop', 'Komite Claim Prop', 'KLAIM', 'komiteclaimprop', 8, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'komiteclaimprop')
/

-- FACULTATIVE
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'nbfacin', 'NB FacIn', 'FACULTATIVE', 'nbfacin', 1, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'nbfacin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'rnwfacin', 'RNW Fac In', 'FACULTATIVE', 'rnwfacin', 2, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'rnwfacin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'endorsmentfacin', 'Endorsment Fac In', 'FACULTATIVE', 'endorsmentfacin', 3, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'endorsmentfacin')
/

-- TREATY
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'nbtreatyin', 'NB Treaty In', 'TREATY', 'nbtreatyin', 1, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'nbtreatyin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'edmtreatyin', 'EDM Treaty In', 'TREATY', 'edmtreatyin', 2, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'edmtreatyin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'treatyin', 'Treaty In', 'TREATY', 'treatyin', 3, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'treatyin')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'treatyinadjustment', 'Treaty In Adjustment', 'TREATY', 'treatyinadjustment', 4, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'treatyinadjustment')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'premiumlistlife', 'PremiumList Life', 'TREATY', 'premiumlistlife', 5, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'premiumlistlife')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'endorsementlife', 'Endorsement Life', 'TREATY', 'endorsementlife', 6, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'endorsementlife')
/

-- MASTER
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'mastercontractretrolife', 'Master Contract Retro Life', 'MASTER', 'mastercontractretrolife', 1, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'mastercontractretrolife')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'masterproductnamelife', 'Master Product Name Life', 'MASTER', 'masterproductnamelife', 2, '0' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'masterproductnamelife')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, NULL, 'treatycontractout', 'Treaty Contract Out', 'MASTER', 'treatycontractout', 3, '1' FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU WHERE KODE = 'treatycontractout')
/

-- ---------------------------------------------------------------------------
-- Butir menu - 5 baris. KODE = kunci halaman frontend (`modul/<nama>/menu.ts`);
-- LABEL VERBATIM dari label frontend yang sudah ada. GROUPMENU, MODUL, dan
-- DIMIGRASI DIWARISI dari kelompok induknya - dibaca dari baris induk, bukan
-- ditulis ulang, supaya keduanya tidak dapat berbeda. Beranda TIDAK di sini:
-- ia selalu tampil dan bukan milik modul.
-- ---------------------------------------------------------------------------

INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, k.ID, 'inbox', 'Inbox Claim Life', k.GROUPMENU, k.MODUL, 1, k.DIMIGRASI
FROM {skema}.M_NAV_MENU k
WHERE k.KODE = 'claimlife' AND k.PARENT_ID IS NULL
AND NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU b WHERE b.KODE = 'inbox')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, k.ID, 'register', 'Register', k.GROUPMENU, k.MODUL, 2, k.DIMIGRASI
FROM {skema}.M_NAV_MENU k
WHERE k.KODE = 'claimlife' AND k.PARENT_ID IS NULL
AND NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU b WHERE b.KODE = 'register')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, k.ID, 'premiumlist', 'PremiumList', k.GROUPMENU, k.MODUL, 1, k.DIMIGRASI
FROM {skema}.M_NAV_MENU k
WHERE k.KODE = 'premiumlistlife' AND k.PARENT_ID IS NULL
AND NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU b WHERE b.KODE = 'premiumlist')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, k.ID, 'komite', 'Inbox Komite', k.GROUPMENU, k.MODUL, 1, k.DIMIGRASI
FROM {skema}.M_NAV_MENU k
WHERE k.KODE = 'komiteclaimlife' AND k.PARENT_ID IS NULL
AND NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU b WHERE b.KODE = 'komite')
/
INSERT INTO {skema}.M_NAV_MENU (ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI)
SELECT {skema}.SEQ_M_NAV_MENU.NEXTVAL, k.ID, 'tco-tahun', 'Treaty Contract Out', k.GROUPMENU, k.MODUL, 1, k.DIMIGRASI
FROM {skema}.M_NAV_MENU k
WHERE k.KODE = 'treatycontractout' AND k.PARENT_ID IS NULL
AND NOT EXISTS (SELECT 1 FROM {skema}.M_NAV_MENU b WHERE b.KODE = 'tco-tahun')
/

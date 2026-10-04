-- M_LOGIN_GO_MENU - menu yang boleh dibuka setiap akun (Kelola User,
-- keputusan work owner 01-10-2026).
--
-- Permintaan work owner: "tambahkan menu untuk kelola user, CRUD user;
-- didalam itu nanti bisa tambahkan menu apa aja yang bisa diakses sama user
-- nya". Keputusan: akses PER AKUN (bukan per workbasket); menu yang tidak
-- diberikan DISEMBUNYIKAN di sidebar DAN DITOLAK server (403).
--
-- ⛔ SATU BARIS = SATU MENU SATU AKUN. `MENU_KODE` = KODE menu: nama modul
-- (kode baris modul tabel menu navigasi, migrasi 900) atau menu aplikasi
-- `kelolauser` (Kelola User - hidup di kode, `inti/backend/menu`). Pemegang
-- `kelolauser` adalah admin.
--
-- ⛔ TANPA FK ke tabel menu navigasi: menu aplikasi bukan barisnya, dan
-- aplikasi memeriksa saat menyimpan bahwa setiap KODE dikenal. FK ke
-- M_LOGIN_GO tanpa ON DELETE: "Hapus permanen" membuang baris anak lebih
-- dulu, dalam SATU transaksi bersama akunnya.
--
-- ⛔ ISI AWAL: setiap akun yang SUDAH ADA mendapat SEMUA menu - dua puluh
-- modul dan Kelola User (keputusan work owner: "semua menu"). Tanpanya akun
-- lama kehilangan seluruh layarnya begitu gerbang menyala, dan tidak ada
-- seorang pun yang dapat membuka Kelola User untuk memperbaikinya. Daftar
-- di bawah dijaga `inti/backend/penjaga` sama dengan menu modul hasil bersih
-- migrasi menu + menu aplikasi.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
CREATE TABLE {skema}.M_LOGIN_GO_MENU (
  LOGIN_ID   VARCHAR2(64) NOT NULL,
  MENU_KODE  VARCHAR2(64) NOT NULL,
  TGL_CREATE DATE DEFAULT SYSDATE NOT NULL,
  CONSTRAINT PK_M_LOGIN_GO_MENU PRIMARY KEY (LOGIN_ID, MENU_KODE),
  CONSTRAINT FK_M_LOGIN_GO_MENU_LOGIN FOREIGN KEY (LOGIN_ID) REFERENCES {skema}.M_LOGIN_GO (LOGIN_ID)
)
/
INSERT INTO {skema}.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE)
SELECT l.LOGIN_ID, k.KODE
  FROM {skema}.M_LOGIN_GO l
 CROSS JOIN (
  SELECT 'claimfacin' AS KODE FROM DUAL UNION ALL
  SELECT 'claimlife' FROM DUAL UNION ALL
  SELECT 'claimnonprop' FROM DUAL UNION ALL
  SELECT 'claimprop' FROM DUAL UNION ALL
  SELECT 'komiteclaimfacin' FROM DUAL UNION ALL
  SELECT 'komiteclaimlife' FROM DUAL UNION ALL
  SELECT 'komiteclaimnonprop' FROM DUAL UNION ALL
  SELECT 'komiteclaimprop' FROM DUAL UNION ALL
  SELECT 'nbfacin' FROM DUAL UNION ALL
  SELECT 'rnwfacin' FROM DUAL UNION ALL
  SELECT 'endorsmentfacin' FROM DUAL UNION ALL
  SELECT 'nbtreatyin' FROM DUAL UNION ALL
  SELECT 'edmtreatyin' FROM DUAL UNION ALL
  SELECT 'treatyin' FROM DUAL UNION ALL
  SELECT 'treatyinadjustment' FROM DUAL UNION ALL
  SELECT 'premiumlistlife' FROM DUAL UNION ALL
  SELECT 'endorsementlife' FROM DUAL UNION ALL
  SELECT 'mastercontractretrolife' FROM DUAL UNION ALL
  SELECT 'masterproductnamelife' FROM DUAL UNION ALL
  SELECT 'treatycontractout' FROM DUAL UNION ALL
  SELECT 'kelolauser' FROM DUAL
 ) k
/

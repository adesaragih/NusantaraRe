-- T_WORK_POLIS seragam dengan T_WORK_CLAIM - brief seragam kolom 01-10-2026.
--
-- `[DIPUTUSKAN; work owner 01-10-2026]` "samakan kolom antara T_WORK_POLIS
-- dan T_WORK_CLAIM: CREATE_OP, CREATE_OP_NAME, TGL_UPDATE, TGL_CREATE. Nama
-- sama antara T_WORK_POLIS dan T_WORK_CLAIM. COVER_KEY tetap ada, jangan
-- dihapus, di dua tabel itu."
--
-- ⛔ STATUS -> STATUS_WORK. Keduanya `pyWorkStatus` (klaim: migrasi 017;
-- polis: ralat 28-09-2026 di repository/polis_work.go), jadi namanya satu.
-- Panjangnya TIDAK diperkecil ke VARCHAR2(32) milik klaim: memperkecil dapat
-- gagal pada data yang lebih panjang.
--
-- ⛔ POSITION TIDAK diganti nama. Polis menyimpan `Position` (layar
-- `Offer`/`Premium`, Activity/ProtectAccept.xml b1207, b2288), klaim
-- `pyPosition` (peran). Menyamakannya dengan PY_POSITION menyatakan dua
-- properti Pega yang berbeda sebagai satu.
--
-- ⚠️ BUKAN `RENAME COLUMN`, walau hasil akhirnya sama. Pembanding STRUKTUR
-- dan penjaga tipe di `inti/` hanya mengenal CREATE TABLE, `ALTER ... ADD (`,
-- dan DROP COLUMN (`migrasi.KolomCreateTable`, `KolomAlterTambah`,
-- `KolomAlterBuang`); RENAME tidak terlihat olehnya, sehingga STATUS akan
-- tetap "ada" dan STATUS_WORK "tidak ada" bagi setiap penjaga. Bentuk yang
-- terlihat: tambah STATUS_WORK, salin isinya, buang STATUS - pola 901.
--
-- ⛔ COVER_KEY menunjuk T_WORK_POLIS(ID) induknya, TANPA `ON DELETE` - sama
-- persis dengan klaim butir d (001): aturan hapus bawaan Oracle MENOLAK
-- menghapus induk yang masih ditunjuk. XML PremiumList dan Endorsement tidak
-- memakai penunjuk induk (nol `pxCoverInsKey`), jadi kolomnya KOSONG sampai
-- ada modul yang terbukti mengisinya. Index-nya pola klaim
-- (IX_WORK_CLAIM_COVER_KEY): FK tanpa index memindai seluruh tabel tiap kali
-- induk dihapus.
--
-- ⚠️ PERNYATAAN 1 SATU-SATUNYA YANG TIDAK BERPELINDUNG KATALOG. Kolom baru
-- wajib lewat `ALTER ... ADD (` biasa - penjaga `pelanggaranBlokPLSQL`
-- menolak kolom baru di dalam blok, sebab blok tidak terlihat pembanding
-- STRUKTUR. Pernyataan 2-6 aman diulang. BILA 059 GAGAL SESUDAH PERNYATAAN 1,
-- percobaan ulang berhenti di ORA-01430 (kolom sudah ada) dan T_MIGRASI belum
-- mencatat 059. Pulihkan (DBA): jalankan seluruh pernyataan 059_down secara
-- manual - masing-masing berpelindung, aman pada keadaan separuh jadi - lalu
-- ulangi `-migrate`.
--
-- ⛔ PENGISIAN BARIS LAMA (pernyataan 6) - penghubungnya `T_PREMIUM_LIST.ID_PEGA
-- = T_WORK_POLIS.ID`: kasus baru menyisipkan header polisnya dengan
-- `ID = ID_PEGA = pengenal work` (repository/polis_kasus.go
-- sqlSisipPremiumListKosong), dan kotak masuk menggabung `p.ID_PEGA = w.ID`
-- (repository/polis_inbox.go). Hanya baris dengan TEPAT SATU pasangan.
--   - TGL_CREATE     <- TGL_INPUT
--   - CREATE_OP_NAME <- CREATE_OP_NAME, hanya bila muat VARCHAR2(128) - nilai
--                       yang lebih panjang (sumbernya VARCHAR2(255)) DIBIARKAN
--                       kosong, tidak dipotong.
--   - CREATE_OP      TETAP KOSONG: T_PREMIUM_LIST tidak menyimpan akun
--                    pembuat, dan mengarang nilainya dilarang (ADR-U-0027).
--   - TGL_UPDATE     TETAP KOSONG: T_PREMIUM_LIST tidak punya kolom itu.
-- Baris yang TGL_CREATE atau CREATE_OP_NAME-nya sudah terisi tidak disentuh.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
ALTER TABLE {skema}.T_WORK_POLIS ADD (
  STATUS_WORK    VARCHAR2(255),
  COVER_KEY      VARCHAR2(32),
  CREATE_OP      VARCHAR2(64),
  CREATE_OP_NAME VARCHAR2(128),
  TGL_CREATE     DATE,
  TGL_UPDATE     DATE
)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_POLIS' AND COLUMN_NAME = 'STATUS';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.T_WORK_POLIS SET STATUS_WORK = STATUS';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_POLIS' AND COLUMN_NAME = 'STATUS';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_POLIS DROP COLUMN STATUS';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_POLIS' AND CONSTRAINT_NAME = 'FK_WORK_POLIS_COVER_KEY';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_POLIS ADD CONSTRAINT FK_WORK_POLIS_COVER_KEY FOREIGN KEY (COVER_KEY) REFERENCES {skema}.T_WORK_POLIS (ID)';
  END IF;
END;
/
CREATE INDEX {skema}.IX_WORK_POLIS_COVER_KEY ON {skema}.T_WORK_POLIS (COVER_KEY)
/
UPDATE {skema}.T_WORK_POLIS w
   SET (TGL_CREATE, CREATE_OP_NAME) = (
         SELECT p.TGL_INPUT,
                CASE WHEN LENGTH(p.CREATE_OP_NAME) <= 128 THEN p.CREATE_OP_NAME END
           FROM {skema}.T_PREMIUM_LIST p
          WHERE p.ID_PEGA = w.ID)
 WHERE w.TGL_CREATE IS NULL
   AND w.CREATE_OP_NAME IS NULL
   AND (SELECT COUNT(*) FROM {skema}.T_PREMIUM_LIST p WHERE p.ID_PEGA = w.ID) = 1
/

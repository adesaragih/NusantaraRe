-- T_WORK_CLAIM diseragamkan - keputusan work owner 01-10-2026.
--
-- Lanjutan brief seragam kolom T_WORK_POLIS / T_WORK_CLAIM (migrasi 059
-- PremiumList). Jawaban work owner atas `select * from t_work_claim`:
--   - PY_POSITION diganti nama menjadi POSITION;
--   - TYPE dipindah ke header klaim T_GENERAL_CLAIM;
--   - CASE_ID dihapus, ID dipakai sebagai gantinya;
--   - TAHAP TETAP (tangga kerja klaim).
--
-- ⚠️ POSITION klaim tetap berisi `pyPosition` - NAMA PERAN (`ReasLifeAdmin`,
-- `ReasLifeMedicalAdvisor`, `ReasLifeSPV`). POSITION polis berisi `Position`
-- - layar `Offer`/`Premium`. Namanya kini sama, ARTINYA BERBEDA: jangan
-- membandingkan kolom POSITION kedua tabel kerja (keputusan work owner,
-- membalik brief seragam §1).
--
-- ⛔ TYPE pindah ke T_GENERAL_CLAIM: `kontrak.KlaimKomite.TypeKlaim` dan
-- enam alur yang membacanya (akseptasi, DOL, Komite, simpan RNM) membaca
-- header klaim. Baris kasus Komite (`KMTLF-`) tidak punya header klaim; TYPE
-- mereka salinan TYPE klaim induk (`pxAddChildWork`) yang tidak pernah
-- dibaca, jadi tidak ada nilai yang hilang.
--
-- ⛔ CASE_ID: kasus baru selalu `CASE_ID = ID` (services/pendaftaran.go).
-- Pernyataan 1 - PENGAMAN - memasang CHECK (CASE_ID IS NULL OR CASE_ID = ID):
-- bila satu baris saja berbeda, Oracle menolak (ORA-02293) dan migrasi
-- berhenti SEBELUM apa pun berubah. Perbaiki datanya dulu (DBA, persetujuan
-- work owner), baru ulangi `-migrate`. CHECK itu tak bernama, jadi aman
-- diulang, dan ikut terbuang bersama kolomnya (pernyataan 8, CASCADE
-- CONSTRAINTS). Migrasi klaim lama (tiket 13) memakai CASEID warisan
-- sebagai ID (repository/migrasidata.go BongkarBarisLama).
--
-- ⚠️ BUKAN `RENAME COLUMN` untuk PY_POSITION: pembanding STRUKTUR dan penjaga
-- tipe `inti/` hanya mengenal CREATE TABLE, `ALTER ... ADD (`, dan DROP
-- COLUMN - pola 059 (tambah, salin, buang).
--
-- ⚠️ PERNYATAAN 2 DAN 3 SATU-SATUNYA YANG TIDAK BERPELINDUNG KATALOG: kolom
-- baru wajib lewat `ALTER ... ADD (` biasa (penjaga `pelanggaranBlokPLSQL`).
-- BILA 023 GAGAL SESUDAH PERNYATAAN 2, percobaan ulang berhenti di ORA-01430
-- dan T_MIGRASI belum mencatat 023. Pulihkan (DBA): jalankan seluruh
-- pernyataan 023_down secara manual - masing-masing berpelindung, aman pada
-- keadaan separuh jadi - lalu ulangi `-migrate`.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'CASE_ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_CLAIM ADD CHECK (CASE_ID IS NULL OR CASE_ID = ID)';
  END IF;
END;
/
ALTER TABLE {skema}.T_WORK_CLAIM ADD (
  POSITION VARCHAR2(64)
)
/
ALTER TABLE {skema}.T_GENERAL_CLAIM ADD (
  TYPE VARCHAR2(32)
)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'PY_POSITION';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.T_WORK_CLAIM SET POSITION = PY_POSITION';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'PY_POSITION';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_CLAIM DROP COLUMN PY_POSITION';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'TYPE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.T_GENERAL_CLAIM g SET g.TYPE = (SELECT w.TYPE FROM {skema}.T_WORK_CLAIM w WHERE w.ID = g.ID) WHERE g.TYPE IS NULL';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'TYPE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_CLAIM DROP COLUMN TYPE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_WORK_CLAIM' AND COLUMN_NAME = 'CASE_ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_WORK_CLAIM DROP COLUMN CASE_ID CASCADE CONSTRAINTS';
  END IF;
END;
/

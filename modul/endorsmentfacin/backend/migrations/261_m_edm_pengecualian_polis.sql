-- 261 - M_EDM_PENGECUALIAN_POLIS: daftar nomor polis yang dikecualikan pemeriksaan alur masuk Endorsment Fac In
-- (keputusan work owner 11-10-2026: tabel konfigurasi baru, isinya diisi DBA, ADR-0035; kosong = tanpa bypass).
--
-- Pengganti literal nomor polis yang tertanam di prakondisi Pega `[terverifikasi]` (nilainya tidak disalin ke repo,
-- CLAUDE.md §4.10):
--   JENIS = 'TANGGAL' - `Endorsment Fac In\Activity\CheckEDMPolisDate.xml` langkah 2 (1 literal: lompat ke END,
--                       validasi tanggal dilewati);
--   JENIS = 'BATAL'   - `Endorsment Fac In\Activity\SetValueToEDMWork.xml` langkah 5 (2 literal: kasus tetap
--                       dibuat walau EdmType polis lama = 1).
-- Cacah dua cara (urai XML / cacah baris atas `pyStepsPreCondParamsWhen`): 1 dan 2, irisan 0.
--
-- PK (NOPOLIS, JENIS): satu polis boleh tercatat di kedua daftar (keputusan agent; kontrak mengusulkan PK NOPOLIS
-- saja). CHECK JENIS menolak salah ketik DBA - salah ketik tanpa CHECK = pengecualian yang diam-diam tidak berlaku.
-- NOPOLIS VARCHAR2(100) = lebar JSON_POLIS.NOPOLIS / FACINPRODUCTION.NOPOLIS (`D:\migrasi\RNM\DDL`).
-- Nomor 260 DICADANGKAN untuk T_WORK_POLIS/sequence kasus EDM- (I3). Ditulis, TIDAK dijalankan agent.
-- Nol COMMIT (ADR-U-0029).
CREATE TABLE {skema}.M_EDM_PENGECUALIAN_POLIS (
  NOPOLIS    VARCHAR2(100) NOT NULL,
  JENIS      VARCHAR2(10) NOT NULL,
  KETERANGAN VARCHAR2(1000),
  TGL_BUAT   DATE DEFAULT SYSDATE,
  CONSTRAINT PK_M_EDM_PENGECUALIAN_POLIS PRIMARY KEY (NOPOLIS, JENIS),
  CONSTRAINT CK_M_EDM_PENGECUALIAN_JENIS CHECK (JENIS IN ('TANGGAL', 'BATAL'))
)
/

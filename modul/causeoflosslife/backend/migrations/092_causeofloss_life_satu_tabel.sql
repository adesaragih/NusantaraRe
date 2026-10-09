-- 092 - Cause Of Loss Life SATU tabel CAUSEOFLOSS_LIFE (langkah 3 dari 3).
--
-- Keputusan work owner 08-10-2026 K1/K2. Tiga blok berpelindung katalog (semuanya hanya selama JSONDATA masih ada):
--   1. ISI: CAUSEOFLOSS dari kunci JSON "CauseofLoss" APA ADANYA - ekspresi notasi titik PERSIS teks view lama
--      (`a.JSONDATA.CauseofLoss`, 090_down), jadi nilainya identik dengan yang view keluarkan. Notasi titik JSON peka
--      huruf: `CauseofLoss` = nama kunci Pega. Nilai yang tidak muat VARCHAR2(200) menghentikan langkah ORA-12899
--      (tidak dipotong). Baris 100001 (`"CauseofLoss":""`) terbaca NULL oleh view DAN oleh ekspresi ini - kolomnya NULL,
--      sama dengan view; baris itu tidak diisi dan tidak dihapus (K1). Huruf TIDAK diubah (XML tanpa pengubah huruf).
--      Bukti baca-saja sebelum -migrate: modul/causeoflosslife/docs/sql/col_bukti_k1.sql.
--   2. PEMERIKSAAN K1.4 - BERHENTI sebelum JSONDATA dibuang HANYA bila nilai kolom BERBEDA dari yang view baca: kolom
--      NULL padahal view tidak, kolom terisi padahal view NULL, atau keduanya terisi tetapi tidak sama. NULL = NULL
--      dianggap SAMA (100001 lolos). Bentuknya UPDATE yang SENGAJA gagal (pola 944): baris seperti itu diberi ID NULL
--      -> ORA-01407 (ID NOT NULL, PK SYS_C008825) -> pelari berhenti, 092 tidak tercatat, JSONDATA utuh. Nol baris =
--      nol perubahan. Tanpa tanda kutip (teks EXECUTE IMMEDIATE blok berpelindung tidak boleh memuatnya -
--      polaPerintahKatalog). Pemaksa menulis ID, BUKAN CAUSEOFLOSS: nilai yang diperiksa tidak pernah diubah
--      pemeriksaannya sendiri (bentuk RPAD 948 tidak dipakai - ID di sini sudah NOT NULL sejak Pega);
--   3. buang JSONDATA + constraint IS JSON (ENSURE_M_CAUSEOFLOSS_LIFE_JSON) lewat CASCADE CONSTRAINTS.
-- Kunci pxObjClass dan pyRuleHarness (bukan kolom view) TIDAK dipindah - hanya di cadangan P3. Prosedur
-- PEGA_M_CAUSEOFLOSS_LIFE menjadi INVALID (K3, diterima WO); sequence M_CAUSEOFLOSS_LIFE_SEQ TETAP.
--
-- SETIAP pernyataan aman DIULANG (isi menulis nilai yang sama; pemeriksaan tidak mengubah apa pun; ketiganya melewati
-- diri sesudah JSONDATA dibuang). Volume DEV 4 baris. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.CAUSEOFLOSS_LIFE m SET m.CAUSEOFLOSS = m.JSONDATA.CauseofLoss';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.CAUSEOFLOSS_LIFE m SET m.ID = NULL WHERE (m.CAUSEOFLOSS IS NULL AND m.JSONDATA.CauseofLoss IS NOT NULL) OR (m.CAUSEOFLOSS IS NOT NULL AND m.JSONDATA.CauseofLoss IS NULL) OR m.CAUSEOFLOSS <> m.JSONDATA.CauseofLoss';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.CAUSEOFLOSS_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/

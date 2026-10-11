-- 086 - Cover Life SATU tabel M_COVER_LIFE (langkah 2 dari 2) - nama TETAP, TANPA RENAME (keputusan work owner
-- 08-10-2026 C1).
--
-- Empat blok berpelindung katalog, BERURUTAN, SETIAP pernyataan aman DIULANG:
--   1. ISI COVER dan NOTE dari JSONDATA APA ADANYA, SELAGI view COVER_LIFE dan JSONDATA masih ada - ekspresi notasi titik
--      PERSIS teks view lama (`a.JSONDATA.Cover`, `a.JSONDATA.Note`; 086_down), peka huruf: `Cover` / `Note` = nama kunci
--      Pega. Jadi nilainya identik dengan yang view keluarkan. Nilai yang tidak muat lebar 085 menghentikan langkah
--      ORA-12899 (tidak dipotong). NOTE NULL di 4/4 baris DEV (kunci Note tidak ada) - sama dengan view. Huruf TIDAK
--      diubah (XML tanpa pengubah huruf). Bukti baca-saja sebelum -migrate: modul/coverlife/docs/sql/cov_bukti.sql;
--   2. PEMERIKSAAN C1.3 - BERHENTI sebelum JSONDATA dibuang HANYA bila kolom BERBEDA dari yang view baca (COVER atau
--      NOTE: kolom NULL padahal view tidak, kolom terisi padahal view NULL, atau keduanya terisi tetapi tidak sama).
--      NULL = NULL dianggap SAMA (NOTE 4/4). Bentuknya UPDATE yang SENGAJA gagal (pola 944 / 092): baris seperti itu
--      diberi ID NULL -> ORA-01407 (ID NOT NULL, PK SYS_C009203) -> pelari berhenti, 086 tidak tercatat, JSONDATA dan
--      view utuh. Nol baris = nol perubahan. Tanpa tanda kutip (polaPerintahKatalog). Pemaksa menulis ID, BUKAN
--      COVER / NOTE: nilai yang diperiksa tidak pernah diubah pemeriksaannya sendiri;
--   3. buang JSONDATA + constraint IS JSON (ENSURE_M_COVER_LIFE_JSON) lewat CASCADE CONSTRAINTS;
--   4. DROP VIEW COVER_LIFE - TERAKHIR, sesudah kolom terbukti sama; blok SYS.ALL_VIEWS (hanya bila objeknya VIEW).
-- Hasil: SATU tabel M_COVER_LIFE (ID VARCHAR2(10) PK SYS_C009203, COVER, NOTE), urutan kolom = view lama, tanpa
-- JSONDATA, dan TIDAK ADA lagi objek bernama COVER_LIFE.
-- Kunci pxObjClass dan pyRuleHarness (bukan kolom view) TIDAK dipindah - hanya di cadangan P3. Prosedur
-- PEGA_M_COVER_LIFE menjadi INVALID (C2, diterima WO); sequence M_COVER_LIFE_SEQ TETAP (ID baru = rumus Pega). Objek
-- COVERAGE* / COVERNOTE* milik aplikasi lain TIDAK disentuh.
-- Volume DEV 4 baris. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_COVER_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.M_COVER_LIFE m SET m.COVER = m.JSONDATA.Cover, m.NOTE = m.JSONDATA.Note';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_COVER_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.M_COVER_LIFE m SET m.ID = NULL WHERE (m.COVER IS NULL AND m.JSONDATA.Cover IS NOT NULL) OR (m.COVER IS NOT NULL AND m.JSONDATA.Cover IS NULL) OR m.COVER <> m.JSONDATA.Cover OR (m.NOTE IS NULL AND m.JSONDATA.Note IS NOT NULL) OR (m.NOTE IS NOT NULL AND m.JSONDATA.Note IS NULL) OR m.NOTE <> m.JSONDATA.Note';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_COVER_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_COVER_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS
   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'COVER_LIFE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP VIEW {skema}.COVER_LIFE';
  END IF;
END;
/

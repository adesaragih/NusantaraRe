-- 944 - Benefit SATU tabel BENEFIT_LIFE (langkah 3 dari 3).
--
-- Keputusan work owner 08-10-2026 K1/K2. Tiga blok berpelindung katalog (semuanya hanya selama JSONDATA masih ada):
--   1. ISI: BENEFIT dari kunci JSON "Benefit" APA ADANYA - ekspresi notasi titik PERSIS teks view lama
--      (`a.JSONDATA.Benefit`, 942_down), jadi nilainya identik dengan yang view keluarkan. Notasi titik JSON peka
--      huruf: `Benefit` (B besar) = nama kunci Pega. Nilai yang tidak muat VARCHAR2(200) menghentikan langkah ORA-12899
--      (tidak dipotong). Bukti baca-saja sebelum -migrate: modul/benefitlife/docs/sql/benefit_bukti_k2.sql.
--   2. PEMERIKSAAN K2 - BERHENTI sebelum JSONDATA dibuang bila ada baris yang JSON-nya memuat Benefit tetapi BENEFIT-nya
--      NULL (atau berbeda). Bentuknya UPDATE yang SENGAJA gagal: baris seperti itu diberi ID NULL -> ORA-01407 (ID
--      NOT NULL, PK SYS_C009031) -> pelari berhenti, langkah 944 tidak tercatat, JSONDATA utuh. Nol baris = nol
--      perubahan. Ditulis tanpa tanda kutip karena teks EXECUTE IMMEDIATE blok berpelindung tidak boleh memuatnya
--      (polaPerintahKatalog). Batasnya: "punya Benefit" dibaca lewat notasi titik yang sama (nilai JSON null / objek /
--      lebih dari 4000 byte terbaca NULL) - bukti K2 menghitung JSON_EXISTS terpisah;
--   3. buang JSONDATA + constraint IS JSON (ENSURE_M_BENEFIT_LIFE_JSON) lewat CASCADE CONSTRAINTS.
-- Kunci Number (5 baris, selalu = ID), pxObjClass, pyRuleHarness, dan kunci px* lain (bukan kolom view) TIDAK
-- dipindah - hanya di cadangan P3. Prosedur PEGA_M_BENEFIT_LIFE menjadi INVALID (K3, diterima WO); sequence
-- M_BENEFIT_LIFE_SEQ TETAP.
--
-- SETIAP pernyataan aman DIULANG (isi menulis nilai yang sama; pemeriksaan tidak mengubah apa pun; ketiganya melewati
-- diri sesudah JSONDATA dibuang). Volume DEV 10 baris. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.BENEFIT_LIFE m SET m.BENEFIT = m.JSONDATA.Benefit';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.BENEFIT_LIFE m SET m.ID = NULL WHERE m.JSONDATA.Benefit IS NOT NULL AND (m.BENEFIT IS NULL OR m.BENEFIT <> m.JSONDATA.Benefit)';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.BENEFIT_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/

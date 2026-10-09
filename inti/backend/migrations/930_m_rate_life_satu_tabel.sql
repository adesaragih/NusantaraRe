-- 930 - rincian R/I Rate Life menjadi tabel flat M_RATE_LIFE (langkah 2 dari 2; langkah 1 = 929).
--
-- Keputusan work owner 07-10-2026 (lihat 929). Isi kolom dipindah dari JSONDATA APA ADANYA (notasi titik = definisi view
-- RATE_LIFE, tanpa konversi; nilai yang tidak muat lebar 929 menghentikan langkah dengan ORA-12899, tidak dipotong),
-- lalu JSONDATA + constraint IS JSON (ENSURE_M_RATE_LIFE_JSON) dibuang, indeks IDUSEDBY dibuat, dan view RATE_LIFE
-- dibuang TERAKHIR. FLAG (ada di JSON, tidak di view) TIDAK dipindah (RALAT R5). Baris yatim (IDUSEDBY tanpa ringkasan)
-- ikut dipindah apa adanya. Prosedur PEGA_M_RATE_LIFE menjadi INVALID (diterima WO).
-- ⛔ PRASYARAT: cadangan CSV ID + JSONDATA (LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md langkah (a)).
--
-- SETIAP pernyataan aman DIULANG sesudah gagal di tengah (pelari tanpa transaksi): pengisian dan pembuangan JSONDATA
-- lewat blok berpelindung katalog (jalan hanya selama JSONDATA masih ada), CREATE INDEX (ORA-00955 ditoleransi), dan
-- DROP VIEW terakhir (bila pelari mati SESUDAH DROP VIEW tetapi sebelum mencatat 930, pemulihannya mencatat 930 manual).
-- Volume DEV ±98 ribu baris: satu UPDATE, undo puluhan MB, perkiraan detik-hingga-menit.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RATE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.M_RATE_LIFE m SET m.IDUSEDBY = m.JSONDATA.IDUSEDBY, m.USEDBY = m.JSONDATA.USEDBY, m.TYPE = m.JSONDATA.TYPE, m.GENDER = m.JSONDATA.GENDER, m.CONTRACT = m.JSONDATA.CONTRACT, m.AGE = m.JSONDATA.AGE, m.RATE = m.JSONDATA.RATE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RATE_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RATE_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
-- Indeks rincian per ringkasan: grid Rate Detail, kembar, Delete berantai, salinan nama; dan pembaca lain (Claim Life,
-- PremiumList Life, Product Name Life, Contract Retro Life) yang memilih baris lewat IDUSEDBY.
CREATE INDEX {skema}.IX_M_RATE_LIFE_IDUSEDBY ON {skema}.M_RATE_LIFE (IDUSEDBY)
/
-- TERAKHIR: view warisan RATE_LIFE dibuang (definisinya dipulihkan 930_down).
DROP VIEW {skema}.RATE_LIFE
/

-- 937 - ringkasan R/I Risk SATU tabel RIRISK_LIFE_SUMMARY (langkah 3 dari 3).
--
-- Keputusan work owner 08-10-2026 K1/K2. Isi kolom dari JSONDATA APA ADANYA (notasi titik = definisi view lama; nilai
-- yang tidak muat lebar 936 menghentikan langkah ORA-12899, tidak dipotong), lalu JSONDATA + constraint IS JSON
-- (ENSURE_M_RIRISK_LIFE_SUMMARY_JSON, 33 byte - tidak ditulis di sini) dibuang lewat CASCADE CONSTRAINTS, lalu indeks
-- nama (pemeriksa nama kembar Upload, `SqlPemakaiNama`). Kunci pxObjClass (bukan kolom view) TIDAK dipindah - hanya di
-- cadangan CSV P3. Prosedur PEGA_M_RIRISK_LIFE_SUMMARY menjadi INVALID (K3, diterima WO).
--
-- SETIAP pernyataan aman DIULANG: blok berpelindung katalog (isi SEBELUM buang) dan CREATE INDEX (ORA-00955
-- ditoleransi). Volume DEV 117 baris. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE_SUMMARY' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.RIRISK_LIFE_SUMMARY m SET m.USEDBY = m.JSONDATA.USEDBY, m.MODIFIEDDATE = m.JSONDATA.MODIFIEDDATE, m.OPERATORID = m.JSONDATA.OPERATORID';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE_SUMMARY' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.RIRISK_LIFE_SUMMARY DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
-- Indeks nama: pemeriksa nama kembar ririsklife (SqlPemakaiNama, UPPER(TRIM(USEDBY))) setiap Simpan Upload.
CREATE INDEX {skema}.IX_RIRISK_LIFE_SUMMARY_NAMA ON {skema}.RIRISK_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))
/

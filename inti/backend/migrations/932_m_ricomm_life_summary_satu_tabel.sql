-- 932 - ringkasan R/I Comm Life SATU tabel M_RICOMM_LIFE_SUMMARY (langkah 2 dari 2; langkah 1 = 931).
--
-- Keputusan work owner 08-10-2026 (lihat 931). Isi kolom dipindah dari JSONDATA APA ADANYA (notasi titik = definisi
-- view RICOMM_LIFE_SUMMARY, tanpa konversi; nilai yang tidak muat lebar 931 menghentikan langkah dengan ORA-12899, tidak
-- dipotong), lalu JSONDATA + constraint IS JSON dibuang (CASCADE CONSTRAINTS - nama constraint warisan 33 byte tidak
-- ditulis di sini), indeks nama dibuat, dan view RICOMM_LIFE_SUMMARY dibuang TERAKHIR. Kunci pxObjClass (ada di JSON,
-- tidak di view) TIDAK dipindah - hanya di cadangan CSV. Prosedur PEGA_M_RICOMM_LIFE_SUMMARY menjadi INVALID (diterima
-- WO). PRASYARAT: Pega R/I Comm dan backend dihentikan + cadangan CSV ID + JSONDATA (LANGKAH-WO (a)/(b)).
--
-- SETIAP pernyataan aman DIULANG sesudah gagal di tengah: pengisian dan pembuangan JSONDATA lewat blok berpelindung
-- katalog (jalan hanya selama JSONDATA masih ada; pengisian SEBELUM pembuangan), CREATE INDEX (ORA-00955 ditoleransi),
-- dan DROP VIEW terakhir (bila pelari mati SESUDAH DROP VIEW tetapi sebelum mencatat 932, pemulihannya mencatat 932
-- manual). Volume DEV 3 baris. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RICOMM_LIFE_SUMMARY' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.M_RICOMM_LIFE_SUMMARY m SET m.USEDBY = m.JSONDATA.USEDBY, m.MODIFIEDDATE = m.JSONDATA.MODIFIEDDATE, m.OPERATORID = m.JSONDATA.OPERATORID';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RICOMM_LIFE_SUMMARY' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RICOMM_LIFE_SUMMARY DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
-- Indeks nama: pemeriksa nama kembar ricommlife (SqlPemakaiNama, UPPER(TRIM(USEDBY))) setiap Save dan Upload.
CREATE INDEX {skema}.IX_M_RICOMM_LIFE_SUMMARY_NAMA ON {skema}.M_RICOMM_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))
/
-- TERAKHIR: view warisan RICOMM_LIFE_SUMMARY dibuang (definisinya dipulihkan 932_down).
DROP VIEW {skema}.RICOMM_LIFE_SUMMARY
/

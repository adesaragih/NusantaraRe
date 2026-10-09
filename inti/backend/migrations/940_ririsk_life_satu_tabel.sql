-- 940 - rincian R/I Risk SATU tabel RIRISK_LIFE (langkah 3 dari 3).
--
-- Keputusan work owner 08-10-2026 K1/K2. Sumber kebenaran = JSONDATA (yang dibaca view lama): kolom datar warisan BASI
-- (IDUSEDBY NULL di 86 baris, MONTH hanya 1.912 baris - fakta WO), jadi SEMUA kolom view DITIMPA dari JSONDATA:
--   1. isi ulang IDUSEDBY, USEDBY, AGE, YEAR, MONTH, CONTRACT apa adanya (notasi titik) dan RISK teks -> NUMBER TANPA
--      bergantung NLS sesi: TRANSLATE memetakan koma DAN titik ke pemisah desimal sesi itu sendiri (karakter pertama
--      TO_CHAR(0.5), mis. `.5` / `,5`), lalu TO_NUMBER tanpa format - "921,9" -> 921.9, "580.894351210924" ->
--      580.894351210924 di NLS mana pun. Teks yang bukan angka (atau berpemisah ribuan) menghentikan langkah
--      ORA-01722 - tidak ditebak. Ditulis tanpa tanda kutip karena teks EXECUTE IMMEDIATE blok berpelindung tidak boleh
--      memuatnya (CHR(44) = koma, CHR(46) = titik); bukti baca-saja: docs/sql/ririsk_bukti_konversi.sql;
--   2. PK_RIRISK_LIFE (ID) - pola rincian ricommlife (M_RICOMM_LIFE ber-PK); data DEV nol ID ganda. ID NULL / ganda
--      menghentikan langkah SEBELUM JSONDATA dibuang (LANGKAH-WO (a) memeriksanya lebih dulu);
--   3. buang JSONDATA + constraint IS JSON (ENSURE_M_RIRISK_LIFE_JSON) lewat CASCADE CONSTRAINTS;
--   4. indeks IDUSEDBY HANYA bila belum ada indeks berkolom pertama IDUSEDBY (bentuk ketiga BacaPerintahKatalog,
--      SYS.ALL_IND_COLUMNS): kolom INDEX4 warisan tidak tercatat di repo - bila INDEX4 sudah berkolom pertama IDUSEDBY,
--      ia dipakai dan tidak ada indeks kembar (ORA-01408); LANGKAH-WO (a) mencatat kolomnya.
-- Kunci pxObjClass (bukan kolom view) TIDAK dipindah - hanya di cadangan CSV P3. Prosedur PEGA_M_RIRISK_LIFE menjadi
-- INVALID (K3, diterima WO). M_RIRISK_LIFE_TEMP TIDAK disentuh.
--
-- SETIAP pernyataan aman DIULANG (semuanya blok berpelindung katalog; isi SEBELUM buang). Volume DEV 16.398 baris:
-- satu UPDATE, perkiraan detik. NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.RIRISK_LIFE m SET m.IDUSEDBY = m.JSONDATA.IDUSEDBY, m.USEDBY = m.JSONDATA.USEDBY, m.AGE = m.JSONDATA.AGE, m.YEAR = m.JSONDATA.YEAR, m.MONTH = m.JSONDATA.MONTH, m.CONTRACT = m.JSONDATA.CONTRACT, m.RISK = TO_NUMBER(TRANSLATE(TRIM(m.JSONDATA.RISK), CHR(44) || CHR(46), SUBSTR(TO_CHAR(0.5), 1, 1) || SUBSTR(TO_CHAR(0.5), 1, 1)))';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_CONSTRAINTS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE' AND CONSTRAINT_NAME = 'PK_RIRISK_LIFE';
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.RIRISK_LIFE ADD CONSTRAINT PK_RIRISK_LIFE PRIMARY KEY (ID)';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.RIRISK_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_IND_COLUMNS
   WHERE TABLE_OWNER = UPPER('{skema}') AND TABLE_NAME = 'RIRISK_LIFE' AND COLUMN_NAME = 'IDUSEDBY' AND COLUMN_POSITION = 1;
  IF n = 0 THEN
    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.IX_RIRISK_LIFE_IDUSEDBY ON {skema}.RIRISK_LIFE (IDUSEDBY)';
  END IF;
END;
/

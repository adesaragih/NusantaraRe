-- 928 - ringkasan R/I Rate Life menjadi SATU tabel M_RATE_LIFE_SUMMARY (langkah 2 dari 2; langkah 1 = 927).
--
-- Keputusan work owner 07-10-2026 (menggantikan RALAT R4 riratelife): tabel flat RATE_LIFE_SUMMARY (926) adalah SUMBER
-- KEBENARAN sejak aplikasi menulis ke sana (DEV 07-10-2026 09:09); isinya dipindah ke kolom 927, lalu JSONDATA dan
-- tabel flat dibuang. Pega tidak lagi menyimpan R/I Rate (prosedur PEGA_M_RATE_LIFE_SUMMARY dan
-- PEGA_M_PLAN_LIFE_SUMMARY menjadi INVALID - diterima WO). ⛔ PRASYARAT: cadangan CSV ID + JSONDATA (satu-satunya
-- salinan FLAG lama) - LANGKAH-WO-RIRATELIFE-SATU-TABEL.md langkah (a).
--
-- SETIAP pernyataan aman DIULANG sesudah gagal di tengah (pelari tanpa transaksi): UPDATE dan DELETE menurut tabel flat
-- (hasil sama), pembuangan kolom lewat blok berpelindung katalog (bentuk 901), INSERT `NOT EXISTS`, CREATE INDEX
-- (ORA-00955 ditoleransi), dan DROP TABLE TERAKHIR. Urutan menjaga NOT NULL / IS JSON lama: baris yang hanya ada di
-- flat DISISIP SESUDAH kolom lama dibuang, sehingga tidak perlu isi buatan. Constraint IS JSON
-- (ENSURE_M_RATE_LIFE_SUMMARY_JSON) ikut terbuang bersama kolomnya (`CASCADE CONSTRAINTS`).
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
UPDATE {skema}.M_RATE_LIFE_SUMMARY m
   SET (USEDBY, TYPE, MODIFIEDDATE, OPERATORID) =
       (SELECT f.USEDBY, f.TYPE, f.MODIFIEDDATE, f.OPERATORID FROM {skema}.RATE_LIFE_SUMMARY f WHERE f.ID = m.ID)
 WHERE EXISTS (SELECT 1 FROM {skema}.RATE_LIFE_SUMMARY f WHERE f.ID = m.ID)
/
-- Ringkasan yang sudah dihapus aplikasi dari tabel flat tidak boleh hidup lagi (isinya tetap di cadangan CSV (a)).
DELETE FROM {skema}.M_RATE_LIFE_SUMMARY m
 WHERE NOT EXISTS (SELECT 1 FROM {skema}.RATE_LIFE_SUMMARY f WHERE f.ID = m.ID)
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RATE_LIFE_SUMMARY' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RATE_LIFE_SUMMARY DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
INSERT INTO {skema}.M_RATE_LIFE_SUMMARY (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID)
SELECT f.ID, f.USEDBY, f.TYPE, f.MODIFIEDDATE, f.OPERATORID FROM {skema}.RATE_LIFE_SUMMARY f
 WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_RATE_LIFE_SUMMARY m WHERE m.ID = f.ID)
/
-- Indeks nama (pemeriksa nama kembar riratelife); nama baru - IX_RATE_LIFE_SUMMARY_NAMA milik tabel flat sampai DROP.
CREATE INDEX {skema}.IX_M_RATE_LIFE_SUMMARY_NAMA ON {skema}.M_RATE_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))
/
-- TERAKHIR: tabel flat 926 dibuang (indeks dan PK-nya ikut). Bila pelari mati SESUDAH ini tetapi sebelum mencatat 928,
-- pemulihannya mencatat 928 secara manual (LANGKAH-WO, bab Pemulihan).
DROP TABLE {skema}.RATE_LIFE_SUMMARY CASCADE CONSTRAINTS
/

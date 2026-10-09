-- 934 - rincian R/I Comm Life SATU tabel M_RICOMM_LIFE (langkah 2 dari 2; langkah 1 = 933); tabel RICOMM_LIFE (924)
-- dibuang.
--
-- Keputusan work owner 08-10-2026 (lihat 933). Sumber kebenaran rincian = tabel RICOMM_LIFE (aplikasi menulis ke sana
-- sejak 924; DEV 2 baris). M_RICOMM_LIFE (JSON) 0 baris di DEV - LANGKAH-WO (a) BERHENTI bila tidak nol; baris JSON yang
-- ternyata ada tetap diisi kolomnya dari JSON (langkah 1, notasi titik = view RICOMM_LIFE lama) supaya tidak hilang.
--
-- URUTAN (alasan): (1) isi kolom dari JSON selama JSONDATA ada; (2) buang JSONDATA + constraint IS JSON (CASCADE
-- CONSTRAINTS); (3) timpa kolom baris yang ID-nya ada di RICOMM_LIFE (flat = sumber kebenaran); (4) sisip baris
-- RICOMM_LIFE yang belum ada (NOT EXISTS) - SESUDAH JSONDATA dibuang, karena NULLABLE asli JSONDATA tidak diketahui:
-- bila NOT NULL, INSERT tanpa JSONDATA mati di ORA-01400 (pola 928 langkah 4); (5) indeks IDUSEDBY; (6) DROP TABLE
-- RICOMM_LIFE TERAKHIR. Constraint CHECK SYS_C0015437 milik RICOMM_LIFE TIDAK disalin: DDL 924 hanya punya
-- ID ... NOT NULL (Oracle mencatatnya sebagai CHECK "ID" IS NOT NULL) dan ID M_RICOMM_LIFE sudah PK; LANGKAH-WO (a)
-- mencatat SEARCH_CONDITION-nya untuk dipastikan WO.
--
-- SETIAP pernyataan aman DIULANG sampai DROP TABLE: blok berpelindung katalog, UPDATE dari flat (hasil sama), INSERT
-- NOT EXISTS, CREATE INDEX (ORA-00955 ditoleransi). Bila pelari mati SESUDAH DROP TABLE tetapi sebelum mencatat 934,
-- pemulihannya mencatat 934 manual. Prosedur PEGA_M_RICOMM_LIFE menjadi INVALID (diterima WO).
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RICOMM_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'UPDATE {skema}.M_RICOMM_LIFE m SET m.IDUSEDBY = m.JSONDATA.IDUSEDBY, m.USEDBY = m.JSONDATA.USEDBY, m.CONTRACT = m.JSONDATA.CONTRACT, m.YEAR = m.JSONDATA.YEAR, m.COMM = m.JSONDATA.COMM';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RICOMM_LIFE' AND COLUMN_NAME = 'JSONDATA';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RICOMM_LIFE DROP COLUMN JSONDATA CASCADE CONSTRAINTS';
  END IF;
END;
/
UPDATE {skema}.M_RICOMM_LIFE m
   SET (IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM) =
       (SELECT f.IDUSEDBY, f.USEDBY, f.CONTRACT, f.YEAR, f.COMM FROM {skema}.RICOMM_LIFE f WHERE f.ID = m.ID)
 WHERE EXISTS (SELECT 1 FROM {skema}.RICOMM_LIFE f WHERE f.ID = m.ID)
/
INSERT INTO {skema}.M_RICOMM_LIFE (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM)
SELECT f.ID, f.IDUSEDBY, f.USEDBY, f.CONTRACT, f.YEAR, f.COMM FROM {skema}.RICOMM_LIFE f
 WHERE NOT EXISTS (SELECT 1 FROM {skema}.M_RICOMM_LIFE m WHERE m.ID = f.ID)
/
-- Indeks rincian per ringkasan: grid R/I COMM DETAIL, kembar CONTRACT+YEAR, Delete berantai, salinan nama.
CREATE INDEX {skema}.IX_M_RICOMM_LIFE_IDUSEDBY ON {skema}.M_RICOMM_LIFE (IDUSEDBY)
/
-- TERAKHIR: tabel flat 924 dibuang (PK dan indeksnya ikut); bentuknya dipulihkan 934_down.
DROP TABLE {skema}.RICOMM_LIFE CASCADE CONSTRAINTS
/

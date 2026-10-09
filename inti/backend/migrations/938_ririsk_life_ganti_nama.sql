-- 938 - rincian R/I Risk: view RIRISK_LIFE dibuang, tabel Pega M_RIRISK_LIFE diganti nama menjadi RIRISK_LIFE
-- (langkah 1 dari 3; 939 kolom AGE, 940 isi ulang dari JSONDATA + PK + buang JSONDATA + indeks).
--
-- Keputusan work owner 08-10-2026 K1 (pola 935): DROP VIEW hanya bila RIRISK_LIFE masih VIEW (SYS.ALL_VIEWS), RENAME
-- hanya selama M_RIRISK_LIFE masih ada (ALL_TAB_COLUMNS tabel sumber; target yang sudah TABLE = ORA-00955, gagal keras).
-- Data (16.398 baris DEV), indeks warisan INDEX4, dan grant (skema GL punya SELECT - prasyarat P2) ikut nama baru.
-- Pembaca lain (premiumlistlife: View R/I Risk dan hitung QR) tetap memakai nama RIRISK_LIFE.
-- ⛔ PRASYARAT P1-P3 (LANGKAH-WO-RIRISKLIFE.md (a)). Aman DIULANG. NOL COMMIT (ADR-U-0029). -migrate oleh work owner.
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS
   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'RIRISK_LIFE';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'DROP VIEW {skema}.RIRISK_LIFE';
  END IF;
END;
/
DECLARE
  n NUMBER;
BEGIN
  SELECT COUNT(*) INTO n FROM SYS.ALL_TAB_COLUMNS
   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'M_RIRISK_LIFE' AND COLUMN_NAME = 'ID';
  IF n > 0 THEN
    EXECUTE IMMEDIATE 'ALTER TABLE {skema}.M_RIRISK_LIFE RENAME TO RIRISK_LIFE';
  END IF;
END;
/

-- 943 - BENEFIT_LIFE (tabel Pega M_BENEFIT_LIFE, berganti nama di 942) mendapat kolom BENEFIT (langkah 2 dari 3;
-- 944 isi + pemeriksaan K2 + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K1: kolom akhir = kolom view lama (ID, BENEFIT); lebar BENEFIT = pola ririsklife
-- untuk kolom nama (USEDBY 936 = 200) karena XML tidak menetapkan batas (pxTextArea tanpa pyMaxChars, InboxBenefit.xml
-- b1175) [data DEV 08-10-2026: maks 48 byte].
-- Berdiri sendiri: ADD diulang mati di ORA-01430 dan penjaga inti menolak kolom baru lewat blok di jalur maju.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.BENEFIT_LIFE ADD (
  BENEFIT       VARCHAR2(200)
)
/

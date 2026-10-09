-- 936 - RIRISK_LIFE_SUMMARY (tabel Pega M_RIRISK_LIFE_SUMMARY, berganti nama di 935) mendapat kolom ringkasan
-- (langkah 2 dari 3; 937 isi + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K2: kolom akhir = kolom view lama (ID, USEDBY, MODIFIEDDATE, OPERATORID); tipe dan
-- lebar kolom baru = pola ricommlife 931: USEDBY VARCHAR2(200) [terverifikasi data DEV 08-10-2026: maks 81 byte; sama
-- dengan models.BatasNama], MODIFIEDDATE VARCHAR2(50) (bentuk Pega YYYYMMDDTHHMMSS.mmm GMT), OPERATORID VARCHAR2(200)
-- [maks 17 byte].
-- Berdiri sendiri: ADD diulang mati di ORA-01430 dan penjaga inti menolak kolom baru lewat blok di jalur maju.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.RIRISK_LIFE_SUMMARY ADD (
  USEDBY        VARCHAR2(200),
  MODIFIEDDATE  VARCHAR2(50),
  OPERATORID    VARCHAR2(200)
)
/

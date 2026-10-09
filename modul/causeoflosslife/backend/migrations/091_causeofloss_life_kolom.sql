-- 091 - CAUSEOFLOSS_LIFE (tabel Pega M_CAUSEOFLOSS_LIFE, berganti nama di 090) mendapat kolom CAUSEOFLOSS (langkah 2
-- dari 3; 092 isi + pemeriksaan K1.4 + buang JSONDATA).
--
-- Keputusan work owner 08-10-2026 K1: kolom akhir = kolom view lama (ID, CAUSEOFLOSS); lebar 200 = pola Benefit 943
-- karena XML tidak menetapkan batas (pxTextInput tanpa pyMaxChars, InboxCauseofLossLife.xml b1029) [data DEV
-- 08-10-2026: maks 9 byte]. NULLABLE: baris 100001 (`"CauseofLoss":""`) terbaca NULL - tidak diisi, tidak dihapus.
-- Berdiri sendiri: ADD diulang mati di ORA-01430 dan penjaga inti menolak kolom baru lewat blok di jalur maju.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.CAUSEOFLOSS_LIFE ADD (
  CAUSEOFLOSS   VARCHAR2(200)
)
/

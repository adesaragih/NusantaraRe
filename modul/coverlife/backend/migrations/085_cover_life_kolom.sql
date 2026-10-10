-- 085 - Cover Life (coverlife): tabel Pega M_COVER_LIFE mendapat kolom COVER dan NOTE (langkah 1 dari 2; 086 isi +
-- pemeriksaan C1.3 + buang JSONDATA + buang view COVER_LIFE).
--
-- Keputusan work owner 08-10-2026 C1: nama tabel TETAP M_COVER_LIFE (TANPA RENAME), dijadikan flat di tempat; kolom
-- akhir = kolom view lama COVER_LIFE bernama sama, urutan sama (ID warisan, COVER, NOTE).
--   - COVER VARCHAR2(200) (C1.1; data DEV maks 17 byte). NULLABLE: nilai dari view apa adanya; wajib dijaga aplikasi.
--   - NOTE  VARCHAR2(1000): XML tidak menetapkan batas (pxTextArea tanpa pyMaxChars, InboxCoverLife.xml b1058) -> 1000
--     (C1.1 "lebar dari XML atau 1000"). NULL di 4/4 baris DEV (kunci Note tidak ada di JSONDATA).
-- Berdiri sendiri: ADD diulang mati di ORA-01430 dan penjaga inti menolak kolom baru lewat blok di jalur maju.
-- K0: migrasi MODUL (rentang 085-089, dipinjam dari jatah premiumlistlife) - di skema baru berjalan SEBELUM 900.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner (LANGKAH-WO-COVERLIFE.md).
ALTER TABLE {skema}.M_COVER_LIFE ADD (
  COVER   VARCHAR2(200),
  NOTE    VARCHAR2(1000)
)
/

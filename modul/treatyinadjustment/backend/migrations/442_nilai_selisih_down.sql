-- Jalur mundur untuk 442_nilai_selisih.sql
--
-- ⛔ WAJIB dijalankan BERSAMA `treatyin/.../435_cabut_nilai_selisih_down.sql`,
-- yang membangun keduanya kembali di modul `treatyin`. Menjalankan yang ini
-- sendirian meninggalkan tiket 76 dan 77 tanpa tabel sama sekali.
--
-- ⚠️ Membongkar sequence menghilangkan nilai terakhirnya; memasangnya
-- kembali memulai dari 1 dan menawarkan ulang pengenal yang sudah terpakai
-- (INV-03). Aman selama keduanya nol baris.
DROP TABLE {skema}.NILAI_SEBELUM_PRO_RATE
/
DROP TABLE {skema}.NILAI_SELISIH
/
DROP SEQUENCE {skema}.SEQ_TRIA_NILAI_SBL_PRORATA
/
DROP SEQUENCE {skema}.SEQ_TRIA_NILAI_SELISIH
/

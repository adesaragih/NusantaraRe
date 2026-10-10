-- Jalur mundur untuk 426_nilai_selisih.sql
--
-- ⚠️ Membongkar ini mengembalikan lubang spec yang menahan tiket 06, 11, dan
-- 13 - dan ketiganya berstatus `aktif`, sehingga lubangnya kembali tidak
-- terlihat di pencacah papan.
--
-- ⚠️ Membongkar sequence menghilangkan nilai terakhirnya; memasangnya kembali
-- memulai dari 1 dan menawarkan ulang pengenal yang sudah terpakai (INV-03).
DROP TABLE {skema}.NILAI_SEBELUM_PRO_RATE
/
DROP TABLE {skema}.NILAI_SELISIH
/
DROP SEQUENCE {skema}.SEQ_TRIA_NILAI_SBL_PRORATA
/
DROP SEQUENCE {skema}.SEQ_TRIA_NILAI_SELISIH
/

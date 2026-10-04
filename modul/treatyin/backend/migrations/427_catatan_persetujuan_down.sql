-- Jalur mundur untuk 427_catatan_persetujuan.sql
--
-- ⚠️ Membongkar ini membuang jejak persetujuan. Perilaku hapus `tolak` pada
-- kunci asingnya sengaja melindungi baris itu dari penghapusan versi;
-- DROP TABLE melewati perlindungan itu.
--
-- ⚠️ Membongkar sequence menghilangkan nilai terakhirnya; memasangnya kembali
-- memulai dari 1 dan menawarkan ulang pengenal yang sudah terpakai (INV-03).
DROP TABLE {skema}.CATATAN_PERSETUJUAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_CATATAN_PERSETUJUAN
/

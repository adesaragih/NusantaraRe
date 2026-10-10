-- Jalur mundur untuk 423_pencapaian.sql
--
-- ⚠️ Membongkar ini menghapus angka yang pernah dibukukan bila tabelnya sudah
-- berisi. Perilaku hapus `tolak` pada kunci asingnya sengaja melindungi baris
-- itu dari penghapusan kontrak; DROP TABLE melewati perlindungan itu.
DROP TABLE {skema}.PENCAPAIAN
/

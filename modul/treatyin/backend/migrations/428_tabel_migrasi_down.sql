-- Jalur mundur untuk 428_tabel_migrasi.sql
--
-- ⚠️ Membongkar ini membuang jalan menelusuri baris baru kembali ke asalnya.
-- Bila tiket 44 sudah memuat data, penelusuran itu TIDAK DAPAT dibangun ulang
-- dari apa pun - sistem lama mungkin sudah dimatikan saat kekeliruannya
-- ketahuan.
--
-- Urutan terbalik dari jalur maju: nilai ditolak lebih dulu, sebab ia
-- menggantung pada pendaratan.
DROP TABLE {skema}.MIGRASI_NILAI_DITOLAK
/
DROP TABLE {skema}.MIGRASI_PENDARATAN
/
DROP TABLE {skema}.MIGRASI_KORELASI
/
DROP SEQUENCE {skema}.SEQ_TRIN_MIG_NILAI_DITOLAK
/
DROP SEQUENCE {skema}.SEQ_TRIN_MIGRASI_PENDARATAN
/
DROP SEQUENCE {skema}.SEQ_TRIN_MIGRASI_KORELASI
/

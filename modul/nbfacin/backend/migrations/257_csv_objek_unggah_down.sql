-- Jalur mundur 257 - unggah CSV objek: tiga tabel + sequence. T_LOCATIONLIST.SEQ_NO TIDAK dipersempit kembali ke NUMBER(5)
-- (Oracle menolak mempersempit presisi kolom berisi data, ORA-01440; dan lokasi > 99999 bisa sudah tersimpan).
DROP TABLE {skema}.T_NB_CSV_PESAN CASCADE CONSTRAINTS
/
DROP TABLE {skema}.T_NB_CSV_BARIS CASCADE CONSTRAINTS
/
DROP SEQUENCE {skema}.SEQ_T_NB_CSV_UNGGAH
/
DROP TABLE {skema}.T_NB_CSV_UNGGAH CASCADE CONSTRAINTS
/

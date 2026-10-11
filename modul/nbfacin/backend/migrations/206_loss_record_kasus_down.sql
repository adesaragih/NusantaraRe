-- Jalur mundur 206 - dua tabel Loss Record tingkat case dan sequence-nya (anak -> induk). Data Loss Record HILANG.
DROP SEQUENCE {skema}.SEQ_T_LOSSRECORDLIST
/
DROP TABLE {skema}.T_LOSSRECORDLIST
/
DROP SEQUENCE {skema}.SEQ_T_LOSSRECORD
/
DROP TABLE {skema}.T_LOSSRECORD
/

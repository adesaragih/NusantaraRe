-- Jalur mundur untuk 443_kunci_alami_nilai.sql
--
-- ⚠️ Melepas kunci alami mengembalikan keadaan yang jalur gagal kedua tiket
-- sebut: dua baris berbesaran dan bermata-uang sama pada satu versi menjadi
-- DAPAT masuk, dan tidak ada yang akan berbunyi.
--
-- ⚠️ Mengembalikan lebar kolom ke `VARCHAR2(1000 CHAR)` juga mengembalikan
-- sebab kunci alami itu tidak dapat dipasang: pada AL32UTF8 kuncinya menjadi
-- 8.022 bita, melewati batas panjang kunci index Oracle.
--
-- Urutan terbalik: constraint dulu, baru lebarnya — kolom yang memimpin
-- sebuah UNIQUE tidak dapat diubah lebarnya selama constraintnya berdiri.
ALTER TABLE {skema}.NILAI_SEBELUM_PRO_RATE DROP CONSTRAINT UQ_NILAI_SEBELUM_PRO_RATE
/
ALTER TABLE {skema}.NILAI_SELISIH DROP CONSTRAINT UQ_NILAI_SELISIH
/
ALTER TABLE {skema}.NILAI_SEBELUM_PRO_RATE MODIFY (KODE_MATA_UANG VARCHAR2(1000 CHAR))
/
ALTER TABLE {skema}.NILAI_SEBELUM_PRO_RATE MODIFY (KODE_BESARAN VARCHAR2(1000 CHAR))
/
ALTER TABLE {skema}.NILAI_SELISIH MODIFY (KODE_MATA_UANG VARCHAR2(1000 CHAR))
/
ALTER TABLE {skema}.NILAI_SELISIH MODIFY (KODE_BESARAN VARCHAR2(1000 CHAR))
/

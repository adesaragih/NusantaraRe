-- Mundur 929: buang ketujuh kolom rincian dari M_RATE_LIFE. Berjalan SESUDAH 930_down (urutan mundur menurun), yang
-- lebih dulu membangun ulang JSONDATA dan view RATE_LIFE dari kolom ini.
--
-- ⛔ PELINDUNG GAGAL-KERAS (pernyataan pertama): kolom ini hanya boleh dibuang selama JSONDATA masih ada. Bila 930
-- sudah membuang JSONDATA tetapi TIDAK tercatat di T_MIGRASI (pelari mati di tengah 930), Bongkar melewati 930_down lalu
-- akan membuang satu-satunya salinan isi rincian. UPDATE nol baris di bawah dirujuk Oracle saat parse: tanpa kolom
-- JSONDATA (atau kolom itu setengah-terbuang) ia gagal ORA-00904 "JSONDATA": invalid identifier, Bongkar berhenti
-- SEBELUM DROP dan SEBELUM menghapus catatan 929. Bila JSONDATA ada, ia tidak mengubah apa pun (WHERE 1 = 0).
-- Pemulihan keadaan itu: selesaikan langkah MAJU (LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md, bab Pemulihan).
UPDATE {skema}.M_RATE_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.M_RATE_LIFE DROP (IDUSEDBY, USEDBY, TYPE, GENDER, CONTRACT, AGE, RATE)
/

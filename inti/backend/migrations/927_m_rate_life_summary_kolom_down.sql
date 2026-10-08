-- Mundur 927: buang keempat kolom ringkasan dari M_RATE_LIFE_SUMMARY. Berjalan SESUDAH 928_down (urutan mundur
-- menurun), yang lebih dulu mengembalikan tabel flat RATE_LIFE_SUMMARY dan JSONDATA dari kolom ini.
ALTER TABLE {skema}.M_RATE_LIFE_SUMMARY DROP (USEDBY, TYPE, MODIFIEDDATE, OPERATORID)
/

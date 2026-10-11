-- Jalur mundur 211 - posisi case Fac In ReasFacInJuniorUnderwritingB kembali ke nama lama ReasFacInJuniorUnderwriting.
-- ⚠️ Baris yang SEJAK AWAL ditulis biner baru dengan kode B ikut dikembalikan (tidak dapat dibedakan); jalankan hanya bersama
-- biner lama.
UPDATE {skema}.T_WORK_POLIS SET POSITION = 'ReasFacInJuniorUnderwriting'
 WHERE LINI = 'FAC' AND POSITION = 'ReasFacInJuniorUnderwritingB'
/

-- Mundur 936: buang ketiga kolom ringkasan dari RIRISK_LIFE_SUMMARY. Berjalan SESUDAH 937_down, yang lebih dulu
-- membangun ulang JSONDATA dari kolom ini.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 929_down / 931_down): tanpa JSONDATA (937 membuangnya tetapi tidak
-- tercatat) UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM DROP dan sebelum menghapus
-- catatan 936. Bila JSONDATA ada, tidak ada yang berubah (WHERE 1 = 0).
UPDATE {skema}.RIRISK_LIFE_SUMMARY SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.RIRISK_LIFE_SUMMARY DROP (USEDBY, MODIFIEDDATE, OPERATORID)
/

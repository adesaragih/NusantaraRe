-- Mundur 939: buang AGE dari RIRISK_LIFE. Berjalan SESUDAH 940_down, yang lebih dulu membangun ulang JSONDATA.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 929_down): tanpa JSONDATA (940 membuangnya tetapi tidak tercatat)
-- UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti sebelum DROP dan sebelum menghapus catatan 939.
UPDATE {skema}.RIRISK_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.RIRISK_LIFE DROP (AGE)
/

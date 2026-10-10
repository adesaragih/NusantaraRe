-- Mundur 943: buang kolom BENEFIT dari BENEFIT_LIFE. Berjalan SESUDAH 944_down, yang lebih dulu membangun ulang
-- JSONDATA dari kolom ini.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 936_down / 939_down): tanpa JSONDATA (944 membuangnya tetapi tidak
-- tercatat) UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM DROP dan sebelum menghapus
-- catatan 943. Bila JSONDATA ada, tidak ada yang berubah (WHERE 1 = 0).
UPDATE {skema}.BENEFIT_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.BENEFIT_LIFE DROP (BENEFIT)
/

-- Mundur 091: buang kolom CAUSEOFLOSS dari CAUSEOFLOSS_LIFE. Berjalan SESUDAH 092_down, yang lebih dulu membangun ulang
-- JSONDATA dari kolom ini.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 943_down): tanpa JSONDATA (092 membuangnya tetapi tidak tercatat)
-- UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM DROP dan sebelum menghapus catatan 091.
-- Bila JSONDATA ada, tidak ada yang berubah (WHERE 1 = 0).
UPDATE {skema}.CAUSEOFLOSS_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.CAUSEOFLOSS_LIFE DROP (CAUSEOFLOSS)
/

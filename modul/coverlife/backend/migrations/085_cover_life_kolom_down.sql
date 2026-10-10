-- Mundur 085: buang kolom COVER dan NOTE dari M_COVER_LIFE. Berjalan SESUDAH 086_down, yang lebih dulu membangun ulang
-- JSONDATA dari kedua kolom ini dan view COVER_LIFE.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 091_down / 943_down): tanpa JSONDATA (086 membuangnya tetapi tidak
-- tercatat) UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM DROP - isi hanya tinggal di kolom
-- ini dan tidak boleh hilang. Bila JSONDATA ada, tidak ada yang berubah (WHERE 1 = 0).
UPDATE {skema}.M_COVER_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.M_COVER_LIFE DROP (COVER, NOTE)
/

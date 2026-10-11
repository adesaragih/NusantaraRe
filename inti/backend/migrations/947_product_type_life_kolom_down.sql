-- Mundur 947: buang kelima kolom dari PRODUCT_TYPE_LIFE. Berjalan SESUDAH 948_down, yang lebih dulu membangun ulang
-- JSONDATA dari kolom ini.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 943_down): tanpa JSONDATA (948 membuangnya tetapi tidak tercatat)
-- UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM DROP. Bila JSONDATA ada, tidak ada yang
-- berubah (WHERE 1 = 0).
UPDATE {skema}.PRODUCT_TYPE_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.PRODUCT_TYPE_LIFE DROP (COVERNAME, BUSINESS, BUSINESSID, BENEFIT, BENEFITID)
/

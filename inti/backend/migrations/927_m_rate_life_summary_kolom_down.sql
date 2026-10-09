-- Mundur 927: buang keempat kolom ringkasan dari M_RATE_LIFE_SUMMARY. Berjalan SESUDAH 928_down (urutan mundur
-- menurun), yang lebih dulu mengembalikan tabel flat RATE_LIFE_SUMMARY dan JSONDATA dari kolom ini.
--
-- ⛔ PELINDUNG GAGAL-KERAS (pernyataan pertama, pola sama dengan 929_down): kolom ini hanya boleh dibuang selama
-- JSONDATA masih ada. Bila 928 sudah membuang JSONDATA tetapi TIDAK tercatat di T_MIGRASI (pelari mati di tengah 928),
-- Bongkar melewati 928_down lalu akan membuang kolom yang - sesudah DROP TABLE flat - satu-satunya salinan ringkasan.
-- Tanpa kolom JSONDATA UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM DROP dan SEBELUM
-- menghapus catatan 927. Bila JSONDATA ada (928 gagal sebelum membuangnya), tidak ada yang berubah (WHERE 1 = 0) dan
-- mundur aman: JSON + tabel flat masih ada (kecuali baris yang sudah dibuang langkah 2 DELETE 928 - hanya di cadangan
-- CSV LANGKAH-WO-RIRATELIFE-SATU-TABEL.md (a)).
UPDATE {skema}.M_RATE_LIFE_SUMMARY SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.M_RATE_LIFE_SUMMARY DROP (USEDBY, TYPE, MODIFIEDDATE, OPERATORID)
/

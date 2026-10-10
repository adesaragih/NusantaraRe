-- Mundur 931: buang ketiga kolom ringkasan dari M_RICOMM_LIFE_SUMMARY. Berjalan SESUDAH 932_down (urutan mundur
-- menurun), yang lebih dulu membangun ulang JSONDATA dan view RICOMM_LIFE_SUMMARY dari kolom ini.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 929_down): kolom ini hanya boleh dibuang selama JSONDATA masih ada.
-- Bila 932 sudah membuang JSONDATA tetapi TIDAK tercatat di T_MIGRASI, Bongkar melewati 932_down lalu akan membuang
-- satu-satunya salinan ringkasan. Tanpa kolom JSONDATA (atau setengah-terbuang) UPDATE nol baris ini gagal ORA-00904
-- saat parse; Bongkar berhenti SEBELUM DROP dan SEBELUM menghapus catatan 931. Bila JSONDATA ada, tidak ada yang
-- berubah (WHERE 1 = 0). Pemulihan: selesaikan langkah MAJU (LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md, Pemulihan).
UPDATE {skema}.M_RICOMM_LIFE_SUMMARY SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.M_RICOMM_LIFE_SUMMARY DROP (USEDBY, MODIFIEDDATE, OPERATORID)
/

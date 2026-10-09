-- Mundur 933: buang kelima kolom rincian dari M_RICOMM_LIFE. Berjalan SESUDAH 934_down, yang lebih dulu membangun
-- ulang tabel RICOMM_LIFE (berisi data) dan JSONDATA dari kolom ini.
--
-- PELINDUNG GAGAL-KERAS (pernyataan pertama, pola 929_down): tanpa kolom JSONDATA (934 membuangnya tetapi tidak
-- tercatat) UPDATE nol baris ini gagal ORA-00904 saat parse; Bongkar berhenti SEBELUM DROP dan SEBELUM menghapus
-- catatan 933. Bila JSONDATA ada, tidak ada yang berubah (WHERE 1 = 0).
UPDATE {skema}.M_RICOMM_LIFE SET JSONDATA = JSONDATA WHERE 1 = 0
/
ALTER TABLE {skema}.M_RICOMM_LIFE DROP (IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM)
/

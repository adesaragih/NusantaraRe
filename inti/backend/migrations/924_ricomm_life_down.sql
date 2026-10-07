-- Mundur 924: buang tabel flat RICOMM_LIFE (indeksnya ikut), lalu pulihkan VIEW warisan RICOMM_LIFE persis definisi
-- DEV (dicek work owner 06-10-2026) atas M_RICOMM_LIFE yang tidak pernah disentuh.
-- ⚠️ Baris yang ditulis aplikasi ke tabel flat (tidak ada di M_RICOMM_LIFE) HILANG - jalur mundur ini hanya aman
-- sebelum aplikasi menulis (modul/ricommlife/MODUL.md, bab urutan langkah).
DROP TABLE {skema}.RICOMM_LIFE CASCADE CONSTRAINTS
/
CREATE VIEW {skema}.RICOMM_LIFE AS
SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.CONTRACT, a.JSONDATA.YEAR, a.JSONDATA.COMM
FROM {skema}.M_RICOMM_LIFE a
/

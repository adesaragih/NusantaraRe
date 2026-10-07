-- Mundur 926: buang tabel flat RATE_LIFE_SUMMARY (indeksnya ikut), lalu pulihkan VIEW warisan RATE_LIFE_SUMMARY persis
-- definisi DEV (ALL_VIEWS, dibaca WO 06-10-2026) atas M_RATE_LIFE_SUMMARY yang tidak pernah disentuh - termasuk FLAG
-- (tabel flat tidak memuatnya: keputusan work owner 06-10-2026, FLAG tidak digunakan).
-- ⚠️ Ringkasan yang ditulis aplikasi ke tabel flat (tidak ada di M_RATE_LIFE_SUMMARY) HILANG - jalur mundur ini hanya
-- aman sebelum aplikasi menulis (modul/riratelife/docs/LANGKAH-WO-RIRATELIFE-FLAT.md).
DROP TABLE {skema}.RATE_LIFE_SUMMARY CASCADE CONSTRAINTS
/
CREATE VIEW {skema}.RATE_LIFE_SUMMARY AS
SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID, a.JSONDATA.FLAG
FROM {skema}.M_RATE_LIFE_SUMMARY a
/

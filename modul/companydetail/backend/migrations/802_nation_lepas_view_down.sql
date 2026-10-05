-- Mundur 802: buat ulang view NATION persis teks aslinya (katalog DEV 04-10-2026).
CREATE OR REPLACE VIEW {skema}.NATION AS
SELECT a.JSONDATA.ID, a.OLDID, a.JSONDATA.Note, a.JSONDATA.NationInitial
  FROM {skema}.M_NATION a
/

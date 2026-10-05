-- 804 - salin isi NATION dari M_NATION SEKALI, dengan ekspresi yang SAMA dengan view lamanya (tiga kolom dari
-- dokumen M_NATION, OLDID dari kolomnya). Satu pernyataan - atomik.
--
-- ⚠️ Sesudah langkah ini NATION tidak lagi mengikuti M_NATION: negara yang ditambah atau diubah lewat form Pega
-- (InputNation, UpdateMasterNation) masuk M_NATION, tidak ke NATION. Disampaikan ke work owner 04-10-2026.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
INSERT INTO {skema}.NATION (ID, OLDID, NOTE, NATIONINITIAL)
SELECT a.JSONDATA.ID, a.OLDID, a.JSONDATA.Note, a.JSONDATA.NationInitial
  FROM {skema}.M_NATION a
/

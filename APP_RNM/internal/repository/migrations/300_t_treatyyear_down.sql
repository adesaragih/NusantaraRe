-- Jalur mundur 300 - T_TREATYYEAR.
--
-- CASCADE CONSTRAINTS membuang FK dari T_TREATYCONTRACT bila anaknya belum
-- sempat dibongkar, sehingga jalur mundur tidak bergantung pada urutan.
DROP TABLE {skema}.T_TREATYYEAR CASCADE CONSTRAINTS
/

DROP SEQUENCE {skema}.SEQ_T_TREATYYEAR
/

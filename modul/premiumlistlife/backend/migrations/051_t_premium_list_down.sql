-- Jalur mundur 051 - T_PREMIUM_LIST.
--
-- Index ikut terbuang bersama tabelnya. CASCADE CONSTRAINTS membuang
-- rujukan dari tabel anak yang mungkin belum sempat dibuang, sehingga
-- jalur mundur tidak bergantung pada urutan.
DROP TABLE {skema}.T_PREMIUM_LIST CASCADE CONSTRAINTS
/

-- Jalur mundur 056 - T_VIEW_SUGGEST.
--
-- Index ikut terbuang bersama tabelnya. CASCADE CONSTRAINTS membuang
-- rujukan dari tabel anak yang mungkin belum sempat dibuang, sehingga
-- jalur mundur tidak bergantung pada urutan.
DROP TABLE {skema}.T_VIEW_SUGGEST CASCADE CONSTRAINTS
/

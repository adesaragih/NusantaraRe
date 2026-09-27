-- Jalur mundur 055 - T_PREMIUM_LIST_SUMMARY.
--
-- Index ikut terbuang bersama tabelnya. CASCADE CONSTRAINTS membuang
-- rujukan dari tabel anak yang mungkin belum sempat dibuang, sehingga
-- jalur mundur tidak bergantung pada urutan.
DROP TABLE {skema}.T_PREMIUM_LIST_SUMMARY CASCADE CONSTRAINTS
/

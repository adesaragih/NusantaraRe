-- Mundur 760: buang LOG_TIME dan AKSES_LOGIN dari MARKETINGOFFICER_LOG (isinya hilang). Trigger warisan tidak pernah
-- diubah, jadi tidak ada yang dikembalikan. Tanda ACTION = 'UPDATE-GO' pada baris log tetap tinggal (kolom warisan).
ALTER TABLE {skema}.MARKETINGOFFICER_LOG DROP (LOG_TIME, AKSES_LOGIN)
/

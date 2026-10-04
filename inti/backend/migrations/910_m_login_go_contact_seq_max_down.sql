-- Mundur 910: TIDAK ada yang dibatalkan - sequence yang sudah mengikuti CON tertinggi tetap dipakai. Mengembalikannya
-- ke posisi lama justru membuat nomor CON akun baru bentrok dengan akun yang sudah ada (UX_M_LOGIN_GO_CONTACT_ID).
-- Pernyataan di bawah tidak mengubah apa pun (sequence-nya memang NOCACHE); mundur 905 sesudahnya yang membuangnya.
ALTER SEQUENCE {skema}.M_LOGIN_GO_CONTACT_SEQ NOCACHE
/

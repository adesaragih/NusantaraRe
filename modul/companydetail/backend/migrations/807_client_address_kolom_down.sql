-- Mundur 807: buang kolom telfax dari CLIENT_ADDRESS (isinya hilang; baris tambahan per nomor tetap tinggal).
ALTER TABLE {skema}.CLIENT_ADDRESS DROP (TELFAX_TYPE, TELFAX_CODE, TELFAX_NO)
/

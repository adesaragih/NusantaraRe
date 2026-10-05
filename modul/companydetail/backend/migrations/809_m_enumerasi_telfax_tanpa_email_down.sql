-- Mundur 809: EMAIL kembali menjadi pilihan Phone and Fax.
UPDATE {skema}.M_ENUMERASI SET AKTIF = '1'
WHERE JENIS = 'telfax' AND KODE = '6'
/

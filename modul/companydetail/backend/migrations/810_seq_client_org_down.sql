-- Mundur 810: buang sequence nomor ORG. Sesudahnya Create menjawab "belum dimigrasi" (ORA-02289) sampai 810
-- dijalankan lagi - nilai awalnya dihitung ulang dari nomor tertinggi saat itu.
DROP SEQUENCE {skema}.SEQ_CLIENT_ORG
/

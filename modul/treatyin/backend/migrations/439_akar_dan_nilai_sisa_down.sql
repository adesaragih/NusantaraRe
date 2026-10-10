-- Pembalikan `439_akar_dan_nilai_sisa.sql` -- keempat tabel beserta
-- sequence-nya dibongkar.
--
-- ⚠ Urutannya KEBALIKAN dari yang membangun, dan itu wajib:
-- `T_TREATY_LIMIT_MEASURE` menunjuk `T_TREATY_LIMITS` lewat kunci asing, jadi
-- ia harus pergi lebih dulu. `CASCADE CONSTRAINTS` menutup sisanya.
--
-- Selama keempatnya nol baris, pembalikan ini tidak kehilangan apa pun.

DROP TABLE {skema}.T_TREATY_TOTAL CASCADE CONSTRAINTS
/
DROP TABLE {skema}.T_TREATY_LIMIT_SUMMARY CASCADE CONSTRAINTS
/
DROP TABLE {skema}.T_TREATY_LIMIT_MEASURE CASCADE CONSTRAINTS
/
DROP TABLE {skema}.T_TREATY_REVISION CASCADE CONSTRAINTS
/
DROP SEQUENCE {skema}.SEQ_TT_TOTAL
/
DROP SEQUENCE {skema}.SEQ_TT_LIMIT_SUMMARY
/
DROP SEQUENCE {skema}.SEQ_TT_LIMIT_MEASURE
/
DROP SEQUENCE {skema}.SEQ_TT_REVISION
/

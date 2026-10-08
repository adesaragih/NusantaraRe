-- Pembalikan `438_nilai_anak_treaty_in.sql` — ketiga tabel nilai dan
-- sequence-nya dicabut.
--
-- ⚠ Urutannya tidak mengikat di sini: ketiganya DAUN, nol tabel menunjuk
-- mereka. Yang mengikat pada `437` — anak sebelum induk — tidak berlaku,
-- sebab ketiganya tidak punya anak.
--
-- ⛔ Nol data hilang: ketiga tabel ini lahir KOSONG di `438` dan pemuatnya
-- belum dibangun.

DROP TABLE {skema}.T_TREATY_FAC_SHARE_AMOUNT CASCADE CONSTRAINTS
/
DROP TABLE {skema}.T_TREATY_SHARE_AMOUNT CASCADE CONSTRAINTS
/
DROP TABLE {skema}.T_TREATY_LIMIT_AMOUNT CASCADE CONSTRAINTS
/
DROP SEQUENCE {skema}.SEQ_TT_FAC_SHARE_AMOUNT
/
DROP SEQUENCE {skema}.SEQ_TT_SHARE_AMOUNT
/
DROP SEQUENCE {skema}.SEQ_TT_LIMIT_AMOUNT
/

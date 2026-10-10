-- Mundur 641 - ADJUSTMENT_ID kembali wajib isi untuk baris BARU (NOT NULL NOVALIDATE): baris kasus komite TT3 / TT4
-- ber-ADJUSTMENT_ID kosong yang sudah ada TIDAK diperiksa dan tidak dihapus, sehingga mundur tetap berhasil walau baris
-- itu ada (NOT NULL biasa akan gagal ORA-02296). Penyerahan TT3 / TT4 berikutnya ditolak basis data. Untuk skema uji.
--
-- NOL COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_KOMITE MODIFY (ADJUSTMENT_ID NOT NULL NOVALIDATE)
/

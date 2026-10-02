-- Jalur mundur untuk 441_nomor_urut_versi_boleh_kosong.sql
--
-- Mengembalikan NOT NULL yang tiket 14 pasang.
--
-- ⚠️ Jalur mundur ini GAGAL bila sudah ada baris bernomor urut kosong - dan itu
-- benar, bukan cacat: mengetatkan kolom yang memuat NULL tidak punya jawaban
-- yang aman. Oracle menolak dengan ORA-02296. Yang hendak mundur sesudah data
-- warisan masuk harus mengisi nomornya lebih dulu (tiket 10).
ALTER TABLE {skema}.VERSI_KONTRAK MODIFY (
  NOMOR_URUT_VERSI  NUMBER(10) NOT NULL
)
/

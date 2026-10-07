-- 483 - sub-jenis EDM Type layar Create Addendum (OQ-EDM-005 DITUTUP, keputusan
-- work owner 07-10-2026).
--
-- `EndorsmentLife_Section.xml` b1645 `EDM Type Perubahan Data` (`.EdmTypePerubahanData`,
-- tampil bila EDM Type 1) dan b1937 `EDM Type Batal` (`.EdmTypeBatal`, tampil bila
-- EDM Type 3). Opsinya tidak diekspor korpus; daftarnya diberikan work owner dari
-- layar Pega dan disimpan sebagai TEKS LABEL persis (kode Pega tidak diketahui) -
-- `models.OpsiEdmTypePerubahanData` / `models.OpsiEdmTypeBatal`.
--
-- ⛔ Dua kolom NULLABLE tanpa DEFAULT: baris new business PremiumList dan kasus EDM
-- yang lahir sebelum 483 tetap kosong (nol isi mundur). Hanya ditulis jalur buat
-- kasus Endorsement Life; PremiumList Life tidak menulis maupun membacanya.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). Batas transaksi milik pelari migrasi.
ALTER TABLE {skema}.T_PREMIUM_LIST ADD (
  EDM_TYPE_PERUBAHAN_DATA VARCHAR2(100),
  EDM_TYPE_BATAL          VARCHAR2(100)
)
/

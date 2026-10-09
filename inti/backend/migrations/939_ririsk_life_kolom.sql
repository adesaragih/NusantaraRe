-- 939 - RIRISK_LIFE (tabel Pega M_RIRISK_LIFE, berganti nama di 938) mendapat kolom AGE (langkah 2 dari 3).
--
-- Keputusan work owner 08-10-2026 K2: kolom akhir = kolom view lama RIRISK_LIFE (ID, IDUSEDBY, USEDBY, AGE, YEAR,
-- MONTH, RISK, CONTRACT). Kolom datar IDUSEDBY VARCHAR2(100), USEDBY VARCHAR2(1000), YEAR/MONTH/CONTRACT VARCHAR2(10),
-- RISK NUMBER SUDAH ADA (tipe warisan dipakai apa adanya - RISK tanpa skala supaya 580.894351210924 tidak terpotong);
-- yang belum ada hanya AGE. AGE VARCHAR2(10) = tipe teks saudaranya (YEAR / MONTH / CONTRACT) dan riratelife 929 -
-- [terverifikasi data DEV 08-10-2026] kunci AGE tidak ada di satu baris JSON pun (view mengeluarkan NULL); tidak ada
-- medan AGE di XML InboxRIRisk, jadi modul tidak menulisnya.
-- Berdiri sendiri (alasan sama dengan 936). NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.RIRISK_LIFE ADD (
  AGE  VARCHAR2(10)
)
/

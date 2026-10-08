-- 927 - M_RATE_LIFE_SUMMARY mendapat kolom ringkasan (langkah 1 dari 2; langkah 2 = 928).
--
-- Keputusan work owner 07-10-2026 (menggantikan RALAT R4 riratelife): "ringkasan R/I Rate Life cukup SATU tabel -
-- M_RATE_LIFE_SUMMARY". Sasaran akhir: M_RATE_LIFE_SUMMARY berkolom persis seperti tabel flat RATE_LIFE_SUMMARY (926):
-- ID VARCHAR2(10) PK + USEDBY, TYPE, MODIFIEDDATE, OPERATORID (lebar = 926); JSONDATA dibuang (928); tabel flat 926
-- dibuang (928); nol tabel baru.
--
-- Mengapa dua langkah: pelari menjalankan pernyataan TANPA transaksi dan hanya menoleransi ORA-00955 pada CREATE.
-- `ALTER TABLE … ADD (` diulang sesudah gagal mati di ORA-01430, dan penjaga inti MENOLAK kolom baru lewat blok
-- berpelindung katalog di jalur maju (kolom harus terlihat `KolomAlterTambah`). Maka ADD berdiri sendiri di langkah
-- ini (satu pernyataan: berhasil seluruhnya atau tidak sama sekali, T_MIGRASI mencatatnya), dan seluruh pernyataan 928
-- aman diulang. Urutan WO + pemulihan: modul/riratelife/docs/LANGKAH-WO-RIRATELIFE-SATU-TABEL.md.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.M_RATE_LIFE_SUMMARY ADD (
  USEDBY        VARCHAR2(500),
  TYPE          VARCHAR2(100),
  MODIFIEDDATE  VARCHAR2(50),
  OPERATORID    VARCHAR2(200)
)
/

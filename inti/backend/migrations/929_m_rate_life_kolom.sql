-- 929 - M_RATE_LIFE (rincian R/I Rate Life) mendapat kolom bernama (langkah 1 dari 2; langkah 2 = 930).
--
-- Keputusan work owner 07-10-2026: "detail R/I Rate Life (M_RATE_LIFE) menjadi tabel flat, satu tabel saja (pola sama
-- dengan ringkasan, migrasi 927/928); M_RATE_LIFE TIDAK dihapus". Sasaran: M_RATE_LIFE berkolom PERSIS seperti view
-- RATE_LIFE (`SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER,
-- a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM M_RATE_LIFE a`): ID VARCHAR2(10) PK + tujuh kolom di bawah.
-- Tipe TETAP TEKS (VARCHAR2) seperti view - nilai Pega utuh (RATE berkoma maupun bertitik desimal), pembaca tidak berubah
-- perilaku. Lebar dari data DEV 07-10-2026 (panjang maksimum: IDUSEDBY 7, USEDBY 95, TYPE kosong, GENDER 1, CONTRACT 3,
-- AGE 3, RATE 20) dengan cadangan, sejalan ringkasan 927 (IDUSEDBY = ID ringkasan VARCHAR2(10), USEDBY 500, TYPE 100).
--
-- Mengapa dua langkah: `ALTER TABLE … ADD (` diulang sesudah gagal mati di ORA-01430, dan penjaga inti menolak kolom baru
-- lewat blok berpelindung di jalur maju - ADD berdiri sendiri (satu pernyataan, berhasil seluruhnya atau tidak), 930
-- aman diulang. Urutan WO: modul/riratelife/docs/LANGKAH-WO-RIRATELIFE-DETAIL-FLAT.md. NOL COMMIT (ADR-U-0029).
ALTER TABLE {skema}.M_RATE_LIFE ADD (
  IDUSEDBY  VARCHAR2(10),
  USEDBY    VARCHAR2(500),
  TYPE      VARCHAR2(100),
  GENDER    VARCHAR2(10),
  CONTRACT  VARCHAR2(10),
  AGE       VARCHAR2(10),
  RATE      VARCHAR2(50)
)
/

-- 926 - RATE_LIFE_SUMMARY: tabel FLAT ringkasan R/I Rate Life (modul `riratelife`), menggantikan VIEW warisan bernama
-- sama. Perintah work owner 06-10-2026: "tabel M_RATE_LIFE_SUMMARY buat jadi flat menampilkan data detail yang ada pada
-- tabel view RATE_LIFE_SUMMARY". K-F1: view diganti tabel flat bernama sama (pola 924 RICOMM_LIFE); M_RATE_LIFE_SUMMARY
-- (JSON, data DEV masih berubah) TIDAK disentuh - cadangan dan sumber alat `modul/riratelife/backend/alat/pindahflat`. K-F2:
-- kolom view ikut (`SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID,
-- a.JSONDATA.FLAG FROM M_RATE_LIFE_SUMMARY a`, ALL_VIEWS dibaca WO 06-10-2026), semua teks, isi apa adanya - KECUALI
-- FLAG: keputusan work owner 06-10-2026 "kolom FLAG tidak digunakan" (nilai lama tetap di M_RATE_LIFE_SUMMARY.JSONDATA). Pembaca lain (`mastercontractretrolife`, `masterproductnamelife`: `SELECT ID, USEDBY`) ikut membaca tabel.
--
-- ⛔ PRASYARAT: VIEW POOLDATA.RATE_LIFE_SUMMARY sudah dibuang DBA lewat
-- `modul/riratelife/docs/DBA-LEPAS-VIEW-RATE_LIFE_SUMMARY.sql` SEBELUM `-migrate`. Selama view itu ada, pra-terbang
-- pelari menemukan objek bernama sama dengan kolom yang sama, CREATE TABLE dijawab ORA-00955 dan DILEWATI; CREATE INDEX
-- di bawah lalu gagal atas view (ORA-01702), -migrate BERHENTI KERAS dan 926 TIDAK tercatat. Langkah WO + pemulihan:
-- modul/riratelife/docs/LANGKAH-WO-RIRATELIFE-FLAT.md.
--
-- Lebar (bukti per kolom: modul/riratelife/docs/STRUKTUR-TABEL-RIRATELIFE.md): ID VARCHAR2(10) (RATE_LIFE_SUMMARY.ID,
-- STRUKTUR MPNL b180); USEDBY VARCHAR2(500) (preseden MPNL M_PRODUCTNAME_LIFE_PLAN.RIRATE yang menyalin USEDBY);
-- OPERATORID VARCHAR2(200) (Go memotong akun ke 200 byte); MODIFIEDDATE VARCHAR2(50), TYPE VARCHAR2(100). Lebar cukup
-- untuk data DEV (panjang maksimum isi dicek WO 06-10-2026: ID 7, USEDBY 99, MODIFIEDDATE 23, OPERATORID 17, TYPE kosong).
-- Alat pindah MENOLAK nilai yang tidak muat - nol pemotongan. Semua kolom NULLABLE kecuali PK; NOL COMMIT.
CREATE TABLE {skema}.RATE_LIFE_SUMMARY (
  ID            VARCHAR2(10) NOT NULL,
  USEDBY        VARCHAR2(500),
  TYPE          VARCHAR2(100),
  MODIFIEDDATE  VARCHAR2(50),
  OPERATORID    VARCHAR2(200),
  CONSTRAINT PK_RATE_LIFE_SUMMARY PRIMARY KEY (ID)
)
/
-- Indeks nama: dipakai SETIAP Save / Simpan Upload (`SqlPemakaiNama`: `UPPER(TRIM(USEDBY)) = UPPER(TRIM(:1))`, nama
-- tidak kembar dan pencocokan upload). Sekaligus PENGAMAN: atas view yang belum dibuang, pernyataan ini gagal ORA-01702.
CREATE INDEX {skema}.IX_RATE_LIFE_SUMMARY_NAMA ON {skema}.RATE_LIFE_SUMMARY (UPPER(TRIM(USEDBY)))
/

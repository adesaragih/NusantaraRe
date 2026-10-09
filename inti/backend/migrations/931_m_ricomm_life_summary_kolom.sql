-- 931 - M_RICOMM_LIFE_SUMMARY mendapat kolom ringkasan (langkah 1 dari 2; langkah 2 = 932).
--
-- Keputusan work owner 08-10-2026: "R/I Comm Life dibuat satu tabel per jenis data, sama seperti R/I Rate Life
-- (927/928, 929/930)". Sasaran akhir: M_RICOMM_LIFE_SUMMARY berkolom persis seperti view RICOMM_LIFE_SUMMARY
-- (ID PK warisan + USEDBY, MODIFIEDDATE, OPERATORID); JSONDATA dan view dibuang (932); nol tabel baru.
--
-- Lebar [terverifikasi data DEV 08-10-2026]: USEDBY VARCHAR2(200) - panjang maksimum isi 18 byte; SAMA dengan salinan
-- nama di rincian (RICOMM_LIFE.USEDBY 924 / M_RICOMM_LIFE.USEDBY 933) dan batas Go models.BatasNama 200, sehingga nama
-- ringkasan selalu muat di rinciannya (bukan 500 milik R/I Rate: modul ini menyalin nama ke rincian 200 byte).
-- MODIFIEDDATE VARCHAR2(50) (bentuk Pega YYYYMMDDTHHMMSS.mmm GMT, 23 byte) dan OPERATORID VARCHAR2(200) = 927.
--
-- Mengapa berdiri sendiri: pelari tanpa transaksi; ALTER TABLE ... ADD ( diulang mati di ORA-01430 dan penjaga inti
-- menolak kolom baru lewat blok di jalur maju. Satu pernyataan: berhasil seluruhnya atau tidak sama sekali.
-- Urutan WO + pemulihan: modul/ricommlife/docs/LANGKAH-WO-RICOMMLIFE-SATU-TABEL.md.
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.M_RICOMM_LIFE_SUMMARY ADD (
  USEDBY        VARCHAR2(200),
  MODIFIEDDATE  VARCHAR2(50),
  OPERATORID    VARCHAR2(200)
)
/

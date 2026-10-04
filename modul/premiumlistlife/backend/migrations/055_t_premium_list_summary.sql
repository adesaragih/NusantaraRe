-- T_PREMIUM_LIST_SUMMARY - rekap uang per mata uang (tiket 00).
--
-- ⛔ REKAP UANG PENUH, bukan hanya kode mata uang (AC 37, penyimpangan
-- sadar 2). Menyimpan kodenya saja memaksa setiap pembaca menjumlah ulang
-- dari peserta - dan jumlah yang dihitung ulang di banyak tempat adalah
-- jumlah yang suatu hari berbeda di salah satunya.
--
-- ⚠️ Ia juga kontrak HILIR: Claim Life membaca `M_LIFE_PREMIUM_SUMMARY`
-- warisan. Penulisan tabel warisan itu ditiru di Go DALAM TRANSAKSI YANG SAMA
-- dengan tabel relasional ini (butir pl2) - titik potong transaksi lenyap.
--
-- ⛔ TIPE FISIKNYA KEPUTUSAN KAMI, BUKAN BACAAN DARI KORPUS.
-- STRUKTUR menyebut KATEGORI logis (teks / angka desimal / bilangan bulat /
-- DATE) dan menyatakan presisi fisik `[data DBA]`. Pemetaannya:
--
--   angka desimal   -> NUMBER(38,8)   uang dan share (ADR-0003, revisi-penyimpanan)
--   bilangan bulat  -> NUMBER(5)     umur, periode, PROD_KE, nomor baris
--                     (konvensi yang SUDAH ada di repo ini; penjaga
--                      TestNolNumberTanpaPresisi hanya menerima
--                      NUMBER(19), NUMBER(38,8), dan NUMBER(5))
--   DATE            -> DATE
--   teks            -> VARCHAR2(255); ID dan kolom ber-akhiran _ID -> VARCHAR2(32)
--
-- ⚠️ Ketiganya menunggu pencocokan DBA, dan kedua arah salahnya TIDAK
-- setara: VARCHAR2 yang terlalu pendek MENOLAK data yang sah - kegagalan yang
-- terlihat - sedangkan NUMBER yang terlalu pendek MEMBULATKAN uang diam-diam.
--
-- ⛔ SELURUH kolom nullable kecuali PK. STRUKTUR menyatakannya, dan
-- wajib-isi ditegakkan di Go. NOT NULL yang ditambahkan di sini akan menolak
-- baris warisan yang memang kosong saat migrasi data (tiket 09).
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). Batas transaksi milik Go.
--
-- ⚠️ Nama kolom di bawah DIBANGKITKAN dari STRUKTUR, tidak diketik ulang.
CREATE TABLE {skema}.T_PREMIUM_LIST_SUMMARY (
  ID                            VARCHAR2(32) NOT NULL,
  PREMIUM_LIST_ID               VARCHAR2(32),
  CURRENCY                      VARCHAR2(255),
  BALANCE                       NUMBER(38,8),
  PREMIUM                       NUMBER(38,8),
  COMMISSION                    NUMBER(38,8),
  PROF_COMM                     NUMBER(38,8),
  OVR_COMM                      NUMBER(38,8),
  TAX                           NUMBER(38,8),
  CLAIM                         NUMBER(38,8),
  CLAIM_AMOUNT                  NUMBER(38,8),
  DEDUCTION                     NUMBER(38,8),
  BROKERAGE_FEE                 NUMBER(38,8),
  RI_ADMIN_FEE                  NUMBER(38,8),
  CEDING_RETENTION              NUMBER(38,8),
  SUM_REASURED                  NUMBER(38,8),
  SUM_AT_RISK_GROSS             NUMBER(38,8),
  SHARE_RETRO                   NUMBER(38,8),
  SHARE_NUSANTARA_RE_GROSS      NUMBER(38,8),
  GROSS_PREMIUM_RETRO           NUMBER(38,8),
  NET_PREMIUM_RETRO             NUMBER(38,8),
  DISCOUNT_PREMIUM_RETRO        NUMBER(38,8),
  OVR_COMM_RETRO                NUMBER(38,8),
  BROKERAGE_FEE_RETRO           NUMBER(38,8),
  RI_ADMIN_FEE_RETRO            NUMBER(38,8),
  GROSS_PREMIUM_REFUND          NUMBER(38,8),
  NET_PREMIUM_REFUND            NUMBER(38,8),
  COMM_REFUND                   NUMBER(38,8),
  OVR_COMM_REFUND               NUMBER(38,8),
  TAX_REFUND                    NUMBER(38,8),
  BROKERAGE_FEE_REFUND          NUMBER(38,8),
  DEDUCTION_REFUND              NUMBER(38,8),
  RI_ADMIN_FEE_REFUND           NUMBER(38,8),
  GROSS_PREMIUM_REFUND_RETRO    NUMBER(38,8),
  NET_PREMIUM_REFUND_RETRO      NUMBER(38,8),
  DISCOUNT_PREMIUM_REFUND_RETRO NUMBER(38,8),
  OVR_COMM_REFUND_RETRO         NUMBER(38,8),
  BROKERAGE_FEE_REFUND_RETRO    NUMBER(38,8),
  RI_ADMIN_FEE_REFUND_RETRO     NUMBER(38,8),
  CONSTRAINT PK_T_PL_SUMMARY PRIMARY KEY (ID),
  CONSTRAINT FK_PLSUM_PL FOREIGN KEY (PREMIUM_LIST_ID)
    REFERENCES {skema}.T_PREMIUM_LIST (ID) ON DELETE CASCADE
)
/

CREATE INDEX {skema}.IDX_PLSUM_PL ON {skema}.T_PREMIUM_LIST_SUMMARY (PREMIUM_LIST_ID)
/

-- T_PREMIUM_LIST - header polis (tiket 00 PremiumList Life).
--
-- Satu baris = SATU VERSI polis, new business maupun endorsement. Seluruh
-- versi hidup berdampingan; versi berjalan adalah baris ber-`PROD_KE`
-- TERBESAR (Endorsement Life spec §16, AC 56).
--
-- ⛔ NB DAN EDM MEMAKAI TABEL YANG SAMA. Yang membedakan BARIS, bukan
-- tabel. Kolom EDM (`EDM_TYPE`, `OLD_POLICY_NO`, dan seterusnya) seluruhnya
-- nullable dan KOSONG pada baris new business (AC 57).
--
-- ⛔ NOL KOLOM JSON (AC 32, AC 34, penyimpangan sadar 1). Setiap atribut
-- polis adalah KOLOM BERNAMA. `INSERTJSONPOLISLIFE` tidak ditiru sama sekali
-- dan `@ASM.GetPageJSONString()` tidak direplikasi: dokumen JSON yang
-- disimpan adalah dokumen yang tidak dapat dicari, tidak dapat divalidasi,
-- dan tidak dapat dibandingkan saat rekonsiliasi migrasi.
--
-- ⛔ NOL PROPERTI BAWAAN PEGA (AC 47): tidak ada `px*`, `py*`, `pz*`, dan
-- tidak ada single-page kosong (`Policy`, `Quotation`, `TempError`).
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
CREATE TABLE {skema}.T_PREMIUM_LIST (
  ID                    VARCHAR2(32) NOT NULL,
  ID_PEGA               VARCHAR2(255),
  NO_POLIS              VARCHAR2(255),
  BUSINESS_CODE         VARCHAR2(255),
  BUSINESS_NAME         VARCHAR2(255),
  CEDING_CO             VARCHAR2(255),
  CEDING_CO_NAME        VARCHAR2(255),
  DATE_RECEIVED         DATE,
  MARKETING_CODE        VARCHAR2(255),
  MARKETING_NAME        VARCHAR2(255),
  POLICY_HOLDER         VARCHAR2(255),
  POLICY_HOLDER_NAME    VARCHAR2(255),
  PRO_RATE_TYPE         VARCHAR2(255),
  CREATE_OP_NAME        VARCHAR2(255),
  SOB                   VARCHAR2(255),
  SOB_NAME              VARCHAR2(255),
  TYPE                  VARCHAR2(255),
  TYPE_CEDING           VARCHAR2(255),
  TYPE_CEDING_NAME      VARCHAR2(255),
  MO_ID                 VARCHAR2(32),
  NO_OFFER              VARCHAR2(255),
  RI_SLIP_RNM           VARCHAR2(255),
  RETRO_ID              VARCHAR2(32),
  RETRO_NAME            VARCHAR2(255),
  SECURITY_REINSURER_ID VARCHAR2(32),
  SECURITY_REINSURER    VARCHAR2(255),
  TGL_INPUT             DATE,
  DESCRIPTION           VARCHAR2(255),
  WPC                   DATE,
  PRODUCT_NAME_ID       VARCHAR2(32),
  PRODUCT_NAME          VARCHAR2(255),
  LAYER_1               VARCHAR2(255),
  LAYER_2               VARCHAR2(255),
  LAYER_3               VARCHAR2(255),
  LAYER_4               VARCHAR2(255),
  ANNUITY_INTEREST      NUMBER(38,8),
  PREMIUM_REFUND_FACTOR NUMBER(38,8),
  NO_ENDORS             VARCHAR2(255),
  EDM_TYPE              VARCHAR2(255),
  OLD_POLICY_NO         VARCHAR2(255),
  EDM_DATE              DATE,
  EDM_NOTE              VARCHAR2(255),
  PREMI_PROPOSED        NUMBER(38,8),
  UANG_PERTANGGUNGAN    NUMBER(38,8),
  SUM_INSURED           NUMBER(38,8),
  JENIS_PRODUK          VARCHAR2(255),
  SISTEM_REASURANSI     VARCHAR2(255),
  STATUSS               VARCHAR2(255),
  STATUS_UPDATE         VARCHAR2(255),
  STATUS_SERVICE        VARCHAR2(255),
  START_DATE            DATE,
  END_DATE              DATE,
  PL_NUMBER_EDM         VARCHAR2(255),
  PROD_KE               NUMBER(5),
  EDM_STATUS            VARCHAR2(255),
  STATUS_OLD            VARCHAR2(255),
  CONSTRAINT PK_T_PREMIUM_LIST PRIMARY KEY (ID)
)
/

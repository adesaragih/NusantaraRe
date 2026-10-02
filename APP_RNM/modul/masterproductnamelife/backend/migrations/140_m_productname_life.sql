-- 140 - M_PRODUCTNAME_LIFE: induk tabel flat produk Life (tiket 01 bab 02-10-2026, K2/K5).
--
-- Satu baris = satu produk: sisi umum (halaman Pega `ProductName`) dan sisi inward (`ProductNameInward`) DIGABUNG
-- (grilling Q1b) - dulu dua dokumen di dua tabel warisan, kini kolom bernama (grilling D2). Kedua tabel warisan
-- TIDAK disentuh migrasi ini: tetap ada sebagai cadangan dan sumber alat pindah `backend/alat/pindahflat`.
--
-- Tipe (K6, penjaga inti TestNolNumberTanpaPresisi): uang, persen, rate, faktor -> NUMBER(38,8); bilangan kecil
-- (usia, nomor, hari, kontrak) -> NUMBER(5); tanggal -> DATE; teks -> VARCHAR2 selebar data DEV 02-10-2026 dan ID
-- master rujukannya. Nilai yang tidak muat DITOLAK services berkalimat, tidak pernah dipotong atau dibulatkan.
--
-- Seluruh kolom NULLABLE kecuali PK (tiket 01 AC 43): wajib-isi ditegakkan di Go. Nol sequence baru: ID dari
-- M_PRODUCT_LIFE_SEQ warisan (P6). Batas transaksi milik Go (ADR-U-0029). Kolom dan isinya: docs/STRUKTUR-TABEL-...md.
CREATE TABLE {skema}.M_PRODUCTNAME_LIFE (
  ID                  VARCHAR2(6) NOT NULL,
  PRODUCTNAME         VARCHAR2(1000),
  CEDINGID            VARCHAR2(100),
  CEDING              VARCHAR2(1000),
  SOBID               VARCHAR2(100),
  SOBNAME             VARCHAR2(1000),
  RICOMM              NUMBER(38,8),
  RIRISKID            VARCHAR2(10),
  RIRISK              VARCHAR2(100),
  RIRATEID            VARCHAR2(10),
  RIRATE              VARCHAR2(500),
  INWARDNAME          VARCHAR2(1000),
  TREATYNUMBER        VARCHAR2(100),
  CAUSEID             VARCHAR2(10),
  CAUSE               VARCHAR2(100),
  IS_ORS              NUMBER(5),
  POLICYHOLDER        VARCHAR2(100),
  POLICYHOLDERNAME    VARCHAR2(1000),
  INSURED             VARCHAR2(4000),
  ADDENDUMNO          NUMBER(5),
  ADDENDUMWORD        VARCHAR2(4000),
  AMANDEMENTNO        NUMBER(5),
  AMANDEMENTSCHD      VARCHAR2(4000),
  MAXEXPIREDCLAIM     NUMBER(5),
  BEGIN_DATE          DATE,
  STNC                DATE,
  MATURE              DATE,
  CEDINGRETENTIONNUM  NUMBER(38,8),
  CEDINGLIMIT         NUMBER(38,8),
  BROKERAGE           NUMBER(38,8),
  MINAGE              NUMBER(5),
  MAXAGE              NUMBER(5),
  EXPIRYAGE           NUMBER(5),
  EXTRAPREMI          NUMBER(38,8),
  MINSUMINSURED       NUMBER(38,8),
  MAXSUMINSURED       NUMBER(38,8),
  MAXSUMREASURED      NUMBER(38,8),
  RNMSHARE            NUMBER(38,8),
  RNMLIMITNUM         NUMBER(38,8),
  PREMIUMFACTOR       NUMBER(38,8),
  PAYMENT             VARCHAR2(1),
  SUBJECTTO           VARCHAR2(4000),
  ANNUITYINTEREST     NUMBER(38,8),
  PREMIUMREFUNDFACTOR NUMBER(38,8),
  MAXDATARECEIVE      NUMBER(5),
  BIRTHDAY            VARCHAR2(1),
  CURRENCYID          VARCHAR2(100),
  CURRENCY            VARCHAR2(20),
  EXTRAMORTALITY      NUMBER(38,8),
  MAXCONTRACT         NUMBER(5),
  PROPORTIONALTABLE   NUMBER(38,8),
  CREATEOP            VARCHAR2(100),
  UPDATEOP            VARCHAR2(100),
  CONSTRAINT PK_M_PRODUCTNAME_LIFE PRIMARY KEY (ID),
  CONSTRAINT CK_M_PRODUCTNAME_LIFE_IS_ORS CHECK (IS_ORS IN (0, 1))
)
/

-- T_PREMIUM_LIST_DETAIL - peserta polis (tiket 00 PremiumList Life).
--
-- Satu baris = satu peserta pada satu versi polis. Sumber kolomnya
-- `RDBList/SaveMasterLPDet.xml` (80 kolom, identik byte demi byte di modul
-- PremiumList Life dan Endorsement Life).
--
-- ⚠️ SKALA JUTAAN BARIS. Satu polis grup dapat memuat ribuan peserta.
-- Karena itu SETIAP FK ber-index (AC 49), dan migrasi data dijalankan
-- BERTAHAP - bukan satu transaksi raksasa (tiket 09).
--
-- ⛔ `PARENT_ID` adalah FK SELF-REFERENCE ke peserta versi SEBELUMNYA,
-- nullable, ber-index. Ia dipakai konteks Endorsement Life untuk mencocokkan
-- peserta lama dengan peserta baru (Endorsement spec §16, AC 60). Ketiga
-- tabel anak lainnya TIDAK memilikinya - menambahkannya "supaya seragam"
-- akan membuat tiga kolom yang tidak pernah terisi.
--
-- ⛔ KEDUA JENDELA VALUASI ADA (AC 38): `GROSS_VALUATION_*` DAN
-- `RETRO_VALUATION_*`, beserta `EFFECTIVE_DATE`, `LAPSE_DATE`, `PERIOD_MM`.
-- Menyimpan satu jendela saja membuat klaim yang jatuh di jendela lainnya
-- tidak dapat divalidasi sama sekali - dan itulah gerbang yang dipakai modul
-- Claim Life.
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
CREATE TABLE {skema}.T_PREMIUM_LIST_DETAIL (
  ID                            VARCHAR2(32) NOT NULL,
  PREMIUM_LIST_ID               VARCHAR2(32),
  PARENT_ID                     VARCHAR2(32),
  ID_PEGA                       VARCHAR2(255),
  PL_NUMBER                     VARCHAR2(255),
  POLICY_NO                     VARCHAR2(255),
  POLICY_HOLDER                 VARCHAR2(255),
  CERTIFICATE_NO                VARCHAR2(255),
  NAME_OF_INSURED               VARCHAR2(255),
  DESCRIPTION                   VARCHAR2(255),
  SEX                           VARCHAR2(255),
  DOB                           DATE,
  AGE                           NUMBER(5),
  ENTRY_AGE                     NUMBER(5),
  CURRENT_AGE                   NUMBER(5),
  PLAN                          VARCHAR2(255),
  RISK                          VARCHAR2(255),
  MEDICAL_STATUS                VARCHAR2(255),
  STNC                          VARCHAR2(255),
  WPC                           VARCHAR2(255),
  CURRENCY                      VARCHAR2(255),
  PERIOD_YY                     NUMBER(5),
  PERIOD_MM                     NUMBER(5),
  PASSED_PERIOD                 NUMBER(5),
  BEGIN_DATE                    DATE,
  EFFECTIVE_DATE                DATE,
  EXPIRED_DATE                  DATE,
  LAPSE_DATE                    DATE,
  GROSS_VALUATION_BEGIN_DATE    DATE,
  GROSS_VALUATION_EXPIRED_DATE  DATE,
  RETRO_VALUATION_BEGIN_DATE    DATE,
  RETRO_VALUATION_EXPIRED_DATE  DATE,
  SUM_INSURED                   NUMBER(38,8),
  SUM_REASURED                  NUMBER(38,8),
  SUM_AT_RISK_GROSS             NUMBER(38,8),
  SUM_AT_RISK_RETRO             NUMBER(38,8),
  CEDING_RETENTION              NUMBER(38,8),
  CEDING_CO                     VARCHAR2(255),
  SHARE_NUSANTARA_RE            NUMBER(38,8),
  SHARE_NUSANTARA_RE_GROSS      NUMBER(38,8),
  SHARE_RETRO                   NUMBER(38,8),
  RETROCEDED_SHARE              NUMBER(38,8),
  RATE                          NUMBER(38,8),
  FACTOR                        NUMBER(38,8),
  EM_PERCENT                    NUMBER(38,8),
  PRO_RATE_TYPE                 VARCHAR2(255),
  GROSS_PREMIUM                 NUMBER(38,8),
  NET_PREMIUM                   NUMBER(38,8),
  COMM                          NUMBER(38,8),
  PROF_COMM                     NUMBER(38,8),
  OVR_COMM                      NUMBER(38,8),
  BROKERAGE_FEE                 NUMBER(38,8),
  TAX                           NUMBER(38,8),
  FLEET_DISCOUNT                NUMBER(38,8),
  DEDUCTION                     NUMBER(38,8),
  RI_ADMIN_FEE                  NUMBER(38,8),
  CLAIM                         NUMBER(38,8),
  CLAIM_AMOUNT                  NUMBER(38,8),
  GROSS_PREMIUM_REFUND          NUMBER(38,8),
  NET_PREMIUM_REFUND            NUMBER(38,8),
  COMM_REFUND                   NUMBER(38,8),
  OVR_COMM_REFUND               NUMBER(38,8),
  BROKERAGE_FEE_REFUND          NUMBER(38,8),
  TAX_REFUND                    NUMBER(38,8),
  DEDUCTION_REFUND              NUMBER(38,8),
  RI_ADMIN_FEE_REFUND           NUMBER(38,8),
  GROSS_PREMIUM_RETRO           NUMBER(38,8),
  NET_PREMIUM_RETRO             NUMBER(38,8),
  DISCOUNT_PREMIUM_RETRO        NUMBER(38,8),
  OVR_COMM_RETRO                NUMBER(38,8),
  BROKERAGE_FEE_RETRO           NUMBER(38,8),
  RI_ADMIN_FEE_RETRO            NUMBER(38,8),
  GROSS_PREMIUM_REFUND_RETRO    NUMBER(38,8),
  NET_PREMIUM_REFUND_RETRO      NUMBER(38,8),
  DISCOUNT_PREMIUM_REFUND_RETRO NUMBER(38,8),
  OVR_COMM_REFUND_RETRO         NUMBER(38,8),
  BROKERAGE_FEE_REFUND_RETRO    NUMBER(38,8),
  RI_ADMIN_FEE_REFUND_RETRO     NUMBER(38,8),
  PL_NUMBER_EDM                 VARCHAR2(255),
  EDM_STATUS                    VARCHAR2(255),
  STATUS_OLD                    VARCHAR2(255),
  STATUS                        VARCHAR2(255),
  CONSTRAINT PK_T_PL_DETAIL PRIMARY KEY (ID),
  CONSTRAINT FK_PLD_PL FOREIGN KEY (PREMIUM_LIST_ID)
    REFERENCES {skema}.T_PREMIUM_LIST (ID) ON DELETE CASCADE,
  CONSTRAINT FK_PLD_PARENT FOREIGN KEY (PARENT_ID)
    REFERENCES {skema}.T_PREMIUM_LIST_DETAIL (ID) ON DELETE CASCADE
)
/

CREATE INDEX {skema}.IDX_PLD_PL ON {skema}.T_PREMIUM_LIST_DETAIL (PREMIUM_LIST_ID)
/
CREATE INDEX {skema}.IDX_PLD_PARENT ON {skema}.T_PREMIUM_LIST_DETAIL (PARENT_ID)
/

-- T_CLAIMLF_PREMIUMLIST_DETAIL - peserta yang diklaim
--
-- Migrasi tiket 14 Claim Life. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Aturan yang dijaga seluruh berkas di folder ini:
--   ADR-U-0027  seluruh kolom nullable kecuali kunci utama; wajib-isi di services
--   ADR-U-0003  ADR-U-0016  uang NUMBER(38,8), tidak pernah float
--   ADR-U-0029  nol COMMIT di teks SQL; transaksi dibuka-ditutup aplikasi
--   ADR-U-0006  identitas dari sequence, kecuali yang dinyatakan berformat
--
-- Satu baris mewakili satu peserta di dalam satu klaim.
-- Induknya T_GENERAL_CLAIM lewat CLAIM_ID, 1:N, ON DELETE CASCADE (relasi 3).
-- Cascade berarti menghapus induk ikut menghapus anaknya, dikerjakan Oracle.
--
-- Atribut polis TIDAK disalin ke sini. Yang disimpan hanya nilai yang dipakai
-- dan perlu bertahan. Sembilan kolom tambahan hasil audit ikut: IS_CHECK,
-- STATUS, RECOMMENDATION, STS_REJECT, SOURCE_ID, ketiga tanggal per peserta,
-- dan CEDING_RETENTION.
--
-- Kolom share, persen, dan rate memakai NUMBER(38,8) - keputusan work owner
-- 26 September 2026 (brief ronde 2 bab 2c). Presisinya sama dengan kolom uang,
-- tetapi di Go tipenya Ratio, bukan Money: ia perbandingan, bukan mata uang.
--
-- Dua kolom valuasi retro memakai nama korpus RETRO_VALUATION_* dan bukan nama
-- STRUKTUR RETROCESSION_VALUATION_* - keputusan work owner 26 September 2026
-- (brief ronde 2 bab 2j). Nama STRUKTUR berukuran 33 dan 35 byte, dan Oracle di
-- bawah 12.2 menolak pengenal lebih dari 30 byte dengan ORA-00972. Selisih nama
-- ini dijaga test TestKolomDDLCocokDenganStruktur.
CREATE TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL (
  ID                    VARCHAR2(32) NOT NULL,
  CLAIM_ID              VARCHAR2(32),
  PL_NUMBER             VARCHAR2(64),
  POLICY_NO             VARCHAR2(64),
  POLICY_HOLDER         VARCHAR2(255),
  CERTIFICATE_NO        VARCHAR2(64),
  NAME_OF_INSURED       VARCHAR2(255),
  DOB                   DATE,
  AGE                   NUMBER(5),
  SEX                   VARCHAR2(8),
  PLAN                  VARCHAR2(128),
  DISEASE               VARCHAR2(255),
  ICD_CODE              VARCHAR2(32),
  DESCRIPTION           VARCHAR2(1000),
  NOTES                 VARCHAR2(1000),
  KETERANGAN            VARCHAR2(1000),
  IS_CHECK              VARCHAR2(8),
  STATUS                VARCHAR2(32),
  RECOMMENDATION        VARCHAR2(1000),
  STS_REJECT            VARCHAR2(8),
  SOURCE_ID             VARCHAR2(64),
  DATE_OF_LOSS          DATE,
  RECEIVED_DATE         DATE,
  CONFIRMATION_DATE     DATE,
  COMPLETE_DATE         DATE,
  CLAIM_RECEIVED_DATE   DATE,
  BEGIN_DATE            DATE,
  EFFECTIVE_DATE        DATE,
  EXPIRED_DATE          DATE,
  LAPSE_DATE            DATE,
  GROSS_VALUATION_BEGIN_DATE    DATE,
  GROSS_VALUATION_EXPIRED_DATE  DATE,
  RETRO_VALUATION_BEGIN_DATE    DATE,
  RETRO_VALUATION_EXPIRED_DATE  DATE,
  CURRENCY              VARCHAR2(8),
  SUM_INSURED           NUMBER(38,8),
  SUM_REASURED          NUMBER(38,8),
  GROSS_PREMIUM         NUMBER(38,8),
  NET_PREMIUM           NUMBER(38,8),
  CLAIM_AMOUNT          NUMBER(38,8),
  EM_PERCENT            NUMBER(38,8),
  SHARE_NUSANTARA_RE    NUMBER(38,8),
  SHARE_RETRO           NUMBER(38,8),
  RETROCEDED_SHARE      NUMBER(38,8),
  CEDING_RETENTION      NUMBER(38,8),
  WPC                   VARCHAR2(32),
  STNC_TREATY           VARCHAR2(64),
  CONSTRAINT PK_T_CLAIMLF_PLD PRIMARY KEY (ID),
  CONSTRAINT FK_PLD_CLAIM FOREIGN KEY (CLAIM_ID)
    REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_PLD_CLAIM_ID ON {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL (CLAIM_ID)
/

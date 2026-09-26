-- T_GENERAL_CLAIM - header klaim
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
-- Header klaim Life. Satu baris mewakili satu klaim.
--
-- SHARED PRIMARY KEY: ID di sini SAMA PERSIS dengan T_WORK_CLAIM.ID baris
-- klaim. Shared primary key berarti kunci utama tabel ini sekaligus menjadi
-- kunci tamu ke induknya, sehingga tidak ada kolom penyambung terpisah.
-- Kolom WORK_CLAIM_ID TIDAK ADA: ia dibuang, bukan diganti nama.
--
-- Ketiga penunjuk polis menunjuk KE LUAR, ke modul PremiumList Life. Mereka
-- bukan anak dan tidak ikut kaskade.
-- [terbuka] Isi ketiga penunjuk belum dapat ditulis: tidak satu pun menunjuk
-- kolom yang ada di sisi polis. Kolomnya dibuat; pengisiannya menunggu
-- keputusan work owner. Jangan tebak.
--
-- CASEID, CREATE_OP, CREATE_OP_NAME, dan TGL_UPDATE TIDAK ADA di sini:
-- keempatnya pindah ke T_WORK_CLAIM. PL_NUMBER diganti nama jadi POLICY_NO.
CREATE TABLE {skema}.T_GENERAL_CLAIM (
  ID               VARCHAR2(32) NOT NULL,
  CLAIM_NO         VARCHAR2(64),
  PY_ID            VARCHAR2(64),
  STS_KATASTROFE   VARCHAR2(8),
  KATASTROFE_NOTE  VARCHAR2(1000),
  IS_KPR           VARCHAR2(8),
  STNC_CLAIM       VARCHAR2(64),
  ACCEPTED_NO      VARCHAR2(64),
  STS_REJECT       VARCHAR2(8),
  RI_SLIP_RNM      VARCHAR2(64),
  BUSINESS_NAME    VARCHAR2(255),
  CLAIM_RETRO      NUMBER(38,8),
  CASEID_POLICY    VARCHAR2(64),
  POLICY_NO        VARCHAR2(64),
  ENDORSMENT_NO    VARCHAR2(64),
  CONSTRAINT PK_T_GENERAL_CLAIM PRIMARY KEY (ID),
  CONSTRAINT FK_GENERAL_CLAIM_WORK FOREIGN KEY (ID)
    REFERENCES {skema}.T_WORK_CLAIM (ID)
)
/

-- T_CLAIMLF_ADJUSTMENT - baris keputusan
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
-- Satu baris mewakili satu baris AdjustmentList, dan itulah unit keputusan
-- mesin status (ADR-U-0011). Induknya PESERTA, bukan header klaim.
--
-- KOMITE_ID adalah penunjuk KE ATAS ke T_WORK_CLAIM.ID baris komite, bertipe
-- teks berformat KMT-xxxxxx, nullable, dan ber-index UNIK: satu baris
-- adjustment sama dengan tepat satu kasus komite. Oracle mengizinkan banyak
-- NULL pada index unik, jadi baris yang belum pernah dikirim ke Komite tetap
-- boleh banyak.
-- [terbuka] Apakah KOMITE_ID dipasangi REFERENCES belum diputuskan (tiket 14
-- bab Blocker), sehingga constraint-nya TIDAK dipasang di sini.
CREATE TABLE {skema}.T_CLAIMLF_ADJUSTMENT (
  ID                      VARCHAR2(32) NOT NULL,
  PREMIUM_LIST_DETAIL_ID  VARCHAR2(32),
  SHARE_NUSANTARA_RE      NUMBER,
  CEDING_RETENTION        NUMBER,
  SUM_INSURED             NUMBER(38,8),
  SUM_REASURED            NUMBER(38,8),
  SHARE_RETRO             NUMBER,
  RETROCEDED_SHARE        NUMBER,
  CURRENCY_ID             VARCHAR2(32),
  CURRENCY                VARCHAR2(8),
  CLAIM_AMOUNT            NUMBER(38,8),
  STS_REJECT              VARCHAR2(8),
  ACCEPTED_NO             VARCHAR2(64),
  ACCEPTATION_DATE        DATE,
  NAME_OF_BANK            VARCHAR2(255),
  ID_BANK                 VARCHAR2(64),
  ACCOUNT_NO              VARCHAR2(64),
  KOMITE_ID               VARCHAR2(32),
  CONSTRAINT PK_T_CLAIMLF_ADJ PRIMARY KEY (ID),
  CONSTRAINT FK_ADJ_PLD FOREIGN KEY (PREMIUM_LIST_DETAIL_ID)
    REFERENCES {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_ADJ_PLD_ID ON {skema}.T_CLAIMLF_ADJUSTMENT (PREMIUM_LIST_DETAIL_ID)
/
CREATE UNIQUE INDEX {skema}.UX_ADJ_KOMITE_ID ON {skema}.T_CLAIMLF_ADJUSTMENT (KOMITE_ID)
/

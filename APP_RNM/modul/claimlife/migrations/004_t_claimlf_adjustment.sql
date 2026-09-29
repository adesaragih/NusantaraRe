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
-- Kolom share, persen, dan rate bertipe NUMBER(38,8) - keputusan work owner c,
-- 26 September 2026. Presisinya sama dengan kolom uang, tetapi di Go tipenya
-- Ratio dan bukan Money: ia perbandingan, bukan mata uang, sehingga keduanya
-- tidak pernah terjumlahkan. NUMBER tanpa presisi TIDAK dipakai di mana pun -
-- dijaga TestNolNumberTanpaPresisi.
-- KETETAPAN MODUL CLAIM LIFE atas kolom persen dan rate, SESUAI ADR-U-0016
-- Akibat 2 - bukan penyimpangan darinya. Akibat 2 berbunyi "kolom persen tidak
-- termasuk - ia bukan uang, dan tetap mengikuti ketetapan modulnya", yaitu ia
-- MENYERAHKAN kolom itu ke modul masing-masing. Keputusan work owner c
-- (26-09-2026) adalah ketetapan modul ini, persis yang diserahkan kepadanya.
-- Karena itu nol amandemen ADR diperlukan.
-- `[keputusan work owner 26-09-2026, butir p1]` meralat label "penyimpangan
-- sadar" yang ditulis ronde 3 di tempat ini.
--
-- Satu baris mewakili satu baris AdjustmentList, dan itulah unit keputusan
-- mesin status (ADR-U-0011). Induknya PESERTA, bukan header klaim.
--
-- KOMITE_ID adalah penunjuk KE ATAS ke T_WORK_CLAIM.ID baris komite, bertipe
-- teks berformat KMT-xxxxxx, nullable, dan ber-index UNIK: satu baris
-- adjustment sama dengan tepat satu kasus komite. Oracle mengizinkan banyak
-- NULL pada index unik, jadi baris yang belum pernah dikirim ke Komite tetap
-- boleh banyak.
-- [keputusan work owner 26-09-2026, butir d] KOMITE_ID DIPASANGI REFERENCES ke
-- T_WORK_CLAIM(ID), tetap nullable. Baris yang belum pernah dikirim ke Komite
-- tetap NULL; yang sudah, wajib menunjuk baris komite yang benar-benar ada.
--
-- TANPA "ON DELETE", sebab alasan yang sama dengan COVER_KEY di 001: ini
-- penunjuk KE ATAS, bukan kepemilikan. Menghapus baris komite tidak boleh ikut
-- menghapus baris adjustment yang menunjuknya.
CREATE TABLE {skema}.T_CLAIMLF_ADJUSTMENT (
  ID                      VARCHAR2(32) NOT NULL,
  PREMIUM_LIST_DETAIL_ID  VARCHAR2(32),
  SHARE_NUSANTARA_RE      NUMBER(38,8),
  CEDING_RETENTION        NUMBER(38,8),
  SUM_INSURED             NUMBER(38,8),
  SUM_REASURED            NUMBER(38,8),
  SHARE_RETRO             NUMBER(38,8),
  RETROCEDED_SHARE        NUMBER(38,8),
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
    REFERENCES {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL (ID) ON DELETE CASCADE,
  CONSTRAINT FK_ADJ_KOMITE FOREIGN KEY (KOMITE_ID)
    REFERENCES {skema}.T_WORK_CLAIM (ID)
)
/
CREATE INDEX {skema}.IX_ADJ_PLD_ID ON {skema}.T_CLAIMLF_ADJUSTMENT (PREMIUM_LIST_DETAIL_ID)
/
CREATE UNIQUE INDEX {skema}.UX_ADJ_KOMITE_ID ON {skema}.T_CLAIMLF_ADJUSTMENT (KOMITE_ID)
/

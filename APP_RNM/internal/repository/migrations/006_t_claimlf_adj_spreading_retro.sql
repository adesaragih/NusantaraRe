-- T_CLAIMLF_ADJ_SPREADING_RETRO - pecahan per reinsurer
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
-- Pecahan spreading per reinsurer. Satu baris mewakili satu reinsurer pada
-- satu baris spreading. Ini tingkat keenam, cicit dari baris adjustment.
--
-- Ejaan COMMISION dipertahankan apa adanya dari dokumen struktur, termasuk
-- kekeliruan ejaannya, supaya tidak lahir dua nama untuk satu kolom.
--
-- DITUTUP OLEH XML 26-09-2026. Kedua cabang PREMIUM_SPREADED_NET memang ada,
-- dan yang memilihnya adalah TAHUN POLIS: langkah 8.2.1.9.3 ("Tahun pertama")
-- memakai GROSS - Discount - Comm, langkah 8.2.1.9.2 ("Bukan tahun pertama")
-- memakai GROSS - Comm. Lihat bab "Pembacaan ulang XML" tiket 03.
--
-- RALAT 26-09-2026: COMMISION dan OVR_COMM adalah PERSEN, bukan uang - XML
-- membagi keduanya seratus sebelum memakainya. Di Go keduanya kini Ratio.
-- Tipe kolomnya tetap NUMBER(38,8); nol langkah migrasi baru diperlukan.
CREATE TABLE {skema}.T_CLAIMLF_ADJ_SPREADING_RETRO (
  ID                      VARCHAR2(32) NOT NULL,
  SPREADING_ID            VARCHAR2(32),
  REINSURER_NAME          VARCHAR2(255),
  PERCENT_SHARE           NUMBER(38,8),
  AMOUNT                  NUMBER(38,8),
  RATE                    NUMBER(38,8),
  PREMIUM_SPREADED_GROSS  NUMBER(38,8),
  PREMIUM_SPREADED_NET    NUMBER(38,8),
  COMMISION               NUMBER(38,8),
  OVR_COMM                NUMBER(38,8),
  TREATY_TYPE_ID          VARCHAR2(32),
  TREATY_TYPE_NAME        VARCHAR2(128),
  CONSTRAINT PK_T_CLAIMLF_SPR_RETRO PRIMARY KEY (ID),
  CONSTRAINT FK_SPR_RETRO_SPR FOREIGN KEY (SPREADING_ID)
    REFERENCES {skema}.T_CLAIMLF_ADJUSTMENT_SPREADING (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_SPR_RETRO_SPR_ID
  ON {skema}.T_CLAIMLF_ADJ_SPREADING_RETRO (SPREADING_ID)
/

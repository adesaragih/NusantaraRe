-- T_CLAIMLF_ADJUSTMENT_SPREADING - spreading per treaty-year
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
-- Hasil spreading sebuah baris adjustment, dipecah per treaty-year. Satu baris
-- mewakili satu treaty-year dari satu baris adjustment.
--
-- Nilainya DIBEKUKAN saat adjustment disimpan: perubahan master treaty
-- sesudahnya tidak mengubah angka yang sudah tersimpan.
--
-- RETROCADED_SHARE dan RATE adalah share dan rate, bukan uang. Tipe pastinya
-- masih USULAN di brief sesi bab 2c, jadi dipakai NUMBER berpresisi arbitrer.
CREATE TABLE {skema}.T_CLAIMLF_ADJUSTMENT_SPREADING (
  ID                VARCHAR2(32) NOT NULL,
  ADJUSTMENT_ID     VARCHAR2(32),
  TREATY_TYPE_ID    VARCHAR2(32),
  TREATY_TYPE_NAME  VARCHAR2(128),
  TREATY_YEAR_LIFE  VARCHAR2(16),
  RETROCADED_SHARE  NUMBER,
  RATE              NUMBER,
  IDR               NUMBER(38,8),
  USD               NUMBER(38,8),
  CURRENCY          VARCHAR2(8),
  CONSTRAINT PK_T_CLAIMLF_SPR PRIMARY KEY (ID),
  CONSTRAINT FK_SPR_ADJ FOREIGN KEY (ADJUSTMENT_ID)
    REFERENCES {skema}.T_CLAIMLF_ADJUSTMENT (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_SPR_ADJ_ID ON {skema}.T_CLAIMLF_ADJUSTMENT_SPREADING (ADJUSTMENT_ID)
/

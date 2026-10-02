-- T_PREMIUM_LIST - kolom isian layar Input Offer (tiket 01 bagian 3).
--
-- `[terverifikasi]` sel `Section/InputOfferLife.xml` (salinan korpus
-- `kelvin\PremiumListLife (Done)\`, dibaca 01-10-2026), dan layar Pega lama
-- yang ditunjukkan work owner berisi nilai pada sel-sel ini:
--
--   .BatasUsiaPeserta      "Age Limit"          pxNumber   -> BATAS_USIA_PESERTA
--   .PeriodePertanggungan  "Coverage Period"    pxTextInput-> PERIODE_PERTANGGUNGAN
--   .TanggalPenawaran      "Offering Date"      pxDateTime -> TANGGAL_PENAWARAN
--   .TanggalRespon         "Response Date"      pxDateTime -> TANGGAL_RESPON
--   .TanggalKonfirmasi     "Confirmation Date"  pxDateTime -> TANGGAL_KONFIRMASI
--   .TBC                   "Input TBC"          pxNumber   -> TBC
--   .TanggalTBC            "Max TBC"            (SetMaxTBCLife_Act) -> TANGGAL_TBC
--   .KeteranganMarketing   "Marketing Note"     pxTextArea -> KETERANGAN_MARKETING
--
-- `Sum Insured` dan `Status Update` TIDAK ditambah: kolomnya sudah ada
-- (`SUM_INSURED`, `STATUS_UPDATE`, migrasi 051).
--
-- ⚠️ Spec user story 23 sudah meminta tanggal-tanggal penawaran tersimpan
-- (`INSERTJSONOFFERLIFE`: OFFERING, RESPONSE, CONFIRMATION, ..., MAX_TBC);
-- rancangan tujuh tabel (§12) melewatkannya. Migrasi ini menutup celah itu
-- untuk sel yang terisi di layar lama - bukan seluruh 24 kolom warisan.
--
-- Tipe mengikuti pemetaan 051: teks VARCHAR2(255), tanggal DATE, bilangan
-- bulat NUMBER(5). Nullable seluruhnya; nol `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.T_PREMIUM_LIST ADD (
  BATAS_USIA_PESERTA    NUMBER(5),
  PERIODE_PERTANGGUNGAN VARCHAR2(255),
  TANGGAL_PENAWARAN     DATE,
  TANGGAL_RESPON        DATE,
  TANGGAL_KONFIRMASI    DATE,
  TBC                   NUMBER(5),
  TANGGAL_TBC           DATE,
  KETERANGAN_MARKETING  VARCHAR2(255)
)
/

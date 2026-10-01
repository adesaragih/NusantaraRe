-- T_PREMIUM_LIST - sisa sel layar Input Offer (tiket 01 bagian 3, lanjutan 059).
--
-- `[keputusan work owner 01-10-2026]` "tampilkan semua kolom di gambar, nanti
-- baru saya filter": tujuh sel `Section/InputOfferLife.xml` yang KOSONG di
-- layar Pega lama yang ditunjukkan, sehingga tidak ikut di 059:
--
--   .QQName                  "Insured Name"         pxTextInput -> QQ_NAME
--   .JenisUsaha              "Occupation"           pxTextInput -> JENIS_USAHA
--   .KetentuanUnderwriting   "Underwriting Policy"  pxTextArea  -> KETENTUAN_UNDERWRITING
--   .TanggalKonfirmasiBalik  "Re-Confirmation Date" pxDateTime  -> TANGGAL_KONFIRMASI_BALIK
--   .TanggalRealisasi        "Realization Date"     pxDateTime  -> TANGGAL_REALISASI
--   .TanggalBind             "Binding Date"         pxDateTime  -> TANGGAL_BIND
--   .StatusFinal             "Final Status"         pxTextInput -> STATUS_FINAL
--
-- Migrasi TERPISAH dari 059, bukan suntingan di tempat: 059 mungkin sudah
-- terpasang, dan `T_MIGRASI` mencatat nama berkas yang sudah berjalan.
--
-- Tipe mengikuti pemetaan 051. Nullable seluruhnya; nol `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.T_PREMIUM_LIST ADD (
  QQ_NAME                  VARCHAR2(255),
  JENIS_USAHA              VARCHAR2(255),
  KETENTUAN_UNDERWRITING   VARCHAR2(255),
  TANGGAL_KONFIRMASI_BALIK DATE,
  TANGGAL_REALISASI        DATE,
  TANGGAL_BIND             DATE,
  STATUS_FINAL             VARCHAR2(255)
)
/

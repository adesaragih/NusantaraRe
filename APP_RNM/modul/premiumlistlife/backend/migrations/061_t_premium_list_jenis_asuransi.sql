-- T_PREMIUM_LIST.JENIS_ASURANSI - "Reinsurance Type" layar Input Offer.
--
-- `[terverifikasi]` sel `.JenisAsuransi` pxDropdown "Reinsurance Type"
-- (`Section/InputOfferLife.xml`); nilainya `Proportional` / `Non Proportional`
-- (`DataTransform/SetReinsuranceType`, InputOfferLife_ACT langkah 7).
--
-- Nilainya DITURUNKAN dari System Reinsurance (`TypeCeding "4"` -> Non
-- Proportional, selain itu Proportional) dan disimpan setiap simpan - padanan
-- `Obj-Save` InputOfferLife_ACT langkah 8. Tidak pernah diterima dari klien.
--
-- ⚠️ RALAT 01-10-2026 (komentar saja, SQL tidak berubah): berkas ini lahir
-- untuk dropdown yang DIPILIH pemakai; work owner membatalkannya di hari yang
-- sama. Kolomnya tetap dipakai untuk menyimpan nilai turunan.
--
-- Migrasi TERPISAH: 059 dan 060 mungkin sudah terpasang. Nullable; nol
-- `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.T_PREMIUM_LIST ADD (
  JENIS_ASURANSI VARCHAR2(255)
)
/

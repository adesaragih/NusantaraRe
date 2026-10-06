-- T_PREMIUM_LIST_SUMMARY.PL_NUMBER - nomor PL pada baris rekap.
--
-- `[keputusan work owner 03-10-2026]` "tambahkan PL_NUMBER juga di tabel
-- T_PREMIUM_LIST_SUMMARY". Menggantikan keputusan tiket 05a ("tetap tanpa
-- kolom PL_NUMBER"). Rekapnya dihitung saat Save (nomor belum ada, kolom
-- NULL); kolom ini diisi saat Confirm, di transaksi yang sama dengan
-- penerbitan nomornya (repository/polis_summary.go `TulisNomorRekap`).
--
-- Tipe sama dengan T_PREMIUM_LIST_DETAIL.PL_NUMBER (052). Nullable; nol
-- `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.T_PREMIUM_LIST_SUMMARY ADD (
  PL_NUMBER VARCHAR2(255)
)
/

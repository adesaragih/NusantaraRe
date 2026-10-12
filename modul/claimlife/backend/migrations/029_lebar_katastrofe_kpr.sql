-- T_GENERAL_CLAIM.STS_KATASTROFE dan IS_KPR diperlebar 8 -> 20 - keputusan work owner 11-10-2026 ("Perlebar kolom"),
-- persetujuan migrasi diberikan 11-10-2026. Bug DEV: ORA-12899 value too large for column STS_KATASTROFE (actual: 15,
-- maximum: 8) saat menjalankan pengisian bawaan.
--
-- Migrasi 002 menulis keduanya VARCHAR2(8) ("keputusan tiket 14", tanpa bukti lebar). Nilai yang sebenarnya:
--   `[terverifikasi]` Section/InputOSClaimLife.xml (korpus kelvin/InboxClaimLife (Done)) b3728 pyDefaultValue
--   "Non-Catastrophe" = 15 byte; "Catastrophe" = 11 byte (claimprop grilling ronde 2 P12);
--   Is KPR "KPR" / "Non-KPR" = 3 / 7 byte (keputusan work owner 11-10-2026).
-- 20 memberi ruang tanpa mengarang daftar nilai lain. Layanan menolak isian > 20 byte (400) sebelum bind.
--
-- ⚠️ Jatah migrasi Claim Life 001-029 HABIS dengan berkas ini (MODUL.md).
--
-- ⛔ NOL COMMIT (ADR-U-0029). Dijalankan work owner sesudah giliran selesai; WAJIB sebelum kode Register / Save grid
-- 11-10-2026 dipasang.

ALTER TABLE {skema}.T_GENERAL_CLAIM MODIFY (
  STS_KATASTROFE VARCHAR2(20),
  IS_KPR         VARCHAR2(20)
)
/

-- T_CLAIMLF_PREMIUMLIST_DETAIL.IS_ACCEPT - keputusan work owner 10-10-2026 (layar Outstanding dirombak mengikuti Pega;
-- persetujuan migrasi diberikan 10-10-2026).
--
-- `[terverifikasi]` `Section/InputOSClaimLife.xml` (korpus kelvin/InboxClaimLife (Done)): grid Participant Details
-- `pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail` b15924, kolom Accept = checkbox `.IsAccept` b17582,
-- tombol kepala kolom `Select All` -> `SelectAllClaimLife_act` b16633. Dibaca `SaveOutStandingLife_Act` langkah 3:
-- peserta `.IsAccept=="true"` wajib punya dokumen ("The document hasn't been uploaded person number N").
--
-- ⛔ BUKAN `IS_CHECK`: IS_CHECK = dipilih untuk diklaim saat Register (`.IsCheck` SavePesertaClaim 7.8). Dua
-- penanda, dua kolom.
--
-- Teks apa adanya seperti IS_CHECK ('true' / 'false' / NULL). Nullable (ADR-U-0027), tanpa pengisian ulang: baris
-- lama = belum di-Accept.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Dijalankan work owner sesudah giliran selesai.

ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD (
  IS_ACCEPT VARCHAR2(8)
)
/

-- T_CLAIMLF_JEJAK.KOMENTAR - OQ-M5 DITUTUP 29-09-2026 (GILIRAN-17), keputusan
-- work owner: alasan penolakan Admin disimpan di kolom komentar BARU jejak klaim.
--
-- `[terverifikasi]` dialog Reject Outstanding `Section/RejectOSClaimLife_Sec.xml`:
--
--   Remarks  b1680 (label, `pyUseLabelDesc` b1674) -> `.DataCommitteeTreaty.Remarks`
--            b1687, pxTextArea b1690, WAJIB b1653/b1695/b1699
--   Submit   b3098 -> `RejectOSClaimLife_Act` b3117
--
-- `RejectOSClaimLife_Act` langkah 6 (b2147, `pyStepsBlockName` kosong) menyalin
-- Remarks ke `KomiteList(<APPEND>).KomiteComment` (b2262-b2263) bersama
-- `IDKomite = "Claim Admin"` (b2241-b2242). Baris "Claim Admin" bukan tingkat
-- komite, dan `T_KOMITE_KOMITELIST` terikat FK ke `T_GENERAL_KOMITE`; karena itu
-- tempatnya jejak klaim, bukan riwayat Komite.
--
-- ⚠️ LEBAR 4000: korpus tidak mengekspor rule Property `Remarks` maupun
-- komentar Komite, jadi lebarnya diambil dari preseden kolom yang memuat nilai
-- yang sama: kolom komentar anggota `T_KOMITE_KOMITELIST`, VARCHAR2(4000) (013).
-- Layanan menolak alasan yang lebih panjang (`services.BatasKomentarJejak`),
-- supaya kelebihan lebar menjadi 400, bukan ORA-12899.
--
-- Nullable: baris jejak lain (transisi, jalur balik tahap) tidak berkomentar.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Dijalankan work owner sesudah giliran selesai.

ALTER TABLE {skema}.T_CLAIMLF_JEJAK ADD (
  KOMENTAR VARCHAR2(4000)
)
/

-- 213 - Form tahap MARKETING (BINDING) (tiket 64): ConfirmBinding + Received RI Slip, di T_GENERAL_POLIS.
--
-- Keputusan work owner 08-10-2026 (diteruskan sesi 9d, D1 "2 kolom"). CONFIRM_BINDING (`.OfferFacIn.ConfirmBinding`, teks Pega
-- "1" / "0"; diisi tombol Binding `UpdateTglBindJsonOffer_Act`, direset Revise `SetReviseProposal`) dan RECEIVED_RI_SLIP
-- (`.OfferFacIn.ReceivedRiSlip`, teks "true" / "false"; checkbox `Section\EmailSectionCeding.xml`) = kolom RANCANGAN
-- T_GENERAL_POLIS (`loader/skema_gen.go`, nama dan tipe apa adanya). COMMENT_CEDING TIDAK ditambah (komentar Binding ke
-- riwayat KETERANGAN). Kolom BARU lewat ALTER biasa (pola 208 / 210 / 212). Ditulis, TIDAK dijalankan agent. Nol COMMIT
-- (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  CONFIRM_BINDING  VARCHAR2(50),
  RECEIVED_RI_SLIP VARCHAR2(50)
)
/

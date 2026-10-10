-- 210 - Submit tahap Team Leader (tiket 58): isian Deductible Acceptance + tanggal ke UW, di T_GENERAL_POLIS.
--
-- IS_DEDUCTIBLE_ACCEPTANCE (`.OfferFacIn.IsDeductibleAcceptance`, teks Pega "true" / "false" apa adanya - butir 68.1) dan
-- COMMENT_DEDUCTIBLE (`.OfferFacIn.CommentDeductible`) = kolom RANCANGAN T_GENERAL_POLIS (`loader/skema_gen.go`, tipe apa
-- adanya), `Section\Periode.xml` + Protection_Act langkah 9-12. DATE_TO_UW / SUBMIT_TO_UW (`.OfferFacIn.DateToUW` /
-- `.SubmitToUW`, `DataTransform\IsJuniorUW_B.xml`, teks Pega DateTime) ada di rancangan TETAPI di T_WORK_POLIS - tabel milik
-- premiumlistlife (K-064) - jadi ditaruh di T_GENERAL_POLIS (berbagi PK dengan T_WORK_POLIS) lewat amandemen loader, tipe =
-- rancangan VARCHAR2(30) `[keputusan agent]`. Kolom BARU lewat ALTER biasa (pola 195 / 203 / 208). Ditulis, TIDAK
-- dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  IS_DEDUCTIBLE_ACCEPTANCE VARCHAR2(10),
  COMMENT_DEDUCTIBLE       VARCHAR2(500),
  DATE_TO_UW               VARCHAR2(30),
  SUBMIT_TO_UW             VARCHAR2(30)
)
/

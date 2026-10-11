-- 212 - Submit tahap Underwriting (tiket 61): Validity Of Acceptance + Comment, di T_GENERAL_POLIS.
--
-- Keputusan work owner 08-10-2026 (diteruskan sesi 9d, D1 "3 kolom + COMMENT_TEXT"). DAYS_VALIDITY (`.OfferFacIn.DaysValidity`,
-- teks bulat), DATE_VALIDITY / WAITING_BIND_DATE (`.OfferFacIn.DateValidity` / `.WaitingBindDate`, teks Pega DateTime
-- "yyyyMMddT170000.000 GMT" = 00:00 WIB, seperti ekspor DDL\CONTOH) dan COMMENT_TEXT (`.OfferFacIn.Comment`) = kolom
-- RANCANGAN T_GENERAL_POLIS (`loader/skema_gen.go`, tipe apa adanya); `Section\EmailSection_IsUW.xml` +
-- `Activity\GetDateValidity_ACT.xml`. Kolom BARU lewat ALTER biasa (pola 208 / 210). Ditulis, TIDAK dijalankan agent. Nol
-- COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  DAYS_VALIDITY     VARCHAR2(50),
  DATE_VALIDITY     VARCHAR2(30),
  WAITING_BIND_DATE VARCHAR2(30),
  COMMENT_TEXT      VARCHAR2(500)
)
/

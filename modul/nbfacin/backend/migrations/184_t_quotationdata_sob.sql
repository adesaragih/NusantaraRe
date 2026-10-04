-- 184 - T_QUOTATIONDATA.SOURCE_OF_BUSINESS: kode SOB popup Change SOB (tiket 33).
--
-- Kolom RANCANGAN (bukan kolom baru karangan): `loader/skema_gen.go` T_QUOTATIONDATA
-- SOURCE_OF_BUSINESS VARCHAR2(50) <- `.QuotationData.SourceOfBusiness`. `[terverifikasi]` isinya kode:
-- 5 fixture `services/premium/testdata/kasus` berbentuk huruf + 7 digit (G..., B...), sejalan
-- ekspor case yang dikutip tiket 33 dan pasangan DDL FACINPRODUCTION "SOB" (nama) + "SOBID" (kode).
-- Nama SOB tetap di SOB_NAME (183). Ditambah lewat ALTER karena 183 membuat tabel sebagian
-- (butir 78.4). Nullable. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_QUOTATIONDATA ADD (
  SOURCE_OF_BUSINESS VARCHAR2(50)
)
/

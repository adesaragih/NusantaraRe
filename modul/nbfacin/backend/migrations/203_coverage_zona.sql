-- 203 - tab Coverage FIRE (tiket 52): T_COVERAGELIST + Zone dan 4.2 Construction.
--
-- ZONE ada di rancangan T_COVERAGELIST (`.Zone`, VARCHAR2(50)) tetapi dilebarkan VARCHAR2(1000) = lebar sumbernya
-- ZONEOJK.DESCRIPTION (DDL `ZONEOJK.txt`; autocomplete `BrowseZones_RD` menulis ZONES.DESCRIPTION ke `.Zone`, dan
-- `RateEarthquake_Sql` mencocokkan `.Zone` ke ZONEOJK.DESCRIPTION) - pola butir 80. Sampel `DDL\CONTOH` sampai 111 byte
-- (melampaui 50). CONSTRUCTION_EARTHQUAKE (`.ConstructionEarthquake`, isian `SearchClassConstructionReas_SQL`) tidak
-- ada di rancangan - amandemen loader; VARCHAR2(500) = kolom konstruksi saudaranya (*_CONSTRUCTION, migrasi 187).
-- Kolom BARU lewat ALTER biasa (pola 195). Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_COVERAGELIST ADD (
  ZONE                    VARCHAR2(1000),
  CONSTRUCTION_EARTHQUAKE VARCHAR2(500)
)
/

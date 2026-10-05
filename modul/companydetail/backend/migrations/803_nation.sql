-- 803 - NATION sebagai tabel datar: nama dan kolom SAMA dengan view yang dibuang 802 (pra-terbang membandingkan
-- kolomnya dengan view itu - keempatnya cocok).
--
-- Lebar dari isi DEV 04-10-2026: ID 6 digit (57 baris, unik, tanpa kosong), OLDID maks. 3, NOTE maks. 18,
-- NATIONINITIAL maks. 17. COUNTRY di CLIENT = OLDID (18.326/18.326 baris Org cocok), COUNTRYNAME = NOTE.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
CREATE TABLE {skema}.NATION (
  ID            VARCHAR2(10) NOT NULL,
  OLDID         VARCHAR2(6),
  NOTE          VARCHAR2(100),
  NATIONINITIAL VARCHAR2(20),
  CONSTRAINT PK_NATION PRIMARY KEY (ID)
)
/

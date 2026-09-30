-- T_TREATYYEAR - tahun treaty (tiket 01 Treaty Contract Out).
--
-- Satu baris = satu tahun treaty per grup: wadah seluruh kontrak tahun itu.
-- Sumber warisan POOLDATA.TREATYYEAR, ditulis PEGA_TREATYYEAR
-- (RDBList/SaveMasterTreatyYear_SQL.xml, 10 parameter + 2 keluaran).
--
-- ⛔ AWALAN T_ (keputusan tco1, veto work owner). Nama warisan TREATYYEAR
-- sudah dipakai tabel hidup di skema POOLDATA yang sama. Tabel warisan
-- TIDAK disentuh: tidak ditulis, tidak di-ALTER, hanya dibaca migrasi data.
--
-- ⛔ NAMA KOLOM VERBATIM dari warisan. Yang berubah hanya TIPE: STARTDATE,
-- ENDDATE, TGLUPDATE dari VARCHAR2 menjadi DATE (AC 53, penyimpangan sadar 6).
--
-- ⚠️ PROPORTION tetap teks: artinya `[terbuka]`; layar mengisinya dari
-- pilihan "Reinsurance Type" (InputDtlTreatyContact.xml b6829). Kolom uang
-- dan persen yang wajib desimal (AC 51/52) tidak menyebutnya.
--
-- ⛔ TIPE FISIK: teks -> VARCHAR2(255), pengenal -> VARCHAR2(32), DATE -> DATE.
-- Seluruh kolom nullable kecuali PK (ADR-U-0027); NOT NULL akan menolak baris
-- warisan yang kosong saat migrasi data.
--
-- ⛔ Anti-dobel (STARTDATE, ENDDATE, TREATYGROUPID) - AC 73 - ditegakkan di
-- Go, BUKAN unique index: data warisan boleh sudah memuat duplikat, dan index
-- unik akan menggagalkan migrasi datanya tanpa menyebut baris mana.
--
-- ⛔ NOL COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_TREATYYEAR (
  ID               VARCHAR2(32) NOT NULL,
  TREATYYEAR       VARCHAR2(255),
  UNDERWRITINGYEAR VARCHAR2(255),
  TREATYGROUPID    VARCHAR2(32),
  TREATYGROUPNAME  VARCHAR2(255),
  USERID           VARCHAR2(255),
  TGLUPDATE        DATE,
  PROPORTION       VARCHAR2(255),
  STARTDATE        DATE,
  ENDDATE          DATE,
  CONSTRAINT PK_T_TREATYYEAR PRIMARY KEY (ID)
)
/

-- Identitas '1' + lpad(6): bentuk warisan TreatyYear_seq dipertahankan
-- (spec §7, ADR-0006). Nilai berjalan diselaraskan migrasi data (AC 69).
CREATE SEQUENCE {skema}.SEQ_T_TREATYYEAR START WITH 1 INCREMENT BY 1 NOCACHE
/

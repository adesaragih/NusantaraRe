-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`/`445`/`446`: rentang `treatyin` (`400-439`) PENUH, dan `448` jatuh di
-- rentang modul ini (`440-479`).
--
-- ---------------------------------------------------------------------
-- Tempat simpan properti yang tombol Save kirim tetapi belum punya kolom
-- ---------------------------------------------------------------------
--
-- ⭐ Keputusan pemilik proses 7 Oktober 2026 ("ini dibuatkan saja"), atas
-- laporan `kunciTakTersimpan` jalur Save (`modul/treatyin/docs/
-- KEPUTUSAN-SASARAN-TULIS.md` §5.1):
--
--   T_TREATY_REVISION (skalar akar `TreatyIn`)
--     RNMSHAREP           TreatyIn.RNMShareP          tab Share Prop
--     BROKERAGEPERCENTP   TreatyIn.BrokeragePercentP  tab Share Prop
--     OPTIONLIMIT         TreatyIn.OptionLimit        tab Share Prop
--     INSTALLMENTNO       TreatyIn.InstallmentNo      tab Installment Non-Prop
--     REVISIONSTATE       TreatyIn.RevisionState      jalur REVISI tangga akseptasi
--     VIEWSTATE           TreatyIn.ViewState          jalur REVISI tangga akseptasi
--
--   T_TREATY_SHARE_SUMMARY (BARU) — dua larik ringkasan tab Share Non-Prop,
--     dibedakan `JENIS`: LimitShareSummaryList, LimitFacShareSummaryList.
--
-- ⚠️ Kesebelas larik TOTAL (`TotalShareRnmProp`, `TotalSpreadedRnmProp`,
-- `TotalShareGrossNP`, `TotalInstallmentNP`, …) TIDAK butuh DDL: bentuknya
-- `{Currency, CurrencyID, Value}`, sama dengan `T_TREATY_TOTAL`, dan kolom
-- `JENIS`-nya yang membedakan — cukup peta pendaratan.
--
-- ⚠️ `VARCHAR2(4000 CHAR)` seperti seluruh kolom pendaratan: nilainya teks
-- apa adanya, ditafsirkan services.
--
-- ⚠️ Penulis dan pembaca `T_TREATY_*` TOLERAN terhadap kolom/tabel yang belum
-- terpasang (`kolomTerpasang`): sebelum berkas ini dijalankan Save tetap
-- berjalan dan melaporkan properti itu sebagai tidak tersimpan.

ALTER TABLE {skema}.T_TREATY_REVISION ADD (
  RNMSHAREP                        VARCHAR2(4000 CHAR),
  BROKERAGEPERCENTP                VARCHAR2(4000 CHAR),
  OPTIONLIMIT                      VARCHAR2(4000 CHAR),
  INSTALLMENTNO                    VARCHAR2(4000 CHAR),
  REVISIONSTATE                    VARCHAR2(4000 CHAR),
  VIEWSTATE                        VARCHAR2(4000 CHAR)
)
/
CREATE TABLE {skema}.T_TREATY_SHARE_SUMMARY (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  JENIS     VARCHAR2(40 CHAR)  NOT NULL,
  LAYERTYPE                        VARCHAR2(4000 CHAR),
  LAYER                            VARCHAR2(4000 CHAR),
  LAYERPARTTYPE                    VARCHAR2(4000 CHAR),
  LAYERPART                        VARCHAR2(4000 CHAR),
  NOTE                             VARCHAR2(4000 CHAR),
  LIMITVAL                         VARCHAR2(4000 CHAR),
  LIMITVAL2                        VARCHAR2(4000 CHAR),
  MDP                              VARCHAR2(4000 CHAR),
  MDP2                             VARCHAR2(4000 CHAR),
  DEDUCTIBLE                       VARCHAR2(4000 CHAR),
  DEDUCTIBLE2                      VARCHAR2(4000 CHAR),
  NETPREMI                         VARCHAR2(4000 CHAR),
  NETPREMI2                        VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_SHARE_SUMMARY PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_SHARE_SUMMARY UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_SHARE_SUMMARY START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

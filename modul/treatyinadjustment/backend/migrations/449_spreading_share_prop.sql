-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`/`445`/`446`/`448`: rentang `treatyin` (`400-439`) PENUH.
--
-- ---------------------------------------------------------------------
-- Tempat simpan rincian Detail tab Share PROPORSIONAL
-- ---------------------------------------------------------------------
--
-- Laporan Save kontrak 1002305 (8 Oktober 2026) — "Belum punya kolom di
-- tabel, jadi TIDAK tersimpan: Currency, NusaReLimit, RNMShareList,
-- RNMSpreadedList, RNMSpreadedListRI, SpreadingList". Pemilik proses:
-- "berikan apa yang anda butuhkan".
--
--   T_TREATY_LIMIT_DETAIL (TreatyIn.Limits.Detail)
--     CURRENCY      .Currency     `TreatyInPropshare` [4]: = .CurrencyIOOLimit
--     NUSARELIMIT   .NusaReLimit  `TreatyInPropshare` [4]: Cession × RNMShareP
--
--   T_TREATY_LIMIT_SPREADING (BARU) — `.Detail.SpreadingList`, anak
--     susunan treaty dari `PROPORTIONALARRG` (`FetchQSfromMaster` [9]).
--     Sebentuk dengan `T_TREATY_SHARE_SPREADING` (cabang Non-Prop), plus
--     `VALUE` (`(.Pct/100) × RNMShareList(1).Value`).
--
-- ⚠️ `RNMShareList`, `RNMSpreadedList`, `RNMSpreadedListRI` TIDAK butuh DDL:
-- bentuknya `{Currency, Value}`, sama dengan `T_TREATY_LIMIT_AMOUNT`, dan
-- kolom `JENIS`-nya (`VARCHAR2(20)`, nama terpanjang 17 aksara) yang
-- membedakan — cukup peta pendaratan.
--
-- ⚠️ TANPA kunci asing ke `T_TREATY_LIMIT_DETAIL` — pola `448`: setiap
-- kaskade baru wajib tercatat di keputusan ERD (`migrasi_invarian_test`), dan
-- ia tidak dibutuhkan: penulis menghapus baris anak per `MASTERID`
-- (`KosongkanKontrak`) sebelum menyisip ulang.
--
-- ⚠️ Penulis dan pembaca `T_TREATY_*` TOLERAN terhadap kolom/tabel yang belum
-- terpasang (`kolomTerpasang`): sebelum berkas ini dijalankan Save tetap
-- berjalan dan melaporkan properti itu sebagai tidak tersimpan.

ALTER TABLE {skema}.T_TREATY_LIMIT_DETAIL ADD (
  CURRENCY                         VARCHAR2(4000 CHAR),
  NUSARELIMIT                      VARCHAR2(4000 CHAR)
)
/
-- T_TREATY_LIMIT_SPREADING  <- TreatyIn.Limits.Detail.SpreadingList
CREATE TABLE {skema}.T_TREATY_LIMIT_SPREADING (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  PARENTREINSTYPEID                VARCHAR2(4000 CHAR),
  PCT                              VARCHAR2(4000 CHAR),
  REINSTYPEID                      VARCHAR2(4000 CHAR),
  REINSTYPENAME                    VARCHAR2(4000 CHAR),
  RP                               VARCHAR2(4000 CHAR),
  USD                              VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_SPREADING PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_SPREADING UNIQUE (IDINDUK, URUTAN)
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_SPREADING_MST ON {skema}.T_TREATY_LIMIT_SPREADING (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_SPREADING START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

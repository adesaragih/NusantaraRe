-- T_TREATYCONTRACT - kontrak treaty di dalam tahun treaty (tiket 01).
--
-- Satu baris = satu jenis reasuransi (REINSTYPEID) yang dibuka pada sebuah
-- tahun treaty. Sumber warisan POOLDATA.TREATYCONTRACT, ditulis
-- PEGA_TREATYCONTRACT (RDBList/SaveMasterTreatyContract_SQL.xml, 8 + 2 out).
--
-- ⛔ AWALAN T_ (tco1); nama kolom VERBATIM. TGLUPDATE warisan VARCHAR2(1000)
-- menjadi DATE; TREATYSTARTDATE/TREATYENDDATE sudah DATE di warisan.
--
-- ⛔ FK ke T_TREATYYEAR TANPA ON DELETE CASCADE. Tidak ada jalur hapus tahun
-- di tiket mana pun; bila kelak ada, Oracle menolak tahun yang masih
-- berkontrak dengan ORA-02292 - gagal terang, bukan hilang diam-diam.
--
-- ⚠️ Kontrak "membuka" kombinasi (TREATYYEAR, TREATYGROUPID, REINSTYPEID)
-- yang dipakai reinsurer, business, dan klausul, tetapi TIDAK memilikinya:
-- ketiganya tidak menyimpan ID kontrak (spec §2, fakta bisnis work owner).
--
-- ⛔ NOL COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_TREATYCONTRACT (
  ID              VARCHAR2(32) NOT NULL,
  IDTREATYYEAR    VARCHAR2(32),
  REINSTYPEID     VARCHAR2(32),
  REINSTYPENAME   VARCHAR2(255),
  TREATYSTARTDATE DATE,
  TREATYENDDATE   DATE,
  USERID          VARCHAR2(255),
  TGLUPDATE       DATE,
  CONSTRAINT PK_T_TREATYCONTRACT PRIMARY KEY (ID),
  CONSTRAINT FK_T_TREATYCONTRACT_TAHUN FOREIGN KEY (IDTREATYYEAR)
    REFERENCES {skema}.T_TREATYYEAR (ID)
)
/

CREATE INDEX {skema}.IDX_T_TREATYCONTRACT_TAHUN ON {skema}.T_TREATYCONTRACT (IDTREATYYEAR)
/

-- Identitas '1' + lpad(6): bentuk warisan treatycontract_seq (spec §7).
CREATE SEQUENCE {skema}.SEQ_T_TREATYCONTRACT START WITH 1 INCREMENT BY 1 NOCACHE
/

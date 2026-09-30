-- T_PROPORTIONALARRG - SATU tabel untuk dua puluh lima jenis klausul (tiket 01).
--
-- Sumber warisan POOLDATA.PROPORTIONALARRG. `[data DBA]` KEDUA procedure
-- penulis - PEGA_PROPORTIONALARRG (35 kolom, 18 jenis induk) dan
-- PEGA_M_PROPORTIONALARRG_CHILD (26 kolom, 7 jenis anak) - menulis ke tabel
-- yang SAMA; nama "_CHILD" menipu (OQ-066). Jenis dibedakan TREATYDESCID.
--
-- ⛔ PENYIMPANGAN SADAR 2: satu tabel generik. Baris anak mengisi NULL pada
-- sembilan kolom khusus induk: ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE,
-- TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD (AC 24, 25).
--
-- ⛔ PENYIMPANGAN SADAR 6: warisan CAMPUR - TREATYLIMIT, COINS_MIN, COINS_MAX,
-- MORERP, MOREUSD sudah NUMBER, tetapi RP, USD, PCT, PCTME VARCHAR2(1000).
-- Di sini seluruh uang dan persen NUMBER(38,8) (AC 51, 52); KURS ikut desimal
-- (tiket 11). RP dan USD tetap DUA kolom terpisah (AC 49).
--
-- ⛔ Kolom warisan PROPORTIONALLIST dan OBJECT TIDAK dibawa (AC 70): tidak
-- di-set procedure mana pun; migrasi data melaporkan bila ada isi hidup.
--
-- ⛔ TANPA FK ke kontrak. Klausul menggantung pada (TREATYYEAR, TREATYGROUPID,
-- REINSTYPEID), milik level tahun/grup/jenis; kaskade hapus kontrak (tiket 10)
-- TIDAK menyentuhnya - desain, bukan bug (fakta bisnis work owner).
--
-- ⚠️ PARENTREINSTYPEID = "00" adalah sentinel baris tanpa induk
-- (Activity/GetPeriode.xml b889); dibawa apa adanya. Istilah LAYER, LAYERPART,
-- LAYERPARTTYPE, LAYERTYPE, LINE, YDCF, METHOD, SPREADINGORDER dipakai apa
-- adanya dari Pega (spec §4), artinya OQ Product + UW.
--
-- ⛔ NOL COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_PROPORTIONALARRG (
  ID                VARCHAR2(32) NOT NULL,
  TREATYYEAR        VARCHAR2(255),
  TREATYYEARID      VARCHAR2(32),
  TREATYGROUPID     VARCHAR2(32),
  TREATYGROUPNAME   VARCHAR2(255),
  TREATYDESCID      VARCHAR2(32),
  TREATYDESCNAME    VARCHAR2(255),
  REINSTYPEID       VARCHAR2(32),
  REINSTYPENAME     VARCHAR2(255),
  LAYER             VARCHAR2(255),
  LAYERPART         VARCHAR2(255),
  LAYERPARTTYPE     VARCHAR2(255),
  LAYERTYPE         VARCHAR2(255),
  KURS              NUMBER(38,8),
  TGLUPDATE         DATE,
  USERID            VARCHAR2(255),
  LINE              VARCHAR2(255),
  PCT               NUMBER(38,8),
  PCTME             NUMBER(38,8),
  YDCF              VARCHAR2(255),
  METHOD            VARCHAR2(255),
  TERRITORIALLIMIT  VARCHAR2(1000),
  PARENTREINSTYPEID VARCHAR2(32),
  SPREADINGORDER    VARCHAR2(255),
  RP                NUMBER(38,8),
  USD               NUMBER(38,8),
  ID_OCCUPATION     VARCHAR2(32),
  OCCUPATION        VARCHAR2(1000),
  ID_CLAUSE         VARCHAR2(32),
  CLAUSE            VARCHAR2(1000),
  TREATYLIMIT       NUMBER(38,8),
  COINS_MIN         NUMBER(38,8),
  COINS_MAX         NUMBER(38,8),
  MORERP            NUMBER(38,8),
  MOREUSD           NUMBER(38,8),
  CONSTRAINT PK_T_PROPORTIONALARRG PRIMARY KEY (ID)
)
/

-- Kunci gabungan yang dipakai kaskade, daftar per kombinasi, dan hilir.
CREATE INDEX {skema}.IDX_T_PROPARRG_KOMB
  ON {skema}.T_PROPORTIONALARRG (TREATYYEAR, TREATYGROUPID, REINSTYPEID)
/

-- Bentuk WHERE GetMasterDescriptionLimitParentList.xml (TreatyTestChildTotal_Act):
-- TREATYYEARID + TREATYGROUPID + TREATYDESCID + PARENTREINSTYPEID.
CREATE INDEX {skema}.IDX_T_PROPARRG_DESC
  ON {skema}.T_PROPORTIONALARRG (TREATYYEARID, TREATYGROUPID, TREATYDESCID, PARENTREINSTYPEID)
/

-- Identitas '1' + lpad(7): TUJUH digit, berbeda dari lima tabel lain
-- (bentuk warisan PROPORTIONALARRG_SEQ, spec §7).
CREATE SEQUENCE {skema}.SEQ_T_PROPORTIONALARRG START WITH 1 INCREMENT BY 1 NOCACHE
/

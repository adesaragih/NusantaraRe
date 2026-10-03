-- 328 - T_POLIS_SUGGEST <- PolicyTreatyIn.SuggestList (catatan dan putusan
-- per tahap; AC 71, 72). RALAT rancangan §4bis.1: riwayat produksi hanya
-- ditulis untuk bisnis fakultatif, jadi daftar treaty butuh tabel sendiri.
-- NOURUT = nomor urut baris di dalam induknya, unik per induk (ID-11, AC 8, 10);
-- kunci pasangan antar generasi, bukan kunci dagang (ID-13, AC 11).
CREATE TABLE {skema}.T_POLIS_SUGGEST (
  ID             VARCHAR2(32) NOT NULL,
  POLIS_ID       VARCHAR2(32) NOT NULL,
  NOURUT         NUMBER(5) NOT NULL,
  SUGGEST        VARCHAR2(4000),
  IS_APPROVED    VARCHAR2(16),
  SUGGEST_DATE   DATE,
  OPERATOR_NAME  VARCHAR2(128),
  OPERATOR_ID    VARCHAR2(64),
  IS_SAVE        VARCHAR2(16),
  CONSTRAINT PK_POLIS_SUGGEST PRIMARY KEY (ID),
  CONSTRAINT FK_POLIS_SUGGEST_POLIS FOREIGN KEY (POLIS_ID) REFERENCES {skema}.T_GENERAL_POLIS (ID),
  CONSTRAINT UQ_POLIS_SUGGEST_NOURUT UNIQUE (POLIS_ID, NOURUT)
)
/

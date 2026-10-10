-- Pembalikan `447_buang_t_treaty_currency.sql`.
--
-- ⭐ BENTUKNYA DIBACA DARI ORACLE, BUKAN DIKARANG. Tabel ini nol pernah
-- dibuat migrasi, jadi nol DDL asli yang dapat disalin. Yang di bawah
-- disusun dari `ALL_TAB_COLUMNS`, `ALL_CONSTRAINTS`, dan `ALL_CONS_COLUMNS`
-- pada 7 Oktober 2026 — delapan kolom, satu PK, satu UQ:
--
--   ID          NUMBER presisi=19 skala=0  NOT NULL   PK_TT_CURRENCY (ID)
--   MASTERID    VARCHAR2 100 CHAR          NOT NULL   UQ_TT_CURRENCY
--   URUTAN      NUMBER presisi=10 skala=0  NOT NULL     (MASTERID, URUTAN)
--   CURRENCY    VARCHAR2 4000 CHAR         NULL
--   CURRENCYID  VARCHAR2 4000 CHAR         NULL
--   CONVERSION  VARCHAR2 4000 CHAR         NULL
--   PERIODSTART VARCHAR2 4000 CHAR         NULL
--   PERIODEND   VARCHAR2 4000 CHAR         NULL
--
-- ⚠️ NOL `PXOBJCLASS` — tabel ini memang tidak punya, beda dengan kedelapan
-- tabel migrasi `430`. Dipertahankan apa adanya: jalur mundur yang
-- "merapikan" bentuk aslinya bukan pembalikan.
--
-- ⛔ DATANYA TIDAK KEMBALI, dan tidak ada yang hilang karenanya: tabelnya
-- 0 baris saat dibuang, dan nol kode pernah mengisinya.

CREATE TABLE {skema}.T_TREATY_CURRENCY (
  ID           NUMBER(19)            NOT NULL,
  MASTERID     VARCHAR2(100 CHAR)    NOT NULL,
  URUTAN       NUMBER(10)            NOT NULL,
  CURRENCY     VARCHAR2(4000 CHAR),
  CURRENCYID   VARCHAR2(4000 CHAR),
  CONVERSION   VARCHAR2(4000 CHAR),
  PERIODSTART  VARCHAR2(4000 CHAR),
  PERIODEND    VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_CURRENCY PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_CURRENCY UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_CURRENCY START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

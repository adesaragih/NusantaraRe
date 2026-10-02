-- 146 - M_PRODUCTNAME_LIFE_OUTWARD: anak M_PRODUCTNAME_LIFE, satu baris = satu baris `OutwardList[*]` berisi.
--
-- PK (PRODUCTID, URUT): URUT = urutan baris di grid, mulai 1. Index PK berawalan PRODUCTID sekaligus melayani FK
-- (nol index terpisah). FK berkaskade: baris anak hidup DI DALAM halaman produk - simpan menulis ulang seluruh anak
-- satu produk di transaksi yang sama dengan induknya. Kolom isi NULLABLE (tiket 01 AC 43).
CREATE TABLE {skema}.M_PRODUCTNAME_LIFE_OUTWARD (
  PRODUCTID        VARCHAR2(6) NOT NULL,
  URUT             NUMBER(5) NOT NULL,
  REINSTYPEID      VARCHAR2(100),
  REINSTYPENAME    VARCHAR2(100),
  TRANSACTIONYEAR  NUMBER(5),
  TREATYCONTRACTID VARCHAR2(100),
  UNDERWRITINGYEAR NUMBER(5),
  OVR_COMM         NUMBER(38,8),
  CONSTRAINT PK_M_PRODUCTNAME_LIFE_OUTWARD PRIMARY KEY (PRODUCTID, URUT),
  CONSTRAINT FK_M_PRODUCTNAME_LIFE_OUTWARD FOREIGN KEY (PRODUCTID)
    REFERENCES {skema}.M_PRODUCTNAME_LIFE (ID) ON DELETE CASCADE
)
/

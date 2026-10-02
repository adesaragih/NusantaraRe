-- 144 - M_PRODUCTNAME_LIFE_FINUW: anak M_PRODUCTNAME_LIFE, satu baris = satu baris `FinancialUnderwritingList[*]`.
--
-- PK (PRODUCTID, URUT): URUT = urutan baris di grid, mulai 1. Index PK berawalan PRODUCTID sekaligus melayani FK
-- (nol index terpisah). FK berkaskade: baris anak hidup DI DALAM halaman produk - simpan menulis ulang seluruh anak
-- satu produk di transaksi yang sama dengan induknya. Kolom isi NULLABLE (tiket 01 AC 43).
CREATE TABLE {skema}.M_PRODUCTNAME_LIFE_FINUW (
  PRODUCTID    VARCHAR2(6) NOT NULL,
  URUT         NUMBER(5) NOT NULL,
  MININSURED   NUMBER(38,8),
  MAXINSURED   NUMBER(38,8),
  EMPLOYEE     VARCHAR2(1000),
  NON_EMPLOYEE VARCHAR2(1000),
  CONSTRAINT PK_M_PRODUCTNAME_LIFE_FINUW PRIMARY KEY (PRODUCTID, URUT),
  CONSTRAINT FK_M_PRODUCTNAME_LIFE_FINUW FOREIGN KEY (PRODUCTID)
    REFERENCES {skema}.M_PRODUCTNAME_LIFE (ID) ON DELETE CASCADE
)
/

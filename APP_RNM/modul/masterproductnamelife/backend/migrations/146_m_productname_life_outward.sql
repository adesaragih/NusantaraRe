-- 146 - M_PRODUCTNAME_LIFE_OUTWARD: anak M_PRODUCTNAME_LIFE, satu baris = satu baris `OutwardList[*]` berisi.
--
-- PK (PRODUCTID, URUT): URUT = urutan baris di grid, mulai 1. Index PK berawalan PRODUCTID sekaligus melayani FK
-- (nol index terpisah). FK berkaskade: baris anak hidup DI DALAM halaman produk - simpan menulis ulang seluruh anak
-- satu produk di transaksi yang sama dengan induknya. Kolom isi NULLABLE (tiket 01 AC 43).
--
-- Dua bentuk baris (uji kering DEV 02-10-2026): baris `On Retention` (REINSTYPEID ... OVR_COMM, 2 baris) dan baris
-- reasuradur outward + rate-nya (OUTWARDNAMEID ... OUTWARDRATE, 189 baris) - keempat kolom terakhir keputusan work owner
-- 02-10-2026 OQ-FLAT-08 "pindahkan". Migrasi ini belum pernah dijalankan di DEV.
CREATE TABLE {skema}.M_PRODUCTNAME_LIFE_OUTWARD (
  PRODUCTID        VARCHAR2(6) NOT NULL,
  URUT             NUMBER(5) NOT NULL,
  REINSTYPEID      VARCHAR2(100),
  REINSTYPENAME    VARCHAR2(100),
  TRANSACTIONYEAR  NUMBER(5),
  TREATYCONTRACTID VARCHAR2(100),
  UNDERWRITINGYEAR NUMBER(5),
  OVR_COMM         NUMBER(38,8),
  OUTWARDNAMEID    VARCHAR2(100),
  OUTWARDNAME      VARCHAR2(1000),
  OUTWARDRATEID    VARCHAR2(100),
  OUTWARDRATE      VARCHAR2(500),
  CONSTRAINT PK_M_PRODUCTNAME_LIFE_OUTWARD PRIMARY KEY (PRODUCTID, URUT),
  CONSTRAINT FK_M_PRODUCTNAME_LIFE_OUTWARD FOREIGN KEY (PRODUCTID)
    REFERENCES {skema}.M_PRODUCTNAME_LIFE (ID) ON DELETE CASCADE
)
/

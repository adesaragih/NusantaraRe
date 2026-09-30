-- T_MTREATYSECURITY - security di bawah seorang reinsurer (tiket 01).
--
-- Sumber warisan POOLDATA.MTREATYSECURITY: `[data DBA]` TANPA primary key,
-- ditulis SQL mentah posisional tujuh nilai (RDBList/InsertToMTreatySecurity.xml)
-- yang mengosongkan TOP_ID, TP_TREATY, USER_ID; diperbarui dan dihapus
-- berkunci REAS_ID + trim(REAS_SECURITY).
--
-- ⛔ PENYIMPANGAN SADAR 5 (keputusan work owner): PK SURROGATE `ID` dari
-- sequence, seluruh kolom bernama, PCT_SHARE desimal, REAS_SECURITY atribut
-- biasa - bukan bagian kunci. trim() di kunci TIDAK dibawa (AC 18-20).
--
-- ⚠️ TOP_ID, TP_TREATY, USER_ID dibawa sebagai kolom bernama `[terbuka]`:
-- DDL-nya diketahui, artinya tidak terbaca korpus. Pembuangannya keputusan
-- terpisah (tiket 06 blocker), bukan inisiatif migrasi.
--
-- ⛔ Satu-satunya FK BERKASKADE di modul ini: security memang milik
-- reinsurernya - DeleteFromTreatyReinsurer_Act.xml menghapus security lalu
-- reinsurer. Kaskade kontrak tiket 10 tetap menghapusnya eksplisit.
--
-- ⛔ NOL COMMIT (ADR-U-0029).
CREATE TABLE {skema}.T_MTREATYSECURITY (
  ID            VARCHAR2(32) NOT NULL,
  THN_TREATY    VARCHAR2(255),
  TOP_ID        VARCHAR2(32),
  TP_TREATY     VARCHAR2(255),
  REAS_ID       VARCHAR2(32),
  PCT_SHARE     NUMBER(38,8),
  USER_ID       VARCHAR2(255),
  REAS_SECURITY VARCHAR2(255),
  CONSTRAINT PK_T_MTREATYSECURITY PRIMARY KEY (ID),
  CONSTRAINT FK_T_MTREATYSECURITY_REAS FOREIGN KEY (REAS_ID)
    REFERENCES {skema}.T_TREATYREINSURER (ID) ON DELETE CASCADE
)
/

CREATE INDEX {skema}.IDX_T_MTREATYSECURITY_REAS ON {skema}.T_MTREATYSECURITY (REAS_ID)
/

-- Sequence BARU: warisan tidak punya identitas baris. Bentuknya mengikuti
-- lima saudaranya ('1' + lpad(6)) supaya rujukan seragam.
CREATE SEQUENCE {skema}.SEQ_T_MTREATYSECURITY START WITH 1 INCREMENT BY 1 NOCACHE
/

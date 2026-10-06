-- 890 - BORDEREAUX_HISTORY: riwayat persetujuan berkas Bordereaux (keputusan work owner 04-10-2026: "tambah table
-- riwayat", nama "BORDEREAUX_HISTORY"). Padanan Pega: page list `CommentList` (Date, OperatorName, IsApproved, Suggest)
-- yang di Pega HANYA tersimpan di dalam JSON M_BORDEREAUX.DATA_JSON - aplikasi ini tidak memakai JSON lagi, jadi
-- riwayatnya mendapat tabel sendiri. Satu baris per Submit / Approve / Reject (`ActionSubmit` langkah 2).
--
-- ID dari SEQ_BORDEREAUX_HISTORY. BDX_ID tanpa FK ke BORDEREAUX (tabel warisan tanpa jaminan FK; berkas dihapus
-- bersama riwayatnya oleh layanan dalam satu transaksi).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
CREATE TABLE {skema}.BORDEREAUX_HISTORY (
  ID          NUMBER(19) NOT NULL,
  BDX_ID      VARCHAR2(50) NOT NULL,
  TANGGAL     DATE DEFAULT SYSDATE NOT NULL,
  PIC         VARCHAR2(64) NOT NULL,
  IS_APPROVED VARCHAR2(1) NOT NULL,
  KOMENTAR    VARCHAR2(2000),
  CONSTRAINT PK_BORDEREAUX_HISTORY PRIMARY KEY (ID),
  CONSTRAINT CK_BORDEREAUX_HISTORY_APPROVED CHECK (IS_APPROVED IN ('0', '1'))
)
/
CREATE INDEX {skema}.IX_BORDEREAUX_HISTORY_BDX ON {skema}.BORDEREAUX_HISTORY (BDX_ID)
/
CREATE SEQUENCE {skema}.SEQ_BORDEREAUX_HISTORY START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

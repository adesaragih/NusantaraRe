-- 535 - T_KATEGORI_DOC_KLAIM: master kategori dokumen klaim (perintah work owner 08-10-2026: "BUATKAN KATEGORI FILE
-- INI, JADIKAN MASTER T_KATEGORI_DOC_KLAIM", berkas kategoriifile.xls kolom ID / LABEL / TYPE KLAIM).
--
-- ID = kode kategori berkas klaim = KATEGORI_1 tabel warisan dokumen klaim (lebar sama, VARCHAR2(100)); kategori
-- lampiran yang AttachmentProtect_ACT periksa (LOD, DLA, SPGR, ADU, Invoice, Salvage) ada di sini. ID yang sama
-- dipakai lebih dari satu TYPE_KLAIM (FAC / PROP / NONPROP), maka kuncinya (TYPE_KLAIM, ID).
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
CREATE TABLE {skema}.T_KATEGORI_DOC_KLAIM (
  ID          VARCHAR2(100) NOT NULL,
  LABEL       VARCHAR2(255) NOT NULL,
  TYPE_KLAIM  VARCHAR2(20)  NOT NULL,
  CONSTRAINT PK_T_KATEGORI_DOC_KLAIM PRIMARY KEY (TYPE_KLAIM, ID)
)
/

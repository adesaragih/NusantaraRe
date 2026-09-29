-- T_TREATYYEAR_LAMPIRAN - lampiran berkas pada tahun treaty (tiket 12).
--
-- ⚠️ FITUR BARU, bukan migrasi (penyimpangan sadar 9, keputusan work owner).
-- Jalur lampiran Pega (`M_ATTACHMENTTREATY_2`, `PEGA_M_ATTACHMENT`) belum
-- rampung dan dikunci ke ID treaty INWARD (`GetAllAttachment2_Sql.xml`
-- `where treatyid = {TreatyIn.ID}`). Di sini lampiran melekat pada TAHUN
-- treaty; tabel warisan itu tetap milik konteks treaty inward, tidak dibaca.
--
-- ⛔ Nama kolom mengikuti tabel warisan di mana maknanya sama
-- (`FILENAME`, `CATEGORY`, `FILEMIMETYPE`, `T_STORAGE_ID` -
-- `GetAllAttachment2_Sql.xml`), kecuali kunci induk: `IDTREATYYEAR`, bukan
-- `TREATYID`.
--
-- ⛔ `IMAGEID` lahir BERSAMA barisnya (acak, `models.ImageIDBaru`) dan tidak
-- pernah berubah: ia kunci berkas di penyimpanan. Pengulangan unggah menulis
-- ke kunci yang SAMA, jadi tidak menggandakan berkas (AC 58).
-- `T_STORAGE_ID` baru terisi sesudah penyimpanan memastikan berkasnya ada -
-- kosong berarti "tertunda".
--
-- ⛔ `CATEGORY` teks `NOTE` master `CATEGORY_ATTACH_REAS`
-- (`SetCategoryAttachTreatyin.xml` b500 `.NOTE`); kolom ID master tidak
-- terbukti di korpus (`CategoryAttach_SQL.xml` `select *`), jadi tidak dipakai.
--
-- ⛔ FK ke T_TREATYYEAR TANPA kaskade: kaskade diam-diam meninggalkan berkas
-- yatim di penyimpanan tanpa satu pun efek hapus.
--
-- ⛔ NOL COMMIT (ADR-U-0029). Nol CLOB/BLOB: isi berkas di penyimpanan.
CREATE TABLE {skema}.T_TREATYYEAR_LAMPIRAN (
  ID           VARCHAR2(32) NOT NULL,
  IDTREATYYEAR VARCHAR2(32) NOT NULL,
  FILENAME     VARCHAR2(255),
  FILEMIMETYPE VARCHAR2(255),
  CATEGORY     VARCHAR2(255),
  IMAGEID      VARCHAR2(64) NOT NULL,
  T_STORAGE_ID VARCHAR2(64),
  UKURAN       NUMBER(19),
  USERID       VARCHAR2(255),
  TGLUPLOAD    DATE,
  CONSTRAINT PK_T_TREATYYEAR_LAMPIRAN PRIMARY KEY (ID),
  CONSTRAINT UQ_T_TYLAMPIRAN_IMAGEID UNIQUE (IMAGEID),
  CONSTRAINT FK_T_TYLAMPIRAN_TAHUN FOREIGN KEY (IDTREATYYEAR)
    REFERENCES {skema}.T_TREATYYEAR (ID)
)
/

CREATE INDEX {skema}.IDX_T_TYLAMPIRAN_TAHUN ON {skema}.T_TREATYYEAR_LAMPIRAN (IDTREATYYEAR)
/

CREATE SEQUENCE {skema}.SEQ_T_TREATYYEAR_LAMPIRAN START WITH 1 INCREMENT BY 1 NOCACHE
/

-- Tiga tabel NILAI anak `437` — struktur dari
-- `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`, kolom dari PENGUKURAN.
--
-- ---------------------------------------------------------------------
-- PENGUKURAN YANG MENDAHULUI BERKAS INI
-- ---------------------------------------------------------------------
-- Disapu 6 Oktober 2026 atas SELURUH 1.855 baris `POOLDATA.M_TREATY_IN`
-- (nol `ROWNUM`, nol sampel; 1.855 dokumen terurai, nol gagal urai):
--
--   Limits[].Detail[].IOOLimitList            2.870 elemen · 1.077 dokumen
--   Limits[].Detail[].RetentionList           2.870 elemen · 1.078 dokumen
--   Limits[].Detail[].CessionList             2.868 elemen · 1.077 dokumen
--   Limits[].Detail[].EPIList                 2.867 elemen · 1.039 dokumen
--                                   jumlah   11.475 elemen
--
--   Share[].GrossPremiumList                  3.049 elemen ·   845 dokumen
--   Share[].NetPremiumList                    3.048 elemen ·   845 dokumen
--                                   jumlah    6.097 elemen
--
--   FacultativeShareList[].GrossPremiumList      12 elemen ·     5 dokumen
--   FacultativeShareList[].NetPremiumList        12 elemen ·     5 dokumen
--                                   jumlah       24 elemen
--
-- ⚠ `DIPETAKAN` BUKAN izin membangun, dan itu sudah terbukti sekali:
--   `T_TREATY_FAC_LIMITS` jalurnya utuh dan lima kolomnya `DIPETAKAN`,
--   tetapi sapuan menemukan NOL elemen — dan tabel tanpa satu pun baris
--   muatan bukan tabel. Ketiga tabel di bawah dibangun sebab cacahnya
--   SUNGGUH lebih besar dari nol, bukan sebab jalurnya bertanda.
--
-- ⚠ Yang PALING TIPIS: `T_TREATY_FAC_SHARE_AMOUNT`, 24 elemen di 5 dokumen
--   dari 1.855. Ia dibangun sebab bukan nol, dan ketipisannya dinyatakan di
--   sini supaya siapa pun yang kelak menimbangnya tidak perlu mengukur lagi.
--
-- ---------------------------------------------------------------------
-- SATU KOLOM YANG TIDAK ADA DI DOKUMEN, DAN SEBABNYA DIPAKSA BENTUK
-- ---------------------------------------------------------------------
-- `JENIS` tidak diturunkan dari satu pun kunci JSON. Ia ada sebab xlsx
-- menaruh EMPAT larik sumber ke dalam SATU tabel:
--
--   T_TREATY_LIMIT_AMOUNT  <- EPIList | RetentionList | CessionList | IOOLimitList
--
-- Tanpa pembeda, keempat larik itu menjadi baris yang TIDAK DAPAT DIBEDAKAN
-- satu sama lain -- `Currency` dan `Value` saja tidak memberi tahu apakah
-- satu baris adalah retensi atau cession. Itu lebih buruk daripada tidak
-- membangun tabelnya.
--
-- ⚠ Jadi `JENIS` DIPAKSA oleh bentuk yang xlsx pilih, bukan dipilih di sini.
--   Nilainya nama lariknya apa adanya: `EPI` · `RETENTION` · `CESSION` ·
--   `IOOLIMIT` · `GROSSPREMIUM` · `NETPREMIUM`.
--
-- ---------------------------------------------------------------------
-- KOLOM MUATAN — HIMPUNAN KUNCI YANG TERUKUR, BUKAN YANG DIKIRA
-- ---------------------------------------------------------------------
-- Kunci yang sungguh muncul di tiap elemen, beserta cacahnya:
--
--   CessionList    Value 2.868 · Currency 2.868 · CurrencyID 2.688 · Note 23
--   EPIList        Value 2.866 · Currency 2.812 · CurrencyID 17
--   IOOLimitList   Currency 2.868 · Value 2.868 · Layer 2.820 · CurrencyID 2.697 · Note 1.361
--   RetentionList  Value 2.869 · Currency 2.869 · Layer 2.820 · CurrencyID 2.688 · Note 868
--   Gross/NetPremiumList (Share)  Currency 3.025 · Value 3.025/3.045
--   Gross/NetPremiumList (Fac)    Currency 12 · Value 12
--
-- ⚠ `LAYER` dan `NOTE` NULLABLE dengan sengaja: keduanya muncul pada
--   sebagian larik saja (`Layer` nol di EPIList dan CessionList), dan kolom
--   yang dipaksa `NOT NULL` akan menolak baris yang sistem lama terima.
--
-- ⛔ `pxObjClass`, `pxListSubscript`, `pxCreateDateTime`, `pxCreateOpName`,
--   `pxCreateOperator`, `pxCreateSystemID` TIDAK diberi kolom. Keenamnya
--   perabot Pega, bukan muatan -- dan `PXOBJCLASS` dilarang keras di modul
--   ini sejak migrasi `436`.
--
-- ---------------------------------------------------------------------
-- TAUTAN
-- ---------------------------------------------------------------------
-- Induk ketiganya tabel MILIK MODUL INI, jadi kunci asingnya NYATA dengan
-- `ON DELETE CASCADE` -- berbeda dari tautan ke akar `TREATY_IN`, yang
-- mustahil (`ORA-02270`, lihat kepala `437`) dan karena itu tetap `MASTERID`
-- teks.
--
--   T_TREATY_LIMIT_AMOUNT      -> T_TREATY_LIMIT_DETAIL (ID)
--   T_TREATY_SHARE_AMOUNT      -> T_TREATY_SHARE        (ID)
--   T_TREATY_FAC_SHARE_AMOUNT  -> T_TREATY_FAC_SHARE    (ID)
--
-- ⚠ `URUTAN` dihitung PER INDUK, bukan global -- karena itu `UQ` nya
--   `(IDINDUK, URUTAN)`, bukan `(URUTAN)`.

-- T_TREATY_LIMIT_AMOUNT  <- Limits[].Detail[].{EPI,Retention,Cession,IOOLimit}List
--                           (11.475 elemen, 4 larik sumber)
CREATE TABLE {skema}.T_TREATY_LIMIT_AMOUNT (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  JENIS     VARCHAR2(20 CHAR)  NOT NULL,
  CURRENCY                         VARCHAR2(4000 CHAR),
  CURRENCYID                       VARCHAR2(4000 CHAR),
  LAYER                            VARCHAR2(4000 CHAR),
  NOTE                             VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_AMOUNT PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_AMOUNT UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_LIMIT_AMOUNT FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_LIMIT_DETAIL (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_AMOUNT_MST ON {skema}.T_TREATY_LIMIT_AMOUNT (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_AMOUNT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_SHARE_AMOUNT  <- Share[].{Gross,Net}PremiumList  (6.097 elemen)
CREATE TABLE {skema}.T_TREATY_SHARE_AMOUNT (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  JENIS     VARCHAR2(20 CHAR)  NOT NULL,
  CURRENCY                         VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_SHARE_AMOUNT PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_SHARE_AMOUNT UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_SHARE_AMOUNT FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_SHARE (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_SHARE_AMOUNT_MST ON {skema}.T_TREATY_SHARE_AMOUNT (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_SHARE_AMOUNT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_FAC_SHARE_AMOUNT  <- FacultativeShareList[].{Gross,Net}PremiumList
--                               (24 elemen di 5 dokumen -- yang PALING TIPIS)
CREATE TABLE {skema}.T_TREATY_FAC_SHARE_AMOUNT (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  JENIS     VARCHAR2(20 CHAR)  NOT NULL,
  CURRENCY                         VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_FAC_SHARE_AMOUNT PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_FAC_SHARE_AMOUNT UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_FAC_SHARE_AMOUNT FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_FAC_SHARE (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_FAC_SHARE_AMT_MST ON {skema}.T_TREATY_FAC_SHARE_AMOUNT (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_FAC_SHARE_AMOUNT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

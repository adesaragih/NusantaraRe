-- Tiga belas anak `TREATY_IN` dibangun — struktur dari
-- `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`, kolom dari PENGUKURAN.
--
-- Keputusan pemilik proses 6 Oktober 2026: puncaknya tetap `TREATY_IN`, dan
-- anak-anaknya mengikuti xlsx itu.
--
-- ---------------------------------------------------------------------
-- DUA HAL YANG DIPUTUSKAN DI BERKAS INI, DAN SEBABNYA
-- ---------------------------------------------------------------------
-- 1. TAUTAN KE AKAR MEMAKAI `MASTERID` TEKS, BUKAN KUNCI ASING.
--
--    xlsx menuliskan `FK TREATY_IN_ID -> TREATY_IN.ID CASCADE` pada tiap
--    anak tingkat pertama, dan menyebut `TREATY_IN   PK ID` di sel C6.
--    Terukur di POOLDATA 6 Oktober 2026: `TREATY_IN` punya NOL kunci utama,
--    NOL kunci unik, dan NOL index unik. Kunci asing ke sana ditolak Oracle
--    dengan `ORA-02270`.
--
--    Datanya SEBENARNYA unik -- 1.854 baris, nol `ID` kembar -- jadi kuncinya
--    DAPAT dipasang. Yang menghalangi aturan, bukan data: `TREATY_IN` tabel
--    WARISAN, dan nol DDL boleh dikenakan padanya sejak modul ini lahir.
--    Berkas ini karena itu memakai `MASTERID` teks, persis seperti kesembilan
--    tabel pendaratan yang sudah berdiri.
--
--    ⚠ YANG HILANG KARENANYA, dan itu dinyatakan: baris yatim tidak ditolak
--    basis data, dan `CASCADE` dari akar harus ditiru aplikasi.
--
--    Pembalikan: bila `UNIQUE (ID)` dipasang pada `TREATY_IN`, kolom
--    `TREATY_IN_ID` beserta kunci asingnya dapat menggantikan `MASTERID`.
--
-- 2. KUNCI ASING ANTAR TABEL KITA SENDIRI TETAP NYATA.
--
--    Hanya tautan ke `TREATY_IN` yang mustahil. `LIMIT_ID`, `SHARE_ID`,
--    `LIMIT_DETAIL_ID`, `LIMIT_GROUP_ID`, `FAC_SHARE_ID`, `FAC_REINSURER_ID`
--    semuanya menunjuk tabel milik modul ini, jadi keenamnya dipasang sebagai
--    kunci asing sungguhan dengan `ON DELETE CASCADE`, persis bunyi xlsx.
--
-- ---------------------------------------------------------------------
-- YANG TIDAK DIBANGUN DI SINI, DAN SEBABNYA MASING-MASING
-- ---------------------------------------------------------------------
-- ⛔ RALAT 6 Oktober 2026 atas alinea "TUJUH tabel" di bawah: pemotongan sel
--    itu NYATA, tetapi ia BUKAN penghalangnya. Jalur PERTAMA tiap sel utuh,
--    dan putusan per jalur sudah tertulis berbukti di
--    `2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md`. Keadaan sebenarnya ketujuh
--    tabel itu -- tiga dapat dibangun, lima menunggu SATU pertanyaan yang
--    belum diajukan -- ada di `docs/PEMETAAN-36-TABEL.md` §3.3. Alinea di
--    bawah ditinggalkan apa adanya sebagai catatan keadaan saat berkas ini
--    dijalankan; yang berlaku sekarang §3.3.
--
--  • TUJUH tabel yang jalur sumbernya TERPOTONG di dalam xlsx itu sendiri --
--    selnya berbunyi `| Tr`, `| TreatyIn.Limi`, terputus di tengah kata:
--      T_TREATY_LIMIT_AMOUNT   T_TREATY_LIMIT_MEASURE   T_TREATY_SHARE_AMOUNT
--      T_TREATY_TOTAL          T_TREATY_LIMIT_SUMMARY   T_TREATY_FAC_SHARE_AMOUNT
--      T_TREATY_SHARE_SPREAD_AMOUNT
--    Ketujuhnya menggabungkan BEBERAPA larik dokumen, dan daftarnya tidak
--    lengkap di mana pun. Menebak sisanya berarti mengarang kolom.
--
--  • DUA tabel yang NOL ELEMEN di seluruh 1.854 dokumen:
--      T_TREATY_FAC_LIMITS   T_TREATY_FAC_LIMIT_DETAIL
--    Jalurnya utuh dan terbaca; isinya memang tidak pernah ada. Tabel tanpa
--    satu pun kolom muatan bukan tabel.
--
--  • `T_TREATY_CURRENCY` -- xlsx menautkannya ke `MATA_UANG_KONTRAK`, yang
--    migrasi `434` CABUT atas keputusan §16 (kurs dari `TREATYEXCHANGEYEARLY`).
--
--  • `T_TREATY_VALUE_DIFFERENCE` dan `T_TREATY_VALUE_BEFORE_PRORATE` --
--    sudah PINDAH ke modul Adjustment lewat §19, migrasi `442`.
--
--  • `T_TREATY_REVISION` dan `T_TREATY_HAZARD_LIMIT` -- jalur sumbernya
--    `TreatyIn` (akar dokumen), sehingga sapuan mengembalikan ke-140 kunci
--    kepala kontrak. Itu isi `TREATY_IN` sendiri; membangunnya menggandakan
--    akar. Menunggu penegasan medan mana yang dimaksud.
--
-- ⚠ Kunci dokumen `Comment` menjadi kolom `COMMENT_` — `COMMENT` kata cadangan
--    Oracle, dan memakainya menggagalkan seluruh pernyataan dengan `ORA-00904`.
--    Ia satu-satunya yang bentrok dari seluruh kolom yang disapu.
--
-- ⛔ SELURUH kolom muatan `VARCHAR2(4000 CHAR)`. Dokumen lama menyimpan
--    semuanya sebagai teks, dan menebak tipe di sini berarti menolak nilai
--    yang sistem lama terima. Penyempitan tipe adalah pekerjaan tersendiri,
--    sesudah pemuatnya berjalan dan sebarannya terukur.
--
-- Kolomnya DITURUNKAN dari sapuan 1.854 dokumen, bukan dari daftar tangan --
-- xlsx menyatakan "KOLOM SENGAJA TIDAK DIMUAT". Cacah elemen tiap tabel
-- tercatat di komentar atasnya, dan `TURUNAN-KOLOM-28-ANAK.md` memuat
-- sebaran tiap kolomnya.

-- T_TREATY_LIMITS  <- TreatyIn.Limits  (4210 elemen, 30 kolom)
CREATE TABLE {skema}.T_TREATY_LIMITS (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  ADJRATE                          VARCHAR2(4000 CHAR),
  AGREGATELIMIT                    VARCHAR2(4000 CHAR),
  AGREGATELIMIT2                   VARCHAR2(4000 CHAR),
  COVER                            VARCHAR2(4000 CHAR),
  CURRENCY                         VARCHAR2(4000 CHAR),
  CURRENCY2                        VARCHAR2(4000 CHAR),
  CURRENCYID                       VARCHAR2(4000 CHAR),
  CURRENCYRELATION                 VARCHAR2(4000 CHAR),
  DEDUCTIBLE                       VARCHAR2(4000 CHAR),
  DEDUCTIBLE2                      VARCHAR2(4000 CHAR),
  IDX                              VARCHAR2(4000 CHAR),
  ISCOMBINEMDP                     VARCHAR2(4000 CHAR),
  LAYER                            VARCHAR2(4000 CHAR),
  LAYERPART                        VARCHAR2(4000 CHAR),
  LAYERPARTTYPE                    VARCHAR2(4000 CHAR),
  LAYERTYPE                        VARCHAR2(4000 CHAR),
  LIMIT                            VARCHAR2(4000 CHAR),
  LIMIT2                           VARCHAR2(4000 CHAR),
  MDPMINPCT                        VARCHAR2(4000 CHAR),
  MDPPCT                           VARCHAR2(4000 CHAR),
  NORIPCALCULATION                 VARCHAR2(4000 CHAR),
  ROLPCT                           VARCHAR2(4000 CHAR),
  REINSTATEMENTNOTE                VARCHAR2(4000 CHAR),
  REINSTATEMENTPCT                 VARCHAR2(4000 CHAR),
  REINSTATEMENTVALUE               VARCHAR2(4000 CHAR),
  SUMLOSSRATIO                     VARCHAR2(4000 CHAR),
  SUMTOTALACHIEVINCURED            VARCHAR2(4000 CHAR),
  SUMTOTALACHIEVNETPREMIUM         VARCHAR2(4000 CHAR),
  TREATYTYPE                       VARCHAR2(4000 CHAR),
  TREATYTYPEID                     VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMITS PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMITS UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMITS START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_SHARE  <- TreatyIn.Share  (2923 elemen, 9 kolom)
CREATE TABLE {skema}.T_TREATY_SHARE (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  COVER                            VARCHAR2(4000 CHAR),
  LAYER                            VARCHAR2(4000 CHAR),
  LAYERPART                        VARCHAR2(4000 CHAR),
  LAYERPARTTYPE                    VARCHAR2(4000 CHAR),
  LAYERTYPE                        VARCHAR2(4000 CHAR),
  RNMSHARE                         VARCHAR2(4000 CHAR),
  SPREADINGTOTALPCTXOL             VARCHAR2(4000 CHAR),
  SPREADINGTYPEIDXOL               VARCHAR2(4000 CHAR),
  SPREADINGTYPEXOL                 VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_SHARE PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_SHARE UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_SHARE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_RETRO_SHARE  <- TreatyIn.ShareReins  (9 elemen, 4 kolom)
CREATE TABLE {skema}.T_TREATY_RETRO_SHARE (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  LAYER                            VARCHAR2(4000 CHAR),
  REINSID                          VARCHAR2(4000 CHAR),
  REINSNAME                        VARCHAR2(4000 CHAR),
  SHAREPCT                         VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_RETRO_SHARE PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_RETRO_SHARE UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_RETRO_SHARE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_FAC_SHARE  <- TreatyIn.FacultativeShareList  (12 elemen, 12 kolom)
CREATE TABLE {skema}.T_TREATY_FAC_SHARE (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  COVER                            VARCHAR2(4000 CHAR),
  LAYER                            VARCHAR2(4000 CHAR),
  LAYERPART                        VARCHAR2(4000 CHAR),
  LAYERPARTTYPE                    VARCHAR2(4000 CHAR),
  LAYERTYPE                        VARCHAR2(4000 CHAR),
  LIMIT                            VARCHAR2(4000 CHAR),
  LIMIT2                           VARCHAR2(4000 CHAR),
  SPREADINGTOTALPCTXOL             VARCHAR2(4000 CHAR),
  SPREADINGTYPEIDXOL               VARCHAR2(4000 CHAR),
  SPREADINGTYPEXOL                 VARCHAR2(4000 CHAR),
  SPREADINGTYPEXOLRETRO            VARCHAR2(4000 CHAR),
  SPREADINGTYPEXOLRETROID          VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_FAC_SHARE PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_FAC_SHARE UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_FAC_SHARE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_FAC_REINSURER  <- TreatyIn.ShareFacultativeReinsurers  (15 elemen, 4 kolom)
CREATE TABLE {skema}.T_TREATY_FAC_REINSURER (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  LAYER                            VARCHAR2(4000 CHAR),
  REINSID                          VARCHAR2(4000 CHAR),
  REINSNAME                        VARCHAR2(4000 CHAR),
  SHAREPCT                         VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_FAC_REINSURER PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_FAC_REINSURER UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_FAC_REINSURER START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_LIMIT_DETAIL  <- TreatyIn.Limits.Detail  (2868 elemen, 59 kolom)
CREATE TABLE {skema}.T_TREATY_LIMIT_DETAIL (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  ACHIEVEMENTPCT                   VARCHAR2(4000 CHAR),
  BROKERAGE                        VARCHAR2(4000 CHAR),
  CASHLOSS                         VARCHAR2(4000 CHAR),
  CESSIONPCT                       VARCHAR2(4000 CHAR),
  CLAIMCOOPERATION                 VARCHAR2(4000 CHAR),
  CLASSOFBUSINESS                  VARCHAR2(4000 CHAR),
  CURRENCYCASHLOSS                 VARCHAR2(4000 CHAR),
  CURRENCYCLAIMCOOPERATION         VARCHAR2(4000 CHAR),
  CURRENCYEPI                      VARCHAR2(4000 CHAR),
  CURRENCYEARTHQUAKE               VARCHAR2(4000 CHAR),
  CURRENCYFLOODJAB                 VARCHAR2(4000 CHAR),
  CURRENCYFLOODNAT                 VARCHAR2(4000 CHAR),
  CURRENCYID                       VARCHAR2(4000 CHAR),
  CURRENCYPLA                      VARCHAR2(4000 CHAR),
  CURRENCYRSMD                     VARCHAR2(4000 CHAR),
  EPI                              VARCHAR2(4000 CHAR),
  EARTHQUAKE                       VARCHAR2(4000 CHAR),
  FLOODJAB                         VARCHAR2(4000 CHAR),
  FLOODNATION                      VARCHAR2(4000 CHAR),
  IDX                              VARCHAR2(4000 CHAR),
  IOOPCT                           VARCHAR2(4000 CHAR),
  LOSSRATIO                        VARCHAR2(4000 CHAR),
  LOWERBAND                        VARCHAR2(4000 CHAR),
  PLA                              VARCHAR2(4000 CHAR),
  PARENTID                         VARCHAR2(4000 CHAR),
  PERIODE                          VARCHAR2(4000 CHAR),
  PREMIUMRESERVEPCT                VARCHAR2(4000 CHAR),
  PROFITCOMMISION                  VARCHAR2(4000 CHAR),
  PROFITME                         VARCHAR2(4000 CHAR),
  PROFITYDCF                       VARCHAR2(4000 CHAR),
  QSPCT                            VARCHAR2(4000 CHAR),
  RIOGR                            VARCHAR2(4000 CHAR),
  RIONR                            VARCHAR2(4000 CHAR),
  RNMSHARE                         VARCHAR2(4000 CHAR),
  RSMDLIMIT                        VARCHAR2(4000 CHAR),
  REISUREDPARTICIPANT              VARCHAR2(4000 CHAR),
  RETENTIONPCT                     VARCHAR2(4000 CHAR),
  SHARENOTE                        VARCHAR2(4000 CHAR),
  SPREADINGTOTALPCT                VARCHAR2(4000 CHAR),
  SPREADINGTYPE                    VARCHAR2(4000 CHAR),
  SPREADINGTYPEID                  VARCHAR2(4000 CHAR),
  SUMLOSSRATIO                     VARCHAR2(4000 CHAR),
  SUMTOTALACHIEVINCURED            VARCHAR2(4000 CHAR),
  SUMTOTALACHIEVNETPREMIUM         VARCHAR2(4000 CHAR),
  SURPLUS                          VARCHAR2(4000 CHAR),
  TOTALACHAFTERCLAIM               VARCHAR2(4000 CHAR),
  TOTALACHCASHCALL                 VARCHAR2(4000 CHAR),
  TOTALACHESTCASHCALL              VARCHAR2(4000 CHAR),
  TOTALACHINCURED                  VARCHAR2(4000 CHAR),
  TOTALACHNETPREMIUM               VARCHAR2(4000 CHAR),
  TOTALACHOSCASHCALL               VARCHAR2(4000 CHAR),
  TOTALACHOSCLAIM                  VARCHAR2(4000 CHAR),
  TOTALACHOUTS                     VARCHAR2(4000 CHAR),
  TOTALACHPAID                     VARCHAR2(4000 CHAR),
  TOTALACHPREMIUM                  VARCHAR2(4000 CHAR),
  TREATYGROUP                      VARCHAR2(4000 CHAR),
  TREATYGROUPID                    VARCHAR2(4000 CHAR),
  TREATYTYPE                       VARCHAR2(4000 CHAR),
  UPPERBAND                        VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_DETAIL PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_DETAIL UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_LIMIT_DETAIL FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_LIMITS (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_DETAIL_MST ON {skema}.T_TREATY_LIMIT_DETAIL (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_DETAIL START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_LIMIT_GROUP  <- TreatyIn.Limits.TreatyGroupList  (8433 elemen, 5 kolom)
CREATE TABLE {skema}.T_TREATY_LIMIT_GROUP (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  CLASSOFBUSINESS                  VARCHAR2(4000 CHAR),
  CLASSOFBUSINESSID                VARCHAR2(4000 CHAR),
  ISROLPROFILE                     VARCHAR2(4000 CHAR),
  TREATYGROUP                      VARCHAR2(4000 CHAR),
  TREATYGROUPID                    VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_GROUP PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_GROUP UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_LIMIT_GROUP FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_LIMITS (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_GROUP_MST ON {skema}.T_TREATY_LIMIT_GROUP (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_GROUP START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_SHARE_SPREADING  <- TreatyIn.Share.SpreadingListXOL  (5550 elemen, 6 kolom)
CREATE TABLE {skema}.T_TREATY_SHARE_SPREADING (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  PARENTREINSTYPEID                VARCHAR2(4000 CHAR),
  PCT                              VARCHAR2(4000 CHAR),
  REINSTYPEID                      VARCHAR2(4000 CHAR),
  REINSTYPENAME                    VARCHAR2(4000 CHAR),
  RP                               VARCHAR2(4000 CHAR),
  USD                              VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_SHARE_SPREADING PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_SHARE_SPREADING UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_SHARE_SPREADING FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_SHARE (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_SHARE_SPREADING_MST ON {skema}.T_TREATY_SHARE_SPREADING (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_SHARE_SPREADING START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_SHARE_DEDUCTION  <- TreatyIn.Share.DeductionList  (2419 elemen, 5 kolom)
CREATE TABLE {skema}.T_TREATY_SHARE_DEDUCTION (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  COMMENT_                          VARCHAR2(4000 CHAR),
  CURRENCY                         VARCHAR2(4000 CHAR),
  DEDUCTION                        VARCHAR2(4000 CHAR),
  DEDUCTIONPCT                     VARCHAR2(4000 CHAR),
  DEDUCTIONPCTCALCULATE            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_SHARE_DEDUCTION PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_SHARE_DEDUCTION UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_SHARE_DEDUCTION FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_SHARE (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_SHARE_DEDUCTION_MST ON {skema}.T_TREATY_SHARE_DEDUCTION (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_SHARE_DEDUCTION START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_FAC_SHARE_DEDUCTION  <- TreatyIn.FacultativeShareList.DeductionList  (12 elemen, 5 kolom)
CREATE TABLE {skema}.T_TREATY_FAC_SHARE_DEDUCTION (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  COMMENT_                          VARCHAR2(4000 CHAR),
  CURRENCY                         VARCHAR2(4000 CHAR),
  DEDUCTION                        VARCHAR2(4000 CHAR),
  DEDUCTIONPCT                     VARCHAR2(4000 CHAR),
  DEDUCTIONPCTCALCULATE            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_FAC_SHARE_DED PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_FAC_SHARE_DED UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_FAC_SHARE_DED FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_FAC_SHARE (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_FAC_SHARE_DED_MST ON {skema}.T_TREATY_FAC_SHARE_DEDUCTION (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_FAC_SHARE_DED START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_LIMIT_COB  <- TreatyIn.Limits.Detail.COBList  (6547 elemen, 4 kolom)
CREATE TABLE {skema}.T_TREATY_LIMIT_COB (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  CLASSOFBUSINESS                  VARCHAR2(4000 CHAR),
  CLASSOFBUSINESSID                VARCHAR2(4000 CHAR),
  TREATYGROUP                      VARCHAR2(4000 CHAR),
  TREATYGROUPID                    VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_COB PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_COB UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_LIMIT_COB FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_LIMIT_DETAIL (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_COB_MST ON {skema}.T_TREATY_LIMIT_COB (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_COB START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_LIMIT_ACHIEVEMENT  <- TreatyIn.Limits.Detail.AchievementLists  (2031 elemen, 44 kolom)
CREATE TABLE {skema}.T_TREATY_LIMIT_ACHIEVEMENT (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  BROKERAGE                        VARCHAR2(4000 CHAR),
  CASHCALL                         VARCHAR2(4000 CHAR),
  CASHCALLTOIDR                    VARCHAR2(4000 CHAR),
  CONVERSION                       VARCHAR2(4000 CHAR),
  CURRENCY                         VARCHAR2(4000 CHAR),
  CURRENCYID                       VARCHAR2(4000 CHAR),
  ESTIMATIONCASHCALL               VARCHAR2(4000 CHAR),
  ESTIMATIONCASHCALLTOIDR          VARCHAR2(4000 CHAR),
  IDPEGA                           VARCHAR2(4000 CHAR),
  INCUREDCLAIM                     VARCHAR2(4000 CHAR),
  INCUREDCLAIMTOIDR                VARCHAR2(4000 CHAR),
  LOSSRATIO                        VARCHAR2(4000 CHAR),
  NETPREMIUM                       VARCHAR2(4000 CHAR),
  NETPREMIUMTOIDR                  VARCHAR2(4000 CHAR),
  NOOFFER                          VARCHAR2(4000 CHAR),
  NOPOLIS                          VARCHAR2(4000 CHAR),
  OUTSTANDINGCASHCALL              VARCHAR2(4000 CHAR),
  OUTSTANDINGCASHCALLTOIDR         VARCHAR2(4000 CHAR),
  OUTSTANDINGCLAIM                 VARCHAR2(4000 CHAR),
  OUTSTANDINGCLAIMTOIDR            VARCHAR2(4000 CHAR),
  PREMIUM                          VARCHAR2(4000 CHAR),
  PREMIUMTOIDR                     VARCHAR2(4000 CHAR),
  PAIDCLAIM                        VARCHAR2(4000 CHAR),
  PAIDCLAIMTOIDR                   VARCHAR2(4000 CHAR),
  PERIOD                           VARCHAR2(4000 CHAR),
  QUARTERYEAR                      VARCHAR2(4000 CHAR),
  QUARTER                          VARCHAR2(4000 CHAR),
  RICOMM                           VARCHAR2(4000 CHAR),
  SOBNAME                          VARCHAR2(4000 CHAR),
  TREATYGROUPNAME                  VARCHAR2(4000 CHAR),
  TREATYTYPE                       VARCHAR2(4000 CHAR),
  TOTAL                            VARCHAR2(4000 CHAR),
  TOTALACHIEVCASHCALL              VARCHAR2(4000 CHAR),
  TOTALACHIEVESTCASHCALL           VARCHAR2(4000 CHAR),
  TOTALACHIEVINCURED               VARCHAR2(4000 CHAR),
  TOTALACHIEVNETPREMIUM            VARCHAR2(4000 CHAR),
  TOTALACHIEVOSCASHCALL            VARCHAR2(4000 CHAR),
  TOTALACHIEVOSCLAIM               VARCHAR2(4000 CHAR),
  TOTALACHIEVOUTS                  VARCHAR2(4000 CHAR),
  TOTALACHIEVPAID                  VARCHAR2(4000 CHAR),
  TOTALACHIEVPREMIUM               VARCHAR2(4000 CHAR),
  TOTALAFTERCLAIM                  VARCHAR2(4000 CHAR),
  TOTALTOIDR                       VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_ACHIEVE PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_ACHIEVE UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_LIMIT_ACHIEVE FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_LIMIT_DETAIL (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_ACHIEVE_MST ON {skema}.T_TREATY_LIMIT_ACHIEVEMENT (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_ACHIEVE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- T_TREATY_LIMIT_GROUP_COB  <- TreatyIn.Limits.TreatyGroupList.ClassOfBusinessList  (15746 elemen, 4 kolom)
CREATE TABLE {skema}.T_TREATY_LIMIT_GROUP_COB (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  CLASSOFBUSINESS                  VARCHAR2(4000 CHAR),
  CLASSOFBUSINESSID                VARCHAR2(4000 CHAR),
  TREATYGROUP                      VARCHAR2(4000 CHAR),
  TREATYGROUPID                    VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_GRP_COB PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_GRP_COB UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_LIMIT_GRP_COB FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_LIMIT_GROUP (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_GRP_COB_MST ON {skema}.T_TREATY_LIMIT_GROUP_COB (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_GRP_COB START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- EMPAT tabel terakhir yang menutup larangan `JSONDATA` -- untuk Treaty In
-- DAN untuk Treaty In Adjustment.
--
-- Keputusan pemilik proses 6 Oktober 2026: "jangan ada dri jsondata lagi,
-- begitu juga treaty in adjustment". Sesudah migrasi ini, nol nilai di kedua
-- layar datang dari `M_TREATY_IN.JSONDATA` maupun `M_TREATY_IN_EDM.JSONDATA`.
--
-- Nama dan pohonnya dari `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx`,
-- kolomnya dari PENGUKURAN (`docs/TURUNAN-KOLOM-28-ANAK.md`), bentuknya
-- mengikuti `437` dan `438`.
--
-- ---------------------------------------------------------------------
-- TIGA HAL YANG DIPUTUSKAN DI BERKAS INI
-- ---------------------------------------------------------------------
-- 1. `T_TREATY_REVISION` MENAMPUNG SKALAR AKAR DOKUMEN, 1:1 per MASTERID.
--
--    xlsx menulis jalurnya `<- TreatyIn` -- akar dokumen -- sehingga sapuan
--    mengembalikan ke-134 kunci akar, larik dan skalar bercampur. Yang
--    didaratkan di sini hanya SKALAR yang kedua layar sungguh baca; larik
--    akar sudah punya tabelnya sendiri sejak `430`/`437`/`438`.
--
--    xlsx menamainya `[VERSI_KONTRAK . bagian addendum]`, dan itu cocok:
--    medan kepala sebuah kontrak memang milik versinya.
--
--    ⛔ `T_TREATY_CURRENCY` TIDAK DIBANGUN, dan itu keputusan pemilik proses
--    yang ditegaskan dua kali (4 dan 6 Oktober 2026): pengganti
--    `MATA_UANG_KONTRAK` adalah `TREATYEXCHANGEYEARLY`, tabel WARISAN yang
--    sudah hidup -- 140 baris, 25 mata uang, terbagi per `TREATYYEAR`.
--
--    Dan spesifikasi membenarkannya dengan bukti:
--    `KOREKSI-ERD-VERSUS-POOLDATA.md` §1.2 menyatakan `CurrencyList` di
--    `M_TREATY_IN.JSONDATA` adalah SALINAN tabel itu, bukan sumbernya.
--    Mendaratkan salinan sebagai kalau ia sumber adalah cara tercepat
--    melahirkan dua angka kurs yang berselisih tanpa ada yang tahu mana yang
--    benar.
--
--    `T_TREATY_HAZARD_LIMIT` TIDAK dibangun. xlsx memberinya jalur akar yang
--    SAMA, tetapi nilainya -- RSMD, gempa, banjir -- sudah berumah di
--    `T_TREATY_LIMIT_DETAIL` sejak `437`. Dua rumah untuk satu nilai berarti
--    dua tempat untuk salah, dan yang kedua selalu yang basi.
--
-- 2. SISI `OLDDATA` DIBEDAKAN LEWAT `MASTERID`, BUKAN LEWAT KOLOM BARU.
--
--    Dokumen Adjustment membawa DUA halaman: akar (sisi `New`) dan `OLDDATA`
--    (sisi `Old`). Keduanya berbentuk sama dan harus mendarat di tabel yang
--    sama -- xlsx menandai setiap kotak `BERSAMA -> Prop, Non Prop, EDM Prop,
--    EDM Non Prop`.
--
--    Pembedanya `MASTERID`: sisi `New` memakai pengenal dokumennya apa
--    adanya, sisi `Old` memakai pengenal yang sama berakhiran `#LAMA`.
--
--    ALTERNATIFNYA kolom `SISI` pada ke-24 tabel yang sudah berdiri, beserta
--    pembangunan ulang setiap `UNIQUE (MASTERID, URUTAN)`. Yang dipilih yang
--    pertama sebab ia NOL perubahan pada tabel yang sudah dijalankan, dan
--    `MASTERID` memang sudah `VARCHAR2(100)` berisi pengenal bukan-angka
--    (`<asal>/R<nn>`). Tetapannya hidup di SATU tempat di Go
--    (`repository.AkhiranSisiLama`), dipakai pemuat dan pembaca bersama.
--
--    ⚠ YANG HILANG KARENANYA, dan itu dinyatakan: basis data tidak dapat
--    membedakan kedua sisi tanpa membaca akhiran teks. Kueri ad-hoc yang
--    lupa menyaringnya akan melihat baris kedua sisi sekaligus.
--
-- 3. `JENIS` PADA DUA TABEL, mengikuti keputusan migrasi `438`.
--
--    `T_TREATY_LIMIT_MEASURE` menampung tiga larik (`MDPList`,
--    `PremiumEarnedList`, `EgnpiTotalList`) dan `T_TREATY_TOTAL` menampung
--    enam. Tanpa `JENIS` barisnya tidak dapat dibedakan sesudah mendarat.
--    Kolom itu TIDAK ADA di dokumen; ia dipaksa oleh bentuk xlsx.

CREATE TABLE {skema}.T_TREATY_REVISION (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  IDX                              VARCHAR2(4000 CHAR),
  OLDID                            VARCHAR2(4000 CHAR),
  PROPORTIONTYPE                   VARCHAR2(4000 CHAR),
  TREATYCONTRACTNAME               VARCHAR2(4000 CHAR),
  TERITORIALSCOPE                  VARCHAR2(4000 CHAR),
  TREATYYEAR                       VARCHAR2(4000 CHAR),
  CEDING                           VARCHAR2(4000 CHAR),
  LEADINGREINSSOURCE               VARCHAR2(4000 CHAR),
  COMMENCEMENT                     VARCHAR2(4000 CHAR),
  TERMINATION                      VARCHAR2(4000 CHAR),
  CONTRACTREFNO                    VARCHAR2(4000 CHAR),
  BORDEAUX                         VARCHAR2(4000 CHAR),
  BORDEREAUXNOTE                   VARCHAR2(4000 CHAR),
  ACCOUNTINGMODE                   VARCHAR2(4000 CHAR),
  ACCOUNTINGMODENONPROP            VARCHAR2(4000 CHAR),
  TREATYLEADER                     VARCHAR2(4000 CHAR),
  ISMULTIPLERETRO                  VARCHAR2(4000 CHAR),
  EDMSTATE                         VARCHAR2(4000 CHAR),
  EDMMATERIALTYPE                  VARCHAR2(4000 CHAR),
  EDMEFFECTIVE                     VARCHAR2(4000 CHAR),
  STATUSAKSEPTASI                  VARCHAR2(4000 CHAR),
  ISPRORATE                        VARCHAR2(4000 CHAR),
  PRORATEDAYS                      VARCHAR2(4000 CHAR),
  PRORATETOTALDAYS                 VARCHAR2(4000 CHAR),
  PRORATEPERCENT                   VARCHAR2(4000 CHAR),
  ISEDITDATA                       VARCHAR2(4000 CHAR),
  TOTALEGNPIAMOUNT                 VARCHAR2(4000 CHAR),
  TOTALEGNPIPROPORTION             VARCHAR2(4000 CHAR),
  TOTALLIMITSROL                   VARCHAR2(4000 CHAR),
  FACULTATIVESHARE                 VARCHAR2(4000 CHAR),
  FACULTATIVESHAREBROKERAGE        VARCHAR2(4000 CHAR),
  VALUEDIFF_RNMSHARE               VARCHAR2(4000 CHAR),
  VALUEDIFF_BROKERAGEPCT           VARCHAR2(4000 CHAR),
  EXCLUSIONS                       CLOB,
  EXCLUSIONSP                      CLOB,
  SPECIALCONDITIONS                CLOB,
  SPECIALCONDITIONSP               CLOB,
  SPECIALCONDITIONSLC              CLOB,
  CONSTRAINT PK_TT_REVISION PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_REVISION UNIQUE (MASTERID, URUTAN)
)
/
-- ⛔ LIMA kolom CLOB, dan itu diukur bukan dijaga-jaga: teks tab
-- `Special Conditions` terpanjang 23.453 aksara, sementara `VARCHAR2` Oracle
-- berhenti di 4.000 bita. Kolom yang terlalu pendek MEMOTONG tanpa bersuara.
--
-- ⚠ `SPECIALCONDITIONSLC` memuat ejaan `SpecialConditionsp` -- huruf `p`
-- KECIL. Ia ejaan KETIGA yang sungguh ada di 292 dokumen, dan isinya BERBEDA
-- dari `SpecialConditionsP`. Oracle tidak membedakan besar-kecil huruf pada
-- nama kolom, jadi pembedanya dipindah ke akhiran `LC` (lower case).
CREATE SEQUENCE {skema}.SEQ_TT_REVISION START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- `T_TREATY_LIMIT_MEASURE` <- Limits[].MDPList | PremiumEarnedList |
-- EgnpiTotalList. Induknya `T_TREATY_LIMITS`, bukan `_DETAIL`: ketiga larik
-- ini hidup di elemen `Limits[]`, satu tingkat di atas `Detail[]`.
CREATE TABLE {skema}.T_TREATY_LIMIT_MEASURE (
  ID        NUMBER(19)         NOT NULL,
  IDINDUK   NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  JENIS     VARCHAR2(20 CHAR)  NOT NULL,
  CURRENCY                         VARCHAR2(4000 CHAR),
  CURRENCYID                       VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_MEASURE PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_MEASURE UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_TT_LIMIT_MEASURE FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.T_TREATY_LIMITS (ID) ON DELETE CASCADE
)
/
CREATE INDEX {skema}.IX_TT_LIMIT_MEASURE_MST ON {skema}.T_TREATY_LIMIT_MEASURE (MASTERID)
/
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_MEASURE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- `T_TREATY_LIMIT_SUMMARY` <- TreatyIn.LimitSummaryList, 2.855 elemen
-- terukur. Larik AKAR, jadi ia menaut ke kontrak lewat `MASTERID`.
CREATE TABLE {skema}.T_TREATY_LIMIT_SUMMARY (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  LAYER                            VARCHAR2(4000 CHAR),
  LAYERPART                        VARCHAR2(4000 CHAR),
  LAYERPARTTYPE                    VARCHAR2(4000 CHAR),
  LAYERTYPE                        VARCHAR2(4000 CHAR),
  CURRENCY                         VARCHAR2(4000 CHAR),
  CURRENCY2                        VARCHAR2(4000 CHAR),
  LIMITVAL                         VARCHAR2(4000 CHAR),
  LIMITVAL2                        VARCHAR2(4000 CHAR),
  AGGREGATELIMIT                   VARCHAR2(4000 CHAR),
  AGGREGATELIMIT2                  VARCHAR2(4000 CHAR),
  DEDUCTIBLE                       VARCHAR2(4000 CHAR),
  DEDUCTIBLE2                      VARCHAR2(4000 CHAR),
  MDP                              VARCHAR2(4000 CHAR),
  MDP2                             VARCHAR2(4000 CHAR),
  NOTE                             VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_LIMIT_SUMMARY PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_LIMIT_SUMMARY UNIQUE (MASTERID, URUTAN)
)
/
-- ⚠ `LIMITVAL`/`LIMITVAL2`, bukan `LIMIT`/`LIMIT2`. `T_TREATY_LIMITS` sudah
-- memakai `LIMIT` dan Oracle menerimanya, tetapi di tabel ini ia berdampingan
-- dengan `AGGREGATELIMIT` -- dan dua kolom yang hanya berbeda awalan membuat
-- `SELECT LIMIT` di tabel yang keliru lolos diam-diam.
CREATE SEQUENCE {skema}.SEQ_TT_LIMIT_SUMMARY START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

-- `T_TREATY_TOTAL` <- enam larik total akar yang layar Adjustment baca:
-- TotalRetentionAmountNP, TotalEgnpiAmountNP, TotalLimitIOONP,
-- TotalLimitDeductblNP, TotalLimitPremiEarnNP, TotalLimitMDPNP.
-- Keenamnya berbentuk sama (`Currency` + `Value`); `JENIS` yang membedakan.
CREATE TABLE {skema}.T_TREATY_TOTAL (
  ID        NUMBER(19)         NOT NULL,
  MASTERID  VARCHAR2(100 CHAR) NOT NULL,
  URUTAN    NUMBER(10)         NOT NULL,
  JENIS     VARCHAR2(40 CHAR)  NOT NULL,
  CURRENCY                         VARCHAR2(4000 CHAR),
  CURRENCYID                       VARCHAR2(4000 CHAR),
  VALUE                            VARCHAR2(4000 CHAR),
  CONSTRAINT PK_TT_TOTAL PRIMARY KEY (ID),
  CONSTRAINT UQ_TT_TOTAL UNIQUE (MASTERID, URUTAN)
)
/
CREATE SEQUENCE {skema}.SEQ_TT_TOTAL START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

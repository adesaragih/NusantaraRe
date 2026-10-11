-- Delapan tabel tab Treaty In - pendaratan larik `M_TREATY_IN.JSONDATA`
--
-- Migrasi tiket tab Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- ⛔ BUKAN tabel model baru. Kedelapan tabel di bawah MENDARATKAN larik yang
-- hari ini hidup di dalam dokumen `POOLDATA.M_TREATY_IN.JSONDATA`, apa adanya
-- dan sebagai TEKS, supaya layar berhenti mengurai CLOB pada setiap pembacaan.
-- Model baru (`PERIODE_PELAPORAN`, `PORTOFOLIO`, `EGNPI`, `RETENSI_CEDANT`,
-- `TERMIN`, `RINCIAN_ANGSURAN`, `PERIODE_AKUMULASI`) sudah ada sejak migrasi
-- 403-424 dan TIDAK diganti oleh berkas ini; keduanya hidup berdampingan
-- sampai pemindahan tiket 44 selesai. Polanya sama dengan `MIGRASI_PENDARATAN`
-- (migrasi 428).
--
-- =====================================================================
-- UKURAN - SELURUH 1.854 DOKUMEN, BUKAN CONTOH
-- =====================================================================
--   Disapu 3 Oktober 2026 dengan mengurai ke-1.854 dokumen secara utuh
--   (nol dokumen gagal urai). Cacah BARIS yang akan masuk tiap tabel:
--
--     M_TREATYIN_REPORTINGPERIOD   4.548 baris dari 1.137 kontrak
--     M_TREATYIN_PORTFOLIO         1.925 baris dari   845 kontrak
--     M_TREATYIN_ACCUMULATION         60 baris dari    18 kontrak
--     M_TREATYIN_EGNPI             2.298 baris dari   846 kontrak
--     M_TREATYIN_RETENTION         2.511 baris dari   808 kontrak
--     M_TREATYIN_INSTALLMENT         796 baris dari   768 kontrak
--     M_TREATYIN_INSTALLMENTITEM   3.033 baris dari   768 kontrak
--     M_TREATYIN_COMMENT          11.365 baris dari 1.837 kontrak
--                                 ------
--                                 26.536 baris
--
--   ⚠️ Angka EGNPI MENGOREKSI sapuan sebelumnya di ronde ini sendiri, yang
--   melaporkan 155. Sapuan itu mencari teks `"EGNPI":[ ]` di dalam CLOB, dan
--   kunci `EGNPI` muncul BERKALI-KALI di satu dokumen - juga bersarang. Yang
--   kosong di kedalaman membuat yang berisi di permukaan ikut terhitung
--   kosong. Pengurai JSON tidak punya kebingungan itu; 846 yang berlaku.
--
-- =====================================================================
-- ⛔ NOL KUNCI ASING KE `TREATY_IN` - DAN ITU BUKAN PILIHAN
-- =====================================================================
--   Rancangannya menyebut `MASTERID` sebagai kunci asing ke `TREATY_IN.ID`.
--   Itu TIDAK DAPAT dibuat. Diukur 3 Oktober 2026 di POOLDATA:
--
--     all_constraints  TREATY_IN  tipe P atau U  -> NOL BARIS
--     all_indexes      INDEX_ID (ID)             -> NONUNIQUE
--     all_tab_columns  ID VARCHAR2(100)          -> NULLABLE = Y
--
--   Oracle menolak `REFERENCES TREATY_IN (ID)` dengan ORA-02270 selama kolom
--   itu tidak memimpin kunci utama maupun UNIQUE. Memberinya satu berarti DDL
--   terhadap tabel warisan, dan ronde ini melarangnya dengan kalimat yang
--   tidak bercabang: nol tulisan ke tabel warisan mana pun.
--
--   Jadi `MASTERID` adalah KOLOM BERINDEX, bukan kunci asing - persis seperti
--   `MIGRASI_KORELASI.ID_KONTRAK_BARU` dan `MIGRASI_PENDARATAN.KUNCI_WARISAN`,
--   yang memilih hal yang sama dengan sebab yang berbeda. Ongkosnya dibayar
--   dan dinyatakan: basis data TIDAK menjaga keterhubungan ke kontraknya;
--   yang menjaganya pemuat dan uji rekonsiliasinya.
--
--   Syarat pembalikan tercatat di `docs/KEPUTUSAN-PENYELARASAN-REPO.md` §12.
--
-- =====================================================================
-- PERILAKU HAPUS `ikut hapus`, DIWUJUDKAN DI DUA TEMPAT YANG BERBEDA
-- =====================================================================
--   Antara `M_TREATYIN_INSTALLMENT` dan anaknya `M_TREATYIN_INSTALLMENTITEM`
--   keduanya tabel BARU, jadi di sana `ikut hapus` adalah `ON DELETE CASCADE`
--   sungguhan - satu-satunya kunci asing di berkas ini. Dari kontrak ke
--   kedelapan tabel, `ikut hapus` dijalankan oleh pemuat, sebab tidak ada
--   kunci asing yang dapat menjalankannya (lihat di atas).
--
-- =====================================================================
-- SELURUH KOLOM ISI BERTIPE TEKS, DAN ITU DISENGAJA
-- =====================================================================
--   Setiap nilai di dalam `JSONDATA` adalah string JSON - termasuk yang
--   terlihat seperti angka dan tanggal. Menafsirkannya saat memuat berarti
--   memutuskan hal yang belum diputuskan, dan dua ukuran menunjukkan mahalnya:
--
--     `Installment.AmountTotal` terpanjang 61 aksara
--        -> `1323411750.0000008394305684816601000000000000000000000000000`
--        NUMBER(38,8) akan MEMBULATKANNYA, diam-diam.
--     `Retention.Currency` terpanjang 9 aksara, isinya `1/04/2023`
--        -> kolom mata uang yang memuat tanggal. Nyata, dan ia harus
--        tersimpan apa adanya supaya dapat ditemukan dan diperbaiki.
--
--   Lebar kolom di bawah = panjang TERUKUR dengan kelonggaran. Panjang
--   terukur tiap kolom didaftar per tabel di
--   `docs/STRUKTUR-TABEL-TREATY-IN.md`, supaya kolom yang suatu hari
--   kesempitan dapat ditelusuri ke angka yang menentukannya.
--
-- `Date` -> `TANGGAL`: `DATE` kata cadangan Oracle (ORA-00923),
-- `TestNolKataCadanganOracleSebagaiKolom`. Kesepuluh nama kolom lain yang
-- berisiko - TYPE, PERIOD, NOTE, AMOUNT, INSTALLMENT, DESCRIPTION, CURRENCY,
-- SUGGEST, PROPORTION, URUTAN - diuji satu per satu dengan
-- `SELECT 1 AS <nama> FROM DUAL` di Oracle DEV, dan kesepuluhnya lolos.
--
-- INV-01 kunci utama. INV-02 pengenal dari sequence (migrasi 431).
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
--
-- ⛔ `AchievementLists` TIDAK ADA DI SINI, dan sebabnya diukur. Sebagai kunci
-- puncak ia muncul NOL kali di ke-1.854 dokumen. Jalur sebenarnya
-- `Limits[].Detail[].AchievementLists` - 1.961 kemunculan, 859 di antaranya
-- kosong - dan tiga lagi di `RevisionHistory[].Limits[].Detail[]`. Butirnya
-- milik satu baris `Limits[].Detail[]`, bukan milik kontrak, jadi tabel
-- berinduk `MASTERID` saja TIDAK DAPAT menyatakan ia milik siapa. Induknya ada
-- di `M_TREATY_IN2`, tabel warisan yang lain. Uraiannya di
-- `docs/STRUKTUR-TABEL-TREATY-IN.md`.
--
-- ⛔ UNIQUE (MASTERID, URUTAN) adalah invarian PEMUAT, bukan invarian bisnis,
-- jadi ia tidak menagih nomor di `SPEC-INVARIAN.md` (aturan
-- `Z00_KUNCI_ALAMI.sql` berlaku untuk yang kedua). Ia ada sebab ronde ini
-- menuntut pemuat yang IDEMPOTEN, dan idempotensi yang hanya dijaga oleh kode
-- pemanggil benar sampai dua pemuat berjalan bersamaan. Kunci ini juga
-- melayani index untuk `MASTERID`, sehingga tujuh tabel tidak perlu index
-- kedua.
CREATE TABLE {skema}.M_TREATYIN_REPORTINGPERIOD (
  ID               NUMBER(19)           NOT NULL,
  MASTERID         VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN           NUMBER(10)           NOT NULL,
  AUTOCALCULATE    VARCHAR2(50 CHAR),
  CONFIRMATIONDUE  VARCHAR2(50 CHAR),
  INITIALDATE      VARCHAR2(50 CHAR),
  PERIOD           VARCHAR2(50 CHAR),
  SETTLEMENTDUE    VARCHAR2(50 CHAR),
  SUBMISSIONDUE    VARCHAR2(50 CHAR),
  PXOBJCLASS       VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_REPORTINGPERIOD PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_REPORTINGPERIOD UNIQUE (MASTERID, URUTAN)
)
/
CREATE TABLE {skema}.M_TREATYIN_PORTFOLIO (
  ID             NUMBER(19)           NOT NULL,
  MASTERID       VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN         NUMBER(10)           NOT NULL,
  DESCRIPTION    VARCHAR2(4000 CHAR),
  TYPE           VARCHAR2(100 CHAR),
  TYPEPORTFOLIO  VARCHAR2(100 CHAR),
  PXOBJCLASS     VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_PORTFOLIO PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_PORTFOLIO UNIQUE (MASTERID, URUTAN)
)
/
CREATE TABLE {skema}.M_TREATYIN_ACCUMULATION (
  ID          NUMBER(19)           NOT NULL,
  MASTERID    VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN      NUMBER(10)           NOT NULL,
  PERIOD      VARCHAR2(50 CHAR),
  REPORTDATE  VARCHAR2(50 CHAR),
  SUBDAYS     VARCHAR2(50 CHAR),
  SUBDUEDATE  VARCHAR2(50 CHAR),
  PXOBJCLASS  VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_ACCUMULATION PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_ACCUMULATION UNIQUE (MASTERID, URUTAN)
)
/
CREATE TABLE {skema}.M_TREATYIN_EGNPI (
  ID                        NUMBER(19)           NOT NULL,
  MASTERID                  VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN                    NUMBER(10)           NOT NULL,
  AMOUNT                    VARCHAR2(100 CHAR),
  AMOUNTIDR                 VARCHAR2(100 CHAR),
  ASDATE                    VARCHAR2(50 CHAR),
  CLASSOFBUSINESS           VARCHAR2(200 CHAR),
  CURRENCY                  VARCHAR2(50 CHAR),
  CURRENCYID                VARCHAR2(50 CHAR),
  NOTE                      VARCHAR2(4000 CHAR),
  PROPORTION                VARCHAR2(100 CHAR),
  TREATYGROUP               VARCHAR2(200 CHAR),
  TREATYGROUPID             VARCHAR2(50 CHAR),
  PYTEMPLATERICHTEXTEDITOR  VARCHAR2(50 CHAR),
  PXOBJCLASS                VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_EGNPI PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_EGNPI UNIQUE (MASTERID, URUTAN)
)
/
CREATE TABLE {skema}.M_TREATYIN_RETENTION (
  ID               NUMBER(19)           NOT NULL,
  MASTERID         VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN           NUMBER(10)           NOT NULL,
  AMOUNT           VARCHAR2(100 CHAR),
  CLASSOFBUSINESS  VARCHAR2(200 CHAR),
  CURRENCY         VARCHAR2(50 CHAR),
  CURRENCYID       VARCHAR2(50 CHAR),
  NOTE             VARCHAR2(4000 CHAR),
  TREATYGROUP      VARCHAR2(200 CHAR),
  TREATYGROUPID    VARCHAR2(50 CHAR),
  PXOBJCLASS       VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_RETENTION PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_RETENTION UNIQUE (MASTERID, URUTAN)
)
/
CREATE TABLE {skema}.M_TREATYIN_INSTALLMENT (
  ID               NUMBER(19)           NOT NULL,
  MASTERID         VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN           NUMBER(10)           NOT NULL,
  AMOUNTTOTAL      VARCHAR2(200 CHAR),
  CURRENCY         VARCHAR2(50 CHAR),
  PCTTOTAL         VARCHAR2(100 CHAR),
  PXLISTSUBSCRIPT  VARCHAR2(50 CHAR),
  PXOBJCLASS       VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_INSTALLMENT PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_INSTALLMENT UNIQUE (MASTERID, URUTAN)
)
/
CREATE TABLE {skema}.M_TREATYIN_INSTALLMENTITEM (
  ID              NUMBER(19)           NOT NULL,
  IDINDUK         NUMBER(19)           NOT NULL,
  MASTERID        VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN          NUMBER(10)           NOT NULL,
  AMOUNT          VARCHAR2(200 CHAR),
  CURRENCY        VARCHAR2(50 CHAR),
  DUEDATE         VARCHAR2(50 CHAR),
  INSTALLMENT     VARCHAR2(50 CHAR),
  INSTALLMENTPCT  VARCHAR2(100 CHAR),
  PAYMENTDATE     VARCHAR2(50 CHAR),
  WPC             VARCHAR2(50 CHAR),
  PXOBJCLASS      VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_INSTALLMENTITEM PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_INSTALLMENTITEM UNIQUE (IDINDUK, URUTAN),
  CONSTRAINT FK_MTI_INSTALLMENTITEM_1 FOREIGN KEY (IDINDUK)
    REFERENCES {skema}.M_TREATYIN_INSTALLMENT (ID) ON DELETE CASCADE
)
/
CREATE TABLE {skema}.M_TREATYIN_COMMENT (
  ID            NUMBER(19)           NOT NULL,
  MASTERID      VARCHAR2(100 CHAR)   NOT NULL,
  URUTAN        NUMBER(10)           NOT NULL,
  TANGGAL       VARCHAR2(50 CHAR),
  CONVERTDATE   VARCHAR2(50 CHAR),
  HASHISTORY    VARCHAR2(50 CHAR),
  ISAPPROVED    VARCHAR2(100 CHAR),
  OPERATORNAME  VARCHAR2(200 CHAR),
  SUGGEST       VARCHAR2(4000 CHAR),
  PXOBJCLASS    VARCHAR2(200 CHAR),
  CONSTRAINT PK_MTI_COMMENT PRIMARY KEY (ID),
  CONSTRAINT UQ_MTI_COMMENT UNIQUE (MASTERID, URUTAN)
)
/
CREATE INDEX {skema}.IX_MTI_INSTALLMENTITEM_MST ON {skema}.M_TREATYIN_INSTALLMENTITEM (MASTERID)
/

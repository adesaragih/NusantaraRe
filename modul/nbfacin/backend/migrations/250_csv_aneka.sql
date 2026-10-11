-- 250 - unggah CSV keluarga ANEKA, mulai Marine Hull (tiket 73 P6, keputusan work owner 11-10-2026 K73-52 opsi 1) atas infrastruktur
-- tiket 75 / 74 P6. Rentang KETIGA nbfacin 250-254 (keputusan work owner 11-10-2026, diambil dari rnwfacin yang belum memakai
-- satu nomor pun). Ditulis, TIDAK dijalankan agent; work owner / DBA yang menjalankan. Nol COMMIT (ADR-U-0029).
--
-- T_NB_CSV_BARIS_ANEKA - penampung baris pekerjaan berjenis 'aneka' (InsertUploadAneka_act / ProtectUploadAneka_Act, 35 judul dibaca
-- per NAMA): satu kolom VARCHAR2(4000) per judul, awalan C_ (pola T_NB_CSV_BARIS 257 / 258). TABEL SENDIRI, bukan ALTER
-- T_NB_CSV_BARIS: pelari migrasi mengurutkan menurut NAMA berkas, sehingga 250 berjalan SEBELUM 257 yang membuat T_NB_CSV_BARIS /
-- T_NB_CSV_UNGGAH - karena itu juga TANPA FK ke T_NB_CSV_UNGGAH (pembersihan lewat kode). T_NB_CSV_UNGGAH.JENIS VARCHAR2(10) tanpa
-- CHECK sudah menampung 'aneka' (tanpa perubahan).
CREATE TABLE {skema}.T_NB_CSV_BARIS_ANEKA (
  UNGGAH_ID              NUMBER(19) NOT NULL,
  NO_BARIS               NUMBER(10) NOT NULL,
  C_NO                   VARCHAR2(4000),
  C_TYPE_ADDRESS         VARCHAR2(4000),
  C_ADDRESS              VARCHAR2(4000),
  C_BUILDING_NO          VARCHAR2(4000),
  C_POSTAL_CODE          VARCHAR2(4000),
  C_COUNTRY              VARCHAR2(4000),
  C_OCCUPATION_NAME      VARCHAR2(4000),
  C_OBJECT_DESCRIPTION   VARCHAR2(4000),
  C_SECTION              VARCHAR2(4000),
  C_CURRENCY             VARCHAR2(4000),
  C_TSI                  VARCHAR2(4000),
  C_YEAR                 VARCHAR2(4000),
  C_PLACE_OF_MANUFACTURE VARCHAR2(4000),
  C_OBJECT_TYPE          VARCHAR2(4000),
  C_NRT                  VARCHAR2(4000),
  C_GRT                  VARCHAR2(4000),
  C_DWT                  VARCHAR2(4000),
  C_FLAGNAME             VARCHAR2(4000),
  C_CONSTRUCTION         VARCHAR2(4000),
  C_CLASSIFICATION       VARCHAR2(4000),
  C_COVERAGE_NAME        VARCHAR2(4000),
  C_CALCULATE_METHOD     VARCHAR2(4000),
  C_DAY                  VARCHAR2(4000),
  C_PCT_LOADING          VARCHAR2(4000),
  C_PCT_SCALE            VARCHAR2(4000),
  C_PCT_RATE             VARCHAR2(4000),
  C_PREMIUM              VARCHAR2(4000),
  C_PCT_DISCOUNT         VARCHAR2(4000),
  C_DISCOUNT             VARCHAR2(4000),
  C_PCT_LOL              VARCHAR2(4000),
  C_LOL                  VARCHAR2(4000),
  C_SERIALNO             VARCHAR2(4000),
  C_CHASSISNUMBER        VARCHAR2(4000),
  C_ENGINENUMBER         VARCHAR2(4000),
  C_KELURAHAN            VARCHAR2(4000),
  CONSTRAINT PK_T_NB_CSV_BARIS_ANEKA PRIMARY KEY (UNGGAH_ID, NO_BARIS)
)
/

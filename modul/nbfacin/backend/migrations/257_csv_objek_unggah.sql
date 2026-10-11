-- 257 - unggah CSV objek TANPA BATAS baris (tiket 75 F1 + F2, keputusan work owner 10-10-2026 K75-1 / K75-2). Rentang KEDUA
-- nbfacin 255-259. Ditulis, TIDAK dijalankan agent; work owner / DBA yang menjalankan. Nol COMMIT (ADR-U-0029).
--
-- T_NB_CSV_UNGGAH - satu pekerjaan unggah per berkas: status memvalidasi / valid / gagal_validasi / menyimpan / selesai / gagal,
-- FASE validasi / simpan (Ulangi melanjutkan fase yang gagal), kemajuan (JUMLAH_BARIS, BARIS_DIPROSES, BARIS_TERSIMPAN,
-- KELOMPOK_TERAKHIR = kelompok simpan terakhir yang COMMIT), header polis langkah 5.3 (NO_OFFER_SLIP / BEGIN_DATE / END_DATE teks CSV,
-- baris terakhir yang terisi menang - D1 (b) tiket 66), galat proses (GALAT_PROSES, GALAT_DI_BARIS). DIUBAH = detak pekerjaan:
-- pekerjaan berjalan tanpa detak lama dianggap terputus (gagal, dapat diulang). CASE_ID tanpa FK (case tanpa General tetap boleh
-- mengunggah; General dipastikan saat simpan, pola T_LOCATIONLIST).
-- T_NB_CSV_BARIS - penampung baris CSV apa adanya: satu kolom VARCHAR2(4000) per kolom template (KolomCsvObjek, 61 kolom, awalan
-- C_), urut NO_BARIS 1..n; dibaca per kelompok (validasi / simpan) dan per halaman (View Upload).
-- T_NB_CSV_PESAN - galat (JENIS 'galat') dan pesan (JENIS 'pesan', mis. "KodePos … tidak ditemukan!") per baris, urut URUT.
-- Penampung TIDAK dihapus otomatis (pembersihan pekerjaan lama = usulan, menunggu keputusan work owner).
-- T_LOCATIONLIST.SEQ_NO NUMBER(5) -> NUMBER(10): urutan lokasi satu case (186 L18) - NUMBER(5) memuat paling banyak 99999 lokasi,
-- 1 juta baris CSV = 1 juta lokasi (K75 temuan c3). Memperlebar presisi NUMBER aman walau berisi data.
CREATE TABLE {skema}.T_NB_CSV_UNGGAH (
  ID                NUMBER(19) NOT NULL,
  CASE_ID           VARCHAR2(32) NOT NULL,
  STATUS            VARCHAR2(20) NOT NULL,
  FASE              VARCHAR2(20) NOT NULL,
  NAMA_BERKAS       VARCHAR2(500),
  UKURAN            NUMBER(19),
  JUMLAH_BARIS       NUMBER(10) DEFAULT 0 NOT NULL,
  BARIS_DIPROSES    NUMBER(10) DEFAULT 0 NOT NULL,
  BARIS_TERSIMPAN   NUMBER(10) DEFAULT 0 NOT NULL,
  KELOMPOK_TERAKHIR NUMBER(10) DEFAULT 0 NOT NULL,
  JUMLAH_GALAT      NUMBER(10) DEFAULT 0 NOT NULL,
  JUMLAH_PESAN      NUMBER(10) DEFAULT 0 NOT NULL,
  NO_OFFER_SLIP     VARCHAR2(4000),
  BEGIN_DATE        VARCHAR2(4000),
  END_DATE          VARCHAR2(4000),
  GALAT_PROSES      VARCHAR2(4000),
  GALAT_DI_BARIS    NUMBER(10),
  DIBUAT_OLEH       VARCHAR2(100),
  DIBUAT            TIMESTAMP NOT NULL,
  DIUBAH            TIMESTAMP NOT NULL,
  CONSTRAINT PK_T_NB_CSV_UNGGAH PRIMARY KEY (ID)
)
/
CREATE INDEX {skema}.IX_T_NB_CSV_UNGGAH_CASE ON {skema}.T_NB_CSV_UNGGAH (CASE_ID, DIBUAT)
/
CREATE SEQUENCE {skema}.SEQ_T_NB_CSV_UNGGAH START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.T_NB_CSV_BARIS (
  UNGGAH_ID               NUMBER(19) NOT NULL,
  NO_BARIS                NUMBER(10) NOT NULL,
  C_NO                    VARCHAR2(4000),
  C_NO_OFFER_SLIP         VARCHAR2(4000),
  C_CLASS_OF_BUSINESS     VARCHAR2(4000),
  C_BEGIN_DATE            VARCHAR2(4000),
  C_END_DATE              VARCHAR2(4000),
  C_INSURED_NAME          VARCHAR2(4000),
  C_SOB_NAME              VARCHAR2(4000),
  C_CEDING_CO_NAME        VARCHAR2(4000),
  C_MARKETING_NAME        VARCHAR2(4000),
  C_OCCUPATION_CODE       VARCHAR2(4000),
  C_OCCUPATION_NAME       VARCHAR2(4000),
  C_CLASS_OF_CONSTRUCTION VARCHAR2(4000),
  C_TYPE_ADDRESS          VARCHAR2(4000),
  C_ADDRESS               VARCHAR2(4000),
  C_BUILDING_NO           VARCHAR2(4000),
  C_POSTAL_CODE           VARCHAR2(4000),
  C_KELURAHAN             VARCHAR2(4000),
  C_OBJECT_TYPE           VARCHAR2(4000),
  C_INTEREST              VARCHAR2(4000),
  C_TOP_RISK              VARCHAR2(4000),
  C_ITEM_TYPE             VARCHAR2(4000),
  C_ITEM_TYPE_1           VARCHAR2(4000),
  C_ITEM_TYPE_2           VARCHAR2(4000),
  C_ITEM_TYPE_3           VARCHAR2(4000),
  C_ITEM_TYPE_4           VARCHAR2(4000),
  C_ITEM_TYPE_5           VARCHAR2(4000),
  C_ITEM_TYPE_6           VARCHAR2(4000),
  C_TSI_OBJECT_ITEM       VARCHAR2(4000),
  C_TSI_OBJECT_ITEM_1     VARCHAR2(4000),
  C_TSI_OBJECT_ITEM_2     VARCHAR2(4000),
  C_TSI_OBJECT_ITEM_3     VARCHAR2(4000),
  C_TSI_OBJECT_ITEM_4     VARCHAR2(4000),
  C_TSI_OBJECT_ITEM_5     VARCHAR2(4000),
  C_TSI_OBJECT_ITEM_6     VARCHAR2(4000),
  C_ACCUMULATION_CODE     VARCHAR2(4000),
  C_ACCUMULATION_NAME     VARCHAR2(4000),
  C_ADJUSTABLE_PCT        VARCHAR2(4000),
  C_CURRENCY              VARCHAR2(4000),
  C_INDEMNITY_MONTH       VARCHAR2(4000),
  C_LOL                   VARCHAR2(4000),
  C_RATE_FLEXAS           VARCHAR2(4000),
  C_NET_RATE_FLEXAS       VARCHAR2(4000),
  C_DISCOUNT_RATE_FLEXAS  VARCHAR2(4000),
  C_RATE_RSMDCC           VARCHAR2(4000),
  C_NET_RATE_RSMDCC       VARCHAR2(4000),
  C_DISCOUNT_RATE_RSMDCC  VARCHAR2(4000),
  C_RATE_FWTWD            VARCHAR2(4000),
  C_NET_RATE_FWTWD        VARCHAR2(4000),
  C_DISCOUNT_RATE_FWTWD   VARCHAR2(4000),
  C_RATE_EQ               VARCHAR2(4000),
  C_NET_RATE_EQ           VARCHAR2(4000),
  C_DISCOUNT_RATE_EQ      VARCHAR2(4000),
  C_RATE_OTHERS           VARCHAR2(4000),
  C_NET_RATE_OTHERS       VARCHAR2(4000),
  C_DISCOUNT_RATE_OTHERS  VARCHAR2(4000),
  C_KODE_POS              VARCHAR2(4000),
  C_KECAMATAN             VARCHAR2(4000),
  C_KABUPATEN_KOTA        VARCHAR2(4000),
  C_PROVINSI              VARCHAR2(4000),
  C_SUMBER_KODE_POS       VARCHAR2(4000),
  C_TINGKAT_KEYAKINAN     VARCHAR2(4000),
  CONSTRAINT PK_T_NB_CSV_BARIS PRIMARY KEY (UNGGAH_ID, NO_BARIS),
  CONSTRAINT FK_NB_CSV_BARIS_UNGGAH FOREIGN KEY (UNGGAH_ID) REFERENCES {skema}.T_NB_CSV_UNGGAH (ID)
)
/
CREATE TABLE {skema}.T_NB_CSV_PESAN (
  UNGGAH_ID NUMBER(19) NOT NULL,
  URUT      NUMBER(10) NOT NULL,
  JENIS     VARCHAR2(10) NOT NULL,
  NO_BARIS  NUMBER(10),
  TEKS      VARCHAR2(4000) NOT NULL,
  CONSTRAINT PK_T_NB_CSV_PESAN PRIMARY KEY (UNGGAH_ID, URUT),
  CONSTRAINT FK_NB_CSV_PESAN_UNGGAH FOREIGN KEY (UNGGAH_ID) REFERENCES {skema}.T_NB_CSV_UNGGAH (ID)
)
/
ALTER TABLE {skema}.T_LOCATIONLIST MODIFY (SEQ_NO NUMBER(10))
/

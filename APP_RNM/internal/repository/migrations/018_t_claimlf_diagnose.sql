-- T_CLAIMLF_DIAGNOSE - diagnosa per peserta (butir bd, 27-09-2026).
--
-- ⛔ BANYAK diagnosa per peserta, dan XML yang menjawabnya - bukan tebakan.
-- `Section/ClaimLifeDetailGCNM.xml` b3923 `pyPageListProperty .DiagnoseList`
-- disajikan RepeatGrid berkelas `Data-DiagnoseLife` b3914, dan gridnya punya
-- DUA tombol yang hanya masuk akal pada daftar:
--
--   b4690 `Add`    -> `addRow`    b4700/b4841
--   b6160 `Delete` -> `deleteRow` b6170
--
-- ⚠️ RALAT ATAS BACAAN KAMI SENDIRI. Ronde sebelumnya menyimpulkan "satu
-- lawan banyak, tidak dapat diputuskan" (OQ-K.2) karena melihat grid tetapi
-- tidak membaca TOMBOLNYA. Grid dapat berarti tampilan satu baris; grid
-- ber-Add dan ber-Delete tidak dapat. Jawabannya ada dua baris di bawah
-- tempat kami berhenti membaca.
--
-- ⛔ `URUTAN` adalah padanan `.pxListSubscript`, mulai 1. Ia BUKAN hiasan:
-- `SetSTS_Reject.xml` b241 memutar `.DiagnoseList` dan b257 menulis
-- `.STS_REJECT = Primary.STS_REJECT` pada tiap barisnya, jadi urutan grid
-- adalah urutan yang orang lihat dan urutan yang rule itu tempuh.
--
-- ⛔ LEBAR `DISEASE` 1000, DAN SEBABNYA TERUKUR. Nilainya datang dari
-- `POOLDATA.DISEASE_LIFE.DISEASE` `VARCHAR2(1000)`, yang isi terpanjangnya
-- `[data DBA - katalog DEV 27-09-2026]` **290 karakter**. Kolom peserta
-- `T_CLAIMLF_PREMIUMLIST_DETAIL.DISEASE` `VARCHAR2(255)` karena itu TERBUKTI
-- KURANG - ia akan menolak nama penyakit yang sah, dan migrasi ini
-- melebarkannya pula.
--
-- ⚠️ `ICD_CODE` 100 mengikuti sumbernya (`DISEASE_LIFE.ICD_CODE`
-- `VARCHAR2(100)`, isi terpanjang 7). Lebar sumber dipakai apa adanya:
-- menyempitkannya ke 32 "karena cukup" berarti menebak bahwa kode ICD tidak
-- akan pernah lebih panjang, dan tabel sumbernya sendiri tidak menebak itu.
--
-- ⛔ `GROUP_DIAGNOSE` - nilainya `[tidak ada di korpus]`, lebarnya `[data DBA]`.
-- Dropdown b5860/b5863 ber-`pyListSource` **`associated`**, artinya daftar
-- pilihannya hidup pada RULE PROPERTI `.GROUPDIAGNOSE` di kelas
-- `Data-DiagnoseLife` - dan rule itu TIDAK ADA di ekspor: `GROUPDIAGNOSE`
-- muncul di TEPAT SATU berkas korpus, yaitu section ini sendiri. Kolomnya
-- dibuat, nilainya tidak dikarang, dan layar menyatakan daftarnya belum ada.
--
-- ⛔ `STS_REJECT` bertipe sama dengan kolom peserta (`VARCHAR2(8)`, migrasi
-- 003) - teks, tidak pernah angka (ADR-U-0022).
--
-- ⛔ KASKADE, dan buktinya bukan selera. `.DiagnoseList` hidup DI DALAM
-- halaman peserta: `SetDisease.xml` b389 menutup dengan `Obj-Save pyWorkPage`,
-- bukan menyimpan halaman diagnosa sendiri. Menghapus peserta karena itu
-- menghapus daftarnya. Daftar `TestKaskadeHanyaPadaEmpatRelasi` diperbarui
-- DENGAN BUKTI ini - pelajaran tiket 00 Komite: penjaga yang ditulis sebelum
-- relasi berikutnya lahir berhenti menjadi penjaga dan mulai menjadi pagar.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).

CREATE TABLE {skema}.T_CLAIMLF_DIAGNOSE (
  ID                     NUMBER(19) NOT NULL,
  PREMIUM_LIST_DETAIL_ID VARCHAR2(32) NOT NULL,
  URUTAN                 NUMBER(5) NOT NULL,
  ICD_CODE               VARCHAR2(100),
  DISEASE                VARCHAR2(1000),
  GROUP_DIAGNOSE         VARCHAR2(255),
  STS_REJECT             VARCHAR2(8),
  CONSTRAINT PK_T_CLAIMLF_DIAGNOSE PRIMARY KEY (ID),
  CONSTRAINT FK_DIAGNOSE_PESERTA FOREIGN KEY (PREMIUM_LIST_DETAIL_ID)
    REFERENCES {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL (ID) ON DELETE CASCADE
)
/

CREATE INDEX {skema}.IX_DIAGNOSE_PESERTA ON {skema}.T_CLAIMLF_DIAGNOSE (PREMIUM_LIST_DETAIL_ID)
/

CREATE SEQUENCE {skema}.SEQ_CLAIMLF_DIAGNOSE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/

ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL MODIFY (
  DISEASE VARCHAR2(1000)
)
/

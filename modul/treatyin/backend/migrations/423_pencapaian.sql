-- PENCAPAIAN triwulanan per kontrak
--
-- Migrasi tiket 69 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Nama kolom: `2-to-spec/KAMUS-KOLOM.md` §10.27.
--
-- ---------------------------------------------------------------------
-- INDUKNYA `KONTRAK`, DAN ERD HTML MENULIS SEBALIKNYA. INI KEPUTUSANNYA.
-- ---------------------------------------------------------------------
--   `ERD-TREATY-IN-DAN-EDM.html` baris relasi 6 menempatkan
--   `T_TREATY_LIMIT_ACHIEVEMENT` di bawah `T_TREATY_LIMIT_DETAIL` lewat
--   `LIMIT_DETAIL_ID`. Berkas ini TIDAK mengikutinya, dan sebabnya datang dari
--   ERD itu sendiri:
--
--     1. ERD menandai bukti relasi itu **DAUN-RELATIF**, dan legenda ERD
--        sendiri menyebut tingkat itu *"lemah, mungkin memungut nama milik
--        entitas lain"* - yang ditemukan hanya ruas terakhir
--        (`.AchievementLists`), bukan jalur penuhnya.
--     2. Tabel Oracle yang sungguh hidup berkunci TINGKAT KONTRAK:
--        `RDBList/GetAchievement.xml` memilih dari `POOLDATA.ACHIEVEMENT`
--        dengan `WHERE SUBSTR(NOOFFER,1,7) = SUBSTR(:CARI1,1,7)` - nomor
--        penawaran, bukan pengenal layer mana pun. Enam belas kolomnya memuat
--        `NOPOLIS`, `NOOFFER`, `QUARTER`, `QUARTERYEAR`, `IDCURRENCY`; nol
--        kolom menunjuk rincian layer.
--     3. ERD menyatakan dirinya **"POTRET SISTEM LAMA"** dan memotret POHON
--        HALAMAN Pega. Halaman bersarang tidak berarti datanya bersarang -
--        itu pola yang sama dengan "20 kolom lawan 99 skalar".
--
--   Jadi keduanya benar pada pertanyaan yang berbeda: ERD benar tentang di
--   mana LAYARNYA menampilkannya, SQL benar tentang di mana FAKTANYA
--   disimpan. Model baru menyimpan fakta di tempat fakta itu hidup.
--
--   SYARAT PEMBALIKAN: satu baris `ACHIEVEMENT` yang kolom pengenal layernya
--   terisi, atau satu aturan yang menulisnya per rincian layer. Keduanya akan
--   membatalkan bacaan ini, dan tiket 69 berubah induk - murah dilakukan
--   SELAMA tiket 44 belum memuat data.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS: TOLAK, bukan ikut hapus
-- ---------------------------------------------------------------------
--   `ERD.md` §2.9 - dokumen MENGIKAT untuk perilaku hapus:
--   `KONTRAK 1--o< PENCAPAIAN [hapus: tolak]`.
--
--   `tolak` disengaja: pencapaian adalah ANGKA YANG PERNAH DIBUKUKAN, dan
--   kontrak yang punya pencapaian tidak boleh lenyap bersamanya. Alasannya
--   sebangun dengan `KONTRAK 1--< VERSI_KONTRAK [hapus: tolak]` di §2.1.
--   Diwujudkan dengan TIDAK menulis klausa ON DELETE - bawaan Oracle menolak,
--   dan di sini bawaan itu DIPILIH, bukan dibiarkan.
--
--   ⚠️ Penanda `[G3]` di §2.9 digantikan `REKONSILIASI-XML-VS-DDL.md` §3.1
--   (2 Oktober 2026) yang memindahkannya ke gelombang 1. Yang digantikan hanya
--   penanda gelombangnya; perilaku hapusnya tetap.
--
-- ---------------------------------------------------------------------
-- `LOG_ACHIEVEMENT` TIDAK MENJADI TABEL KEDUA
-- ---------------------------------------------------------------------
--   Sistem lama punya DUA tabel - `ACHIEVEMENT` (dibaca `GetAchievement.xml`)
--   dan `LOG_ACHIEVEMENT` (ditulis `InsertToLogAchievement_SQL.xml`) - dan
--   keduanya BUKAN salinan satu sama lain. Daftar kolomnya dibandingkan, dan
--   `LOG_ACHIEVEMENT` membawa EMPAT yang `ACHIEVEMENT` tidak punya. Keempatnya
--   dipisah menurut sifatnya, bukan menurut nama tabelnya:
--
--     CASHCALLCLAIM             -> kolom CASH_CALL_KLAIM di bawah; ia FAKTA
--                                  tersendiri, tidak dapat dihitung dari yang lain
--     INCUREDCLAIM, TOTAL,      -> TIDAK disimpan. Turunan - INV-58. Hilir
--     LOSSRATIO                    menghitungnya saat membaca
--     PXCREATEOPNAME, INSERTDATE-> JEJAK_PERUBAHAN (tiket 39), yang sudah ada
--
--   SYARAT PEMBALIKAN untuk ketiga turunan: satu baris `LOG_ACHIEVEMENT` yang
--   `TOTAL`-nya TIDAK sama dengan jumlah kolom penyusunnya. Bila itu ada,
--   ketiganya berhenti menjadi turunan - ia angka yang pernah dibukukan, dan
--   INV-58 sendiri memberi jalannya. Pemeriksaan itu menuntut ISI tabel lama;
--   ditagih ke tiket 44.
--
-- ---------------------------------------------------------------------
-- NOL KOLOM SALINAN - INV-59
-- ---------------------------------------------------------------------
--   `GetAchievement.xml` juga memilih `SOBNAME`, `TREATYGROUPNAME`, dan
--   `TREATYTYPE`. Ketiganya ATRIBUT KONTRAK dan LAYER yang disalin ke dalam
--   baris pencapaian, dan ketiganya TIDAK dibawa ke sini: ia dibaca lewat
--   kunci asingnya. Menyalinnya berarti nama kelompok yang diperbaiki di
--   master tidak pernah merambat ke baris pencapaian yang sudah ada.
--
-- PERNYATAAN KEPUTUSAN - kunci alami (ID_KONTRAK, TRIWULAN, TAHUN_TRIWULAN,
-- KODE_MATA_UANG) TIDAK dipasang sebagai UNIQUE: `SPEC-INVARIAN.md` berhenti
-- di INV-71 dan belum menomorinya. Aturan `Z00_KUNCI_ALAMI.sql` berlaku -
-- constraint tanpa invarian tidak punya tempat untuk gagal. Tagihannya di
-- `SPEC-INVARIAN.md` §13.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE TABLE {skema}.PENCAPAIAN (
  ID_PENCAPAIAN      NUMBER(19)          NOT NULL,
  ID_KONTRAK         NUMBER(19)          NOT NULL,
  TRIWULAN           NUMBER(10)          NOT NULL,
  TAHUN_TRIWULAN     NUMBER(10)          NOT NULL,
  KODE_MATA_UANG     VARCHAR2(1000 CHAR) NOT NULL,
  PREMI              NUMBER(38,8),
  KOMISI_REASURANSI  NUMBER(38,8),
  BROKERAGE          NUMBER(38,8),
  PREMI_BERSIH       NUMBER(38,8),
  KLAIM_DIBAYAR      NUMBER(38,8),
  CASH_CALL_KLAIM    NUMBER(38,8),
  KLAIM_OUTSTANDING  NUMBER(38,8),
  CONSTRAINT PK_PENCAPAIAN PRIMARY KEY (ID_PENCAPAIAN),
  CONSTRAINT FK_PENCAPAIAN_1 FOREIGN KEY (ID_KONTRAK)
    REFERENCES {skema}.KONTRAK (ID_KONTRAK)
)
/
CREATE INDEX {skema}.IX_PENCAPAIAN_KONTRAK ON {skema}.PENCAPAIAN (ID_KONTRAK)
/

-- Tiga perkakas pemindahan - MIGRASI_KORELASI, MIGRASI_PENDARATAN,
-- MIGRASI_NILAI_DITOLAK
--
-- Migrasi tiket 71, 72, 73 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- ⛔ BERTENGGAT TERHADAP TIKET 44. Tiket itu memuat data pertama, dan
-- SESUDAHNYA ketiganya kehilangan sebagian gunanya selamanya: tanpa tabel
-- korelasi tidak ada jalan menelusuri baris baru kembali ke asalnya, dan
-- tiket 61 ("baris warisan diperbaiki ke keadaan sah yang dipilih eksplisit")
-- menuntut justru penelusuran itu.
--
-- ---------------------------------------------------------------------
-- KENAPA KETIGANYA TIDAK PUNYA KOTAK DI ERD, DAN ITU WAJAR
-- ---------------------------------------------------------------------
--   `ERD-TREATY-IN-DAN-EDM.html` menyatakan dirinya "POTRET SISTEM LAMA,
--   BUKAN RANCANGAN", sementara ketiganya PERKAKAS SISTEM BARU untuk
--   pemindahan. Perkakas migrasi memang tidak akan muncul di potret sistem
--   yang sedang dipindahkan.
--
--   Dan ERD TIDAK membuangnya: ketiganya disebut di daftar "7 tabel tanpa
--   jalur Pega ... tidak digambar, dilaporkan di sini", dan satu di antaranya
--   bahkan punya baris relasi - baris 38,
--   `T_TREATY_MIG_LANDING -> T_TREATY_MIG_REJECTED` lewat `LANDING_ID`.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS
-- ---------------------------------------------------------------------
--   `MIGRASI_KORELASI` dan `MIGRASI_PENDARATAN`: NOL kunci asing, jadi nol
--   perilaku hapus. Sebabnya di bawah, dan ia keputusan - bukan kelalaian.
--
--   `MIGRASI_NILAI_DITOLAK` -> `MIGRASI_PENDARATAN`: **TOLAK**, diwujudkan
--   dengan tidak menulis klausa `ON DELETE`.
--
--   ⚠️ Tiket 73 yang ditulis ronde 5 menyebut "ikut hapus". Itu DIKOREKSI di
--   sini, dan sumbernya ERD: baris 38 menulis kolom `ON DELETE`-nya **"di
--   Go"** - ia TIDAK meresepkan aturan tingkat basis data. Dari dua bacaan
--   yang mungkin, `tolak` yang dipilih, dengan alasan yang sudah dipakai
--   migrasi 425 untuk arsip: catatan forensik yang lenyap bersama induknya
--   berhenti menjadi catatan forensik tepat saat ia paling dibutuhkan.
--   Pembuangan yang disengaja tetap mungkin - lewat jalur yang membuang
--   anaknya lebih dulu, dan jalur itu terlihat di kode.
--
-- ---------------------------------------------------------------------
-- NOL KUNCI ASING PADA KORELASI DAN PENDARATAN - keputusan, bukan lupa
-- ---------------------------------------------------------------------
--   Jejak asal-usul harus BERTAHAN MELEWATI penghapusan barisnya. Kunci
--   asing memaksa salah satu dari dua, dan keduanya merusak: `ikut hapus`
--   menghapus jejaknya bersama barisnya, `tolak` membuat baris yang salah
--   muat TIDAK DAPAT DIBUANG. Keduanya bertentangan dengan guna tabel ini.
--
--   Dan pada `MIGRASI_PENDARATAN` ada sebab kedua yang berdiri sendiri:
--   muatan yang GAGAL DIURAI justru belum punya baris baru untuk ditunjuk.
--   Kunci asing akan menolak baris yang paling penting disimpan.
--
--   ⚠️ Ongkosnya dibayar dan dinyatakan: basis data TIDAK menjaga
--   keterhubungannya - baris korelasi dapat menunjuk pengenal yang sudah
--   tidak ada. Yang menjaganya uji rekonsiliasi di tiket 44, bukan
--   constraint.
--
-- ---------------------------------------------------------------------
-- INV-61 berlaku atas MUATAN
-- ---------------------------------------------------------------------
--   Kolom `MUATAN` tidak punya jalur baca aplikasi; ia bukan sumber kanonik.
--   Yang membacanya pemindahan dan manusia yang menyelidiki. Basis data
--   tidak dapat menegakkannya; yang menjaganya tinjauan kode dan
--   `TestArsipTidakPunyaJalurBaca`, yang menyapu kolom bernama `MUATAN` di
--   seluruh berkas Go modul ini.
--
-- ⚠️ `MUATAN` dan `NILAI_MENTAH` bertipe `VARCHAR2(4000 CHAR)`, bukan CLOB -
-- seluruh repo ini nol CLOB. BATAS YANG DINYATAKAN: muatan yang lebih
-- panjang AKAN DITOLAK saat disisipkan, dan penolakan LEBIH BAIK daripada
-- pemotongan diam-diam: pendaratan yang terpotong adalah pendaratan yang
-- berbohong. Bila `M_TREATY_IN.JSONDATA` ternyata melampauinya, kolomnya
-- dinaikkan menjadi CLOB lewat migrasi tersendiri - dan ukurannya DIUKUR,
-- bukan ditaksir.
--
-- PERNYATAAN KEPUTUSAN - kunci alami ketiganya TIDAK dipasang sebagai
-- `UNIQUE`: `SPEC-INVARIAN.md` berhenti di INV-71 dan belum menomorinya.
-- Aturan `ddl-usulan/Z00_KUNCI_ALAMI.sql` berlaku - "constraint tanpa
-- invarian tidak punya tempat untuk gagal". Tagihannya di
-- `SPEC-MODEL-DATA.md` §13.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE SEQUENCE {skema}.SEQ_TRIN_MIGRASI_KORELASI START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_MIGRASI_PENDARATAN START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_TRIN_MIG_NILAI_DITOLAK START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE TABLE {skema}.MIGRASI_KORELASI (
  ID_MIGRASI_KORELASI  NUMBER(19)          NOT NULL,
  KUNCI_PEGA           VARCHAR2(1000 CHAR),
  ID_PEGA              VARCHAR2(1000 CHAR),
  ID_WARISAN           VARCHAR2(1000 CHAR),
  ID_KONTRAK_BARU      NUMBER(19),
  DIPINDAHKAN_PADA     DATE                NOT NULL,
  CONSTRAINT PK_MIGRASI_KORELASI PRIMARY KEY (ID_MIGRASI_KORELASI)
)
/
CREATE TABLE {skema}.MIGRASI_PENDARATAN (
  ID_MIGRASI_PENDARATAN  NUMBER(19)          NOT NULL,
  KUNCI_WARISAN          VARCHAR2(1000 CHAR) NOT NULL,
  MUATAN                 VARCHAR2(4000 CHAR) NOT NULL,
  MENDARAT_PADA          DATE                NOT NULL,
  CONSTRAINT PK_MIGRASI_PENDARATAN PRIMARY KEY (ID_MIGRASI_PENDARATAN)
)
/
CREATE TABLE {skema}.MIGRASI_NILAI_DITOLAK (
  ID_MIGRASI_NILAI_DITOLAK  NUMBER(19)          NOT NULL,
  ID_MIGRASI_PENDARATAN     NUMBER(19)          NOT NULL,
  JALUR_SIMPUL              VARCHAR2(1000 CHAR) NOT NULL,
  NILAI_MENTAH              VARCHAR2(4000 CHAR) NOT NULL,
  SEBAB_DITOLAK             VARCHAR2(1000 CHAR) NOT NULL,
  CONSTRAINT PK_MIGRASI_NILAI_DITOLAK PRIMARY KEY (ID_MIGRASI_NILAI_DITOLAK),
  CONSTRAINT FK_MIGRASI_NILAI_DITOLAK_1 FOREIGN KEY (ID_MIGRASI_PENDARATAN)
    REFERENCES {skema}.MIGRASI_PENDARATAN (ID_MIGRASI_PENDARATAN)
)
/
CREATE INDEX {skema}.IX_MIG_KORELASI_WARISAN ON {skema}.MIGRASI_KORELASI (ID_WARISAN)
/
CREATE INDEX {skema}.IX_MIG_KORELASI_PEGA ON {skema}.MIGRASI_KORELASI (ID_PEGA)
/
CREATE INDEX {skema}.IX_MIG_PENDARATAN_KUNCI ON {skema}.MIGRASI_PENDARATAN (KUNCI_WARISAN)
/
CREATE INDEX {skema}.IX_MIG_DITOLAK_PENDARATAN ON {skema}.MIGRASI_NILAI_DITOLAK (ID_MIGRASI_PENDARATAN)
/

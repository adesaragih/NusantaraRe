-- ARSIP_MUATAN_KELUAR - muatan yang dikirim ke hilir, disimpan apa adanya
--
-- Migrasi tiket 74 Treaty In, dan prasyarat tiket 42. Satu berkas = satu
-- langkah migrasi. Pernyataan dipisahkan oleh baris yang hanya berisi tanda
-- garis miring.
--
-- Asal struktur: `ERD-TREATY-IN-DAN-EDM.html` baris relasi 39 -
-- `TREATY_IN -> T_TREATY_OUTBOUND_ARCHIVE` lewat `TREATY_IN_ID`, berindex.
-- Nama kolom: `2-to-spec/KAMUS-KOLOM.md` §10.29.
--
-- Tiket 42 ("arsip JSON sistem lama disimpan dan dapat ditunjukkan") berdiri
-- di papan sejak batch 1 dengan tempat penyimpanan yang tidak ada. Yang
-- membuatnya mendesak bukan kerapian: hilir menerima SATU muatan berisi 175
-- argumen, dan ketika hilir berkata "angkanya salah", satu-satunya cara
-- memisahkan "kami mengirim salah" dari "kalian membaca salah" adalah
-- MENUNJUKKAN apa yang dikirim.
--
-- ---------------------------------------------------------------------
-- INV-18 - PERILAKU HAPUS: TOLAK. Dan ERD menyerahkannya ke Go.
-- ---------------------------------------------------------------------
--   ERD baris 39 menulis kolom ON DELETE-nya **"di Go"**, bukan `CASCADE` -
--   satu dari empat relasi yang ia tandai begitu (34 `T_VIEW_COMMENT`,
--   35 `DOCUMENT_TREATY_IN`, 38 `T_TREATY_MIG_REJECTED`, 39 relasi ini).
--   Artinya ERD TIDAK meresepkan aturan hapus tingkat basis data untuknya.
--
--   Di sini aturan itu diwujudkan sebagai TOLAK - kunci asing TANPA klausa
--   `ON DELETE`, yang pada Oracle berarti penghapusan induk ditolak selama
--   anaknya ada. Alasannya: arsip yang lenyap bersama kontraknya berhenti
--   menjadi arsip tepat pada saat ia paling dibutuhkan, sebab perselisihan
--   dengan hilir sering muncul SESUDAH barisnya dibersihkan. Sebangun dengan
--   `KONTRAK 1--< VERSI_KONTRAK [hapus: tolak]`, `ERD.md` §2.1.
--
--   ⚠️ "di Go" TIDAK dibaca sebagai "tanpa kunci asing". Kunci asing tetap
--   dipasang - ia yang menjamin arsip tidak menunjuk kontrak yang tidak ada.
--   Yang diserahkan ke Go adalah APA YANG TERJADI saat seseorang hendak
--   menghapus kontrak berarsip: lapisan aplikasi yang menjelaskan sebabnya,
--   basis data yang menolaknya.
--
--   SYARAT PEMBALIKAN: bila pemilik proses menetapkan masa simpan arsip,
--   baris yang lewat masa itu dibuang lewat jalur pembuangan tersendiri -
--   BUKAN dengan melonggarkan kunci asingnya menjadi `CASCADE`.
--
-- ---------------------------------------------------------------------
-- INV-61 - ARSIP TIDAK PUNYA JALUR BACA
-- ---------------------------------------------------------------------
--   Kolom `MUATAN` tidak boleh dibaca oleh satu pun jalur aplikasi: ia bukan
--   sumber kanonik. Yang membacanya pemindahan dan manusia yang menyelidiki.
--   Basis data TIDAK dapat menegakkan ini - ia tidak dapat melarang `SELECT`
--   atas kolomnya sendiri tanpa mencabut akses ke tabelnya (`SPEC-INVARIAN.md`
--   §7). Yang menjaganya tinjauan kode, dan uji `TestArsipTidakPunyaJalurBaca`.
--
-- ⚠️ `MUATAN` bertipe `VARCHAR2(4000 CHAR)`, bukan CLOB. Alasannya: seluruh
-- repo ini nol CLOB, dan menambahkan tipe baru menuntut jalur baca/tulis yang
-- berbeda di driver. BATAS YANG DINYATAKAN: muatan yang lebih panjang dari
-- 4000 bita AKAN DITOLAK basis data saat disisipkan - dan penolakan itu LEBIH
-- BAIK daripada pemotongan diam-diam, sebab arsip yang terpotong adalah arsip
-- yang berbohong. Bila 175 argumen ternyata melampauinya, kolomnya dinaikkan
-- menjadi CLOB lewat migrasi tersendiri, dan ukurannya diukur - bukan ditaksir.
--
-- PERNYATAAN KEPUTUSAN - kunci alami (ID_KONTRAK, TUJUAN, DIKIRIM_PADA) TIDAK
-- dipasang sebagai UNIQUE: `SPEC-INVARIAN.md` berhenti di INV-71. Dan di sini
-- ada sebab kedua yang berdiri sendiri: dua pengiriman ke tujuan yang sama
-- pada detik yang sama adalah keadaan yang MUNGKIN dan SAH - arsip mencatat
-- apa yang terjadi, ia tidak mengadilinya.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
CREATE TABLE {skema}.ARSIP_MUATAN_KELUAR (
  ID_ARSIP_MUATAN_KELUAR  NUMBER(19)          NOT NULL,
  ID_KONTRAK              NUMBER(19)          NOT NULL,
  TUJUAN                  VARCHAR2(1000 CHAR) NOT NULL,
  DIKIRIM_PADA            DATE                NOT NULL,
  BERHASIL                VARCHAR2(40 CHAR)   NOT NULL,
  MUATAN                  VARCHAR2(4000 CHAR) NOT NULL,
  CONSTRAINT PK_ARSIP_MUATAN_KELUAR PRIMARY KEY (ID_ARSIP_MUATAN_KELUAR),
  CONSTRAINT FK_ARSIP_MUATAN_KELUAR_1 FOREIGN KEY (ID_KONTRAK)
    REFERENCES {skema}.KONTRAK (ID_KONTRAK)
)
/
CREATE INDEX {skema}.IX_ARSIP_MUATAN_KONTRAK ON {skema}.ARSIP_MUATAN_KELUAR (ID_KONTRAK)
/

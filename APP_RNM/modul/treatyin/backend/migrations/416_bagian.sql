-- BAGIAN - bagian NuRe atas sebuah layer, beserta premi brutonya per mata uang
--
-- Migrasi tiket 33 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/33-*.md`, `2-to-spec/KAMUS-KOLOM.md` §10.5,
-- `ddl-usulan/{24_BAGIAN,30_NILAI_PREMI_BRUTO,31_NILAI_PREMI_BRUTO_MINIMUM}.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- INV-64: satu BAGIAN per LAYER - ID_LAYER unik.
--
-- ⭐ §12.3 MENGOREKSI §8: `BAGIAN` BUKAN satu entitas untuk kedua cabang. Di
-- cabang non-proporsional ia kelas tersendiri dengan baris `Share[]`-nya
-- sendiri; di cabang proporsional TIDAK ADA kelas tersendiri sama sekali -
-- bagian NuRe di sana adalah ATRIBUT pada baris `DETAIL_PROPORSIONAL` (417).
--
-- ⛔ Sepuluh daftar `RNMSpreadedList…XOL` TIDAK DISIMPAN - turunan (ADR-0037,
-- §12.3). Itu pernyataan keputusan: siapa pun yang menyimpannya menanam kolom
-- turunan yang INV-58 larang.
--
-- ---------------------------------------------------------------------
-- ⛔ INV-33 BELUM DITEGAKKAN, dan sebabnya harus dibaca sebelum seseorang
-- memasangnya dengan tergesa.
-- ---------------------------------------------------------------------
--   APA     : tidak ada yang menolak BAGIAN di bawah kontrak ber-SIFAT_PROPORSI
--             `PROPORSIONAL`.
--   BENTUK  : tiket 33 menggolongkannya INDEKS UNIK. SIFAT_PROPORSI ada DUA
--             tabel di atas (BAGIAN -> LAYER -> VERSI_KONTRAK -> KONTRAK),
--             sehingga di Oracle bentuk itu berarti MATERIALIZED VIEW
--             ber-REFRESH ON COMMIT dengan indeks unik di atasnya.
--   KENAPA BELUM: `2-to-spec/F-13-MV-TANPA-PEMANTAU-KEBASIAN.md` menyatakan
--             penegakan lewat MV di modul ini BELUM punya pemantau kebasian,
--             dan bahwa kewajiban itu "tidak pernah menjadi kemampuan" - nol
--             `P-nn` yang dapat dinyatakan selesai. Kalimat F-13 sendiri:
--             "MV yang gagal me-refresh berhenti menegakkan tanpa satu galat
--             pun" - ia "terlihat terpasang di setiap pemeriksaan yang pernah
--             dijalankan orang, dan tidak menolak apa pun".
--             Memasangnya SEKARANG, tanpa Oracle yang dapat dijangkau (`L-3`)
--             untuk menguji perilaku refresh-nya, menanam persis jebakan itu.
--   AKIBAT  : bagian pada kontrak proporsional diterima. Jalur simpan harus
--             menolaknya, dan pesannya menyebut sifat proporsi kontraknya.
--   DITAGIH : pemilik F-13 - pemantau kebasian lebih dulu, MV sesudahnya.
--             Hal yang sama berlaku bagi INV-32 (417), INV-47, INV-50, INV-51.
--
-- ⛔ INV-36 (premi bruto terisi tanpa mata uang -> ditolak) dan INV-52 (premi
-- bruto dikurangi seluruh potongan sama dengan premi bersih) juga belum
-- ditegakkan. INV-52 memang DIBAWA sebagai invarian di sini; penegakannya
-- irisan 37 menurut tiket 33 sendiri.
CREATE TABLE {skema}.BAGIAN (
  ID_BAGIAN              NUMBER(19)   NOT NULL,
  ID_LAYER               NUMBER(19)   NOT NULL,
  CAKUPAN                VARCHAR2(1000 CHAR),
  PERSEN_BAGIAN_NURE     NUMBER(38,8),
  PREMI_BRUTO            NUMBER(38,8) NOT NULL,
  PREMI_BRUTO_MINIMUM    NUMBER(38,8),
  ID_SUSUNAN_RETRO       NUMBER(19),
  PERSEN_BAGIAN_DIPAKAI  NUMBER(38,8),
  CONSTRAINT PK_BAGIAN PRIMARY KEY (ID_BAGIAN),
  CONSTRAINT UQ_BAGIAN UNIQUE (ID_LAYER),
  CONSTRAINT FK_BAGIAN_1 FOREIGN KEY (ID_LAYER)
    REFERENCES {skema}.LAYER (ID_LAYER)
)
/
CREATE TABLE {skema}.NILAI_PREMI_BRUTO (
  ID_NILAI_PREMI_BRUTO  NUMBER(19)          NOT NULL,
  ID_BAGIAN             NUMBER(19)          NOT NULL,
  KODE_MATA_UANG        VARCHAR2(1000 CHAR) NOT NULL,
  NILAI                 NUMBER(38,8)        NOT NULL,
  CONSTRAINT PK_NILAI_PREMI_BRUTO PRIMARY KEY (ID_NILAI_PREMI_BRUTO),
  CONSTRAINT FK_NILAI_PREMI_BRUTO_1 FOREIGN KEY (ID_BAGIAN)
    REFERENCES {skema}.BAGIAN (ID_BAGIAN)
)
/
CREATE TABLE {skema}.NILAI_PREMI_BRUTO_MINIMUM (
  ID_NILAI_PREMI_BRUTO_MINIMUM  NUMBER(19)          NOT NULL,
  ID_BAGIAN                     NUMBER(19)          NOT NULL,
  KODE_MATA_UANG                VARCHAR2(1000 CHAR) NOT NULL,
  NILAI                         NUMBER(38,8)        NOT NULL,
  CONSTRAINT PK_NILAI_PREMI_BRUTO_MINIMUM PRIMARY KEY (ID_NILAI_PREMI_BRUTO_MINIMUM),
  CONSTRAINT FK_NILAI_PREMI_BRUTO_MIN_1 FOREIGN KEY (ID_BAGIAN)
    REFERENCES {skema}.BAGIAN (ID_BAGIAN)
)
/
CREATE INDEX {skema}.IX_NILAI_PB_BAGIAN ON {skema}.NILAI_PREMI_BRUTO (ID_BAGIAN)
/
CREATE INDEX {skema}.IX_NILAI_PB_MIN_BAGIAN ON {skema}.NILAI_PREMI_BRUTO_MINIMUM (ID_BAGIAN)
/

-- POTONGAN - potongan atas premi, aturan sama di kedua pelekatannya
--
-- Migrasi tiket 37 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/37-*.md`, `2-to-spec/KAMUS-KOLOM.md`,
-- `ddl-usulan/33_POTONGAN.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- ---------------------------------------------------------------------
-- INDUK POLIMORFIK - KTV-B: DUA KOLOM BERNAMA, TEPAT SATU TERISI
-- ---------------------------------------------------------------------
--   Potongan melekat pada BAGIAN (cabang non-proporsional) ATAU pada
--   DETAIL_PROPORSIONAL (cabang proporsional) - tidak pernah keduanya, tidak
--   pernah nol.
--
--   Bentuk ini MENGGANTIKAN kolom tunggal `ID_INDUK_POTONGAN`, yang tidak
--   dapat punya kunci asing sama sekali: sasarannya bergantung pada nilai
--   kolom lain (INV-17). Dua kolom bernama, MASING-MASING berkunci asing,
--   membuat baris yatim MUSTAHIL - bukan sekadar tidak dianjurkan.
--
--   Dasar dan syarat pembalikan: `KEPUTUSAN-TANPA-VERIFIKASI.md` §7 butir
--   KTV-B. DITINJAU ULANG bila induknya bertambah menjadi TIGA - pada titik
--   itu CHECK dua cabang menjadi tiga cabang, dan bentuk tabel pemetaan
--   mungkin lebih murah.
--
-- ⭐ `CK_POTONGAN_INDUK` ADALAH SATU-SATUNYA `CHECK` DI MODUL INI, dan itu
-- tidak melanggar apa pun. ADR-0056 (K-4) melarang TRIGGER dan PROCEDURE
-- pembawa aturan bisnis; ADR-0038 melarang `CHECK` yang memuat DAFTAR NILAI
-- bagi himpunan yang dapat bertambah. CHECK di bawah bukan keduanya - ia
-- menyatakan BENTUK baris, bukan nilai yang sah, dan `ddl-usulan/` sendiri
-- memuatnya (hanya ia dan `CK_PENYEBARAN_INDUK`, tiket 38).
--
-- INV-15 terpasang sebagai DUA `UNIQUE`, satu per pelekatan - bukan satu yang
-- mencampur keduanya. Oracle tidak membandingkan NULL, sehingga baris yang
-- melekat pada BAGIAN tidak pernah bentrok dengan baris yang melekat pada
-- DETAIL_PROPORSIONAL.
--
-- ⛔ INV-63 TIDAK dapat ditegakkan basis data sama sekali: ia menuntut rumus
-- potongan ditemukan TEPAT SATU KALI di seluruh basis kode. Daftar periksa
-- tiket 37 menyebutnya "diperiksa lewat TINJAUAN SKEMA TERTULIS", bukan lewat
-- constraint. DITAGIH: tiket lapisan aplikasi yang menulis rumusnya - dan
-- rumusnya harus lahir SATU KALI, dipakai kedua pelekatan.
--
-- ⛔ INV-52 (premi bruto dikurangi seluruh potongan sama dengan premi bersih)
-- dibawa tiket 33 dan penegakannya diserahkan ke irisan ini - tetapi ia
-- menjumlahkan BANYAK BARIS terhadap induknya, yang CHECK tidak dapat
-- nyatakan. Bentuknya sama dengan INV-47/50/51: materialized view, dan karena
-- itu tertahan `F-13` (lihat berkas 416). BELUM DITEGAKKAN.
CREATE TABLE {skema}.POTONGAN (
  ID_POTONGAN             NUMBER(19)        NOT NULL,
  ID_BAGIAN               NUMBER(19),
  ID_DETAIL_PROPORSIONAL  NUMBER(19),
  ID_JENIS_POTONGAN       NUMBER(19)        NOT NULL,
  DASAR_PERHITUNGAN       VARCHAR2(40 CHAR) NOT NULL,
  PERSEN_POTONGAN         NUMBER(38,8)      NOT NULL,
  CONSTRAINT PK_POTONGAN PRIMARY KEY (ID_POTONGAN),
  CONSTRAINT UQ_POTONGAN UNIQUE (ID_BAGIAN, ID_JENIS_POTONGAN),
  CONSTRAINT UQ_POTONGAN_2 UNIQUE (ID_DETAIL_PROPORSIONAL, ID_JENIS_POTONGAN),
  CONSTRAINT FK_POTONGAN_1 FOREIGN KEY (ID_BAGIAN)
    REFERENCES {skema}.BAGIAN (ID_BAGIAN),
  CONSTRAINT FK_POTONGAN_2 FOREIGN KEY (ID_DETAIL_PROPORSIONAL)
    REFERENCES {skema}.DETAIL_PROPORSIONAL (ID_DETAIL_PROPORSIONAL),
  CONSTRAINT FK_POTONGAN_3 FOREIGN KEY (ID_JENIS_POTONGAN)
    REFERENCES {skema}.JENIS_POTONGAN (ID_JENIS_POTONGAN)
)
/
-- CK_POTONGAN_INDUK berdiri sebagai pernyataan TERSENDIRI, bukan di dalam
-- CREATE TABLE - sama seperti `ddl-usulan/33_POTONGAN.sql`. Sebabnya bukan
-- selera: pembaca kolom penjaga (`inti/backend/penjaga`, polaCreateTabel)
-- mengurai tiap baris di dalam CREATE TABLE sebagai definisi kolom, dan baris
-- lanjutan CHECK yang diawali `OR` terbaca sebagai kolom bernama "OR".
ALTER TABLE {skema}.POTONGAN ADD CONSTRAINT CK_POTONGAN_INDUK
  CHECK ( (ID_BAGIAN IS NOT NULL AND ID_DETAIL_PROPORSIONAL IS NULL)
       OR (ID_BAGIAN IS NULL     AND ID_DETAIL_PROPORSIONAL IS NOT NULL) )
/
CREATE INDEX {skema}.IX_POTONGAN_JENIS ON {skema}.POTONGAN (ID_JENIS_POTONGAN)
/

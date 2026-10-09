-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`/`445`/`446`/`448`/`449`/`450`: rentang `treatyin` (`400-439`) PENUH,
-- dan `451` jatuh di rentang modul ini (`440-479`).
--
-- ---------------------------------------------------------------------
-- `RevisionDate` — satu-satunya properti yang Save penyesuaian tolak simpan
-- ---------------------------------------------------------------------
--
-- Laporan pemakai 8 Oktober 2026, sesudah menekan Save di layar Treaty In
-- Adjustment:
--
--   Data Sudah Disimpan Dengan ID : 1000001/R01
--   TIDAK tersimpan (belum punya kolom/tabel): RevisionDate
--
-- ⭐ Asalnya BUKAN karangan layar. `TreatyInEDMNew` langkah [3]
-- ("Set EDM Properties") menyetelnya bersama `EDMState` dan
-- `EDMMaterialType`:
--
--   TreatyIn.RevisionDate = @CurrentDateTime()
--
-- dan `services.SusunDraf` menyalinnya apa adadanya
-- (`draf_penyesuaian.go`). Jadi SETIAP draf penyesuaian membawanya, dan
-- setiap Save melaporkannya hilang — bukan kasus pinggir.
--
-- ⚠️ Ia TANGGAL REVISI, bukan `EDMEffective` dan bukan `Commencement`:
-- kapan draf itu DIBUAT, yang menjadi jejak waktu satu-satunya bagi versi
-- penyesuaian. Tanpa kolomnya, dua revisi kontrak yang sama nol pembeda
-- waktu di basis data.
--
-- ⛔ `REVISIONSTATE` (migrasi 448) BERBEDA dan sudah ada — ia KEADAAN jalur
-- revisi, bukan tanggalnya. Nama yang berdekatan, isi yang tidak
-- bersinggungan; disebut di sini supaya tidak ada yang mengira 448 sudah
-- menutupnya.
--
-- ⚠️ `VARCHAR2(4000 CHAR)` seperti SELURUH kolom pendaratan, termasuk kolom
-- tanggal lain (`COMMENCEMENT`, `TERMINATION`, `EDMEFFECTIVE`). Nilainya
-- stempel Pega apa adanya — `20261008T083015.000 GMT` — dan services yang
-- menafsirkannya. Memakai `DATE` di sini akan membuat satu kolom tanggal
-- berbeda dari enam saudaranya, dan pemuatnya menulis TEKS ke semuanya.
--
-- ⚠️ Penulis dan pembaca `T_TREATY_*` TOLERAN terhadap kolom yang belum
-- terpasang (`kolomTerpasang`): sebelum berkas ini dijalankan Save tetap
-- berjalan dan melaporkan properti itu sebagai tidak tersimpan. Itulah yang
-- pemakai lihat.

ALTER TABLE {skema}.T_TREATY_REVISION ADD (
  REVISIONDATE                     VARCHAR2(4000 CHAR)
)
/

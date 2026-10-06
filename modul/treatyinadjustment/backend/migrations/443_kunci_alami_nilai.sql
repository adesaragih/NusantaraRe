-- Kunci alami NILAI_SELISIH dan NILAI_SEBELUM_PRO_RATE — tiket 76 dan 77.
--
-- Migrasi tiket 76 dan 77, papan Treaty In Adjustment. Satu berkas = satu
-- langkah migrasi. Pernyataan dipisahkan oleh baris yang hanya berisi tanda
-- garis miring.
--
-- ⛔ INILAH yang kedua tiket sebut dan yang `442` belum punya. Daftar
-- periksa tiket 76 berbunyi *"`NILAI_SELISIH` berdiri, dan MATA UANG ADA DI
-- DALAM KUNCI ALAMINYA"*; jalur gagal kedua tiket berbunyi *"dua baris
-- berbesaran dan bermata-uang sama pada satu versi -> DITOLAK kunci
-- alaminya"*. Tanpa constraint ini, jalur gagal itu tidak punya tempat
-- untuk gagal — dan uji negatifnya tidak dapat ditulis.
--
-- `ADR-0048` butir 3: kunci bisnis **termasuk mata uang**. `ADR-0053`: mata
-- uang adalah daftar.
--
-- =====================================================================
-- ⛔ KOLOMNYA DIPERSEMPIT LEBIH DULU, DAN ITU WAJIB — BUKAN KERAPIAN
-- =====================================================================
--   `442` (salinan `426`) memberi `KODE_BESARAN` dan `KODE_MATA_UANG`
--   lebar `VARCHAR2(1000 CHAR)`. Pada basis data AL32UTF8 — terukur, ini
--   charset POOLDATA — satu CHAR memesan 4 bita, jadi kedua kolom itu
--   4.000 bita masing-masing.
--
--   Kunci alami di bawah mencakup keduanya ditambah `NUMBER(19)`:
--
--     22 + 4.000 + 4.000  =  8.022 bita
--
--   Batas panjang kunci index Oracle kira-kira 75% ukuran blok dikurangi
--   overhead — 6.398 bita pada blok 8K. **Constraint itu akan ditolak
--   ORA-01450**, dan tiket 76 maupun 77 tidak dapat ditutup selama
--   lebarnya begitu.
--
--   Lebar barunya DIUKUR, bukan ditebak. Sapuan 5 Oktober 2026 atas pohon
--   `TreatyIn.ValueDifference` di 1.612 dokumen:
--
--     nama simpul terpanjang    24 aksara  (`LimitFacShareSummaryList`)
--     JALUR penuh terpanjang    49 aksara  (`TreatyIn.ValueDifference.LimitFacShareSummaryList`)
--     kode mata uang terpanjang  9 aksara  (terukur di tabel pendaratan;
--                                           lazimnya 3, yang 9 adalah
--                                           `1/04/2023` — tanggal di kolom
--                                           mata uang, data kotor yang nyata)
--
--   `VARCHAR2(200 CHAR)` empat kali jalur terpanjang; `VARCHAR2(50 CHAR)`
--   lima kali kode mata uang terpanjang yang pernah ada. Kunci barunya
--   22 + 800 + 200 = 1.022 bita — muat pada ukuran blok mana pun.
--
--   ⚠️ Mempersempit AMAN hari ini dan hanya hari ini: kedua tabel **nol
--   baris**, diukur ulang sebelum berkas ini ditulis. Sesudah tiket 06
--   mengisinya, penyempitan menuntut pemindahan data.
--
-- =====================================================================
-- KENAPA `UNIQUE`, PADAHAL `Z00_KUNCI_ALAMI.sql` MENUNTUT NOMOR INVARIAN
-- =====================================================================
--   Aturan itu berbunyi *"constraint tanpa invarian tidak punya tempat
--   untuk gagal"*, dan kepala `426` memakainya untuk MENUNDA kunci alami
--   ini. Penundaan itu benar ketika tabelnya belum bertiket.
--
--   Kini ia bertiket, dan kedua tiket menuntutnya di daftar periksa mereka
--   sendiri. `SPEC-INVARIAN.md` berhenti di `INV-71`; kedua tiket
--   menyatakan apa adanya *"kunci alaminya belum punya nomor invarian"*
--   sebagai persyaratan yang diketahui, bukan sebagai penghalang.
--
--   ⛔ NOMORNYA TIDAK DIKARANG DI SINI. Memberi `INV-72`/`INV-73` dari
--   berkas migrasi berarti mengambil wewenang pemilik `SPEC-INVARIAN.md`,
--   dan nomor yang ditetapkan dua tempat akan bertabrakan. Yang ditulis di
--   sini ATURANNYA; nomornya ditagih — lihat `MODUL.md` bab kunci alami.
--
-- =====================================================================
-- ⚠️ SATU TUNTUTAN TIKET 77 YANG TIDAK DAPAT DIPENUHI
-- =====================================================================
--   Tiket 77 menuntut kunci asing kedua:
--
--     MATA_UANG  1--<  NILAI_SEBELUM_PRO_RATE   [hapus: tolak]
--
--   `MATA_UANG` **DICABUT** migrasi `434` (KEPUTUSAN §16, 4 Oktober 2026):
--   kurs dan daftar mata uang kini dibaca dari `TREATYEXCHANGEYEARLY`.
--   Kunci asing ke tabel yang tidak ada tidak dapat dipasang, dan
--   membangun ulang `MATA_UANG` dilarang §16.
--
--   Jadi `KODE_MATA_UANG` berdiri sebagai TEKS di dalam kunci alami, tanpa
--   kunci asing — setengah dari yang tiket 77 tuntut, dan setengahnya
--   DINYATAKAN alih-alih disamarkan. Pola dan alasannya sama dengan
--   `KODE_BESARAN` terhadap `BESARAN_DAPAT_DISESUAIKAN` (lihat kepala
--   `426`). Ditagih ke pemilik proses di `MODUL.md`.
--
-- Pembalikan: `443_kunci_alami_nilai_down.sql` melepas kedua constraint dan
-- mengembalikan lebar kolomnya. Aman selama kedua tabel nol baris.
ALTER TABLE {skema}.NILAI_SELISIH MODIFY (KODE_BESARAN VARCHAR2(200 CHAR))
/
ALTER TABLE {skema}.NILAI_SELISIH MODIFY (KODE_MATA_UANG VARCHAR2(50 CHAR))
/
ALTER TABLE {skema}.NILAI_SEBELUM_PRO_RATE MODIFY (KODE_BESARAN VARCHAR2(200 CHAR))
/
ALTER TABLE {skema}.NILAI_SEBELUM_PRO_RATE MODIFY (KODE_MATA_UANG VARCHAR2(50 CHAR))
/
ALTER TABLE {skema}.NILAI_SELISIH ADD CONSTRAINT UQ_NILAI_SELISIH
  UNIQUE (ID_VERSI_KONTRAK, KODE_BESARAN, KODE_MATA_UANG)
/
ALTER TABLE {skema}.NILAI_SEBELUM_PRO_RATE ADD CONSTRAINT UQ_NILAI_SEBELUM_PRO_RATE
  UNIQUE (ID_VERSI_KONTRAK, KODE_BESARAN, KODE_MATA_UANG)
/

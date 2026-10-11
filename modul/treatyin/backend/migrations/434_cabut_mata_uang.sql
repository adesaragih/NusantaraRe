-- MATA_UANG dan MATA_UANG_KONTRAK DICABUT.
--
-- Keputusan pemilik proses 4 Oktober 2026: kurs dan daftar mata uang diambil
-- dari `TREATYEXCHANGEYEARLY` -- tabel warisan yang sudah hidup, 140 baris,
-- 25 mata uang, terbagi per TREATYYEAR -- seperti yang layar lama lakukan.
-- Dua tabel model baru ini karena itu tidak dipakai lagi.
--
-- Pencabutan AMAN HARI INI, dan itu terukur: keduanya NOL BARIS pada
-- 4 Oktober 2026, dan kunci asing yang masuk hanya DUA -
-- `FK_VERSI_KONTRAK_MATA_UANG` dan `FK_MATA_UANG_KONTRAK_2`; yang kedua ikut
-- terbawa bersama tabelnya.
--
-- ---------------------------------------------------------------------
-- ⛔ SATU KEWAJIBAN YANG LAHIR, DAN IA BELUM DITUTUP: TIKET 57
-- ---------------------------------------------------------------------
--   Tiket 57 menuntut kurs DIBEKUKAN pada versi yang disetujui -- "angka
--   rupiah sebuah versi DISETUJUI sama hari ini dan tahun depan", dan jalur
--   gagalnya berbunyi: "ubah baris kurs tahunan lalu buka versi disetujui ->
--   angkanya TIDAK berubah".
--
--   `TREATYEXCHANGEYEARLY` adalah justru "baris kurs tahunan" itu: ia master
--   bersama yang boleh berubah. `MATA_UANG_KONTRAK.KURS` adalah tempat
--   pembekuannya, dan berkas ini mencabutnya.
--
--   Jadi sesudah migrasi ini, tiket 57 TIDAK PUNYA RUMAH. Itu dinyatakan di
--   sini, bukan ditemukan orang lain nanti: siapa pun yang mengerjakan 57
--   harus memutuskan lebih dulu di mana kurs dibekukan.
--
-- ⚠️ `VERSI_KONTRAK.KODE_MATA_UANG_KONTRAK` TETAP ADA, dan kini tanpa
--   penjaga. Kolomnya dulu menunjuk `MATA_UANG.ID_MATA_UANG`; sesudah ini ia
--   angka yang tidak dirujuk ke mana pun. Mengubah artinya menjadi kode mata
--   uang `TREATYEXCHANGEYEARLY` adalah pekerjaan tiket 20, bukan berkas ini -
--   dan `INV-44` yang dulu dijaga `FK_VERSI_KONTRAK_MATA_UANG` kini TIDAK
--   DITEGAKKAN di basis data.
--
-- Pembalikan: `434_cabut_mata_uang_down.sql` membangun keduanya kembali utuh
-- beserta kunci asingnya. Selama keduanya nol baris, pembalikan tidak
-- kehilangan apa pun.

ALTER TABLE {skema}.VERSI_KONTRAK DROP CONSTRAINT FK_VERSI_KONTRAK_MATA_UANG
/
DROP TABLE {skema}.MATA_UANG_KONTRAK CASCADE CONSTRAINTS
/
DROP TABLE {skema}.MATA_UANG CASCADE CONSTRAINTS
/
DROP SEQUENCE {skema}.SEQ_TRIN_MATA_UANG_KONTRAK
/
DROP SEQUENCE {skema}.SEQ_TRIN_MATA_UANG
/

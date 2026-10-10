-- NILAI_SELISIH dan NILAI_SEBELUM_PRO_RATE DICABUT DARI MODUL INI.
--
-- ⛔ BUKAN dibuang — PINDAH. Keduanya dibangun ulang utuh oleh
-- `modul/treatyinadjustment/backend/migrations/442_nilai_selisih.sql`,
-- dengan bentuk yang disalin apa adanya: presisi `NUMBER(38,8)`, nama
-- sequence, nama kunci asing, nama index. Migrasi yang membangun bentuk
-- berbeda bukan pemindahan.
--
-- Keputusan pemilik proses 4 Oktober 2026 — kepemilikan dijernihkan:
-- `docs/KEPUTUSAN-PENYELARASAN-REPO.md` §19. Tiket 76 dan 77 ada di papan
-- Adjustment, dan sejak ronde ini tabelnya ada di sana juga.
--
-- ---------------------------------------------------------------------
-- PENCABUTAN AMAN, DAN ITU DIUKUR ULANG SEBELUM BERKAS INI DITULIS
-- ---------------------------------------------------------------------
--   4 Oktober 2026, langsung ke POOLDATA:
--
--     SELECT COUNT(*) FROM POOLDATA.NILAI_SELISIH           -> 0
--     SELECT COUNT(*) FROM POOLDATA.NILAI_SEBELUM_PRO_RATE  -> 0
--
--   Nol baris berarti pemindahan tidak kehilangan apa pun. Bila salah
--   satunya kelak berisi, pencabutan ini TIDAK boleh diulang tanpa
--   memindahkan barisnya lebih dulu.
--
-- ---------------------------------------------------------------------
-- ⚠️ URUTAN MIGRASINYA MENGIKAT
-- ---------------------------------------------------------------------
--   Kedua tabel berkunci asing keluar ke `VERSI_KONTRAK`, yang milik modul
--   ini (migrasi 401). Rentang `440`+ modul Adjustment berjalan SESUDAH
--   rentang `400`–`439` modul ini, jadi `VERSI_KONTRAK` sudah berdiri
--   ketika `442` membangunnya kembali. Ketergantungan itu dinyatakan juga
--   di kepala `442`.
--
--   ⛔ Nomor `435` berada SESUDAH `426` yang membuatnya, dan itu wajib:
--   pelari migrasi berjalan menurut nomor, dan mencabut sebelum membuat
--   akan gagal pada pemasangan baru.
--
-- ---------------------------------------------------------------------
-- RIWAYATNYA TIDAK DIHAPUS
-- ---------------------------------------------------------------------
--   `426_nilai_selisih.sql` TETAP DI TEMPATNYA, lengkap dengan seluruh
--   alasan rancangannya — kardinalitas 1:N, penghalang
--   `BESARAN_DAPAT_DISESUAIKAN`, induk `VERSI_KONTRAK` lawan `KONTRAK`, dan
--   INV-58. Pola yang sama dipakai `434`: yang dicabut bendanya, bukan
--   catatannya. `442` menunjuk balik ke `426` untuk alasan-alasan itu
--   alih-alih menyalinnya dan membiarkan salinannya membeku.
--
-- Pembalikan: `435_cabut_nilai_selisih_down.sql` membangun keduanya kembali
-- di modul ini. Selama keduanya nol baris, pembalikan tidak kehilangan apa
-- pun — dan ia HARUS dijalankan bersama pembalikan `442`, sebab dua modul
-- yang sama-sama membuat tabel yang sama akan bertabrakan.
DROP TABLE {skema}.NILAI_SEBELUM_PRO_RATE CASCADE CONSTRAINTS
/
DROP TABLE {skema}.NILAI_SELISIH CASCADE CONSTRAINTS
/
DROP SEQUENCE {skema}.SEQ_TRIA_NILAI_SBL_PRORATA
/
DROP SEQUENCE {skema}.SEQ_TRIA_NILAI_SELISIH
/

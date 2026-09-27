-- Tahap kasus dan waktu buat pada work object - butir at dan au.
--
-- Pemilik: F0.4 (brief lanjutan 8 §2.2, `[DIPUTUSKAN 27-09-2026]`, diturunkan
-- dari XML; work owner dapat memveto sebelum migrasi dijalankan di Oracle
-- mana pun).
--
-- ⛔ BUTIR at - kenapa `PY_POSITION` tidak cukup.
--
-- `Register_Flow.xml` punya EMPAT assignment, dan keempatnya adalah KEADAAN
-- kasus: `Input Register` (baris 358), `Outstanding Claim` (343),
-- `Medical Check` (268), `Claim Analis` (313). Dua di antaranya dipegang
-- peran yang SAMA (`ReasLifeAdmin`), sehingga `PY_POSITION` tidak dapat
-- membedakannya - `models.TahapDariPeran` selama ini memetakan Admin SELALU
-- ke Outstanding dan mencatatnya `[terbuka]`.
--
-- Keadaan `Input Register` bukan teori: ia DAPAT DITUJU KEMBALI.
-- `[terverifikasi]` `Section/InputOSClaimLife.xml` baris 21404
-- `<pyLabel>Send Back to Register</pyLabel>`, baris 21433
-- `<pyLocalAction>SendtoAdmin</pyLocalAction>`. Tanpa kolom ini, kasus yang
-- dikembalikan ke Register tidak dapat dibedakan dari kasus yang sedang
-- Outstanding - dan kotak masuk Admin akan menyatukan dua antrian yang di
-- Pega terpisah.
--
-- ⚠️ Isinya nama assignment VERBATIM `pyTaskName`, bukan angka. Angka akan
-- menuntut peta kedua di suatu tempat, dan peta kedua adalah tempat kedua
-- untuk salah. `PY_POSITION` TETAP: ia peran pemegangnya (ADR-U-0002), dan
-- kedua kolom menjawab dua pertanyaan berbeda.
--
-- ⛔ BUTIR au - `TGL_UPDATE` bukan waktu buat.
--
-- `TGL_UPDATE` DITIMPA setiap perpindahan, sehingga sesudah satu serah terima
-- ia tidak lagi menyatakan kapan kasusnya lahir. Pega menyimpan keduanya
-- terpisah: `[terverifikasi]` `ReportDefinition/InboxPremiumList.xml` baris
-- 736 mengurutkan kotak masuk dengan `pxCreateDateTime`. Kotak masuk yang
-- diurutkan waktu UBAH akan melompat-lompat setiap kali seseorang memindah
-- kasus lain.
--
-- ⚠️ KEDUANYA NULLABLE, dan itu disengaja (ADR-U-0027): baris yang sudah ada
-- belum punya nilainya, dan kosong BERBEDA dari tebakan. Pengisian baris lama
-- dilakukan di KODE sebagai langkah bernama (migrasi data tiket 13), bukan di
-- DDL ini - `TAHAP` dari `PY_POSITION` lewat `TahapDariPeran`, `TGL_CREATE`
-- dari `TGL_UPDATE`, keduanya dicatat sebagai turunan, bukan sebagai fakta.
ALTER TABLE {skema}.T_WORK_CLAIM ADD (
  TAHAP      VARCHAR2(32),
  TGL_CREATE DATE
)
/

-- =====================================================================
-- Z02_KUNCI_PEMILIK.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- LANGKAH TERAKHIR URUTAN PEMASANGAN. Dijalankan sesudah:
--   1. 00_SKEMA_DAN_AKUN.sql        skema, akun, peran, hak sistem
--   2. 01_*.sql .. 22_*.sql         seluruh tabel
--   3. V01_*.sql .. V09_*.sql       seluruh view
--   4. V00_HAK_AKSES.sql            hak atas objek
--   5. Z01_PENGAWASAN_TULIS.sql     kebijakan pengawasan
--   6. >>> BERKAS INI <<<
--
-- AWALAN Z, BUKAN V. Di folder ini awalan V berarti VIEW — V01 sampai V09
-- adalah kesembilan view. Berkas ini bukan view, jadi ia memakai awalan
-- tersendiri yang juga menempatkannya di urutan terakhir.
--
-- Catatan atas V00_HAK_AKSES.sql: berkas itu juga bukan view, dan awalan V
-- padanya adalah warisan penamaan yang tidak konsisten. Ia TIDAK diganti
-- nama di sini karena sudah dirujuk di ADR-0028, tiket 01, dan tiket 35 —
-- mengganti nama berarti memutus ketiganya. Didaftarkan sebagai kerapian
-- yang ditunda, bukan sebagai cacat yang dilewatkan.
--
-- Menjalankannya lebih awal menghentikan pemasangan itu sendiri.
-- =====================================================================


-- ---------------------------------------------------------------------
-- Kenapa ini berkas tersendiri, bukan komentar di berkas lain
-- ---------------------------------------------------------------------
-- Versi pertama 00_SKEMA_DAN_AKUN.sql menaruh perintah ini sebagai baris
-- yang dikomentari, dengan niat "dinyalakan nanti". Itu dicabut.
--
-- Proyek ini sudah mendaftarkan cacat yang sama dua kali di sistem lama:
--
--   * Dasar ADR-0024 adalah blok yang dikomentari.
--   * POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP selalu INSERT dan tidak pernah
--     UPDATE, karena seluruh logika upsert-nya dikomentari.
--
-- Kode yang dimatikan dengan niat dinyalakan nanti tidak pernah dinyalakan.
-- Maka penguncian ini berdiri sebagai langkah pemasangan yang punya nomor
-- urut, bukan sebagai niat yang menunggu seseorang ingat.


-- ---------------------------------------------------------------------
-- Apa yang dikunci, dan apa yang tidak
-- ---------------------------------------------------------------------
-- DIKUNCI  : KLAIMNP        akun pemilik objek.
-- TIDAK    : KLAIMNP_APP    akun aplikasi; ia yang bekerja sehari-hari.
-- TIDAK    : KLAIMNP_HILIR  akun sistem hilir; SELECT atas view.
--
-- Mengunci akun pemilik TIDAK menghapus objeknya dan TIDAK mencabut hak
-- yang sudah diberikan. Ia hanya menutup pintu masuk akun itu. Tabel tetap
-- ada, view tetap jalan, dan KLAIMNP_APP tetap bekerja atas objek milik
-- akun yang terkunci itu.

ALTER USER KLAIMNP ACCOUNT LOCK;


-- ---------------------------------------------------------------------
-- Siapa yang membukanya kembali, dan dengan syarat apa
-- ---------------------------------------------------------------------
-- Yang membuka: DBA, bukan tim aplikasi, bukan pemegang KLAIMNP_APP.
--
-- Kapan boleh dibuka: HANYA untuk perubahan skema yang sudah punya berkas
-- DDL-nya sendiri dan sudah melewati ADR bila ia mengubah keputusan.
--
-- Berapa lama: satu perubahan, satu pembukaan, satu penguncian kembali.
-- Akun yang dibiarkan terbuka "sementara" adalah akun yang terbuka.
--
-- Cara membukanya: DBA menjalankan ALTER USER atas akun KLAIMNP dengan
-- klausa ACCOUNT UNLOCK, lalu menjalankan kembali berkas ini setelah
-- perubahan skemanya selesai.
--
-- Perintahnya sengaja ditulis sebagai kalimat, bukan sebagai baris SQL yang
-- dikomentari. Berkas pemasangan ini memuat NOL baris perintah yang menunggu
-- dinyalakan — dan itu justru alasan berkas ini ada.


-- ---------------------------------------------------------------------
-- Pemeriksaan sesudah dijalankan
-- ---------------------------------------------------------------------
-- Harapan: KLAIMNP berstatus LOCKED; kedua akun lain OPEN.
--
--   SELECT username, account_status, lock_date
--     FROM dba_users
--    WHERE username LIKE 'KLAIMNP%'
--    ORDER BY username;
--
-- Harapan: objek tetap milik KLAIMNP, jumlahnya tidak berubah oleh
-- penguncian.
--
--   SELECT owner, object_type, COUNT(*)
--     FROM dba_objects
--    WHERE owner = 'KLAIMNP'
--    GROUP BY owner, object_type
--    ORDER BY object_type;

-- =====================================================================
-- Akhir V01_KUNCI_PEMILIK.sql
-- =====================================================================

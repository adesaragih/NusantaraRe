-- =====================================================================
-- Z01_PENGAWASAN_TULIS.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Dijalankan sesudah V00_HAK_AKSES.sql, sebelum Z02_KUNCI_PEMILIK.sql.
--
-- Dasar: ADR-0017 (accepted), dengan rumusan klaimnya yang diperbaiki
-- 19 September 2026. Lihat badan ADR-0017 bagian "Apa yang benar-benar
-- dapat dijanjikan skema ini".
-- =====================================================================


-- ---------------------------------------------------------------------
-- Apa yang berkas ini lakukan, dan apa yang tidak
-- ---------------------------------------------------------------------
-- TIDAK MENCEGAH. Skema tidak dapat mencegah tulisan dari pemegang hak
-- sistem berakhiran ANY; pencabutan hak itu di tingkat instance, oleh DBA
-- (REQ-037).
--
-- MEMPERLIHATKAN. Setiap INSERT, UPDATE, dan DELETE atas tabel skema ini
-- yang datang dari akun SELAIN KLAIMNP_APP meninggalkan jejak.
--
-- Pembedaan itu bukan pelemahan klaim; ia klaim yang berbeda dan dapat
-- dipertahankan. "Satu pintu tulis" tanpa syarat tidak dapat dipertahankan
-- selama REQ-037 terbuka. "Tidak ada tulisan dari luar pintu yang tidak
-- meninggalkan jejak" dapat.


-- ---------------------------------------------------------------------
-- PRASYARAT YANG BELUM DIKETAHUI — dua, dan keduanya masuk REQ-037
-- ---------------------------------------------------------------------
-- 1. Unified Auditing aktif?
--    Oracle 12c dapat berjalan dalam mode campuran (mixed mode), dan dalam
--    mode itu kebijakan di bawah dibuat tetapi tidak selalu berlaku penuh.
--      SELECT value FROM v$option WHERE parameter = 'Unified Auditing';
--
-- 2. Siapa yang boleh membuat kebijakan pengawasan?
--    CREATE AUDIT POLICY menuntut peran AUDIT_ADMIN. Akun KLAIMNP
--    TIDAK memilikinya dan TIDAK diberi olehnya di 00_SKEMA_DAN_AKUN.sql —
--    memberikannya berarti pemilik skema dapat mematikan pengawasan atas
--    dirinya sendiri, dan pengawasan yang dapat dimatikan yang diawasi
--    bukan pengawasan.
--
-- Maka berkas ini DIJALANKAN DBA, bukan pemilik skema. Itu bukan kelemahan
-- rancangan; itu justru bentuknya.


-- ---------------------------------------------------------------------
-- Kebijakan
-- ---------------------------------------------------------------------
-- Nama kebijakan KLAIMNP_PENGAWASAN_TULIS = 24 byte. Ia bukan nama
-- constraint dan bukan nama index, jadi gerbang penamaan tidak menahannya.

CREATE AUDIT POLICY KLAIMNP_PENGAWASAN_TULIS
  ACTIONS
    INSERT ON KLAIMNP.KLAIM,                        UPDATE ON KLAIMNP.KLAIM,                        DELETE ON KLAIMNP.KLAIM,
    INSERT ON KLAIMNP.NILAI_KLAIM_MATA_UANG,        UPDATE ON KLAIMNP.NILAI_KLAIM_MATA_UANG,        DELETE ON KLAIMNP.NILAI_KLAIM_MATA_UANG,
    INSERT ON KLAIMNP.ALOKASI_LAYER,                UPDATE ON KLAIMNP.ALOKASI_LAYER,                DELETE ON KLAIMNP.ALOKASI_LAYER,
    INSERT ON KLAIMNP.RETENSI_CEDANT,               UPDATE ON KLAIMNP.RETENSI_CEDANT,               DELETE ON KLAIMNP.RETENSI_CEDANT,
    INSERT ON KLAIMNP.AKSEPTASI,                    UPDATE ON KLAIMNP.AKSEPTASI,                    DELETE ON KLAIMNP.AKSEPTASI,
    INSERT ON KLAIMNP.ADJUSTMENT,                   UPDATE ON KLAIMNP.ADJUSTMENT,                   DELETE ON KLAIMNP.ADJUSTMENT,
    INSERT ON KLAIMNP.PREMI_PEMULIHAN,              UPDATE ON KLAIMNP.PREMI_PEMULIHAN,              DELETE ON KLAIMNP.PREMI_PEMULIHAN,
    INSERT ON KLAIMNP.REKENING_PENERIMA,            UPDATE ON KLAIMNP.REKENING_PENERIMA,            DELETE ON KLAIMNP.REKENING_PENERIMA,
    INSERT ON KLAIMNP.OBJEK_PERTANGGUNGAN,          UPDATE ON KLAIMNP.OBJEK_PERTANGGUNGAN,          DELETE ON KLAIMNP.OBJEK_PERTANGGUNGAN,
    INSERT ON KLAIMNP.PEMBAGIAN_KERUGIAN,           UPDATE ON KLAIMNP.PEMBAGIAN_KERUGIAN,           DELETE ON KLAIMNP.PEMBAGIAN_KERUGIAN,
    INSERT ON KLAIMNP.PENYEBARAN,                   UPDATE ON KLAIMNP.PENYEBARAN,                   DELETE ON KLAIMNP.PENYEBARAN,
    INSERT ON KLAIMNP.ESTIMASI_AWAL,                UPDATE ON KLAIMNP.ESTIMASI_AWAL,                DELETE ON KLAIMNP.ESTIMASI_AWAL,
    INSERT ON KLAIMNP.KRONOLOGI_KLAIM,              UPDATE ON KLAIMNP.KRONOLOGI_KLAIM,              DELETE ON KLAIMNP.KRONOLOGI_KLAIM,
    INSERT ON KLAIMNP.DOKUMEN_KLAIM,                UPDATE ON KLAIMNP.DOKUMEN_KLAIM,                DELETE ON KLAIMNP.DOKUMEN_KLAIM,
    INSERT ON KLAIMNP.TUTUP_BUKU,                   UPDATE ON KLAIMNP.TUTUP_BUKU,                   DELETE ON KLAIMNP.TUTUP_BUKU,
    INSERT ON KLAIMNP.TARIF_BERLAKU,                UPDATE ON KLAIMNP.TARIF_BERLAKU,                DELETE ON KLAIMNP.TARIF_BERLAKU,
    INSERT ON KLAIMNP.KOREKSI_NILAI,                UPDATE ON KLAIMNP.KOREKSI_NILAI,                DELETE ON KLAIMNP.KOREKSI_NILAI,
    INSERT ON KLAIMNP.KLAIM_PENJAGA_TANGGAL, UPDATE ON KLAIMNP.KLAIM_PENJAGA_TANGGAL, DELETE ON KLAIMNP.KLAIM_PENJAGA_TANGGAL,
    INSERT ON KLAIMNP.MIGRASI_KORELASI,             UPDATE ON KLAIMNP.MIGRASI_KORELASI,             DELETE ON KLAIMNP.MIGRASI_KORELASI,
    INSERT ON KLAIMNP.MIGRASI_NILAI_DITOLAK,        UPDATE ON KLAIMNP.MIGRASI_NILAI_DITOLAK,        DELETE ON KLAIMNP.MIGRASI_NILAI_DITOLAK,
    INSERT ON KLAIMNP.MIGRASI_PENDARATAN,           UPDATE ON KLAIMNP.MIGRASI_PENDARATAN,           DELETE ON KLAIMNP.MIGRASI_PENDARATAN,
    INSERT ON KLAIMNP.ARSIP_MUATAN_KELUAR,          UPDATE ON KLAIMNP.ARSIP_MUATAN_KELUAR,          DELETE ON KLAIMNP.ARSIP_MUATAN_KELUAR
  WHEN 'SYS_CONTEXT(''USERENV'', ''SESSION_USER'') <> ''KLAIMNP_APP'''
  EVALUATE PER STATEMENT;

AUDIT POLICY KLAIMNP_PENGAWASAN_TULIS;


-- ---------------------------------------------------------------------
-- Kenapa syaratnya SESSION_USER, bukan CURRENT_USER
-- ---------------------------------------------------------------------
-- Ini bukan kehalusan; ia justru inti berkas ini.
--
-- CURRENT_USER adalah pemilik kode yang sedang berjalan. Di dalam prosedur
-- definer's rights milik KLAIMNP_APP, CURRENT_USER menjadi KLAIMNP_APP
-- siapa pun yang memanggilnya — dan jejaknya hilang.
--
-- SESSION_USER adalah akun yang benar-benar masuk. Ia tidak berubah oleh
-- prosedur siapa pun.
--
-- Ini bukan kemungkinan teoretis. Sistem lama memakai pola itu:
-- POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP adalah prosedur yang menulis atas
-- nama pemiliknya. Siapa pun yang memegang EXECUTE atasnya menulis sebagai
-- POOLDATA. Memakai CURRENT_USER berarti membangun ulang lubang itu di
-- dalam pengawasan yang seharusnya menutupinya.


-- ---------------------------------------------------------------------
-- Apa yang TIDAK tertangkap kebijakan ini
-- ---------------------------------------------------------------------
-- Ditulis supaya tidak dibaca sebagai jaring yang rapat:
--
--   * SELECT tidak diawasi. Kebijakan ini tentang perubahan, bukan
--     pembacaan. Pembacaan oleh akun tak berhak adalah pokok tersendiri.
--   * Pemegang AUDIT_ADMIN dapat mematikan kebijakan ini. Itu sebabnya
--     KLAIMNP tidak diberi peran itu, dan sebabnya berkas ini dijalankan
--     DBA.
--   * Perubahan lewat DDL — TRUNCATE, DROP, ALTER — tidak tercakup
--     ACTIONS di atas. Bila diperlukan, ia kebijakan kedua.
--   * Bila Unified Auditing tidak aktif, kebijakan ini dibuat tetapi tidak
--     berlaku penuh. Lihat prasyarat 1.


-- ---------------------------------------------------------------------
-- Pemeriksaan
-- ---------------------------------------------------------------------
-- Kebijakan ada dan menyala:
--   SELECT policy_name, enabled_option, entity_name
--     FROM audit_unified_enabled_policies
--    WHERE policy_name = 'KLAIMNP_PENGAWASAN_TULIS';
--
-- Jejak tulisan dari luar pintu — harapan: NOL baris.
--   SELECT event_timestamp, dbusername, action_name, object_name
--     FROM unified_audit_trail
--    WHERE object_schema = 'KLAIMNP'
--      AND action_name IN ('INSERT','UPDATE','DELETE')
--      AND dbusername <> 'KLAIMNP_APP'
--    ORDER BY event_timestamp DESC;
--
-- Baris apa pun di sini adalah tulisan yang melewati pintu, dan setiap
-- barisnya menyebut siapa, kapan, dan atas tabel mana.

-- =====================================================================
-- Akhir Z01_PENGAWASAN_TULIS.sql
-- =====================================================================

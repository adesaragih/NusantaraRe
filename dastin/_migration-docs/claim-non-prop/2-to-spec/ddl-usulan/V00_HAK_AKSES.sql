-- =====================================================================
-- V00_HAK_AKSES.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 35. ADR-0017 + ADR-0028: satu pintu tulis, ditegakkan lewat HAK
-- AKSES — bukan lewat kesepakatan.
--
-- PRASYARAT: 00_SKEMA_DAN_AKUN.sql lebih dulu (skema, tiga akun, tiga peran,
-- hak SISTEM), lalu KEDUA PULUH DUA berkas tabel. Berkas ini mengurus hak atas
-- OBJEK, dan objeknya harus sudah ada.
--
-- CATATAN URUTAN — dikoreksi 19 September 2026. Versi sebelumnya berbunyi
-- "sesudah seluruh tabel DAN VIEW terpasang". Itu keliru, dan keliru dengan
-- cara yang berbahaya: namanya V00, jadi ia berjalan SEBELUM V01..V09 pada
-- urutan abjad mana pun. Untungnya ia memang tidak memerlukan view — seluruh
-- REVOKE di bawah hanya menyentuh TABEL, dan hibah SELECT atas view diberikan
-- di dalam berkas view masing-masing.
--
-- Yang benar: 00 → 01..22 → V00 → V01..V09 → Z00 → Z01 → Z02. Urutan abjad
-- berkas SUDAH menghasilkan urutan itu; tidak ada langkah yang perlu diingat
-- di luar mengurutkan namanya.
--
-- TIDAK ADA DML DI BERKAS INI.
-- =====================================================================


-- ---------------------------------------------------------------------
-- BATAS YANG HARUS DIBACA BERSAMA BERKAS INI
-- ---------------------------------------------------------------------
-- Seluruh REVOKE di bawah berlaku pada tingkat OBJEK. Hak SISTEM berakhiran
-- ANY — SELECT ANY TABLE, INSERT ANY TABLE, UPDATE ANY TABLE, DELETE ANY
-- TABLE, ALTER ANY TABLE — MENGATASI hak objek dan TIDAK TERSENTUH satu
-- baris pun di sini.
--
-- Dan bukan hanya itu. REQ-037 mendaftar ENAM jalur tulis yang melewati hak
-- objek: hak ANY; prosedur definer's rights milik akun berhak tulis; hibah
-- ke PUBLIC; CREATE ANY TRIGGER; peran yang MEMUAT hak ANY; serta
-- GRANT ANY OBJECT PRIVILEGE dan GRANT ANY PRIVILEGE.
--
-- Jalur kedua BUKAN kemungkinan teoretis — ia pola yang SUDAH DIPAKAI di
-- instance ini: POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP menulis atas nama
-- pemiliknya, dan siapa pun yang memegang EXECUTE atasnya menulis sebagai
-- POOLDATA.
--
-- Berkas ini menutup PINTU OBJEK, bukan seluruh pintu. Yang menutup sisanya
-- ada di tingkat instance, oleh DBA, dan itu REQ-037.
--
-- Yang dapat dijanjikan sekarang, dan yang TIDAK:
--   TIDAK: "tidak ada tulisan dari luar pintu."
--   YA   : "tidak ada tulisan dari luar pintu yang tidak meninggalkan
--           jejak" — ditegakkan Z01_PENGAWASAN_TULIS.sql.
-- Lihat badan ADR-0017; rumusan klaimnya sudah diperbaiki 19 Sep 2026.


-- ---------------------------------------------------------------------
-- KENAPA REVOKE-nya DIBUNGKUS, dan bukan 44 baris polos
-- ---------------------------------------------------------------------
-- Versi sebelumnya menulis 44 baris REVOKE ALL polos. Di instance BERSIH,
-- tidak satu pun hak itu pernah diberikan — dan Oracle MENOLAK pencabutan
-- hak yang tidak ada: ORA-01927 (cannot REVOKE privileges you did not
-- grant). Baris pertama gagal, dan pemasangan berhenti di sana.
--
-- Ini KELAS CACAT YANG SAMA dengan sembilan REVOKE di 00_SKEMA_DAN_AKUN.sql
-- yang sudah diperbaiki lebih dulu, dan ia terulang karena keduanya ditulis
-- dengan anggapan yang sama: bahwa mencabut sesuatu yang tidak ada tidak
-- berakibat apa-apa. Di Oracle, itu keliru.
--
-- Bentuk di bawah mencabut bila ada, melewati bila tidak ada, dan MERAMBAT
-- untuk galat apa pun selain itu — sehingga kegagalan yang sesungguhnya
-- tetap menghentikan pemasangan.

BEGIN
  FOR r IN (
    SELECT 'REVOKE ALL ON KLAIMNP.' || tabel || ' FROM ' || akun AS perintah
      FROM ( SELECT 'KLAIM' AS tabel FROM dual
                 UNION ALL SELECT 'NILAI_KLAIM_MATA_UANG'
                 UNION ALL SELECT 'ALOKASI_LAYER'
                 UNION ALL SELECT 'RETENSI_CEDANT'
                 UNION ALL SELECT 'AKSEPTASI'
                 UNION ALL SELECT 'ADJUSTMENT'
                 UNION ALL SELECT 'PREMI_PEMULIHAN'
                 UNION ALL SELECT 'REKENING_PENERIMA'
                 UNION ALL SELECT 'OBJEK_PERTANGGUNGAN'
                 UNION ALL SELECT 'PEMBAGIAN_KERUGIAN'
                 UNION ALL SELECT 'PENYEBARAN'
                 UNION ALL SELECT 'ESTIMASI_AWAL'
                 UNION ALL SELECT 'KRONOLOGI_KLAIM'
                 UNION ALL SELECT 'DOKUMEN_KLAIM'
                 UNION ALL SELECT 'TUTUP_BUKU'
                 UNION ALL SELECT 'TARIF_BERLAKU'
                 UNION ALL SELECT 'KOREKSI_NILAI'
                 UNION ALL SELECT 'KLAIM_PENJAGA_TANGGAL'
                 UNION ALL SELECT 'MIGRASI_KORELASI'
                 UNION ALL SELECT 'MIGRASI_NILAI_DITOLAK'
                 UNION ALL SELECT 'MIGRASI_PENDARATAN'
                 UNION ALL SELECT 'ARSIP_MUATAN_KELUAR' ),
           ( SELECT 'POOLDATA' AS akun FROM dual
                 UNION ALL SELECT 'KLAIMNP_HILIR' )
  ) LOOP
    BEGIN
      EXECUTE IMMEDIATE r.perintah;
    EXCEPTION
      WHEN OTHERS THEN
        IF SQLCODE = -1927 THEN NULL; ELSE RAISE; END IF;
    END;
  END LOOP;
END;
/

-- 22 tabel x 2 akun = 44 pencabutan. Yang tersisa bagi hilir hanyalah
-- SELECT atas view, diberikan di berkas V01..V09.


-- ---------------------------------------------------------------------
-- SATU PENGECUALIAN, DAN IA DISEBUT DI SINI SUPAYA TERBACA DI SATU TEMPAT
-- ---------------------------------------------------------------------
-- Akun aplikasi memegang SELECT, INSERT, UPDATE, DELETE atas 21 tabel.
-- Atas ARSIP_MUATAN_KELUAR ia memegang SELECT dan INSERT SAJA — TANPA
-- UPDATE, TANPA DELETE. Hibahnya ada di 22_ARSIP_MUATAN_KELUAR.sql; yang
-- ditulis di sini hanya FAKTANYA, supaya pengecualian itu tidak hanya hidup
-- di satu berkas tabel yang tidak dibaca orang saat meninjau hak akses.
--
-- Sebabnya: arsip menjawab "apa yang BENAR-BENAR dikirim". Arsip yang dapat
-- diubah menjawab pertanyaan lain, dan itu bukan catatan. Dan barisnya
-- TAMBAHAN, bukan keadaan — satu UPDATE mengubah setiap selisih yang pernah
-- dihitung sesudahnya (tiket 25, tiket 29).


-- ---------------------------------------------------------------------
-- PEMERIKSAAN — dijalankan sesudah berkas ini
-- ---------------------------------------------------------------------
-- 1. Hak objek yang tersisa bagi POOLDATA dan hilir atas TABEL.
--    Harapan: NOL baris.
--      SELECT grantee, table_name, privilege
--        FROM dba_tab_privs
--       WHERE owner = 'KLAIMNP'
--         AND grantee IN ('POOLDATA','KLAIMNP_HILIR')
--         AND table_name IN (SELECT table_name FROM dba_tables
--                             WHERE owner = 'KLAIMNP');
--
-- 2. Hak akun aplikasi atas arsip. Harapan: TEPAT dua baris, INSERT dan
--    SELECT. Baris UPDATE atau DELETE di sini adalah kegagalan.
--      SELECT privilege FROM dba_tab_privs
--       WHERE owner = 'KLAIMNP' AND table_name = 'ARSIP_MUATAN_KELUAR'
--         AND grantee = 'KLAIMNP_APP';
--
-- 3. Yang TIDAK dapat diperiksa dari sini: keenam jalur REQ-037. Ketiadaan
--    baris pada pemeriksaan 1 TIDAK berarti tidak ada yang dapat menulis.
-- =====================================================================

-- =====================================================================
-- 00_SKEMA_DAN_AKUN.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 01. Dasar: ADR-0028 (accepted 19 Sep 2026), ADR-0017 (accepted).
-- Dijalankan PALING AWAL — sebelum 01_KLAIM.sql dan seluruh berkas lain.
-- Pasangannya di ujung lain: V00_HAK_AKSES.sql, yang mengurus hak atas
-- OBJEK dan baru dapat jalan sesudah objeknya ada (tiket 35).
--
-- PENAMAAN: FINAL. Gerbang dicabut 19 September 2026.
-- Aturannya SPEC-MODEL-DATA.md bagian 16: kata utuh bahasa Indonesia, batas
-- 30 byte sebagai ATURAN TETAP (sah di setiap versi Oracle, bukan pengamanan
-- sementara sampai COMPATIBLE diketahui). REQ-032 turun jadi verifikasi.
--
-- Pengenal di berkas ini, dihitung:
--   KLAIMNP                     7    skema dan pemilik
--   KLAIMNP_APP                11    akun aplikasi
--   KLAIMNP_HILIR              13    akun baca hilir
--   KLAIMNP_DATA               12    tablespace
--   KLAIMNP_INDEKS             14    tablespace
--   KLAIMNP_PERAN_HILIR        19    peran
--   KLAIMNP_PERAN_APLIKASI     22    peran
--   KLAIMNP_PERAN_PEMASANGAN   24    peran  <- terpanjang
-- Di atas 30 byte: NOL.
--
-- DUA NAMA YANG BUKAN KATA UTUH, DICATAT TERBUKA:
--   KLAIMNP  memampatkan "Klaim Non Prop".
--   _APP     memampatkan "aplikasi".
-- Keduanya nama SKEMA dan AKUN, bukan nama tabel maupun kolom, sehingga
-- aturan 1 tidak menjangkaunya. Keduanya juga sudah dirujuk ADR-0017,
-- ADR-0028, seluruh GRANT di 22 berkas tabel, kesembilan view, dan
-- kebijakan pengawasan Z01 — menggantinya sekarang menyentuh lebih banyak
-- daripada yang diperbaikinya. Dicatat supaya pilihannya terbaca, bukan
-- supaya terlupa: bila kelak nama skema disentuh, KLAIMNP_APLIKASI dan
-- KLAIMNP_HILIR lebih konsisten, dan keduanya muat.
--
-- Berkas ini tidak memuat nama constraint maupun nama index — bukan karena
-- gerbang, melainkan karena tidak ada objek yang membutuhkannya di sini.
--
-- TIDAK ADA DML DI BERKAS INI. Tidak ada INSERT, UPDATE, maupun DELETE.
-- =====================================================================


-- ---------------------------------------------------------------------
-- BAGIAN 0 — yang harus diputuskan DBA sebelum berkas ini berarti
-- ---------------------------------------------------------------------
-- Tiga hal sengaja TIDAK ditulis di sini, karena menebaknya lebih buruk
-- daripada mengosongkannya:
--
--   1. Letak dan ukuran datafile. Ditulis sebagai penanda &&.
--   2. Kata sandi ketiga akun. TIDAK PERNAH ditulis di berkas mana pun,
--      termasuk berkas ini. Ditulis sebagai penanda &&, diisi saat
--      pemasangan, dan tidak disimpan sesudahnya.
--   3. Profil kata sandi (PROFILE). Kebijakan organisasi, bukan kebijakan
--      modul ini.
--
-- Penanda && adalah variabel substitusi SQL*Plus. Bila berkas ini kelak
-- benar-benar dijalankan, ia akan meminta ketiganya saat itu juga.


-- ---------------------------------------------------------------------
-- BAGIAN 1 — tablespace
-- ---------------------------------------------------------------------
-- Dipisah dua supaya data dan index tidak bersaing pada I/O yang sama.
-- Nama KLAIMNP_INDEKS adalah nama TABLESPACE, bukan nama index. Ia memakai
-- kata utuh bahasa Indonesia (INDEKS), sesuai aturan 1.

CREATE TABLESPACE KLAIMNP_DATA
  DATAFILE '&&letak_datafile_data'
  SIZE &&ukuran_awal_data
  AUTOEXTEND ON NEXT &&pertambahan_data MAXSIZE &&batas_data
  EXTENT MANAGEMENT LOCAL
  SEGMENT SPACE MANAGEMENT AUTO;

CREATE TABLESPACE KLAIMNP_INDEKS
  DATAFILE '&&letak_datafile_indeks'
  SIZE &&ukuran_awal_indeks
  AUTOEXTEND ON NEXT &&pertambahan_indeks MAXSIZE &&batas_indeks
  EXTENT MANAGEMENT LOCAL
  SEGMENT SPACE MANAGEMENT AUTO;


-- ---------------------------------------------------------------------
-- BAGIAN 2 — tiga akun, dan alasan masing-masing terpisah
-- ---------------------------------------------------------------------
-- KLAIMNP        pemilik objek. Membuat tabel dan view saat pemasangan,
--                lalu DIKUNCI (bagian 6). Tidak dipakai aplikasi.
-- KLAIMNP_APP    akun aplikasi. Tidak memiliki objek apa pun dan tidak
--                dapat membuatnya. Hanya menyentuh objek milik KLAIMNP,
--                lewat hak yang diberikan V00_HAK_AKSES.sql.
-- KLAIMNP_HILIR  akun sistem hilir. SELECT atas view saja, dan itu pun
--                baru diberikan sesudah view-nya ada.
--
-- Pemisahan pemilik dari aplikasi bukan formalitas. Bila keduanya satu
-- akun, aplikasi dapat ALTER dan DROP tabelnya sendiri, dan ADR-0017
-- kehilangan artinya: pintu tulis kedua tidak perlu dicari dari luar,
-- ia sudah ada di dalam.

CREATE USER KLAIMNP
  IDENTIFIED BY "&&kata_sandi_pemilik"
  DEFAULT TABLESPACE KLAIMNP_DATA
  TEMPORARY TABLESPACE TEMP
  QUOTA UNLIMITED ON KLAIMNP_DATA
  QUOTA UNLIMITED ON KLAIMNP_INDEKS;

CREATE USER KLAIMNP_APP
  IDENTIFIED BY "&&kata_sandi_aplikasi"
  DEFAULT TABLESPACE KLAIMNP_DATA
  TEMPORARY TABLESPACE TEMP
  QUOTA 0 ON KLAIMNP_DATA
  QUOTA 0 ON KLAIMNP_INDEKS;

CREATE USER KLAIMNP_HILIR
  IDENTIFIED BY "&&kata_sandi_hilir"
  DEFAULT TABLESPACE KLAIMNP_DATA
  TEMPORARY TABLESPACE TEMP
  QUOTA 0 ON KLAIMNP_DATA
  QUOTA 0 ON KLAIMNP_INDEKS;

-- QUOTA 0 pada kedua akun bukan kehati-hatian berlebih. Tanpa itu, akun
-- aplikasi dapat membuat segmen di tablespace ini — yaitu memiliki objek
-- yang tidak pernah masuk DDL mana pun dan tidak pernah terlihat siapa pun.


-- ---------------------------------------------------------------------
-- BAGIAN 3 — peran
-- ---------------------------------------------------------------------
-- Hak diberikan ke PERAN, bukan langsung ke akun. Alasannya audit:
-- satu pertanyaan "peran ini boleh apa" menjawab seluruh akun yang
-- memegangnya, dan penambahan akun kelak tidak menambah tempat memeriksa.

CREATE ROLE KLAIMNP_PERAN_PEMASANGAN;
CREATE ROLE KLAIMNP_PERAN_APLIKASI;
CREATE ROLE KLAIMNP_PERAN_HILIR;


-- ---------------------------------------------------------------------
-- BAGIAN 4 — hak SISTEM, seminimal yang dapat bekerja
-- ---------------------------------------------------------------------
-- Hak atas OBJEK tidak ada di sini. Ia di V00_HAK_AKSES.sql, dan baru
-- berarti sesudah objeknya berdiri.

-- Peran pemasangan: membuat objek. Hanya itu, dan hanya selama pemasangan.
GRANT CREATE SESSION   TO KLAIMNP_PERAN_PEMASANGAN;
GRANT CREATE TABLE     TO KLAIMNP_PERAN_PEMASANGAN;
GRANT CREATE VIEW      TO KLAIMNP_PERAN_PEMASANGAN;
GRANT CREATE SEQUENCE  TO KLAIMNP_PERAN_PEMASANGAN;

-- Peran aplikasi: MASUK SAJA. Tidak ada CREATE apa pun.
-- Aplikasi tidak pernah membuat objek. Bila kelak ia perlu, itu perubahan
-- desain yang harus melewati ADR, bukan penambahan satu baris di sini.
GRANT CREATE SESSION TO KLAIMNP_PERAN_APLIKASI;

-- Peran hilir: MASUK SAJA.
GRANT CREATE SESSION TO KLAIMNP_PERAN_HILIR;

GRANT KLAIMNP_PERAN_PEMASANGAN TO KLAIMNP;
GRANT KLAIMNP_PERAN_APLIKASI   TO KLAIMNP_APP;
GRANT KLAIMNP_PERAN_HILIR      TO KLAIMNP_HILIR;

ALTER USER KLAIMNP       DEFAULT ROLE ALL;
ALTER USER KLAIMNP_APP   DEFAULT ROLE ALL;
ALTER USER KLAIMNP_HILIR DEFAULT ROLE ALL;


-- ---------------------------------------------------------------------
-- BAGIAN 5 — yang TIDAK diberikan, ditulis supaya dapat diaudit
-- ---------------------------------------------------------------------
-- Hak yang tidak tertulis tidak dapat diaudit. Maka ketiadaan pun ditulis.
--
-- TIDAK diberikan ke akun mana pun di skema ini:
--
--   CONNECT, RESOURCE    peran bawaan Oracle yang isinya berubah antarversi.
--                        Memakainya berarti hak akun ini ditentukan versi
--                        basis data, bukan ditentukan dokumen ini.
--   UNLIMITED TABLESPACE ikut terbawa RESOURCE, dan membatalkan QUOTA 0.
--   DBA, SYSDBA          tidak perlu dijelaskan.
--   CREATE ANY *         lihat BAGIAN 5b.
--   SELECT ANY TABLE     lihat BAGIAN 5b.
--   CREATE DATABASE LINK dilarang ADR-0016. Lihat catatan di bawah.
--
-- KOREKSI 19 Sep 2026. Versi pertama berkas ini menulis kesembilan REVOKE
-- sebagai perintah telanjang, dengan catatan "pada instance bersih ia tidak
-- melakukan apa-apa". ITU SALAH: Oracle mengangkat galat, bukan mengabaikan.
--   ORA-01951  ROLE 'CONNECT' not granted to ...
--   ORA-01952  system privileges not granted to ...
-- Pada instance bersih — yaitu keadaan yang justru diharapkan — baris
-- pertama menggugurkan seluruh pemasangan.
--
-- Bentuk di bawah mencabut bila ada dan diam bila tidak ada. Ia menelan
-- HANYA kedua galat itu; galat lain tetap diangkat.

BEGIN
  FOR r IN (
    SELECT 'REVOKE ' || hak || ' FROM ' || akun AS perintah
      FROM ( SELECT 'CONNECT'              AS hak FROM dual
             UNION ALL SELECT 'RESOURCE'
             UNION ALL SELECT 'UNLIMITED TABLESPACE'
             UNION ALL SELECT 'CREATE DATABASE LINK' ),
           ( SELECT 'KLAIMNP'       AS akun FROM dual
             UNION ALL SELECT 'KLAIMNP_APP'
             UNION ALL SELECT 'KLAIMNP_HILIR' )
  ) LOOP
    BEGIN
      EXECUTE IMMEDIATE r.perintah;
    EXCEPTION
      WHEN OTHERS THEN
        IF SQLCODE IN (-1951, -1952) THEN NULL; ELSE RAISE; END IF;
    END;
  END LOOP;
END;
/

-- CREATE DATABASE LINK tidak diberikan kepada satu akun pun di sini.
-- ADR-0016 melarangnya. Perlu dicatat terang: instance ini SUDAH memuat
-- satu database link yang berjalan produksi — hrdasm.v_hrd_mst@asmd...
-- dipakai POOLDATA.V_MST_USER_TEKNIS (D20). Larangan ini berlaku untuk
-- skema baru; nasib yang lama adalah K1, dan K1 belum diputuskan.
-- Pencabutannya ikut di dalam blok di atas, bersama ketiga hak lain,
-- dengan penanganan galat yang sama.


-- ---------------------------------------------------------------------
-- BAGIAN 5b — LUBANG YANG TIDAK DAPAT DITUTUP DARI SINI
-- ---------------------------------------------------------------------
-- Ini bagian terpenting berkas ini, dan ia tidak memuat satu perintah pun.
--
-- V00_HAK_AKSES.sql mencabut hak POOLDATA atas tiap tabel, satu per satu.
-- Pencabutan itu TIDAK BERARTI APA-APA bila POOLDATA — atau akun lain —
-- memegang hak sistem berakhiran ANY:
--
--     SELECT ANY TABLE     membaca seluruh tabel di seluruh skema
--     INSERT ANY TABLE     menulis ke seluruh tabel di seluruh skema
--     UPDATE ANY TABLE     ...
--     DELETE ANY TABLE     ...
--     ALTER ANY TABLE      ...
--
-- Hak ANY mengatasi hak objek. Satu akun dengan INSERT ANY TABLE membuat
-- seluruh ADR-0017 tidak berlaku, dan tidak ada satu baris pun di skema
-- ini yang dapat mencegahnya — hak itu dicabut di tingkat instance, oleh
-- DBA, bukan oleh pemilik skema.
--
-- Selama ini belum pernah diperiksa. Tidak disebut di ADR-0017, tidak di
-- ADR-0028, tidak di V00_HAK_AKSES.sql, dan tidak di REQ mana pun.
-- Diajukan sebagai REQ-037.
--
-- Sampai REQ-037 kembali, klaim "satu pintu tulis" berlaku pada tingkat
-- OBJEK saja, dan itu harus dikatakan apa adanya.


-- ---------------------------------------------------------------------
-- BAGIAN 6 — akun pemilik dikunci sesudah pemasangan: BERKAS TERSENDIRI
-- ---------------------------------------------------------------------
-- Penguncian TIDAK ada di berkas ini, dan tidak ada di sini sebagai
-- komentar. Ia di V01_KUNCI_PEMILIK.sql, langkah terakhir urutan pemasangan.
--
-- Alasannya sejarah proyek ini sendiri. Dua kali sudah tercatat bahwa kode
-- yang dimatikan dengan niat dinyalakan nanti tidak pernah dinyalakan:
--   • dasar ADR-0024 adalah blok yang dikomentari;
--   • PEGA_JSON_OS_AKSEP_KLAIMTNP selalu INSERT dan tidak pernah UPDATE,
--     karena logika upsert-nya dikomentari seluruhnya.
-- Berkas pemasangan ini karena itu memuat NOL baris SQL yang dikomentari
-- dengan maksud dijalankan kelak.


-- ---------------------------------------------------------------------
-- BAGIAN 7 — pemeriksaan penerimaan
-- ---------------------------------------------------------------------
-- Seluruhnya SELECT. Tidak mengubah apa pun.
-- Ketiganya menjawab kriteria terima tiket 01.

-- (1) Akun aplikasi dapat membuat sesi dan melihat skema.
--     Harapan: satu baris, PRIVILEGE = CREATE SESSION.
--
-- SELECT grantee, privilege
--   FROM dba_sys_privs
--  WHERE grantee IN ('KLAIMNP_PERAN_APLIKASI')
--  ORDER BY privilege;

-- (2) Akun hilir tidak melihat satu objek pun.
--     Harapan: NOL baris.
--
-- SELECT owner, table_name, privilege
--   FROM dba_tab_privs
--  WHERE grantee = 'KLAIMNP_HILIR';

-- (3) Akun pemilik terpisah dari akun aplikasi.
--     Harapan: TIGA baris dengan USERNAME berbeda, dan
--     KLAIMNP_APP tidak memiliki satu objek pun.
--
-- SELECT username, account_status, default_tablespace
--   FROM dba_users
--  WHERE username LIKE 'KLAIMNP%'
--  ORDER BY username;
--
-- SELECT owner, COUNT(*) AS jumlah_objek
--   FROM dba_objects
--  WHERE owner LIKE 'KLAIMNP%'
--  GROUP BY owner;

-- (4) Tidak ada akun skema ini yang memegang hak berakhiran ANY.
--     Harapan: NOL baris.
--
-- SELECT grantee, privilege
--   FROM dba_sys_privs
--  WHERE grantee LIKE 'KLAIMNP%'
--    AND privilege LIKE '% ANY %';

-- (5) REQ-037 — siapa pun di instance ini yang memegang hak ANY.
--     Ini BUKAN pemeriksaan skema ini; ia pemeriksaan instance, dan
--     hasilnya menentukan apakah ADR-0017 berlaku sungguhan.
--     Harapan: tidak ada harapan. Yang ada hanya kenyataan.
--
-- SELECT grantee, privilege
--   FROM dba_sys_privs
--  WHERE privilege LIKE '%ANY TABLE'
--     OR privilege LIKE '%ANY INDEX'
--  ORDER BY grantee, privilege;

-- =====================================================================
-- Akhir 00_SKEMA_DAN_AKUN.sql
-- =====================================================================

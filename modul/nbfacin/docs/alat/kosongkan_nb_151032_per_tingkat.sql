-- =====================================================================================================================
-- Kosongkan NB-151032 - versi PER TINGKAT (sesi 0f, 11-10-2026), untuk data sebesar hasil `D:\migrasi\RNM\DDL\DATA NB.csv`
-- (BAGIAN 0.3 kosongkan_nb_151032.sql dari work owner): HANYA pohon objek FIRE yang berisi - T_LOCATIONLIST 1.010.688,
-- T_PROPERTY / T_RISKLOCATION / T_BUILDINGCONSTRUCTION / T_SURROUNDINGRISK / T_PROPERTYITEMLIST / T_OCCUPATIONLIST /
-- T_TABLEOFLIMIT / T_LISTCAUSEOFLOSS / T_COINSDATA masing-masing 1.010.688, T_COVERAGELIST 4.994.304, T_NB_CSV_BARIS 1.010.688,
-- T_NB_CSV_UNGGAH 1; semua tabel lain 0. Total ± 16 juta baris.
--
-- Kenapa versi ini (bukan 1A / 1B kosongkan_nb_151032.sql):
--   * 1A satu transaksi ± 16 juta baris -> hampir pasti melebihi undo DEV;
--   * setiap TINGKAT dihapus untuk SELURUH case (cucu, lalu anak, lalu induk), dipotong per v_potong lokasi (rentang SEQ_NO, indeks
--     IX_LOCATIONLIST_PARENT) dengan COMMIT per potongan - undo per transaksi terbatas, dan tabel anak di setiap tingkat dipindai
--     sekali per potongan (hash join), bukan per lokasi.
--
-- ⛔ PRASYARAT: MIGRASI 259 (indeks FK, keputusan work owner 11-10-2026) SUDAH DIJALANKAN. RALAT 11-10-2026 (tinjauan c3,
--    diverifikasi 0f): foreign key anak yang TIDAK berindeks hanya TIGA - T_PROPERTYITEMLIST.PARENT_ID dan T_LISTCAUSEOFLOSS.PARENT_ID
--    (anak T_PROPERTY) serta T_FEALIST.PARENT_ID (anak T_LOCATIONLIST); yang lain berindeks lewat UNIQUE(PARENT_ID) (186-191) atau
--    CREATE INDEX. Tanpa indeks itu, menghapus ±1 juta baris T_PROPERTY memindai PENUH kedua tabel anak itu untuk SETIAP baris
--    (pemindaian penuh membaca sampai high-water mark, yang TIDAK turun oleh DELETE - menghapus anak lebih dulu tidak menolong).
--    Klaim lama di kepala ini ("sebagian besar tabel anak tak berindeks; anak kosong -> pemeriksaan murah") KELIRU dan dicabut.
--    Penjaga -20007 di bawah menolak jalan bila indeks 259 belum ada.
--
-- ⛔ DITULIS SAJA. BELUM PERNAH DIJALANKAN. Dijalankan DBA di DEV sesudah persetujuan work owner. Bukan migrasi bernomor.
-- ⚠️ COMMIT per potongan: TIDAK atomik. Bila berhenti di tengah, sebagian sudah terhapus; blok AMAN dijalankan ulang (setiap
--    pernyataan hanya menghapus yang tersisa; sapuan akhir tanpa rentang menangkap sisa).
-- DIPERTAHANKAN: T_WORK_POLIS, T_NB_OPPORTUNITY, T_GENERAL_POLIS, T_QUOTATIONDATA, T_CEDINGCOLIST (tidak ada DELETE ke sana).
-- Sebelum: ganti SEMUA `{skema}` dengan skema objek nbfacin (env OracleSchema). Sesudah: jalankan BAGIAN 2 kosongkan_nb_151032.sql.
-- Pernyataan hapus = sqlHapusObjek kode repository (sama dengan kosongkan_nb_151032.sql), hanya diberi filter rentang SEQ_NO.
-- =====================================================================================================================
-- SET SERVEROUTPUT ON   (perintah SQL*Plus - hanya untuk SQL*Plus / SQL Developer F5; di aplikasi lain hapus baris ini dan aktifkan panel "DBMS Output")
DECLARE
  v_kasus  CONSTANT VARCHAR2(32) := 'NB-151032';
  v_potong CONSTANT NUMBER := 100000;
  v_n      NUMBER;
  v_maks   NUMBER;
  v_dari   NUMBER;
  v_sampai NUMBER;
  v_total  NUMBER;

  TYPE t_sql IS TABLE OF VARCHAR2(4000);
  -- Satu tingkat = satu pernyataan; bind :1 kasus, :2 SEQ_NO dari, :3 SEQ_NO sampai. Urutan = anak sebelum induk.
  v_tingkat t_sql := t_sql(
    -- cucu Property
    'DELETE FROM {skema}.T_COINSDATA WHERE PARENT_ID IN (SELECT c.ID FROM {skema}.T_LISTCAUSEOFLOSS c JOIN {skema}.T_PROPERTY p ON p.ID = c.PARENT_ID JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_TABLEOFLIMIT WHERE PARENT_ID IN (SELECT o.ID FROM {skema}.T_OCCUPATIONLIST o JOIN {skema}.T_PROPERTY p ON p.ID = o.PARENT_ID JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE o.PARENT_TABLE = ''T_PROPERTY'' AND l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_DEDUCTIBLELIST WHERE PARENT_ID IN (SELECT v.ID FROM {skema}.T_COVERAGELIST v JOIN {skema}.T_PROPERTYITEMLIST i ON i.ID = v.PARENT_ID JOIN {skema}.T_PROPERTY p ON p.ID = i.PARENT_ID JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE v.PARENT_TABLE = ''T_PROPERTYITEMLIST'' AND l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_SPREADINGLIST WHERE PARENT_ID IN (SELECT v.ID FROM {skema}.T_COVERAGELIST v JOIN {skema}.T_PROPERTYITEMLIST i ON i.ID = v.PARENT_ID JOIN {skema}.T_PROPERTY p ON p.ID = i.PARENT_ID JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE v.PARENT_TABLE = ''T_PROPERTYITEMLIST'' AND l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_COVERAGELIST WHERE PARENT_TABLE = ''T_PROPERTYITEMLIST'' AND PARENT_ID IN (SELECT i.ID FROM {skema}.T_PROPERTYITEMLIST i JOIN {skema}.T_PROPERTY p ON p.ID = i.PARENT_ID JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    -- anak Property (cucunya sudah habis)
    'DELETE FROM {skema}.T_LISTCAUSEOFLOSS WHERE PARENT_ID IN (SELECT p.ID FROM {skema}.T_PROPERTY p JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_OCCUPATIONLIST WHERE PARENT_TABLE = ''T_PROPERTY'' AND PARENT_ID IN (SELECT p.ID FROM {skema}.T_PROPERTY p JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_PROPERTYITEMLIST WHERE PARENT_ID IN (SELECT p.ID FROM {skema}.T_PROPERTY p JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_RISKLOCATION WHERE PARENT_ID IN (SELECT p.ID FROM {skema}.T_PROPERTY p JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_BUILDINGCONSTRUCTION WHERE PARENT_ID IN (SELECT p.ID FROM {skema}.T_PROPERTY p JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_SURROUNDINGRISK WHERE PARENT_ID IN (SELECT p.ID FROM {skema}.T_PROPERTY p JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    -- anak Location
    'DELETE FROM {skema}.T_PROPERTY WHERE PARENT_ID IN (SELECT l.ID FROM {skema}.T_LOCATIONLIST l WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    'DELETE FROM {skema}.T_FEALIST WHERE PARENT_ID IN (SELECT l.ID FROM {skema}.T_LOCATIONLIST l WHERE l.PARENT_ID = :1 AND l.SEQ_NO BETWEEN :2 AND :3)',
    -- Location
    'DELETE FROM {skema}.T_LOCATIONLIST WHERE PARENT_ID = :1 AND SEQ_NO BETWEEN :2 AND :3'
  );
BEGIN
  -- Penjaga
  SELECT COUNT(*) INTO v_n FROM {skema}.T_WORK_POLIS WHERE ID = v_kasus AND LINI = 'FAC';
  IF v_n <> 1 THEN RAISE_APPLICATION_ERROR(-20001, 'NB-151032 tidak ada di T_WORK_POLIS LINI FAC'); END IF;
  SELECT COUNT(*) INTO v_n FROM {skema}.T_NB_OPPORTUNITY WHERE ID = v_kasus AND UPPER(TRIM(CLASS_OF_BUSINESS)) = 'MARINE CARGO';
  IF v_n <> 1 THEN RAISE_APPLICATION_ERROR(-20002, 'CLASS_OF_BUSINESS NB-151032 bukan MARINE CARGO - periksa BAGIAN 0.1 kosongkan_nb_151032.sql'); END IF;
  SELECT COUNT(*) INTO v_n FROM {skema}.T_NB_CSV_UNGGAH WHERE CASE_ID = v_kasus AND STATUS IN ('memvalidasi', 'menyimpan');
  IF v_n <> 0 THEN RAISE_APPLICATION_ERROR(-20003, 'ada pekerjaan unggah CSV berjalan di NB-151032'); END IF;
  -- Prasyarat kinerja: indeks FK migrasi 259 (lihat kepala berkas).
  SELECT COUNT(DISTINCT c.TABLE_NAME) INTO v_n FROM ALL_IND_COLUMNS c
   WHERE c.TABLE_OWNER = UPPER('{skema}') AND c.COLUMN_NAME = 'PARENT_ID' AND c.COLUMN_POSITION = 1
     AND c.TABLE_NAME IN ('T_PROPERTYITEMLIST', 'T_LISTCAUSEOFLOSS', 'T_FEALIST');
  IF v_n <> 3 THEN RAISE_APPLICATION_ERROR(-20007, 'indeks PARENT_ID migrasi 259 belum lengkap (' || v_n || '/3) - jalankan migrasi 259 dulu'); END IF;
  -- Versi ini hanya untuk keadaan DATA NB.csv: pohon Marine Cargo dan Aneka KOSONG. Bila tidak, pakai kosongkan_nb_151032.sql.
  SELECT COUNT(*) INTO v_n FROM {skema}.T_CARGOLIST WHERE PARENT_ID = v_kasus;
  IF v_n <> 0 THEN RAISE_APPLICATION_ERROR(-20004, 'T_CARGOLIST NB-151032 tidak kosong - pakai kosongkan_nb_151032.sql'); END IF;
  SELECT COUNT(*) INTO v_n FROM {skema}.T_OCCUPATIONLIST o WHERE o.PARENT_TABLE = 'T_RISKLOCATION'
     AND o.PARENT_ID IN (SELECT r.ID FROM {skema}.T_RISKLOCATION r JOIN {skema}.T_PROPERTY p ON p.ID = r.PARENT_ID
                          JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = v_kasus);
  IF v_n <> 0 THEN RAISE_APPLICATION_ERROR(-20005, 'okupasi Aneka NB-151032 tidak kosong - pakai kosongkan_nb_151032.sql'); END IF;
  SELECT COUNT(*) INTO v_n FROM {skema}.T_ANEKALIST a WHERE a.PARENT_TABLE = 'T_RISKLOCATION'
     AND a.PARENT_ID IN (SELECT r.ID FROM {skema}.T_RISKLOCATION r JOIN {skema}.T_PROPERTY p ON p.ID = r.PARENT_ID
                          JOIN {skema}.T_LOCATIONLIST l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = v_kasus);
  IF v_n <> 0 THEN RAISE_APPLICATION_ERROR(-20006, 'objek Golf NB-151032 tidak kosong - pakai kosongkan_nb_151032.sql'); END IF;

  SELECT NVL(MAX(SEQ_NO), 0) INTO v_maks FROM {skema}.T_LOCATIONLIST WHERE PARENT_ID = v_kasus;
  DBMS_OUTPUT.PUT_LINE('SEQ_NO lokasi maks = ' || v_maks);

  -- Per tingkat, untuk SELURUH case, dipotong per v_potong lokasi; lalu sapuan tanpa rentang (SEQ_NO NULL / di luar rentang).
  FOR t IN 1 .. v_tingkat.COUNT LOOP
    v_total := 0;
    v_dari := 1;
    WHILE v_dari <= v_maks LOOP
      v_sampai := v_dari + v_potong - 1;
      EXECUTE IMMEDIATE v_tingkat(t) USING v_kasus, v_dari, v_sampai;
      v_total := v_total + SQL%ROWCOUNT;
      COMMIT;
      v_dari := v_sampai + 1;
    END LOOP;
    EXECUTE IMMEDIATE REPLACE(REPLACE(v_tingkat(t), 'l.SEQ_NO BETWEEN :2 AND :3', '1 = 1'), 'SEQ_NO BETWEEN :2 AND :3', '1 = 1') USING v_kasus;
    v_total := v_total + SQL%ROWCOUNT;
    COMMIT;
    DBMS_OUTPUT.PUT_LINE('tingkat ' || t || ': ' || v_total || ' baris - ' || SUBSTR(v_tingkat(t), 1, INSTR(v_tingkat(t), ' WHERE') - 1));
  END LOOP;

  -- Tabel tingkat case (DATA NB.csv: semua 0; tetap dijalankan supaya case benar-benar kosong)
  DELETE FROM {skema}.T_LISTINSTALLMENT WHERE PARENT_ID IN (SELECT ID FROM {skema}.T_CURRENCYLIST WHERE PARENT_ID = v_kasus);
  DELETE FROM {skema}.T_CURRENCYLIST WHERE PARENT_ID = v_kasus;
  DELETE FROM {skema}.T_CEDINGCEDANTLIST WHERE PARENT_ID = v_kasus;
  DELETE FROM {skema}.T_CLAUSELIST WHERE PARENT_ID = v_kasus;
  DELETE FROM {skema}.T_SCORING_OPTION WHERE PARENT_ID IN (SELECT f.ID FROM {skema}.T_SCORING_FACTOR f WHERE f.PARENT_ID IN (SELECT d.ID FROM {skema}.T_DATASCORINGRISKLIST d WHERE d.PARENT_ID IN (SELECT s.ID FROM {skema}.T_SCORINGRISK s WHERE s.PARENT_ID = v_kasus)));
  DELETE FROM {skema}.T_SCORING_FACTOR WHERE PARENT_ID IN (SELECT d.ID FROM {skema}.T_DATASCORINGRISKLIST d WHERE d.PARENT_ID IN (SELECT s.ID FROM {skema}.T_SCORINGRISK s WHERE s.PARENT_ID = v_kasus));
  DELETE FROM {skema}.T_DATASCORINGRISKLIST WHERE PARENT_ID IN (SELECT s.ID FROM {skema}.T_SCORINGRISK s WHERE s.PARENT_ID = v_kasus);
  DELETE FROM {skema}.T_SCORINGRESULT WHERE PARENT_ID IN (SELECT s.ID FROM {skema}.T_SCORINGRISK s WHERE s.PARENT_ID = v_kasus);
  DELETE FROM {skema}.T_SCORINGRISK WHERE PARENT_ID = v_kasus;
  DELETE FROM {skema}.T_LOSSRECORDLIST WHERE PARENT_ID IN (SELECT k.ID FROM {skema}.T_LOSSRECORD k WHERE k.PARENT_ID = v_kasus);
  DELETE FROM {skema}.T_LOSSRECORD WHERE PARENT_ID = v_kasus;
  DELETE FROM {skema}.T_RECIPIENTLIST WHERE PARENT_ID IN (SELECT k.ID FROM {skema}.T_CORRESPONDENCELIST k WHERE k.PARENT_ID = v_kasus);
  DELETE FROM {skema}.T_CORRESPONDENCELIST WHERE PARENT_ID = v_kasus;
  COMMIT;

  -- Penampung unggah CSV: baris per rentang NO_BARIS (PK (UNGGAH_ID, NO_BARIS)), lalu pesan + pekerjaan
  SELECT NVL(MAX(b.NO_BARIS), 0) INTO v_maks FROM {skema}.T_NB_CSV_BARIS b
   WHERE b.UNGGAH_ID IN (SELECT u.ID FROM {skema}.T_NB_CSV_UNGGAH u WHERE u.CASE_ID = v_kasus);
  v_dari := 1;
  WHILE v_dari <= v_maks LOOP
    v_sampai := v_dari + v_potong - 1;
    DELETE FROM {skema}.T_NB_CSV_BARIS WHERE NO_BARIS BETWEEN v_dari AND v_sampai
       AND UNGGAH_ID IN (SELECT u.ID FROM {skema}.T_NB_CSV_UNGGAH u WHERE u.CASE_ID = v_kasus);
    COMMIT;
    v_dari := v_sampai + 1;
  END LOOP;
  DELETE FROM {skema}.T_NB_CSV_BARIS WHERE UNGGAH_ID IN (SELECT u.ID FROM {skema}.T_NB_CSV_UNGGAH u WHERE u.CASE_ID = v_kasus);
  DELETE FROM {skema}.T_NB_CSV_PESAN WHERE UNGGAH_ID IN (SELECT u.ID FROM {skema}.T_NB_CSV_UNGGAH u WHERE u.CASE_ID = v_kasus);
  DELETE FROM {skema}.T_NB_CSV_UNGGAH WHERE CASE_ID = v_kasus;
  COMMIT;
  DBMS_OUTPUT.PUT_LINE('SELESAI (sudah COMMIT per potongan) - jalankan BAGIAN 2 kosongkan_nb_151032.sql');
END;
/
